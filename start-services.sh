
# --components-path ./components

# Start User Service
dapr run --config ./config/user-service-config.json -- go run user-service/main.go

# Start Product Service
dapr run --config ./config/product-service-config.json -- go run product-service/main.go

# Start Order Service
dapr run --config ./config/order-service-config.json -- go run order-service/main.go

# Start Payment Service
dapr run --config ./config/payment-service-config.json -- go run payment-service/main.go

# Start Inventory Service
dapr run --config ./config/inventory-service-config.json -- go run inventory-service/main.go

# Start Shipping Service
dapr run --config ./config/shipping-service-config.json -- go run shipping-service/main.go

# Start Notification Service
dapr run --config ./config/notification-service-config.json -- go run notification-service/main.go