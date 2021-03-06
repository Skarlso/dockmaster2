package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
)

//APIBASE The base of the API that this agent provides
const APIBASE = "/api/1"

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
	ExpiredAfterSeconds int          `json:"expireAfterSeconds"`
	IPAddress           string       `json:"ip"`
	Port                string       `json:"port"`
	Containers          []Containers `json:"containers"`
}

var (
	serverURL          string
	agentID            string
	refreshRate        int
	expireAfterSeconds int
	port               string
)

func init() {
	flag.StringVar(&serverURL, "server", "http://localhost:8989", "The server uri where dockmaster is located.")
	flag.StringVar(&agentID, "agent", "localhost", "The name of an Agent. Example: TestQA1.")
	flag.IntVar(&refreshRate, "refresh", 60, "The rate at which this agent should check for changes in seconds.")
	flag.IntVar(&expireAfterSeconds, "expireAfterSeconds", 60, "The rate at which data sent by this agent should expire in seconds.")
	flag.StringVar(&port, "port", "9999", "The port number on which this agent is running on.")

	flag.Parse()
	go startDiscovering()
}

func main() {
	router := gin.Default()
	v1 := router.Group(APIBASE)
	{
		v1.POST("/stop", stopContainer)
		v1.POST("/stopAll", stopAllContainers)
		v1.GET("/inspect/:id", inspectContainer)
	}
	router.Run(":" + port)
}

func startDiscovering() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	log.Println("Started listening... Refresh rate is :", refreshRate)
	for {
		log.Println("Assembling Post...")
		post := Post{AgentID: agentID, ExpiredAfterSeconds: expireAfterSeconds}
		containers := []Containers{}
		runningContainers, err := client.ListContainers(docker.ListContainersOptions{All: false})
		if err != nil {
			panic("Failed to connect to Docker Client." + err.Error())
		}
		for _, v := range runningContainers {
			c := Containers{}
			c.ID = v.ID
			c.Name = strings.Join(v.Names, ",")
			for _, p := range v.Ports {
				c.Port += p.IP + ":" + p.Type + ":" + strconv.Itoa(int(p.PrivatePort)) + ":" + strconv.Itoa(int(p.PublicPort))
			}
			c.Command = v.Command
			containers = append(containers, c)
		}
		post.Containers = containers
		post.IPAddress = getLocalIP()
		post.Port = port
		postString, err := json.Marshal(post)
		if err != nil {
			log.Println("Error occured while trying ot marshal POST:", err.Error())
		}
		req, err := http.NewRequest("POST", serverURL+"/api/1/add", bytes.NewBuffer(postString))
		if err != nil {
			log.Println("Failed to create post request... Trying again later.")
			time.Sleep(time.Second * time.Duration(refreshRate))
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		log.Println("Posting JSON:", string(postString))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Failed to receive from server... Trying again later.")
			time.Sleep(time.Second * time.Duration(refreshRate))
			continue
		}
		defer resp.Body.Close()
		//TODO: Verify the response
		time.Sleep(time.Second * time.Duration(refreshRate))
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("Couldn't determine ip address. Sending empty string. Error:", err.Error())
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
