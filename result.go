package option

import (
	"errors"
	"fmt"
	"reflect"
)

// Result[T] 是一个用于包装值和错误的类型。它要么是 Ok(T)（成功并包含一个值），要么是 Err(error)（错误并包含一个错误值）
type Result[T any] struct {
	val *T
	err error
}

// 构造一个 Result[T] 的 Ok 变体
func Ok[T any](value T) Result[T] {
	return Result[T]{val: &value, err: nil}
}

// 构造一个 Result[T] 的 Err 变体
func Err[T any](err error) Result[T] {
	if err == nil {
		panic("Err() called with nil error")
	}
	return Result[T]{val: nil, err: err}
}

func (r Result[T]) From(val T ,err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(val)
}

func (r Result[T]) value() T { return *r.val }

func (r Result[T]) String() string {
	if r.IsOk() {
		return fmt.Sprintf("Ok[%T](%v)", r.value(), r.value())
	}
	typ := reflect.TypeFor[T]()
	return fmt.Sprintf("Err[%v](%v)", typ, r.err)
}

// 如果 Result 是 Ok 值，则返回 true
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// 如果 Result 是 Err 值，则返回 true
func (r Result[T]) IsErr() bool {
	return !r.IsOk()
}

// 当 Result 是 Ok 且内部值等于v时返回 true
func (r Result[T]) Has(v T) bool {
	return r.IsOk() && reflect.DeepEqual(r.value(), v)
}

// 当 Result 是 Ok 且内部值满足f时返回 true
func (r Result[T]) HasFunc(f func(T) bool) bool {
	return r.IsOk() && f(r.value())
}

// 当 Result 是 Err 且内部错误等于e时返回 true
func (r Result[T]) HasErr(e error) bool {
	return !r.IsOk() && errors.Is(r.err, e)
}

// 当 Result 是 Err 且内部错误满足f时返回 true。
func (r Result[T]) HasErrFunc(f func(error) bool) bool {
	return !r.IsOk() && f(r.err)
}

// Do 如果 Result 是 Ok，则对其包含的值调用 f。
func (r Result[T]) Do(f func(T)) Result[T] {
	if r.IsOk() {
		f(r.value())
	}
	return r
}

// ElseDo 如果 Result 是 Err，则对其包含的错误调用 f。
func (r Result[T]) ElseDo(f func(error)) Result[T] {
	if !r.IsOk() {
		f(r.err)
	}
	return r
}

// 如果Result是Ok，则返回原来的Result，否则返回一个新的Result
func (r Result[T]) Or(b Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	return b
}

// 如果Result是Ok，则返回原来的Result，否则调用f并返回
func (r Result[T]) OrFunc(f func(error) Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	return f(r.err)
}

// 如果是Ok，则返回原来的 Result，否则使用f转换错误，并返回一个新的 Result
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if !r.IsOk() {
		return Err[T](f(r.err))
	}
	return r
}

// =========================== 获取值或错误 ============================

// 如果 Result 是 Ok，则返回其包含的值。否则 panic
func (r Result[T]) Get() T {
	if !r.IsOk() {
		panic(fmt.Sprintf("called Result.Unwrap() on an Err value: %v", r.err))
	}
	return r.value()
}

// 如果 Result 是 Ok，则返回其包含的值。否则返回默认值 v
func (r Result[T]) GetOr(v T) T { return r.Val().GetOr(v) }

// 如果 Result 是 Ok，则返回其包含的值。否则返回类型的零值
func (r Result[T]) GetOrZero() T { return r.Val().GetOrZero() }

// 如果 Result 是 Ok，则返回其包含的值。否则调用f并返回其结果
func (r Result[T]) GetOrFunc(f func(error) T) T {
	if r.IsOk() {
		return r.value()
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

// 同时返回值和错误
func (r Result[T]) GetWithErr() (T, error) {
	if r.IsOk() {
		return r.value(), nil
	}
	return *new(T), r.err
}

// ========================== 和 Option 转换 ============================

// 将 Result[T] 的值转换为 Option[T]
func (r Result[T]) Val() Option[T] {
	if r.IsOk() {
		return Some(r.value())
	}
	return None[T]()
}

// 将 Result[T] 的错误转换为 Option[error]
func (r Result[T]) Err() Option[error] {
	if !r.IsOk() {
		return Some(r.err)
	}
	return None[error]()
}

// ========================== 逻辑与 ============================

// 如果 Result 是 Ok，则返回 b，否则返回原来的Err
func RAnd[T any, U any](a Result[T], b Result[U]) Result[U] {
	if a.IsOk() {
		return b
	}
	return Err[U](a.err)
}

// 如果 Result 是 Ok，则调用f得到一个新的Result并返回，否则返回原来的Err
func RAndFunc[T any, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.IsOk() {
		return f(r.value())
	}
	return Err[U](r.err)
}

// ==========================  Map操作 ============================
// 如果 Result 是 Ok，则使用f转换其值构造一个新的 Result 并返回，否则返回原来的Err
func RMap[T any, U any](r Result[T], f func(T) U) Result[U] {
	if r.IsOk() {
		return Ok(f(r.value()))
	}
	return Err[U](r.err)
}

// 如果 Result 是 Ok，则使用f转换其值并返回，否则返回默认值 v
func RMapOr[T any, U any](r Result[T], f func(T) U, v U) U {
	if r.IsOk() {
		return f(r.value())
	}
	return v
}

// 如果 Result 是 Ok，则使用okFn并返回其结果，否则调用errFn并返回其结果
func RMapOrFunc[T any, U any](r Result[T], okFn func(T) U, errFn func(error) U) U {
	if r.IsOk() {
		return okFn(r.value())
	}
	return errFn(r.err)
}
