package handler

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type user struct {
	Username string `form:"username" json:"username" binding:"required"`
	UserID   string `form:"userid" json:"userid" binding:"required"`
}

func (h *Handler) newAuthMiddleWare() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte("pulsarjwtsupersecretkeynah"),
		Timeout:         time.Hour * 24 * 7,
		MaxRefresh:      time.Hour * 24 * 7,
		Authenticator:   h.authenticator, //при логине
		PayloadFunc:     payloadFunc,     //при логине
		LoginResponse:   loginResponse,   //при логине
		IdentityHandler: identityHandler, //при запросе
		Authorizator:    authorizator,    //при запросе
		Unauthorized:    unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		//SendCookie:      true,
	})
	if err != nil {
		return nil, err
	}
	return authMiddleware, nil
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*user); ok {
		return jwt.MapClaims{
			"username": v.Username,
			"userid":   v.UserID,
		}
	}
	return jwt.MapClaims{}
}

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	u := &user{
		Username: claims["username"].(string),
		UserID:   claims["userid"].(string),
	}
	c.Set("username", u.Username)
	c.Set("userid", u.UserID)
	return u
}

func authorizator(data interface{}, c *gin.Context) bool {
	if _, ok := data.(*user); ok {
		return true
	}

	return false
}
