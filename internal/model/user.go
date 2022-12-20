package model

import "time"

type User struct {
	Id       int64
	UserName string
	Address  string
}

type SpaceInfo struct {
	Id      int64
	Owner   string
	Total   int64
	Used    int64
	Rest    time.Duration
	Buckets []string
	State   string
}
