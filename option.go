package option

import (
	"fmt"
	"reflect"
	
	"github.com/viocha/go-option/util"
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

func FromFunc[T any](f func() T) Option[T] {
	var result Option[T]
	if nil == util.DoWithPanic(func() {
		result = Val(f())
	}) {
		return result
	}
	return Nul[T]()
}

// ========================== 方法 =============================

// 存在值
func (o Option[T]) IsVal() bool {
	return o.exists
}

// 不存在值
func (o Option[T]) IsNul() bool {
	return !o.exists
}

func (o Option[T]) String() string {
	if o.IsVal() {
		return fmt.Sprintf("Some[%T](%v)", o.Get(), o.Get())
	}
	typ := reflect.TypeFor[T]()
	return fmt.Sprintf("None[%v]()", typ)
}

// 判断是否存在指定值，使用 reflect.DeepEqual 进行比较
func (o Option[T]) Has(value T) bool {
	if o.IsNul() {
		return false
	}
	return reflect.DeepEqual(o.Get(), value)
}

// 判断是否存在满足条件的值
func (o Option[T]) HasFunc(f func(T) bool) bool {
	return o.IsVal() && f(o.Get())
}

// ============================= 获取值或 error ================================

// 如果存在值，则返回该值。否则 panic。
func (o Option[T]) Get() T {
	if o.IsNul() {
		panic("called Option.Unwrap() on a None value")
	}
	return *o.val
}

func (o Option[T]) GetOr(value T) T {
	if o.IsVal() {
		return o.Get()
	}
	return value
}

func (o Option[T]) GetOrFunc(f func() T) T {
	if o.IsVal() {
		return o.Get()
	}
	return f()
}

func (o Option[T]) GetOrZero() T {
	if o.IsVal() {
		return o.Get()
	}
	return *new(T)
}

func (o Option[T]) ToPtr() *T {
	if o.IsVal() {
		return o.val
	}
	return nil
}

func (o Option[T]) ToErr(e error) error {
	if o.IsVal() {
		return nil
	}
	return e
}

func (o Option[T]) Unwrap(err error) (T, error) {
	if o.IsVal() {
		return o.Get(), nil
	}
	return *new(T), err
}

// ============================= 链式方法 ================================

func (o Option[T]) Try(f func(T)) Option[T] {
	if o.IsNul() {
		return o
	}
	if nil == util.DoWithPanic(func() {
		f(o.Get())
	}) {
		return o
	}
	return Nul[T]()
}

func (o Option[T]) Catch(f func()) Option[T] {
	if o.IsVal() {
		return o
	}
	if nil == util.DoWithPanic(f) {
		return o
	}
	return Nul[T]()
}

func (o Option[T]) Finally(f func()) Option[T] {
	if nil == util.DoWithPanic(f) {
		return o
	}
	return Nul[T]()
}

func (o Option[T]) Filter(f func(T) bool) Option[T] {
	if o.IsNul() {
		return Nul[T]()
	}
	result := o
	if nil == util.DoWithPanic(func() {
		if !f(o.Get()) {
			result = Nul[T]()
		}
	}) {
		return result
	}
	return Nul[T]()
}

// 不存在值时，执行给定的函数构造一个新的 Option
func (o Option[T]) Else(f func() Option[T]) Option[T] {
	if o.IsVal() {
		return o
	}
	var result Option[T]
	if nil == util.DoWithPanic(func() {
		result = f()
	}) {
		return result
	}
	return Nul[T]()
}

func (o Option[T]) ElseVal(f func() T) Option[T] {
	if o.IsVal() {
		return o
	}
	var result Option[T]
	if nil == util.DoWithPanic(func() {
		result = Val(f())
	}) {
		return result
	}
	return Nul[T]()
}

// ================================ 常用类型的逻辑与方法 =============================

func (o Option[T]) ThenT(f func(T) Option[T]) Option[T]                 { return Then(o, f) }
func (o Option[T]) ThenInt(f func(T) Option[int]) Option[int]           { return Then(o, f) }
func (o Option[T]) ThenFloat(f func(T) Option[float64]) Option[float64] { return Then(o, f) }
func (o Option[T]) ThenStr(f func(T) Option[string]) Option[string]     { return Then(o, f) }
func (o Option[T]) ThenBool(f func(T) Option[bool]) Option[bool]        { return Then(o, f) }

// ================================ 常用类型的Map方法 =============================

func (o Option[T]) MapT(f func(T) T) Option[T]                 { return Map(o, f) }
func (o Option[T]) MapInt(f func(T) int) Option[int]           { return Map(o, f) }
func (o Option[T]) MapFloat(f func(T) float64) Option[float64] { return Map(o, f) }
func (o Option[T]) MapStr(f func(T) string) Option[string]     { return Map(o, f) }
func (o Option[T]) MapBool(f func(T) bool) Option[bool]        { return Map(o, f) }

// ================================ 逻辑与  =============================

// 如果存在值，使用f构造一个新的 Option
func Then[T any, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.IsNul() {
		return Nul[U]()
	}
	var result Option[U]
	if nil == util.DoWithPanic(func() {
		result = f(o.Get())
	}) {
		return result
	}
	return Nul[U]()
}

// =============================== Map操作 =============================

// 若存在值，则使用f转换该值，构造一个 Option
func Map[T any, U any](o Option[T], f func(T) U) Option[U] {
	if o.IsNul() {
		return Nul[U]()
	}
	var result Option[U]
	if nil == util.DoWithPanic(func() {
		result = Val(f(o.Get()))
	}) {
		return result
	}
	return Nul[U]()
}

// 若存在值，则使用f转换该值并返回，否则返回给定的默认值
func MapOr[T any, U any](o Option[T], f func(T) U, v U) U {
	var result U
	if o.IsNul() {
		return v
	}
	if nil == util.DoWithPanic(func() {
		result = f(o.Get())
	}) {
		return result
	}
	return v
}

// 若存在值，则使用okFn转换该值并返回，否则调用defaultFn并返回其结果，defaultFn 中的panic不会被捕获
func MapOrFunc[T any, U any](o Option[T], okFn func(T) U, defaultFn func() U) U {
	if o.IsNul() {
		return defaultFn()
	}
	var result U
	if nil == util.DoWithPanic(func() {
		result = okFn(o.Get())
	}) {
		return result
	}
	return defaultFn()
}
