GoBunny
=======

Golang RabbitMQ Client implementation for 100DaysOfCode

Do it
-----

1. Clone this repo
2. `docker-compose -f docker-compose-rmq.yml up # Starts an RMQ instance`
3. `./build-me.sh && docker-compose up # Starts an RMQ receiver, connected to instance, in a docker container`
4. `./pkg/gobunny send # Sends "Hello World" message to the RMQ instance (no docker container)`
