package model

// ResolvedValue is the output of ProfileResolver.
// It carries BOTH value and type so generators don't guess.
type ResolvedValue struct {
	// Constant key.
	Key string
	// Actual value.
	Value any
}