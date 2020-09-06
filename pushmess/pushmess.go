package pushmess

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bdenning/go-pushover"
	"github.com/gorilla/mux"

	"database/sql"
	//_ "github.com/lib/pq"
)

// Message содержит информацию о пересылаемом сообщении
type Message struct {
	ID     int
	Token  string
	User   string
	Text   string `json:"text"`
	Status int
	Sent   time.Time
}

// CountCodes содержит коды статусов отправленных сообщений
type CountCodes struct {
	Sent   int `json:"sent"`
	Failed int `json:"failed"`
}

// Token содержит токен приложения
var Token *string

// User содержит ключ пользователя
var User *string

// DB СУБД Postgres
var DB *sql.DB

// GetCount возвращает JSON-структуру с количеством отправленных сообщений с определенного временного периода
func GetCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	//fmt.Println("from1:", params["from"])

	from, _ := time.Parse(time.RFC3339, params["from"])

	//fmt.Println("from2:", from)

	rows, err := DB.Query("SELECT * FROM messages")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	failed, sent := 0, 0
	for rows.Next() {
		message := Message{}
		t := ""
		err := rows.Scan(&message.ID, &message.Token, &message.User, &message.Text, &message.Status, &t)
		message.Sent, _ = time.Parse(time.UnixDate, t)

		fmt.Println("from3:", message.Sent)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if message.Sent.Before(from) {
			switch message.Status {
			case 1:
				sent++
			default:
				failed++
			}
		}
	}

	json.NewEncoder(w).Encode(CountCodes{Sent: sent, Failed: failed})
}

// Send принимает JSON-структуру с текстом сообщения и перенаправляет его посредством push-сообщения адресату
func Send(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var counter int

	DB.QueryRow("SELECT count(*) FROM messages").Scan(&counter)

	var message = Message{ID: 1000 + counter, Token: *Token, User: *User, Text: "", Status: 0, Sent: time.Now()}

	_ = json.NewDecoder(r.Body).Decode(&message)

	m := pushover.NewMessage(message.Token, message.User)
	resp, _ := m.Push(message.Text)

	message.Status = resp.Status

	_, err := DB.Exec("insert into messages (id, token, userr, textm, status, sent) values ($1, $2, $3, $4, $5, $6)",
		message.ID, message.Token, message.User, message.Text, message.Status, message.Sent.Format(time.UnixDate))
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(message)

	fmt.Println("Получено сообщение:", message)
}
