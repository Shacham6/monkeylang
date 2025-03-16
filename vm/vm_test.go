package vm_test

import (
	"testing"

	"monkey/object"
	"monkey/vm/internal/vmtest"
)

func objInteger(val int) *object.Integer {
	return &object.Integer{Value: int64(val)}
}

func mustHash(t *testing.T, obj object.Object) object.HashKey {
	hash, err := obj.HashKey()
	if err != nil {
		t.Fatalf("obj is not hashable: %q", err)
	}
	return hash
}

func TestIntegerArithmetic(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("1", 1),
		vmtest.New("2", 2),
		vmtest.New("1 + 2", 3),
		vmtest.New("1 - 2", -1),
		vmtest.New("1 * 2", 2),
		vmtest.New("4 / 2", 2),
		vmtest.New("((1 + 2 - 1) * 2) / 2", 2),
		vmtest.New("4 / 2 * 2 + 2 - 2", 4),
		vmtest.New("-5", -5),
		vmtest.New("-10", -10),
		vmtest.New("-50 + 100 + -50", 0),
		vmtest.New("(5 + 10 * 2 + 15 / 3) * 2 + -10", 50),
	})
}

func TestBooleanExpressions(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("true", true),
		vmtest.New("false", false),
		vmtest.New("1 < 2", true),
		vmtest.New("1 > 2", false),
		vmtest.New("1 < 1", false),
		vmtest.New("1 > 1", false),
		vmtest.New("1 == 1", true),
		vmtest.New("1 != 1", false),
		vmtest.New("1 == 2", false),
		vmtest.New("1 != 2", true),
		vmtest.New("true == true", true),
		vmtest.New("false == true", false),
		vmtest.New("false == false", true),
		vmtest.New("true == false", false),
		vmtest.New("true != false", true),
		vmtest.New("false != true", true),
		vmtest.New("(1 > 2) == false", true),
		vmtest.New("(1 < 2) == false", false),
		vmtest.New("(1 > 2) == true", false),
		vmtest.New("(1 < 2) == true", true),
		vmtest.New("!true", false),
		vmtest.New("!false", true),
		vmtest.New("!5", false),
		vmtest.New("!!true", true),
		vmtest.New("!!false", false),
		vmtest.New("!!5", true),
		vmtest.New("!null", true),
		vmtest.New("!!null", false),
	})
}

func TestNil(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("null", nil),
	})
}

func TestConditionals(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("if (true) { 10 }", 10),
		vmtest.New("if (true) { 10 } else { 20 }", 10),
		vmtest.New("if (false) { 10 } else { 20 }", 20),
		vmtest.New("if (true) { 5 + 5 } else { 20 }", 10),
		vmtest.New("if (false) { 10 } else { 10 + 10 }", 20),
		vmtest.New("if (1) {10}", 10),
		vmtest.New("if (1 - 1) {10} else {20}", 20),
		vmtest.New("if (1 < 2) {10} else {20}", 10),
		vmtest.New("if (1 > 2) {10} else {20}", 20),
		vmtest.New("if (true) {10}; 20", 20),
		vmtest.New("if (false) {false;}", nil),
		vmtest.New("if (null) {10} else {20}", 20),
		vmtest.New("if (null == null) {10}", 10),
	})
}

func TestGlobalLetStatements(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("let a = 1; a", 1),
		vmtest.New("let a = 1; let b = 2; a + b", 3),
		vmtest.New("let a = 1; let b = a + a; a + b", 3),
	})
}

func TestStringExpressions(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New(`"lol"`, "lol"),
		vmtest.New(`"mon" + "key"`, "monkey"),
		vmtest.New(`"mon" + "key" + "banana"`, "monkeybanana"),
	})
}

func TestArrayExpressions(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("[]", []int{}),
		vmtest.New("[1, 2, 3]", []int{1, 2, 3}),
		vmtest.New("[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}),
	})
}

func TestHashLiterals(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New(
			"{}",
			map[object.HashKey]int64{},
		),
		vmtest.New(
			"{1: 2, 2: 3}",
			map[object.HashKey]int64{
				mustHash(t, objInteger(1)): 2,
				mustHash(t, objInteger(2)): 3,
			},
		),
		vmtest.New(
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[object.HashKey]int64{
				mustHash(t, objInteger(2)): 4,
				mustHash(t, objInteger(6)): 16,
			},
		),
	})
}

