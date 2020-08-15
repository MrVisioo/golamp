package main

import (
	"encoding/json"
	lamp "github.com/mrvisioo/golamp/publishlamp"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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

	if len(os.Args) > 1 {
		switch arg := os.Args[1]; arg {
		case "on":
			lamp.Publish(host, bulbs, lamp.On())
		case "off":
			lamp.Publish(host, bulbs, lamp.Off())
		case "dim":
			{
				if len(os.Args) > 2 {
					bright, _ := strconv.ParseInt(os.Args[2], 10, 16)
					if bright < 0 {
						if bright > -100 {
							bright = 100 + bright
						} else {
							bright = 1
						}
					}
					lamp.Publish(host, bulbs, lamp.Dim(byte(bright)))
				} else {
					lamp.Publish(host, bulbs, lamp.Dim(100))
				}
			}
		default:
			log.Fatal("Not enough arguments")
		}
	} else {
		log.Fatal("Not enough arguments")
	}
}
