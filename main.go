package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rishabh625/websocket-example/auth"
)

var dataMap map[string]string

type UserDetails struct {
	UserData []User `json:"userdata"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	if dataMap == nil {
		dataMap = make(map[string]string)
	}
	plan, _ := ioutil.ReadFile("userdetails.json")
	var data UserDetails
	err := json.Unmarshal(plan, &data)
	if err != nil {
		panic(err)
	}
	for _, v := range data.UserData {
		dataMap[v.Username] = v.Password
	}
}

func main() {
	go h.run()
	router := gin.New()

	router.GET("/token", func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")

		if _, ok := dataMap[username]; !ok {
			c.JSON(http.StatusUnauthorized, "{\"message\":\"Incorrect Username and Password\"}")
			c.Abort()
		} else {
			if dataMap[username] != password {
				c.JSON(http.StatusUnauthorized, "{\"message\":\"Incorrect Username and Password\"}")
				c.Abort()
			} else {
				c.Next()
			}
		}

	}, func(c *gin.Context) {
		username := c.Param("username")
		token, err := auth.CreateToken(username)
		if err == nil {
			c.JSON(http.StatusAccepted, token)
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
	})

	router.GET("/ws", auth.AuthenticateHeader, func(c *gin.Context) {
		roomId := "1"
		clientid := c.MustGet("clientid").(float64)
		serveWs(c.Writer, c.Request, roomId, clientid)
	})

	router.Run("0.0.0.0:8080")
}
