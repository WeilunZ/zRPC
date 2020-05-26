package selector

import (
	"sync"
	"time"
)

type consistentHashBalancer struct {
	pickers  *sync.Map
	duration time.Duration
}
