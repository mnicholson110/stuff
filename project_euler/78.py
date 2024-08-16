from sympy import npartitions as p

n = 5
while p(n) % 10**6 != 0:
    n += 1

print(n)

