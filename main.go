package main

import (
	"time"

	// "iot-project/db/model"
	"iot-project/db/model"
	"iot-project/db/repository"
	"iot-project/db/storage"
	"iot-project/mqtt-go"
	"log"
	"net/http"

	// "time"

	rice "github.com/GeertJohan/go.rice"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	gintemplate "github.com/foolin/gin-template"
	"github.com/foolin/gin-template/supports/gorice"
	"github.com/gin-gonic/gin"
	"gopkg.in/robfig/cron.v2"
)

func main() {
	//Connect mqtt broker and execute topic subscribe
	mess := make(chan string)
	c := cron.New()
	m := cron.New()
	opts := mqtt.NewClientOptions()
	client, err := mqtt_go.ConnectMQTTBroker(opts)
	if err != nil {
		log.Printf("Error connect broker: %v\n" ,err)
	}
	go func() {
		mqtt_go.Sub(*client)
		mqtt_go.Publish(c, *client)
		(*client).Disconnect(2500)
		c.Stop()
	}()
	//Connect MYSQL/PostgresSQL , insert&& read record
	psqlDB, err := storage.NewPSQLManager()
	if err != nil {
		log.Fatalf("Error when connecting database, err: %v", err)
	}
	defer psqlDB.Close()
	var data *model.Sensor
	sensorRepository := repository.NewSensorRepository(psqlDB)
	time.Sleep(time.Second * 10)
	m.AddFunc("@every 0h0m5s", func() {
		then := time.Now().Add(time.Duration(-10) * time.Second).Format("2006-01-02 15:04")
		currentTime, _ := time.Parse("2006-01-02 15:04", then)
		data, err = sensorRepository.GetSensorByTime(currentTime)
		if err != nil {
			log.Fatalf("Error when get sensor database, err: %v", err)
		}
	})
	m.Start()
		// fmt.Println(data)
		router := gin.Default()
		router.HTMLRender = gintemplate.Default()
		staticBox := rice.MustFindBox("static")
		router.StaticFS("/static", staticBox.HTTPBox())
		router.HTMLRender = gorice.New(rice.MustFindBox("views"))
		n := cron.New()
		n.AddFunc("@every 0h0m30s", func() {
			router.GET("/", func(c *gin.Context){
				c.HTML(http.StatusOK, "index", gin.H{
					"title": "Smart Home",
					"answer1": data.Flight,
					"answer2": data.Temperature,
					"answer3": data.Humidity,
					"answer4": data.UpdatedAt.Format("2006-01-02 15:04"),
				})
			})
		})
		n.Start()
		router.Run(":9090")
		<-mess
		m.Stop()
		n.Stop()
}

