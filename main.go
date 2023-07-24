package main

import (
	"ap_config_web_ui/config"
	"net/http"
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
	MAX_INPUT_LENGTH = 2 << 14
	YAML_FOLDER = "config_files/"
	NETWORK_ENV = YAML_FOLDER + "network.yaml"
	OPERATION_ENV = YAML_FOLDER + "operation.yaml"
	VALIDATE_YAML_CHANGES = false //flag to validate YAML changes using gRPC
	GRPC_SUCCESS_TOKEN_NETWORK = 1
	GRPC_SUCCESS_TOKEN_OPERATION = 2
	GRPC_FAIL_TOKEN = 0
)

var (
	networkEnv config.NetworkEnv
	operationEnv config.OperationEnv
	client validate.ValidateClient
	errorMessage string

)

func addQuotationLiteral(s string) string {
	return strings.ReplaceAll(s, "\"", "&quot;")
}

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
		log.Println("readNetworkConfig Error:", err)
		errorMessage = "Failed to retrieve data.\nPlease check to see if the file matches the required specifications and try again later."
		ctx.Redirect(http.StatusFound, "/error")
		return
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
		log.Println("updateNetworkConfig Error:", err)
		errorMessage = "Failed to submit data.\nPlease try again later."
		ctx.Redirect(http.StatusFound, "/error")
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
				errorMessage = "Failed to handle submitted data.\nPlease try again later."
				ctx.Redirect(http.StatusFound, "/error")
				return
			}
			value.SetUint(uint64(num))
		}
	}

	err = config.SetConfigEnv(NETWORK_ENV, &newNetworkStruct)

	if err != nil {
		log.Println("updateNetworkConfig  SetConfigEnv Error:", err)
		errorMessage = "Failed to submit data.\nPlease try again later."
		ctx.Redirect(http.StatusFound, "/error")
		return
	}
	if VALIDATE_YAML_CHANGES {
		verResponse, err := client.Verify(context.Background(), &validate.VerifyRequest{Token: GRPC_SUCCESS_TOKEN_NETWORK})
		if err != nil {
			log.Printf("updateNetworkConfig Verify failed: %v", err)
		}
		if verResponse.IsValid {
			log.Printf("updateNetworkConfig Verify result: Valid")
		} else {
			log.Printf("updateNetworkConfig Verify result: Invalid")
		}
	}
	ctx.Redirect(http.StatusSeeOther, "/network")
}

func readOperationConfig(ctx *gin.Context) {
	//parse YAML from yaml directory
	parsedOperationEnv, err := config.GetConfigEnv(OPERATION_ENV, &operationEnv)
	if err != nil {
		log.Println("readOperationConfig Error:", err)
		errorMessage = "Failed to retrieve data.\nPlease check to see if the file matches the required specifications and try again later."
		ctx.Redirect(http.StatusFound, "/error")
		return
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
		errorMessage = "Failed to parse form.\nPlease try again later."
		ctx.Redirect(http.StatusFound, "/error")
		return
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
				log.Printf("updateOperationConfig Atoi Err:%v", err)
				errorMessage = "Failed to handle submitted data.\nPlease try again later."
		 		ctx.Redirect(http.StatusFound, "/error")	
				return
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
				log.Printf("updateOperationConfig Atoi Err:%v", err)
				errorMessage = "Failed to handle submitted data.\nPlease try again later."
				ctx.Redirect(http.StatusFound, "/error")
				return
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

	err = config.SetConfigEnv(OPERATION_ENV, &newOperationStruct)

	if err != nil {
		log.Println("updateOperationConfig SetConfig Env Error:", err)
		errorMessage = "Failed to submit data.\nPlease try again later."
		ctx.Redirect(http.StatusFound, "/error")
		return
	}
	if VALIDATE_YAML_CHANGES {
		verResponse, err := client.Verify(context.Background(), &validate.VerifyRequest{Token: GRPC_SUCCESS_TOKEN_OPERATION})
		if err != nil {
			log.Printf("updateOperationConfig Verify failed: %v", err)
		}
		if verResponse.IsValid {
			log.Printf("updateOperationConfig Verify result: Valid")
		} else {
			log.Printf("updateOperationConfig Verify result: Invalid")
		}
	}
	ctx.Redirect(http.StatusSeeOther, "/operation")
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
			v := addQuotationLiteral(value.String())
			formHTML += fmt.Sprintf("<input class = \"form-control\" type=\"text\" name=\"%s\" value=\"%s\" maxlength=\"%v\" required>\n", fieldName, v, MAX_INPUT_LENGTH)
		case reflect.Uint32:
			formHTML += fmt.Sprintf("<input class = \"form-control\" type=\"number\" name=\"%s\" value=\"%d\" min=\"0\" max=\"4294967295\" step=\"1\" required>\n", fieldName, value.Uint())
		case reflect.Int:
			formHTML += fmt.Sprintf("<input class = \"form-control\" type=\"number\" name=\"%s\" value=\"%d\" min=\"-2147483648\" max=\"2147483647\" step=\"1\" required>\n", fieldName, value.Int())
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
			formHTML += fmt.Sprintf("<input class = \"form-control\"type = \"text\" name = \"%s\" id = \"%s\" value = \"%v\" maxlength=\"%v\" pattern = \"%s\" title = \"Please enter a space seperated int list of length 6, example: [4 3 2 6 0 1]\"  required>", 
			fieldName, fieldName, value, MAX_INPUT_LENGTH, regex)
		case reflect.Slice:
			regex := regexp.MustCompile(`^\[\s*[\w\d.]+(?:\s+[\w\d.]+)*\s*\]$`)
			formHTML += fmt.Sprintf("<input class = \"form-control\"type = \"text\" name = \"%s\" id = \"%s\" value = \"%v\" maxlength=\"%v\" pattern = \"%s\" title = \"Please enter a space seperated list, example: [4 2 0]\"  required>", 
			fieldName, fieldName, value, MAX_INPUT_LENGTH, regex)
		}
		formHTML += "</div>"
	}
	return template.HTML(formHTML)

}

func handleErrors(ctx *gin.Context) {
	log.Println(errorMessage)
	ctx.HTML(http.StatusOK, "error.html", gin.H{"errorMessage" : errorMessage,})
}

func handle404(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "404.html", nil)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	if VALIDATE_YAML_CHANGES {
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to connect: %v", err)
		}
		defer conn.Close()
		client = validate.NewValidateClient(conn)
	}
	
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"generateHTMLForm": generateHTMLForm,
	})
	router.LoadHTMLGlob("templates/*")


	//custom middleware to handle error 500
	router.Use(func (ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Header("Content-Type", "text/html; charset=utf-8")
				ctx.HTML(http.StatusInternalServerError, "500.html", nil)
				ctx.Abort()
			}
		}()
		ctx.Next()
	})

	//initial load
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusSeeOther, "/network")
	})

	// display network YAML data in webpage
	router.GET("/network", readNetworkConfig)

	// write to network YAML when form is submitted, update webpage as well
	router.POST("/network", updateNetworkConfig)

	// diplay operation YAML data in webpage
	router.GET("/operation", readOperationConfig)

	// write to operation YAML when form is submitted, update webpage as well
	router.POST("/operation", updateOperationConfig)
	
	// error handling
	router.GET("/error", handleErrors)

	// handle error 404
	router.NoRoute(handle404)
	
	router.Run(":8080")

}