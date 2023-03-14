package snowflake

import (
	"sync"
	"time"
)

// start time is 22:22:22 2/22/2022 (it's also a tuesday)
const start int64 = 1645568542000

var mut = sync.Mutex{}
var sequence uint16 = 0

var previous int64 = 0

type Snowflake int64

func Next() Snowflake {
	t := time.Now().UnixMilli() - start
	for t == previous && sequence == 0xFFFF {
		t = time.Now().UnixMilli() - start
	}

	mut.Lock()
	defer mut.Unlock()
	id := t | (int64(sequence) << 48)

	if previous == t {
		sequence++
	} else {
		sequence = 0
	}
	previous = t
	return Snowflake(id)
}

func (s Snowflake) Time() time.Time {
	return time.UnixMilli(start + int64(s&0x0000FFFFFFFFFFFF))
}

func (s Snowflake) Sequence() uint16 {
	return uint16(s >> 48)
}
