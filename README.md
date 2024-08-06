# Monkey

This is yet another implementation of the Monkey Programming Language authored by the incredible Thorsten Ball.

## Purpose of this project

This project was not made with challenging the status quo in mind, but rather of a personal want to see through a project like this to the end.

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
