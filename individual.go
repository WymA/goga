// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Individual implements one individual in a population
type Individual struct {

	// data
	Ovas      []float64 // objective values
	Oors      []float64 // out-of-range values: sum of positive distances from constraints
	Demerit   float64   // quantity for comparing individuals. 0=good 1=bad 2=worse(oor) 3=worst(oor)
	Nfltgenes int       // number of floats == number of float64 genes
	Nbases    int       // number of bases to split Floats

	// chromosome
	Ints    []int     // integers
	Floats  []float64 // floats [nFLTgenes * nbases]
	Strings []string  // strings
	Keys    []byte    // 1D bytes
	Bytes   [][]byte  // 2D bytes
	Funcs   []Func_t  // functions
}

// NewIndividual allocates a new individual
//  Input:
//   nbases -- used to split genes of floats into smaller parts
//   slices -- slices of ints, floats, strings, bytes, []bytes, and/or Func_t
//  Notes:
//   1) the slices in 'genes' can all be combined to define genes with mixed data;
//   2) the slices can also be nil, except for one of them.
func NewIndividual(nova, noor, nbases int, slices ...interface{}) (o *Individual) {
	o = new(Individual)
	o.Ovas = make([]float64, nova)
	o.Oors = make([]float64, noor)
	for _, slice := range slices {
		switch s := slice.(type) {
		case []int:
			o.Ints = make([]int, len(s))
			copy(o.Ints, s)

		case []float64:
			o.Nfltgenes = len(s)
			o.Nbases = nbases
			if o.Nbases > 1 {
				o.Floats = SimpleChromo(s, nbases)
			} else {
				o.Floats = make([]float64, o.Nfltgenes*o.Nbases)
				copy(o.Floats, s)
			}

		case []string:
			o.Strings = make([]string, len(s))
			copy(o.Strings, s)

		case []byte:
			o.Keys = make([]byte, len(s))
			copy(o.Keys, s)

		case [][]byte:
			o.Bytes = make([][]byte, len(s))
			for i, x := range s {
				o.Bytes[i] = make([]byte, len(x))
				copy(o.Bytes[i], x)
			}

		case []Func_t:
			o.Funcs = make([]Func_t, len(s))
			copy(o.Funcs, s)
		}
	}
	return
}

// GetCopy returns a copy of this individual
func (o Individual) GetCopy() (x *Individual) {

	x = new(Individual)
	x.Ovas = make([]float64, len(o.Ovas))
	x.Oors = make([]float64, len(o.Oors))
	copy(x.Ovas, o.Ovas)
	copy(x.Oors, o.Oors)
	x.Demerit = o.Demerit
	x.Nfltgenes = o.Nfltgenes
	x.Nbases = o.Nbases

	if o.Ints != nil {
		x.Ints = make([]int, len(o.Ints))
		copy(x.Ints, o.Ints)
	}

	if o.Floats != nil {
		x.Floats = make([]float64, len(o.Floats))
		copy(x.Floats, o.Floats)
	}

	if o.Strings != nil {
		x.Strings = make([]string, len(o.Strings))
		copy(x.Strings, o.Strings)
	}

	if o.Keys != nil {
		x.Keys = make([]byte, len(o.Keys))
		copy(x.Keys, o.Keys)
	}

	if o.Bytes != nil {
		x.Bytes = make([][]byte, len(o.Bytes))
		for i, b := range o.Bytes {
			x.Bytes[i] = make([]byte, len(b))
			copy(x.Bytes[i], b)
		}
	}

	if o.Funcs != nil {
		x.Funcs = make([]Func_t, len(o.Funcs))
		copy(x.Funcs, o.Funcs)
	}
	return
}

