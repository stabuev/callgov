package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

func (gsrv *GServer) getSession(r *http.Request) (sid string) {
	// cookies and runtime session must exist at the same time
	cookie, err := r.Cookie("sid")
	if err == nil {
		sid = cookie.Value
		_, ok := gsrv.session[sid]
		if !ok {
			sid = ""
		}
	}
	// check token
	if sid == "" {
		token, ok := r.URL.Query()["token"]
		if ok {
			for idx := range gsrv.session {
				if gsrv.session[idx].Token == token[0] {
					sid = idx
					break
				}
			}
		}
	}
	return
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

	var now string = time.Now().UTC().Format("2006-01-02 15:04:05")

	err := gsrv.db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	gsrv.muSession.Lock()
	var sid string = gsrv.getSession(r)

	if (r.URL.Path == "/" || r.URL.Path == "/register" || r.URL.Path == "/obr" || r.URL.Path == "/detail") &&
		r.Method == http.MethodGet || (r.URL.Path == "/login" || r.URL.Path == "/logout") && r.Method == http.MethodPost {

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
				err = json.Unmarshal(body, &req)
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
					cookie, _ := r.Cookie("sid")
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
				token = gsrv.session[sid].Token
			}

			gsrv.tmpl = template.Must(template.ParseGlob(baseDir + "/vhosts/front/*.html"))

			var page string = r.URL.Path[1:]
			if len(page) == 0 {
				page = "index"
			}

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

	} else if r.URL.Path == "/json/obr" {

		w.Header().Set("Content-Type", "application/json")
		var status string = "error"
		if sid != "" {
			var req struct {
				Id      uint32 `json:"id"`
				Title   string `json:"title"`
				Content string `json:"content"`
				Public  uint8  `json:"public"`
				State   string `json:"state"`
				Address string `json:"address"`
			}

			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			log.Println(string(body))
			err = json.Unmarshal(body, &req)
			if err == nil {
				if req.Id == 0 {
					// Create obr
					_, err = gsrv.db.Exec("INSERT INTO obr (title, content, account_id, public, state, address) VALUES(?,?,?,?,?,?)",
						req.Title, req.Content, gsrv.session[sid].Id, req.Public, req.State, req.Address)
					if err == nil {
						status = "ok"
					}
				} else {
					_, err = gsrv.db.Exec("UPDATE obr SET title=?, content=?, public=?, state=?, address=?, dtlast=? WHERE id=?",
						req.Title, req.Content, req.Public, req.State, req.Address, now, req.Id)
					if err == nil {
						status = "ok"
					}
				}
			} else {
				log.Println(err)
			}
		}
		w.Write([]byte("{\"status\":\"" + status + "\"}"))

	} else if r.URL.Path == "/json/comment" {

		w.Header().Set("Content-Type", "application/json")
		var status string = "error"
		if sid != "" {
			var req struct {
				Id      uint32 `json:"id"`
				Content string `json:"content"`
			}

			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			log.Println(string(body))
			err = json.Unmarshal(body, &req)
			if err == nil && req.Id != 0 {
				var state string
				gsrv.db.QueryRow("SELECT state FROM obr WHERE id=?", req.Id).Scan(&state)
				_, err = gsrv.db.Exec("INSERT INTO comment (obr_id, account_id, content, type) VALUES(?,?,?,?)",
					req.Id, gsrv.session[sid].Id, req.Content, state)
				if err == nil {
					status = "ok"
				}
			}
		}
		w.Write([]byte("{\"status\":\"" + status + "\"}"))

	} else if r.URL.Path == "/json/sign" {

		w.Header().Set("Content-Type", "application/json")
		var status string = "error"
		if sid != "" {
			var req struct {
				Id uint32 `json:"id"`
			}

			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			err = json.Unmarshal(body, &req)
			if err == nil {
				if req.Id != 0 {
					// Create obr
					_, err = gsrv.db.Exec("INSERT INTO obr_sign (obr_id, account_id) VALUES(?,?)", req.Id, gsrv.session[sid].Id)
					if err == nil {
						status = "ok"
					}
				}
			}
		}
		w.Write([]byte("{\"status\":\"" + status + "\"}"))

	} else if r.URL.Path == "/json/obrlist" {
		var rows *sql.Rows
		w.Header().Set("Content-Type", "application/json")

		var where string

		var req struct {
			Id uint32 `json:"id"`
		}

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		//log.Println(string(body))
		err = json.Unmarshal(body, &req)
		if err == nil && req.Id != 0 {
			where = fmt.Sprintf("WHERE o.id=%d", req.Id)
		}
		var signed string
		if sid != "" {
			signed = fmt.Sprintf("(SELECT COUNT(*) FROM obr_sign os WHERE o.id=os.obr_id AND os.account_id=%d)", gsrv.session[sid].Id)
		} else {
			signed = "3"
		}

		var jresp bytes.Buffer
		jresp.WriteString("{\"obr\":[")
		var stmt string = fmt.Sprintf(`SELECT
			o.id,
			o.title,
			o.content,
			a.name,
			o.public,
			o.state,
			o.address,
			o.dtreg,
			o.dtlast,
			(SELECT COUNT(*) FROM obr_sign os WHERE o.id=os.obr_id),
			%s,
			(SELECT COUNT(*) FROM comment cc WHERE o.id=cc.obr_id)
			FROM obr o
			LEFT JOIN account a ON o.account_id=a.id
			%s ORDER by o.dtlast DESC`, signed, where)
		log.Println(stmt)

		rows, err = gsrv.db.Query(stmt)
		log.Println(err)
		if err == nil {
			defer rows.Close()
			var cnt int
			for rows.Next() {
				var id, title, content, name, public, state, address, dtreg, dtlast, totalsign, sign, ccnt string
				err = rows.Scan(&id, &title, &content, &name, &public, &state, &address, &dtreg, &dtlast, &totalsign, &sign, &ccnt)
				if err == nil {
					if cnt > 0 {
						jresp.WriteString(",")
					}
					jresp.WriteString("[\"")
					jresp.WriteString(strings.Join([]string{id, title, content, name, public, state, address, dtreg, dtlast, totalsign, sign, ccnt}, "\",\""))
					jresp.WriteString("\"]")
					cnt++
				}
			}
		}
		jresp.WriteString("]}")
		w.Write(jresp.Bytes())

	} else if r.URL.Path == "/json/commentlist" {
		var rows *sql.Rows
		w.Header().Set("Content-Type", "application/json")

		var req struct {
			Id uint32 `json:"id"`
		}

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		log.Println(string(body))
		err = json.Unmarshal(body, &req)
		log.Println(err)

		var jresp bytes.Buffer
		jresp.WriteString("{\"comment\":[")

		if err == nil && req.Id != 0 {

			var stmt string = fmt.Sprintf(`SELECT
			c.content,
			a.name,
			c.dt
			FROM comment c
			LEFT JOIN account a ON c.account_id=a.id
			WHERE c.obr_id=%d
			ORDER by c.dt DESC`, req.Id)
			log.Println(stmt)
			rows, err = gsrv.db.Query(stmt)
			if err == nil {
				defer rows.Close()
				var cnt int
				for rows.Next() {
					var content, name, dt string
					err = rows.Scan(&content, &name, &dt)
					if err == nil {
						if cnt > 0 {
							jresp.WriteString(",")
						}
						jresp.WriteString("[\"")
						jresp.WriteString(strings.Join([]string{content, name, dt}, "\",\""))
						jresp.WriteString("\"]")
						cnt++
					}
				}
			}
		}

		jresp.WriteString("]}")
		w.Write(jresp.Bytes())

	} else if r.URL.Path == "/json/signlist" {
		var rows *sql.Rows
		w.Header().Set("Content-Type", "application/json")

		var req struct {
			Id uint32 `json:"id"`
		}

		var jresp bytes.Buffer
		jresp.WriteString("{\"sign\":[")

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &req)
		if err == nil && req.Id != 0 {

			var stmt string = fmt.Sprintf(`SELECT
			a.id,
			a.name,
			os.dt
			FROM obr_sign os
			LEFT JOIN account a ON os.account_id=a.id
			WHERE os.obr_id=%d ORDER BY os.dt ASC`, req.Id)
			rows, err = gsrv.db.Query(stmt)
			if err == nil {
				defer rows.Close()
				var cnt int
				for rows.Next() {
					var id, name, dt string
					err = rows.Scan(&id, &name, &dt)
					if err == nil {
						if cnt > 0 {
							jresp.WriteString(",")
						}
						jresp.WriteString("[\"")
						jresp.WriteString(strings.Join([]string{id, name, dt}, "\",\""))
						jresp.WriteString("\"]")
						cnt++
					}
				}
			}
		}
		jresp.WriteString("]}")
		w.Write(jresp.Bytes())

	} else if strings.HasPrefix(r.URL.Path, "/json/") {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)

	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	}

	gsrv.muSession.Unlock()

}
