package backend

import (
	"PORTal/types"
	"context"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func (b Backend) Login(username, password string) (types.Member, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Attempting to login member", slog.String("username", username))
	member, err := b.memberProvider.GetMember(username, ByUsername)
	if err != nil {
		return types.Member{}, ErrAuthenticationFailed
	}
	if err = bcrypt.CompareHashAndPassword([]byte(member.Hash), []byte(password)); err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Password validation failed")
		return types.Member{}, ErrAuthenticationFailed
	}
	return member, nil
}
