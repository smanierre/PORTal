package templates

import "PORTal/types"

type NavData struct {
	Show         bool
	OobSwap      bool
	Member       types.Member
	Subordinates bool
}

func (n NavData) GetNavItems() []NavItem {
	// All users get a dashboard
	items := []NavItem{
		{
			Href:    "/dashboard",
			Display: "Dashboard",
		},
	}
	if n.Subordinates {
		items = append(items, NavItem{
			Href:    "/members",
			Display: "Members",
		})
	}
	if n.Member.Admin {
		items = append(items, NavItem{
			Href:    "/admin",
			Display: "Admin",
		})
	}
	return items
}

type NavItem struct {
	Href    string
	Display string
}
