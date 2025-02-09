/*
 * Fractale - Self-organisation for humans.
 * Copyright (C) 2022 Fractale Co
 *
 * This file is part of Fractale.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Fractale.  If not, see <http://www.gnu.org/licenses/>.
 */

package graph

import (
    //"fmt"
    "context"
    "strings"
    "strconv"
    "github.com/99designs/gqlgen/graphql"
    "fractale/fractal6.go/graph/model"
    "fractale/fractal6.go/web/auth"
	. "fractale/fractal6.go/tools"
)



func addReactionInputHook(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
    // Get User context
    ctx, uctx, err := auth.GetUserContext(ctx)
    if err != nil { return nil, LogErr("Access denied", err) }

    data, err := next(ctx)
    if err != nil {
        return data, err
    }

    // Move pendingCandidate to candidate if email exists in User
    newData := data.([]*model.AddReactionInput)
    for i, input := range newData {
        ids := []string{uctx.Username, *input.Comment.ID, strconv.Itoa(input.Type)}
        reactionid := strings.Join(ids, "#")
        input.Reactionid = reactionid
        newData[i] = input
    }

    return newData, err
}
