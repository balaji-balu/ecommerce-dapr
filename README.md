Each directiry:
```sh
go mod init example.com/m
go mod tidy
```
Run go binary
```sh   
go run user-service/main.go
```

Run dapr service 
```sh
dapr run --config ./config/user-service-config.json -- go run user-service/main.go

dapr run --config ./config/order-service-config.json --components-path ./components -- go run order-service/main.go
```