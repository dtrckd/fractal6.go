package db

import (
    "fmt"
    "log"
    "bytes"
    "strings"
    "context"
    "reflect"
    "net/http"
    "encoding/json"
    "text/template"
    "github.com/spf13/viper"
    "github.com/mitchellh/mapstructure"
    //"github.com/vektah/gqlparser/v2/gqlerror"
    "github.com/dgraph-io/dgo/v200"
    "github.com/dgraph-io/dgo/v200/protos/api"
    "google.golang.org/grpc"

    "zerogov/fractal6.go/graph/model"
    "zerogov/fractal6.go/tools"
)

// Draph database clients
type Dgraph struct {
    // HTTP/Graphql and GPRC/DQL client address
    gqlUrl string
    grpcUrl string

    // Auth token
    gqlToken string

    // HTTP/Graphql and GPRC/DQL client template
    gqlTemplates map[string]*QueryString
    dqlTemplates map[string]*QueryString
}

//
// GRPC/Graphql+- response
//

type DqlResp struct {
    All []map[string]interface{} `json:"all"`
}

type DqlRespCount struct {
    All []map[string]int `json:"all"`
}

//
// HTTP/Graphql response
//

type GqlRes struct {
    Data   model.JsonAtom `json:"data"`
    Errors []model.JsonAtom `json:"errors"` // message, locations, path, extensions
}

type GraphQLError struct {
    msg string
}

func (e *GraphQLError) Error() string {
    return fmt.Sprintf("%s", e.msg)
}

//
// Query String Interface
//

type QueryString struct {
    Q string
    Template *template.Template
}

// Init clean the query to be compatible in application/json format.
func (q *QueryString) Init() {
    d := q.Q
    q.Q = tools.CleanString(d, false)
    // Load the template @DEBUG: Do we need a template name ?
    q.Template = template.Must(template.New("graphql").Parse(q.Q))
}

func (q QueryString) Format(maps map[string]string) string {
    buf := bytes.Buffer{}
    q.Template.Execute(&buf, maps)
    return buf.String()
}

//
// Initialization
//

// Database client
var DB *Dgraph

func init () {
    DB = initDB()
}

func GetDB() *Dgraph {
    return DB
}

func initDB() *Dgraph {
    tools.InitViper()
    HOSTDB := viper.GetString("db.host")
    PORTDB := viper.GetString("db.port_graphql")
    PORTGRPC := viper.GetString("db.port_grpc")
    APIDB := viper.GetString("db.api")
    dgraphApiUrl := "http://"+HOSTDB+":"+PORTDB+"/"+APIDB
    grpcUrl := HOSTDB+":"+PORTGRPC

    if HOSTDB == "" {
        panic("Viper error: not host found")
    } else {
        fmt.Println("Dgraph Graphql addr:", dgraphApiUrl)
        fmt.Println("Dgraph Grpc addr:", grpcUrl)
    }

    gqlToken := ""

    // HTTP/Graphql Request Template
    gqlQueries := map[string]string{
        // QUERIES
        "query": `{
            "query": "query {{.Args}} {{.QueryName}} {
                {{.QueryName}} {
                    {{.QueryGraph}}
                }
            }"
        }`,
        "rawQuery": `{
            "query": "{{.RawQuery}}"
        }`,

        // MUTATIONS
        "add": `{
            "query": "mutation {{.QueryName}}($input:[{{.InputType}}!]!) {
                {{.QueryName}}(input: $input) {
                    {{.QueryGraph}}
                }
            }",
            "variables": {
                "input": {{.InputPayload}}
            }
        }`,
        "update": `{
            "query": "mutation {{.QueryName}}($input:{{.InputType}}!) {
                {{.QueryName}}(input: $input) {
                    {{.QueryGraph}}
                }
            }",
            "variables": {
                "input": {{.InputPayload}}
            }
        }`,
        "delete": `{
            "query": "mutation {{.QueryName}}($input:{{.InputType}}!) {
                {{.QueryName}}(filter: $input) {
                    {{.QueryGraph}}
                }
            }",
            "variables": {
                "input": {{.InputPayload}}
            }
        }`,
    }

    dqlT := map[string]*QueryString{}
    gqlT := map[string]*QueryString{}

    for op, q := range(dqlQueries) {
        dqlT[op] = &QueryString{Q:q}
        dqlT[op].Init()
    }
    for op, q := range(gqlQueries) {
        gqlT[op] = &QueryString{Q:q}

        gqlT[op].Init()
    }

    return &Dgraph{
        gqlUrl: dgraphApiUrl,
        grpcUrl: grpcUrl,
        gqlToken: gqlToken,
        dqlTemplates: dqlT,
        gqlTemplates: gqlT,
    }
}

