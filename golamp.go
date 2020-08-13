package main

import (
	"encoding/json"
	"github.com/mrvisioo/golamp/publishlamp"
	"io/ioutil"
	"log"
	"os"
)

type Bulbs struct {
	Bulbs []string `json:"Bulbs"`
}

func ShareSecrets() []string {
	jsonFile, err := os.Open("secret.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var bulbs Bulbs

	json.Unmarshal(byteValue, &bulbs)

	return bulbs.Bulbs
}

func main() {

	bulbs := ShareSecrets()

	host := "tcp://cloud.qh-tek.com:1883"

	off := "\xfa\x24\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x24\xfb"
	on := "\xfa\x23\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x23\xfb"

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "on" {
			publishlamp.PublishPaho(host, bulbs, on)
		} else if arg == "off" {
			publishlamp.PublishPaho(host, bulbs, off)
		}
	} else {
		log.Fatal("Not enough arguments")
	}
}
