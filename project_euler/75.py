from usefulStuff import pythagTriples
from collections import Counter

L = []
R = []
for i in range(2,65000,2):
    for tri in pythagTriples(i):
            L.append(sum(tri))
print("L done")

for item, count in Counter(L).most_common():
    if count == 1 and item <= 1500000:
        R.append(item)

print(len(R))
