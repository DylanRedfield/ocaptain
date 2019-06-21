package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func GetEnvValues() EnvValues {
	jsonFile, err := os.Open("../env_values.json")

	if err != nil {
		log.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var envValues EnvValues

	if err := json.Unmarshal([]byte(byteValue), &envValues); err != nil {
		log.Print(err)
	}

	return envValues
}
