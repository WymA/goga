// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
)

func Test_pop01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("pop01")

	rnd.Init(0)

	// genes
	genes := [][]float64{
		{1, 5, -200},
		{1, 3, -300},
		{5, 7, -400},
		{1, 2, -500},
		{2, 4, -300},
	}

	// objective values and fitness values
	ovs := []float64{11, 21, 10, 12, 13}

	// init population
	nbases := 2
	pop := NewPopFloatChromo(nbases, genes)
	for i, ind := range pop {
		ind.ObjValue = ovs[i]
	}
	io.Pforan("%v\n", pop.Output(nil, false))

	// check floats and subfloats
	for i, ind := range pop {
		for j := 0; j < ind.Nfltgenes; j++ {
			chk.Scalar(tst, io.Sf("before: i%dg%d", i, j), 1e-12, ind.GetFloat(j), genes[i][j])
		}
	}

	// print bases
	io.Pf("\nbases (before)\n")
	io.Pf("%s\n", pop.OutFloatBases("%10.4f"))

	// change subfloats
	bases := [][]float64{
		{10, 1, 14, 5, -10, -1},
		{10, 1, 12, 1, -20, -1},
		{10, 5, 12, 1, -30, -1},
		{10, 1, 12, 2, -40, -1},
		{10, 2, 11, 1, -20, -1},
	}
	for i, b := range bases {
		copy(pop[i].Floats, b)
	}

	// print bases
	io.Pfyel("bases (after)\n")
	io.Pfyel("%s\n", pop.OutFloatBases("%4g"))

	// checkf floats
	chk.Scalar(tst, "after: i0g0", 1e-16, pop[0].GetFloat(0), 11)
	chk.Scalar(tst, "after: i0g1", 1e-16, pop[0].GetFloat(1), 19)
	chk.Scalar(tst, "after: i0g2", 1e-16, pop[0].GetFloat(2), -11)

	chk.Scalar(tst, "after: i1g0", 1e-16, pop[1].GetFloat(0), 11)
	chk.Scalar(tst, "after: i1g1", 1e-16, pop[1].GetFloat(1), 13)
	chk.Scalar(tst, "after: i1g2", 1e-16, pop[1].GetFloat(2), -21)

	chk.Scalar(tst, "after: i2g0", 1e-16, pop[2].GetFloat(0), 15)
	chk.Scalar(tst, "after: i2g1", 1e-16, pop[2].GetFloat(1), 13)
	chk.Scalar(tst, "after: i2g2", 1e-16, pop[2].GetFloat(2), -31)

	chk.Scalar(tst, "after: i3g0", 1e-16, pop[3].GetFloat(0), 11)
	chk.Scalar(tst, "after: i3g1", 1e-16, pop[3].GetFloat(1), 14)
	chk.Scalar(tst, "after: i3g2", 1e-16, pop[3].GetFloat(2), -41)

	chk.Scalar(tst, "after: i4g0", 1e-16, pop[4].GetFloat(0), 12)
	chk.Scalar(tst, "after: i4g1", 1e-16, pop[4].GetFloat(1), 12)
	chk.Scalar(tst, "after: i4g2", 1e-16, pop[4].GetFloat(2), -21)
}

func Test_pop02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("pop02")

	rnd.Init(0)

	genes := [][]float64{
		{1, 5}, // 0
		{1, 3}, // 1
		{5, 7}, // 2
		{1, 2}, // 3
		{2, 4}, // 4
		{3, 6}, // 5
		{4, 8}, // 6
		{4, 6}, // 7
		{1, 3}, // 8
		{0, 0}, // 9
	}

	// objective values and fitness values
	//                0   1   2   3   4   5   6     7     8     9
	ovs := []float64{11, 21, 10, 12, 13, 31, 41, 11.1, 31.5, 11.5}

	// init population
	nbases := 2
	pop := NewPopFloatChromo(nbases, genes)
	for i, ind := range pop {
		ind.ObjValue = ovs[i]
	}
	io.Pforan("%v\n", pop.Output(nil, false))

	pop.Sort()

	io.Pfyel("%v\n", pop.Output(nil, false))

	genes_sorted := [][]float64{
		{5, 7}, // 2
		{1, 5}, // 0
		{4, 6}, // 7
		{0, 0}, // 9
		{1, 2}, // 3
		{2, 4}, // 4
		{1, 3}, // 1
		{3, 6}, // 5
		{1, 3}, // 8
		{4, 8}, // 6
	}

	for i, ind := range pop {
		for j := 0; j < ind.Nfltgenes; j++ {
			chk.Scalar(tst, io.Sf("i%dg%d", i, j), 1e-14, ind.GetFloat(j), genes_sorted[i][j])
		}
	}
}

func Test_pop03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("pop03")

	rnd.Init(0)

	nbases := 3
	ind := get_individual(0, nbases)

	fmts := [][]string{{"%4d"}, {"%8.2f"}, {" %6.6s"}, {" %x"}, {"%6.6s"}, {"%3s"}}

	ninds := 5
	bingo := NewExampleBingo()
	pop := NewPopRandom(ninds, ind, bingo)
	io.Pf("\n%v\n", pop.Output(fmts, false))

	ngenes := 3
	flts := la.MatAlloc(ninds, ngenes)
	for i := 0; i < ninds; i++ {
		for j := 0; j < ngenes; j++ {
			flts[i][j] = pop[i].GetFloat(j)
		}
	}
	chk.Matrix(tst, "flts", 1e-13, flts, [][]float64{
		{-123, -1, 0},
		{-12, -0.5, 0.25},
		{99, 0, 0.5},
		{210, 0.5, 0.75},
		{321, 1, 1},
	})
}
