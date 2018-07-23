// See LICENSE.txt for licensing information.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

import (
	"github.com/gorilla/mux"
)

var webuiLog *log.Logger

func handleGetList(w http.ResponseWriter, r *http.Request) {
	j := json.NewEncoder(w)
	w.Header()["Content-Type"] = []string{"application/json"}
	j.Encode(cfg)
}

func handleGetGameBackup(w http.ResponseWriter, r *http.Request) {
	g, ok := getGame(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		webuiLog.Println("FAILED to parse form")
		http.Error(w, "Couldn't parse form.", 400)
		return
	}
	note := r.Form.Get("note")

	sv, err := g.Backup()
	if err != nil {
		webuiLog.Println("FAILED to save: ", err)
		http.Error(w, fmt.Sprintln("Failed to save: ", err), 501)
		return
	}
	if note != "" {
		sv.Note = note
	}

	outputJSON(w, sv)
}

func handleGetGameRestore(w http.ResponseWriter, r *http.Request) {
	g, ok := getGame(w, r)
	if !ok {
		return
	}
	v := mux.Vars(r)
	i, _ := strconv.Atoi(v["id"])
	sv, err := g.Restore(i)
	if err != nil {
		webuiLog.Println("FAILED to restore: ", err)
		http.Error(w, fmt.Sprintln("Failed to restore: ", err), 501)
		return
	}
	g.Stamp = time.Now()

	outputJSON(w, sv)
}

func loggingMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webuiLog.Printf("Request to %s\n", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func startWebUI() {
	webuiLog = log.New(os.Stdout, "", log.Ltime|log.Lshortfile)

	r := mux.NewRouter()
	r.HandleFunc("/list", handleGetList).Methods("GET")
	r.HandleFunc("/{game}/backup", handleGetGameBackup).Methods("GET")
	r.HandleFunc("/{game}/restore/{id:[0-9]+}", handleGetGameRestore).Methods("GET")
	r.Use(loggingMw)
	http.Handle("/", r)

	addr := fmt.Sprintf("127.0.0.1:%d", flagPort)
	webuiLog.Printf("Starting Web UI at http://%s\n", addr)
	http.ListenAndServe(addr, nil)
}

func getGame(w http.ResponseWriter, r *http.Request) (*Game, bool) {
	v := mux.Vars(r)
	game := v["game"]
	g := cfg.GetGame(game)
	if g == nil {
		webuiLog.Printf("Game '%s' not found\n", game)
		http.NotFound(w, r)
		return nil, false
	}
	return g, true
}

func outputJSON(w http.ResponseWriter, data interface{}) {
	j := json.NewEncoder(w)
	w.Header()["Content-Type"] = []string{"application/json"}
	j.Encode(data)
}
