// user-service/main.go
package main

import (
	"context"
	"fmt"
    "encoding/json"
//    "net/http"
//    "github.com/gorilla/mux"
    "github.com/dapr/go-sdk/service/common"
    "github.com/dapr/go-sdk/service/http"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var users = map[string]User{}

func main() {
    svc := http.NewService(":8080")

    svc.AddServiceInvocationHandler("/createUser", createUser)
    svc.AddServiceInvocationHandler("/getUser", getUser)

    if err := svc.Start(); err != nil {
        panic(err)
    }
}

func createUser(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
    var user User
    if err := json.Unmarshal(in.Data, &user); err != nil {
        return nil, err
    }
    users[user.ID] = user
    return &common.Content{
        Data:        in.Data,
        ContentType: "application/json",
    }, nil
}

func getUser(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
    id := string(in.Data)
    user, exists := users[id]
	fmt.Println("user: ", user, id)
    if !exists {
        return &common.Content{
            Data:        nil,
            ContentType: "application/json",
        }, nil
    }
    userData, _ := json.Marshal(user)
    return &common.Content{
        Data:        userData,
        ContentType: "application/json",
    }, nil
}
