package handler

import (
	"avito-shop/pkg/storage/models"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) authenticator(c *gin.Context) (interface{}, error) {
	var l login
	if err := c.ShouldBind(&l); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	if (l.Username == "") || (l.Password == "") {
		return "", jwt.ErrMissingLoginValues
	}
	u, err := h.storage.User.Get(h.storage.NewQuery("/checkpwd", l.Username))
	if err != nil {
		return nil, err
	}
	storedUser := u.(models.User)
	if storedUser.ID == "" {
		storedUser.Username = l.Username
		storedUser.Password = l.Password
		_, err := h.storage.User.Create(storedUser)
		if err != nil {
			return nil, err
		}

	}
	u, err = h.storage.User.Get(h.storage.NewQuery("/checkpwd", l.Username))
	if err != nil {
		return nil, err
	}
	storedUser = u.(models.User)
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(l.Password))
	if err != nil {
		return nil, err
	}
	claim := &user{
		Username: storedUser.Username,
		UserID:   storedUser.ID,
	}
	return claim, nil
}
