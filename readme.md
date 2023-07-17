# Notification service

## System diagram
![alt text](https://github.com/iliyanm/notification/blob/master/notification_service.png?raw=true)

## Description
Notification service allows you to send notifications to different channels.

Channels implemented:
1. SMS (implemented Twilio)
2. Email (implemented Mailjet)
3. Slack

The service is open for extensibility, and adding a new channel is a straightforward process, as well as adding implementations of providers (gateways).

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

```
sms   // if you want to test sms, you need to create account in twilio and put the api keys in the config file
email // if you want to test email, you need to create account in mailjet nad put the api keys in the config file
slack // if you want to test slack, you need to create an organization and configure a bot, then put the boy key in the config file
```

#### JSON
```json
{
   "email": {
      "from": "string",
      "from_name": "string",
      "reply_to": "string",
      "subject": "string",
      "template_data": [
         {
            "key": "string",
            "value": "string"
         }
      ],
      "template_name": "string",
      "to": "string"
   },
   "slack_message": {
      "bot_name": "string",
      "channel_name": "string",
      "message": "string"
   },
   "sms": {
      "mobile_number": "string",
      "text": "string"
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
