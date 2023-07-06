package main

import (
	"ap_config_web_ui/config"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/fatih/structs"
	"log"
	"strconv"
)

const (
	YAML_FOLDER = "config_files/"
	NETWORK_ENV = YAML_FOLDER + "network.yaml"
)


var network_env config.NetworkEnv

func updateConfig(ctx *gin.Context) {

	//Retrieve all form fields
	formFields := ctx.Request.PostForm
	SiteCode, err := strconv.Atoi(formFields.Get("SiteCode"))
	if err != nil {
		log.Println("Error:", err)
		SiteCode = 1
	}

	file, _ := os.OpenFile("form.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Flags() | log.LstdFlags)
	log.Printf("start")
	for key, values := range formFields {
		log.Printf("printing key")
		log.Printf("%s:\n", key)
		for _, value := range values {
			log.Printf("   printing value")
			log.Printf("-s %s\n", value)
		}
	}

	new_network_struct := config.NetworkEnv{   
		SiteId : formFields.Get("SiteId"),
		SiteCode: uint32(SiteCode),
		StoreCode: formFields.Get("StoreCode"),
		Ip : formFields.Get("Ip"),                
		DefaultGwIP : formFields.Get("DefaultGwIP"),       
		Netmask : formFields.Get("Netmask"),           
		NameServers : append(make([]string, 0), formFields.Get("NameServers")),       
		TimeZone : formFields.Get("TimeZone"),          
		TimeServerUrls : append(make([]string, 0), formFields.Get("TimeServerUrls")),    
		InterApPort : formFields.Get("InterApPort"),       
		InterApPortTarget : formFields.Get("InterApPortTarget"), 
		ApBrokerUrl : formFields.Get("ApBrokerUrl"),       
		EthernetInterface : formFields.Get("EthernetInterface")}
	err = config.SetConfigEnv(NETWORK_ENV, &new_network_struct)
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
	ctx.Redirect(http.StatusSeeOther, "/")
}




func main() {


	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// display YAML data in webpage
	router.GET("/", func(ctx *gin.Context) {
		log.Println("GET")
		//parse YAML from yaml directory
		parsed_network_env, err := config.GetConfigEnv(NETWORK_ENV, &network_env)
		if err != nil {
			log.Println("Error:", err)
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