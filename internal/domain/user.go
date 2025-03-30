package domain

import "time"

// User 领域对象, DDD中的entity
type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Ctime    time.Time
	Utime    time.Time
}
