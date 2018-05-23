package main

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

var mailServer = "localhost"
var mailPort   = "5672"
var mailUser   = "guest"

var amqpUrl    = fmt.Sprintf("amqp://%s:%s@%s:%s/",
	mailUser,
	mailUser,
	mailServer,
	mailPort,
)

var mailQueue  = "goBunnyQ"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("s: %s\n", msg, err))
	}
}

func connect(url string) (conn *amqp.Connection) {
	fmt.Printf("Connecting to: %s\n", url)
	conn, err := amqp.Dial(url)

	failOnError(err, "Connection Failed")
	defer conn.Close()

	return conn
}

func openChannel(conn *amqp.Connection) (channel *amqp.Channel){
	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	return channel
}

func declareQueue(channel *amqp.Channel) (queue amqp.Queue){
	queue, err := channel.QueueDeclare(
		mailQueue,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   //arguments
	)
	failOnError(err, "Failed to declare a queue")

	return queue
}

func sendMail(channel *amqp.Channel, queue amqp.Queue, message string) {
	err := channel.Publish(
		"",         //exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:				 []byte(message),
 	  },
  )
	failOnError(err, "Failed to publish message")
}

func main() {
	action := "sends"
	fmt.Printf("GoBunny %v!\n", action)

	conn    := connect(amqpUrl)
	channel := openChannel(conn)
	queue   := declareQueue(channel)

  sendMail(channel, queue, "Hello World!")

}
