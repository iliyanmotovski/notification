# Notification service

## System diagram  
![alt text](https://github.com/iliyanm/notification/blob/master/notification_service.png?raw=true)

## The service exposes HTTP Rest server

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
