from sympy import factorint
from functools import reduce
from operator import mul

def numSols(n):
    pf = factorint(n)
    return (reduce(mul, [2*value + 1 for key, value in pf.items()])+1)/2

n = 100
while True:
    if numSols(n) > 1000:
        print(n)
        break
    n += 1
