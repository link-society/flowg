package hlc

import (
	"fmt"

	"math"

	"sync"
	"time"
)

type Timestamp struct {
	WallTime int64  `json:"wall"`
	Logical  uint32 `json:"logical"`
	NodeID   string `json:"node"`
}

func (t Timestamp) Compare(other Timestamp) int {
	switch {
	case t.WallTime < other.WallTime:
		return -1
	case t.WallTime > other.WallTime:
		return 1
	case t.Logical < other.Logical:
		return -1
	case t.Logical > other.Logical:
		return 1
	case t.NodeID < other.NodeID:
		return -1
	case t.NodeID > other.NodeID:
		return 1
	default:
		return 0
	}
}

func (t Timestamp) After(other Timestamp) bool  { return t.Compare(other) > 0 }
func (t Timestamp) Before(other Timestamp) bool { return t.Compare(other) < 0 }
func (t Timestamp) Equal(other Timestamp) bool  { return t.Compare(other) == 0 }
func (t Timestamp) IsZero() bool                { return t == Timestamp{} }

func (t Timestamp) String() string {
	return fmt.Sprintf("%d.%d@%s", t.WallTime, t.Logical, t.NodeID)
}

type Clock struct {
	nodeID string
	now    func() int64

	mu   sync.Mutex
	last Timestamp
}

func NewClock(nodeID string) *Clock {
	return &Clock{
		nodeID: nodeID,
		now:    func() int64 { return time.Now().UnixNano() },
	}
}

func bumpLogical(wall int64, logical uint32) (int64, uint32) {
	if logical == math.MaxUint32 {
		return wall + 1, 0
	}
	return wall, logical + 1
}

func (c *Clock) Now() Timestamp {
	c.mu.Lock()
	defer c.mu.Unlock()

	physical := c.now()

	if physical > c.last.WallTime {
		c.last = Timestamp{WallTime: physical, Logical: 0, NodeID: c.nodeID}
	} else {
		wall, logical := bumpLogical(c.last.WallTime, c.last.Logical)
		c.last = Timestamp{WallTime: wall, Logical: logical, NodeID: c.nodeID}
	}

	return c.last
}

func (c *Clock) Update(remote Timestamp) Timestamp {
	c.mu.Lock()
	defer c.mu.Unlock()

	physical := c.now()
	wall := max(c.last.WallTime, remote.WallTime, physical)

	var logical uint32
	switch {
	case wall == c.last.WallTime && wall == remote.WallTime:
		wall, logical = bumpLogical(wall, max(c.last.Logical, remote.Logical))
	case wall == c.last.WallTime:
		wall, logical = bumpLogical(wall, c.last.Logical)
	case wall == remote.WallTime:
		wall, logical = bumpLogical(wall, remote.Logical)
	default:
		logical = 0
	}

	c.last = Timestamp{WallTime: wall, Logical: logical, NodeID: c.nodeID}
	return c.last
}
