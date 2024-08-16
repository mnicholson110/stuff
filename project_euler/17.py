singles = ["one", "two", "three", "four", "five", "six", "seven", "eight", "nine"]

teens = ["eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"]

tens = ["twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"]

def build20to100():
  return [a + " " + b for a in tens for b in singles]

def build101to999():
  return [a + " hundred and " + b for a in singles for b in singles + teens + build20to100()]

def buildtherest():
  return ["ten"] + [a + " hundred and ten" for a in singles] + [a + " hundred" for a in singles] + [a + " hundred and " + b for a in singles for b in tens] + ["one thousand"]

countemup = 0

for num in (singles + teens + tens + build20to100() + build101to999() + buildtherest()):
  countemup += sum([len(x) for x in num.split(" ")])

print(countemup)
