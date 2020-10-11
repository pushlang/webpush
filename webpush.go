package webpush

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"database/sql"

	"webpush/pushmess"

	_ "github.com/lib/pq"
)

var webServer = ":8000" //http://127.0.0.1:8000/getcount/2020-10-24T18:50:23.541Z
var tokenDef = "areaynwu2roqijy8na5hs14gmwytp5"
var userDef = ""

func run() {
	log.Println("Запуск сервиса Webpush (" + webServer + ")...")
	pushmess.Token = flag.String("token", tokenDef, "Application token")
	pushmess.User = flag.String("user", userDef, "User key")
	//ucryge6j8mr9jnyhkef5jkab71y7sm chrome win, firefox ubu

	flag.Parse()

	if *pushmess.Token == tokenDef {
		log.Println("Предупреждение: Используется токен приложения по умолчанию")
	}

	if *pushmess.User == "" || *pushmess.Token == "" {
		log.Fatal("Ошибка: Токен приложения или ключ пользователя не заданы")
	}

	//export DB_DSN="user=postgres password=password dbname=pushover sslmode=disable"
	//set DB_DSN="user=postgres password=password dbname=pushover sslmode=disable"
	connStr := os.Getenv("DB_DSN")
	if connStr == "" {
		log.Fatal("Ошибка: Переменная среды DB_DSN не задана")
	}

	//CREATE DATABASE pushover;
	//\c pushover;
	//CREATE TABLE messages(id integer PRIMARY KEY, token text, userr text, textm text, status integer, sent text);
	var err error
	connStr = "user=postgres password=12481 dbname=pushover sslmode=disable"
	pushmess.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка: Подключение к БД невозможно, код ошибки - ", err)
	}

	log.Println("Подключение к БД...")

	defer func() {
		pushmess.DB.Close()
		if err != nil {
			log.Fatal("Ошибка: Невозможно закрыть БД, код ошибки - ", err)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/getcount/{from}", pushmess.GetCount).Methods("GET")
	r.HandleFunc("/send", pushmess.Send).Methods("POST", "OPTIONS")

	log.Println("Запуск web-сервера...")
	err = http.ListenAndServe(webServer, r)
	if err != nil {
		log.Fatal("Ошибка: Запуск web-сервера не выполнен, код ошибки - ", err)
	}

}
