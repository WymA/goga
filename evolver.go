// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/graph"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Evolver realises the evolutionary process
type Evolver struct {

	// data
	C       *ConfParams // configuration parameters
	Islands []*Island   // islands
	Best    *Individual // best individual among all in all islands

	// callbacks
	OutFcn func(t int, evo *Evolver) // output function callback

	// migration
	OvaMin []float64     // min ova
	OvaMax []float64     // max ova
	Mdist  [][]float64   // match distances
	Match  graph.Munkres // matches

	// migration v2: TODO: remove this
	Fsizes      []int           // front sizes
	Fronts      [][]*Individual // non-dominated fronts
	Competitors []*Individual   // competitors
}

// NewEvolverPop creates a new evolver based on given populations
func NewEvolver(C *ConfParams) (o *Evolver) {

	// check input
	C.check_input()

	// data
	o = new(Evolver)
	o.C = C
	o.Islands = make([]*Island, o.C.Nisl)
	if o.C.Pll {
		done := make(chan int, o.C.Nisl)
		for i := 0; i < o.C.Nisl; i++ {
			go func(idx int) {
				o.Islands[idx] = NewIsland(idx, o.C)
				done <- 1
			}(i)
		}
		for i := 0; i < o.C.Nisl; i++ {
			<-done
		}
	} else {
		for i := 0; i < o.C.Nisl; i++ {
			o.Islands[i] = NewIsland(i, o.C)
		}
	}
	o.Best = o.Islands[0].Pop[0]

	// migration
	o.OvaMin = make([]float64, o.C.Nova)
	o.OvaMax = make([]float64, o.C.Nova)
	o.Mdist = la.MatAlloc(o.C.Nimig, o.C.Nimig)
	o.Match.Init(o.C.Nimig, o.C.Nimig)

	// migration v2: TODO: remove this
	nc := o.C.Nisl * o.C.Ninds // number of competitors
	o.Fsizes = make([]int, nc)
	o.Fronts = make([][]*Individual, nc)
	for i := 0; i < nc; i++ {
		o.Fronts[i] = make([]*Individual, nc)
	}
	o.Competitors = make([]*Individual, nc)
	k := 0
	for _, isl := range o.Islands {
		for _, ind := range isl.Pop {
			o.Competitors[k] = ind
			k++
		}
	}
	return
}

// Run runs the evolution process
func (o *Evolver) Run() {

	// check
	if o.C.Nisl < 1 {
		return
	}
	if o.C.Ninds < o.C.Nisl {
		chk.Panic("number of individuals must be greater than the number of islands")
	}

	// first output
	t := 0
	for _, isl := range o.Islands {
		isl.WritePopToReport(t, 0)
	}
	if o.C.Verbose {
		o.print_legend()
		io.Pf("\nrunning ...\n")
	}
	if o.C.PostProc != nil {
		for _, isl := range o.Islands {
			o.C.PostProc(isl.Id, 0, isl.Pop)
		}
	}

	// time loop
	t = 1
	tmig := o.C.Dtmig
	nomig := false
	if tmig > o.C.Tf {
		tmig = o.C.Tf
		nomig = true
	}
	done := make(chan int, o.C.Nisl)
	for t < o.C.Tf {

		// evolve up to migration time
		if o.C.Pll {
			for i := 0; i < o.C.Nisl; i++ {
				go func(isl *Island) {
					for time := t; time < tmig; time++ {
						report := o.calc_report(time)
						isl.Run(time, report, (o.C.Verbose && isl.Id == 0))
						if o.C.Verbose && isl.Id == 0 {
							o.print_time(time, report)
						}
					}
					done <- 1
				}(o.Islands[i])
			}
			for i := 0; i < o.C.Nisl; i++ {
				<-done
			}
		} else {
			for _, isl := range o.Islands {
				for time := t; time < tmig; time++ {
					report := o.calc_report(time)
					isl.Run(time, report, (o.C.Verbose && isl.Id == 0))
					if o.C.Verbose && isl.Id == 0 {
						o.print_time(time, report)
					}
				}
			}
		}

		// update time
		t = tmig
		tmig += o.C.Dtmig
		if tmig > o.C.Tf {
			tmig = o.C.Tf
			nomig = true
		}

		// output
		if o.OutFcn != nil {
			o.OutFcn(t, o)
		}

		// skip migration
		if nomig {
			continue
		}

		// migration
		if o.C.Verbose {
			io.Pfyel(" %d", t)
		}
		o.migration(t)
	}

	// best individual
	o.FindBestFromAll()

	// message
	if o.C.Verbose {
		io.Pf("\n... end\n\n")
	}

	// write reports
	if o.C.DoReport {
		for _, isl := range o.Islands {
			isl.SaveReport(true)
		}
	}
	return
}

