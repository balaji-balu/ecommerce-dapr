// order-service/main.go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"

    "github.com/dapr/go-sdk/client"
    "github.com/dapr/go-sdk/service/common"
    "github.com/dapr/go-sdk/service/http"
)

type Order struct {
    ID       string `json:"id"`
    UserID   string `json:"user_id"`
    ProductID string `json:"product_id"`
    Quantity int    `json:"quantity"`
    Status   string `json:"status"`
}

func main() {
    svc := http.NewService(":8082")

    svc.AddServiceInvocationHandler("/createOrder", createOrder)
    svc.AddServiceInvocationHandler("/getOrder", getOrder)

    if err := svc.Start(); err != nil {
        log.Fatalf("Failed to start the service: %v", err)
    }
}

func createOrder(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
    var order Order
    if err := json.Unmarshal(in.Data, &order); err != nil {
        return nil, err
    }
    order.Status = "created"

    // Save order state to Redis
    daprClient, err := client.NewClient()
    if err != nil {
        return nil, err
    }
    defer daprClient.Close()

    orderData, _ := json.Marshal(order)
    err = daprClient.SaveState(ctx, "statestore", order.ID, orderData)
    if err != nil {
        return nil, err
    }

    // Publish an event after order is created
    err = daprClient.PublishEvent(ctx, "order-pubsub", "order-created", orderData)
    if err != nil {
        return nil, err
    }

    return &common.Content{
        Data:        orderData,
        ContentType: "application/json",
    }, nil
}

func getOrder(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
    id := string(in.Data)

    // Retrieve order state from Redis
    daprClient, err := client.NewClient()
    if err != nil {
        return nil, err
    }
    defer daprClient.Close()

    item, err := daprClient.GetState(ctx, "statestore", id)
    if err != nil {
        return nil, err
    }
    if item == nil || len(item.Value) == 0 {
        return &common.Content{
            Data:        nil,
            ContentType: "application/json",
        }, nil
    }

    return &common.Content{
        Data:        item.Value,
        ContentType: "application/json",
    }, nil
}
