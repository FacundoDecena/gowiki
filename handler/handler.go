// Package handler implements different handlers for http requests.
package handler

import (
	"log"
	"net/http"
	"strings"

	models "github.com/FacundoDecena/gowiki/models"
)

// ViewHandler shows a page
func ViewHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, err := models.LoadPage(title)

	if err != nil {
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound) // 302
		return
	}

	page.SBody = strings.Split(page.Body, "\n")

	renderTemplate(writer, "view", page)
}

// EditHandler let the user edit a page
func EditHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, err := models.LoadPage(title)

	if err != nil {
		page = &models.Page{Title: title}
	}

	renderTemplate(writer, "edit", page)
}

// SaveHandler saves the edited page
func SaveHandler(writer http.ResponseWriter, request *http.Request, title string) {
	body := request.FormValue("body")

	page := &models.Page{Title: title, Body: body}

	err := page.Save()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(writer, request, "/view/"+title, http.StatusFound)

}

// HomeHandler redirects to /view/FrontPage when someone try to connect to /
func HomeHandler(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/view/FrontPage", http.StatusFound)
}

// Gets the title and call the function function with it
func MakeHandler(function func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		err := request.ParseForm()
		if err != nil {
			log.Fatal("Efeo el ParseFrom")
		}

		search := request.Form.Get("search")

		if search != "" {
			match := models.ValidPath.FindStringSubmatch("/view/" + search)
			function(writer, request, match[2])

		} else {
			match := models.ValidPath.FindStringSubmatch(request.URL.Path)
			if match == nil {
				http.NotFound(writer, request)
				return
			}

			function(writer, request, match[2])
		}
	}
}

func renderTemplate(writer http.ResponseWriter, filename string, page *models.Page) {

	err := models.Templates.ExecuteTemplate(writer, filename+".html", page)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

}