// FindBestFromAll finds best individual from all islands
//  Output: o.Best will point to the best individual
func (o *Evolver) FindBestFromAll() {
	if len(o.Islands) < 1 {
		return
	}
	o.Best = o.Islands[0].Pop[0]
	for i := 1; i < o.C.Nisl; i++ {
		_, other_is_better := IndCompareDet(o.Best, o.Islands[i].Pop[0])
		if other_is_better {
			o.Best = o.Islands[i].Pop[0]
		}
	}
}

// GetFeasible returns all feasible individuals from all islands
func (o *Evolver) GetFeasible() (feasible []*Individual) {
	for _, isl := range o.Islands {
		for _, ind := range isl.Pop {
			unfeasible := false
			for _, oor := range ind.Oors {
				if oor > 0 {
					unfeasible = true
				}
			}
			if !unfeasible {
				feasible = append(feasible, ind)
			}
		}
	}
	return
}

// GetResults returns all ovas and oors from a subset of individuals
//  Output:
//   ovas -- [len(subset)][nova] objective values
//   oors -- [len(subset)][noor] out-of-range values
func (o *Evolver) GetResults(subset []*Individual) (ovas, oors [][]float64) {
	ninds := len(subset)
	ovas = utl.DblsAlloc(ninds, o.C.Nova)
	oors = utl.DblsAlloc(ninds, o.C.Noor)
	for i, ind := range subset {
		for j := 0; j < o.C.Nova; j++ {
			ovas[i][j] = ind.Ovas[j]
		}
		for j := 0; j < o.C.Noor; j++ {
			oors[i][j] = ind.Oors[j]
		}
	}
	return
}

// GetParetoFront returns all feasible individuals on the Pareto front
// Note: input data can be obtained from GetFeasible and GetResults
func (o *Evolver) GetParetoFront(feasible []*Individual, ovas, oors [][]float64) (ovafront, oorfront []*Individual) {
	chk.IntAssert(len(feasible), len(ovas))
	ovaf := utl.ParetoFront(ovas)
	ovafront = make([]*Individual, len(ovaf))
	for i, id := range ovaf {
		ovafront[i] = feasible[id]
	}
	if len(oors) > 0 {
		chk.IntAssert(len(feasible), len(oors))
		oorf := utl.ParetoFront(oors)
		oorfront = make([]*Individual, len(oorf))
		for i, id := range oorf {
			oorfront[i] = feasible[id]
		}
	}
	return
}

// GetFrontOvas collects 2 ova results from Pareto front
//  Input:
//   r and s -- 2 selected objective functions; e.g. r=0 and s=1 for 2D problems
func (o *Evolver) GetFrontOvas(r, s int, front []*Individual) (x, y []float64) {
	x = make([]float64, len(front))
	y = make([]float64, len(front))
	for i, ind := range front {
		x[i] = ind.Ovas[r]
		y[i] = ind.Ovas[s]
	}
	return
}

// GetPopulations returns all populations from all islands
func (o *Evolver) GetPopulations() (pops []Population) {
	pops = make([]Population, o.C.Nisl)
	for i, isl := range o.Islands {
		pops[i] = isl.Pop.GetCopy()
	}
	return
}

