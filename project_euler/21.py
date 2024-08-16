from sympy import divisors

def isAmicable(n):
  a1 = sum(divisors(n)) - n
  a2 = sum(divisors(a1)) - a1
  if a1 == a2:
    return False
  elif a2 == n:
    return True
  return False

def func21(arg):
  amiCheck = [False]*arg
  amiCheck[0] = amiCheck[1] = True

  for (n, checked) in enumerate(amiCheck):
    if not checked:
      if isAmicable(n):
        yield n
        yield sum(divisors(n)) - n
        amiCheck[int(sum(divisors(n)) - n)] = True
      amiCheck[n] = True

print(sum([i for i in func21(10000)]))
