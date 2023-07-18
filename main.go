package main

import (
	"ap_config_web_ui/config"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
	"reflect"
	"fmt"
	"html/template"
	"regexp"
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
	GRPC_SUCCESS_TOKEN_NETWORK = 1
	GRPC_SUCCESS_TOKEN_OPERATION = 2
	GRPC_FAIL_TOKEN = 0
)

var (
	networkEnv config.NetworkEnv
	operationEnv config.OperationEnv
	client validate.ValidateClient

)


func removeWhiteSpace(Arr []string) []string {
	result := make([]string, 0)
	for i, val := range Arr {
		if val != "" {
			result = append(result, Arr[i])
		}
	}
	return result
}

func readNetworkConfig(ctx *gin.Context) {
	//parse YAML from yaml directory
	parsedNetworkEnv, err := config.GetConfigEnv(NETWORK_ENV, &networkEnv)
	if err != nil {
		log.Println("Error:", err)
		log.Println("readNetworkConfig")
		os.Exit(1)
	}

	//type assertion to networkenv
	networkEnvStruct := parsedNetworkEnv.(*config.NetworkEnv)

	ctx.HTML(http.StatusOK, "network.html", gin.H{
		"YAMLData" : *networkEnvStruct,
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

	//convert NameServers, TimeserverUrls into []string
	nameServerRaw := formFields.Get("NameServers")
	timeServerRaw := formFields.Get("TimeServerUrls")
	nameServerRaw = strings.Trim(nameServerRaw, "[]")
	timeServerRaw = strings.Trim(timeServerRaw, "[]")

	newNetworkStruct := config.NetworkEnv{      
		NameServers : removeWhiteSpace(strings.Split(nameServerRaw, " ")),                
		TimeServerUrls : removeWhiteSpace(strings.Split(timeServerRaw, " ")),    
	}

	netStructReflect := reflect.ValueOf(&newNetworkStruct).Elem()

	//loop through struct fields and retreive values from form
	for i := 0; i < netStructReflect.NumField(); i++ {
		field := netStructReflect.Type().Field(i)
		value := netStructReflect.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			value.SetString(formFields.Get(field.Name))
		case reflect.Uint32:
			num, err := strconv.Atoi(formFields.Get(field.Name))
			if err != nil {
				log.Println("updateNetworkConfigAtoi: during loop", err)
				os.Exit(1)
			}
			value.SetUint(uint64(num))
		}
	}

	err = config.SetConfigEnv(NETWORK_ENV, &newNetworkStruct)

	if err != nil {
		log.Println("Error:", err)
		log.Println("updateNetworkConfig  SetConfigEnv")
		os.Exit(1)
	}
	if VALIDATE_YAML_CHANGES {
		verResponse, err := client.Verify(context.Background(), &validate.VerifyRequest{Token: GRPC_SUCCESS_TOKEN_NETWORK})
		if err != nil {
			log.Fatalf("updateNetworkConfig Verify failed: %v", err)
		}
		if verResponse.IsValid {
			log.Printf("updateNetworkConfig Verify result: Valid")
		} else {
			log.Printf("updateNetworkConfig Verify result: Invalid")
		}
	}
	ctx.Redirect(http.StatusSeeOther, "/network.html")
}

func loadHomePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func readOperationConfig(ctx *gin.Context) {
	//parse YAML from yaml directory
	parsedOperationEnv, err := config.GetConfigEnv(OPERATION_ENV, &operationEnv)
	if err != nil {
		log.Println("Error:", err)
		log.Println("readOperationConfig")
		os.Exit(1)
	}

	//type assertion to operationenv
	operationEnvStruct := parsedOperationEnv.(*config.OperationEnv)
	
	ctx.HTML(http.StatusOK, "operation.html", gin.H{
		"YAMLData" : *operationEnvStruct,
	})
}

func updateOperationConfig(ctx *gin.Context) {
	err := ctx.Request.ParseForm()
	if err != nil {
		log.Println("Error: ctx.Request.ParseForm()", err)
		// Handle the error, possibly by returning an error response
		os.Exit(1)
	}

	//Retrieve all form fields
	formFields := ctx.Request.PostForm
	
	//convert ScanProfile into [6]int
	ScanProfileRaw := formFields.Get("ScanProfile")
	ScanProfileRaw = strings.Trim(ScanProfileRaw, "[]")
	ScanProfileArr := strings.Split(ScanProfileRaw, " ")
	ScanProfileVal := make([]int, len(ScanProfileArr))

	
	for i, val := range ScanProfileArr {
		if val != "" {
			ScanProfileVal[i], err = strconv.Atoi(val)
			if err != nil {
				log.Println("updateNetworkConfigAtoi:", err)
				os.Exit(1)
			}
		}
	}
	newOperationStruct := config.OperationEnv{}
	opStructReflect := reflect.ValueOf(&newOperationStruct).Elem()

	//loop through struct fields and retreive values from form
	for i := 0; i < opStructReflect.NumField(); i++ {
		field := opStructReflect.Type().Field(i)
		value := opStructReflect.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			value.SetString(formFields.Get(field.Name))
		case reflect.Int:
			num, err := strconv.Atoi(formFields.Get(field.Name))
			if err != nil {
				log.Println("updateNetworkConfigAtoi: during loop", err)
				os.Exit(1)
			}
			value.SetInt(int64(num))
		case reflect.Bool:
			if formFields.Get(field.Name) == "" {
				value.SetBool(false)
			} else {
				value.SetBool(true)
			}
		}
	}
	newScanProfile := [6]int{}
	copy(newScanProfile[:], ScanProfileVal)
	newOperationStruct.ScanProfile = newScanProfile

	if err != nil {
		log.Println("Error: readOperationConfig, ", err)
		os.Exit(1)
	}

	err = config.SetConfigEnv(OPERATION_ENV, &newOperationStruct)

	if err != nil {
		log.Println("Error:", err)
		log.Println("updateOperationConfig  SetConfigEnv")
		os.Exit(1)
	}
	if VALIDATE_YAML_CHANGES {
		verResponse, err := client.Verify(context.Background(), &validate.VerifyRequest{Token: GRPC_SUCCESS_TOKEN_OPERATION})
		if err != nil {
			log.Fatalf("updateOperationConfig Verify failed: %v", err)
		}
		if verResponse.IsValid {
			log.Printf("updateOperationConfig Verify result: Valid")
		} else {
			log.Printf("updateOperationConfig Verify result: Invalid")
		}
	}
	ctx.Redirect(http.StatusSeeOther, "/operation.html")
}

