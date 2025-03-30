package serverutils

import (
	"PORTal/types"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

var domain string

// SetDomain should be called once on the initial setup of the server. Subsequent calls after it has been set will have no effect.
func SetDomain(d string) {
	if domain == "" {
		domain = d
	}
}

func CheckHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") != ""
}

func GetSessionId(r *http.Request) (string, error) {
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func GetMemberFromContext(ctx context.Context) (types.Member, error) {
	mVal := ctx.Value(MemberContextKey)
	m, ok := mVal.(types.Member)
	if !ok {
		return types.Member{}, errors.New("unable to get member from context")
	}
	return m, nil
}

func HandleRenderError(ctx context.Context, logger *slog.Logger, err error) {
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "Error rendering template", slog.String("error", err.Error()))
	}
}

func MakeCookie(name, value string, expiration time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   domain,
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func RemoveCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Domain:  domain,
		Expires: time.Now(),
		Path:    "/",
	})
}
