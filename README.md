# websocket-example

## Build code
```
git clone https://github.com/rishabh625/websocket-example.git
cd websocket-example
go mod vendor
go build -o websocket
```
## Run code
```
./websocket
```

## Test code
Go to https://www.postman.com/winter-capsule-279801/workspace/new-team-workspace/
1) Run Get Auth Token from TestChatroom
    in query params can change username and password following userdetails.json

2) Add Bearer token in websocket Authorization Header and connect to server
follow Websocket call collection