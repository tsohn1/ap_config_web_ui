package main

import (
	"fmt"
	"os"
	"ap_config_web_ui/config"

)

var network_env config.NetworkEnv

func main() {
	parsed_network_env, err := config.GetConfigEnv("config_files/network.yaml", &network_env)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", parsed_network_env)

}