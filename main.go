package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	var hub = newHub()
	r := gin.Default()

	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/create", func(c *gin.Context) {
		createRoomID(hub, c)
	})
	r.GET("/connect", connectRoom)
	go hub.Run()
	r.Run()
}
