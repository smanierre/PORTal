package admin

import (
	"PORTal/templates/components/admin"
	"context"
	"fmt"
	"github.com/a-h/templ"
	"log/slog"
	"net/http"
)

type PageBuilder struct {
	logger            *slog.Logger
	rootComponent     Component
	rootComponentData any
	oobComponents     map[Component]any
}

type Component int

const (
	List Component = iota
)

func (p *PageBuilder) Render(ctx context.Context, w http.ResponseWriter) error {
	var component templ.Component
	switch p.rootComponent {
	case List:
		if v, ok := p.rootComponentData.(admin.ListData); !ok {
			return fmt.Errorf("invalid type for List component data, expected admin.ListData")
		} else {
			component = admin.SearchableList(v)
		}
	default:
		return fmt.Errorf("invalid root component type: %T", p.rootComponent)
	}
	err := component.Render(ctx, w)
	if err != nil {
		return err
	}
	return nil
}

func (p *PageBuilder) SetRootComponent(component Component, componentData any) *PageBuilder {
	switch component {
	case List:
		p.rootComponentData = componentData
		p.rootComponent = component
	}
	return p
}

func (p *PageBuilder) AddOobComponent(component Component, componentData any) error {
	switch component {
	case List:
		p.oobComponents[component] = componentData
	}
	return nil
}
