package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ServeHTTP - обработка запросов
func (o *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.String() == "/websocket.json" {
		http.Header.Add(w.Header(), "Access-Control-Allow-Origin", "*")
		http.Header.Add(w.Header(), "Content-Type", "application/json;charset=utf-8")
		http.Header.Add(w.Header(), "X-Content-Type-Options", "nosniff")
		http.Header.Add(w.Header(), "Cache-Control", "no-cache")
		o.hWEBSocket(w, r)
		return
	}

	if r.Method == "GET" {

		if r.URL.String() == "/" || r.URL.String() == "/index.html" {
			http.Header.Add(w.Header(), "X-Content-Type-Options", "nosniff")
			http.Header.Add(w.Header(), "Cache-Control", "no-cache")
			o.hRoot(w)
			return
		}

		if r.URL.String() == "/analize.html" {
			http.Header.Add(w.Header(), "X-Content-Type-Options", "nosniff")
			http.Header.Add(w.Header(), "Cache-Control", "no-cache")
			o.hAnalize(w)
			return
		}

		if r.URL.String() == "/favicon.ico" {
			http.Header.Add(w.Header(), "X-Content-Type-Options", "nosniff")
			http.Header.Add(w.Header(), "Cache-Control", "max-age=31536000,immutable")
			http.ServeFile(w, r, o.cfg.Root+"/files/favicon.ico")
			return
		}

		if strings.HasPrefix(r.URL.String(), "/files/") {
			http.Header.Add(w.Header(), "X-Content-Type-Options", "nosniff")
			http.Header.Add(w.Header(), "Cache-Control", "max-age=31536000,immutable")
			http.ServeFile(w, r, o.cfg.Root+r.URL.String())
			return
		}

		if strings.HasPrefix(r.URL.String(), "/sql.json") {
			http.Header.Add(w.Header(), "Access-Control-Allow-Origin", "*")
			http.Header.Add(w.Header(), "Content-Type", "application/json;charset=utf-8")
			http.Header.Add(w.Header(), "X-Content-Type-Options", "nosniff")
			http.Header.Add(w.Header(), "Cache-Control", "no-cache")
			w.WriteHeader(http.StatusOK)
			o.hSQL(w, r)
			return
		}
		
	}

	w.WriteHeader(http.StatusNotFound)
	if ex, err := os.Executable(); err == nil {
		io.WriteString(w, "Path: "+filepath.Dir(ex)+"\n")
	}
	io.WriteString(w, "URL: "+o.cfg.Root+r.URL.String())
}

// hError - отвечаем страницей с ошибкой
func (o *Server) hError(w http.ResponseWriter, arg string) {
	tmp, err := template.ParseFiles(o.cfg.Root + "/templates/error.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "ERROR: ["+fmt.Sprint(err)+"]")
	} else {
		w.WriteHeader(http.StatusOK)
		if err = tmp.Execute(w, struct{ ErrorDescription string }{arg}); err != nil {
			log.Printf("HTTP ERROR [%s]", arg)
		}
	}
}

// hRoot - отвечаем главной страницей
func (o *Server) hRoot(w http.ResponseWriter) {
	tmp, err := template.ParseFiles(o.cfg.Root + "/templates/index.html")
	if err != nil {
		o.hError(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		if err = tmp.Execute(w, struct {
			ADDRESS string
			PORT    string
		}{
			o.cfg.Address,
			o.cfg.Port,
		}); err != nil {
			o.hError(w, err.Error())
		}
	}
}

// hRoot - отвечаем страницей анализа тревог
func (o *Server) hAnalize(w http.ResponseWriter) {
	tmp, err := template.ParseFiles(o.cfg.Root + "/templates/analize.html")
	if err != nil {
		o.hError(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		if err = tmp.Execute(w, struct {
			ADDRESS string
			PORT    string
		}{
			o.cfg.Address,
			o.cfg.Port,
		}); err != nil {
			o.hError(w, err.Error())
		}
	}
}