// PlotPareto plots pareto front for i and j objectives
func (o *Evolver) PlotPareto(i, j int) {
	feasible := o.GetFeasible()
	ovas, _ := o.GetResults(feasible)
	ovafront, _ := o.GetParetoFront(feasible, ovas, nil)
	f0front, f1front := o.GetFrontOvas(i, j, ovafront)
	f0fin := utl.DblsGetColumn(i, ovas)
	f1fin := utl.DblsGetColumn(j, ovas)
	plt.Plot(f0fin, f1fin, "'r.', clip_on=0, ms=3")
	plt.Plot(f0front, f1front, "'ko',markerfacecolor='none',ms=4, clip_on=0")
}

// GetNfeval returns the number of evaluations
func (o *Evolver) GetNfeval() (nfeval int) {
	for _, isl := range o.Islands {
		nfeval += isl.Nfeval
	}
	return
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

func (o *Evolver) migration(t int) {

	// compute metrics in each island and compute global ova range
	for i, isl := range o.Islands {
		Metrics(isl.OvaMin, isl.OvaMax, isl.Fsizes, isl.Fronts, isl.Pop)
		isl.Pop.SortByRank()
		if i == 0 {
			for j := 0; j < o.C.Nova; j++ {
				o.OvaMin[j] = isl.OvaMin[j]
				o.OvaMax[j] = isl.OvaMax[j]
			}
		} else {
			for j := 0; j < o.C.Nova; j++ {
				o.OvaMin[j] = utl.Min(o.OvaMin[j], isl.OvaMin[j])
				o.OvaMax[j] = utl.Max(o.OvaMax[j], isl.OvaMax[j])
			}
		}
	}

	// loop over pair of islands
	l := o.C.Ninds - o.C.Nimig
	for I := 0; I < o.C.Nisl; I++ {
		Pbest := o.Islands[I].Pop[:o.C.Nimig]
		Pwors := o.Islands[I].Pop[l:]
		for J := I + 1; J < o.C.Nisl; J++ {
			Qbest := o.Islands[J].Pop[:o.C.Nimig]
			Qwors := o.Islands[J].Pop[l:]

			// compute match distances
			for i := 0; i < o.C.Nimig; i++ {
				A := Pbest[i]
				for j := 0; j < o.C.Nimig; j++ {
					B := Qbest[j]
					o.Mdist[i][j] = IndDistance(A, B, o.OvaMin, o.OvaMax)
				}
			}

			// match competitors
			o.Match.SetCostMatrix(o.Mdist)
			o.Match.Run()

			// matches
			for i := 0; i < o.C.Nimig; i++ {
				A := Pbest[i]
				B := Qbest[o.Match.Links[i]]
				if A.Fight(B) {
					A.CopyInto(Qwors[i])
				} else {
					B.CopyInto(Pwors[i])
				}
			}
		}
	}
}

func (o *Evolver) migration_v2(t int) { // TODO: remove this
	Metrics(o.OvaMin, o.OvaMax, o.Fsizes, o.Fronts, o.Competitors)
	nc := o.C.Nisl * o.C.Ninds
	for i := 0; i < nc; i++ {
		A := o.Competitors[i]
		//io.Pforan("cdist = %v\n", A.DistCrowd)
		for j := i + 1; j < nc; j++ {
			B := o.Competitors[j]
			if A.Fight(B) {
				A.Score++
			} else {
				B.Score++
			}
		}
	}
	Population(o.Competitors).SortByScore()
	k := 0
	for i := 0; i < o.C.Nisl; i++ {
		for j := 0; j < o.C.Ninds; j++ {
			o.Competitors[k].CopyInto(o.Islands[i].Pop[j])
			k++
		}
	}
}

func (o Evolver) calc_report(t int) bool {
	return t%o.C.Dtout == 0
}

func (o Evolver) print_legend() {
	io.Pf("\nLEGEND\n")
	io.Pfgrey(" 00 -- generation number (time)\n")
	io.Pfblue(" 00 -- reporting time\n")
	io.Pfyel(" 00 -- migration time\n")
}

func (o Evolver) print_time(time int, report bool) {
	io.Pf(" ")
	if report {
		io.Pfblue("%v", time)
		return
	}
	io.Pfgrey("%v", time)
}
