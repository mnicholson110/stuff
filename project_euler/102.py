class point:
    def __init__(self, x, y):
        self.x = int(x)
        self.y = int(y)
    def __str__(self):
        return "(" + str(self.x) + "," + str(self.y) + ")"
    def __repr__(self):
        return self.__str__()

# returns true if the path from p0 to p1 to p2 results in a counterclockwise movement
def ccw(p0, p1, p2):
    dx1 = p1.x - p0.x
    dy1 = p1.y - p0.y
    dx2 = p2.x - p0.x
    dy2 = p2.y - p0.y
    if (dx1*dy2 > dy1*dx2): return True
    if (dx1*dy2 < dy1*dx2): return False
    if ((dx1*dx2 < 0) or (dy1*dy2 < 0)): return False
    if (dx1 * dy2 > dy1 * dx2): return True
    if ((dx1*dx1 + dy1*dy1) < (dx2*dx2 + dy2*dy2)): return True

# returns true if the segments ab and cd intersect
def intersect(a, b, c, d):
    return ccw(a,c,d) != ccw(b,c,d) and ccw(a,b,c) != ccw(a,b,d)

f = open('triangles.txt')
coords = []
for row in f:
    tmp = row.strip().split(',')
    coords += [[point(tmp[0],tmp[1]),point(tmp[2],tmp[3]),point(tmp[4],tmp[5])]]
f.close()

s = 0
origin = point(0,0)
max_x = point(1200,0)

for tri in coords:
    ints = 0
    if intersect(tri[0],tri[1],origin,max_x): ints +=1
    if intersect(tri[1],tri[2],origin,max_x): ints +=1
    if intersect(tri[0],tri[2],origin,max_x): ints +=1
    if ints % 2 == 1:
        s +=1

print(s)

 
