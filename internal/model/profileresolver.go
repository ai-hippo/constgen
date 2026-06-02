package model

import "fmt"

// ProfileResolver resolves constants directly from profiles.
type ProfileResolver struct {

	// validated spec
	Spec *Spec

	// active profile
	//
	// Example:
	// string
	// integer
	//
	Profile string
}

// NewProfileResolver creates resolver instance.
func NewProfileResolver(
	spec *Spec,
	profile string,
) *ProfileResolver {

	return &ProfileResolver{
		Spec:    spec,
		Profile: profile,
	}
}

// Resolve fetches raw value directly from profile.
//
// Example:
//
// Resolve("YES")
// → "YES"
//
// Resolve("ENABLED")
// → true
//
func (r *ProfileResolver) Resolve(key string) (any, error) {

	// locate profile
	profile, ok := r.Spec.Profiles[r.Profile]
	if !ok {

		return nil, fmt.Errorf(
			"profile not found: %s",
			r.Profile,
		)
	}

	// locate constant
	value, ok := profile[key]
	if !ok {

		return nil, fmt.Errorf(
			"constant '%s' not found in profile '%s'",
			key,
			r.Profile,
		)
	}

	return value, nil
}