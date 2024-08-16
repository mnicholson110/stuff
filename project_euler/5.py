from sympy import factorint
from operator import mul
from functools import reduce

facts = {}

for i in range(2,21):
  tmp = factorint(i)
  for d in tmp.keys():
    try:
      facts[d] = max(facts[d],tmp[d])
    except KeyError:
      facts[d] = tmp[d]

print(reduce(mul, [d**facts[d] for d in facts]))
