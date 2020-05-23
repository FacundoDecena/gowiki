// Package models collects the models used in the site
package models

import (
	"html/template"
	"io/ioutil"
	"regexp"
)

// Page contains all the information of the wiki
type Page struct {
	Title string
	Body  string
	SBody []string
}

// Templates are the available templates to render
var Templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))

// regexp for the valid paths
var ValidPath = regexp.MustCompile("^/(edit|view|save)/([a-zA-Z0-9]+)$")

// Save Saves a Page into bd
func (p *Page) Save() error {
	filename := "./data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, []byte(p.Body), 0600)
}

// LoadPage loads a Page from the bd if existm else return error
func LoadPage(title string) (*Page, error) {
	filename := "./data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: string(body)}, nil
}
