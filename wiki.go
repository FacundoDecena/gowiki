package main

import (
	"html/template"
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

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/view/"):]
	page, _ := loadPage(title)

	renderTemplate(writer, "view", page)
}

func editHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(writer, "edit", page)
}

// Render section

func renderTemplate(writer http.ResponseWriter, filename string, page *Page) {
	templ, _ := template.ParseFiles(filename + ".html")

	templ.Execute(writer, page)
}

// Main section

func main() {

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	// http.HandleFunc("/save/", saveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
