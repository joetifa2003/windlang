<p align="center">
  <img src="assets/WindLogo.png" />
</p>

# WindLang, A simple programming language built with golang 🍃

[![Go Report Card](https://goreportcard.com/badge/github.com/joetifa2003/windlang)](https://goreportcard.com/report/github.com/joetifa2003/windlang)

- [WindLang, A simple programming language built with golang 🍃](#windlang-a-simple-programming-language-built-with-golang-)
  - [What is wind?](#what-is-wind)
  - [Playground](#playground)
  - [Cool but why?](#cool-but-why)
  - [Why Golang, not Rust?](#why-golang-not-rust)
  - [How to use it?](#how-to-use-it)
  - [So what can it do?](#so-what-can-it-do)
    - [Hello world?](#hello-world)
    - [Variables](#variables)
    - [Data types](#data-types)
    - [Arrays](#arrays)
      - [Array.push(element) -> any[]](#arraypushelement---any)
      - [Array.pop() -> any](#arraypop---any)
      - [Array.len() -> int](#arraylen---int)
      - [Array.join(separator) -> string](#arrayjoinseparator---string)
      - [Array.map(function) -> any[]](#arraymapfunction---any)
      - [Array.filter(function) -> any[]](#arrayfilterfunction---any)
      - [Array.reduce(fn(accumulator, element), initialValue) -> any](#arrayreducefnaccumulator-element-initialvalue---any)
    - [Strings](#strings)
      - [String.len(separator) -> int](#stringlenseparator---int)
      - [String.charAt(index) -> string](#stringcharatindex---string)
      - [String.contains(substr) -> string](#stringcontainssubstr---string)
      - [String.containsAny(substr) -> string](#stringcontainsanysubstr---string)
      - [String.count(substr) -> string](#stringcountsubstr---string)
      - [String.replace(old, new) -> string](#stringreplaceold-new---string)
      - [String.replaceAll(old, new) -> string](#stringreplaceallold-new---string)
      - [String.replaceN(old, new, n) -> string](#stringreplacenold-new-n---string)
      - [String.changeAt(index, new) -> string](#stringchangeatindex-new---string)
      - [String.indexOf(substr) -> string](#stringindexofsubstr---string)
      - [String.lastIndexOf(substr) -> string](#stringlastindexofsubstr---string)
      - [String.split(separator) -> string[]](#stringsplitseparator---string)
      - [String.trim() -> string](#stringtrim---string)
      - [String.toLowerCase() -> string](#stringtolowercase---string)
      - [String.toUpperCase() -> string](#stringtouppercase---string)
    - [Functions](#functions)
    - [Closures](#closures)
    - [If expressions](#if-expressions)
    - [Include statement](#include-statement)
    - [For loops](#for-loops)
    - [While loops](#while-loops)
    - [HashMaps](#hashmaps)
  - [Todos](#todos)

## What is wind?

Wind is an interpreted language written in golang.

## Playground

The easiest way to play with WindLang is to use the [playground](https://windlang.vercel.app/playground)

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
Install the vscode extension [here](https://marketplace.visualstudio.com/items?itemName=YoussefAhmed.windlang)!

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
{
    "name": "Youssef",
    "age": 18
}              // HashMaps
```

### Arrays

```swift
let arr = [1, "2", true];

println(arr[0]); // 1
println(arr[1]); // 2
println(arr[2]); // true

arr.push("3") // [1, "2", true, "3"]
println(arr[3]); // 3
arr.pop()     // [1, "2", true]
```

Arrays in Wind can take any type of data, and can be of any size
You can append to an array by using the `append` function
You can remove an element by index using the `remove` function.

#### Array.push(element) -> any[]

```swift
let x = [1, 2, 3];
x.push(4) // [1, 2, 3, 4]
```

Array push function adds an element to the end of the array

#### Array.pop() -> any

```swift
let x = [1, 2, 3, 4];
x.push(4); // [1, 2, 3, 4]

let y = x.pop(); // [1, 2, 3]
println(y); // 4
```

Array pop function removes the last element of the array and returns it

#### Array.len() -> int

```swift
let x = [1, "hi", true];
println(x.len()); // 3
```

Array len function returns the length of the array

#### Array.join(separator) -> string

```swift
let x = [1, "hi", 3, 4];

println(x.join("-")); // 1-hi-3-4
```

Array join function returns a string with all the elements of the array separated by the separator

#### Array.map(function) -> any[]

```swift
let x = [1, 2, 3];

println(x.map(fn(x) {x * 2})); // [2, 4, 6]
```

Array map function applies the function to each element of the array and returns a new array with the results

#### Array.filter(function) -> any[]

```swift
let x = [1, 2, 3, 4];
let even = x.filter(fn(x) { x % 2 == 0});

println(even); // [2, 4]]
```

Array filter function applies the function to each element of the array and if the function returns true, the element is added to the new array

#### Array.reduce(fn(accumulator, element), initialValue) -> any

```swift
let x = [1, 2, 3, 4, 5];
let sum = x.reduce(fn(acc, x) { acc + x}, 0);

println(sum); // 15
```

Array reduce function applies the function to each element of the array and returns a single value

### Strings

Strings in Wind start and end with a double quote `"` and can contain any character and can be multi-line

#### String.len(separator) -> int

```swift
let x = "Hello";

println(x.len()); // 5
```

String len function returns the length of the string

#### String.charAt(index) -> string

```swift
let name = "youssef";
println(name.charAt(0)); // y
```

String charAt function returns the character at the specified index

#### String.contains(substr) -> string

```swift
let name = "youssef";

println(name.contains("ss")); // true
```

String contains function returns true if the string contains the exact substring

#### String.containsAny(substr) -> string

```swift
let vowels = "aeiou";
let name = "youssef";

println(name.contains(vowels)); // true
```

String contains function returns true if the string contains any character of the substring

#### String.count(substr) -> string

```swift
let name = "youssef";

println(name.count("s")); // 2
```

String count function returns the number of times the substring appears in the string

#### String.replace(old, new) -> string

```swift
let name = "John Doe";

println(name.replace("o", "x")); // Jxhn Doe
```

String replace function returns a new string after replacing one old substring with a new substring

#### String.replaceAll(old, new) -> string

```swift
let name = "John Doe";

println(name.replaceAll("o", "x")); // Jxhn Dxe
```

String replace function returns a new string after replacing all old substring with a new substring

#### String.replaceN(old, new, n) -> string

```swift
let name = "Youssef";

println(name.replaceN("s", "x", 1)); // Youxsef
println(name.replaceN("s", "x", 2)); // Youxxef
```

String replace function returns a new string after replacing n of old substring with a new substring

#### String.changeAt(index, new) -> string

```swift
let name = "Ahmed";

println(name.changeAt(0, "a")); // ahmed
```

String changeAt function returns a new string after changing the character at the specified index

#### String.indexOf(substr) -> string

```swift
let name = "John Doe";

println(name.indexOf("o")); // 1
```

String indexOf function returns the index of the first occurrence of the substring

#### String.lastIndexOf(substr) -> string

```swift
let name = "John Doe";

println(name.lastIndexOf("o")); // 6
```

String indexOf function returns the index of the last occurrence of the substring

#### String.split(separator) -> string[]

```swift
let name = "Youssef Ahmed";
let names = name.split(" "); // ["Youssef", "Ahmed"]
let firstName = names[0];
let lastName = names[1];

println("First name is: " + firstName); // First name is: Youssef
println("Last name is: " + lastName); // Last name is: Ahmed
```

String join function returns an array of strings by splitting the string by the separator

#### String.trim() -> string

```swift
let name = " John Doe   ";

println(name.trim()); // "John Doe"
```

String trim function removes whitespace from the start/end of the string and returns a new string

#### String.toLowerCase() -> string

```swift
let name = "JoHn dOe";

println(name.toLowerCase()); // john doe
```

String toLowerCase function returns a new string with all the characters in lower case

#### String.toUpperCase() -> string

```swift
let name = "JoHn dOe";

println(name.toUpperCase()); // JOHN DOE
```

String toLowerCase function returns a new string with all the characters in upper case

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
let addConstructor = fn(x) {
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


let addTwo = addConstructor(2); // This will return a function

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

### Include statement

```swift
// test.wind

let msg = "Hello 👋";

let greeter = fn() {
    println(msg);
};
```

```swift
// main.wind
include "./test.wind";

greeter(); // Hello 👋

// You can also alias includes
include "./test.wind" as test;

test.greeter(); // Hello 👋
```

Include statements allow you to include other Wind scripts, It initializes them once and can be used by multiple files at the same time while preserving state.

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

### HashMaps

```swift
let person = {
    "name": "Youssef",
    "age": 18,
    "incrementAge": fn() {
        person.age++;
    }
};

println(person["name"]); // Youssef
println(person.age); // 18
person.incrementAge();
println(person.age); // 19
```

Hashmaps are like js object and can store key value pairs, Keys can be integers, strings and booleans. Values can be any type.

## Todos

-   ~~Named include statements~~

-   ~~HashMaps (Javascript objects~~

-   A bytecode interpreter maybe
