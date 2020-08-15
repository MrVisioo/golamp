package publishlamp

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os/exec"
)

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

func PublishCmd(host string, bulb_id string, message string) {
	cmd_pub := exec.Command("mosquitto_pub", "-h", host, "-t", bulb_id, "-f", message)
	log.Printf(cmd_pub.String())
	err := cmd_pub.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting to publish on %s", bulb_id)
	err = cmd_pub.Wait()
	log.Printf("Publish on %s is success", bulb_id)
}

func PublishmanyCmd(host string, bulb_ids []string, message string) {
	cmds := make([]*exec.Cmd, len(bulb_ids))
	errs := make([]error, len(bulb_ids))
	for i, value := range bulb_ids {
		cmds[i] = exec.Command("mosquitto_pub", "-h", host, "-t", value, "-f", message)
		errs[i] = cmds[i].Start()
		log.Printf("Publishing on %s, %v", value, errs[i])
		if errs[i] != nil {
			log.Fatal(errs[i])
		}
	}
	for i, value := range bulb_ids {
		errs[i] = cmds[i].Wait()
		log.Printf("Publish on %s is success, %v", value, errs[i])
	}
}
