from usefulStuff import divisorPairs

for r in range(2,1000,2):
  divs = divisorPairs((r**2)/2)
  for pair in divs:
    x = r + pair[0]
    y = r + pair[1]
    z = r + pair[0] + pair[1]
    if x+y+z == 1000:
      print x*y*z
      break
