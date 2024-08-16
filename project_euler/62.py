def getDigits(n):
    p = []
    while n:
        p += [n % 10]
        n //= 10
    p.sort()
    return tuple(p)

cubes = [[getDigits(n**3), n**3] for n in range(345,10000)]

d = {}
for c in cubes:
     d.setdefault(c[0], []).append(c[1])

for k in d.keys():
    if len(d[k]) == 5:
        print(min(d[k]))
        break
