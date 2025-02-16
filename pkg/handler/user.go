package handler

import (
	"avito-shop/pkg/storage/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transaction struct {
	Destination string `json:"toUser"`
	Amount      int    `json:"amount"`
}

// Получение информации о монетах, инвентаре и истории транзакций
func (h *Handler) getInfo(c *gin.Context) {
	id, _ := c.Get("userid")
	info, err := h.storage.User.Get(h.storage.NewQuery("/info", id.(string)))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, info)
}

// Покупка предмета за монеты
func (h *Handler) buyItem(c *gin.Context) {
	id, _ := c.Get("userid")
	if _, ok := h.storage.Merch[c.Param("item")]; !ok {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("этот товар отсутствует в продаже").Error())
	}
	b := models.BuyItem{
		UserID: id.(string),
		Item:   c.Param("item"),
	}
	err := h.storage.User.Update(b)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// Транзакция монеток
func (h *Handler) transaction(c *gin.Context) {
	username, _ := c.Get("username")
	userid, _ := c.Get("userid")
	t := transaction{}
	err := c.ShouldBind(&t)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("bind error: %w", err).Error())
		return
	}
	if username == t.Destination {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("не получится отправить монетки самому себе").Error())
		return
	}
	if t.Amount <= 0 {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("перевод должен быть больше 0").Error())
		return
	}
	transaction := models.Transaction{
		FromUser: userid.(string),
		ToUser:   t.Destination,
		Amount:   t.Amount,
	}
	err = h.storage.User.Update(transaction)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}
