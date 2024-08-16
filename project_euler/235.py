def u(r):
    return sum([(900 - 3*k)*(r**(k-1)) for k in range(1,5001)])

t = 1.003
b = 1.002

for i in range(500):
    r = (t+b)/2.0
    if u(r) < -600000000000:
        t = r
    else:
        b = r

print(r)
