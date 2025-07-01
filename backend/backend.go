package backend

const MinimumPwLength = 8

type ProviderMethod int

const (
	ById ProviderMethod = iota
	ByUsername
)
