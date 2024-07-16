// product-service/main.go
package main

import (
    "context"
    "encoding/json"
    "github.com/dapr/go-sdk/service/common"
    "github.com/dapr/go-sdk/service/http"
)

type Product struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

var products = map[string]Product{}

func main() {
    svc := http.NewService(":8081")

    svc.AddServiceInvocationHandler("/createProduct", createProduct)
    svc.AddServiceInvocationHandler("/getProduct", getProduct)

    if err := svc.Start(); err != nil {
        panic(err)
    }
}

func createProduct(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
    var product Product
    if err := json.Unmarshal(in.Data, &product); err != nil {
        return nil, err
    }
    products[product.ID] = product
    return &common.Content{
        Data:        in.Data,
        ContentType: "application/json",
    }, nil
}

func getProduct(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
    id := string(in.Data)
    product, exists := products[id]
    if !exists {
        return &common.Content{
            Data:        nil,
            ContentType: "application/json",
        }, nil
    }
    productData, _ := json.Marshal(product)
    return &common.Content{
        Data:        productData,
        ContentType: "application/json",
    }, nil
}
