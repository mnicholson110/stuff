def rollUpMax(tri):
    for r in range(len(tri)-2,-1,-1):
        for i in range(0,r+1):
            tri[r][i] += max(tri[r+1][i],tri[r+1][i+1])
    return tri[0]

f = open('triangle.txt')
tri = []
for row in f:
    tri += [[int(i) for i in row.split()]]
f.close()

