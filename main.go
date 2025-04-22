package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

func GetHome(w http.ResponseWriter, r *http.Request) {
	GetTemplate(w, "home.page.html")
}

func GetAbout(w http.ResponseWriter, r *http.Request) {
	GetTemplate(w, "about.page.html")
}

func GetTemplate(w http.ResponseWriter, htmlFile string) {
	tmpCache, err := CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	tmp, ok := tmpCache[htmlFile]
	if !ok {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	_ = tmp.Execute(buf, nil)
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to the browser ", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmpSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			tmpSet, err = tmpSet.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = tmpSet
	}
	return myCache, nil
}

func main() {
	http.HandleFunc("/", GetHome)
	http.HandleFunc("/about", GetAbout)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
