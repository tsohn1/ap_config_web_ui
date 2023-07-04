package config

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var cfgLogger *zap.SugaredLogger

var ConfigOperationEnv *OperationEnv
var ConfigDatabaseEnv *DatabaseEnv
var ConfigNetworkEnv *NetworkEnv
var ConfigGrmEnv *GrmEnv

func GetConfigEnv(configFile string, configStruct interface{}) (interface{}, error) {
	buf, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	
	switch v := configStruct.(type) {
	case *DatabaseEnv:
		err = yaml.Unmarshal(buf, v)
		ConfigDatabaseEnv = v
		return v, err
	case *OperationEnv:
		err = yaml.Unmarshal(buf, v)
		ConfigOperationEnv = v
		return v, err
	case *GrmEnv:
		err = yaml.Unmarshal(buf, v)
		ConfigGrmEnv = v
		return v, err
	case *NetworkEnv:
		err = yaml.Unmarshal(buf, v)
		ConfigNetworkEnv = v
		return v, err
	case *LastTaskIdBackup:
		err = yaml.Unmarshal(buf, v)
		return v, err
	}

	return nil, errors.New("configStruct type is not valid")
}

func SetConfigEnv(configFile string, configStruct interface{}) error {
	data, err := yaml.Marshal(configStruct)
	if err != nil {
		ec := fmt.Sprintf("SetConfigEnv: yaml.Marshal failed: %v, err:%v", configStruct, err)
		cfgLogger.Error(ec)
		return errors.New(ec)
	}
	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		ec := fmt.Sprintf("SetConfigEnv: os.WriteFile failed: %v<-%v, err:%v", configFile, data, err)
		cfgLogger.Error(ec)
		return errors.New(ec)
	}
	switch v := configStruct.(type) {
	case *DatabaseEnv:
		ConfigDatabaseEnv = v
	case *OperationEnv:
		ConfigOperationEnv = v
	case *GrmEnv:
		ConfigGrmEnv = v
	case *NetworkEnv:
		ConfigNetworkEnv = v
	}
	return nil
}

func SetDefaultLogger(logger *zap.SugaredLogger) {
	cfgLogger = logger
}

func GetDefaultLogger() *zap.SugaredLogger {
	return cfgLogger
}
