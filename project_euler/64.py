from math import floor

def contfracsqrt(n):
    m = 0
    d = 1
    a = floor(n**0.5)
    h = []
    r = []
    while (m,d,a) not in h:
        h += [(m,d,a)]
        r += [a]
        m = d*a - m
        d = (n - m**2)/d
        a = floor(((n**0.5)+m)/d)
    return r

c = 0
for n in range(1,10000):
    if floor(n**0.5) == n**0.5:
        pass
    else:
        if len(contfracsqrt(n)) % 2 == 0:
            c += 1

print(c)
