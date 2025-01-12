package gode

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestGode_Call(t *testing.T) {
	gode, err := New("id = function(v) { return v; }")
	if err != nil {
		t.Fatal(err)
	}

	res, err := gode.Call("id", "bar")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, res, "bar")
}

func TestGode_NestedCall(t *testing.T) {
	gode, err := New("a = {}; a.b = {}; a.b.id = function(v) { return v; }")
	if err != nil {
		t.Fatal(err)
	}

	res, err := gode.Call("a.b.id", "bar")
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, res, "bar")
}

func TestGode_CallMissingFunction(t *testing.T) {
	gode, err := New()
	if err != nil {
		t.Fatal(err)
	}

	_, err = gode.Call("missing")
	if err == nil {
		t.Fatal("no err")
	}
}

func TestGode_Exec(t *testing.T) {
	gode, err := New()
	if err != nil {
		t.Fatal(err)
	}

	res, err := gode.exec("1")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, nil, res)

	res, err = gode.exec("return")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, nil, res)

	res, err = gode.exec("return null")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, nil, res)

	res, err = gode.exec("return function() {}")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, nil, res)

	res, err = gode.exec("return 0")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, 0, res)

	res, err = gode.exec("return true")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, true, res)

	res, err = gode.exec("return 'hello'")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "hello", res)

	res, err = gode.exec("return [1, 2]")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, []int{1, 2}, res)

	res, err = gode.exec("return {a:1,b:2}")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, map[string]interface{}{"a": 1, "b": 2}, res)

	res, err = gode.exec("return \"\u3042\"")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "\u3042", res)

	res, err = gode.exec(`return "\u3042"`)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "\u3042", res)

	res, err = gode.exec("return '\\\\'")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "\\", res)
}

func TestGode_Eval(t *testing.T) {

	gode, err := New()
	if err != nil {
		t.Fatal(err)
	}

	res, err := gode.Eval("")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, nil, res)

	res, err = gode.Eval(" ")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, nil, res)

	res, err = gode.Eval("null")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, nil, res)

	res, err = gode.Eval("function(){}")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, nil, res)

	res, err = gode.Eval("0")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, 0, res)

	res, err = gode.Eval("true")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, true, res)

	res, err = gode.Eval("[1,2]")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, []int{1, 2}, res)

	res, err = gode.Eval("[1, function() {}]")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, []interface{}{1, nil}, res)

	res, err = gode.Eval(`"hello"`)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "hello", res)

	res, err = gode.Eval(`'hello'`)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "hello", res)

	res, err = gode.Eval("'red yellow blue'.split(' ')")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, []string{"red", "yellow", "blue"}, res)

	res, err = gode.Eval("{a:1, b:2}")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, map[string]interface{}{"a": 1, "b": 2}, res)

	res, err = gode.Eval("{a:true,b:function (){}}")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, map[string]interface{}{"a": true}, res)

	res, err = gode.Eval("\"\u3042\"")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "\u3042", res)

	res, err = gode.Eval(`"\u3042"`)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "\u3042", res)

	res, err = gode.Eval(`"\\\\"`)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, `\\`, res)
}

func TestGode_New(t *testing.T) {
	gode, err := New(`foo = function() { return "bar"; }`)
	if err != nil {
		t.Fatal(err)
	}

	res, err := gode.exec("return foo()")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "bar", res)

	res, err = gode.Eval("foo()")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "bar", res)

	res, err = gode.Call("foo")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "bar", res)
}

func TestGode_GlobalScope(t *testing.T) {
	gode, err := New()
	if err != nil {
		t.Fatal(err)
	}

	res, err := gode.Eval("this === (function() {return this})()")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, true, res)

	res, err = gode.exec("return this === (function() {return this})()")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, true, res)
}

func TestGode_LargeScripts(t *testing.T) {
	scripts := strings.Repeat("var foo = 'bar';\n", int(math.Pow10(4)))
	code := fmt.Sprintf("function foo() {\n%s\n};\nreturn true", scripts)

	gode, err := New()
	if err != nil {
		t.Fatal(err)
	}

	res, err := gode.exec(code)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, true, res)
}

func TestGode_SyntaxError(t *testing.T) {
	gode, err := New()
	if err != nil {
		t.Fatal(err)
	}

	_, err = gode.exec(")")
	if err == nil {
		t.Fatal("no err")
	}
}

func TestGode_ThrowException(t *testing.T) {
	gode, err := New()
	if err != nil {
		t.Fatal(err)
	}

	_, err = gode.exec("throw 'hello'")
	if err == nil {
		t.Fatal("no err")
	}
}

func TestGode_BrokenSubstitutions(t *testing.T) {
	gode, err := New()
	if err != nil {
		t.Fatal(err)
	}
	s := "#{source}"
	res, err := gode.Eval(fmt.Sprintf(`"%s"`, s))
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, s, res)
}

func assertEqual(t *testing.T, x, y interface{}) {
	if fmt.Sprintf("%+v", x) != fmt.Sprintf("%+v", y) {
		t.Fatal(fmt.Sprintf("x:%+v y:%+v not equal", x, y))
	}
}
