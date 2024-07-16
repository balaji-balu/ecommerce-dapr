// notification-service/main.go
package main

import (
    "context"
    "encoding/json"
    "log"

    "github.com/dapr/go-sdk/service/common"
    "github.com/dapr/go-sdk/service/http"
)

type Notification struct {
    UserID  string `json:"user_id"`
    Message string `json:"message"`
}

func main() {
    svc := http.NewService(":8086")

    svc.AddTopicEventHandler("order-pubsub", "order-created", orderCreatedHandler)
    svc.AddTopicEventHandler("order-pubsub", "payment-processed", paymentProcessedHandler)
    svc.AddTopicEventHandler("order-pubsub", "shipment-created", shipmentCreatedHandler)

    if err := svc.Start(); err != nil {
        log.Fatalf("Failed to start the service: %v", err)
    }
}

func orderCreatedHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
    var order Order
    if err := json.Unmarshal(e.Data, &order); err != nil {
        return false, err
    }

    notification := Notification{
        UserID:  order.UserID,
        Message: "Your order has been created successfully.",
    }
    sendNotification(notification)

    return false, nil
}

func paymentProcessedHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
    var payment Payment
    if err := json.Unmarshal(e.Data, &payment); err != nil {
        return false, err
    }

    notification := Notification{
        UserID:  payment.UserID,
        Message: "Your payment has been processed successfully.",
    }
    sendNotification(notification)

    return false, nil
}

func shipmentCreatedHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
    var shipment Shipping
    if err := json.Unmarshal(e.Data, &shipment); err != nil {
        return false, err
    }

    notification := Notification{
        UserID:  shipment.UserID,
        Message: "Your order has been shipped.",
    }
    sendNotification(notification)

    return false, nil
}

func sendNotification(notification Notification) {
    // Simulate sending a notification (e.g., email, SMS)
    log.Printf("Sending notification to user %s: %s", notification.UserID, notification.Message)
}
