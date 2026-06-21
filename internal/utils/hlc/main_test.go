package hlc

import "testing"

func mockClock(nodeID string, t *int64) *Clock {
	c := NewClock(nodeID)
	c.now = func() int64 { return *t }
	return c
}

func TestNow_LogicalIncrementsWhenWallStalls(t *testing.T) {
	wall := int64(100)
	c := mockClock("a", &wall)

	t1 := c.Now()
	t2 := c.Now()

	if t1.WallTime != 100 || t1.Logical != 0 {
		t.Fatalf("t1 = %s; want 100.0", t1)
	}
	if t2.WallTime != 100 || t2.Logical != 1 {
		t.Fatalf("t2 = %s; want 100.1", t2)
	}
	if !t2.After(t1) {
		t.Fatalf("expected t2 > t1")
	}
}

func TestNow_LogicalResetsWhenWallAdvances(t *testing.T) {
	wall := int64(100)
	c := mockClock("a", &wall)

	c.Now()
	c.Now()
	wall = 200
	t3 := c.Now()

	if t3.WallTime != 200 || t3.Logical != 0 {
		t.Fatalf("t3 = %s; want 200.0", t3)
	}
}

func TestNow_MonotonicWhenWallGoesBackwards(t *testing.T) {
	wall := int64(200)
	c := mockClock("a", &wall)

	t1 := c.Now()
	wall = 150
	t2 := c.Now()

	if !t2.After(t1) {
		t.Fatalf("expected monotonic: t2 %s should be after t1 %s", t2, t1)
	}
	if t2.WallTime != 200 || t2.Logical != 1 {
		t.Fatalf("t2 = %s; want 200.1", t2)
	}
}

func TestUpdate_AdoptsRemoteWall(t *testing.T) {
	wall := int64(100)
	c := mockClock("a", &wall)

	got := c.Update(Timestamp{WallTime: 500, Logical: 7, NodeID: "b"})

	if got.WallTime != 500 || got.Logical != 8 {
		t.Fatalf("got = %s; want 500.8", got)
	}
	if got.NodeID != "a" {
		t.Fatalf("got.NodeID = %s; want a", got.NodeID)
	}
}

func TestUpdate_TakesMaxLogicalOnWallTie(t *testing.T) {
	wall := int64(300)
	c := mockClock("a", &wall)

	c.Now()
	c.Now() // last = 300.1

	got := c.Update(Timestamp{WallTime: 300, Logical: 5, NodeID: "b"})

	if got.WallTime != 300 || got.Logical != 6 {
		t.Fatalf("got = %s; want 300.6", got)
	}
}

func TestUpdate_PhysicalWinsAndResetsLogical(t *testing.T) {
	wall := int64(1000)
	c := mockClock("a", &wall)

	got := c.Update(Timestamp{WallTime: 500, Logical: 9, NodeID: "b"})

	if got.WallTime != 1000 || got.Logical != 0 {
		t.Fatalf("got = %s; want 1000.0", got)
	}
}

func TestCompare_NodeIDTiebreak(t *testing.T) {
	a := Timestamp{WallTime: 100, Logical: 2, NodeID: "node-a"}
	b := Timestamp{WallTime: 100, Logical: 2, NodeID: "node-b"}

	if a.Compare(b) >= 0 {
		t.Fatalf("expected node-a < node-b on tiebreak")
	}
	if !b.After(a) {
		t.Fatalf("expected b after a")
	}
	if !a.Equal(a) {
		t.Fatalf("expected a equal to itself")
	}
}

func TestCausality_RemoteThenLocalAreOrdered(t *testing.T) {
	wallA := int64(100)
	a := mockClock("a", &wallA)

	wallB := int64(100)
	b := mockClock("b", &wallB)

	// b sends an event, a receives it, then a emits a new event.
	sent := b.Now()
	a.Update(sent)
	after := a.Now()

	if !after.After(sent) {
		t.Fatalf("causality violated: %s should be after %s", after, sent)
	}
}
