package result

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	
	"github.com/viocha/go-option/option"
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
	if r.HasErrFunc(func(e error) bool { return e.Error() == "other" }) {
		t.Errorf("Expected predicate not to match error")
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
	
	ok := Ok(1)
	okMapped := ok.MapErr(func(e error) error { return errors.New("should not run") })
	if !okMapped.IsOk() || okMapped.Get() != 1 {
		t.Errorf("Expected MapErr on Ok to return original Ok")
	}
}

func TestGetValErr(t *testing.T) {
	ok := Ok("abc")
	v, err := ok.ToValErr()
	if err != nil || v != "abc" {
		t.Errorf("Expected value 'abc' with nil error")
	}
	
	r := Err[string](errors.New("fail"))
	_, err = r.ToValErr()
	if err == nil {
		t.Errorf("Expected error from GetWithErr")
	}
}

func TestMap_Result(t *testing.T) {
	r := Ok(2)
	mapped := Map(r, func(x int) string { return fmt.Sprintf("%d!", x) })
	if !mapped.IsOk() || mapped.Get() != "2!" {
		t.Errorf("Expected mapped value to be '2!'")
	}
	
	errR := Err[int](errors.New("fail map"))
	mappedErr := Map(errR, func(x int) string { return "should not run" })
	if mappedErr.IsOk() {
		t.Errorf("Expected mapped error to remain an error")
	}
	if !mappedErr.HasErr(errR.GetErr()) {
		t.Errorf("Expected mapped error to carry original error")
	}
}

