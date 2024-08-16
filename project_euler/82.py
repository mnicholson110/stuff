matrix = []
f = open('matrix.txt')
for row in f:
    matrix += [[int(i) for i in row.split(',')]]
f.close()

def colMin(m,row,col):
    elem = m[row][col]
    try:
        res = [elem + m[row][col+1]]
    except IndexError:
        return m[row][col]
    n = row - 1
    tmp = 0
    while n >= 0:
        tmp += m[n][col]
        res += [elem + tmp + m[n][col+1]]
        n -= 1
    n = row + 1
    tmp = 0
    while n < len(m):
        tmp += m[n][col]
        res += [elem + tmp + m[n][col+1]]
        n += 1
    return min(res)


for col in range(len(matrix)-2,-1,-1):
    tmp = []
    for row in range(len(matrix)):
        tmp += [colMin(matrix,row,col)]
    for row in range(len(matrix)):
        matrix[row][col] = tmp[row]

res = 10**15
for row in range(len(matrix)):
    if matrix[row][0] < res:
        res = matrix[row][0]

print(res)

