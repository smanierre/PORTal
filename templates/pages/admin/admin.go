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
	DisabledMembers  bool
}

type MemberEditorData struct {
	SelectedMember     types.Member
	Ranks              map[types.Grade]string
	SupervisorListData SupervisorListData
	NewMember          bool
}

type SupervisorListData struct {
	PotentialSupervisors []types.Member
	SelectedMember       types.Member
}
