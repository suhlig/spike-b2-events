package b2

import (
	"fmt"
	"strconv"
	"time"
)

type Timestamp time.Time

func (ts Timestamp) MarshalJSON() ([]byte, error) {
	data := []byte(fmt.Sprint(time.Time(ts).UTC().UnixMilli()))
	return data, nil
}

func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	millisecondsSinceEpoch, e := strconv.ParseInt(string(data), 10, 64)

	secondsSinceEpoch := millisecondsSinceEpoch / 1000
	nanoSecondsSinceEpoch := millisecondsSinceEpoch % 1000 * 1000 * 1000

	*ts = Timestamp(time.Unix(secondsSinceEpoch, nanoSecondsSinceEpoch).UTC())
	return e
}

func (ts Timestamp) String() string {
	return time.Time(ts).UTC().Format(time.RFC3339Nano)
}
