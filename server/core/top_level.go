package core

import (
	"PORTal/types"
	"net/http"
)

type TopLevelData struct {
	HasSubordinates   bool
	ShowNav           bool
	Admin             bool
	MemberDisplayName string
	Organization      string
}

// TopLevel returns the data necessary to render the root template. This is used when a page is navigated to directly so HTMX isn't loaded yet.
// If member is nil, it will attempt to get the member from the session cookie. If that fails, it will return backend.ErrSessionNotFound
func (c Core) TopLevel(r *http.Request, optionalMember *types.Member) (TopLevelData, error) {
	tld := TopLevelData{}
	//tld.Organization = c.Organization
	//if optionalMember == nil {
	//	sessionID, err := c.getSessionId(r)
	//	m, err := c.memberStore.GetMemberFromSession(sessionID)
	//	if err != nil {
	//		c.logger.LogAttrs(r.Context(), slog.LevelWarn, "Unable to find member from session")
	//		return tld, err
	//	}
	//	optionalMember = &m
	//}
	//if optionalMember != nil {
	//	tld.ShowNav = true
	//}
	//tld.Admin = optionalMember.Admin
	//tld.MemberDisplayName = optionalMember.Display(c.Service)
	//
	//subordinates := c.memberStore.GetSubordinates(optionalMember.ID)
	//tld.HasSubordinates = len(subordinates) > 0

	return tld, nil
}
