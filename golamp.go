package main

import (
	"encoding/json"
	"fmt"
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

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [command] [value] \n", os.Args[0])
	fmt.Fprintf(os.Stderr, "commands:\n")
	fmt.Fprintf(os.Stderr, "  on	Turn on\n")
	fmt.Fprintf(os.Stderr, "  off	Turn off\n")
	fmt.Fprintf(os.Stderr, "  dim [value]	Dim TO value. Default value is 100. Negative values subtract from 100 (-80 is 20).\n")
	os.Exit()
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
			usage()
		}
	} else {
		usage()
	}
}
