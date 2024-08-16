from sympy import totient
from sympy.ntheory.generate import primorial


def R(n):
    return totient(n)/(n-1)

test = [primorial(n) for n in range(1,10)]

m = 10**25
for p1 in test:
    for p2 in test:
        for p3 in test:
            if R(p1*p2*p3) < 15499/94744:
                m = min(m,p1*p2*p3)

print(m)
                
        
