package main

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
)

func StartServer(db *Database, addr string) error {
	if addr == "auto" {
		addr = ":0"
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err.Error())
	}

	http.Handle("/debug/", http.StripPrefix("/debug/", http.FileServer(http.Dir("./frontend/dist"))))
	http.HandleFunc("/debug/dbus/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%v", db.data)
	})

	url := fmt.Sprintf("http://%s/debug", l.Addr())
	if err = exec.Command("xdg-open", url).Run(); err != nil {
		fmt.Printf("Auto open page failed: %v \nPlease visit %q manually\n", err, url)
	}
	return http.Serve(l, nil)
}
