package snowflake

import (
	"errors"
	"sync/atomic"
	"time"
)

const (
	timestampBits = 41
	nodeBits      = 10
	sequenceBits  = 12

	maxNode     = -1 ^ (-1 << nodeBits)     // 1023
	maxSequence = -1 ^ (-1 << sequenceBits) // 4095

	timeShift       = sequenceBits + nodeBits // 22
	nodeShift       = sequenceBits            // 12
	epoch     int64 = 1633024800000
)

type Node struct {
	nodeID      int64
	lastTime    int64
	sequence    int64
	timeBackoff bool
}

func NewNode(nodeID int64, timeBackoff bool) (*Node, error) {
	if nodeID < 0 || nodeID > maxNode {
		return nil, errors.New("node ID is out of range")
	}
	return &Node{
		nodeID:      nodeID,
		timeBackoff: timeBackoff,
	}, nil
}

func (n *Node) Generate() int64 {
	for {
		now := time.Now().UnixMilli()
		last := atomic.LoadInt64(&n.lastTime)
		seq := atomic.LoadInt64(&n.sequence)

		if now < last {
			if n.timeBackoff {
				time.Sleep(time.Duration(last-now) * time.Millisecond)
				continue
			}
			return -1
		}

		if now == last {
			seq = (seq + 1) & maxSequence
			if seq == 0 {
				for now <= last {
					now = time.Now().UnixMilli()
				}
			}
		} else {
			seq = 0
		}

		if atomic.CompareAndSwapInt64(&n.lastTime, last, now) &&
			atomic.CompareAndSwapInt64(&n.sequence, seq-1, seq) {
			return (now-epoch)<<timeShift | n.nodeID<<nodeShift | seq
		}
	}
}

func ParseID(id int64) (time.Time, int64, int64) {
	timestamp := (id >> timeShift) + epoch
	nodeID := (id >> nodeShift) & maxNode
	sequence := id & maxSequence
	return time.UnixMilli(timestamp), nodeID, sequence
}