// CopyInto copies this individual's data into another individual
func (o Individual) CopyInto(x *Individual) {

	copy(x.Ovas, o.Ovas)
	copy(x.Oors, o.Oors)
	x.Demerit = o.Demerit
	x.Nfltgenes = o.Nfltgenes
	x.Nbases = o.Nbases

	if o.Ints != nil {
		copy(x.Ints, o.Ints)
	}

	if o.Floats != nil {
		copy(x.Floats, o.Floats)
	}

	if o.Strings != nil {
		copy(x.Strings, o.Strings)
	}

	if o.Keys != nil {
		copy(x.Keys, o.Keys)
	}

	if o.Bytes != nil {
		for i, b := range o.Bytes {
			copy(x.Bytes[i], b)
		}
	}

	if o.Funcs != nil {
		copy(x.Funcs, o.Funcs)
	}
	return
}

// IndCompare compares individual 'A' with another one 'B'
func IndCompare(A, B *Individual, φ float64) (A_dominates bool) {
	var A_is_unfeasible, B_is_unfeasible bool
	for i := 0; i < len(A.Oors); i++ {
		if A.Oors[i] > 0 {
			A_is_unfeasible = true
		}
		if B.Oors[i] > 0 {
			B_is_unfeasible = true
		}
	}
	if A_is_unfeasible {
		if B_is_unfeasible {
			return utl.DblsParetoMinProb(A.Oors, B.Oors, φ)
		}
		return false // B dominates
	}
	if B_is_unfeasible {
		return true // A dominates
	}
	return utl.DblsParetoMinProb(A.Ovas, B.Ovas, φ)
}

// IndDistance computes a distance measure from individual 'A' to another individual 'B'
func IndDistance(A, B *Individual) (dist float64) {
	nints := len(A.Ints)
	dints := 0.0
	for i := 0; i < nints; i++ {
		dints += math.Abs(float64(A.Ints[i] - B.Ints[i]))
	}
	if nints > 0 {
		dints /= float64(nints)
	}
	nflts := len(A.Floats)
	dflts := 0.0
	for i := 0; i < nflts; i++ {
		dflts += math.Abs(A.Floats[i] - B.Floats[i])
	}
	if nflts > 0 {
		dflts /= float64(nflts)
	}
	return dints + dflts
}

// genetic algorithm routines //////////////////////////////////////////////////////////////////////

// IndCrossover performs the crossover between chromosomes of two individuals A and B
// resulting in the chromosomes of other two individuals a and b
//  Input:
//   A and B -- parents
//   ncuts   -- number of cuts. keys are: 'int', 'flt', 'str', 'key', 'byt', 'fun'
//              ncuts can be nil if 'cuts' is provided
//   cuts    -- positions for cuts in the augmented/whole chromosome
//              len(cuts) == 6: {int, flt, str, key, byt, fun
//              cuts == nil indicates ncuts is to be used instead
//   probs   -- probabilities. use nil for default values
//   cxfucns -- crossover functions. use nil for default ones
//  Output:
//   a and b -- offspring
func IndCrossover(a, b, A, B *Individual, time int, ncuts map[string]int, cuts map[string][]int, probs map[string]float64, extra map[string]interface{},
	cxint CxIntFunc_t, cxflt CxFltFunc_t, cxstr CxStrFunc_t, cxkey CxKeyFunc_t, cxbyt CxBytFunc_t, cxfun CxFunFunc_t) {

	// default values
	pc := func(t string) float64 {
		if val, ok := probs[t]; ok {
			return val
		}
		return 0.8
	}

	// default functions
	if cxint == nil {
		cxint = IntCrossover
	}
	if cxflt == nil {
		cxflt = FltCrossover
	}
	if cxstr == nil {
		cxstr = StrCrossover
	}
	if cxkey == nil {
		cxkey = KeyCrossover
	}
	if cxbyt == nil {
		cxbyt = BytCrossover
	}
	if cxfun == nil {
		cxfun = FunCrossover
	}

	// perform crossover
	if A.Ints != nil {
		cxint(a.Ints, b.Ints, A.Ints, B.Ints, time, ncuts["int"], cuts["int"], pc("int"), extra["int"])
	}
	if A.Floats != nil {
		cxflt(a.Floats, b.Floats, A.Floats, B.Floats, time, ncuts["flt"], cuts["flt"], pc("flt"), extra["flt"])
	}
	if A.Strings != nil {
		cxstr(a.Strings, b.Strings, A.Strings, B.Strings, time, ncuts["str"], cuts["str"], pc("str"), extra["str"])
	}
	if A.Keys != nil {
		cxkey(a.Keys, b.Keys, A.Keys, B.Keys, time, ncuts["key"], cuts["key"], pc("key"), extra["key"])
	}
	if A.Bytes != nil {
		cxbyt(a.Bytes, b.Bytes, A.Bytes, B.Bytes, time, ncuts["byt"], cuts["byt"], pc("byt"), extra["byt"])
	}
	if A.Funcs != nil {
		cxfun(a.Funcs, b.Funcs, A.Funcs, B.Funcs, time, ncuts["fun"], cuts["fun"], pc("fun"), extra["fun"])
	}
}