//
// Internals
//

func (dg Dgraph) getDqlQuery(op string, m map[string]string) string {
    var q string
    if _q, ok := dg.dqlTemplates[op]; ok {
        q = _q.Format(m)
    } else {
        panic("unknonw DQL query op: " + op)
    }
    return q
}

func (dg Dgraph) getGqlQuery(op string, m map[string]string) string {
    var q string
    if _q, ok := dg.gqlTemplates[op]; ok {
        q = _q.Format(m)
    } else {
        panic("unknonw GQL query op: " + op)
    }
    return q
}

// Get the grpc Dgraph client.
func (dg Dgraph) getDgraphClient() (dgClient *dgo.Dgraph, cancelFunc func()) {
    conn, err := grpc.Dial(dg.grpcUrl, grpc.WithInsecure())
    if err != nil {
        log.Fatal("While trying to dial gRPC")
    }

    dgClient = dgo.NewDgraphClient(api.NewDgraphClient(conn))
    //ctx := context.Background()

    //// Perform login call. If the Dgraph cluster does not have ACL and
    //// enterprise features enabled, this call should be skipped.
    //for {
    //	// Keep retrying until we succeed or receive a non-retriable error.
    //	err = dgClient.Login(ctx, "groot", "password")
    //	if err == nil || !strings.Contains(err.Error(), "Please retry") {
    //		break
    //	}
    //	time.Sleep(time.Second)
    //}
    //if err != nil {
    //	log.Fatalf("While trying to login %v", err.Error())
    //}

    cancelFunc =  func() {
        if err := conn.Close(); err != nil {
            log.Printf("Error while closing connection:%v", err)
        }
    }
    return
}

