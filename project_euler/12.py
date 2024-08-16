from sympy import divisors

def triangle(n):
  return sum(range(n+1))

i = 1
while len(divisors(triangle(i))) < 501:
  i += 1
print(triangle(i))
