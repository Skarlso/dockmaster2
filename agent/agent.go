package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
)

//Containers representing information about the running containers
type Containers struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Command string `json:"command"`
	Port    string `json:"port"`
}

//Post represents a post to the server
type Post struct {
	AgentID             string       `json:"agentid"`
	ExpiredAfterSeconds int          `json:"expiredAfterSeconds"`
	Containers          []Containers `json:"containers"`
}

var (
	serverURL          string
	agentID            string
	refreshRate        int
	expireAfterSeconds int
)

func init() {
	flag.StringVar(&serverURL, "server", "http://localhost:8989", "The server uri where dockmaster is located.")
	flag.StringVar(&agentID, "agent", "localhost", "The name of an Agent. Example: TestQA1.")
	flag.IntVar(&refreshRate, "refresh", 60, "The rate at which this agent should check for changes in seconds.")
	flag.IntVar(&expireAfterSeconds, "expireAfterSeconds", 60, "The rate at which data sent by this agent should expire in seconds.")

	flag.Parse()
}

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	for {
		log.Println("Started listening... Refresh rate is :", refreshRate)
		post := Post{AgentID: agentID, ExpiredAfterSeconds: expireAfterSeconds}
		containers := []Containers{}
		runningContainers, _ := client.ListContainers(docker.ListContainersOptions{All: false})
		for _, v := range runningContainers {
			c := Containers{}
			c.ID = v.ID
			c.Name = strings.Join(v.Names, ",")
			for _, p := range v.Ports {
				c.Port += p.IP + ":" + p.Type
			}
			c.Command = v.Command
			containers = append(containers, c)
		}
		post.Containers = containers

		postString, err := json.Marshal(post)
		if err != nil {
			log.Println("Error occured while trying ot marshal POST:", err.Error())
			continue
		}
		req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(postString))
		if err != nil {
			log.Println("Failed to create post request... Trying again later.")
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Failed to receive from server... Trying again later.")
			continue
		}
		defer resp.Body.Close()

		time.Sleep(time.Second * time.Duration(refreshRate))
	}
}
