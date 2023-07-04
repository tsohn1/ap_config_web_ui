package main

import (
	"fmt"
	"os"
	"ap_config_web_ui/config"

)

var network_env *config.NetworkEnv

func main() {
	network_env, err := config.GetConfigEnv("config_files/network.yaml", config.ConfigNetworkEnv)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	
	fmt.Printf("%+v\n", network_env)

}