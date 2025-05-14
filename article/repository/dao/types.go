package dao

import (
	"errors"
)

var ErrPossibleIncorrectAuthor = errors.New("用户在尝试操作非本人数据")
