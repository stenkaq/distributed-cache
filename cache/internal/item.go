package cache

import "time"

type Data struct {
	Value    string
	ExpireAt time.Time
}
