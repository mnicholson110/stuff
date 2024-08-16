from sympy import factorint, factorial, perfect_power, isprime, primerange
from math import log, floor

# see: [1] 'http://mathworld.wolfram.com/SmarandacheFunction.html' for details
# computes quickly: sum([mu(n) for n in range(2,1001)])
# computes quickly: sum([mu(n) for n in range(2,10001)])
# computes quickly: sum([mu(n) for n in range(2,10**5 + 1)])
# computes slowly: sum([mu(n) for n in range(2,10**6 + 1)])
# does not compute in reasonable time: sum([mu(n) for n in range(2,10**7 + 1)])
# the answer: sum([mu(n) for n in range(2,10**8 + 1)])

# previously, I was using sympy.smoothness to check whether n | gpf(n)!, as it is known that the set of n in N s.t. n | gpf(n)! has density 1
# this method might need to be used again (perhaps with something better than sympy.smoothness to get the gpf of n

# pre-memoize small values of p**a
# a is capped at the minimum of p and the largest a s.t. p**a < 10**8
# note: for any prime p > 10**4: p**2 > 10**8 (maybe I should pre-memoize all primes below 10**8? is that better than just checking primality?)

memo = {p**a: a*p for p in primerange(1,10**4) for a in range(1,min(p,int(log(10**8,p)))+1)}


#this takes a sec
for p in primerange(10**4,10**7):
    memo[p] = p

# pre-memoize n! < 10**8

for n in range(1,13):
    memo[factorial(n)] = n

def mu(n):

    # check if the value is memoized
    check = memo.get(n)
    if check is not None:
        return memo[n]

    # sympy.perfect_power returns (base, exp) or False
    tmp = perfect_power(n)
    if tmp:

        # if n is a prime power not pre-memoized, then exp > base, and we use kempner's algorithm from [1]
        if isprime(tmp[0]):
            return (tmp[0] - 1)*tmp[1] + sumk(tmp[0],tmp[1])

        # sympy.factorint returns a dictionary of { primefactor : power }
        # we recursively check the value of mu(primefactor**power), returning the max
        # this should be very fast most of the time, as small primes have all been pre-memoized
        else:
            f = factorint(n)
            memo[n] = max([mu(p**f[p]) for p in f.keys()])
            return memo[n]
    else:

        # any prime not pre-memoized : mu(p) = p
        if isprime(n):
            memo[n] = n
            return memo[n]
        
        # else recursively check the prime powers in n's factorization as above
        f = factorint(n)
        memo[n] = max([mu(p**f[p]) for p in f.keys()])
        return memo[n]
    

def sumk(p,alpha):
    if alpha <= p:
        return alpha
    else:
        v = floor(log(1 + alpha*(p-1),p))
        a = (p**v - 1)/(p-1)
        k, r = divmod(alpha, a)
        s = k
        while r != 0:
            v -= 1
            a = (p**v - 1)/(p-1)
            k, r = divmod(r, a)
            s += k
        return s
