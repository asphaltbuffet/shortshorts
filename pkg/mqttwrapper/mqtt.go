// Package mqttwrapper provides a wrapper around the Eclipse Paho MQTT client.
package mqttwrapper

import (
	"context"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Start starts the MQTT client. It should be called on the service start up.
func Start(ctx context.Context, datastream chan [2]string) (mqtt.Client, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := cfg.Build()

	topic := "Hotpants/data"                  // The topic name to/from which to publish/subscribe
	broker := "tcp://test.mosquitto.org:1883" // The broker URI. ex: tcp://10.10.1.1:1883
	// password := "" // The password (optional)
	// user := "" // The User (optional)
	id := "shortshorts" // The ClientID (optional)
	cleanSess := false  // Set Clean Session (default false)
	qos := 0            // The Quality of Service 0,1,2 (default 0)
	store := ":memory:" // "The Store Directory (default use memory store)")

	if topic == "" {
		logger.Error("invalid topic", zap.String("topic", topic))
		return nil, fmt.Errorf("invalid topic: %s", topic)
	}

	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID(id)
	options.SetCleanSession(cleanSess)

	if store != ":memory:" {
		options.SetStore(mqtt.NewFileStore(store))
	}

	options.SetOnConnectHandler(func(client mqtt.Client) {
		logger.Info("Connected to broker", zap.String("broker", broker))
	})

	options.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logger.Error("connection lost", zap.Error(err))
	})

	cli := mqtt.NewClient(options)
	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		logger.Error("error connecting to broker", zap.Error(token.Error()))
		return nil, token.Error()
	}

	f := func(client mqtt.Client, msg mqtt.Message) {
		datastream <- [2]string{msg.Topic(), string(msg.Payload())}
	}

	if token := cli.Subscribe(topic, byte(qos), f); token.Wait() && token.Error() != nil {
		logger.Error("error subscribing to topic", zap.Error(token.Error()))
		return nil, token.Error()
	}

	logger.Info("successfully subscribed", zap.String("topic", topic), zap.String("broker", broker))

	return cli, nil
}
