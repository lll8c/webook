package domain

import "time"

// User 领域对象，是DDD中的entity
// 这个User是业务意义上的User
type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Nickname string
	// YYYY-MM-DD
	Birthday time.Time
	AboutMe  string

	Ctime time.Time
	Utime time.Time
	//不要组合，因为将来还有其他登录方式，里面可能有同名字段
	WechatInfo WechatInfo
}
