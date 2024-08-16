from sympy import fibonacci as fib

num = 3
s = 0
while fib(num) < 4000000:
  s += fib(num)
  num+=3

print(s)
