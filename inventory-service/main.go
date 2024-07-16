// inventory-service/main.go
package main

import (
    "encoding/json"
    "context"
    "github.com/dapr/go-sdk/service/common"
    "github.com/dapr/go-sdk/service/http"
    "github.com/dapr/go-sdk/client"
    "net/http"
)

type Inventory struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

var inventories = map[string]Inventory{}

func main() {
    svc := http.NewService(":8084")

    svc.AddServiceInvocationHandler("/updateInventory", updateInventory)
    svc.AddTopicEventHandler("order-pubsub", "order-created", orderCreatedHandler)

    if err := svc.Start(); err != nil {
        panic(err)
    }
}

func updateInventory(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
    var inventory Inventory
    if err := json.Unmarshal(in.Data, &inventory); err != nil {
        return nil, err
    }
    inventories[inventory.ProductID] = inventory
    return &common.Content{
        Data:        in.Data,
        ContentType: "application/json",
    }, nil
}

func orderCreatedHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
    var order Order
    if err := json.Unmarshal(e.Data, &order); err != nil {
        return false, err
    }

    // Update inventory
    inventory, exists := inventories[order.ProductID]
    if exists {
        inventory.Quantity -= order.Quantity
        inventories[order.ProductID] = inventory
    }

    return false, nil
}
