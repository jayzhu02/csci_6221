package main

import (
	githubtrending2 "6221/githubtrending"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	fmt.Println("runnning the spider")
	filename := "./static/githubtrending_Any_daily"
	githubtrending2.TrendingStart("", "daily")
	githubtrending2.DrawBar(filename)
	githubtrending2.DrawWordCloud(filename)
	githubtrending2.DrawOverlap(filename)

	r := gin.Default()
	r.StaticFS("/static", http.Dir("./static"))
	r.LoadHTMLFiles("./static/login.html")

	r.GET("/", func(c *gin.Context) {

		c.HTML(200, "login.html", "github dashboard")
	})

	r.Run(":8080")




}
