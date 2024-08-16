from sympy import sieve, primepi

def numSemiprimes(x):
    base = list(sieve.primerange(1,int(x**0.5)+1))
    return sum([primepi(int(x/base[k-1])+1) - k + 1 for k in range(1,primepi(int(x**0.5)))])

print(numSemiprimes(10**8))    