// Post send a post request to the Graphql client.
func (dg Dgraph) postql(data []byte, res interface{}) error {
    req, err := http.NewRequest("POST", dg.gqlUrl, bytes.NewBuffer(data))
    req.Header.Set("Content-Type", "application/json")

    // Set dgraph token
    req.Header.Set("X-FRAC6-AUTH", dg.gqlToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return json.NewDecoder(resp.Body).Decode(res)
}

//
// DQL (ex GraphQL+-) Interface
//


// QueryDql runs a query on dgraph (...QueryDql)
func (dg Dgraph) QueryDql(op string, maps map[string]string) (*api.Response, error) {
    // init client
    dgc, cancel := dg.getDgraphClient()
    defer cancel()
    ctx := context.Background()
    txn := dgc.NewTxn()
    defer txn.Discard(ctx)

    // Get the Query
    q := dg.getDqlQuery(op, maps)
    // Send Request
    fmt.Println(op)
    //fmt.Println(string(q))
    res, err := txn.Query(ctx, q)
    //fmt.Println(res)
    return res, err
}

//MutateWithQueryDql runs an upsert block mutation by first querying query
//and then mutate based on the result.
func (dg Dgraph) MutateWithQueryDql(query string, mu *api.Mutation) (error) {
    // init client
    dgc, cancel := dg.getDgraphClient()
    defer cancel()
    ctx := context.Background()
    txn := dgc.NewTxn()
    defer txn.Discard(ctx)

    req := &api.Request{
        Query: query,
        Mutations: []*api.Mutation{mu},
        CommitNow: true,
    }

    _, err := txn.Do(ctx, req)
    return err
}

//MutateUpsertDql adds a new object in the database if it doesn't exist
func (dg Dgraph) MutateUpsertDql_(object map[string]interface{}, dtype string, upsertField string, upsertVal string) error {
    // init client
    dgc, cancel := dg.getDgraphClient()
    defer cancel()
    ctx := context.Background()
    txn := dgc.NewTxn()
    defer txn.Discard(ctx)

    // make the template here.
    template := template.Must(template.New("graphql").Parse(`{
        all(func: eq({{.dtype}}.{{.upsertField}}, "{{.upserVal}}")) {
            v as uid
        }
    }`))
    buf := bytes.Buffer{}
    template.Execute(&buf, map[string]string{"dtype":dtype, "upsertField":upsertField, "upsertVal":upsertVal})
    query := buf.String()

    object["dgraph.type"] = []string{dtype}
    object["uid"] = "uid(v)"
    js, err := json.Marshal(object)
    if err != nil { return err }
    mu := &api.Mutation{SetJson: js}

    req := &api.Request{
        Query: query,
        Mutations: []*api.Mutation{mu},
        CommitNow: true,
    }
    // User do instead of Mutate here ?
    _, err = txn.Do(ctx, req)
    return err
}

//Push adds a new object in the database.
func (dg Dgraph) Push_(object map[string]interface{}, dtype string) (string, error) {
    // init client
    dgc, cancel := dg.getDgraphClient()
    defer cancel()
    ctx := context.Background()
    txn := dgc.NewTxn()
    defer txn.Discard(ctx)

    object["dgraph.type"] = []string{dtype}
    uid, ok := object["uid"]
    if !ok {
        uid = "_:new_obj"
        object["uid"] = uid
    }
    js, err := json.Marshal(object)
    if err != nil { return "", err }

    mu := &api.Mutation{
        CommitNow: true,
        SetJson: js,
    }
    r, err := txn.Mutate(ctx, mu)
    if err != nil { return "", err }

    uid = r.Uids[uid.(string)]
    return uid.(string), nil
}

//ClearNodes remove nodes from theirs uid and their edges in dgraph.
// refers to https://dgraph.io/docs/mutations/json-mutation-format/#deleting-edges
func (dg Dgraph) ClearNodes(uids ...string) (error) {
    // init client
    dgc, cancel := dg.getDgraphClient()
    defer cancel()
    ctx := context.Background()
    txn := dgc.NewTxn()
    defer txn.Discard(ctx)

    d :=  map[string]string{}
    for _, uid := range uids { d["uid"] = uid }
    js, err := json.Marshal(d)
    if err != nil { return err }

    mu := &api.Mutation{
        CommitNow: true,
        DeleteJson: js,
    }

    _, err = txn.Mutate(ctx, mu)
    return err
}

//ClearEdges remove edges from their uid
func (dg Dgraph) ClearEdges(key string, value string, delMap map[string]interface{}) (error) {
    // init client
    dgc, cancel := dg.getDgraphClient()
    defer cancel()
    ctx := context.Background()
    txn := dgc.NewTxn()
    defer txn.Discard(ctx)

    query := fmt.Sprintf(`query {
        o as var(func: eq(%s, "%s"))
    }`, key, value)

    var mu string
    for k, _ := range delMap {
        mu = mu + fmt.Sprintf(`uid(o) <%s> * .`, k) + "\n"
    }

    mutation := &api.Mutation{
        DelNquads: []byte(mu),
    }

    err := dg.MutateWithQueryDql(query, mutation)
    return err
}

//
// GraphQL Interface
//

// QueryGql query the Dgraph Graphql endpoint by following a http request.
// It map the result in to given data structure
func (dg Dgraph) QueryGql(op string, reqInput map[string]string, data interface{}) error {
    // Get the query
    queryName := reqInput["QueryName"]
    q := dg.getGqlQuery(op, reqInput)

    // Send the dgraph request and follow the results
    res := &GqlRes{}
    //fmt.Println("request ->", string(q))
    err := dg.postql([]byte(q), res)
    //fmt.Println("response ->", res)
    if err != nil {
        return err
    } else if res.Errors != nil {
        err, _ := json.Marshal(res.Errors)
        //return fmt.Errorf(string(err))
        return &GraphQLError{string(err)}
    }

    switch v := data.(type) {
    case model.JsonAtom:
        for k, val := range res.Data {
            v[k] = val
        }
    default: // Payload data type
        var config *mapstructure.DecoderConfig
        if op == "query" || op == "rawQuery" {
            // Decoder config to handle aliased request
            // @DEBUG: see bug #3c3f1f7
            config = &mapstructure.DecoderConfig{
                Result: data,
                TagName: "json",
                DecodeHook: func(from, to reflect.Kind, v interface{}) (interface{}, error) {
                    if to == reflect.Struct {
                        nv := tools.CleanAliasedMap(v.(map[string]interface{}))
                        return nv, nil
                    }
                    return v, nil
                },
            }
        } else {
            config = &mapstructure.DecoderConfig{TagName: "json", Result: data}
        }

        decoder, err := mapstructure.NewDecoder(config)
        if err != nil { return err }
        err = decoder.Decode(res.Data[queryName])
    }

    return err
}


//
// Graphql requests
//

// Add a new vertex
func (dg Dgraph) Add(vertex string, input interface{}) (string, error) {
    Vertex := strings.Title(vertex)
    queryName := "add" + Vertex
    inputType := "Add" + Vertex + "Input"
    queryGraph := vertex + ` { id }`

    // Just One Node
    var ifaces []interface{} = make([]interface{}, 1)
    ifaces[0] = input
    inputs, _ := json.Marshal(ifaces)

    // Build the string request
    reqInput := map[string]string{
        "QueryName": queryName, // function name (e.g addUser)
        "InputType": inputType, // input type name (e.g AddUserInput)
        "QueryGraph": tools.CleanString(queryGraph, true), // output data
        "InputPayload": string(inputs), // inputs data
    }

    // Send request
    payload := make(model.JsonAtom, 1)
    err := dg.QueryGql("add", reqInput, payload)
    if err != nil { return "", err }
    // Extract id result
    res := payload[queryName].(model.JsonAtom)[vertex].([]interface{})[0].(model.JsonAtom)["id"]
    return res.(string), err
}

// Update a vertex
func (dg Dgraph) Update(vertex string, input interface{}) error {
    Vertex := strings.Title(vertex)
    queryName := "update" + Vertex
    inputType := "Update" + Vertex + "Input"
    queryGraph := vertex + ` { id }`

    // Just One Node
    inputs, _ := json.Marshal(input)

    // Build the string request
    reqInput := map[string]string{
        "QueryName": queryName, // function name (e.g addUser)
        "InputType": inputType, // input type name (e.g AddUserInput)
        "QueryGraph": tools.CleanString(queryGraph, true), // output data
        "InputPayload": string(inputs), // inputs data
    }

    // Send request
    payload := make(model.JsonAtom, 1)
    err := dg.QueryGql("update", reqInput, payload)
    return err
}

// Delete a vertex
func (dg Dgraph) Delete(vertex string, input interface{}) error {
    Vertex := strings.Title(vertex)
    queryName := "delete" + Vertex
    inputType :=  Vertex + "Filter"
    queryGraph := vertex + ` { id }`

    // Just One Node
    inputs, _ := json.Marshal(input)

    // Build the string request
    reqInput := map[string]string{
        "QueryName": queryName, // function name (e.g addUser)
        "InputType": inputType, // input type name (e.g AddUserInput)
        "QueryGraph": tools.CleanString(queryGraph, true), // output data
        "InputPayload": string(inputs), // inputs data
    }

    // Send request
    payload := make(model.JsonAtom, 1)
    err := dg.QueryGql("delete", reqInput, payload)
    return err
}


