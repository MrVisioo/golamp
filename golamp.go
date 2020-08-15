package main

import (
	"fmt"
	lamp "github.com/mrvisioo/golamp/publishlamp"
	"os"
	"strconv"
)

func usage(err int) {
	fmt.Fprintf(os.Stderr, "usage: %s [command] [value] \n", os.Args[0])
	fmt.Fprintf(os.Stderr, "commands:\n")
	fmt.Fprintf(os.Stderr, "  on         	Turn on\n")
	fmt.Fprintf(os.Stderr, "  off        	Turn off\n")
	fmt.Fprintf(os.Stderr, "  dim [value]	Dim TO value. Default value is 100. Negative values subtract from 100 (-80 is 20). Turns on lightbulb if it's off\n")
	os.Exit(err)
}

func main() {

	bulbs := lamp.ShareSecrets("secret.json")

	if len(os.Args) > 1 {
		switch arg := os.Args[1]; arg {
		case "on":
			lamp.Publish(lamp.QhtekHost(), bulbs.IDs, lamp.On())
		case "off":
			lamp.Publish(lamp.QhtekHost(), bulbs.IDs, lamp.Off())
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
					lamp.Publish(lamp.QhtekHost(), bulbs.IDs, lamp.Dim(byte(bright)))
				} else {
					lamp.Publish(lamp.QhtekHost(), bulbs.IDs, lamp.Dim(100))
				}
			}
		case "help":
			usage(0)
		default:
			usage(1)
		}
	} else {
		usage(2)
	}
}
