package option

import (
	"fmt"
	"reflect"
)

// Option 表示一个可选值：每个 Option[T] 要么是 Some 并包含一个值，要么是 None 且不包含值。
type Option[T any] struct {
	val    *T
	exists bool
}

// ========================== 构造函数 =============================

// Some 构造一个 Option[T] 的 Some 变体。
func Some[T any](value T) Option[T] {
	return Option[T]{val: &value, exists: true}
}

// None 构造一个 Option[T] 的 None 变体。
func None[T any]() Option[T] {
	return Option[T]{val: nil, exists: false}
}

func From[T any](val T, err error) Option[T] {
	if err != nil {
		return None[T]()
	}
	return Some(val)
}

func FromPtr[T any](val *T) Option[T] {
	if val == nil {
		return None[T]()
	}
	return Some(*val)
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
func (o Option[T]) IsSome() bool {
	return o.exists
}

// 不存在值
func (o Option[T]) IsNone() bool {
	return !o.exists
}

// 判断是否存在值并且值等于给定的值，使用 reflect.DeepEqual 进行比较
func (o Option[T]) Has(value T) bool {
	if !o.exists {
		return false
	}
	return reflect.DeepEqual(o.Get(), value)
}

// 判断是否存在值且满足给定的条件函数
func (o Option[T]) HasFunc(f func(T) bool) bool {
	return o.exists && f(o.Get())
}

// 存在值则对其包含的值调用 f，并返回原来的 Option
func (o Option[T]) Do(f func(T)) Option[T] {
	if o.exists {
		f(o.Get())
	}
	return o
}

// 不存在值则调用f
func (o Option[T]) ElseDo(f func()) {
	if !o.exists {
		f()
	}
}

// 如果存在值且满足条件f，则返回原来的 Option，否则返回 None
func (o Option[T]) Filter(f func(T) bool) Option[T] {
	if o.exists && f(o.Get()) {
		return o
	}
	return None[T]()
}

// 如果存在值，则返回原来的 Option，否则返回给定的 Option
func (o Option[T]) Or(b Option[T]) Option[T] {
	if o.exists {
		return o
	}
	return b
}

// 如果存在值，则返回原来的 Option，否则调用f并返回其结果
func (o Option[T]) OrFunc(f func() Option[T]) Option[T] {
	if o.exists {
		return o
	}
	return f()
}

// 仅当其中一个 Option 存在值时，返回该 Option，否则返回 None
func (o Option[T]) Xor(b Option[T]) Option[T] {
	if o.exists && !b.exists {
		return o
	}
	if !o.exists && b.exists {
		return b
	}
	return None[T]()
}

// ============================= 获取值或 error ================================

// 如果存在值，则返回该值。否则 panic。
func (o Option[T]) Get() T {
	if !o.exists {
		panic("called Option.Unwrap() on a None value")
	}
	return *o.val
}

// 存在值则返回该值，否则返回给定的值。
func (o Option[T]) GetOr(value T) T {
	if o.exists {
		return o.Get()
	}
	return value
}

// 存在值则返回该值，否则调用给定的函数并返回其结果。
func (o Option[T]) GetOrFunc(f func() T) T {
	if o.exists {
		return o.Get()
	}
	return f()
}

// 存在值则返回该值，否则返回类型 T 的零值
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

// 如果不存在值，则返回给定的错误，否则返回 nil
func (o Option[T]) ToErr(e error) error {
	if o.exists {
		return nil
	}
	return e
}

// 同时返回value和error
func (o Option[T]) ToValErr(err error) (T, error) {
	if o.exists {
		return o.Get(), nil
	}
	return *new(T), err
}

// ================================ 逻辑与  =============================

// 如果两个 Option 都存在值，则返回第二个 Option，否则返回 None
func And[T any, U any](a Option[T], b Option[U]) Option[U] {
	if a.exists && b.exists {
		return b
	}
	return None[U]()
}

// 如果存在值，使用f处理该值并返回新的 Option
func AndFunc[T any, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.exists {
		return f(o.Get())
	}
	return None[U]()
}

// =============================== Map操作 =============================

// 若存在值，则使用f转换该值构造一个 Option 并返回
func Map[T any, U any](o Option[T], f func(T) U) Option[U] {
	if o.exists {
		return Some(f(o.Get()))
	}
	return None[U]()
}

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
