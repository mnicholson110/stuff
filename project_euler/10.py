from sympy import sieve

print(sum([i for i in sieve.primerange(2,2000000)]))
