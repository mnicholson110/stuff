from sympy.solvers.diophantine import diop_DN as sol

largest_x = 0

for D in range(1001):
  x = sol(D,1)[0][0]
  largest_x = max(x,largest_x)
  if x == largest_x: res = D

print(res)
