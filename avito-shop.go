package main

import (
	"encoding/json"
	"log"
	"os"

	"avito-shop/pkg/handler"
	"avito-shop/pkg/server"
	"avito-shop/pkg/storage"
	"avito-shop/pkg/storage/postgresql"
)

var elog = log.New(os.Stderr, "[Service error]\t", log.Ldate|log.Ltime|log.Lshortfile)
var ilog = log.New(os.Stdout, "[Service info]\t", log.Ldate|log.Ltime)

// Конфигурация приложения
type config struct {
	PostgresConfig postgresql.PostgresConfig `json:"postgres_settings"`
}

func main() {
	connString := postgresql.GetConnectionStringEnv()
	db, err := postgresql.New(connString)
	if err != nil {
		elog.Fatal(err)
	}
	s, err := storage.NewStorage(db)
	if err != nil {
		elog.Fatal(err)
	}
	handlers, err := handler.NewHandler(s)
	if err != nil {
		elog.Fatal(err)
	}
	srv := new(server.Server)
	ilog.Println("service is starting")
	elog.Fatal(srv.Run(server.HTTP_PORT, handlers.InitRoutes()))
}

// Чтение JSON файла конфигурации
func readConfig(path string) (*config, error) {
	c, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config config
	err = json.Unmarshal(c, &config)
	return &config, err
}
