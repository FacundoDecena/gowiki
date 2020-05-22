package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Structs section

type Page struct {
	Title string
	Body  []byte
}

// Persistence section

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// Handlers section

func viewHandler(resW http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(resW, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

// Main section

func main() {
	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
