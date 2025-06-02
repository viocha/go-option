package util

import (
	"errors"
	"fmt"
	
	"github.com/viocha/go-option/internal/common"
)

var (
	ErrMust = fmt.Errorf("expected panic occurred") // 用于标识预期的 panic，会被捕获
)

// 将一个值包装成ErrMustFnPanic错误，并返回错误
func WrapMust(v any) error {
	return errors.Join(ErrMust, common.ToError(v))
}

// 强制获取值，如果有错误则 panic
func MustGet[T any](v T, err error) T {
	if err != nil {
		panic(common.WrapSub(err, ErrMust, "panic in MustGet"))
	}
	return v
}

// 强制获取两个值，如果有错误则 panic
func MustGet2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(common.WrapSub(err, ErrMust, "panic in MustGet2"))
	}
	return v1, v2
}

// 捕获所有 panic，如果存在，则转换成 ErrMust 错误，然后panic
func WrapPanic(f func()) {
	err := common.SafeDo(f)
	if err != nil {
		if !errors.Is(err, ErrMust) {
			err = WrapMust(err)
		}
		panic(err) // 如果发生错误，panic
	}
}

// 捕获所有 panic，如果存在，则转换成 ErrMust 错误，然后panic。否则返回原来的结果
func WrapPanicGet[T any](f func() T) T {
	var result T
	if err := common.SafeDo(func() {
		result = f()
	}); err != nil {
		if !errors.Is(err, ErrMust) {
			err = WrapMust(err)
		}
		panic(err) // 如果发生错误，panic
	} else {
		return result // 如果没有错误，返回结果
	}
}
