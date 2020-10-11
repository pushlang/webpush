package pushmess

import (
	"encoding/json"
	"log"
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

// GetCount возвращает JSON-структуру с количеством отправленных сообщений с заданного времени
// http://127.0.0.1:8000/getcount/2020-10-24T18:50:23.541Z
func GetCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)
	
	params := mux.Vars(r)

	from, err := time.Parse(time.RFC3339, params["from"])
	if err != nil {
		log.Println("Предупреждение: Парсинг даты не выполнен, код ошибки - ", err)
		return
	}
	log.Println("Получен запрос GetCount:", from)

	rows, err := DB.Query("SELECT * FROM messages")
	if err != nil {
		log.Println("Предупреждение: Запрос к БД не выполнен, код ошибки - ", err)
		return
	}
	defer rows.Close()

	failed, sent := 0, 0
	for rows.Next() {
		message := Message{}
		t := ""
		err := rows.Scan(&message.ID, &message.Token, &message.User, &message.Text, &message.Status, &t)
		if err != nil {
			log.Println("Предупреждение: Некорректно сформирован запрос к БД, код ошибки - ", err)
			return
		}
		message.Sent, err = time.Parse(time.UnixDate, t)
		if err != nil {
			log.Println("Предупреждение: Некорректный формат даты, код ошибки - ", err)
			return
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

	err = json.NewEncoder(w).Encode(CountCodes{Sent: sent, Failed: failed})
	if err != nil {
		log.Println("Предупреждение: Ошибка кодирования JSON-структуры -", err)
	}
}

// Send принимает JSON-структуру с текстом сообщения и перенаправляет его посредством push-сообщения адресату
func Send(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//enableCors(&w)

	var counter int

	DB.QueryRow("SELECT count(*) FROM messages").Scan(&counter)

	var message = Message{ID: 1000 + counter, Token: *Token, User: *User, Text: "", Status: 0, Sent: time.Now()}
	log.Println("Получен запрос Send:", message.ID)

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		log.Println("Предупреждение: Декодирование JSON-структуры не выполнено -", err)
		return
	}

	m := pushover.NewMessage(message.Token, message.User)
	resp, err := m.Push(message.Text)
	if err != nil {
		log.Println("Предупреждение: Отправка push-сообщения не выполнена -", err)
		//return
	}

	message.Status = resp.Status

	_, err = DB.Exec("insert into messages (id, token, userr, textm, status, sent) values ($1, $2, $3, $4, $5, $6)",
		message.ID, message.Token, message.User, message.Text, message.Status, message.Sent.Format(time.UnixDate))
	if err != nil {
		log.Println("Предупреждение: Отправка push-сообщения не выполнена -", err)
		return
	}

	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Println("Предупреждение: Ошибка кодирования JSON-структуры -", err)
	}

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
