package selector

import (
	"sync"
	"time"
)

// TODO
type consistentHashBalancer struct {
	pickers  *sync.Map
	duration time.Duration
}
