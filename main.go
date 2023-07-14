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
	OPERATION_ENV = YAML_FOLDER + "operation.yaml"
	VALIDATE_YAML_CHANGES = true //flag to validate YAML changes using gRPC
	GRPC_SUCCESS_TOKEN = 1
	GRPC_FAIL_TOKEN = 0
)

var (
	network_env config.NetworkEnv
	operation_env config.OperationEnv
	client validate.ValidateClient
	operationFieldTypes = map[string]string{
		"HomeDir" :           "string",
		"ConfigDir" :         "string",
		"CertFile" :          "string",
		"LogDir" :            "string",
		"TmpDir" :            "string",
		"TagImageDir" :       "string",
		"TemplateDir" :       "string",
		"FontDir" :           "string",
		"ImageDir" :          "string",
		"ExternalBinaryDir" : "string",
		"DataBackupDir" :     "string",
		"LogLevel" :                  "string",
		"ApBrokerRetryTimingSecond" : "int",
		"ImgGenThreadCount" :         "int",
		"ImgGenReqPort" :             "string",
		"ImgGenRespPort" :            "string",
		"ImgGenPubPort" :             "string",
		"EslApTimerReqPort" :         "string",
		"DeassignThreadCount" :       "int",
		"FontFacePreloadCount" :      "int",
		"LastTaskIdBackupFile" :  "string",
		"ProductDataBackupFile" : "string",
		"AssignDataBackupFile" :  "string",
		"NfcDataBackupFile" :     "string",
		"EventFrameTxTiming" :    "int",
		"TagImageTxTimging" :     "int",
		"ScanProfile" :           "[6]int",
		"BackoffBase" :           "int",
		"BackoffMulFactor" :      "int",
		"FreezerTagMultiplier" :  "int",
		"TagDistributionMinute" : "int",
		"PageRotationMacPage" :   "int",
		"LogMaxSizeMb" :  "int",
		"LogMaxBackup" :  "int",
		"LogMaxAgeDays" : "int",
		"LogCompress" :   "bool",
		"GRpcMaxSize" :   "int",
	}
)

type OperationData struct {
	Data  map[string]interface{}
	Types map[string]string
}

func readNetworkConfig(ctx *gin.Context) {
	//parse YAML from yaml directory
	parsed_network_env, err := config.GetConfigEnv(NETWORK_ENV, &network_env)
	if err != nil {
		log.Println("Error:", err)
		log.Println("readNetworkConfig")
		os.Exit(1)
	}

	//type assertion to networkenv
	network_env_struct := parsed_network_env.(*config.NetworkEnv)

	//change interface to map type to pass to html
	network_env_map := structs.Map(network_env_struct)
	ctx.HTML(http.StatusOK, "network.html", gin.H{
		"YAMLData" : network_env_map,
	})
}

func updateNetworkConfig(ctx *gin.Context) {
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
		log.Println("updateNetworkConfigAtoi")
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
		NameServers : strings.Split(nameServerRaw, " "),       
		TimeZone : formFields.Get("TimeZone"),          
		TimeServerUrls : strings.Split(timeServerRaw, " "),    
		InterApPort : formFields.Get("InterApPort"),       
		InterApPortTarget : formFields.Get("InterApPortTarget"), 
		ApBrokerUrl : formFields.Get("ApBrokerUrl"),       
		EthernetInterface : formFields.Get("EthernetInterface"),
	}
	err = config.SetConfigEnv(NETWORK_ENV, &new_network_struct)

	if err != nil {
		log.Println("Error:", err)
		log.Println("updateNetworkConfig  SetConfigEnv")
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
	ctx.Redirect(http.StatusSeeOther, "/network.html")
}

func loadHomePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func readOperationConfig(ctx *gin.Context) {
	//parse YAML from yaml directory
	parsed_operation_env, err := config.GetConfigEnv(OPERATION_ENV, &operation_env)
	if err != nil {
		log.Println("Error:", err)
		log.Println("readOperationConfig")
		os.Exit(1)
	}

	//type assertion to operationenv
	operation_env_struct := parsed_operation_env.(*config.OperationEnv)

	//change interface to map type to pass to html
	operation_env_map := structs.Map(operation_env_struct)
	
	ctx.HTML(http.StatusOK, "operation.html", gin.H{
		"YAMLData" : OperationData{
				Data : operation_env_map,
				Types: operationFieldTypes,
		},
	})
}

func updateOperationConfig(ctx *gin.Context) {
	return
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
	//initial load
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusSeeOther, "/index.html")
	})

	//homepage
	router.GET("/index.html", loadHomePage)

	// display network YAML data in webpage
	router.GET("/network.html", readNetworkConfig)

	// write to network YAML when form is submitted, update webpage as well
	router.POST("/network.html", updateNetworkConfig)

	// diplay operation YAML data in webpage
	router.GET("/operation.html", readOperationConfig)

	// write to operation YAML when form is submitted, update webpage as well
	router.POST("/operation.html", updateOperationConfig)
	
	router.Run(":8080")

}