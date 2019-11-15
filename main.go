package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/Trans/:text", func(c *gin.Context) {
		text := c.Param("text")
		resp, err := http.Get("http://translate.google.cn/translate_a/single?client=gtx&dt=t&dj=1&ie=UTF-8&sl=auto&tl=zh_TW&q=" + url.QueryEscape(text))
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}
		c.JSON(200, string(body))
	})

	router.GET("/TopStories/:number", func(c *gin.Context) {
		number, err := strconv.Atoi(c.Param("number"))
		resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty")
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}
		var TopStories []int
		json.Unmarshal([]byte(string(body)), &TopStories)
		c.JSON(200, TopStories[0:number])
	})

	router.Run(":" + port)
}
