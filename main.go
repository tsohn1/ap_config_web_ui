package main

import (
	"ap_config_web_ui/config"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/fatih/structs"
	"log"
	"strconv"
	"strings"
	"google.golang.org/grpc"
	"context"
	"google.golang.org/grpc/credentials/insecure"
	validate "ap_config_web_ui/validate"
)

const (
	YAML_FOLDER = "config_files/"
	NETWORK_ENV = YAML_FOLDER + "network.yaml"
	VALIDATE_YAML_CHANGES = true //flag to validate YAML changes using gRPC
	GRPC_SUCCESS_TOKEN = 1
	GRPC_FAIL_TOKEN = 0
)

var (
	network_env config.NetworkEnv
	client validate.ValidateClient
)

func updateConfig(ctx *gin.Context) {
	err := ctx.Request.ParseForm()
	if err != nil {
		// Handle the error, possibly by returning an error response
		return
	}

	//Retrieve all form fields
	formFields := ctx.Request.PostForm

	//convert SiteCode to uint32
	SiteCode, err := strconv.Atoi(formFields.Get("SiteCode"))

	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}

	//convert NameServers, TimeserverUrls into []string
	nameServerRaw := formFields.Get("NameServers")
	timeServerRaw := formFields.Get("TimeServerUrls")
	nameServerRaw = strings.Trim(nameServerRaw, "[]")
	timeServerRaw = strings.Trim(timeServerRaw, "[]")

	new_network_struct := config.NetworkEnv{   
		SiteId : formFields.Get("SiteId"),
		SiteCode: uint32(SiteCode),
		StoreCode: formFields.Get("StoreCode"),
		Ip : formFields.Get("Ip"),                
		DefaultGwIP : formFields.Get("DefaultGwIP"),       
		Netmask : formFields.Get("Netmask"),           
		NameServers : strings.Split(nameServerRaw, ","),       
		TimeZone : formFields.Get("TimeZone"),          
		TimeServerUrls : strings.Split(timeServerRaw, ","),    
		InterApPort : formFields.Get("InterApPort"),       
		InterApPortTarget : formFields.Get("InterApPortTarget"), 
		ApBrokerUrl : formFields.Get("ApBrokerUrl"),       
		EthernetInterface : formFields.Get("EthernetInterface")}

	err = config.SetConfigEnv(NETWORK_ENV, &new_network_struct)

	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
	if VALIDATE_YAML_CHANGES {
		verResponse, err := client.Verify(context.Background(), &validate.VerifyRequest{Token: GRPC_SUCCESS_TOKEN})
		if err != nil {
			log.Fatalf("Verify failed: %v", err)
		}
		if verResponse.IsValid {
			log.Printf("Verify result: Valid")
		} else {
			log.Printf("Verify result: Invalid")
		}
	}
	ctx.Redirect(http.StatusSeeOther, "/")
}




func main() {

	if VALIDATE_YAML_CHANGES {
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()
		client = validate.NewValidateClient(conn)
	}
	
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// display YAML data in webpage
	router.GET("/", func(ctx *gin.Context) {
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