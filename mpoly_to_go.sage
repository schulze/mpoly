def mpoly_go_head(n, nmod):
    """
    Returns as a string a file header of Go code for computation modulo nmod with
    multivariate polynomials in n variables.
    """
    go_head = """
package main

const (
    N = {0} // the number of variables
    Nmod = {1} // the modulus
    // FIXME: Have to DegBound manually for now; could get this from the the polys in Sage.
    DegBound = 40 
)

""".format(n, nmod)
    return go_head

def mpoly_to_go(f, nmod):
    """
    Convert the multivariate polynomial f in n variables into a Go
    mpoly. We can then compute/count points in Go.
    Reduces coefficients mod nmod.
    Returns a string of Go code for the polynomial.
    """
    monoms = f.monomials()
    coeffs = [f.monomial_coefficient(m) for m in monoms]
    coeffs = [Mod(c, nmod) for c in coeffs] # reduce mod nmod if necessary.

    go_coeffs = ''.join(["uint64({0}),".format(c) for c in coeffs])
    go_monoms = ''.join(["{0},".format(m.exponents()[0]) for m in monoms])
    go_monoms = go_monoms.replace('(', "{")
    go_monoms = go_monoms.replace(')', "}")

    return "&Poly{0}".format("{") + "[]uint64{0}{2}{1}, []Monom{0}{3}{1}".format("{", "}", go_coeffs, go_monoms) + "}"

def list_of_mpoly(poly_list, nmod):
    """
    Convert a list of polynomials into the Go representation thereof.
    Returns a string of Go code for the slice of polynomials.
    """
    return "var Polys = []*Poly{0}".format("{") + ''.join([mpoly_to_go(poly, nmod) + ", " for poly in poly_list]) + "}"

def ToGo(filename, poly_list, n, nmod):
    """"
    Convert the poly_list into a Go representation and write the resulting string to filename.
    """
    f = open(filename, "w")
    f.write(mpoly_go_head(n, nmod))
    f.write(list_of_mpoly(poly_list, nmod))
    f.close
