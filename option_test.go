package go_option

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestSomeAndNone(t *testing.T) {
	s := Some(42)
	if !s.IsSome() || s.IsNone() {
		t.Error("Some should be IsSome() == true and IsNone() == false")
	}
	if v := s.Get(); v != 42 {
		t.Errorf("Expected 42, got %v", v)
	}

	n := None[int]()
	if n.IsSome() || !n.IsNone() {
		t.Error("None should be IsSome() == false and IsNone() == true")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic on Get() from None")
		}
	}()
	n.Get() // should panic
}

func TestGetOr(t *testing.T) {
	opt := None[int]()
	if val := opt.GetOr(10); val != 10 {
		t.Errorf("Expected 10, got %v", val)
	}

	opt2 := Some(5)
	if val := opt2.GetOr(10); val != 5 {
		t.Errorf("Expected 5, got %v", val)
	}
}

func TestHasAndHasFunc(t *testing.T) {
	opt := Some(100)
	if !opt.Has(100) {
		t.Error("Expected Has to return true for 100")
	}
	if opt.Has(200) {
		t.Error("Expected Has to return false for 200")
	}
	if !opt.HasFunc(func(x int) bool { return x > 50 }) {
		t.Error("Expected HasFunc to return true for >50")
	}
}

func TestDoAndElseDo(t *testing.T) {
	called := false
	opt := Some("Go")
	opt.Do(func(val string) {
		if val != "Go" {
			t.Errorf("Expected value 'Go', got %v", val)
		}
		called = true
	})
	if !called {
		t.Error("Expected Do to call function")
	}

	none := None[string]()
	none.Else(func() {
		called = true
	})
	if !called {
		t.Error("Expected ElseDo to call function")
	}
}

func TestFilter(t *testing.T) {
	opt := Some(10)
	result := opt.Filter(func(v int) bool { return v > 5 })
	if result.IsNone() {
		t.Error("Expected Filter to keep value")
	}
	result2 := opt.Filter(func(v int) bool { return v < 5 })
	if result2.IsSome() {
		t.Error("Expected Filter to remove value")
	}
}

func TestOrAndOrFunc(t *testing.T) {
	none := None[string]()
	some := Some("hello")

	if val := none.Or(some); !val.Has("hello") {
		t.Error("Expected Or to return 'hello'")
	}
	if val := none.OrFunc(func() Option[string] { return Some("world") }); !val.Has("world") {
		t.Error("Expected OrFunc to return 'world'")
	}
}

func TestXor(t *testing.T) {
	a := Some(1)
	b := None[int]()
	if res := a.Xor(b); !res.Has(1) {
		t.Error("Expected Xor to return Some(1)")
	}
	if res := b.Xor(a); !res.Has(1) {
		t.Error("Expected Xor (symmetric) to return Some(1)")
	}

	c := None[int]()
	d := Some(2)
	if res := c.Xor(d); !res.Has(2) {
		t.Error("Expected Xor to return Some(2)")
	}

	e := Some(3)
	f := Some(4)
	if res := e.Xor(f); res.IsSome() {
		t.Error("Expected Xor to return None")
	}

	g := None[int]()
	h := None[int]()
	if res := g.Xor(h); res.IsSome() {
		t.Error("Expected Xor(None, None) to return None")
	}
}

func TestOption_ToErr(t *testing.T) {
	errMsg := errors.New("value is missing")
	none := None[int]()
	if err := none.ToErr(errMsg); err == nil || !errors.Is(err, errMsg) {
		t.Error("Expected error from ToErr on None")
	}

	some := Some(5)
	if err := some.ToErr(errMsg); err != nil {
		t.Error("Expected nil from ToErr on Some")
	}
}

func TestOption_GetValErr(t *testing.T) {
	errMsg := errors.New("value is missing for GetValErr")
	none := None[int]()
	val, err := none.ToValErr(errMsg)
	if err == nil || !errors.Is(err, errMsg) || val != 0 {
		t.Error("Expected error and zero value from GetValErr on None")
	}

	some := Some(5)
	valSome, errSome := some.ToValErr(errMsg)
	if errSome != nil || valSome != 5 {
		t.Errorf("Expected value 5 and nil error from GetValErr on Some, got val %v, err %v", valSome, errSome)
	}
}

func TestOption_Map(t *testing.T) {
	opt := Some(2)
	mapped := Map(opt, func(x int) string { return fmt.Sprintf("num:%d", x) })
	if !mapped.Has(fmt.Sprintf("num:%d", 2)) {
		t.Error("Expected mapped Option to be Some(\"num:2\")")
	}

	noneOpt := None[int]()
	mappedNone := Map(noneOpt, func(x int) string { t.Error("Map func called on None"); return "fail" })
	if mappedNone.IsSome() {
		t.Error("Expected Map on None to return None")
	}
}

