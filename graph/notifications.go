package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"fractale/fractal6.go/db"
	"fractale/fractal6.go/graph/auth"
	"fractale/fractal6.go/graph/model"
	. "fractale/fractal6.go/tools"
	"fractale/fractal6.go/web/email"
)

var ctx context.Context = context.Background()

//
// Publisher functions (Redis)
//

// Will trigger processTensionNotification in cmd/notifier.go
// and PushEventNotifications
func PublishTensionEvent(notif model.EventNotif) error {
    payload, _ := json.Marshal(notif)
    if err := cache.Publish(ctx, "api-tension-notification", payload).Err(); err != nil {
        fmt.Printf("Redis publish error: %v", err)
        panic(err)
    }

    return nil
}

// Will trigger processContractNotification in cmd/notifier.go
// and PushContractNotifications
func PublishContractEvent(notif model.ContractNotif) error {
    payload, _ := json.Marshal(notif)
    if err := cache.Publish(ctx, "api-contract-notification", payload).Err(); err != nil {
        fmt.Printf("Redis publish error: %v", err)
        panic(err)
    }

    return nil
}

// Will trigger processNotifNotification in cmd/notifier.go
// and PushNotifNotifications
func PublishNotifEvent(notif model.NotifNotif) error {
    payload, _ := json.Marshal(notif)
    if err := cache.Publish(ctx, "api-notif-notification", payload).Err(); err != nil {
        fmt.Printf("Redis publish error: %v", err)
        panic(err)
    }

    return nil
}

//
// User helpers
//

// GetUserToNotify returns a list of user should receive notifications uponf tension updates.
func GetUsersToNotify(tid string, withAssignees, withSubscribers bool) (map[string]model.UserNotifInfo, error) {
    users := make(map[string]model.UserNotifInfo)
    user_selection := "User.username User.email User.name User.notifyByEmail"

    {
        // Get Coordos
        coordos, err := auth.GetCoordosFromTid(tid)
        if err != nil { return users, err }
        for _, user := range coordos {
            if _, ex := users[user.Username]; ex { continue }
            users[user.Username] = model.UserNotifInfo{User: user, Reason: model.ReasonIsCoordo}
        }
    }

    {
        // Get First-link
        res, err := db.GetDB().GetSubSubFieldById(tid, "Tension.receiver", "Node.first_link", user_selection)
        if err != nil { return users, err }
        if res != nil {
            var user model.User
            if err := Map2Struct(res.(model.JsonAtom), &user); err == nil {
                if _, ex := users[user.Username]; !ex {
                    users[user.Username] = model.UserNotifInfo{User: user, Reason: model.ReasonIsFirstLink}
                }
            }
        }
    }

    if withAssignees {
        // Get Assignees
        res, err := db.GetDB().GetSubFieldById(tid, "Tension.assignees", user_selection)
        if err != nil { return users, err }
        if assignees, ok := InterfaceSlice(res); ok {
            for _, u := range assignees {
                var user model.User
                if err := Map2Struct(u.(model.JsonAtom), &user); err == nil {
                    if _, ex := users[user.Username]; ex { continue }
                    users[user.Username] = model.UserNotifInfo{User: user, Reason: model.ReasonIsAssignee}
                }
            }
        }
    }

    if withSubscribers {
        // Get Subscribers
        res, err := db.GetDB().GetSubFieldById(tid, "Tension.subscribers", user_selection)
        if err != nil { return users, err }
        if subscribers, ok := InterfaceSlice(res); ok {
            for _, u := range subscribers {
                var user model.User
                if err := Map2Struct(u.(model.JsonAtom), &user); err == nil {
                    if _, ex := users[user.Username]; ex { continue }
                    users[user.Username] = model.UserNotifInfo{User: user, Reason: model.ReasonIsSubscriber}
                }
            }
        }
    }


    return users, nil
}


//
// Notifiers functions
//

/* INTERNAL (websocket, platform notification etc) */

