package result

import (
	"errors"
	"fmt"
	"reflect"
	
	opt "github.com/viocha/go-option"
	"github.com/viocha/go-option/common"
)

type Result[T any] struct {
	val *T
	err error
}

// ========================== 构造函数 =============================

func Ok[T any](value T) Result[T] {
	return Result[T]{val: &value, err: nil}
}

func Err[T any](err error) Result[T] {
	if err == nil {
		panic("Err() called with nil error")
	}
	return Result[T]{val: nil, err: err}
}

// 将 T 和 error 转换为 Result[T]
func From[T any](val T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(val)
}

// 将 Option[T] 和 error 转换为 Result[T]
func FromOption[T any](o opt.Option[T], err error) Result[T] {
	if o.IsVal() {
		return Ok[T](o.Get())
	}
	return Err[T](err)
}

// ========================== 方法 =============================

func (r Result[T]) String() string {
	if r.IsOk() {
		return fmt.Sprintf("Ok[%T](%v)", r.Get(), r.Get())
	}
	typ := reflect.TypeFor[T]()
	return fmt.Sprintf("Err[%v](%v)", typ, r.err)
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return !r.IsOk()
}

func (r Result[T]) Has(v T) bool {
	return r.IsOk() && reflect.DeepEqual(r.Get(), v)
}

func (r Result[T]) HasFunc(f func(T) bool) bool {
	return r.IsOk() && f(r.Get())
}

func (r Result[T]) HasErr(e error) bool {
	return !r.IsOk() && errors.Is(r.err, e)
}

func (r Result[T]) HasErrFunc(f func(error) bool) bool {
	return !r.IsOk() && f(r.err)
}

// ========================== 链式方法 ============================

func (r Result[T]) Try(f func(T)) Result[T] {
	if r.IsErr() {
		return r
	}
	if err := common.DoSafe(func() {
		f(r.Get())
	}); err != nil {
		return Err[T](err)
	}
	return r
}

func (r Result[T]) Catch(f func(error)) Result[T] {
	if r.IsOk() {
		return r
	}
	if err := common.DoSafe(func() {
		f(r.err)
	}); err != nil {
		return Err[T](err)
	}
	return r
}

func (r Result[T]) Finally(f func()) Result[T] {
	if err := common.DoSafe(f); err != nil {
		return Err[T](err)
	}
	return r
}

// 如果Result是Err，则调用f并返回一个新的Result[T]
func (r Result[T]) Else(f func(error) Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	var newResult Result[T]
	if err := common.DoSafe(func() {
		newResult = f(r.err)
	}); err != nil {
		return Err[T](err)
	}
	return newResult
}

func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.IsOk() {
		return r
	}
	var newResult Result[T]
	if err := common.DoSafe(func() {
		newResult = Err[T](f(r.err))
	}); err != nil {
		return Err[T](err)
	}
	return newResult
}

// =========================== 获取值或错误 ============================

// 如果 Result 是 Ok，则返回其包含的值。否则 panic
func (r Result[T]) Get() T {
	if !r.IsOk() {
		panic(fmt.Sprintf("called Result.Unwrap() on an Err value: %v", r.err))
	}
	return *r.val
}

func (r Result[T]) GetOr(v T) T { return r.Val().GetOr(v) }

func (r Result[T]) GetOrZero() T { return r.Val().GetOrZero() }

func (r Result[T]) GetOrFunc(f func(error) T) T {
	if r.IsOk() {
		return r.Get()
	}
	return f(r.err)
}

// 如果 Result 是 Err，则返回其包含的错误。否则 panic
func (r Result[T]) GetErr() error {
	if r.IsOk() {
		panic("called Result.UnwrapErr() on an Ok value")
	}
	return r.err
}

func (r Result[T]) GetValErr() (T, error) {
	if r.IsOk() {
		return r.Get(), nil
	}
	return *new(T), r.err
}

func (r Result[T]) ToPtr() *T {
	if r.IsOk() {
		return r.val
	}
	return nil
}

// ========================== 和 opt.Option 转换 ============================

func (r Result[T]) Val() opt.Option[T] {
	if r.IsOk() {
		return opt.Val(r.Get())
	}
	return opt.Nul[T]()
}

func (r Result[T]) Err() opt.Option[error] {
	if !r.IsOk() {
		return opt.Val(r.err)
	}
	return opt.Nul[error]()
}

// ========================== 逻辑与 ============================

// Ok时则调用f得到一个新的Result
func Then[T any, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.IsErr() {
		return Err[U](r.err)
	}
	var newResult Result[U]
	if err := common.DoSafe(func() {
		newResult = f(r.Get())
	}); err != nil {
		return Err[U](err)
	}
	return newResult
}

// ==========================  Map操作 ============================
// Ok时使用f转换其值，构造一个新的 Result
func Map[T any, U any](r Result[T], f func(T) U) Result[U] {
	if r.IsErr() {
		return Err[U](r.err)
	}
	var newResult Result[U]
	if err := common.DoSafe(func() {
		newResult = Ok(f(r.Get()))
	}); err != nil {
		return Err[U](err)
	}
	return newResult
}

// ==========================  带有默认值的Map操作 ============================

// Ok时则使用f转换其值并返回，否则返回默认值 v
func MapOr[T any, U any](r Result[T], f func(T) U, v U) U {
	if r.IsErr() {
		return v
	}
	var val U
	if err := common.DoSafe(func() {
		val = f(r.Get())
	}); err != nil {
		return v
	}
	return val
}

// Ok时调用okFn并返回其结果，否则调用errFn并返回其结果
func MapOrFunc[T any, U any](r Result[T], okFn func(T) U, errFn func(error) U) U {
	if r.IsErr() {
		return errFn(r.err)
	}
	var val U
	if err := common.DoSafe(func() {
		val = okFn(r.Get())
	}); err != nil {
		return errFn(err)
	}
	return val
}
