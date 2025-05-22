package option

import (
	"errors"
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
	none.ElseDo(func() {
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
}

func TestToErrAndGetWithErr(t *testing.T) {
	errMsg := errors.New("value is missing")
	none := None[int]()
	if err := none.ToErr(errMsg); err == nil {
		t.Error("Expected error from ToErr")
	}

	val, err := none.GetWithErr(errMsg)
	if err == nil || val != 0 {
		t.Error("Expected error from GetWithErr")
	}

	some := Some(5)
	if err := some.ToErr(errMsg); err != nil {
		t.Error("Expected nil from ToErr on Some")
	}
}

func TestMapFunctions(t *testing.T) {
	opt := Some(2)
	mapped := Map(opt, func(x int) string { return "num" })
	if !mapped.Has("num") {
		t.Error("Expected mapped Option to be Some(\"num\")")
	}

	res := MapOr(opt, func(x int) string { return "yes" }, "no")
	if res != "yes" {
		t.Errorf("Expected 'yes', got %v", res)
	}

	none := None[int]()
	res2 := MapOrFunc(none, func(x int) string { return "yes" }, func() string { return "fallback" })
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
}
