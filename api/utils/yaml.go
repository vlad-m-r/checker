package utils

import (
	"io/ioutil"
	"log"
)

func ReadYaml(yamlFile string) []byte {
	yamlContent, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	return yamlContent
}
