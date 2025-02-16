package storage

import (
	"avito-shop/pkg/storage/models"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type uuid string

var elog = log.New(os.Stderr, "[Postgresql error]\t", log.Ldate|log.Ltime|log.Lshortfile)
var ilog = log.New(os.Stdout, "[Postgresql info]\t", log.Ldate|log.Ltime)

type Query struct {
	Command string
	Param   string
}

type Crud interface {
	Create(d any) (uuid, error)
	Get(q Query) (any, error)
	Update(u interface{}) error
	Delete(id uuid) error
}

type Storage struct {
	User  Crud
	Merch map[string]models.Merch
}

var storage *Storage

func NewStorage(db *pgxpool.Pool) (*Storage, error) {
	if db == nil {
		err := errors.New("db is nil")
		elog.Println(err)
		return nil, err
	}
	m, err := getMerchList(db)
	if err != nil {
		return nil, err
	}
	storage = &Storage{
		User:  newUserCrud(db),
		Merch: m,
	}
	return storage, nil
}

func (s Storage) NewQuery(c, p string) Query {
	return Query{
		Command: c,
		Param:   p,
	}
}
