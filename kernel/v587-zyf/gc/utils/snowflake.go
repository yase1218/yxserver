package utils

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

var (
	// Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
	// You may customize this to set a different epoch for your application.
	// Epoch int64 = 1420070400000

	// 2020-01-01 08:00:00 UTC/GMT  in millseconds
	Epoch int64 = 1577836800000

	// NodeBits holds the number of bits to use for Snowflake
	// Remember, you have a total 22 bits to share between Snowflake/Step
	NodeBits uint8 = 14

	// StepBits holds the number of bits to use for Step
	// Remember, you have a total 22 bits to share between Snowflake/Step
	StepBits uint8 = 9

	mu        sync.Mutex
	nodeMax   int64 = -1 ^ (-1 << NodeBits) // = 1<<NodeBits - 1
	nodeMask        = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift       = NodeBits + StepBits
	nodeShift       = StepBits
)

type Snowflake struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
}

type ID uint64

// NewNode returns a new snowflake node that can be used to generate snowflake
// IDs
func NewSnowflake(node int64) (*Snowflake, error) {
	// re-calc in case custom NodeBits or StepBits were set
	// DEPRECATED: the below block will be removed in a future release.
	mu.Lock()
	nodeMax = -1 ^ (-1 << NodeBits)
	nodeMask = nodeMax << StepBits
	stepMask = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits
	nodeShift = StepBits
	mu.Unlock()

	n := Snowflake{}
	n.node = node
	n.nodeMax = -1 ^ (-1 << NodeBits)
	n.nodeMask = n.nodeMax << StepBits
	n.stepMask = -1 ^ (-1 << StepBits)
	n.timeShift = NodeBits + StepBits
	n.nodeShift = StepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Snowflake number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

// GenerateGoIds creates and returns a unique snowflake ID
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node ID
func (n *Snowflake) Generate() ID {

	n.mu.Lock()

	now := time.Since(n.epoch).Nanoseconds() / 1000000

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	r := ID((now)<<n.timeShift |
		(n.node << n.nodeShift) |
		(n.step),
	)

	n.mu.Unlock()
	return r
}

// Int64 returns an int64 of the snowflake ID
func (f ID) Uint64() uint64 {
	return uint64(f)
}

// ParseInt64 converts an int64 into a snowflake ID
func ParseInt64(id int64) ID {
	return ID(id)
}

func ParseUint64(id uint64) ID {
	return ID(id)
}

// ParseString converts a string into a snowflake ID
func ParseString(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return ID(i), err
}

// ParseBase2 converts a Base2 string into a snowflake ID
func ParseBase2(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 2, 64)
	return ID(i), err
}

// String returns a string of the snowflake ID
func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

// Base2 returns a string base2 of the snowflake ID
func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

// Time returns an int64 unix timestamp in milliseconds of the snowflake ID time
func (f ID) TimeUnix() int64 {
	return (int64(f) >> timeShift) + Epoch
}

func (f ID) TimeSecond() time.Time {
	return time.Unix(f.TimeUnix()/1000, 0)
}

func (f ID) Time() time.Time {
	return time.Unix(0, f.TimeUnix()*1000000)
}

// Snowflake returns an int64 of the snowflake ID node number
// DEPRECATED: the below function will be removed in a future release.
func (f ID) Node() int64 {
	return int64(f) & nodeMask >> nodeShift
}

// Step returns an int64 of the snowflake step (or sequence) number
// DEPRECATED: the below function will be removed in a future release.
func (f ID) Step() int64 {
	return int64(f) & stepMask
}
