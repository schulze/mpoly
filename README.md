mpoly
=====

Brute force solving of systems of multivariate polynomials in
F_p.

`mpoly_to_go.sage` contains code to write out a list of
polynomials as Go code.

`mpoly.go` contains the code to search for solutions by
enumerating points in the affine space `A^n(F_p)`.

`mpoly_to_go_test.sage` contains an small example (that in fact
could have been solved by other means). The example is taken from
the notes of Noam Elkies lectures at the  Clay Mathematics
Institute's Summer School on arithmetic geometry at Goettingen,
2006 (http://www.math.harvard.edu/~elkies/gottingen_notes.txt).

To run the example do:

	> sage mpoly_to_go_test.sage 
	> go build .
	> ./mpoly
	[...]
	[[0 0 2],
	[10 1 2],  <-- this is the solution we are looking for
	[0 0 3],
	[...]
	> 

In general it is easy to filter out the 'singular' solutions back
in the Sage world.
