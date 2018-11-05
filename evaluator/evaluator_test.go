package evaluator

import (
	"javascript_interpreter/lexer"
	"javascript_interpreter/object"
	"javascript_interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestArrayExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected []float64
	}{
		{"let x = [2,3,4]", []float64{2, 3, 4}},
		{"[2,3,4]", []float64{2, 3, 4}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated != nil {

		}
		// result, ok := evaluated.(*object.Array)
		// if !ok {

		// }
		// testArrayObject(t, *result, tt.expected)
	}
}

func TestStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected []float64
	}{
		{`let x = "stringdata"; x;`, []float64{2, 3, 4}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated != nil {

		}
		// result, ok := evaluated.(*object.Array)
		// if !ok {

		// }
		// testArrayObject(t, *result, tt.expected)
	}
}

func testArrayObject(t *testing.T, obj object.Array, expected []float64) bool {

	for i, ele := range obj.Elements {
		result, ok := ele.(*object.Number)
		if !ok {
			t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
			return false
		}
		if result.Value != expected[i] {
			t.Errorf("")
		}
	}
	return true

}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNumberObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Number)

	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{" 1 < 1 ", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != true", false},
		{"false != true", true},
		{"(1 < 2) == true", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not boolean. got=%T, (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) {10} ", nil},
		{"if (1) { 10 } ", 10},
		{"if (1 < 2){ 10 }", 10},
		{"if (1 > 2 ){ 10 } else { 20 }  ", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	return obj == NULL
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"return 10;", 10},
		{"return 10; 9", 10},
		{"return 2 * 5; 8", 10},
		{"9; return 10; 4", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}

}

func TestLetStatements(t *testing.T) {

	tests := []struct {
		input    string
		expected float64
	}{
		{"let a = 5; a", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
	}

	for _, tt := range tests {
		testNumberObject(t, testEval(tt.input), tt.expected)
	}

}

func TestFunctionObject(t *testing.T) {

	input := "fn(x){ x + 2 };"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}

}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"let identity = fn(x){ x; } identity(5);", 5},
		{"let identity = fn(x){ return x;};  identity(5);", 5},
		{"let double = fn(x){ x * 2 }; double(5);", 10},
		{"let add = fn(x, y){ x + y }; add(5, 5);", 10},
		{"let add = fn(x, y){ x + y }; add( 5 + 5,  add(5, 5));", 20},
	}

	for _, tt := range tests {
		testNumberObject(t, testEval(tt.input), tt.expected)
	}
}
