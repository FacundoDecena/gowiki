package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Structs section

type Page struct {
	Title string
	Body  string//[]byte
	SBody []string
}

var templates = template.Must(template.ParseFiles("./templates/edit.html", "./templates/view.html"))

var validPath = regexp.MustCompile("^/(edit|view|save)/([a-zA-Z0-9]+)$")

// Persistence section

func (p *Page) save() error {
	filename := "./data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, []byte(p.Body), 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "./data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: string(body)}, nil
}

// Handlers section

func viewHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, err := loadPage(title)

	page.SBody = strings.Split(page.Body, "\n")

	if err != nil {
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound) // 302
		return
	}

	renderTemplate(writer, "view", page)
}

func editHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, err := loadPage(title)

	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request, title string) {
	body := request.FormValue("body")

	page := &Page{Title: title, Body: body}

	err := page.save()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(writer, request, "/view/"+title, http.StatusFound)

}

func homeHandler(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/view/FrontPage", http.StatusFound)
}

func makeHandler(function func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		match := validPath.FindStringSubmatch(request.URL.Path)
		if match == nil {
			http.NotFound(writer, request)
			return
		}

		function(writer, request, match[2])
	}
}

// Render section

func renderTemplate(writer http.ResponseWriter, filename string, page *Page) {

	err := templates.ExecuteTemplate(writer, filename+".html", page)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

}

// Main section

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	var port string
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":" + port, nil))

}
