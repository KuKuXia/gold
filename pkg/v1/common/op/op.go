package op

import (
	g "gorgonia.org/gorgonia"
)

// Clip the node value.
func Clip(value *g.Node, min, max float64) (retVal *g.Node, err error) {
	minNode := g.NewScalar(value.Graph(), g.Float64, g.WithValue(min), g.WithName("clip_min"))
	maxNode := g.NewScalar(value.Graph(), g.Float64, g.WithValue(max), g.WithName("clip_max"))

	// check if its the min value.
	minMask, err := g.Lt(value, minNode, true)
	if err != nil {
		return nil, err
	}
	minVal, err := g.HadamardProd(minNode, minMask)
	if err != nil {
		return nil, err
	}

	// check if its the given value.
	isMaskGt, err := g.Gt(value, minNode, true)
	if err != nil {
		return nil, err
	}
	isMaskLt, err := g.Lt(value, maxNode, true)
	if err != nil {
		return nil, err
	}
	isMask, err := g.HadamardProd(isMaskGt, isMaskLt)
	if err != nil {
		return nil, err
	}
	isVal, err := g.HadamardProd(value, isMask)
	if err != nil {
		return nil, err
	}

	// check if its the max value.
	maxMask, err := g.Gt(value, maxNode, true)
	if err != nil {
		return nil, err
	}
	maxVal, err := g.HadamardProd(maxNode, maxMask)
	if err != nil {
		return nil, err
	}

	return g.ReduceAdd(g.Nodes{minVal, isVal, maxVal})
}

// Min value between the nodes. If values are equal the first value is returned.
func Min(a *g.Node, b *g.Node) (retVal *g.Node, err error) {
	aMask, err := g.Lte(a, b, true)
	if err != nil {
		return nil, err
	}
	aVal, err := g.HadamardProd(a, aMask)
	if err != nil {
		return nil, err
	}

	bMask, err := g.Lt(b, a, true)
	if err != nil {
		return nil, err
	}
	bVal, err := g.HadamardProd(b, bMask)
	if err != nil {
		return nil, err
	}
	return g.Add(aVal, bVal)
}

// Max value between the nodes. If values are equal the first value is returned.
func Max(a *g.Node, b *g.Node) (retVal *g.Node, err error) {
	aMask, err := g.Gte(a, b, true)
	if err != nil {
		return nil, err
	}
	aVal, err := g.HadamardProd(a, aMask)
	if err != nil {
		return nil, err
	}

	bMask, err := g.Gt(b, a, true)
	if err != nil {
		return nil, err
	}
	bVal, err := g.HadamardProd(b, bMask)
	if err != nil {
		return nil, err
	}
	return g.Add(aVal, bVal)
}