func TestIndexExpressions(t *testing.T) {
	// I've removed some cases I considered problematic.
	// Problematic in the "I disagree with the language design pretty heavily here"
	// way and so I don't wanna do it. Specifically regarding out-of-bounds tolerance.
	//
	// Here are the cases I've removed:
	//     vmtest.New("[][0]", nil),
	//     vmtest.New("[1, 2, 3, 4][888]", nil),
	//     vmtest.New(`{}["hello"]`, nil),

	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("[0][0]", 0),
		vmtest.New("[0, 1][0 + 1]", 1),
		vmtest.New("[0, 1, 2, 3, 4][1 + 2]", 3),
		vmtest.New("[[1, 2]][0]", []int{1, 2}),
		vmtest.New(`{"hello": "world"}["hello"]`, "world"),
		vmtest.New(`{"hel" + "lo": "world"}["hello"]`, "world"),
		vmtest.New(`{"hello": "wor" + "ld"}["hello"]`, "world"),
	})
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New(
			`
			let fivePlusTen = fn() {5 + 10;};
			fivePlusTen();
			`,
			15,
		),
		vmtest.New(
			`
			let one = fn() {1;};
			let two = fn() {2;};
			one() + two()
			`,
			3,
		),
		vmtest.New(
			`
			let a = fn() { 1 };
			let b = fn() { a() + 1 };
			let c = fn() { b() + 1 };
			c()
			`,
			3,
		),
		vmtest.New(
			`
			let get_five = fn() { return 5; }
			get_five()
			`,
			5,
		),
		vmtest.New(
			`
			let get = fn() {
				if (true) { return 1; }
				return 2;
			}
			get();
			`,
			1,
		),
		vmtest.New(
			`
			let f = fn() { }
			f()
			`,
			nil,
		),
		vmtest.New(
			`
			let noReturn1 = fn() {};
			let noReturn2 = fn() { noReturn1(); };
			noReturn2();
			`,
			nil,
		),
		vmtest.New(
			`
			let retOne = fn() { 1; };
			let retRetOne = fn() { retOne; };
			retRetOne()();
			`,
			1,
		),
	})
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New(
			`
			let one = fn() { let one = 1; one; };
			one();
			`,
			1,
		),
		vmtest.New(
			`
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two};
			oneAndTwo()
			`,
			3,
		),
		vmtest.New(
			`
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two };
			let threeAndFour = fn() { let three = 3; let four = 4; three + four };
			oneAndTwo() + threeAndFour()
			`,
			10,
		),
		vmtest.New(
			`
			let firstFoobar = fn() { let foobar = 50; foobar }
			let secondFoobar = fn() { let foobar = 100; foobar }
			firstFoobar() + secondFoobar()
			`,
			150,
		),
		vmtest.New(
			`
			let globalSeed = 50;
			let minusOne = fn() {
				let num = 1;
				globalSeed - num
			}
			let minusTwo = fn() {
				let num = 2;
				globalSeed - num
			}
			minusOne() + minusTwo()
			`,
			97,
		),
	})
}

func TestFirstClassFunctions(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New(
			`
			let returnsOneReturner = fn() {
				let returnsOne = fn() { 1; }
				returnsOne;
			}
			returnsOneReturner()();
			`,
			1,
		),
	})
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New(
			`
			let identity = fn(a) { a; }
			identity(4);
			`,
			4,
		),
		vmtest.New(
			`
			let sum = fn(a, b) {
				a + b;
			};
			sum(1, 2);
			`,
			3,
		),
	})
}

// func TestGoal(t *testing.T) {
// 	vmtest.RunVmTests(t, []vmtest.VmTestCase{
// 		vmtest.New(
// 			`
// 			let globalNum = 10;
//
// 			let sum = fn(a, b) {
// 				let c = a + b;
// 				c + globalNum;
// 			}
//
// 			let outer = fn() {
// 				sum(1, 2) + sum(3, 4) + globalNum;
// 			}
//
// 			outer() + globalNum;
// 			`,
// 			40,
// 		),
// 	})
// }