// PushHistory publish event to a tension history.
func PushHistory(notif *model.EventNotif) error {
    var inputs []model.AddEventInput
    for _, e := range notif.History {
        // Build AddtensionInput
        var temp model.AddEventInput
        StructMap(e, &temp)
        temp.Tension = &model.TensionRef{ID: &notif.Tid}

        // Push AddtensionInput
        inputs = append(inputs, temp)
    }
    // Push events
    ids, err := db.GetDB().AddMany(*notif.Uctx, "event", inputs)
    if err != nil { return err }
    // Set event ids for further notifications
    for i, id := range ids {
        notif.History[i].ID = &id
    }
    return err
}

/* EXTERNAL (email, chat, etc) */

// Notify users for Event events, where events can be batch of event.
func PushEventNotifications(notif model.EventNotif) error {
    // Only the event with an ID will be notified.
    var eventBatch []*model.EventKindRef
    var createdAt string
    for i, e := range notif.History {
        if i == 0 {
            createdAt = *e.CreatedAt
        }
        if *e.ID != "" {
            eventBatch = append(eventBatch, &model.EventKindRef{EventRef: &model.EventRef{ID: e.ID}})
        }
    }
    if len(eventBatch) == 0 {
        return nil
    }

    // Get people to notify
    users := make(map[string]model.UserNotifInfo)
	var err error
    var isAlert bool
    var receiverid string
    if notif.HasEvent(model.TensionEventCreated) {
        if t, err := db.GetDB().GetFieldById(notif.Tid, "Tension.type_ Tension.receiverid"); err != nil {
            return err
        } else if t != nil {
            tension := t.(model.JsonAtom)
            isAlert = tension["type_"].(string) == string(model.TensionTypeAlert)
            receiverid = tension["receiverid"].(string)
        }
    }
    // --
    if isAlert {
        // Alert tension created: Notify everyone
        user_selection := "User.username User.email User.name User.notifyByEmail"
        if data, err := db.GetDB().GetSubMembers("nameid", receiverid, user_selection); err != nil {
            return err
        } else {
            for _, n := range data {
                user := *n.FirstLink
                if _, ex := users[user.Username]; ex { continue }
                users[user.Username] = model.UserNotifInfo{User: user, Reason: model.ReasonIsAlert}
            }
        }
    } else {
        // Notify only suscribers and relative.
        users, err = GetUsersToNotify(notif.Tid, true, true)
        if err != nil { return err }
    }

    // Push user event notification
    for u, ui := range users {
        // Don't self notify.
        if u == notif.Uctx.Username { continue }
        // Pending user has no history yet
        if ui.Reason == model.ReasonIsPendingCandidate { continue }

        // User Event
        eid, err := db.GetDB().Add(db.GetDB().GetRootUctx(), "userEvent", &model.AddUserEventInput{
            User: &model.UserRef{Username: &u},
            IsRead: false,
            CreatedAt: createdAt,
            Event: eventBatch,
        })
        if err != nil { return err }

        // Email
         if notif.Uctx.Rights.HasEmailNotifications && ui.User.NotifyByEmail {
             ui.Eid = eid
             err = email.SendEventNotificationEmail(ui, notif)
             if err != nil { return err }
        }
    }

    return err
}

