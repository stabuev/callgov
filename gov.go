package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"path/filepath"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	Id    uint32
	Token string
}

type GServer struct {
	tmpl      *template.Template
	db        *sql.DB
	localSKey []byte
	session   map[string]*Account
	muSession sync.Mutex
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

	localSKey := make([]byte, 32)
	io.ReadFull(rand.Reader, localSKey)

	// Init structure
	srv := &GServer{
		tmpl:      tmpl,
		db:        db,
		localSKey: localSKey,
		session:   make(map[string]*Account),
	}

	return srv
}

//
// setSession
//
func (gsrv *GServer) setSession(id uint32) (sid string) {
	// look for active session
	for idx := range gsrv.session {
		if gsrv.session[idx].Id == id {
			sid = idx
			break
		}
	}

	if sid == "" {
		b := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, b); err == nil {
			sid = base64.URLEncoding.EncodeToString(b)
			mac := hmac.New(sha256.New, gsrv.localSKey)
			mac.Write(b)
			gsrv.session[sid] = new(Account)
			gsrv.session[sid].Id = id
			gsrv.session[sid].Token = hex.EncodeToString(mac.Sum(nil)[0:16])
		}
	}
	return
}

func (gsrv *GServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//var url string

	err := gsrv.db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	if r.URL.Path == "/" && r.Method == http.MethodGet || (r.URL.Path == "/login" || r.URL.Path == "/logout") && r.Method == http.MethodPost {

		gsrv.muSession.Lock()

		var sid string

		// cookies and runtime session must exist at the same time
		cookie, err := r.Cookie("sid")
		if err == nil {
			sid = cookie.Value
			_, ok := gsrv.session[sid]
			if !ok {
				sid = ""
			}
		}

		// Login & logout
		if r.Method == http.MethodPost {
			if r.URL.Path == "/login" {
				// Login
				var id uint32
				var name string
				var atype string
				var token string

				w.Header().Set("Content-Type", "application/json")

				var req struct {
					Login    string `json:"login"`
					Password string `json:"password"`
				}

				defer r.Body.Close()
				body, _ := ioutil.ReadAll(r.Body)
				log.Print(string(body))
				err = json.Unmarshal(body, &req)
				log.Println(err)
				if err == nil {
					err = gsrv.db.QueryRow("SELECT id, name, type FROM account WHERE login=? AND password=?", req.Login, req.Password).Scan(&id, &name, &atype)
					if err == nil {
						sid = gsrv.setSession(id)
						if sid != "" {
							//newcookie := http.Cookie{Name: "sid", Value: sid, Path: "/", Domain: r.URL.Host , Expires: time.Now().UTC().Add(time.Hour * 24)}
							newcookie := http.Cookie{Name: "sid", Value: sid, Path: "/", Domain: r.URL.Host}
							http.SetCookie(w, &newcookie)
							token = gsrv.session[sid].Token
						}
					}
				}
				w.Write([]byte("{\"token\":\"" + token + "\"}"))

			} else if r.URL.Path == "/logout" {
				// Logout
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				if sid != "" {
					delete(gsrv.session, sid)
					sid = ""
					// delete cookie
					cookie.Expires = time.Now().UTC().AddDate(0, 0, -1)
					http.SetCookie(w, cookie)

				}
				w.Header().Set("Location", "http://45.128.204.157/")
				w.WriteHeader(http.StatusFound)
				w.Write([]byte("Redirect to http://45.128.204.157/"))
			}
		} else {

			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			var id uint32
			var name string
			var atype, token string
			if sid != "" {
				id = gsrv.session[sid].Id
				err = gsrv.db.QueryRow(`SELECT name, type FROM account WHERE id=?`, id).Scan(&name, &atype)
				log.Println(err)
				token = gsrv.session[sid].Token
			}

			gsrv.tmpl = template.Must(template.ParseGlob(baseDir + "/vhosts/front/*.html"))

			var page string = r.URL.Path[1:]
			if len(page) == 0 {
				page = "index"
			}

			log.Println(id, name, atype, token)
			if gsrv.tmpl.Lookup(page+".html") != nil {
				err = gsrv.tmpl.ExecuteTemplate(w, page+".html", struct {
					ID    uint32
					Name  string
					AType string
					Token string
				}{id, name, atype, token})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					log.Println(err)
				}
			} else {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			}
		}

		gsrv.muSession.Unlock()

	} else if strings.HasPrefix(r.URL.Path, "/json/") {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)

	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	}

}
