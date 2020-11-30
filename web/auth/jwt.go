package auth

import (
    //"fmt"
    "log"
    "time"
    "errors"
    "context"
    "encoding/json"
    "net/http"
    //"github.com/mitchellh/mapstructure"
    jwt "github.com/dgrijalva/jwt-go"

    . "zerogov/fractal6.go/tools"
    "zerogov/fractal6.go/web/middleware/jwtauth"
	"zerogov/fractal6.go/graph/model"
    "zerogov/fractal6.go/db"
)

var tkMaster *Jwt
var buildMode string

func init () {
    tkMaster = Jwt{}.New()

    // Get env mode
    if buildMode == "" {
        buildMode = "DEV"
    } else {
        buildMode = "PROD"
    }
}

type Jwt struct {
    // @FIX: How to initialize mapClaims with another map
    // in order to node decode evething at each request
	tokenAuth  *jwtauth.JWTAuth
    tokenClaim string
    tokenClaimErr string
}

// New create a token auth master
func (Jwt) New() *Jwt {
    secret := "frctl6"
	tk := &Jwt{
        tokenClaim: "user_ctx",
        tokenClaimErr: "user_ctx_err",
		tokenAuth: jwtauth.New("HS256", []byte(secret), nil),
	}
    uctx := model.UserCtx{
        Username: "yoda",
        Rights: model.UserRights{CanLogin:false, CanCreateRoot:true},
        Roles: []model.Role{
            {Rootnameid:"SKU", Nameid:"SKU", RoleType:model.RoleTypeCoordinator},
            {Rootnameid:"open-chaos", Nameid:"open-chaos", RoleType:model.RoleTypeCoordinator},
        },
    }
    token, _ := tk.issue(uctx, time.Hour*72)
	log.Println("DEBUG JWT:", Unpack64(token))
	return tk
}

func (tk Jwt) GetAuth() *jwtauth.JWTAuth {
    return tk.tokenAuth
}

// Issue generate and encode a new token
func (tk *Jwt) issue(d model.UserCtx, t time.Duration) (string, error) {
    claims := jwt.MapClaims{ tk.tokenClaim: d }
    jwtauth.SetIssuedNow(claims)
    jwtauth.SetExpiry(claims, time.Now().UTC().Add(t))
	_, tokenString, err := tk.tokenAuth.Encode(claims)
	return Pack64(tokenString), err
}

//
// Global functions
//

func GetTokenMaster() *Jwt {
    return tkMaster
}

// NewUserToken create a new user token from master key
func NewUserToken(userCtx model.UserCtx) (string, error) {
    var token string
    var err error
    if buildMode == "PROD" {
        token, err = tkMaster.issue(userCtx, time.Hour*24*30)
    } else {
        token, err = tkMaster.issue(userCtx, time.Hour*24)
        //token, err = tkMaster.issue(userCtx, time.Second*30)
    }
    return token, err
}

// NexuserCookie create an http cookie that embed a token
func NewUserCookie(userCtx model.UserCtx) (*http.Cookie, error) {
    tokenString, err := NewUserToken(userCtx)
    if err != nil {
        return nil, err
    }
    var httpCookie http.Cookie
    if buildMode == "PROD" {
        httpCookie = http.Cookie{
            Name: "jwt",
            Value: tokenString,
            Path: "/",
            HttpOnly: true,
            Secure: true,
            SameSite: 2, // https://golang.org/src/net/http/cookie.go
            //Expires: expirationTime,
            //MaxAge: 90000,
        }
    } else {
        httpCookie = http.Cookie{
            Name: "jwt",
            Value: tokenString,
            Path: "/",
            Secure: true,
            SameSite: 2,
        }
    }

    return &httpCookie, nil
}

func ContextWithUserCtx(ctx context.Context) context.Context {
    token, claims, err := jwtauth.FromContext(ctx)
    if err != nil {
        //errMsg := fmt.Errorf("%v", err)
        switch err {
        case jwtauth.ErrUnauthorized:
        case jwtauth.ErrExpired:
        case jwtauth.ErrNBFInvalid:
        case jwtauth.ErrIATInvalid:
        case jwtauth.ErrNoTokenFound:
        case jwtauth.ErrAlgoInvalid:
        }
        //http.Error(w, http.StatusText(401), 401)
        //return
    } else if token == nil || !token.Valid {
        err = errors.New("jwtauth: token is invalid")
    } else if claims[tkMaster.tokenClaim] == nil {
        err = errors.New("auth: user claim is invalid")
    }

    if err != nil { // Set the user error token
        ctx = context.WithValue(ctx, tkMaster.tokenClaimErr, err)
    }

    // Set the user token
    userCtx := model.UserCtx{}
    uRaw, err := json.Marshal(claims[tkMaster.tokenClaim])
    if err != nil { panic(err) }
    json.Unmarshal(uRaw, &userCtx)
    ctx = context.WithValue(ctx, tkMaster.tokenClaim, userCtx)

    // Set the Iat
    if claims["iat"] != nil {
        var iat int64 = int64(claims["iat"].(float64))
        return context.WithValue(ctx, "iat", time.Unix(iat, 0).Format(time.RFC3339))
    }
    return ctx
}

func UserCtxFromContext(ctx context.Context) (model.UserCtx, error) {
    uctx := ctx.Value(tkMaster.tokenClaim).(model.UserCtx)
    userCtxErr := ctx.Value(tkMaster.tokenClaimErr)
    if userCtxErr != nil { return uctx, userCtxErr.(error) }

    uctx.Iat = ctx.Value("iat").(string)
    return uctx, nil
}

// CheckUserCtxIat update the user token if the given
// node has been updated after that the token's iat.
func CheckUserCtxIat(uctx model.UserCtx, nid string) (model.UserCtx, error) {
    var u *model.UserCtx
    var e error
    var updatedAt string
    DB := db.GetDB()
    updatedAt_, e := DB.GetFieldByEq("Node.nameid", nid, "Node.updatedAt")
    if e != nil { return uctx, e }
    if updatedAt_ == nil {
        updatedAt = uctx.Iat
    } else {
        updatedAt = updatedAt_.(string)
    }
    if IsOlder(uctx.Iat, updatedAt) {
        u, e = DB.GetUser("username", uctx.Username)
        if e != nil { return uctx, e }
        return *u, e
    }
    return uctx, e
}
