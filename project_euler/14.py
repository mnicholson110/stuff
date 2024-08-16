def collatzCount(n):
  count = 1
  while n != 1:
    if n % 2 == 0:
      n = n/2
      count += 1
    else:
      n = 3*n + 1
      count += 1
  return count

m = 0
res = 0
for i in range(1,1000000,2):
  if collatzCount(i) > m:
    m = collatzCount(i)
    res = i

print(res)
