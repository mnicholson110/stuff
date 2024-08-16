p_set = range(1,5)

p_space = [a+b+c+d+e+f+g+h+i for a in p_set for b in p_set for c in p_set for d in p_set for e in p_set for f in p_set for g in p_set for h in p_set for i in p_set]

c_set = range(1,7)

c_space = [a+b+c+d+e+f for a in c_set for b in c_set for c in c_set for d in c_set for e in c_set for f in c_set]

p_den = {}

for res in range(9, 9*4 + 1):
    p_den[res] = p_space.count(res)/4**9

c_den = {}

for res in range(6, 6*6 + 1):
    c_den[res] = c_space.count(res)/6**6

prob = 0

for y in p_den.keys():
    for x in c_den.keys():
        if x < y:
            prob += p_den[y]*c_den[x]

print(round(prob,7))
