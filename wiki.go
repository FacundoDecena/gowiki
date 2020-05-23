package main

import (
	"log"
	"net/http"
	"os"

	handler "github.com/FacundoDecena/gowiki/handler"
)

func main() {

	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/view/", handler.MakeHandler(handler.ViewHandler))
	http.HandleFunc("/edit/", handler.MakeHandler(handler.EditHandler))
	http.HandleFunc("/save/", handler.MakeHandler(handler.SaveHandler))

	var port string
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
