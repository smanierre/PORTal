package core

import "PORTal/types"

func (c Core) UpdateMember(m types.Member) (types.Member, error) {
	return c.memberStore.UpdateMember(m)
}
