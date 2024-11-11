package consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"router-template/delivery"
	"router-template/delivery/broker"
	"time"

	"github.com/randyardiansyah25/libpkg/util/env"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ** Modul untuk consume notifikasi yang dipublish dari service-service yang memproses transaksi
func Start() {
	reconnectPeriod := env.GetInt("rabbit.reconnect_period", 10)

	retryCfg := env.GetInt("rabbit.consumer_add_retry_periode", 5)
	if retryCfg == 0 {
		retryCfg = 5
	}

	retryPeriode := time.Duration(reconnectPeriod+retryCfg) * time.Second

	for {
		if er := startConsumer(); er != nil {
			delivery.PrintLog(fmt.Sprintf("error consume : %s. \n\nretry in %d second...", er.Error(), reconnectPeriod+retryCfg))
			time.Sleep(retryPeriode)
		}
	}
}

func startConsumer() (er error) {
	delivery.PrintLog("Start Consume()...")
	queueName := env.GetString("rabbit.queue_name")
	messages, er := broker.BrokerChannel.Consume(
		queueName, // queue name
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if er != nil {
		return er
	}

	for message := range messages {
		go messageProcessor(message)
	}

	return
}

func messageProcessor(message amqp.Delivery) {
	var prettyJson bytes.Buffer
	if er := json.Indent(&prettyJson, message.Body, "", "    "); er != nil {
		delivery.PrintErrorf("We got message, it seems the message is not in json format : %s, error : %v\n", string(message.Body), er)
		message.Reject(false)
	} else {
		delivery.PrintLogf("A notification message was received :\n%s\n", prettyJson.String())
		message.Ack(false)
	}
}
