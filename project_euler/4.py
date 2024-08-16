def isPalindrome(n):
  if str(n) == str(n)[::-1]: return True
  else: return False

def productOfThreeDigits(n):
  for i in range(100,1000):
    if n / i in range(100,1000):
      return True
  return False

for n in range(999*999,0,-1):
  if (isPalindrome(n) and productOfThreeDigits(n)):
    print(n)
    break

