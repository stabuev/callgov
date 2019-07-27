package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type GServer struct {
	tmpl *template.Template
	db   *sql.DB
}

func newGServer() *GServer {

	// Connect to DB
	db, err := sql.Open("mysql", "srv:callgov@unix(/var/lib/mysql/mysql.sock)/srv")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("SET NAMES utf8mb4 COLLATE utf8mb4_bin")
	if err != nil {
		log.Fatal(err)
	}

	// Parse template
	tmpl := template.Must(template.ParseGlob(baseDir + "/vhosts/front/*.html"))

	// Init structure
	srv := &GServer{
		tmpl: tmpl, // Templates
		db:   db,   // DB connection
	}

	return srv
}

func (gsrv *GServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//var url string

	err := gsrv.db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	if r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/stat") && filepath.Base(r.URL.Path) == r.URL.Path[1:] {
		//url = "https://" + xoConf.Domain.Site + r.URL.Path

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		gsrv.tmpl = template.Must(template.ParseGlob(baseDir + "/vhosts/front/*.html"))

		var page string = r.URL.Path[1:]
		if len(page) == 0 {
			page = "index"
		}

		if gsrv.tmpl.Lookup(page+".html") != nil {

			var query string
			//var token string = r.FormValue("token")
			//if ok, _, _ := checkStatTokenCore(token, xoConf.Key); ok {
			//	query = "?token=" + token
			//}

			err = gsrv.tmpl.ExecuteTemplate(w, page+".html", struct{ Query string }{query})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println(err)
			}
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}

	} else if strings.HasPrefix(r.URL.Path, "/json/") {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)

	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	}

}
