package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"database/sql"

	"webpush/pushmess"

	_ "github.com/lib/pq"
)

func main() {
	pushmess.Token = flag.String("token", "areaynwu2roqijy8na5hs14gmwytp5", "Application token")
	pushmess.User = flag.String("user", "ucryge6j8mr9jnyhkef5jkab71y7sm", "User key")

	flag.Parse()

	connStr := "user=postgres password=12481 dbname=pushover sslmode=disable"

	var err error
	pushmess.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer pushmess.DB.Close()

	r := mux.NewRouter()
	r.HandleFunc("/getcount/{from}", pushmess.GetCount).Methods("GET")
	r.HandleFunc("/send", pushmess.Send).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
