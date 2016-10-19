package main

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"time"

	"github.com/snyh/dbus-profiler/frontend"
	"strconv"
)

type Server struct {
	db       *Database
	StartAt  time.Time
	listener net.Listener
}

func NewServer(db *Database, addr string) *Server {
	if addr == "auto" {
		addr = ":0"
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	return &Server{
		db:       db,
		listener: l,
	}
}

func (s *Server) Main(w http.ResponseWriter, r *http.Request) {
	top, _ := strconv.ParseInt(r.URL.Query().Get("top"), 0, 8)
	if top == 0 {
		top = 20
	}

	var dur time.Duration

	if since := r.URL.Query().Get("since"); since != "" {
		var err error
		dur, err = time.ParseDuration(since)
		if err != nil {
			http.Error(w, err.Error(), 403)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	s.db.Render(w, int(top), dur)
}

func (s *Server) Info(w http.ResponseWriter, r *http.Request) {
	s.db.RenderGlobalInfo(w)
}

func ResourceBundle(dir string, debug bool) http.Handler {
	if debug {
		return http.FileServer(http.Dir(dir))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if url == "" {
			url = "index.html"
		}
		data, err := frontend.Asset(url)
		if err != nil {
			fmt.Fprintf(w, "ERR: %v for fetch %q\n", err, url)
			return
		}
		w.Write(data)
	})
}

func (s *Server) RenderInterfaceDetail(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	err := s.db.RenderInterfaceDetail(name, w)
	if err != nil {
		fmt.Fprintf(w, "ERR: %v for fetch %q\n", err, name)
	}
}

func (s *Server) Run(debug bool) error {
	http.Handle("/static/", http.StripPrefix("/static/", ResourceBundle("./frontend", debug)))
	http.HandleFunc("/dbus/api/main", s.Main)
	http.HandleFunc("/dbus/api/info", s.Info)
	http.HandleFunc("/dbus/api/interface", s.RenderInterfaceDetail)
	return http.Serve(s.listener, nil)
}

func (s *Server) OpenBrowser(auto bool) {
	url := fmt.Sprintf("http://%s/static/index.html", s.listener.Addr())
	if auto {
		fmt.Printf("Auto open page disabled \nPlease visit %q manually\n", url)
		return
	}
	if err := exec.Command("xdg-open", url).Run(); err != nil {
		fmt.Printf("Auto open page failed: %v \nPlease visit %q manually\n", err, url)
	}
}
