package server

import (
	"PORTal/types"
	"context"
	"errors"
)

type ContextKey string

const (
	MemberContextKey ContextKey = "member"
)

func MemberFromContext(ctx context.Context) (types.Member, error) {
	mVal := ctx.Value(MemberContextKey)
	m, ok := mVal.(types.Member)
	if !ok {
		return types.Member{}, errors.New("unable to get member from context")
	}
	return m, nil
}
