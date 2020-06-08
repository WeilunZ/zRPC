package selector

type Node struct {
	Key    string
	Value  []byte
	weight int
	Hash   string
}

type Balancer interface {
	Balance(serviceName string, nodes []*Node) *Node
}

var (
	balancerMap                = make(map[string]Balancer, 0)
	DefaultLoadBalancer        = newRandomBalancer()
	RoundRobinBalancer         = newRoundRobinBalancer()
	WeightedRoundRobinBalancer = newWeightedRoundRobinBalancer()
)

const (
	Random             = "random"
	RoundRobin         = "roundRobin"
	WeightedRoundRobin = "weightedRoundRobin"
	ConsistentHash     = "consistentHash"
)

func init() {
	RegisterBalancer(Random, DefaultLoadBalancer)
	RegisterBalancer(RoundRobin, RoundRobinBalancer)
	RegisterBalancer(WeightedRoundRobin, WeightedRoundRobinBalancer)
}

func RegisterBalancer(name string, balancer Balancer) {
	if balancerMap == nil {
		balancerMap = make(map[string]Balancer)
	}
	balancerMap[name] = balancer
}

func GetBalancer(name string) Balancer {
	if balancer, ok := balancerMap[name]; ok {
		return balancer
	}
	return DefaultLoadBalancer
}
