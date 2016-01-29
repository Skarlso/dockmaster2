package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

//Container a single container
type Container struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	RunningCmd string `json:"command"`
	Port       string `json:"port"`
}

//Agent post data from an agent with ID and containers it has.
type Agent struct {
	AgentID            string      `json:"agentid"`
	ExpireAfterSeconds int         `json:"expireAfterSeconds"`
	Containers         []Container `json:"containers"`
}

func index(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("http://localhost:8989/api/1/list")
	agents := []Agent{}
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&agents)
	if err != nil {
		fmt.Fprintf(w, "error occured:"+err.Error())
		return
	}

	lp := path.Join("templates", "layout.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		fmt.Fprintf(w, "error occured:"+err.Error())
		return
	}
	tmpl.Execute(w, agents)
}

func main() {
	http.HandleFunc("/", index)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":9191", nil)
}
