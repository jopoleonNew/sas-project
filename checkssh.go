package main

import (
	"net/http"

	"fmt"

	//"log"

	"github.com/gorilla/mux"
)

var port = ":8080"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler) // GET
	log.Println("Server started at port ", port)
	err := http.ListenAndServe(port, r)

	if err != nil {
		log.Fatalln(err)
	}
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello ", r.RemoteAddr)
}
