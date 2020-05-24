package selector

import (
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	Key    string
	Value  []byte
	Weight int
}

type Balancer interface {
	Balance(serviceName string, nodes []*Node) *Node
}

var DefaultLoadBalancer = newRandomBalancer()

func newRandomBalancer() *randomBalancer {
	return &randomBalancer{}
}

type randomBalancer struct {
}

func (r *randomBalancer) Balance(serviceName string, nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix())
	num := rand.Intn(len(nodes))
	return nodes[num]
}

type roundRobinBalancer struct {
	pickers  *sync.Map
	duration time.Duration // time duration to update again
}

type roundRobinPicker struct {
	length         int           // service nodes length
	lastUpdateTime time.Time     // last update time
	duration       time.Duration // time duration to update again
	lastIndex      int           // last accessed index
}

func newRoundRobinBalancer() *roundRobinBalancer {
	return &roundRobinBalancer{
		pickers:  new(sync.Map),
		duration: 3 * time.Minute,
	}
}

func (rp *roundRobinPicker) pick(nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}

	// update picker after timeout
	if time.Now().Sub(rp.lastUpdateTime) > rp.duration ||
		len(nodes) != rp.length {
		rp.length = len(nodes)
		rp.lastUpdateTime = time.Now()
		rp.lastIndex = 0
	}

	if rp.lastIndex == len(nodes)-1 {
		rp.lastIndex = 0
		return nodes[0]
	}

	rp.lastIndex += 1
	return nodes[rp.lastIndex]
}

func (r *roundRobinBalancer) Balance(serviceName string, nodes []*Node) *Node {
	var picker *roundRobinPicker

	if p, ok := r.pickers.Load(serviceName); !ok {
		picker = &roundRobinPicker{
			lastUpdateTime: time.Now(),
			duration:       r.duration,
			length:         len(nodes),
		}
	} else {
		picker = p.(*roundRobinPicker)
	}

	node := picker.pick(nodes)
	r.pickers.Store(serviceName, picker)
	return node
}
