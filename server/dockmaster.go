package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
)

//APIVERSION Is the current API version
const APIVERSION = "1"

//APIBASE Defines the API base URI
const APIBASE = "api/" + APIVERSION

var mdb MongoDBConnection

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
	CreatedAt          time.Time   `bson:"createdAt"`
	IP                 string      `json:"ip"`
	Port               string      `json:"port"`
}

func init() {
	mdb = MongoDBConnection{}
	go func() {
		for {
			//This will remove old data periodically at every minute.
			//Mongodb's own cleanser also ticks at every 60 seconds.
			mdb.startCleansing()
			time.Sleep(time.Minute)
		}
	}()
}

func main() {
	router := gin.Default()
	v1 := router.Group(APIBASE)
	{
		v1.GET("/list", listContainers)
		v1.POST("/add", addContainers)
		v1.POST("/delete", deleteContainers)
		v1.GET("/inspect/:agentID/:containerID", inspectContainer)
		v1.OPTIONS("/inspect/:agentID/:containerID", preflight)
		v1.POST("/stopAll", stopAll)
		v1.OPTIONS("/stopAll", preflight)
	}
	router.Run(":8989")
}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	// c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, Authorization, Content-Type")
	c.JSON(http.StatusOK, struct{}{})
}

func stopAll(c *gin.Context) {
	var agentID struct {
		AgentID string `json:"agentid"`
	}
	err := c.BindJSON(&agentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"error binding json: " + err.Error()})
		return
	}

	a, err := mdb.GetAgent(agentID.AgentID)
	resp, _ := http.Post("http://"+a.IP+":"+a.Port+"/"+APIBASE+"/stopAll", "application/json", nil)

	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, resp.Body)
	}

	c.JSON(resp.StatusCode, Message{"All containers stopped."})
}

func inspectContainer(c *gin.Context) {
	cID := c.Param("containerID")
	agentID := c.Param("agentID")
	a, err := mdb.GetAgent(agentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"Failed to load Agent:" + err.Error()})
		return
	}

	container := docker.Container{}
	resp, _ := http.Get("http://" + a.IP + ":" + a.Port + "/" + APIBASE + "/inspect/" + cID)

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&container)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{"Could not retrieve container:" + err.Error()})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, container)
}

//index a humble welcome to a new user
func listContainers(c *gin.Context) {
	agents, err := mdb.Load()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"error while loading containers: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, agents)
}

func addContainers(c *gin.Context) {
	agent := Agent{}
	err := c.BindJSON(&agent)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"error binding json: " + err.Error()})
		return
	}
	err = mdb.Save(agent)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"error while saving container: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, Message{"Containers successfully saved."})

}
func deleteContainers(c *gin.Context) {
	agent := Agent{}
	err := c.BindJSON(&agent)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"error binding json: " + err.Error()})
		return
	}
	err = mdb.Delete(agent)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"error while deleting containers: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, Message{"Containers successfully removed."})
}
