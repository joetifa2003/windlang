<p align="center">
  <img src="assets/WindLogo.png" />
</p>

# WindLang, A simple programming language built with golang 🍃

- [WindLang, A simple programming language built with golang 🍃](#windlang-a-simple-programming-language-built-with-golang-)
  - [What is wind?](#what-is-wind)
  - [Cool but why?](#cool-but-why)
  - [Why Golang, not Rust?](#why-golang-not-rust)
  - [How to use it?](#how-to-use-it)
  - [So what can it do?](#so-what-can-it-do)
    - [Hello world?](#hello-world)
    - [Variables](#variables)
    - [Data types](#data-types)
    - [Functions](#functions)
    - [Closures](#closures)
    - [If expressions](#if-expressions)
    - [For loops](#for-loops)
    - [While loops](#while-loops)
  - [Todos](#todos)

## What is wind?

Wind is an interpreted language written in golang.

## Cool but why?

I'm working on this as a side project and it's the 5th attempt at writing my own programming language, And it got me my first job too! it's awesome 💙

## Why Golang, not Rust?

I already tried to build this with Rust, but I battled a lot with the Rust compiler and felt like golang is a better middle ground I got so much faster with golang, and it was surprising that the golang implementation was actually faster!
Not because golang is faster than rust, it's my lack of knowledge, and the way I code is not "Rusty"

## How to use it?

Head up to the [releases](https://github.com/joetifa2003/windlang/releases) page and download Wind for your operating system

And execute windlang executable

```
WindLang, A simple programming language written in golang

Usage:
  windlang [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Run a Wind script

Flags:
  -h, --help     help for windlang
  -t, --toggle   Help message for toggle

Use "windlang [command] --help" for more information about a command.
```

This is the Wind cli you can use the run command to run a Wind script file

## So what can it do?

Let me demonstrate the features implemented as of now

### Hello world?

```swift
println("Hello world");
```

Yes, Wind can print Hello world, Surprising huh?

### Variables

```swift
let x = 1;
println(x); // 1


x = "Hello world";
println(x); // Hello World
```

You can declare variables in Wind with the `let` keyword, Variables in wind are dynamically typed and u can reassign them to any type.

### Data types

```swift
1              // int
3.14           // float
true, false    // boolean
nil            // null
"Hello World"  // String
[1, "2", true] // Arrays
```

### Functions

```swift
let add = fn(x, y) { x + y; };


// This is the same as
let add = fn(x, y) {
    return x + y;
};

println(add(1, 2)); // 3
```

Yes, this looks like Rust functions. The last expression in a function is implicitly returned

### Closures

```swift
let addConsturctor = fn(x) {
    fn(y) {
        x + y;
    };
};


// This is the same as
let addConstructor = fn(x) {
    return fn(y) {
        return x + y;
    };
};


let addTwo = addConsturctor(2); // This will return a function

println(addTwo(3)); // 5
```

This is one of the coolest things implemented in Wind. As I said functions in Wind are expressions, so you can return a function or pass a function to another function.

```swift
let welcome = fn(name, greeter) {
    greeter() + " " + name
};

let greeter = fn() { "Hello 👋"};

println(welcome("Wind 🍃", greeter)); // Hello 👋 Wind 🍃
```

### If expressions

```swift
// Use if as an expression
let grade = 85;
let msg = if (grade > 50) { "You passed" } else { "You didn't pass" };

println(msg); // You passed

// Or use if as a statement
if (grade > 50) {
    println("You passed");
} else {
    println("You didn't pass");
}
```

As you can see here we can use it as an expression or as a statement

Note that after any expression the semicolon is optional. We can type `"You passed";` or `"You passed"` and it will implicitly return it

### For loops

```swift
let names = ["Youssef", "Ahmed", "Soren", "Morten", "Mads", "Jakob"];

for (let i = 0; i < len(names); i++) {
    println("Hello " + names[i]);
}

// Hello Youssef
// Hello Ahmed
// Hello Soren
// Hello Morten
// Hello Mads
// Hello Jakob
```

### While loops

```swift
let x = 0;

while (x < 5) {
    println(x);
    x++;
}

// 0
// 1
// 2
// 3
// 4
```

## Todos

-   Named include statements

-   HashMaps (Javascript objects)

-   A bytecode interpreter maybe
