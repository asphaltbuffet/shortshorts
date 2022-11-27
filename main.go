// Package main is the main package for shortshorts.
package main

import (
	"context"
	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/asphaltbuffet/shortshorts/pkg/logging"
	"github.com/asphaltbuffet/shortshorts/pkg/mqttwrapper"
	"github.com/asphaltbuffet/shortshorts/pkg/servicemanager"
	"github.com/asphaltbuffet/shortshorts/pkg/timescalewrapper"
)

const (
	// MQTTDisconnectTimeout is the timeout for disconnecting from the MQTT broker.
	MQTTDisconnectTimeout uint = 100
)

var logger *zap.Logger

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	tsdb, dstream, client := start(ctx)

	go processLoop(ctx, dstream, tsdb)

	servicemanager.WaitShutdown(func() { shutdown(cancel, tsdb, client) })
}

func start(ctx context.Context) (*timescalewrapper.Database, chan [2]string, mqtt.Client) {
	if err := logging.Start(); err != nil {
		panic(err)
	}

	logger = logging.GetLogger()

	connStr := readConfig()

	tsdb, err := timescalewrapper.NewDatabase(ctx, connStr)
	if err != nil {
		logger.Panic("error connecting to timescale", zap.Error(err))
	}

	dstream := make(chan [2]string)

	c, err := mqttwrapper.Start(ctx, dstream)
	if err != nil {
		logger.Panic("error starting mqtt", zap.Error(err))
	}

	return tsdb, dstream, c
}

func shutdown(cancel context.CancelFunc, tsdb *timescalewrapper.Database, mqttClient mqtt.Client) {
	defer cancel()

	defer logging.Shutdown() //nolint:errcheck // We don't care about the error here.

	mqttClient.Disconnect(MQTTDisconnectTimeout)
	tsdb.Shutdown()

	logger.Info("service stopped")
}

func processLoop(ctx context.Context, ds chan [2]string, tsdb *timescalewrapper.Database) {
	for {
		select {
		case <-ctx.Done():
			return
		case d := <-ds:
			logger.Info("received sensor data", zap.String("topic", d[0]), zap.String("payload", d[1]))

			var reading timescalewrapper.SensorData

			err := json.Unmarshal([]byte(d[1]), &reading)
			if err != nil {
				logger.Error("unmarshalling payload", zap.Error(err), zap.String("payload", d[1]))
			}

			logger.Debug("unmarshalling payload", zap.Any("reading", reading))

			err = tsdb.InsertData(reading)
			if err != nil {
				logger.Error("inserting data", zap.Error(err), zap.Any("reading", reading))
			}
		}
	}
}

func readConfig() string {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Panic("reading config file", zap.Error(err))
	}

	var sb strings.Builder

	sb.WriteString("postgres://")
	sb.WriteString(viper.GetString("timescale.user"))
	sb.WriteString(":")
	sb.WriteString(viper.GetString("timescale.password"))
	sb.WriteString("@")
	sb.WriteString(viper.GetString("timescale.host"))
	sb.WriteString(":")
	sb.WriteString(viper.GetString("timescale.port"))
	sb.WriteString("/")
	sb.WriteString(viper.GetString("timescale.database"))
	sb.WriteString("?sslmode=")
	sb.WriteString(viper.GetString("timescale.sslmode"))

	return sb.String()
}
