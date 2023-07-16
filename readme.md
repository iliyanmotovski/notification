# Notification service

## System diagram
![alt text](https://github.com/iliyanm/notification/blob/master/notification_service.png?raw=true)

## Description
Notification service allows you to send notifications to different channels.

Channels implemented:
1. SMS (implemented Twilio)
2. Email (implemented Mailjet)
3. Slack

The service is easily extensible and adding new channels is easy. Implementing new providers (gateways) to the
existing channels is easy too. 

The service is built on top of BeeORM, MySQL, Redis and Redis Streams.
It guarantees `at-least-once` delivery due to the nature of [Redis Streams](https://redis.com/solutions/use-cases/messaging/) - search for `at-least-once`.

## Configuration
The service config file path is: `./config/config.yaml` 

## How to run
1. `cd docker ; ./build.sh ; cd ..`
2. `open your favorite MySQL GUI client and create 2 empty databases: "notification" and "notification_test" - the rest will be taken care of by the ORM when you start the apps`
3. `make consumer` - this cron will consume notification events and send them, you can scale it using 2 approaches:
   - run`make consumer` in new terminals
   - by running following code in consumer `main.go` in a for loop 5-10-20 times. 
```go
	consumerRunner.RunConsumerMany(
		consumerbeeorm.NotificationsDirtyConsumer(smsService, emailService, slackService),
		queue.OrmDirtyNotificationEntity,
		100,
	)
```
4. `make notification-api` - it fires up an HTTP server on 8080, you can change the port in the config file
5. `open swagger` [doc](http://localhost:8080/doc)
6. `send request to "POST /v1/notification/send-async" with the following body (or change it how you want)`

#### Request body

Fields not to be changed if you want the example to work:
```
email.from // don't change, as mailjet is configured to send from this email address
email.template_name // don't change, as this is mailjet template id that I use for test

slack // if you want to test slack, you need to create an organization and configure a bot
```
#### JSON
```json
{
  "email": {
    "from": "ilqnskiq@abv.bg",
    "from_name": "iliyan",
    "subject": "test",
    "template_name": "4954974",
    "to": "iliyan.motovski@gmail.com"
  },
  "slack_message": {
    "bot_name": "bot",
    "channel_name": "bot",
    "message": "test message"
  },
  "sms": {
    "mobile_number": "+359878697929",
    "text": "test"
  }
}
```

## Stopping an app
It's important to stop the consumer and the HTTP server correctly. They expect signal `syscall.SIGINT` or `syscall.SIGTERM`. The easiest way
to stop an app is to go to the terminal its running and `Ctrl+C`. Don't close a terminal without stopping the application first, as
this will cause the application to continue running inside the docker container, and you need to restart docker in order to stop it.

## Testing
Run the tests: `make test`

Run the tests with coverage: `make test-cover`

## Linter and formatting
Run `make format-check`

## Regenerating swagger (every time you change an API)
Run `make generate-notification-api-doc`
