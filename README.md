# ether

[![CircleCI](https://circleci.com/gh/muiscript/ether.svg?style=svg)](https://circleci.com/gh/muiscript/ether)

A toy programming language implemented in go.

## language features

- ether has `integer`, `boolean`, `array`, and `function` as literals
- One of the most (or maybe, only) notable feature of ether is arrow operator `->`. It works like [Elixir's pipe operator](https://elixir-lang.org/getting-started/enumerables-and-streams.html#the-pipe-operator), which makes successive data transformations readable

## sample code

```ruby
# function
var double = |x| { x * 2 }
puts(double(5)) # 10


# if expression
puts(if (5 > 3) { 5 }) # 5


# combination of function and if expression
var max = |x, y| { if (x > y) { x } else { y } }
puts(max(-4, 5)) # 5


# function of ether is closure
var gen_adder = |x| { |y| { x + y } }
var add_three = gen_adder(3)
puts(add_three(4)) # 7


# builtin function: len
puts(len([3, 2, 7])) # 3


# builtin function: map
var squares = map([1, 2, 3], |x| { x * x }) 
puts(squares) # [1, 4, 9]


# builtin function: filter
var odd_numbers = filter([1, 2, 3], |x| { x % 2 != 0 }) 
puts(odd_numbers) # [1, 3]


# builtin function: reduce
var sum = map([1, 2, 3, 4], 0, |acc, x| { acc + x }) 
puts(sum) # 10

var product = map([1, 2, 3, 4], 1, |acc, x| { acc * x }) 
puts(sum) # 24


# you can pass user defined function as a parameter of builtin functions
var triple = |x| { x * 3 }
var triples = map([1, 2, 3], triple)
puts(triples) # [3, 6, 9]


# arrow function
# x -> f() is equivalent to f(x). x -> f(y) is equivalent to f(x, y)
var sum_of_squares_of_odds_between_ten_and_fifty =
[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
-> filter(|x| { x % 2 != 0 })      # [1, 3, 5, 7, 9]
-> map(|x| { x * x })              # [1, 9, 25, 49, 81]
-> filter(|sq| { sq > 10 })        # [25, 49, 81]
-> filter(|sq| { sq < 50 })        # [25, 49]
-> reduce(0, |acc, x| { acc + x }) # 74

puts(sum_of_squares_of_odds_between_ten_and_fifty) # 74
```