func TestOption_MapOr(t *testing.T) {
	opt := Some(2)
	res := MapOr(opt, func(x int) string { return "yes" }, "no")
	if res != "yes" {
		t.Errorf("Expected 'yes', got %v", res)
	}

	noneOpt := None[int]()
	resNone := MapOr(noneOpt, func(x int) string { t.Error("MapOr func called on None"); return "fail" }, "no_from_none")
	if resNone != "no_from_none" {
		t.Errorf("Expected MapOr on None to return default value, got %s", resNone)
	}
}

func TestOption_MapOrFunc(t *testing.T) {
	opt := Some(3)
	resOpt := MapOrFunc(opt,
		func(x int) string { return fmt.Sprintf("val-%d", x) },
		func() string { t.Error("MapOrFunc defaultFn called on Some"); return "default" },
	)
	if resOpt != "val-3" {
		t.Errorf("Expected MapOrFunc on Some to use okFn, got %s", resOpt)
	}

	none := None[int]()
	res2 := MapOrFunc(none,
		func(x int) string { t.Error("MapOrFunc okFn called on None"); return "yes" },
		func() string { return "fallback" },
	)
	if res2 != "fallback" {
		t.Errorf("Expected fallback, got %v", res2)
	}
}

func TestAndAndFunc(t *testing.T) {
	a := Some("ok")
	b := Some(123)
	result := And(a, b)
	if !result.Has(123) {
		t.Error("Expected And to return Some(123)")
	}

	none := None[string]()
	result2 := And(none, b)
	if result2.IsSome() {
		t.Error("Expected And with None to return None")
	}

	fnResult := AndFunc(Some(2), func(x int) Option[string] {
		return Some("ok")
	})
	if !fnResult.Has("ok") {
		t.Error("Expected AndFunc to return Some(\"ok\")")
	}

	noneAndFunc := AndFunc(None[int](), func(x int) Option[string] {
		t.Error("AndFunc's func called on None")
		return Some("fail")
	})
	if noneAndFunc.IsSome() {
		t.Error("Expected AndFunc on None to return None")
	}

	someAndFuncToNone := AndFunc(Some(3), func(x int) Option[string] {
		return None[string]()
	})
	if someAndFuncToNone.IsSome() {
		t.Error("Expected AndFunc with func returning None to return None")
	}
}

func TestOption_String(t *testing.T) {
	sSome := Some(42).String()
	if !strings.HasPrefix(sSome, "Some[int]") || !strings.Contains(sSome, "42") {
		t.Errorf("Expected Some string representation, got %s", sSome)
	}

	sNone := None[string]().String()
	if !strings.HasPrefix(sNone, "None[string]") {
		t.Errorf("Expected None string representation, got %s", sNone)
	}
}

func TestOption_FromVal(t *testing.T) {
	opt1 := From(10, nil)
	if !opt1.IsSome() || opt1.Get() != 10 {
		t.Errorf("Expected FromVal with nil error to be Some(10), got %v", opt1)
	}

	err := errors.New("fromval error")
	opt2 := From(0, err)
	if opt2.IsSome() {
		t.Errorf("Expected FromVal with error to be None, got %v", opt2)
	}
}

func TestOption_GetOrFunc(t *testing.T) {
	someVal := Some(55)
	valSome := someVal.GetOrFunc(func() int {
		t.Error("GetOrFunc's func called on Some")
		return 0
	})
	if valSome != 55 {
		t.Errorf("Expected GetOrFunc on Some to return value, got %d", valSome)
	}

	noneVal := None[int]()
	valNone := noneVal.GetOrFunc(func() int { return 99 })
	if valNone != 99 {
		t.Errorf("Expected GetOrFunc on None to return func result, got %d", valNone)
	}
}

func TestOption_GetOrZero(t *testing.T) {
	someInt := Some(123)
	if someInt.GetOrZero() != 123 {
		t.Errorf("Expected GetOrZero on Some(int) to return value, got %d", someInt.GetOrZero())
	}

	noneInt := None[int]()
	if noneInt.GetOrZero() != 0 {
		t.Errorf("Expected GetOrZero on None[int] to return 0, got %d", noneInt.GetOrZero())
	}

	someStr := Some("hello")
	if someStr.GetOrZero() != "hello" {
		t.Errorf("Expected GetOrZero on Some(string) to return value, got %s", someStr.GetOrZero())
	}

	noneStr := None[string]()
	if noneStr.GetOrZero() != "" {
		t.Errorf("Expected GetOrZero on None[string] to return \"\", got %s", noneStr.GetOrZero())
	}

	type myStruct struct{ val int }
	someStruct := Some(myStruct{val: 1})
	if someStruct.GetOrZero().val != 1 {
		t.Errorf("Expected GetOrZero on Some(struct) to return value, got %v", someStruct.GetOrZero())
	}
	noneStruct := None[myStruct]()
	if noneStruct.GetOrZero() != (myStruct{}) {
		t.Errorf("Expected GetOrZero on None[struct] to return zero struct, got %v", noneStruct.GetOrZero())
	}
}
