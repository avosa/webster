# Documentation

This README provides an overview of the language's features and syntax.

Keep in mind that this arithmetic interpreter is limited in functionality and does not support the full feature set described in the README. Developing a complete interpreter or compiler for the Webster programming language is a work in progress.

### Variables

Variables in Webster are declared using the let keyword, followed by the variable name, an optional type annotation, and an initial value.

```swift
let x: Int = 42
let name: String = "Alice"
```

### Data Types

Webster currently supports the following basic data types:

1. Int: Integer numbers
2. Float: Floating-point numbers
3. Bool: Boolean values (true or false)
4. String: Text strings
5. Arrays: You can create an array by specifying the element type and providing a list of elements Square brackets `[]`.

For example:

```swift
let numbers: Array<Int> = [1, 2, 3, 4, 5]
let names: Array<String> = ["Alice", "Bob", "Charlie"]
```

To access an element in an array, use the index of the element enclosed in square brackets:

```swift
let firstNumber = numbers[0] // 1
let firstName = names[0] // "Alice"
```

6. Dictionaries: Dictionaries, also known as hash maps or associative arrays, are unordered collections of key-value pairs. In Webster, you can create a dictionary by specifying the key type and value type, and providing a list of key-value pairs enclosed in curly braces {}.

For example:

```swift
let ages: Dictionary<String, Int> = ["Alice": 30, "Bob": 25, "Charlie": 22]
```

To access a value in a dictionary, use its key enclosed in square brackets:

```swift
let aliceAge = ages["Alice"] // 30
```

7. Custom Types: Webster can also support custom types using structures (structs) and classes. Structs are value types, while classes are reference types. You can define properties and methods for both structs and classes.

- Structs

```swift
struct Point {
  let x: Float
  let y: Float

  func distance(to point: Point) -> Float {
    let deltaX = x - point.x
    let deltaY = y - point.y
    return sqrt(deltaX * deltaX + deltaY * deltaY)
  }
}

let pointA = Point(x: 0, y: 0)
let pointB = Point(x: 3, y: 4)
let distance = pointA.distance(to: pointB)
```

- Classes

```swift
class Person {
  let name: String
  let age: Int

  init(name: String, age: Int) {
    self.name = name
    self.age = age
  }

  func greet() {
    print("Hello, my name is \(name) and I am \(age) years old.")
  }
}

let alice = Person(name: "Alice", age: 30)
alice.greet()
```

### Conditional Statements

Webster supports `if` and `if-else` statements for conditional execution:

```swift
if x > 0 {
  print("Positive")
} else {
  print("Non-positive")
}
```

### Loops

Webster supports `for` and `while` loops for iteration:

```swift
for i in 0..10 {
  print(i)
}

while x > 0 {
  x -= 1
}
```

### Functions

Functions in Webster are declared using the `func` keyword, followed by the function name, parameter list, optional return type, and function body.

```swift
func add(a: Int, b: Int) -> Int {
  return a + b
}

let sum = add(3, 4)
```

### Input and Output

Reading input from the user and printing output to the console is done using the `readLine()` and `print()` functions:

```swift
let input = readLine()
print("You entered: \(input)")
```

### Classes and Inheritance

Webster supports classes, inheritance, and polymorphism:

```swift
class Animal {
  func speak() {}
}

class Dog: Animal {
  override func speak() {
    print("Woof!")
  }
}

let dog = Dog()
dog.speak()
```

### Error Handling

Webster supports error handling using `throws`, `try`, `catch`, and custom error types:

```swift
enum DivisionError: Error {
  case zeroDenominator
}

func divide(a: Int, b: Int) throws -> Float {
  if b == 0 {
    throw DivisionError.zeroDenominator
  }
  return a / b
}

do {
  let result = try divide(10, 0)
} catch {
  print("Error: \(error)")
}
```

### Modules and Packages

Webster supports modules and packages for organizing and managing code. To define a module, create a new file with the `.wb` extension and use the module keyword:

```swift
// math.wb
module Math

func add(a: Int, b: Int) -> Int {
  return a + b
}

func subtract(a: Int, b: Int) -> Int {
  return a - b
}
```

To use a module in another file, use the import keyword:

```swift
// main.wb
import Math

let sum = Math.add(3, 4)
```

### To do

1. Concurrency
2. Standard Library
3. Package Management
4. Dynamic Libraries
5. Advanced Language Features
6. Tools and Ecosystem
