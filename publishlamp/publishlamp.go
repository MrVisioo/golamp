package publishlamp

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os/exec"
)

func PublishPaho(host string, bulb_ids []string, message string) {
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
