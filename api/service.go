// Package api provides types and interfaces that define how the Regula API is working.
package api

import (
	"context"

	"github.com/heetch/regula"
	"github.com/heetch/regula/errors"
	"github.com/heetch/regula/rule"
)

// API errors.
const (
	ErrRulesetNotFound      = errors.Error("ruleset not found")
	ErrRulesetNotModified   = errors.Error("not modified")
	ErrSignatureNotFound    = errors.Error("signature not found")
	ErrInvalidContinueToken = errors.Error("invalid continue token")
	ErrAlreadyExists        = errors.Error("already exists")
)

// RulesetService is a service managing rulesets.
type RulesetService interface {
	// Create a ruleset using a signature.
	Create(ctx context.Context, path string, signature *regula.Signature) error
	// Put is used to add a new version of the rules to a ruleset.
	Put(ctx context.Context, path string, rules []*rule.Rule) (string, error)
	// Get returns a ruleset alongside its metadata.
	Get(ctx context.Context, path string) (*regula.Ruleset, error)
	// List returns the latest version of each ruleset whose path starts by the given prefix.
	// If the prefix is empty, it returns all the entries following the lexical order.
	// The listing is paginated and can be customised using the ListOptions type.
	List(ctx context.Context, prefix string, opt *ListOptions) (*Rulesets, error)
	// Watch a prefix for changes and return a list of events.
	Watch(ctx context.Context, prefix string, revision string) (*RulesetEvents, error)
	// Eval evaluates a ruleset given a path and a set of parameters. It implements the regula.Evaluator interface.
	Eval(ctx context.Context, path, version string, params rule.Params) (*regula.EvalResult, error)
}

// ListOptions contains list options.
// If the Limit is lower or equal to 0 or greater than 100, it will be set to 50 by default.
type ListOptions struct {
	Limit         int
	ContinueToken string
	PathsOnly     bool // return only the paths of the rulesets
	AllVersions   bool // return all versions of each rulesets
}

// Rulesets holds a list of rulesets.
type Rulesets struct {
	Rulesets []regula.Ruleset `json:"rulesets"`
	Revision string           `json:"revision"`           // revision when the request was applied
	Continue string           `json:"continue,omitempty"` // token of the next page, if any
}

// List of possible events executed against a ruleset.
const (
	RulesetPutEvent = "PUT"
)

// RulesetEvent describes an event that occured on a ruleset.
type RulesetEvent struct {
	Type    string       `json:"type"`
	Path    string       `json:"path"`
	Version string       `json:"version"`
	Rules   []*rule.Rule `json:"rules"`
}

// RulesetEvents holds a list of events occured on a group of rulesets.
type RulesetEvents struct {
	Events   []RulesetEvent
	Revision string
	Timeout  bool // indicates if the watch did timeout
}
