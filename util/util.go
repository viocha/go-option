package util

import (
	"errors"
	"fmt"
)

var (
	ErrMust = fmt.Errorf("expected panic occurred") // 用于标识预期的 panic，会被捕获
)

// 如果传入的值是 error 类型，则直接返回该 error，否则将其包装为 error 类型
func ToError(r any) error {
	if err, ok := r.(error); ok {
		return err
	}
	return fmt.Errorf("%v", r)
}

// 可以从panic中安全地执行函数，返回是否成功执行。
// errs指定需要捕获的错误类型列表，默认捕获所有错误
func DoSafe(f func(), errs ...error) error {
	var err error
	deferFn := func() {
		if r := recover(); r != nil {
			if len(errs) == 0 { // 如果没有指定错误类型列表
				err = ToError(r) // 直接转换为 error
				return
			}
			// 检查是否是指定的错误类型
			if r, ok := r.(error); ok { // 如果 r 是 error 类型
				for _, e := range errs {
					if errors.Is(r, e) { // 如果 r 是指定的错误类型，则捕获并返回
						err = r
						return
					}
				}
				panic(r) // 如果不是指定的错误类型，则继续 panic
			}
			panic(r) // 如果 r 不是 error 类型，则继续 panic
		}
	}
	
	func() {
		defer deferFn()
		f()
	}()
	return err
}

// 将 [错误，格式化消息] 合并为一个错误
func WrapMsg(err error, format string, args ...any) error {
	return errors.Join(
		err,
		fmt.Errorf(format, args...),
	)
}

// 将 [父错误，格式化消息，子错误] 合并为一个错误，如果不需要格式化消息，可以直接使用 errors.Join(err, sub)
func WrapSub(sub error, parent error, format string, args ...any) error {
	return errors.Join(
		WrapMsg(parent, format, args...),
		sub,
	)
}

// 将一个值包装成ErrMustFnPanic错误，并返回错误
func WrapMust(v any) error {
	return errors.Join(ErrMust, ToError(v))
}

// 强制获取值，如果有错误则 panic
func MustGet[T any](v T, err error) T {
	if err != nil {
		panic(WrapSub(err, ErrMust, "panic in MustGet"))
	}
	return v
}

// 强制获取两个值，如果有错误则 panic
func MustGet2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(WrapSub(err, ErrMust, "panic in MustGet2"))
	}
	return v1, v2
}

// 捕获 ErrMust 错误的panic
func DoWithPanic(f func()) error {
	return DoSafe(f, ErrMust)
}
