package mqtt_go

import (
	"encoding/json"
	"fmt"
	"iot-project/db/model"
	"iot-project/db/repository"
	"iot-project/db/storage"
	"log"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	type Sensor struct {
		Flight      bool    `json:"flight"`
		Temperature float64 `json:"temperature"`
		Humidity    float64 `json:"humidity"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	var sensorData Sensor
	err := json.Unmarshal(msg.Payload(), &sensorData)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	psqlDB, err := storage.NewPSQLManager()
	if err != nil {
		log.Fatalf("Error when connecting database, err: %v", err)
	}
	defer psqlDB.Close()
	sensorRepository := repository.NewSensorRepository(psqlDB)
	dataDB := model.Sensor{
		Flight:      sensorData.Flight,
		Temperature: sensorData.Temperature,
		Humidity:    sensorData.Humidity,
		CreatedAt:   sensorData.CreatedAt,
		UpdatedAt:   sensorData.UpdatedAt,
	}
	err = sensorRepository.InsertSensor(&dataDB)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
	}
}
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
var (
	broker   = "broker.emqx.io"
	port     = 1883
	username = "emqx"
	clientId = "go_mqtt_client"
	password = "public"
	topic    = "tuan"
)

func ConnectMQTTBroker(opts *mqtt.ClientOptions) (*mqtt.Client, error) {
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &client, nil
}
func Sub(client mqtt.Client) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}
func Publish(client mqtt.Client) {
	num := 2
	layout := "2006-01-02 15:04"
	for i := 0; i < num; i++ {
		timeSend := time.Now().Format("2006-01-02 15:04")
		t, err := time.Parse(layout, timeSend)

		if err != nil {
			fmt.Println(err)
		}
		sensor := model.Sensor{
			Flight:      false,
			Temperature: -10.0 + rand.Float64()*(100.0-(-10.0)),
			Humidity:    0.0 + rand.Float64()*(100.0-0),
			CreatedAt:   t,
			UpdatedAt:   t,
		}
		messageJSON, err := json.Marshal(sensor)
		if err != nil {
			log.Printf("Error marshal sensor:%v\n", err)
			os.Exit(1)
		}
		token := client.Publish(topic, 0, false, messageJSON)
		token.Wait()
		time.Sleep(time.Second * 10)
	}
}
