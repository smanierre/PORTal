package components

type DropdownData struct {
	Items            []DropdownItem
	SwapTarget       string
	FragmentBasePath string
}

type DropdownItem struct {
	DisplayName string
	Value       string
}
