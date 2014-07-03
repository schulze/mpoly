package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// Using flint's nmods might be safer for larger moduli, but is
	// slower for small moduli. Instead we pre-compute a table of
	// powers up to some bound and use these to evaluate monomials.
	//nmod "github.com/frithjof-schulze/flint.go/extras"

	"runtime"
	"runtime/pprof"
)

// From the Sage program we get a definition of
// 	N = #of variabeles
// 	Nmod = modulus
// 	DegBound = upper bound for the exponents in the polynomials

//var (
//	Ninv uint64 = nmod.Preinvert(Nmod)
//)

type Point [N]uint64
type Monom [N]uint64

type Poly struct {
	coeffs []uint64
	monoms []Monom
}

var (
	PowersTable [][]uint64 = make([][]uint64, Nmod)
	MulTable    [][]uint64 = make([][]uint64, Nmod)
)

func InitPowers() {
	for i := range PowersTable {
		PowersTable[i] = make([]uint64, DegBound)
		for j := range PowersTable[i] {
			//PowersTable[i][j] = nmod.PowMod2Preinv(uint64(i), int64(j), Nmod, Ninv)
			PowersTable[i][j] = PowMod(uint64(i), uint64(j))
		}
	}
}

func InitMuls() {
	for i := range MulTable {
		MulTable[i] = make([]uint64, Nmod)
		for j := range MulTable[i] {
			MulTable[i][j] = uint64(i*j) % Nmod
		}
	}
}

func PowMod(a, exp uint64) uint64 {
	val := uint64(1)
	for i := uint64(0); i < exp; i++ {
		val = (val * a) % Nmod
	}
	return val
}

func (m Monom) Eval(pt *Point) uint64 {
	// TODO: Can we do some caching here?
	var val, tmp uint64
	val = 1
	for i, exp := range m {
		tmp = PowersTable[pt[i]][exp]
		//val = nmod.MulMod2Preinv(val, tmp, Nmod, Ninv)
		val = MulTable[val][tmp]
	}
	return val
}

func CopyPoint(a Point) *Point {
	b := *new(Point)
	for i := 0; i < N; i++ {
		b[i] = a[i]
	}
	return &b
}

// Generate sends the finite sequence of points to the channel ch.
// After the last point send a single nil.
// Based on Algorithm H in Knuth's 'Generating all n-tuples'.
func Generate() chan *Point {
	// init
	var pt Point
	focus := new([N + 1]int)
	orient := new([N]int)
	for i := 0; i < N; i++ {
		focus[i] = i
		orient[i] = 1
	}
	focus[N] = N
	out := make(chan *Point)

	go func() {
		for {
			out <- CopyPoint(pt)
			j := focus[0]
			focus[0] = 0
			if j == N { // after the last point
				out <- nil
				return
			}
			pt[j] += uint64(orient[j])
			if pt[j] == 0 || pt[j] == Nmod-1 {
				orient[j] = -orient[j]
				focus[j] = focus[j+1]
				focus[j+1] = j + 1
			}
		}
	}()
	return out
}

func FilterForPoly(f *Poly, in chan *Point) chan *Point {
	// Get a point pt, test if f(pt) == 0 and pass pt on if this is true.
	// If you get a point nil, stop.
	out := make(chan *Point)
	filter := func() {
		for pt := <-in; pt != nil; pt = <-in {
			if Eval(f, pt) == 0 {
				out <- pt
			}
		}
		out <- nil
		return
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		go filter()
	}
	return out
}

// Eval evaluates f at pt and return the result mod Nmod.
func Eval(f *Poly, pt *Point) uint64 {
	var val, tmp uint64
	for i, m := range f.monoms {
		tmp = m.Eval(pt)
		//tmp = nmod.MulMod2Preinv(tmp, f.coeffs[i], Nmod, Ninv)
		//val = nmod.AddMod(val, tmp, Nmod)
		tmp = MulTable[tmp][f.coeffs[i]]
		val = val + tmp
		if val >= Nmod {
			val = val - Nmod
		}
	}
	return val
}

func RationalPoints(polys []*Poly) chan *Point {
	out := Generate()
	for _, f := range polys {
		out = FilterForPoly(f, out)
	}
	return out
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	points := make([]Point, 0, 10)
	InitPowers()
	InitMuls()

	//runtime.GOMAXPROCS(runtime.NumCPU())
	sols := RationalPoints(Polys)
	for {
		pt := <-sols
		if pt == nil {
			break
		}
		points = append(points, *pt)
	}
	for _, pt := range points {
		fmt.Print(pt)
		fmt.Println(",")
	}
	return
}
