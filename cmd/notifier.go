package cmd

import (
    "fmt"
    "time"
    "log"
    "context"
    "encoding/json"
    "github.com/go-redis/redis/v8"
    "fractale/fractal6.go/graph/model"
    "fractale/fractal6.go/graph"
    "fractale/fractal6.go/web/middleware"
    //. "fractale/fractal6.go/tools"
)

var cache *redis.Client = redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    //Password: "", // no password set
    //DB:       0,  // use default DB
})

var ctx = context.Background()

func RunNotifier() {
    // Test connection
    if _, err := cache.Ping(ctx).Result(); err != nil {
        log.Fatal("redis error: ", err)
    }

    // Init Suscribe channel
    // Queuing limit, and concurency see:
    // https://stackoverflow.com/questions/27745842/redis-pubsub-and-message-queueing
    // https://github.com/go-redis/redis/issues/653
    subscriber := cache.Subscribe(
        ctx,
        "api-tension-notification",
        "api-contract-notification",
        "api-notif-notification",
    )

    if _, err := subscriber.Receive(ctx); err != nil {
        log.Fatal("Failed to receive from suscriber: ", err)
        return
    }

    for msg := range subscriber.Channel() {
        switch msg.Channel {
        case "api-tension-notification":
            go processTensionNotification(msg)

        case "api-contract-notification":
            go processContractNotification(msg)

        case "api-notif-notification":
            go processNotifNotification(msg)

        }
    }
}

func processTensionNotification(msg *redis.Message) {
    defer middleware.NotifRecover("tension event")
    // Extract message
    var notif model.EventNotif
    if err := json.Unmarshal([]byte(msg.Payload), &notif); err != nil {
        log.Printf("unmarshaling error for channel %s: %v", msg.Channel, err)
    }
    if len(notif.History) == 0 {
        log.Printf("No event in notif.")
        return
    }

    // Push notification
    if err := graph.PushEventNotifications(notif); err != nil {
        log.Printf("PushEventNotifications error: %v", err)
    }

    fmt.Printf("e")
}

func processContractNotification(msg *redis.Message) {
    defer middleware.NotifRecover("contract event")
    // Extract message
    var notif model.ContractNotif
    if err := json.Unmarshal([]byte(msg.Payload), &notif); err != nil {
        log.Printf("unmarshaling error for channel %s: %v", msg.Channel, err)
    }
    if notif.Contract == nil {
        log.Printf("No contract in notif.")
        return
    }

    // Add a little sleep to wait for UpdateContractHook as it writes after publishing
    time.Sleep(1 * time.Second)

    // Push notification
    if err := graph.PushContractNotifications(notif); err != nil {
        log.Printf("PushContractNotification error: %v: ", err)
    }

    fmt.Printf("c")
}

func processNotifNotification(msg *redis.Message) {
    defer middleware.NotifRecover("notif event")
    // Extract message
    var notif model.NotifNotif
    if err := json.Unmarshal([]byte(msg.Payload), &notif); err != nil {
        log.Printf("unmarshaling error for channel %s: %v", msg.Channel, err)
    }
    if len(notif.Msg) == 0 {
        log.Printf("No message in notif.")
        return
    }

    // Push notification
    if err := graph.PushNotifNotifications(notif, false); err != nil {
        log.Printf("PushEventNotifications error: %v", err)
    }

    fmt.Printf("n")
}
