package selector

import (
	"math/rand"
	"time"
)

type randomBalancer struct {
}

func newRandomBalancer() *randomBalancer {
	return &randomBalancer{}
}

func (r *randomBalancer) Balance(serviceName string, nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix())
	num := rand.Intn(len(nodes))
	return nodes[num]
}
