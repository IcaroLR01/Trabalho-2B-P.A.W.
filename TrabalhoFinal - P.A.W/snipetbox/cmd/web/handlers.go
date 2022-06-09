package main

//go run cmd/web/*
import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
  "encoding/json"

	"github.com/rmcs87/cc5m/pkg/models"
)

type Item struct{
  Nome string
  Contato string
  Saida string
}

func (app *application) home(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(rw)
		return
	}

  snippets, err := app.snippets.Latest()
  if err != nil{
    app.serverError(rw, err)
    return
  }

	files := []string{
		"./ui/html/home.page.tmpl.html",
		"./ui/html/base.layout.tmpl.html",
		"./ui/html/footer.partial.tmpl.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(rw, err)
		return
	}
	err = ts.Execute(rw, snippets)
	if err != nil {
		app.serverError(rw, err)
		return
	}

}

//http://localhost:4000/snippet?id=1
func (app *application) showSnippet(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(rw)
		return
	}

  s, err := app.snippets.Get(id)
  if err == models.ErrNoRecord {
    app.notFound(rw)
    return
  }else if err != nil{
    app.serverError(rw, err)
    return
  }
  
  files := []string{
		"./ui/html/show.page.tmpl.html",
		"./ui/html/base.layout.tmpl.html",
		"./ui/html/footer.partial.tmpl.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(rw, err)
		return
	}
	err = ts.Execute(rw, s)
	if err != nil {
		app.serverError(rw, err)
		return
	}
  
}

func (app *application) createSnippet(rw http.ResponseWriter, r *http.Request) {

	var item Item
    err := json.NewDecoder(r.Body).Decode(&item)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

	id, err := app.snippets.Insert(item.Nome, item.Contato, item.Saida)
	if err != nil {
		app.serverError(rw, err)
		return
	}

	http.Redirect(rw, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *application) deleteSnippet(rw http.ResponseWriter, r *http.Request) {

  id, err := strconv.Atoi(r.URL.Query().Get("id"))

	id, err = app.snippets.Delete(id)
	if err != nil {
		app.serverError(rw, err)
		return
	}

	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (app *application) updateSnippet(rw http.ResponseWriter, r *http.Request) {

  id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var item Item
    err := json.NewDecoder(r.Body).Decode(&item)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

	id, err = app.snippets.Update(id, item.Nome, item.Contato)
	if err != nil {
		app.serverError(rw, err)
		return
	}

	http.Redirect(rw, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
