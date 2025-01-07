package admin

import (
	"PORTal/templates/components"
	"PORTal/types"
)

type ReferenceRootData struct {
	DropdownData   components.DropdownData
	ReferencesData ReferenceContentData
}
type ReferenceContentData struct {
	SelectedReference   types.Reference
	References          []types.Reference
	SwapTarget          string
	FragmentBasePath    string
	ReferenceEditorData ReferenceEditorData
}

type ReferenceEditorData struct{}