func generateHTMLForm(data interface{}) template.HTML {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	
	formHTML := ""

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fieldName := field.Name

		if !value.CanInterface() {
			continue
		}
		formHTML += fmt.Sprintf("<div class=\"col-6 container\">\n<label class= \"form-label\" for = \"%s\">%s</label>\n", fieldName, fieldName)
		switch value.Kind() {
		case reflect.String:
			formHTML += fmt.Sprintf("<input class = \"form-control\" type=\"text\" name=\"%s\" value=\"%s\">\n", fieldName, value.String())
		case reflect.Uint32:
			formHTML += fmt.Sprintf("<input class = \"form-control\" type=\"number\" name=\"%s\" value=\"%d\" min=\"0\" max=\"4294967295\" step=\"1\">\n", fieldName, value.Uint())
		case reflect.Int:
			formHTML += fmt.Sprintf("<input class = \"form-control\" type=\"number\" name=\"%s\" value=\"%d\" min=\"-2147483648\" max=\"2147483647\" step=\"1\">\n", fieldName, value.Int())
		case reflect.Bool:
			formHTML += "<div class = \"row\">\n<div class = \"col d-flex align-items-center\">\n<div class = \"form-check form-switch\">\n"
			if value.Bool(){
				formHTML += fmt.Sprintf("<input class=\"form-check-input\" name=\"%s\" type=\"checkbox\" role=\"switch\" id=\"flexSwitchCheckDefault\" checked>\n<label class=\"form-check-label\" for=\"%s\">On: True, Off: False\n</label>", fieldName, fieldName)
			} else{
				formHTML += fmt.Sprintf("<input class=\"form-check-input\" name=\"%s\" type=\"checkbox\" role=\"switch\" id=\"flexSwitchCheckDefault\">\n<label class=\"form-check-label\" for=\"%s\">On: True, Off: False\n</label>", fieldName, fieldName)
			}
			formHTML += "</div>\n</div>\n</div>\n"
		case reflect.Array: //Special case for [6]int array
			regex := regexp.MustCompile(`^\[\d+(?:\s+\d+){5}\]$`)
			formHTML += fmt.Sprintf("<input class = \"form-control\"type = \"text\" name = \"%s\" id = \"%s\" value = \"%v\" pattern = \"%s\" title = \"Please enter a space seperated int list of length 6, example: [4 3 2 6 0 1]\"  required>", 
			fieldName, fieldName, value, regex)
		case reflect.Slice:
			regex := regexp.MustCompile(`^\[\s*[\w\d.]+(?:\s+[\w\d.]+)*\s*\]$`)
			formHTML += fmt.Sprintf("<input class = \"form-control\"type = \"text\" name = \"%s\" id = \"%s\" value = \"%v\" pattern = \"%s\" title = \"Please enter a space seperated list, example: [4 2 0]\"  required>", 
			fieldName, fieldName, value, regex)
		}
		formHTML += "</div>"
	}
	return template.HTML(formHTML)

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
	router.SetFuncMap(template.FuncMap{
		"generateHTMLForm": generateHTMLForm,
	})
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