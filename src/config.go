package main

import (
	"encoding/json"
	"os"
)

var configRootDir, _ = os.UserConfigDir()
var configDirPath = configRootDir + "/.slatusify/"
var configPath = configDirPath + "config.json"

type Config struct {
	OauthToken string
}

var config Config = Config{}

func (self Config) store() {
	dirInfo, _ := os.Stat(configDirPath)
	if !dirInfo.IsDir() {
		dirErr := os.Mkdir(configDirPath, 0770)
		if dirErr != nil {
			logger.Fatalln("Cannot make config dir:", dirErr)
		}
	}
	contents, jsonErr := json.Marshal(self)
	if jsonErr == nil {
		configFile, ioErr := os.OpenFile(configPath, os.O_CREATE|os.O_RDWR, os.FileMode(0770))
		if ioErr == nil {
			byteCount, writeErr := configFile.Write(contents)
			if writeErr != nil {
				logger.Fatalln("Config file write failure:", writeErr, "bytes written:", byteCount)
			}
			configFile.Close()
		} else {
			logger.Fatalln("Cannot open config file:", ioErr)
		}
	} else {
		logger.Fatalln("Config failed to marshal:", jsonErr)
	}
}
