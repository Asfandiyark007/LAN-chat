# Architecture for Local chat:
For this project i picked client/server model. 

# Why:
A centralized server fits the requirements of this terminal chat application better because it provides reliable message handling, user management and persistent storage (keeps message history). For now i am going to do the simplest connection method possible, client connect by entering the server's IP manually, later i would like to add more advanced discovery method. 

# Server:
- When it receives a message from user, it sends the acknowledgement to the user and we use that to control the terminal color (green or red).
- When a user sends message to the server, the server should send the same message to everyone connected to that server with the sender username. (not the sender)
- when a user disconnect the server should remove him from the connected clients.
- When a user connects they should get the most recent messages of previous chat (like 5 most recent).
- Server keeps the Logs of users messages.   
- Should timestamp the messages as well.

# Client:
- should be able to connect to server using known IP
- should be able to send and receive messages.
- should be able to use Acknowledge Signals to change 
- terminal color (maybe red when not connected)

# Message flow:
```
user#1 ---> Send's message ---> Server ---> send it to everone (execpt sender ---> send's message reciveing acknowledgement).
```