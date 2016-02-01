package main

import (
	"net/http"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
)

func stopContainer(c *gin.Context) {
	var container struct {
		ID string `json:"id"`
	}

	//reflect.ValueOf(&t).MethodByName("Foo").Call([]reflect.Value{})
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)

	err := client.StopContainer(container.ID, 1)
	if err != nil {
		e := ErrorResponse{}
		e.ErrorMessage = "error stopping containers:" + err.Error()
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	m := Message{}
	m.Message = "continer stopped successfully"
	c.JSON(http.StatusOK, m)
}

func stopAllContainers(c *gin.Context) {

	var stopErrors [][]string
	for i := range stopErrors {
		stopErrors[i] = make([]string, 0)
	}
	//reflect.ValueOf(&t).MethodByName("Foo").Call([]reflect.Value{})
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	runningContainers, err := client.ListContainers(docker.ListContainersOptions{All: false})

	if err != nil {
		e := ErrorResponse{}
		e.ErrorMessage = "error getting all containers:" + err.Error()
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	for _, v := range runningContainers {
		err := client.StopContainer(v.ID, 1)
		if err != nil {
			stopErrors = append(stopErrors, v.Names)
		}
	}

	if len(stopErrors) != 0 {
		var errCon string
		for _, v := range stopErrors {
			errCon += "[" + strings.Join(v, "|") + "]"
		}
		e := ErrorResponse{}
		e.ErrorMessage = "error stopping containers:" + errCon
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	m := Message{}
	m.Message = "all containers successfully stopped"
	c.JSON(http.StatusOK, m)
}

func inspectContainer(c *gin.Context) {
	cID := c.Param("id")
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	con, err := client.InspectContainer(cID)
	if err != nil {
		e := ErrorResponse{"error while trying to inspect container:" + err.Error()}
		c.JSON(http.StatusInternalServerError, e)
	}

	c.JSON(http.StatusOK, con)
}
