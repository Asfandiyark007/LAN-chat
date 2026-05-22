# Message Protocol
I selected JSON based message protocol and the reason i picked it cause i worked with it in the past and it is known for structured and readable. For this we will use JSON parser. For the timestamp format i picked ```ISO 8601```
## General message protocol:

```JSON
{
    "type": "message",
    "from": "Bob",
    "to": "server",
    "content": "Hi! How is it?",
    "read" : "true",
    "message_id": "uuid",
    "user_id": "123123",
    "current_client_time": "",
    "current_server_time": ""


}

```
## General user joining protocol:

```JSON
{
    "type": "joining",
    "from": "Bob",
    "to": "server",
    "current_client_time": "",
    "user_id": "123123"

}

```
## General user leaving protocol:

```JSON
{
    "type": "leaving",
    "from": "Bob",
    "to": "server",
    "current_client_time": "",
    "user_id" : "123123",

}

```
## General server response protocol:

```JSON
{
    "type": "server_res",
    "from": "server",
    "to": "user",
    "current_server_time": "",
    "user_id" : "123123",
    "acknow" : "true",

}

```

# BOB joins:

```JSON
{
    "type": "joining",
    "from": "Bob",
    "to": "server",
    "user_id": "usr_001",
    "current_client_time": "2026-05-22T10:30:00Z"
} 
```

# BOB sends a message:

```JSON
{
    "type": "message",
    "from": "Bob",
    "content": "Hey everyone, just joined!",
    "message_id": "msg_001",
    "user_id": "usr_001",
    "current_server_time": "2026-05-22T10:30:05Z"
}
```
# Server acknowledges:

```JSON
{
    "type": "server_res",
    "from": "server",
    "to": "Bob",
    "user_id": "usr_001",
    "current_server_time": "2026-05-22T10:30:05Z",
    "acknow": true
}
```
