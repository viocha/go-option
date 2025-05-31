package common

import "fmt"

// 如果传入的值是 error 类型，则直接返回该 error，否则将其包装为 error 类型
func ToError(r any) error {
	if err, ok := r.(error); ok {
		return err
	}
	return fmt.Errorf("%v", r)
}

// 可以从panic中安全地执行函数，返回是否成功执行
func DoSafe(f func()) error {
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = ToError(r)
			}
		}()
		f()
	}()
	return err
}
