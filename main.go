package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	baseDir string
)

func NoAutoIndexHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

type SServer struct {
	gsrv *GServer
}

//
// ServeHTTP
//
func (ssrv *SServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" || r.URL.Path == "/robots.txt" {
		http.NotFound(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/cdn/") {
		NoAutoIndexHandler(http.StripPrefix("/cdn/", http.FileServer(http.Dir(baseDir+"/vhosts/front/cdn")))).ServeHTTP(w, r)
	} else {
		ssrv.gsrv.ServeHTTP(w, r)
	}
}

func main() {
	var err error

	// Get basedir
	var progpath string
	progpath, err = os.Executable()
	if err != nil {
		log.Fatalln("Error while getting path name of process executable.", err)
	}
	baseDir = filepath.Clean(filepath.Dir(progpath) + "/..")

	var handler *SServer = &SServer{gsrv: newGServer()}
	srv := &http.Server{
		Addr:    "45.128.204.157:80",
		Handler: handler,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("Server error: %s\n", err)
	}

}