func TestMapOr_Result(t *testing.T) {
	rOk := Ok(3)
	resOk := MapOr(rOk, func(x int) string { return fmt.Sprintf("val: %d", x) }, "default")
	if resOk != "val: 3" {
		t.Errorf("Expected MapOr on Ok to return mapped value, got %s", resOk)
	}
	
	errVal := errors.New("fail")
	errR := Err[int](errVal)
	defaultVal := "default"
	res := MapOr(errR, func(x int) string { return "should not run" }, defaultVal)
	if res != defaultVal {
		t.Errorf("Expected fallback value '%s', got %s", defaultVal, res)
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

func TestFromVal_Result(t *testing.T) {
	r1 := From(10, nil)
	if !r1.IsOk() || r1.Get() != 10 {
		t.Errorf("Expected Ok(10), got %v", r1)
	}
	
	err := errors.New("fromval error")
	r2 := From(0, err)
	if !r2.IsErr() || !r2.HasErr(err) {
		t.Errorf("Expected Err(fromval error), got %v", r2)
	}
}

func TestFromOpt_Result(t *testing.T) {
	optSome := option.Some(5)
	errConv := errors.New("conversion error")
	
	r1 := FromOption(optSome, errConv)
	if !r1.IsOk() || r1.Get() != 5 {
		t.Errorf("Expected Ok(5) from Some, got %v", r1)
	}
	
	optNone := option.None[int]()
	r2 := FromOption(optNone, errConv)
	if !r2.IsErr() || !r2.HasErr(errConv) {
		t.Errorf("Expected Err(conversion error) from None, got %v", r2)
	}
}

func TestString_Result(t *testing.T) {
	okStr := Ok(123).String()
	if !strings.HasPrefix(okStr, "Ok[int]") || !strings.Contains(okStr, "123") {
		t.Errorf("Expected Ok string representation, got %s", okStr)
	}
	
	err := errors.New("test error for string")
	errStr := Err[string](err).String()
	if !strings.HasPrefix(errStr, "Err[string]") || !strings.Contains(errStr, "test error for string") {
		t.Errorf("Expected Err string representation, got %s", errStr)
	}
}

func TestGetOr_Result(t *testing.T) {
	rOk := Ok(10)
	if rOk.GetOr(0) != 10 {
		t.Errorf("Expected GetOr on Ok to return value, got %d", rOk.GetOr(0))
	}
	
	rErr := Err[int](errors.New("err"))
	if rErr.GetOr(0) != 0 {
		t.Errorf("Expected GetOr on Err to return default, got %d", rErr.GetOr(0))
	}
}

func TestGetOrZero_Result(t *testing.T) {
	rOk := Ok(10)
	if rOk.GetOrZero() != 10 {
		t.Errorf("Expected GetOrZero on Ok to return value, got %d", rOk.GetOrZero())
	}
	
	rErrStr := Err[string](errors.New("err"))
	if rErrStr.GetOrZero() != "" {
		t.Errorf("Expected GetOrZero on Err to return zero value for string, got '%s'", rErrStr.GetOrZero())
	}
	
	rErrInt := Err[int](errors.New("err"))
	if rErrInt.GetOrZero() != 0 {
		t.Errorf("Expected GetOrZero on Err to return zero value for int, got %d", rErrInt.GetOrZero())
	}
}

func TestGetOrFunc_Result(t *testing.T) {
	rOk := Ok(10)
	valOk := rOk.GetOrFunc(func(e error) int { t.Error("GetOrFunc's func called on Ok"); return 0 })
	if valOk != 10 {
		t.Errorf("Expected GetOrFunc on Ok to return value, got %d", valOk)
	}
	
	errVal := errors.New("err for getorfunc")
	rErr := Err[int](errVal)
	valErr := rErr.GetOrFunc(func(e error) int {
		if !errors.Is(e, errVal) {
			t.Errorf("Error in GetOrFunc's func mismatch: got %v, want %v", e, errVal)
		}
		return 5
	})
	if valErr != 5 {
		t.Errorf("Expected GetOrFunc on Err to return func result, got %d", valErr)
	}
}

func TestVal_Result(t *testing.T) {
	optSome := Ok(100).Val()
	if !optSome.IsSome() || optSome.Get() != 100 {
		t.Errorf("Expected Val on Ok to return Some(100), got %v", optSome)
	}
	
	optNone := Err[int](errors.New("err")).Val()
	if optNone.IsSome() {
		t.Errorf("Expected Val on Err to return None, got %v", optNone)
	}
}

func TestErr_Result_Method(t *testing.T) { // Renamed to avoid conflict with constructor
	optNone := Ok(100).Err()
	if optNone.IsSome() {
		t.Errorf("Expected Err method on Ok to return None, got %v", optNone)
	}
	
	errVal := errors.New("err for Err method")
	optSomeErr := Err[int](errVal).Err()
	if !optSomeErr.IsSome() || !errors.Is(optSomeErr.Get(), errVal) {
		t.Errorf("Expected Err method on Err to return Some(error), got %v", optSomeErr)
	}
}

func TestAnd_Result(t *testing.T) {
	ok1 := Ok(1)
	ok2 := Ok("hello")
	err1 := Err[int](errors.New("err1"))
	err2 := Err[string](errors.New("err2"))
	
	res1 := And(ok1, ok2)
	if !res1.IsOk() || res1.Get() != "hello" {
		t.Errorf("Expected And(Ok, Ok) to be Ok(value from second), got %v", res1)
	}
	
	res2 := And(err1, ok2)
	if !res2.IsErr() || !res2.HasErr(err1.GetErr()) {
		t.Errorf("Expected And(Err, Ok) to be Err(from first), got %v", res2)
	}
	
	res3 := And(ok1, err2)
	if !res3.IsErr() || !res3.HasErr(err2.GetErr()) {
		t.Errorf("Expected And(Ok, Err) to be Err(from second), got %v", res3)
	}
	
	res4 := And(err1, err2) // Though err2 is Err, err1 is returned
	if !res4.IsErr() || !res4.HasErr(err1.GetErr()) {
		t.Errorf("Expected And(Err, Err) to be Err(from first), got %v", res4)
	}
}

func TestAndFunc_Result(t *testing.T) {
	okVal := Ok(5)
	errVal := Err[int](errors.New("andfunc err"))
	
	res1 := AndFunc(okVal, func(v int) Result[string] {
		return Ok(fmt.Sprintf("val:%d", v))
	})
	if !res1.IsOk() || res1.Get() != "val:5" {
		t.Errorf("Expected AndFunc on Ok to execute func and return Ok, got %v", res1)
	}
	
	res2 := AndFunc(errVal, func(v int) Result[string] {
		t.Error("AndFunc's func called on Err")
		return Ok("should not happen")
	})
	if !res2.IsErr() || !res2.HasErr(errVal.GetErr()) {
		t.Errorf("Expected AndFunc on Err to return original Err, got %v", res2)
	}
	
	expectedErrFromFunc := errors.New("err from func")
	res3 := AndFunc(okVal, func(v int) Result[string] {
		return Err[string](expectedErrFromFunc)
	})
	if !res3.IsErr() || !res3.HasErr(expectedErrFromFunc) {
		t.Errorf("Expected AndFunc on Ok with func returning Err to be Err, got %v", res3)
	}
}

func TestMapOrFunc_Result(t *testing.T) {
	okRes := Ok(10)
	valOk := MapOrFunc(okRes,
		func(v int) string { return fmt.Sprintf("ok-%d", v) },
		func(e error) string { t.Error("MapOrFunc's errFn called on Ok"); return "err" },
	)
	if valOk != "ok-10" {
		t.Errorf("Expected MapOrFunc on Ok to use okFn, got %s", valOk)
	}
	
	errContent := errors.New("maporfunc error")
	errRes := Err[int](errContent)
	valErr := MapOrFunc(errRes,
		func(v int) string { t.Error("MapOrFunc's okFn called on Err"); return "ok" },
		func(e error) string {
			if !errors.Is(e, errContent) {
				t.Errorf("Error in MapOrFunc's errFn mismatch: got %v, want %v", e, errContent)
			}
			return fmt.Sprintf("err-%s", e.Error())
		},
	)
	expectedErrStr := fmt.Sprintf("err-%s", errContent.Error())
	if valErr != expectedErrStr {
		t.Errorf("Expected MapOrFunc on Err to use errFn, got %s, want %s", valErr, expectedErrStr)
	}
}
