package capabilities

import (
	"fmt"

	"github.com/ai-hippo/constgen/internal/model"
)

// ProfileResolver resolves constants from profiles.
//
// ARCHITECTURE CHANGE:
//
// We no longer resolve from:
// values → profiles
//
// We now resolve directly from:
//
// profiles → constants
//
type ProfileResolver struct {

	Spec    *model.Spec
	Profile string
}

// Resolve returns final constant value.
func (r *ProfileResolver) Resolve(key string) (*model.ResolvedValue, error) {

	profile, ok := r.Spec.Profiles[r.Profile]
	if !ok {
		return nil, fmt.Errorf(
			"profile not found: %s",
			r.Profile,
		)
	}

	value, ok := profile[key]
	if !ok {
		return nil, fmt.Errorf(
			"constant '%s' not found in profile '%s'",
			key,
			r.Profile,
		)
	}

	return &model.ResolvedValue{
		Key:   key,
		Value: value,
	}, nil
}