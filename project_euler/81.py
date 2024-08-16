matrix = []
f = open('matrix.txt')
for row in f:
    matrix += [[int(i) for i in row.split(',')]]
f.close()

for i in range(len(matrix)-2,-1,-1):
    matrix[len(matrix)-1][i] += matrix[len(matrix)-1][i+1]
    matrix[i][len(matrix)-1] += matrix[i+1][len(matrix)-1]

for col in range(len(matrix)-2,-1,-1):
    for row in range(len(matrix)-2,-1,-1):
        matrix[row][col] += min(matrix[row][col+1],matrix[row+1][col])

print(matrix[0][0])
