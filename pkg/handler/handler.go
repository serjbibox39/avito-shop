package handler

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"avito-shop/pkg/storage"

	"github.com/gin-gonic/gin"
)

type uuid string

var elog = log.New(os.Stderr, "[Handler error]\t", log.Ldate|log.Ltime|log.Lshortfile)
var ilog = log.New(os.Stdout, "[Handler info]\t", log.Ldate|log.Ltime)

// Обработчик HTTP запросов сервера
type Handler struct {
	storage *storage.Storage
	rand    *rand.Rand
}

// Конструктор объекта Handler
func NewHandler(storage *storage.Storage) (*Handler, error) {
	if storage == nil {
		return nil, errors.New("storage is nil")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Handler{
		storage: storage,
		rand:    r,
	}, nil
}

// Инициализация маршрутизатора запросов.
// Регистрация обработчиков запросов
func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	corsHandler := CORSMiddleware()
	r.Use(corsHandler)
	authMiddleware, err := h.newAuthMiddleWare()
	if err != nil {
		elog.Fatalf("ошибка инициализации сервиса аутентификации: %s", err.Error())
	}
	r.POST("/api/auth", authMiddleware.LoginHandler)
	api := r.Group("/api", authMiddleware.MiddlewareFunc())
	{
		api.GET("/info", h.getInfo)
		api.GET("/buy/:item", h.buyItem)
		api.POST("/sendCoin", h.transaction)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})
	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
