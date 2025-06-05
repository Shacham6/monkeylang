# Monkey

This is yet another implementation of the Monkey Programming Language authored by the incredible Thorsten Ball.

## Purpose of this project

This project was not made with challenging the status quo in mind, but rather of a personal want to see through a project like this to the end.

**STATUS UPDATE: IT'S DONE!**

**Finally, after a year and two months, my journey with this project has concluded!** Feel free, whoever you are, to use whatever you want from this project that you want if at all seems useful to you.

I could not recommend the two books by Thorsten Ball: "Writing an Interpreter In Go" and "Writing a Compiler In Go" more. It's astounding how he made Compiler Engineering _approachable_, accessible, readable and digestible.

## Example 1

Given The script `myscript.monkey`:

```
let my_map = {
    "one": fn(x) {
        x + 1
    },

    "one1": fn() {
        puts("this is so cool")
    }
}

puts("This is my programming language")

let result_one = my_map["one"]("one")
my_map[result_one]()
```

```
> go run main.go -file myscript.monkey
This is my programming language
this is so cool
```

## Example 2 - The Macro System

We also have macros!

```
let unless = macro(condition, consequence, alternative) {
	quote(if (!(unquote(condition))) {
		unquote(consequence);
	} else {
		unquote(alternative);
	});
}

unless(10 > 5, puts("not greater"), puts("greater"));
```

Running this will result in `greater` to be printed.

## VM and Tree execution

Both execution modes for tree and vm are available for immediate execution. Given `script.monkey`:

```sh
> go run . -engine vm -file script.monkey
# or...
> go run . -engine tree -file script.monkey
```
