package main

import (
	"ap_config_web_ui/config"
	"fmt"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/fatih/structs"

)

const (
	YAML_FOLDER = "config_files/"
	NETWORK_ENV = YAML_FOLDER + "network.yaml"
)

var network_env config.NetworkEnv
func updateConfig(ctx *gin.Context) {

	//Retrieve all form fields
	formFields := ctx.Request.PostForm

	for key, value := range formFields {
			fmt.Println(key, value)
	}
	ctx.Redirect(http.StatusSeeOther, "/")
}




func main() {


	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// display YAML data in webpage
	router.GET("/", func(ctx *gin.Context) {
		//parse YAML from yaml directory
		parsed_network_env, err := config.GetConfigEnv(NETWORK_ENV, &network_env)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		//type assertion to networkenv
		network_env_struct := parsed_network_env.(*config.NetworkEnv)

		//change interface to map type to pass to html
		network_env_map := structs.Map(network_env_struct)
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"YAMLData" : network_env_map,
		})
	})

	// write to YAML when form is submitted, update webpage as well
	router.POST("/", updateConfig)

	router.Run(":8080")

}