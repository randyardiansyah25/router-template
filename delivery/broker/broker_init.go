package broker

import (
	"fmt"
	"router-template/delivery"
	"time"

	"github.com/kpango/glg"
	"github.com/randyardiansyah25/libpkg/util/env"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	BrokerChannel  *amqp.Channel
	BrokerPrepared = false
	CloseChan      chan *amqp.Error
)

func ConnectToRabbit() {
	url := env.GetString("rabbit.url")

	heartbeatPeriod := env.GetInt("rabbit.heartbeat_period", 60)
	reconnectPeriod := env.GetInt("rabbit.reconnect_period", 6)

	for {
		delivery.PrintLog("Connecting to rabbit..")
		conn, er := amqp.DialConfig(url, amqp.Config{
			Heartbeat: time.Duration(heartbeatPeriod) * time.Second,
		})
		if er == nil {
			CloseChan = conn.NotifyClose(make(chan *amqp.Error))
			BrokerChannel, er = conn.Channel()
			if er == nil {
				if !BrokerPrepared {
					prepareBroker()
					BrokerPrepared = true
				}
				break
			}
			delivery.PrintError("Closing rabbit connection..")
			_ = conn.Close()
		}
		_ = glg.Warnf("Failed to connect to Broker: %s. Retrying in %d second...", er, reconnectPeriod)
		time.Sleep(time.Duration(reconnectPeriod) * time.Second)
	}
	delivery.PrintLog("Connected to rabbit..")
}

func prepareBroker() {

	xchange := env.GetString("rabbit.exchange")
	xchangeType := env.GetString("rabbit.exchange_type")
	routingKey := env.GetString("rabbit.routing_key")
	queueName := env.GetString("rabbit.queue_name")

	er := BrokerChannel.ExchangeDeclare(
		xchange,     //name
		xchangeType, //type
		true,        //durable,
		false,       //auto-delete
		false,       //internal
		false,       //no-wait,
		nil,         //arguments
	)

	if er != nil {
		glg.Fatal("error prepare exchange : ", er)
	}

	//** Prepare untuk consumer
	q, er := BrokerChannel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if er != nil {
		glg.Fatal("error declare consumer queue : ", er)
	}

	er = BrokerChannel.QueueBind(
		q.Name,     // queuname
		routingKey, // routing key
		xchange,    // exchange
		false,      // no-wait
		nil,        // arguments
	)

	if er != nil {
		glg.Fatal("error binding queue : ", er)
	}
}

func BrokerClosedChannelObserver() {
	for {
		<-CloseChan
		reconnectPeriod := env.GetInt("rabbit.reconnect_period", 10)
		delivery.PrintError(fmt.Sprintf("disconnected from rabbit : re-connection in %d seconds..", reconnectPeriod))
		ConnectToRabbit()
	}
}
