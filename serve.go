package main

import (
	"fmt"
	"net/http"
	"time"
	"math/rand"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func write(writer http.ResponseWriter, request *http.Request) {
	database, _ := sql.Open("sqlite3", "./database.sqlite")
    statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS pages (id INTEGER PRIMARY KEY, title TEXT)")
    statement.Exec()

    var title string = randSeq(10)
    statement, _ = database.Prepare("INSERT INTO pages (title) VALUES (?)")
    statement.Exec(title)
    fmt.Fprintf(writer, title)
}

func read(writer http.ResponseWriter, request *http.Request) {
	database, _ := sql.Open("sqlite3", "./database.sqlite")

	rows, _ := database.Query("SELECT * FROM pages WHERE id >= (abs(random()) % (SELECT max(id) FROM pages)) LIMIT 1")
    
    var id int
	var title string
    for rows.Next() {
        rows.Scan(&id, &title)

        fmt.Fprintf(writer, title)
    }

    rows.Close() //good habit to close
}

func sleep(writer http.ResponseWriter, request *http.Request) {
	time.Sleep(time.Second)
	fmt.Fprintf(writer, "sleeped")
}

func main() {
	http.HandleFunc("/sleep", sleep)
	http.HandleFunc("/read", read)
	http.HandleFunc("/write", write)
	http.ListenAndServe(":8080", nil)
}