// Notify users for Contract event.
func PushContractNotifications(notif model.ContractNotif) error {
    // Only the event with an ID will be notified.
    var eventBatch []*model.EventKindRef
    var createdAt string
    if notif.Contract == nil {
        return nil
    }
    createdAt = notif.Contract.CreatedAt
    eventBatch = append(eventBatch, &model.EventKindRef{ContractRef: &model.ContractRef{ID: &notif.Contract.ID}})

    // Get people to notify
    users, err := GetUsersToNotify(notif.Tid, true, false)
    if err != nil { return err }
    // +
    // Add Candidates
    for _, c := range notif.Contract.Candidates {
        users[c.Username] = model.UserNotifInfo{User: *c, Reason: model.ReasonIsCandidate}
    }
    // +
    // Add Pending Candidates
    for _, c := range notif.Contract.PendingCandidates {
        if c.Email == nil { continue }
        users[*c.Email] = model.UserNotifInfo{User: model.User{Email: *c.Email}, Reason: model.ReasonIsPendingCandidate}
    }
    // +
    // Add Participants
    for _, p := range notif.Contract.Participants {
        if _, ex := users[p.Node.FirstLink.Username]; ex { continue }
        users[p.Node.FirstLink.Username] = model.UserNotifInfo{User: *p.Node.FirstLink, Reason: model.ReasonIsParticipant}
    }

    // Push user event notification
    for u, ui := range users {
        // Don't self notify.
        if u == notif.Uctx.Username { continue }

        // User Event
        var eid string
        if ui.Reason == model.ReasonIsPendingCandidate {
            // Update pending users
            err = MaybeSetPendingUserToken(u)
            if err != nil { return err }
            // Link contract
            err = db.GetDB().Update(db.GetDB().GetRootUctx(), "pendingUser", &model.UpdatePendingUserInput{
                Filter: &model.PendingUserFilter{Email: &model.StringHashFilter{Eq: &u}},
                Set: &model.PendingUserPatch{Contracts: []*model.ContractRef{&model.ContractRef{ID: &notif.Contract.ID}}},
            })
            if err != nil { return err }
        } else {
            switch notif.ContractEvent {
            case model.NewContract:
                // Push user event
                eid, err = db.GetDB().Add(db.GetDB().GetRootUctx(), "userEvent", &model.AddUserEventInput{
                    User: &model.UserRef{Username: &u},
                    IsRead: false,
                    CreatedAt: createdAt,
                    Event: eventBatch,
                })
                if err != nil { return err }
            case model.NewComment:
                // Push user notif
                PushNotifNotifications(model.NotifNotif{
                    Uctx: notif.Uctx,
                    Tid: &notif.Tid,
                    Cid: &notif.Contract.ID,
                    Msg: "You have a new comment",
                    To: []string{u},
                })
            }
        }

        // Email
        if notif.Uctx.Rights.HasEmailNotifications && (ui.User.NotifyByEmail || ui.Reason == model.ReasonIsPendingCandidate) &&
        (ui.Reason == model.ReasonIsCandidate ||
        ui.Reason == model.ReasonIsPendingCandidate ||
        ui.Reason == model.ReasonIsParticipant ||
        ui.Reason == model.ReasonIsCoordo ||
        ui.Reason == model.ReasonIsAssignee ) {
        fmt.Println(2)
             ui.Eid = eid
             err = email.SendContractNotificationEmail(ui, notif)
             if err != nil { return err }
        }
    }

    return err
}

// Notify users for Notif events.
func PushNotifNotifications(notif model.NotifNotif) error {
    // Only the event with an ID will be notified.
    var eventBatch []*model.EventKindRef
    var createdAt string = Now()
    var tensionRef model.TensionRef
    var contractRef model.ContractRef
    if notif.Tid != nil {
        tensionRef = model.TensionRef{ID: notif.Tid}
    }
    if notif.Cid != nil {
        contractRef = model.ContractRef{ID: notif.Cid}
    }

    eventBatch = append(eventBatch, &model.EventKindRef{NotifRef: &model.NotifRef{
        CreatedAt: &createdAt,
        CreatedBy: &model.UserRef{Username: &notif.Uctx.Username},
        Message: &notif.Msg,
        Tension: &tensionRef,
        Contract: &contractRef,
    }})


    // Push user event notification
    for _, u := range notif.To {
        // Don't self notify.
        if u == notif.Uctx.Username { continue }

        // User Event
        _, err := db.GetDB().Add(db.GetDB().GetRootUctx(), "userEvent", &model.AddUserEventInput{
            User: &model.UserRef{Username: &u},
            IsRead: false,
            CreatedAt: createdAt,
            Event: eventBatch,
        })
        if err != nil { return err }

        // Email
        // No email for this one
    }

    return nil
}
