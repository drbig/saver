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
	outputJSON(w, cfg)
}

func handleGetGameBackup(w http.ResponseWriter, r *http.Request) {
	g, ok := getGame(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		webuiLog.Println("FAILED to parse form")
		http.Error(w, "Couldn't parse form.", http.StatusBadRequest)
		return
	}
	sv, err := g.Backup()
	if err != nil {
		webuiLog.Println("FAILED to save: ", err)
		http.Error(w, fmt.Sprintln("Failed to save: ", err), http.StatusInternalServerError)
		return
	}
	sv.Note = r.Form.Get("note")

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
		http.Error(w, fmt.Sprintln("Failed to restore: ", err), http.StatusInternalServerError)
		return
	}
	g.Stamp = time.Now()

	outputJSON(w, sv)
}

func handleGetGameDelete(w http.ResponseWriter, r *http.Request) {
	g, ok := getGame(w, r)
	if !ok {
		return
	}
	v := mux.Vars(r)
	f, _ := strconv.Atoi(v["from"])
	t, _ := strconv.Atoi(v["to"])
	_, err := g.Delete(f, t)
	if err != nil {
		webuiLog.Println("FAILED to delete: ", err)
		http.Error(w, fmt.Sprintln("Failed to delete: ", err), http.StatusInternalServerError)
	}

	outputJSON(w, cfg)
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
	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/list", handleGetList).
		Methods(http.MethodGet)
	api.HandleFunc("/{game}/backup", handleGetGameBackup).
		Methods(http.MethodGet)
	api.HandleFunc("/{game}/restore/{id:[0-9]+}", handleGetGameRestore).
		Methods(http.MethodGet)
	api.HandleFunc("/{game}/delete/{from:[0-9]+}-{to:[0-9]+}", handleGetGameDelete).
		Methods(http.MethodGet)
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
