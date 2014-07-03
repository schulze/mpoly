load("mpoly_to_go.sage")

Aff.<t, b0, b1, d0> = AffineSpace(QQ, 4)
K = Aff.coordinate_ring()

d = t + d0
c = (3*b1^2/2)*t^2 + (18*b0 - b1)*b1*t/6 + (9*b1^2*d0 + (9*b0 - b1)^2) / 54
bp = b1*t + b0
b = bp + 2*c*d
a = 1 + 6*d*(bp+c*d)

A = a^2 + 2*b
B = a^3 + 3*a*b + c

E = EllipticCurve(K, [-3*A,2*B])
D = E.discriminant()
f = A^3 - B^2

R.<b0, b1, d0> = PolynomialRing(QQ, 3)

l = [R(f.coefficient({t:i})) for i in range(3)]

ToGo("example_K3.go", l, 3, 13)
