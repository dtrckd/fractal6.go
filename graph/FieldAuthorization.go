package graph

import (
	"fmt"
	"context"
	"github.com/99designs/gqlgen/graphql"

	"fractale/fractal6.go/db"
	"fractale/fractal6.go/graph/model"
	"fractale/fractal6.go/graph/auth"
	. "fractale/fractal6.go/tools"
	webauth "fractale/fractal6.go/web/auth"
)

var FieldAuthorizationFunc map[string]func(context.Context, interface{}, graphql.Resolver, *string, []model.TensionEvent, *int) (interface{}, error)

func init() {

    FieldAuthorizationFunc = map[string]func(context.Context, interface{}, graphql.Resolver, *string, []model.TensionEvent, *int) (interface{}, error){
        "isOwner": isOwner,
        "unique": unique,
        "oneByOne": oneByOne,
        "hasEvent": hasEvent,
        "tensionTypeCheck": tensionTypeCheck,
        "ref": ref,
        "minLen": minLength,
        "maxLen": maxLength,
    }

}



//isOwner Check that object is own by the user.
// If user(u) field is empty, assume a user object, else field should match the user(u) credential.
func isOwner(ctx context.Context, obj interface{}, next graphql.Resolver, f *string, e []model.TensionEvent, n *int) (interface{}, error) {
    // Retrieve userCtx from token
    ctx, uctx, err := webauth.GetUserContext(ctx)
    if err != nil { return nil, LogErr("Access denied", err) }

    // Get attributes and check everything is ok
    userObj := make(model.JsonAtom)
    var userField string
    if f == nil {
        userField = "user"
        userObj[userField] = obj
    } else {
        userField = *f
        userObj = obj.(model.JsonAtom)
    }

    ok, err := CheckUserOwnership(ctx, uctx, userField, userObj)
    if err != nil { return nil, LogErr("Access denied", err) }
    if ok { return next(ctx) }

    return nil, LogErr("Access Denied", fmt.Errorf("bad ownership."))
}

//unique Check uniqueness (@DEBUG follow @unique dgraph field iplementation)
// Ensure the field value is unique. If a field is given, it check the uniqueness on a subset of the parent type.
func unique(ctx context.Context, obj interface{}, next graphql.Resolver, f *string, e []model.TensionEvent, n *int) (interface{}, error) {
    data, err := next(ctx)
    if err != nil { return nil, err }

    var v string
    switch d := data.(type) {
    case *string:
        v = *d
    case string:
        v = d
    }

    field := *graphql.GetPathContext(ctx).Field
    if f != nil {
        // Extract the fieldname and type of the object queried
        _, typeName, _, err := queryTypeFromGraphqlContext(ctx)
        if err != nil { return nil, LogErr("unique", err) }
        fieldName := typeName + "." + field
        filterName := typeName + "." + *f
        s := obj.(model.JsonAtom)[*f]
        if s != nil {
            // *f is present in the inut
            //pass
        } else if ctx.Value("id") != nil {
            s, err = db.GetDB().GetFieldById(ctx.Value("id").(string), filterName)
            if err != nil || s == nil { return nil, LogErr("Internal error", err) }
        } else {
            return nil, LogErr("Value Error", fmt.Errorf("'%s' or id is required.", *f))
        }
        filterValue := s.(string)

        // Check existence
        ex, err :=  db.GetDB().Exists(fieldName, v, &filterName, &filterValue)
        if err != nil { return nil, LogErr("Internal error", err) }
        if !ex {
            return data, err
        }
    } else {
        return nil, fmt.Errorf("@unique alone not implemented.")
    }

    return data, LogErr("Duplicate error", fmt.Errorf("'%s' is already taken", field))
}

//oneByOne ensure that the mutation on the given field should contains at least one element.
func oneByOne(ctx context.Context, obj interface{}, next graphql.Resolver, f *string, e []model.TensionEvent, n *int) (interface{}, error) {
    data, err := next(ctx)
    slice, ok := InterfaceSlice(data)
    if !ok {
        field := *graphql.GetPathContext(ctx).Field
        return nil, fmt.Errorf("Data must be an array '%s'", field)
    }
    if len(slice) > 1 {
        field := *graphql.GetPathContext(ctx).Field
        return nil, LogErr("@oneByOne error", fmt.Errorf("Only one object allowed in slice '%s'", field))
    }
    return data, err
}

