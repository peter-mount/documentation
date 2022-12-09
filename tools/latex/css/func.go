package css

import (
	"context"
	"github.com/peter-mount/documentation/tools/latex/util"
)

var funcMap = map[string]NodeAction{
	"first-child": firstChild,
	"last-child":  lastChild,
	"not":         not,
}

func nop(_ context.Context, _ *Node) (*util.Value, error) {
	return nil, nil
}

// Negate the result
func not(ctx context.Context, n *Node) (*util.Value, error) {
	v, err := n.Left.Do(ctx)
	if err != nil {
		return nil, err
	}

	return util.BoolValue(!v.Bool()), nil
}

func and(ctx context.Context, n *Node) (*util.Value, error) {
	// Eval lhs
	ar, err := n.Left.Do(ctx)
	if err != nil {
		return nil, err
	}

	// If LHS is false then stop, same as a&&b, b not even if a is false
	if !ar.Bool() {
		return util.False(), nil
	}

	// Eval rhs
	br, err := n.Right.Do(ctx)
	if err != nil {
		return nil, err
	}

	// LHS already true so result is RHS
	return br.AsBool(), nil
}

func firstChild(_ context.Context, n *Node) (*util.Value, error) {
	return util.BoolValue(false), nil
}

func lastChild(_ context.Context, n *Node) (*util.Value, error) {
	return util.BoolValue(false), nil
}
