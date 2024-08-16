from sympy import isprime, divisors, sieve

sieve.extend(100000000)

def cond(n):
    if not isprime(int(2 + n/2)): return False
    else:
        d = divisors(n)
        for i in range(2,int(len(d)/2)):
            if not isprime(d[i] + int(n/d[i])): return False
    return True

U = [i-1 for i in sieve.primerange(1,100000000) if (i-1) % 4 != 0]

print(sum([n for n in U if cond(n)]))
