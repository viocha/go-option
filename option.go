package option

import (
	"fmt"
	"reflect"
)

type Option[T any] struct {
	val    *T
	exists bool
}

// ========================== 构造函数 =============================

func Val[T any](value T) Option[T] {
	return Option[T]{val: &value, exists: true}
}

func Nul[T any]() Option[T] {
	return Option[T]{val: nil, exists: false}
}

func From[T any](val T, err error) Option[T] {
	if err != nil {
		return Nul[T]()
	}
	return Val(val)
}

func FromPtr[T any](val *T) Option[T] {
	if val == nil {
		return Nul[T]()
	}
	return Val(*val)
}

// ========================== 方法 =============================

func (o Option[T]) String() string {
	if o.exists {
		return fmt.Sprintf("Some[%T](%v)", o.Get(), o.Get())
	}
	typ := reflect.TypeFor[T]()
	return fmt.Sprintf("None[%v]()", typ)
}

// 存在值
func (o Option[T]) IsVal() bool {
	return o.exists
}

// 不存在值
func (o Option[T]) IsNul() bool {
	return !o.exists
}

// 判断是否存在指定值，使用 reflect.DeepEqual 进行比较
func (o Option[T]) Has(value T) bool {
	if !o.exists {
		return false
	}
	return reflect.DeepEqual(o.Get(), value)
}

// 判断是否存在满足条件的值
func (o Option[T]) HasFunc(f func(T) bool) bool {
	return o.exists && f(o.Get())
}

func (o Option[T]) Try(f func(T)) Option[T] {
	if o.exists {
		f(o.Get())
	}
	return o
}

func (o Option[T]) Catch(f func()) {
	if !o.exists {
		f()
	}
}

func (o Option[T]) Finally(f func()) Option[T] {
	f()
	return o
}

func (o Option[T]) Filter(f func(T) bool) Option[T] {
	if o.exists && f(o.Get()) {
		return o
	}
	return Nul[T]()
}

// 不存在值时，执行给定的函数构造一个新的 Option
func (o Option[T]) Else(f func() Option[T]) Option[T] {
	if o.exists {
		return o
	}
	return f()
}

// ============================= 获取值或 error ================================

// 如果存在值，则返回该值。否则 panic。
func (o Option[T]) Get() T {
	if !o.exists {
		panic("called Option.Unwrap() on a None value")
	}
	return *o.val
}

func (o Option[T]) GetOr(value T) T {
	if o.exists {
		return o.Get()
	}
	return value
}

func (o Option[T]) GetOrFunc(f func() T) T {
	if o.exists {
		return o.Get()
	}
	return f()
}

func (o Option[T]) GetOrZero() T {
	if o.exists {
		return o.Get()
	}
	return *new(T)
}

func (o Option[T]) ToPtr() *T {
	if o.exists {
		return o.val
	}
	return nil
}

func (o Option[T]) ToErr(e error) error {
	if o.exists {
		return nil
	}
	return e
}

func (o Option[T]) ToValErr(err error) (T, error) {
	if o.exists {
		return o.Get(), nil
	}
	return *new(T), err
}

// ================================ 逻辑与  =============================

// 如果存在值，使用f构造一个新的 Option
func Then[T any, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.exists {
		return f(o.Get())
	}
	return Nul[U]()
}

// =============================== Map操作 =============================

// 若存在值，则使用f转换该值，构造一个 Option
func Map[T any, U any](o Option[T], f func(T) U) Option[U] {
	if o.exists {
		return Val(f(o.Get()))
	}
	return Nul[U]()
}

// =============================== 带有默认值的Map操作 =============================

// 若存在值，则使用f转换该值并返回，否则返回给定的默认值
func MapOr[T any, U any](o Option[T], f func(T) U, v U) U {
	if o.exists {
		return f(o.Get())
	}
	return v
}

// 若存在值，则使用okFn转换该值并返回，否则调用defaultFn并返回其结果
func MapOrFunc[T any, U any](o Option[T], okFn func(T) U, defaultFn func() U) U {
	if o.exists {
		return okFn(o.Get())
	}
	return defaultFn()
}
