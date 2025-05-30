package option

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestSomeAndNone(t *testing.T) {
	s := Val(42)
	if !s.IsVal() || s.IsNul() {
		t.Error("Some should be IsSome() == true and IsNone() == false")
	}
	if v := s.Get(); v != 42 {
		t.Errorf("Expected 42, got %v", v)
	}

	n := Nul[int]()
	if n.IsVal() || !n.IsNul() {
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
	opt := Nul[int]()
	if val := opt.GetOr(10); val != 10 {
		t.Errorf("Expected 10, got %v", val)
	}

	opt2 := Val(5)
	if val := opt2.GetOr(10); val != 5 {
		t.Errorf("Expected 5, got %v", val)
	}
}

func TestHasAndHasFunc(t *testing.T) {
	opt := Val(100)
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
	opt := Val("Go")
	opt.Try(func(val string) {
		if val != "Go" {
			t.Errorf("Expected value 'Go', got %v", val)
		}
		called = true
	})
	if !called {
		t.Error("Expected Do to call function")
	}

	none := Nul[string]()
	none.Catch(func() {
		called = true
	})
	if !called {
		t.Error("Expected ElseDo to call function")
	}
}

func TestFilter(t *testing.T) {
	opt := Val(10)
	result := opt.Filter(func(v int) bool { return v > 5 })
	if result.IsNul() {
		t.Error("Expected Filter to keep value")
	}
	result2 := opt.Filter(func(v int) bool { return v < 5 })
	if result2.IsVal() {
		t.Error("Expected Filter to remove value")
	}
}

func TestOption_ToErr(t *testing.T) {
	errMsg := errors.New("value is missing")
	none := Nul[int]()
	if err := none.ToErr(errMsg); err == nil || !errors.Is(err, errMsg) {
		t.Error("Expected error from ToErr on None")
	}

	some := Val(5)
	if err := some.ToErr(errMsg); err != nil {
		t.Error("Expected nil from ToErr on Some")
	}
}

func TestOption_GetValErr(t *testing.T) {
	errMsg := errors.New("value is missing for GetValErr")
	none := Nul[int]()
	val, err := none.ToValErr(errMsg)
	if err == nil || !errors.Is(err, errMsg) || val != 0 {
		t.Error("Expected error and zero value from GetValErr on None")
	}

	some := Val(5)
	valSome, errSome := some.ToValErr(errMsg)
	if errSome != nil || valSome != 5 {
		t.Errorf("Expected value 5 and nil error from GetValErr on Some, got val %v, err %v", valSome, errSome)
	}
}

func TestOption_Map(t *testing.T) {
	opt := Val(2)
	mapped := Map(opt, func(x int) string { return fmt.Sprintf("num:%d", x) })
	if !mapped.Has(fmt.Sprintf("num:%d", 2)) {
		t.Error("Expected mapped Option to be Some(\"num:2\")")
	}

	noneOpt := Nul[int]()
	mappedNone := Map(noneOpt, func(x int) string { t.Error("Map func called on None"); return "fail" })
	if mappedNone.IsVal() {
		t.Error("Expected Map on None to return None")
	}
}

func TestOption_MapOr(t *testing.T) {
	opt := Val(2)
	res := MapOr(opt, func(x int) string { return "yes" }, "no")
	if res != "yes" {
		t.Errorf("Expected 'yes', got %v", res)
	}

	noneOpt := Nul[int]()
	resNone := MapOr(noneOpt, func(x int) string { t.Error("MapOr func called on None"); return "fail" }, "no_from_none")
	if resNone != "no_from_none" {
		t.Errorf("Expected MapOr on None to return default value, got %s", resNone)
	}
}

func TestOption_MapOrFunc(t *testing.T) {
	opt := Val(3)
	resOpt := MapOrFunc(opt,
		func(x int) string { return fmt.Sprintf("val-%d", x) },
		func() string { t.Error("MapOrFunc defaultFn called on Some"); return "default" },
	)
	if resOpt != "val-3" {
		t.Errorf("Expected MapOrFunc on Some to use okFn, got %s", resOpt)
	}

	none := Nul[int]()
	res2 := MapOrFunc(none,
		func(x int) string { t.Error("MapOrFunc okFn called on None"); return "yes" },
		func() string { return "fallback" },
	)
	if res2 != "fallback" {
		t.Errorf("Expected fallback, got %v", res2)
	}
}

func TestOption_String(t *testing.T) {
	sSome := Val(42).String()
	if !strings.HasPrefix(sSome, "Some[int]") || !strings.Contains(sSome, "42") {
		t.Errorf("Expected Some string representation, got %s", sSome)
	}

	sNone := Nul[string]().String()
	if !strings.HasPrefix(sNone, "None[string]") {
		t.Errorf("Expected None string representation, got %s", sNone)
	}
}

func TestOption_FromVal(t *testing.T) {
	opt1 := From(10, nil)
	if !opt1.IsVal() || opt1.Get() != 10 {
		t.Errorf("Expected FromVal with nil error to be Some(10), got %v", opt1)
	}

	err := errors.New("fromval error")
	opt2 := From(0, err)
	if opt2.IsVal() {
		t.Errorf("Expected FromVal with error to be None, got %v", opt2)
	}
}

func TestOption_GetOrFunc(t *testing.T) {
	someVal := Val(55)
	valSome := someVal.GetOrFunc(func() int {
		t.Error("GetOrFunc's func called on Some")
		return 0
	})
	if valSome != 55 {
		t.Errorf("Expected GetOrFunc on Some to return value, got %d", valSome)
	}

	noneVal := Nul[int]()
	valNone := noneVal.GetOrFunc(func() int { return 99 })
	if valNone != 99 {
		t.Errorf("Expected GetOrFunc on None to return func result, got %d", valNone)
	}
}

func TestOption_GetOrZero(t *testing.T) {
	someInt := Val(123)
	if someInt.GetOrZero() != 123 {
		t.Errorf("Expected GetOrZero on Some(int) to return value, got %d", someInt.GetOrZero())
	}

	noneInt := Nul[int]()
	if noneInt.GetOrZero() != 0 {
		t.Errorf("Expected GetOrZero on None[int] to return 0, got %d", noneInt.GetOrZero())
	}

	someStr := Val("hello")
	if someStr.GetOrZero() != "hello" {
		t.Errorf("Expected GetOrZero on Some(string) to return value, got %s", someStr.GetOrZero())
	}

	noneStr := Nul[string]()
	if noneStr.GetOrZero() != "" {
		t.Errorf("Expected GetOrZero on None[string] to return \"\", got %s", noneStr.GetOrZero())
	}

	type myStruct struct{ val int }
	someStruct := Val(myStruct{val: 1})
	if someStruct.GetOrZero().val != 1 {
		t.Errorf("Expected GetOrZero on Some(struct) to return value, got %v", someStruct.GetOrZero())
	}
	noneStruct := Nul[myStruct]()
	if noneStruct.GetOrZero() != (myStruct{}) {
		t.Errorf("Expected GetOrZero on None[struct] to return zero struct, got %v", noneStruct.GetOrZero())
	}
}
