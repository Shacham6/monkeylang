package evaluator_test

import (
	. "monkey/evaluator/internal/evaluatortest"
	"monkey/object"
	"monkey/testutils"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5", 10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := DoEval(tt.input)
		CheckIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"1 == 1", true},
		{"2 > 1", true},
		{"2 < 1", false},
		{"1 != 1", false},
		{"1 >= 1", true},
		{"1 <= 1", true},
		{"1 != 2", true},
		{"(1 + 2) == 3", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := DoEval(tt.input)
			CheckBooleanObject(t, evaluated, tt.expected)
		})
	}
}

func TestEvalNullExpression(t *testing.T) {
	evaluated := DoEval("null")
	CheckNullObject(t, evaluated)
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
		evaluated := DoEval(tt.input)
		CheckBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := DoEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				CheckIntegerObject(t, evaluated, int64(integer))
			} else {
				CheckNullObject(t, evaluated)
			}
		})
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 1;", 1},
		{"return 1; 2", 1},
		{"return 2 * 3; 1", 6},
		{"1; return 2; 3", 2},
		{`
			if (1 > 0) {
				if (2 > 1) {
					return 10;
				}
				return 1;
			}
		`, 10},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := DoEval(tt.input)
			CheckIntegerObject(t, evaluated, tt.expected)
		})
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
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 2) {
					return true + false;
				}
			}
			`, "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`"hello" - "world"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := DoEval(tt.input)

			errObj := testutils.CheckIsA[object.Error](t, evaluated, "evaluated is not an error object.")

			if errObj.Message != tt.expectedMessage {
				t.Errorf("wrong error message. expected = %q, got = %q", tt.expectedMessage, errObj.Message)
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = 5; b + a;", 10},
		{"let a = 5; let b = a; b;", 5},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			CheckIntegerObject(
				t,
				DoEval(tt.input),
				tt.expected,
			)
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) {x + 2;};`
	evaluated := DoEval(input)
	fn := testutils.CheckIsA[object.Function](t, evaluated, "evaluated is not a 'object.Function'")

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters = %+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got = %q", fn.Parameters[0])
	}

	expectedBody := "(block (expr (infix x + 2)))"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got = %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) {return x;}; identity(5);", 5},
		{"let double = fn(x) {x * 2;}; double(5);", 10},
		{"let add = fn(x, y) {x + y;}; add(5, 5);", 10},
		{"let add = fn(x, y) {x + y;}; add(5 + 5, add(5, 5));", 20},
		{"fn(x){x;}(5);", 5},
	}

	for _, tt := range tests {
		CheckIntegerObject(t, DoEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"praise the sun"`

	evaluated := DoEval(input)
	str := testutils.CheckIsA[object.String](t, evaluated, "evaluated is not object.String")

	if str.Value != "praise the sun" {
		t.Errorf("str.Value is not %q. got = %q", "praise the sun", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected CheckEvaluated
	}{
		{`"hello" + " " + "world"`, NewResultInString("hello world")}, // TODO(Jajo): Add str + int
		{`"number " + 1`, NewResultInString("number 1")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := DoEval(tt.input)
			tt.expected.CheckEvaluated(t, evaluated)
		})
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected CheckEvaluated
	}{
		// Tests for `len`
		{`len("")`, NewResultInInt(0)},
		{`len("four")`, NewResultInInt(4)},
		{`len(1)`, NewResultInError("argument to `len` not supported, got INTEGER")},
		{`len("one", "two")`, NewResultInError("wrong number of arguments. got = 2, want = 1")},
		{`len([1, 2, 3])`, NewResultInInt(3)},

		// Tests for `first`
		{`first([1, 2])`, NewResultInInt(1)},
		{`first([2, 1])`, NewResultInInt(2)},
		{`first([])`, NewResultInNil()},
		{`first([], [])`, NewResultInError("wrong number of arguments. got = 2, want = 1")},
		{`first([], [], [])`, NewResultInError("wrong number of arguments. got = 3, want = 1")},
		{`first(123)`, NewResultInError("argument to `first` must be an ARRAY, got INTEGER")},

		// Tests for `last`
		{`last([1, 2])`, NewResultInInt(2)},
		{`last([2, 1])`, NewResultInInt(1)},
		{`last([])`, NewResultInNil()},
		{`last([], [])`, NewResultInError("wrong number of arguments. got = 2, want = 1")},
		{`last([], [], [])`, NewResultInError("wrong number of arguments. got = 3, want = 1")},
		{`last(123)`, NewResultInError("argument to `last` must be an ARRAY, got INTEGER")},

		// Tests for `rest`
		{`rest([1, 2, 3])`, NewResultInArray(NewResultInInt(2), NewResultInInt(3))},
		{`rest([3, 2, 1])`, NewResultInArray(NewResultInInt(2), NewResultInInt(1))},
		{`rest([])`, NewResultInNil()},
		{`rest(123)`, NewResultInError("argument to `rest` must be an ARRAY, got INTEGER")},

		// Tests for `push`
		{`push([], 1)`, NewResultInArray(NewResultInInt(1))},
		{`push([1, 2], 3)`, NewResultInArray(
			NewResultInInt(1), NewResultInInt(2), NewResultInInt(3),
		)},
		{`let arr = []; push(arr, 1)`, NewResultInArray(NewResultInInt(1))},
		{`push([])`, NewResultInError("wrong number of arguments. got = 1, want = 2")},
		{`push([], 12, 12)`, NewResultInError("wrong number of arguments. got = 3, want = 2")},
		{`push(2, 2)`, NewResultInError("first argument to `push` must be ARRAY, got INTEGER")},
		// Tests for `sprintf`
		{`sprintf("1")`, NewResultInString("1")},
		{`sprintf("123")`, NewResultInString("123")},
		{`sprintf("%s", "hello")`, NewResultInString("hello")},
		{`sprintf("%s %s", "hello", "world")`, NewResultInString("hello world")},
		{`sprintf("hello %s", "world")`, NewResultInString("hello world")},
		{`sprintf("1 %d", 2)`, NewResultInString("1 2")},
		{`sprintf()`, NewResultInError("sprintf function requires at least a single argument")},
		{`sprintf(2)`, NewResultInError("first argument to `sprintf` must be STRING, got = INTEGER")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := DoEval(tt.input)
			if !tt.expected.CheckEvaluated(t, evaluated) {
				return
			}
		})
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 1 + 1, 1 + 1 + 1]"
	evaluated := DoEval(input)
	result := testutils.CheckIsA[object.Array](t, evaluated, "evaluated is not a object.Array")
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong amount of elements. got = %d", len(result.Elements))
	}

	CheckIntegerObject(t, result.Elements[0], 1)
	CheckIntegerObject(t, result.Elements[1], 2)
	CheckIntegerObject(t, result.Elements[2], 3)
}

func TestIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected CheckEvaluated
	}{
		{
			"[1, 2, 3][0]",
			NewResultInInt(1),
		},
		{
			"[1, 2, 3][1]",
			NewResultInInt(2),
		},
		{
			"[1, 2, 3][2]",
			NewResultInInt(3),
		},
		{
			"let i = 0; [0][i]",
			NewResultInInt(0),
		},
		{
			"[0, 1, 2][1 + 1]",
			NewResultInInt(2),
		},
		{
			"[1, 2, 3][3]",
			NewResultInNil(),
		},
		{
			"[1, 2, 3][-1]",
			NewResultInNil(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := DoEval(tt.input)
			tt.expected.CheckEvaluated(t, evaluated)
		})
	}
}

func TestHashLiteral(t *testing.T) {
	input := `
	let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := DoEval(input)
	result := testutils.CheckIsA[object.Hash](t, evaluated, "evaluated is not object.Hash")

	expected := map[object.HashKey]int64{
		mustHashKey(t, &object.String{Value: "one"}):   1,
		mustHashKey(t, &object.String{Value: "two"}):   2,
		mustHashKey(t, &object.String{Value: "three"}): 3,
		mustHashKey(t, &object.Integer{Value: 4}):      4,
		mustHashKey(t, &object.Boolean{Value: true}):   5,
		mustHashKey(t, &object.Boolean{Value: false}):  6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has the wrong num of pairs. got = %d, want = %d", len(result.Pairs), len(expected))
	}

	for expectedKey, expectedValue := range expected {
		pairValue, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		CheckIntegerObject(t, pairValue.Value, expectedValue)
	}
}

func mustHashKey(t *testing.T, o object.Object) object.HashKey {
	hashKey, err := o.HashKey()
	if err != nil {
		t.Fatalf("failed to calculate hash key for type %T", o)
	}
	return hashKey
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected CheckEvaluated
	}{
		{`{"one": 1}["one"]`, NewResultInInt(1)},
		{`{"tw" + "o": 2}["two"]`, NewResultInInt(2)},
		{`{"one": 1}["on" + "e"]`, NewResultInInt(1)},
		{`{}["nothing"]`, NewResultInNil()},
		{`let d = {"one": 1}; d["one"]`, NewResultInInt(1)},
		{`let k = "one"; {"one": 1}[k]`, NewResultInInt(1)},
		{`{[1, 2, 3]: 1}`, NewResultInError("unusable as hash key: ARRAY")},
		{`{"one": "one"}[fn(){}]`, NewResultInError("unusable as hash key: FUNCTION")},
		{`{1: 1}[1]`, NewResultInInt(1)},
		{`{true: 1}[true]`, NewResultInInt(1)},
		{`{false: 1}[false]`, NewResultInInt(1)},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := DoEval(tt.input)
			tt.expected.CheckEvaluated(t, evaluated)
		})
	}
}
