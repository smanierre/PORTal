package admin

import (
	"PORTal/templates/components"
	"PORTal/types"
)

type Data struct {
	DropdownData components.DropdownData
	MembersData  MembersData
}

type MembersData struct {
	Members          []types.Member
	SelectedMember   types.Member
	SwapTarget       string
	FragmentBasePath string
	MemberEditorData MemberEditorData
	OobSwap          bool
}

type MemberEditorData struct {
	Members              []types.Member
	SelectedMember       types.Member
	Ranks                map[types.Grade]string
	PotentialSupervisors []types.Member
}
