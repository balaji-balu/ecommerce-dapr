// shipping-service/main.go (updated)
package main

import (
    "encoding/json"
    "context"
    "github.com/dapr/go-sdk/service/common"
    "github.com/dapr/go-sdk/service/http"
    "github.com/dapr/go-sdk/client"
    "net/http"
)

type Shipping struct {
    OrderID   string `json:"order_id"`
    UserID    string `json:"user_id"`
    Address   string `json:"address"`
    Status    string `json:"status"`
}

var shipments = map[string]Shipping{}

func main() {
    svc := http.NewService(":8085")

    svc.AddServiceInvocationHandler("/createShipment", createShipment)
    svc.AddTopicEventHandler("order-pubsub", "payment-processed", paymentProcessedHandler)

    if err := svc.Start(); err != nil {
        panic(err)
    }
}

func createShipment(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
    var shipping Shipping
    if err := json.Unmarshal(in.Data, &shipping); err != nil {
        return nil, err
    }
    shipping.Status = "shipped"
    shipments[shipping.OrderID] = shipping

    // Publish an event after shipment is created
    daprClient, err := client.NewClient()
    if err != nil {
        return nil, err
    }
    defer daprClient.Close()
    shipmentEvent, _ := json.Marshal(shipping)
    err = daprClient.PublishEvent(ctx, "order-pubsub", "shipment-created", shipmentEvent)
    if err != nil {
        return nil, err
    }

    return &common.Content{
        Data:        in.Data,
        ContentType: "application/json",
    }, nil
}

func paymentProcessedHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
    var payment Payment
    if err := json.Unmarshal(e.Data, &payment); err != nil {
        return false, err
    }

    // Create a shipment after payment is processed
    shipment := Shipping{
        OrderID: payment.OrderID,
        UserID:  payment.UserID,
        Address: "User Address", // Assuming address is fetched from user service
        Status:  "ready to ship",
    }
    shipments[shipment.OrderID] = shipment

    // Publish an event after shipment is created
    daprClient, err := client.NewClient()
    if err != nil {
        return false, err
    }
    defer daprClient.Close()
    shipmentEvent, _ := json.Marshal(shipment)
    err = daprClient.PublishEvent(ctx, "order-pubsub", "shipment-created", shipmentEvent)
    if err != nil {
        return false, err
    }

    return false, nil
}
