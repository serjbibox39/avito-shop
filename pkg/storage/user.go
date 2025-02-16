package storage

import (
	"context"
	"errors"

	"avito-shop/pkg/storage/models"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Объект, реализующий интерфейс работы с таблицей users PostgreSQL.
type UserPostgres struct {
	db *pgxpool.Pool
	do map[string]func(*UserPostgres, Query) (any, error)
}

// Конструктор UserPostgres
func newUserCrud(db *pgxpool.Pool) Crud {
	do := make(map[string]func(*UserPostgres, Query) (any, error))
	do["/checkpwd"] = getUserForCheckPwd
	do["/info"] = getInfo
	return &UserPostgres{
		db: db,
		do: do,
	}
}

func (up *UserPostgres) Create(d any) (uuid, error) {
	u := d.(models.User)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	id := ""
	q := `WITH i AS (
		INSERT INTO users(username, password) VALUES ($1, $2) RETURNING ID
		)
		INSERT INTO inventory (userid) SELECT i.id FROM i RETURNING ID;`
	err = up.db.QueryRow(context.Background(), q, u.Username, hashedPassword).Scan(&id)
	return uuid(id), err
}

func (up *UserPostgres) Get(query Query) (any, error) {
	_, ok := up.do[query.Command]
	if !ok {
		return nil, errors.New("астрологи объявили неделю невероятных событий, количество необъяснимого удваивается")
	}
	u, err := up.do[query.Command](up, query)
	return u, err
}

func getInfo(sp *UserPostgres, query Query) (any, error) {
	ctx := context.Background()
	tx, err := sp.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	q := `SELECT unnest(merch) AS type, 
			COUNT(*) AS quantity
			FROM inventory
			WHERE userid = $1
			GROUP BY type
			ORDER BY type;`
	rows, err := tx.Query(ctx, q, query.Param)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	inventory := []models.InventoryOut{}
	for rows.Next() {
		i := models.InventoryOut{}
		err = rows.Scan(
			&i.Name,
			&i.Quantity,
		)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		inventory = append(inventory, i)
	}

	q = `SELECT coins FROM inventory WHERE userid=$1;`
	rows, err = tx.Query(ctx, q, query.Param)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	coins := 0
	for rows.Next() {
		err = rows.Scan(
			&coins,
		)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
	}

	q = `SELECT transactions.to_user, amount, transactions.created_at
		FROM transactions
		WHERE transactions.from_user = $1
		ORDER BY transactions.created_at DESC;
	`
	rows, err = tx.Query(ctx, q, query.Param)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	sent := []models.ToUser{}
	for rows.Next() {
		s := models.ToUser{}
		err = rows.Scan(
			&s.ToUser,
			&s.Amount,
			&s.CreatedAt,
		)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		sent = append(sent, s)
	}
	q = `SELECT transactions.from_user, amount, transactions.created_at
		FROM transactions
		WHERE transactions.to_user = $1
		ORDER BY transactions.created_at DESC;
	`
	rows, err = tx.Query(ctx, q, query.Param)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	rec := []models.FromUser{}
	for rows.Next() {
		r := models.FromUser{}
		err = rows.Scan(
			&r.FromUser,
			&r.Amount,
			&r.CreatedAt,
		)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		rec = append(rec, r)
	}
	ch := models.CoinHistory{
		Received: rec,
		Sent:     sent,
	}
	info := models.Info{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: ch,
	}
	return info, tx.Commit(ctx)
}

func getUserForCheckPwd(sp *UserPostgres, query Query) (any, error) {
	var u models.User
	q := `SELECT id, password FROM users WHERE username=$1`
	rows, err := sp.db.Query(context.Background(), q, query.Param)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&u.ID,
			&u.Password,
		)
		if err != nil {
			return nil, err
		}
	}
	u.Username = query.Param
	return u, rows.Err()
}

func (up *UserPostgres) Update(p any) error {
	item, ok := p.(models.BuyItem)
	if ok {
		err := buyMerch(up, item)
		if err != nil {
			return err
		}
		return nil
	}
	t, ok := p.(models.Transaction)
	if ok {
		err := txCoin(up, t)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("тут должна была быть транзакция, но что-то пошло не так")
}

func txCoin(up *UserPostgres, t models.Transaction) error {
	ctx := context.Background()
	tx, err := up.db.Begin(ctx)
	if err != nil {
		return err
	}
	//Проверяем id принимающего
	toUser := ""
	q := `SELECT id FROM users WHERE username = $1`
	rows, err := tx.Query(ctx, q, t.ToUser)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	for rows.Next() {
		err = rows.Scan(
			&toUser,
		)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}
	// Списание монет
	q = `UPDATE inventory SET coins = coins - $1 WHERE userid = $2`
	_, err = tx.Exec(ctx, q, t.Amount, t.FromUser)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	// Зачисление монет
	q = `UPDATE inventory SET coins = coins + $1 WHERE userid = $2`
	_, err = tx.Exec(ctx, q, t.Amount, toUser)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	id := 0
	q = `INSERT INTO transactions (from_user, to_user, amount) 
			VALUES ($1, $2, $3) RETURNING ID;`
	err = tx.QueryRow(ctx, q, t.FromUser, toUser, t.Amount).Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	ilog.Println("провели транзакцию", t.FromUser, toUser, t.Amount, id)
	return tx.Commit(ctx)
}

func buyMerch(up *UserPostgres, item models.BuyItem) error {
	ctx := context.Background()
	tx, err := up.db.Begin(ctx)
	if err != nil {
		return err
	}
	// Списание монет
	q := `UPDATE inventory SET coins = coins - $1 WHERE userid = $2`
	_, err = tx.Exec(ctx, q, storage.Merch[item.Item].Price, item.UserID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	// Приход мерча
	q = `UPDATE inventory SET merch = array_append(merch,$1) WHERE userid = $2`
	_, err = tx.Query(ctx, q, item.Item, item.UserID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (up *UserPostgres) Delete(id uuid) error {
	return nil
}
