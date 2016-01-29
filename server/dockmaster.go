package main

import (
	"net/http"
	"time"

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

// The main function which starts the rpg
func main() {
	router := gin.Default()
	v1 := router.Group(APIBASE)
	{
		v1.GET("/list", listContainers)
		v1.POST("/add", addContainers)
		v1.POST("/delete", deleteContainers)
	}
	router.Run(":8989")
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
