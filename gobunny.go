package main

// A basic RMQ client in Go, which parses command line args and
// sends them as a message to a localhost RMQ server.

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"flag"
)

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

	return conn
}

func openChannel(conn *amqp.Connection) (channel *amqp.Channel) {
	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return channel
}

func declareQueue(channel *amqp.Channel, mailQueue string) (queue amqp.Queue) {
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
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	failOnError(err, "Failed to publish message")
	log.Printf(" [x] Sent: %s", message)
}

func listenForMail(channel *amqp.Channel, queue amqp.Queue) {
  messages, err := channel.Consume(
  	//"",         // no exchange?
  	queue.Name, // routing key
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
  failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for message := range messages {
      log.Printf("Received a message: %s", message.Body)
		}
	}()
  log.Printf("[*] Waiting for messages.  Press CTRL+C to exit.")
	<-forever
}

func main() {

	// Default connection info
  var mailServer = "localhost"
	var mailPort = "5672"
  var mailUser = "guest"
  var mailQueue = "goBunnyQ"

	// Check to see if a different RMQ server is specified
	if os.Getenv("rmq_server") != "" {
		mailServer = os.Getenv("rmq_server")
	}

	// Construct the AMQP URL based on info above
  var amqpUrl = fmt.Sprintf("amqp://%s:%s@%s:%s/",
  	mailUser,
  	mailUser,
  	mailServer,
  	mailPort,
  )

	// "send" option on the command line, with a default message
	sendCommand := flag.NewFlagSet("send", flag.ExitOnError)
	sendMessagePtr := sendCommand.String("message",
		"Hello World!",
		"message to send")

	// "listen" option on the command line
	listenCommand := flag.NewFlagSet("listen", flag.ExitOnError)

	// Require arguments
	if len(os.Args) < 2 {
			fmt.Println("send or listen sub-command is required")
			os.Exit(1)
	}

	// Parse commands based on input
	switch os.Args[1] {
	case "send":
		sendCommand.Parse(os.Args[2:])
	case "listen":
		listenCommand.Parse(os.Args[2:])
	default:
		 fmt.Println("No command parsed")
		 os.Exit(1)
	}

	// THIS ISN'T WORKING: sendMessagePtr *always* has the
	// default value, when it should contain the content of the
	// string I pass on the command line.  WTH?
	if sendCommand.Parsed() {
	  if *sendMessagePtr == "" {
	  	sendCommand.PrintDefaults()
	  	os.Exit(1)
	  }
	  fmt.Println("GoBunny Sends!")
	  conn := connect(amqpUrl)
	  channel := openChannel(conn)
	  queue := declareQueue(channel, mailQueue)

	  sendMail(channel, queue, *sendMessagePtr)
	  defer channel.Close()
	  defer conn.Close()
	}

	if listenCommand.Parsed() {
	  fmt.Println("GoBunny Listens!")
	  conn := connect(amqpUrl)
	  channel := openChannel(conn)
	  queue := declareQueue(channel, mailQueue)

	  listenForMail(channel, queue)
	  defer channel.Close()
	  defer conn.Close()
	}
}
