from math import log as ln

file = open("base_exp.txt")

z = [line.split(',') for line in file]

file.close()

for a in z:
    a[0] = int(a[0])
    a[1] = int(a[1])

m = 0
l = 1
for a in z:
    if a[1] * ln(a[0]) > m:
        m = a[1] * ln(a[0])
        print(l)
    l += 1

