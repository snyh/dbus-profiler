package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/snyh/dbus-profiler/frontend"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

type Server struct {
	db       *Database
	StartAt  time.Time
	listener net.Listener
	c        Config
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
		c:        NewConfig(),
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
		info, err := frontend.AssetInfo(url)
		if err != nil {
			fmt.Fprintf(w, "ERR: %v for fetch %q\n", err, url)
		}
		data, err := frontend.Asset(url)
		if err != nil {
			fmt.Fprintf(w, "ERR: %v for fetch %q\n", err, url)
			return
		}

		buf := bytes.NewReader(data)
		http.ServeContent(w, r, info.Name(), info.ModTime(), buf)
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

func (s *Server) Test(w http.ResponseWriter, r *http.Request) {
	s.db.Test("org.freedesktop.DBus.Properties", w)
}
func (s *Server) Run(debug bool) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/p/#/", 301)
	})
	http.Handle("/p/", http.StripPrefix("/p/", ResourceBundle("./frontend", debug)))
	http.HandleFunc("/dbus/api/main", s.Main)
	http.HandleFunc("/dbus/api/info", s.Info)
	http.HandleFunc("/dbus/api/interface", s.RenderInterfaceDetail)
	http.HandleFunc("/dbus/api/test", s.Test)
	http.HandleFunc("/config", s.Config)
	return http.Serve(s.listener, nil)
}

func (s *Server) Config(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var err error
	switch q.Get("enable") {
	case "t":
		err = s.c.Enable(true)
	case "f":
		err = s.c.Enable(false)
	default:
	}
	ww := json.NewEncoder(w)

	if err != nil {
		w.WriteHeader(505)
		ww.Encode(err)
		return
	}
	ww.Encode(struct{ Enable bool }{s.c.CheckEnable()})
}

func (s *Server) OpenBrowser(auto bool) {
	url := fmt.Sprintf("http://%s", s.listener.Addr())
	if auto {
		fmt.Printf("Auto open page disabled \nPlease visit %q manually\n", url)
		return
	}

	var cmd *exec.Cmd
	if bin, err := exec.LookPath("google-chrome"); err == nil {
		cmd = exec.Command(bin, "--app="+url)
	} else {
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Auto open page failed: %v \nPlease visit %q manually\n", err, url)
	}
	go cmd.Wait()
}
