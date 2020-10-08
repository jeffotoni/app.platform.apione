package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jeffotoni/gconcat"
	_ "github.com/lib/pq"
)

var (
	once    sync.Once
	err     error
	dbLocal *sql.DB

	database = os.Getenv("DB_NAME")
	host     = os.Getenv("DB_HOST")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	port     = os.Getenv("DB_PORT")
	ssl      = "require"
	source   = "postgres"

	httpPort = ":8080"
)

type Login struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
	Hora string `json:"hora"`
}

func Connect() *sql.DB {
	once.Do(func() {
		if dbLocal != nil {
			return
		}
		connStr := gconcat.Build("host=", host, " port=", port,
			" user=", user, " password=", password, " dbname=", database,
			" sslmode=", ssl)
		if dbLocal, err = sql.Open(source, connStr); err != nil {
			if err != nil {
				log.Println(err)
			}
			dbLocal = nil
			return
		}
	})
	return dbLocal
}

func Get() (lv []Login, err error) {

	Db := Connect()

	rows, err := Db.Query(`SELECT uuid,name,hora FROM public.login`)
	if err != nil {
		return
	}

	var l Login
	for rows.Next() {
		err = rows.Scan(&l.Uuid, &l.Name, &l.Hora)
		if err != nil {
			continue
		}
		lv = append(lv, l)
	}

	return
}

func main() {

	http.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello man!!!"))
	})

	http.HandleFunc("/api/v1/user", func(w http.ResponseWriter, r *http.Request) {

		ol, err := Get()
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			errors := gconcat.Build(`{"msg":"Error:`, err.Error(), ` "}`)
			w.Write([]byte(errors))
			return
		}

		b, err := json.Marshal(ol)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errors := gconcat.Build(`{"msg":"Error:`, err.Error(), ` "}`)
			w.Write([]byte(errors))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	println("Run Server v0.0.1", httpPort)
	http.ListenAndServe(httpPort, nil)

}

// package main
// import "net/http"
// func main() {

// 	http.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("hello man!!!"))
// 	})

// 	println("Run Server:8080")
// 	http.ListenAndServe(":8080", nil)

// }
