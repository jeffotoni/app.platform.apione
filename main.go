package main

import "net/http"

func main() {

	http.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello man!!!"))
	})

	println("Run Server:8080")
	http.ListenAndServe(":8080", nil)

}