// IndMutation performs the mutation operation in the chromosomes of an individual
//  Input:
//   A        -- individual
//   nchanges -- number of changes. keys are: 'int', 'flt', 'str', 'key', 'byt', 'fun'
//               use nil for default values
//   probs    -- probabilities. use nil for default values
//   extra    -- extra arguments for each 'int', 'flt', 'str', 'key', 'byt', 'fun'
//   mutfucns -- mutation functions. use nil for default ones
//  Output: modified individual
func IndMutation(A *Individual, time int, nchanges map[string]int, probs map[string]float64, extra map[string]interface{},
	mtint MtIntFunc_t, mtflt MtFltFunc_t, mtstr MtStrFunc_t, mtkey MtKeyFunc_t, mtbyt MtBytFunc_t, mtfun MtFunFunc_t) {

	// default values
	nc := func(t string) int {
		if val, ok := nchanges[t]; ok {
			return val
		}
		return 1
	}
	pm := func(t string) float64 {
		if val, ok := probs[t]; ok {
			return val
		}
		return 0.01
	}

	// default functions
	if mtint == nil {
		mtint = IntMutation
	}
	if mtflt == nil {
		mtflt = FltMutation
	}
	if mtstr == nil {
		mtstr = StrMutation
	}
	if mtkey == nil {
		mtkey = KeyMutation
	}
	if mtbyt == nil {
		mtbyt = BytMutation
	}
	if mtfun == nil {
		mtfun = FunMutation
	}

	// perform crossover
	if A.Ints != nil {
		mtint(A.Ints, time, nc("int"), pm("int"), extra["int"])
	}
	if A.Floats != nil {
		mtflt(A.Floats, time, nc("flt"), pm("flt"), extra["flt"])
	}
	if A.Strings != nil {
		mtstr(A.Strings, time, nc("flt"), pm("str"), extra["str"])
	}
	if A.Keys != nil {
		mtkey(A.Keys, time, nc("key"), pm("key"), extra["key"])
	}
	if A.Bytes != nil {
		mtbyt(A.Bytes, time, nc("byt"), pm("byt"), extra["byt"])
	}
	if A.Funcs != nil {
		mtfun(A.Funcs, time, nc("fun"), pm("fun"), extra["fun"])
	}
}

// handle bases ////////////////////////////////////////////////////////////////////////////////////

// SetFloat returns the float corresponding to gene 'i'
//  igene -- is the index of gene/float in [0, Nfltgenes]
func (o *Individual) SetFloat(igene int, x float64) {
	if o.Nbases > 1 {
		values := make([]float64, o.Nbases)
		rnd.Float64s(values, 0, 1)
		sum := la.VecAccum(values)
		for j := 0; j < o.Nbases; j++ {
			o.Floats[igene*o.Nbases+j] = x * values[j] / sum
		}
		return
	}
	o.Floats[igene] = x
}

// GetFloat returns the float corresponding to gene 'i'
//  igene -- is the index of gene/float in [0, Nfltgenes]
func (o Individual) GetFloat(igene int) (x float64) {
	if o.Nbases > 1 {
		for j := 0; j < o.Nbases; j++ {
			x += o.Floats[igene*o.Nbases+j]
		}
		return
	}
	return o.Floats[igene]
}

