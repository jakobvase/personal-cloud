package userStorage

import (
	"encoding/json"
	"time"
)

// UnixMillis is a time.Time that unmarshals from a Unix timestamp in milliseconds
type UnixMillis time.Time

func (t *UnixMillis) UnmarshalJSON(b []byte) error {
	var ms int64
	if err := json.Unmarshal(b, &ms); err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(0, ms*int64(time.Millisecond))
	return nil
}

func (t UnixMillis) MarshalJSON() ([]byte, error) {
	ms := time.Time(t).UnixNano() / int64(time.Millisecond)
	return json.Marshal(ms)
}