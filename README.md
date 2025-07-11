# `option` 包

`option` 是一个受 Rust 启发的 Go 泛型库，提供 `Option[T]` 和 `Result[T]` 两个类型，用于处理可能不存在的值或错误，提升代码健壮性。

---

## 📦 安装

```bash
go get github.com/viocha/go-option
```

---

## 📄 示例用法

```go
package main

import (
	"fmt"
	"errors" 

	"github.com/viocha/go-option"
	"github.com/viocha/go-option/result"
)

func main() {
	// Option 示例
	val := option.Val(42)
	if val.IsVal() {
		fmt.Println("Option value:", val.Get()) // Option value: 42
	}

	noneVal := option.Nul[int]()
	fmt.Println("Option value or default:", noneVal.GetOr(0)) // Option value or default: 0

	// Result 示例
	resSuccess := result.Ok("success")
	resSuccess.Try(func(s string) {
		fmt.Println("OK:", s) // OK: success
	}).Catch(func(err error) {
		fmt.Println("Error:", err)
	})

	resFailure := result.Err[string](errors.New("something went wrong"))
	if resFailure.IsErr() {
		fmt.Println("Error found:", resFailure.GetErr()) // Error found: something went wrong
	}
	
	fmt.Println("Result value or default:", resFailure.GetOr("default value")) // Result value or default: default value
}

```

---

## 📘 类型概览

### `Option[T]`

`Option[T]` 表示一个可选值。可以是：

* `Val[T](value)`：表示存在值。
* `Nul[T]()`：表示值不存在。

#### 构造器

* `Val[T](value T) Option[T]`
* `Nul[T]() Option[T]`
* `From[T](val T, err error) Option[T]`
* `FromPtr[T](val *T) Option[T]`
* `FromFunc[T](f func() T) Option[T]`

#### 方法列表

| 方法                         | 返回类型         | 描述                                   |
|----------------------------|--------------|--------------------------------------|
| `String()`                 | `string`     | 返回 Option 的字符串表示                     |
| `IsVal()`                  | `bool`       | 是否包含值                                |
| `IsNul()`                  | `bool`       | 是否为空                                 |
| `Has(value T)`             | `bool`       | 值是否等于指定值                             |
| `HasFunc(f func(T) bool)`  | `bool`       | 值是否满足函数条件                            |
| `Try(f func(T))`           | `Option[T]`  | 如果有值则执行函数                            |
| `Catch(f func())`          | `Option[T]`  | 如果无值则执行函数                            |
| `Finally(f func())`        | `Option[T]`  | 执行函数并返回原 Option (若函数 panic 则返回 None) |
| `Else(f func() Option[T])` | `Option[T]`  | 如果无值则执行函数构造新值                        |
| `ElseVal(f func() T)`      | `Option[T]`  | 如果无值则执行函数构造 Some(value)              |
| `Filter(f func(T) bool)`   | `Option[T]`  | 满足条件则保留，否则返回 None                    |
| `Get()`                    | `T`          | 获取值或 panic                           |
| `GetOr(value T)`           | `T`          | 获取值或默认值                              |
| `GetOrFunc(f func() T)`    | `T`          | 获取值或调用函数返回默认值                        |
| `GetOrZero()`              | `T`          | 获取值或返回零值                             |
| `ToPtr()`                  | `*T`         | 将值转换为指针                              |
| `ToErr(err error)`         | `error`      | 无值返回指定错误，有值返回 `nil`                  |
| `Unwrap(err error)`        | `(T, error)` | 同时返回值和错误                             |

#### 函数列表

| 函数                                                           | 返回类型        | 描述              |
|--------------------------------------------------------------|-------------|-----------------|
| `Then(o Option[T], f func(T) Option[U])`                     | `Option[U]` | 若 o 有值，则使用 f(o) |
| `Map(o Option[T], f func(T) U)`                              | `Option[U]` | 映射值             |
| `MapOr(o Option[T], f func(T) U, v U)`                       | `U`         | 映射或返回默认值        |
| `MapOrFunc(o Option[T], okFn func(T) U, defaultFn func() U)` | `U`         | 映射或调用函数         |

---

### `Result[T]`

`Result[T]` 表示可能成功也可能失败的计算。可以是：

* `Ok[T](value)`：成功，包含值。
* `Err[T](error)`：失败，包含错误。

#### 构造器

* `Ok[T](value T) Result[T]`
* `Err[T](error error) Result[T]`
* `From[T](val T, err error) Result[T]`
* `FromOption[T](o option.Option[T], err error) Result[T]`
* `FromFunc[T](f func() T) Result[T]`

#### 方法列表

| 方法                             | 返回类型                   | 描述                                  |
|--------------------------------|------------------------|-------------------------------------|
| `String()`                     | `string`               | 返回 Result 的字符串表示                    |
| `IsOk()`                       | `bool`                 | 是否成功                                |
| `IsErr()`                      | `bool`                 | 是否失败                                |
| `Has(value T)`                 | `bool`                 | 是否为 Ok 且值相等                         |
| `HasFunc(func(T) bool)`        | `bool`                 | 是否为 Ok 且值满足条件                       |
| `HasErr(error)`                | `bool`                 | 是否为 Err 且错误相等                       |
| `HasErrFunc(func(error) bool)` | `bool`                 | 是否为 Err 且错误满足函数                     |
| `Try(func(T))`                 | `Result[T]`            | 若为 Ok 执行函数                          |
| `Catch(func(error))`           | `Result[T]`            | 若为 Err 执行函数                         |
| `Finally(f func())`            | `Result[T]`            | 执行函数并返回原 Result (若函数 panic 则返回 Err) |
| `Else(func(error) Result[T])`  | `Result[T]`            | 若为 Err 执行函数构造新值                     |
| `ElseMap(func(error) T)`       | `Result[T]`            | 若为 Err 执行函数将错误映射为成功值                |
| `Get()`                        | `T`                    | 获取值或 panic                          |
| `GetOr(v T)`                   | `T`                    | 获取值或返回默认                            |
| `GetOrZero()`                  | `T`                    | 获取值或返回零值                            |
| `GetOrFunc(f func(error) T)`   | `T`                    | 获取值或调用函数                            |
| `GetErr()`                     | `error`                | 获取错误或 panic                         |
| `Unwrap()`                     | `(T, error)`           | 同时获取值和错误                            |
| `ToPtr()`                      | `*T`                   | 将值转换为指针                             |
| `Val()`                        | `option.Option[T]`     | 将 Ok 转为 Some                        |
| `Err()`                        | `option.Option[error]` | 将 Err 转为 Some                       |

#### 函数列表

| 函数                                                            | 返回类型        | 描述                 |
|---------------------------------------------------------------|-------------|--------------------|
| `Then(r Result[T], f func(T) Result[U])`                      | `Result[U]` | 若成功则调用函数           |
| `Map(r Result[T], f func(T) U)`                               | `Result[U]` | 映射成功的值             |
| `MapOr(r Result[T], f func(T) U, v U)`                        | `U`         | 映射或返回默认值           |
| `MapOrFunc(r Result[T], okFn func(T) U, errFn func(error) U)` | `U`         | 成功用 okFn，失败用 errFn |

---

## 📜 License

MIT
