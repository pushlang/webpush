package webpush

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"database/sql"

	_ "github.com/lib/pq"
)

var webServer = ":8000" //http://127.0.0.1:8000/getcount/2020-10-24T18:50:23.541Z
var tokenDef = "areaynwu2roqijy8na5hs14gmwytp5"
var userDef = "ucryge6j8mr9jnyhkef5jkab71y7sm"

func Run() int {
	log.Println("Запуск сервиса Webpush (" + webServer + ")...")
	
	Token = flag.String("token", tokenDef, "Application token")
	User = flag.String("user", userDef, "User key")
	//ucryge6j8mr9jnyhkef5jkab71y7sm chrome win, firefox ubu

	flag.Parse()

	if *Token == tokenDef {
		log.Println("Предупреждение: Используется токен приложения по умолчанию")
	}

	if *User == "" || *Token == "" {
		log.Fatal("Ошибка: Токен приложения или ключ пользователя не заданы")
		return -1
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
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка: Подключение к БД невозможно, код ошибки - ", err)
	}

	log.Println("Подключение к БД...")

	defer func() {
		DB.Close()
		if err != nil {
			log.Fatal("Ошибка: Невозможно закрыть БД, код ошибки - ", err)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/getcount/{from}", GetCount).Methods("GET")
	r.HandleFunc("/send", Send).Methods("POST", "OPTIONS")

	log.Println("Запуск web-сервера...")
	err = http.ListenAndServe(webServer, r)
	if err != nil {
		log.Fatal("Ошибка: Запуск web-сервера не выполнен, код ошибки - ", err)
	}
	
	return 0
}
