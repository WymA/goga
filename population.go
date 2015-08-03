// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"
	"sort"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
)

// Population holds all individuals
type Population []*Individual

// GetCopy returns a copy of this population
func (o Population) GetCopy() (pop Population) {
	ninds := len(o)
	pop = make([]*Individual, ninds)
	for i := 0; i < ninds; i++ {
		pop[i] = o[i].GetCopy()
	}
	return
}

// NewPopFloatChromo allocates a population made entirely of float point numbers
//  Input:
//   nbases -- number of bases in each float point gene
//   genes  -- all genes of all individuals [ninds][ngenes]
//  Output:
//   new population
func NewPopFloatChromo(nbases int, genes [][]float64) (pop Population) {
	ninds := len(genes)
	pop = make([]*Individual, ninds)
	for i := 0; i < ninds; i++ {
		pop[i] = NewIndividual(nbases, genes[i])
	}
	return
}

// NewPopReference creates a population based on a reference individual
//  Input:
//   ninds -- number of individuals to be generated
//   ref   -- reference individual with chromosome structure already set
//  Output:
//   new population
func NewPopReference(ninds int, ref *Individual) (pop Population) {
	pop = make([]*Individual, ninds)
	for i := 0; i < ninds; i++ {
		pop[i] = ref.GetCopy()
	}
	return
}

// NewPopRandom generates random population with individuals based on reference individual
// and gene values randomly drawn from Bingo.
//  Input:
//   ninds -- number of individuals to be generated
//   ref   -- reference individual with chromosome structure already set
//   bingo -- Bingo structure set with pool of values to draw gene values
//  Output:
//   new population
func NewPopRandom(ninds int, ref *Individual, bingo *Bingo) (pop Population) {
	pop = NewPopReference(ninds, ref)
	for i, ind := range pop {
		for j := 0; j < len(ind.Ints); j++ {
			ind.Ints[j] = bingo.DrawInt(i, j, ninds)
		}
		if ind.Floats != nil {
			for j := 0; j < ind.Nfltgenes; j++ {
				ind.SetFloat(j, bingo.DrawFloat(i, j, ninds))
			}
		}
		for j := 0; j < len(ind.Strings); j++ {
			ind.Strings[j] = bingo.DrawString(j)
		}
		for j := 0; j < len(ind.Keys); j++ {
			ind.Keys[j] = bingo.DrawKey(j)
		}
		for j := 0; j < len(ind.Bytes); j++ {
			copy(ind.Bytes[j], bingo.DrawBytes(j))
		}
		for j := 0; j < len(ind.Funcs); j++ {
			ind.Funcs[j] = bingo.DrawFunc(j)
		}
	}
	return
}

// NewPopFloatRandom generates a population of individuals with float point
// numbers only genes for given grid
//  Input:
//   C.Ninds  -- number of individuals to be generated
//   C.Nbases -- number of bases
//   C.Grid   -- whether or not to calc values based on grid;
//               otherwise select randomly between xmin and xmax
//   C.Noise  -- if noise>0, apply noise to move points away from grid nodes
//               noise is a multiplier; e.g. 0.2
//   xmin     -- min values of genes
//   xmax     -- max values of genes. len(xmin) = len(xmax) = ngenes
func NewPopFloatRandom(C *ConfParams, xmin, xmax []float64) (pop Population) {
	ngenes := len(xmin)
	chk.IntAssert(len(xmax), ngenes)
	ref := NewIndividual(C.Nbases, make([]float64, ngenes))
	pop = NewPopReference(C.Ninds, ref)
	pop.GenFloatRandom(C, xmin, xmax)
	return
}

// GenFloatRandom generates a population of individuals with float point
// numbers only genes for given grid
//  Input:
//   C.Ninds  -- number of individuals to be generated
//   C.Nbases -- number of bases
//   C.Grid   -- whether or not to calc values based on grid;
//               otherwise select randomly between xmin and xmax
//   C.Noise  -- if noise>0, apply noise to move points away from grid nodes
//               noise is a multiplier; e.g. 0.2
//   xmin     -- min values of genes
//   xmax     -- max values of genes. len(xmin) = len(xmax) = ngenes
func (o *Population) GenFloatRandom(C *ConfParams, xmin, xmax []float64) {
	if len(*o) < 2 {
		return
	}
	ngenes := len(xmin)
	chk.IntAssert(len(xmax), ngenes)
	chk.IntAssert(ngenes, (*o)[0].Nfltgenes)
	npts := int(math.Pow(float64(C.Ninds), 1.0/float64(ngenes))) // num points in 'square' grid
	ntot := int(math.Pow(float64(npts), float64(ngenes)))        // total num of individuals in grid
	den := 1.0                                                   // denominator to calculate dx
	if npts > 1 {
		den = float64(npts - 1)
	}
	var lfto int // leftover, e.g. n % (nx*ny)
	var rdim int // reduced dimension, e.g. (nx*ny)
	var idx int  // index of gene in grid
	var dx, x, mul float64
	for i := 0; i < C.Ninds; i++ {
		if i < ntot { // on grid
			lfto = i
			for j := 0; j < ngenes; j++ {
				rdim = int(math.Pow(float64(npts), float64(ngenes-1-j)))
				idx = lfto / rdim
				lfto = lfto % rdim
				dx = xmax[j] - xmin[j]
				x = xmin[j] + float64(idx)*dx/den
				if C.Noise > 0 {
					mul = rnd.Float64(0, C.Noise)
					if rnd.FlipCoin(0.5) {
						x += mul * x
					} else {
						x -= mul * x
					}
				}
				(*o)[i].SetFloat(j, x)
			}
		} else { // additional individuals
			for g := 0; g < ngenes; g++ {
				x = rnd.Float64(xmin[g], xmax[g])
				(*o)[i].SetFloat(g, x)
			}
		}
	}
}

