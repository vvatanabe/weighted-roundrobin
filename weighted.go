package weighted

import (
	"math/big"
	"sync"
)

type Node struct {
	Value  interface{}
	Weight uint64
}

func New(nodes []*Node) *RoundRobbin {
	m := make(map[int]*Node)
	for i, v := range nodes {
		m[i] = v
	}

	var rr RoundRobbin
	rr.nodes = m
	rr.lastNodeIndex = -1
	if len(m) == 0 {
		return &rr
	}
	rr.currentNodeWeight = m[0].Weight

	var weights []uint64
	for _, v := range m {
		weights = append(weights, v.Weight)
	}
	rr.weightGCD = calcGCD(weights...)
	return &rr
}

type RoundRobbin struct {
	mu                sync.Mutex
	nodes             map[int]*Node
	lastNodeIndex     int
	currentNodeWeight uint64
	weightGCD         uint64
}

func (rr *RoundRobbin) SetNode(index int, node *Node) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	rr.nodes[index] = node
	rr.weightGCD = calcGCD(rr.getWeights()...)
}

func (rr *RoundRobbin) DeleteNode(index int) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	delete(rr.nodes, index)
	rr.weightGCD = calcGCD(rr.getWeights()...)
}

func (rr *RoundRobbin) GetNode() *Node {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for {
		nodes := len(rr.nodes)
		if nodes == 0 {
			return nil
		}
		rr.lastNodeIndex = (rr.lastNodeIndex + 1) % len(rr.nodes)

		if rr.lastNodeIndex == 0 {
			rr.currentNodeWeight = rr.currentNodeWeight - rr.weightGCD

			if rr.currentNodeWeight <= 0 {
				rr.currentNodeWeight = rr.getMaxWeight()

				if rr.currentNodeWeight == 0 {
					return nil
				}
			}
		}
		if weight := rr.nodes[rr.lastNodeIndex].Weight; weight >= rr.currentNodeWeight {
			return rr.nodes[rr.lastNodeIndex]
		}
	}
}

func (rr *RoundRobbin) Size() int {
	return len(rr.nodes)
}

func calcGCD(values ...uint64) uint64 {
	z := values[0]
	for _, n := range values {
		z = gcd(n, z)
	}
	return z
}

func gcd(m, n uint64) uint64 {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)
	a := new(big.Int).SetUint64(m)
	b := new(big.Int).SetUint64(n)
	return z.GCD(x, y, a, b).Uint64()
}

func (rr *RoundRobbin) getMaxWeight() uint64 {
	var max uint64
	for _, v := range rr.nodes {
		if v.Weight >= max {
			max = v.Weight
		}
	}

	return max
}

func (rr *RoundRobbin) getWeights() []uint64 {
	var weights []uint64
	for _, v := range rr.nodes {
		weights = append(weights, v.Weight)
	}
	return weights
}
