package publishlamp

import (
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"os"
)

type Bulbs struct {
	IDs []string `json:"Bulbs"`
}

func ShareSecrets(jsonname string) Bulbs {
	jsonFile, err := os.Open(jsonname)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var bulbs Bulbs

	json.Unmarshal(byteValue, &bulbs)

	return bulbs
}

func QhtekHost() string {
	return "tcp://cloud.qh-tek.com:1883"
}

//output byte values in Dim(), On(), Off() are captured from the official smartphone app
func Dim(value byte) []byte {
	// values outside 1 to 100 range softlock lightbulbs
	if value > 100 {
		value = 100
	}
	if value < 1 {
		value = 1
	}

	output := []byte{
		0x21,
		0x01,
		value,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		value ^ 1, // Value XOR 1 (checksum?) is found in this place in the app
		0x3a,
	}

	return output
}

func On() []byte {
	output := []byte{
		0xfa,
		0x23,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x23,
		0xfb,
	}

	return output
}

func Off() []byte {
	output := []byte{
		0xfa,
		0x24,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x24,
		0xfb,
	}

	return output
}

func Publish(host string, bulb_ids []string, message []byte) {
	clientOpts := MQTT.NewClientOptions().AddBroker(host)
	client := MQTT.NewClient(clientOpts)
	c_token := client.Connect()
	if c_token.Wait() && c_token.Error() != nil {
		log.Fatal(c_token.Error())
	}

	for _, value := range bulb_ids {
		p_token := client.Publish(value, 0, false, message)
		if p_token.Wait() && p_token.Error() != nil {
			log.Fatal(p_token.Error())
		}
	}
}
