package user

type Scope string

const (
	ScopeFamily       Scope = "FAMILY"
	ScopeUser         Scope = "USER"
	ScopeSubscription Scope = "SUBSCRIPTION"
	ScopeProvider     Scope = "PROVIDER"
	ScopeLabel        Scope = "LABEL"
)