// Len returns the length of the population == number of individuals
func (o Population) Len() int {
	return len(o)
}

// Swap swaps two individuals
func (o Population) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Less returns true if 'i' is "less bad" than 'j'; therefore it can be used
// to sort the population in increasing order of demerits: from best to worst
func (o Population) Less(i, j int) bool {
	return o[i].Demerit < o[j].Demerit
}

// Sort sorts the population from best to worst individuals; i.e. decreasing fitness values
func (o *Population) Sort() {
	sort.Sort(o)
}

// Output generates a nice table with population data
//  Input:
//  fmts      -- [ngenes] formats for int, flt, string, byte, bytes, and func
//               use fmts == nil to choose default ones
//  showBases -- show bases, if any
func (o Population) Output(fmts [][]string, showBases bool) (buf *bytes.Buffer) {

	// check
	if len(o) < 1 {
		return
	}

	// compute sizes and generate formats list
	if fmts == nil {
		sizes := make([][]int, 6)
		for _, ind := range o {
			sz := ind.GetStringSizes()
			for i := 0; i < 6; i++ {
				if len(sizes[i]) == 0 {
					sizes[i] = make([]int, len(sz[i]))
				}
				for j, s := range sz[i] {
					sizes[i][j] = imax(sizes[i][j], s)
				}
			}
		}
		fmts = make([][]string, 6)
		for i, str := range []string{"d", "g", "s", "x", "s", "s"} {
			fmts[i] = make([]string, len(sizes[i]))
			for j, sz := range sizes[i] {
				fmts[i][j] = io.Sf("%%%d%s", sz+1, str)
			}
		}
	}

	// compute sizes of header items
	szova, szoor, szdem := 0, 0, 0
	for _, ind := range o {
		szova = imax(szova, len(io.Sf("%g", ind.Ova)))
		szoor = imax(szoor, len(io.Sf("%g", ind.Oor)))
		szdem = imax(szdem, len(io.Sf("%g", ind.Demerit)))
	}
	szova = imax(szova, 3) // 3 ==> len("Ova")
	szoor = imax(szoor, 3) // 3 ==> len("Oor")
	szdem = imax(szdem, 7) // 7 ==> len("Demerit")

	// print individuals
	fmtova := io.Sf("%%%d", szova+1)
	fmtoor := io.Sf("%%%d", szoor+1)
	fmtdem := io.Sf("%%%d", szdem+1)
	line, sza, szb := "", 0, 0
	for i, ind := range o {
		stra := io.Sf(fmtova+"g", ind.Ova)
		if ind.Oor > 0 {
			stra = io.Sf(fmtova+"s", "n/a")
			stra += io.Sf(fmtoor+"g", ind.Oor)
		} else {
			stra += io.Sf(fmtoor+"s", "n/a")
		}
		stra += io.Sf(fmtdem+"g", ind.Demerit) + " "
		strb := ind.Output(fmts, showBases)
		line += stra + strb + "\n"
		if i == 0 {
			sza, szb = len(stra), len(strb)
		}
	}

	// write to buffer
	fmtgenes := io.Sf(" %%%d.%ds\n", szb, szb)
	n := sza + szb
	buf = new(bytes.Buffer)
	io.Ff(buf, printThickLine(n))
	io.Ff(buf, fmtova+"s", "Ova")
	io.Ff(buf, fmtoor+"s", "Oor")
	io.Ff(buf, fmtdem+"s", "Demerit")
	io.Ff(buf, fmtgenes, "Genes")
	io.Ff(buf, printThinLine(n))
	io.Ff(buf, line)
	io.Ff(buf, printThickLine(n))
	return
}

// OutFloatBases print bases of float genes
func (o Population) OutFloatBases(numFmt string) (l string) {
	for _, ind := range o {
		for _, val := range ind.Floats {
			l += io.Sf(numFmt, val)
		}
		l += "\n"
	}
	return
}