// GetFloats returns all float genes
func (o Individual) GetFloats() (x []float64) {
	x = make([]float64, o.Nfltgenes)
	for i := 0; i < o.Nfltgenes; i++ {
		x[i] = o.GetFloat(i)
	}
	return
}

// output //////////////////////////////////////////////////////////////////////////////////////////

// GetStringSizes returns the sizes of strings representing each gene type
//  sizes -- [6][...] sizes of strings for {int, flt, string, byte, bytes, func}
func (o *Individual) GetStringSizes() (sizes [][]int) {

	sizes = make([][]int, 6)
	if o.Ints != nil {
		sizes[0] = make([]int, len(o.Ints))
		for i, x := range o.Ints {
			sizes[0][i] = utl.Imax(sizes[0][i], len(io.Sf("%v", x)))
		}
	}

	if o.Floats != nil {
		sizes[1] = make([]int, o.Nfltgenes)
		for i := 0; i < o.Nfltgenes; i++ {
			x := o.Floats[i]
			if o.Nbases > 1 {
				x = 0
				for j := 0; j < o.Nbases; j++ {
					x += o.Floats[i*o.Nbases+j]
				}
			}
			sizes[1][i] = utl.Imax(sizes[1][i], len(io.Sf("%v", x)))
		}
	}

	if o.Strings != nil {
		sizes[2] = make([]int, len(o.Strings))
		for i, x := range o.Strings {
			sizes[2][i] = utl.Imax(sizes[2][i], len(io.Sf("%v", x)))
		}
	}

	if o.Keys != nil {
		sizes[3] = make([]int, len(o.Keys))
		for i, x := range o.Keys {
			sizes[3][i] = utl.Imax(sizes[3][i], len(io.Sf("%v", x)))
		}
	}

	if o.Bytes != nil {
		sizes[4] = make([]int, len(o.Bytes))
		for i, x := range o.Bytes {
			sizes[4][i] = utl.Imax(sizes[4][i], len(io.Sf("%v", string(x))))
		}
	}

	if o.Funcs != nil {
		sizes[5] = make([]int, len(o.Funcs))
		for i, x := range o.Funcs {
			sizes[5][i] = utl.Imax(sizes[5][i], len(io.Sf("%v", x(o))))
		}
	}
	return
}

// Output returns a string representation of this individual
//  fmts      -- [6][...] formats of strings for {int, flt, string, byte, bytes, func}
//               use fmts == nil to choose default ones
//  showBases -- show bases, if any
func (o *Individual) Output(fmts [][]string, showBases bool) (l string) {

	if fmts == nil {
		fmts = [][]string{{" %d"}, {" %g"}, {" %q"}, {" %x"}, {" %q"}, {" %q"}}
	}

	fmt := func(itype, idx int) (s string) {
		s = fmts[itype][0]
		if idx < len(fmts[itype]) {
			s = fmts[itype][idx]
		}
		return
	}

	for i, x := range o.Ints {
		l += io.Sf(fmt(0, i), x)
	}

	for i := 0; i < o.Nfltgenes; i++ {
		x := o.Floats[i]
		if o.Nbases > 1 {
			x = 0
			for j := 0; j < o.Nbases; j++ {
				x += o.Floats[i*o.Nbases+j]
			}
		}
		l += io.Sf(fmt(1, i), x)
	}

	for i, x := range o.Strings {
		l += io.Sf(fmt(2, i), x)
	}

	for i, x := range o.Keys {
		l += io.Sf(fmt(3, i), x)
	}

	for i, x := range o.Bytes {
		l += io.Sf(fmt(4, i), string(x))
	}

	for i, x := range o.Funcs {
		l += io.Sf(fmt(5, i), x(o))
	}

	if showBases && len(o.Floats) > 0 {
		for i, x := range o.Floats {
			if i%o.Nbases == 0 {
				if i == 0 {
					l += " ||"
				} else {
					l += " |"
				}
			}
			l += io.Sf("%11.3e", x)
		}
	}

	return
}
