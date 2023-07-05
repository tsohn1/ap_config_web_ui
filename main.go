package main

import (
	"ap_config_web_ui/config"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

)

var network_env config.NetworkEnv

func main() {
	
	parsed_network_env, err := config.GetConfigEnv("config_files/network.yaml", &network_env)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", parsed_network_env)
	})

	router.Run(":8080")

}