//hasEvent ensure the given events are present in the `history` property.
func hasEvent(ctx context.Context, obj interface{}, next graphql.Resolver, f *string, e []model.TensionEvent, n *int) (interface{}, error) {
    var events []interface{}
    events_ := obj.(model.JsonAtom)["history"]
    if events_ != nil {
        events = events_.([]interface{})
    }

    for _, event := range e  {
        for _, eventPresent := range events {
            eventType := eventPresent.(model.JsonAtom)["event_type"].(string)
            if model.TensionEvent(eventType) == event {
                // ok
                return next(ctx)
            }
        }
    }

    // Exception if we got Blob with just an ID. Use to identify the blob user are working on.
    blobs_ := obj.(model.JsonAtom)["blobs"]
    blobs, ok := InterfaceSlice(blobs_)
    if !ok {
        field := *graphql.GetPathContext(ctx).Field
        return nil, fmt.Errorf("Blobs must be an array '%s'", field)
    }
    if len(blobs) == 1 {
        b := blobs[0].(model.JsonAtom)
        if len(b) == 1 && b["id"] != nil {
            // Allows just reference to a blob (@DEBUG: add a 'cut_blob' to remove it as we do for history to prevent blob hack.
            // ok
            return next(ctx)
        }
    }

    field := *graphql.GetPathContext(ctx).Field
    if ctx.Value("hasSet") != nil && ctx.Value("hasRemove") != nil && (field == "labels" || field == "assignees") {
        // @DEBUG: detect events when remove is used in in updates
        // @DEBUG: detect that the number of event is equal to the current field len
        // when removing labels or assigness...
        return next(ctx)
    }
    return nil, LogErr("Event error", fmt.Errorf("missing event for field '%s'", field))
}

//tensionTypeCheck check is the user can use a tension type.
func tensionTypeCheck(ctx context.Context, obj interface{}, next graphql.Resolver, f *string, e []model.TensionEvent, n *int) (interface{}, error) {
    data, err := next(ctx)
    if err != nil { return nil, err }

    if data.(model.TensionType) == model.TensionTypeAlert {
        ctx, uctx, err := webauth.GetUserContext(ctx)

        // Get receiverid
        var receiverid string
        if v := obj.(model.JsonAtom)["receiverid"]; v != nil {
            receiverid = v.(string)
        } else if ctx.Value("id") != nil {
            x, err := db.GetDB().GetFieldById(ctx.Value("id").(string), "Tension.receiverid")
            if err != nil || x == nil { return nil, LogErr("Internal error", err) }
            receiverid = x.(string)
        } else {
            return nil, LogErr("Value Error", fmt.Errorf("'%s' or id is required.", *f))
        }

        // Check auth
        if i := auth.IsCoordo(uctx, receiverid); i >=0 {
            return data, err
        } else if err != nil {
            return nil, err
        }

        return nil, fmt.Errorf("Only coordinator can create Alert tension.")
    }

    return data, err
}

//ref ensure the given objects are just link to an existing one. (@weak: by testing that its size if not equal to one.)
func ref(ctx context.Context, obj interface{}, next graphql.Resolver, f *string, e []model.TensionEvent, n *int) (interface{}, error) {
    data, err := next(ctx)
    if err != nil { return nil, err }
    test :=  func(x interface{}) bool {
        return len(CleanNilMap(Struct2Map(x))) == 1
    }
    var pass bool
    data_list, ok := InterfaceSlice(data)
    if ok {
        for _, d := range data_list {
            if test(d) {
                pass = true
            } else {
                pass = false
                break
            }
        }
    } else {
        pass = test(data)
    }

    if pass {
        return data, err
    }
    field := *graphql.GetPathContext(ctx).Field
    return nil, fmt.Errorf("ref: only referecence allowed for: %s", field)
}

//inputMinLength the that the size of the field is stricly lesser than the given value
func minLength(ctx context.Context, obj interface{}, next graphql.Resolver, f *string, e []model.TensionEvent, n *int) (interface{}, error) {
    var l int
    data, err := next(ctx)
    if err != nil { return nil, err }

    switch d := data.(type) {
    case *string:
        l = len(*d)
    case string:
        l = len(d)
    default:
        field := *graphql.GetPathContext(ctx).Field
        return nil, fmt.Errorf("Type unknwown for field '%s'", field)
    }
    if l < *n {
        field := *graphql.GetPathContext(ctx).Field
        return nil, fmt.Errorf("`%s' to short. Minimum length is '%d'", field, *n)
    }
    return data, err
}

//inputMaxLength the that the size of the field is stricly greater than the given value
func maxLength(ctx context.Context, obj interface{}, next graphql.Resolver, f *string, e []model.TensionEvent, n *int) (interface{}, error) {
    var l int
    data, err := next(ctx)
    if err != nil { return nil, err }

    switch d := data.(type) {
    case *string:
        l = len(*d)
    case string:
        l = len(d)
    default:
        field := *graphql.GetPathContext(ctx).Field
        return nil, fmt.Errorf("Type unknwown for field %s", field)
    }
    if l > *n {
        field := *graphql.GetPathContext(ctx).Field
        return nil, fmt.Errorf("`%s' to short. Maximum length is %d", field, *n)
    }
    return data, err
}
