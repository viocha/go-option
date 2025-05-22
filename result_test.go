package option

import (
	"errors"
	"fmt"
	"testing"
)

func TestOkAndErr(t *testing.T) {
	r1 := Ok(123)
	if !r1.IsOk() || r1.IsErr() {
		t.Errorf("Expected Ok, got Err")
	}
	if r1.Get() != 123 {
		t.Errorf("Expected value 123, got %v", r1.Get())
	}

	errExample := errors.New("test error")
	r2 := Err[int](errExample)
	if !r2.IsErr() || r2.IsOk() {
		t.Errorf("Expected Err, got Ok")
	}
	if !r2.HasErr(errExample) {
		t.Errorf("Expected error to match")
	}
}

func TestHas(t *testing.T) {
	r := Ok("hello")
	if !r.Has("hello") {
		t.Errorf("Expected Has to return true for matching value")
	}
	if r.Has("world") {
		t.Errorf("Expected Has to return false for non-matching value")
	}
}

func TestHasFunc(t *testing.T) {
	r := Ok(42)
	if !r.HasFunc(func(v int) bool { return v > 40 }) {
		t.Errorf("Expected predicate to match")
	}
	if r.HasFunc(func(v int) bool { return v < 0 }) {
		t.Errorf("Expected predicate to fail")
	}
}

func TestHasErrFunc(t *testing.T) {
	err := errors.New("fail")
	r := Err[int](err)
	if !r.HasErrFunc(func(e error) bool { return e.Error() == "fail" }) {
		t.Errorf("Expected predicate to match error")
	}
}

func TestDoAndElseDo_Result(t *testing.T) {
	okCalled := false
	errCalled := false

	Ok("ok").Do(func(v string) { okCalled = true }).ElseDo(func(e error) { errCalled = true })
	if !okCalled || errCalled {
		t.Errorf("Expected Do to be called, ElseDo not to be called")
	}

	okCalled, errCalled = false, false
	Err[string](errors.New("error")).Do(func(v string) { okCalled = true }).ElseDo(func(e error) { errCalled = true })
	if okCalled || !errCalled {
		t.Errorf("Expected ElseDo to be called, Do not to be called")
	}
}

func TestOr(t *testing.T) {
	r1 := Err[int](errors.New("fail"))
	r2 := Ok(99)
	result := r1.Or(r2)
	if !result.IsOk() || result.Get() != 99 {
		t.Errorf("Expected fallback to second Result with value 99")
	}
}

func TestOrFunc(t *testing.T) {
	err := errors.New("fallback")
	r := Err[int](errors.New("original")).OrFunc(func(e error) Result[int] {
		if e.Error() == "original" {
			return Ok(1)
		}
		return Err[int](err)
	})
	if !r.IsOk() || r.Get() != 1 {
		t.Errorf("Expected Ok(1) as fallback from OrFunc")
	}
}

func TestMapErr(t *testing.T) {
	r := Err[int](errors.New("origin"))
	r2 := r.MapErr(func(e error) error {
		return fmt.Errorf("wrapped: %w", e)
	})
	if !r2.HasErrFunc(func(e error) bool { return e.Error() == "wrapped: origin" }) {
		t.Errorf("Expected wrapped error message")
	}
}

func TestGetWithErr(t *testing.T) {
	ok := Ok("abc")
	v, err := ok.GetWithErr()
	if err != nil || v != "abc" {
		t.Errorf("Expected value 'abc' with nil error")
	}

	r := Err[string](errors.New("fail"))
	_, err = r.GetWithErr()
	if err == nil {
		t.Errorf("Expected error from GetWithErr")
	}
}

func TestRMapAndOr(t *testing.T) {
	r := Ok(2)
	mapped := RMap(r, func(x int) string { return fmt.Sprintf("%d!", x) })
	if !mapped.IsOk() || mapped.Get() != "2!" {
		t.Errorf("Expected mapped value to be '2!'")
	}

	err := errors.New("fail")
	errR := Err[int](err)
	defaultVal := "default"
	res := RMapOr(errR, func(x int) string { return "should not run" }, defaultVal)
	if res != defaultVal {
		t.Errorf("Expected fallback value '%s'", defaultVal)
	}
}

func TestPanicOnGet(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on Get() from Err result")
		}
	}()
	Err[string](errors.New("fail")).Get()
}

func TestPanicOnGetErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on GetErr() from Ok result")
		}
	}()
	Ok("ok").GetErr()
}
