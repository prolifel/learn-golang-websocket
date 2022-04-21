package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Weather struct {
	Status []StatusStruct `json:"status"`
}

type StatusStruct struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

// var wind, water int
var isSent bool = true

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var weatherGenerator = func() Weather {
	return Weather{
		Status: []StatusStruct{
			{
				Water: rand.Intn(100),
				Wind:  rand.Intn(100),
			},
		},
	}
}

func init() {
	// Create initial weather data
	weather := weatherGenerator()

	// Write initial data to data.json
	content, err := json.Marshal(weather)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("./data.json", content, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	go updateJSON()

	server := gin.Default()

	server.LoadHTMLFiles("./index.html")
	server.Static("/css", "./css")

	server.GET("/", func(c *gin.Context) {
		// c.File("./index.html")
		c.HTML(http.StatusOK, "index.html", nil)
	})
	server.GET("/ws", websocketEndpoint)

	server.Run(":8080")
}

func updateJSON() {
	for {
		weather := weatherGenerator()

		// Write initial data to data.json
		content, err := json.Marshal(weather)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("./data.json", content, 0644)
		if err != nil {
			panic(err)
		}

		isSent = false
		// time.Sleep(2 * time.Second)
		time.Sleep(15 * time.Second)
	}
}

func websocketEndpoint(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	defer ws.Close()

	for {
		for !isSent {
			content, err := ioutil.ReadFile("./data.json")
			if err != nil {
				panic(err)
			}

			weather := Weather{}
			err = json.Unmarshal(content, &weather)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%+v\n", weather.Status[0])

			err = ws.WriteJSON(weather)
			if err != nil {
				panic(err)
			}
			isSent = true
		}
		if err != nil {
			log.Fatal(err)
		}
	}

}
