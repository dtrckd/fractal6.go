package graph

import (
    //"fmt"
    "zerogov/fractal6.go/graph/model"
    "zerogov/fractal6.go/web/auth"
)

// getRoles returns the list of the users roles below the given node
func getRoles(uctx model.UserCtx, rootnameid string) []model.Role {
    uctx, e := auth.CheckUserCtxIat(uctx, rootnameid)
    if e != nil { panic(e) }

    var roles []model.Role
    for _, r := range uctx.Roles {
        if r.Rootnameid == rootnameid  {
            roles = append(roles, r)
        }
    }

    return roles
}

// usePlayRole return true if the user play the given role (Nameid)
func userPlayRole(uctx model.UserCtx, nameid string) bool {
    uctx, e := auth.CheckUserCtxIat(uctx, nameid)
    if e != nil { panic(e) }

    for _, ur := range uctx.Roles {
        if ur.Nameid == nameid  {
            return true
        }
    }
    return false
}

// useHasRoot return true if the user belongs to the given root
func userHasRoot(uctx model.UserCtx, rootnameid string) bool {
    uctx, e := auth.CheckUserCtxIat(uctx, rootnameid)
    if e != nil { panic(e) }

    for _, ur := range uctx.Roles {
        if ur.Rootnameid == rootnameid {
            return true
        }
    }
    return false
}

// userIsGuest return true if the user is a guest (has only one role) in the given organisation
func userIsGuest(uctx model.UserCtx, rootnameid string) int {
    uctx, e := auth.CheckUserCtxIat(uctx, rootnameid)
    if e != nil { panic(e) }

    for i, r := range uctx.Roles {
        if r.Rootnameid == rootnameid && r.RoleType == model.RoleTypeGuest {
            return i
        }
    }

    return -1
}

// useIsMember return true if the user has at least one role in the given node
func userIsMember(uctx model.UserCtx, nameid string) int {
    uctx, e := auth.CheckUserCtxIat(uctx, nameid)
    if e != nil { panic(e) }

    for i, ur := range uctx.Roles {
        pid, err := nid2pid(ur.Nameid)
        if err != nil {
            panic(err.Error())
        }
        if pid == nameid {
            return i
        }
    }
    return -1
}

// useIsCoordo return true if the user has at least one role of Coordinator in the given node
func userIsCoordo(uctx model.UserCtx, nameid string) int {
    uctx, e := auth.CheckUserCtxIat(uctx, nameid)
    if e != nil { panic(e) }

    for i, ur := range uctx.Roles {
        pid, err := nid2pid(ur.Nameid)
        if err != nil {
            panic("bad nameid format for coordo test: "+ ur.Nameid)
        }
        if pid == nameid && ur.RoleType == model.RoleTypeCoordinator {
            return i
        }
    }

    return -1
}

func userIsOwner(uctx model.UserCtx, rootnameid string) int {
    uctx, e := auth.CheckUserCtxIat(uctx, rootnameid)
    if e != nil { panic(e) }

    for i, r := range uctx.Roles {
        if r.Rootnameid == rootnameid && r.RoleType == model.RoleTypeOwner {
            return i
        }
    }

    return -1
}

