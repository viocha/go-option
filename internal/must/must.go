package must

import (
	"github.com/viocha/go-option/internal/common"
	"github.com/viocha/go-option/util"
)

// 捕获 ErrMust 错误的panic
func CatchMustPanic(f func()) error {
	return common.SafeDo(f, util.ErrMust)
}
