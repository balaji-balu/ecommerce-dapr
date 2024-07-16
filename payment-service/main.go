// payment-service/main.go
package main

import (
    "encoding/json"
    "context"
    "github.com/dapr/go-sdk/service/common"
    "github.com/dapr/go-sdk/service/http"
    "net/http"
)

type Payment struct {
    OrderID   string `json:"order_id"`
    UserID    string `json:"user_id"`
    Amount    float64 `json:"amount"`
    Status    string `json:"status"`
}

var payments = map[string]Payment{}

func main() {
    svc := http.NewService(":8083")

    svc.AddServiceInvocationHandler("/processPayment", processPayment)

    if err := svc.Start(); err != nil {
        panic(err)
    }
}

func processPayment(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
    var payment Payment
    if err := json.Unmarshal(in.Data, &payment); err != nil {
        return nil, err
    }
    payment.Status = "processed"
    payments[payment.OrderID] = payment

    // Publish an event after payment is processed
    daprClient, err := client.NewClient()
    if err != nil {
        return nil, err
    }
    defer daprClient.Close()
    paymentEvent, _ := json.Marshal(payment)
    err = daprClient.PublishEvent(ctx, "order-pubsub", "payment-processed", paymentEvent)
    if err != nil {
        return nil, err
    }

    return &common.Content{
        Data:        in.Data,
        ContentType: "application/json",
    }, nil
}
