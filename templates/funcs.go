package templates

var funcMap = map[string]map[string]interface{}{
	"nav": map[string]interface{}{
		"getNavItems": func() []NavItem { return nil },
	},
}
