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
)

val := option.Some(42)
if val.IsSome() {
  fmt.Println(val.Get()) // 42
}

res := option.Ok("success")
res.Do(func (s string) {
	fmt.Println("OK:", s)
}).ElseDo(func (err error) {
	fmt.Println("Error:", err)
})
```

---

## 📘 类型概览

### `Option[T]`

`Option[T]` 表示一个可选值。可以是：

* `Some[T](value)`：表示存在值。
* `None[T]()`：表示值不存在。

#### 构造器

* `Some[T](value T) Option[T]`
* `None[T]() Option[T]`

#### 方法列表

| 方法                           | 返回类型         | 描述                |
|------------------------------|--------------|-------------------|
| `IsSome()`                   | `bool`       | 是否包含值             |
| `IsNone()`                   | `bool`       | 是否为空              |
| `Has(value T)`               | `bool`       | 值是否等于指定值          |
| `HasFunc(f func(T) bool)`    | `bool`       | 值是否满足函数条件         |
| `Do(f func(T))`              | `Option[T]`  | 如果有值则执行函数         |
| `ElseDo(f func())`           | `void`       | 如果无值则执行函数         |
| `Filter(f func(T) bool)`     | `Option[T]`  | 满足条件则保留，否则返回 None |
| `Or(b Option[T])`            | `Option[T]`  | 若无值则返回备选          |
| `OrFunc(f func() Option[T])` | `Option[T]`  | 若无值则调用函数并返回其结果    |
| `Xor(b Option[T])`           | `Option[T]`  | 仅当某一方有值时返回该值      |
| `Get()`                      | `T`          | 获取值或 panic        |
| `GetOr(value T)`             | `T`          | 获取值或默认值           |
| `GetOrFunc(f func() T)`      | `T`          | 获取值或调用函数返回默认值     |
| `GetOrZero()`                | `T`          | 获取值或返回零值          |
| `ToErr(err error)`           | `error`      | 无值返回错误            |
| `GetWithErr(err error)`      | `(T, error)` | 同时返回值和错误          |
| `ToResult(err error)`        | `Result[T]`  | 转为 `Result` 类型    |

#### 函数列表

| 函数                                                           | 返回类型        | 描述               |
|--------------------------------------------------------------|-------------|------------------|
| `And(a Option[T], b Option[U])`                              | `Option[U]` | 若 a 和 b 均有值，返回 b |
| `AndFunc(o Option[T], f func(T) Option[U])`                  | `Option[U]` | 若 o 有值，则使用 f(o)  |
| `Map(o Option[T], f func(T) U)`                              | `Option[U]` | 映射值              |
| `MapOr(o Option[T], f func(T) U, v U)`                       | `U`         | 映射或返回默认值         |
| `MapOrFunc(o Option[T], okFn func(T) U, defaultFn func() U)` | `U`         | 映射或调用函数          |

---

### `Result[T]`

`Result[T]` 表示可能成功也可能失败的计算。可以是：

* `Ok[T](value)`：成功，包含值。
* `Err[T](error)`：失败，包含错误。

#### 构造器

* `Ok[T](value T) Result[T]`
* `Err[T](error error) Result[T]`

#### 方法列表

| 方法                              | 返回类型            | 描述              |
|---------------------------------|-----------------|-----------------|
| `IsOk()`                        | `bool`          | 是否成功            |
| `IsErr()`                       | `bool`          | 是否失败            |
| `Has(value T)`                  | `bool`          | 是否为 Ok 且值相等     |
| `HasFunc(func(T) bool)`         | `bool`          | 是否为 Ok 且值满足条件   |
| `HasErr(error)`                 | `bool`          | 是否为 Err 且错误相等   |
| `HasErrFunc(func(error) bool)`  | `bool`          | 是否为 Err 且错误满足函数 |
| `Do(func(T))`                   | `Result[T]`     | 若为 Ok 执行函数      |
| `ElseDo(func(error))`           | `Result[T]`     | 若为 Err 执行函数     |
| `Or(Result[T])`                 | `Result[T]`     | 若为 Err 返回备选     |
| `OrFunc(func(error) Result[T])` | `Result[T]`     | 若为 Err 执行函数并返回  |
| `MapErr(func(error) error)`     | `Result[T]`     | 映射错误            |
| `Get()`                         | `T`             | 获取值或 panic      |
| `GetOr(v T)`                    | `T`             | 获取值或返回默认        |
| `GetOrZero()`                   | `T`             | 获取值或返回零值        |
| `GetOrFunc(f func(error) T)`    | `T`             | 获取值或调用函数        |
| `GetErr()`                      | `error`         | 获取错误或 panic     |
| `GetWithErr()`                  | `(T, error)`    | 同时获取值和错误        |
| `Val()`                         | `Option[T]`     | 将 Ok 转为 Some    |
| `Err()`                         | `Option[error]` | 将 Err 转为 Some   |

#### 函数列表

| 函数                                                             | 返回类型        | 描述                 |
|----------------------------------------------------------------|-------------|--------------------|
| `RAnd(a Result[T], b Result[U])`                               | `Result[U]` | 若 a 成功，返回 b        |
| `RAndFunc(r Result[T], f func(T) Result[U])`                   | `Result[U]` | 若成功则调用函数           |
| `RMap(r Result[T], f func(T) U)`                               | `Result[U]` | 映射成功的值             |
| `RMapOr(r Result[T], f func(T) U, v U)`                        | `U`         | 映射或返回默认值           |
| `RMapOrFunc(r Result[T], okFn func(T) U, errFn func(error) U)` | `U`         | 成功用 okFn，失败用 errFn |

---

## 📜 License

MIT

