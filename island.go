// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// ObjFunc_t defines the template for the objective function
type ObjFunc_t func(ind *Individual, time int, best *Individual)

// Island holds one population and performs the reproduction operation
type Island struct {

	// selection/reproduction
	UseRanking  bool    // use ranking for selection process
	RnkPressure float64 // ranking pressure
	Roulette    bool    // use roulette wheel selection; otherwise use stochastic-universal-sampling selection
	Elitism     bool    // perform elitism: keep at least one best individual from previous generation

	// crossover
	CxNcuts map[string]int         // crossover number of cuts for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxCuts  map[string][]int       // crossover specific cuts for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxProbs map[string]float64     // crossover probabilities for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxFuncs map[string]interface{} // crossover functions for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag

	// mutation
	MtNchanges map[string]int         // mutation number of changes for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	MtProbs    map[string]float64     // mutation probabilities for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	MtExtra    map[string]interface{} // mutation extra parameters for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	MtFuncs    map[string]interface{} // mutation functions for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag

	// population
	Pop     Population // pointer to current population
	BkpPop  Population // backup population
	ObjFunc ObjFunc_t  // objective function

	// auxiliary internal data
	fitsrnk []float64 // all fitness values computed by ranking
	fitness []float64 // all fitness values
	prob    []float64 // probabilities
	cumprob []float64 // cumulated probabilities
	selinds []int     // indices of selected individuals
	A, B    []int     // indices of selected parents
}

// NewIsland allocates a new island but with a give population already allocated
// Input:
//  pop   -- the population
//  ofunc -- objective function
func NewIsland(pop Population, ofunc ObjFunc_t) (o *Island) {

	// check
	ninds := len(pop)
	if ninds%2 != 0 {
		chk.Panic("size of population must be even")
	}

	// allocate
	o = new(Island)
	o.Pop = pop
	o.BkpPop = pop.GetCopy()
	o.ObjFunc = ofunc

	// set default control values
	o.UseRanking = true
	o.RnkPressure = 1.2
	o.Elitism = true

	// compute objective values
	for _, ind := range o.Pop {
		o.ObjFunc(ind, 0, nil)
	}

	// sort
	o.Pop.Sort()

	// auxiliary data
	o.fitsrnk = make([]float64, ninds)
	o.fitness = make([]float64, ninds)
	o.prob = make([]float64, ninds)
	o.cumprob = make([]float64, ninds)
	o.selinds = make([]int, ninds)
	o.A = make([]int, ninds/2)
	o.B = make([]int, ninds/2)
	return
}

// SelectAndReprod performs the selection and reproduction processes
func (o *Island) SelectAndReprod(time int) {

	// fitness and probabilities
	ninds := len(o.Pop)
	sumfit := 0.0
	if o.UseRanking {
		sp := o.RnkPressure
		if sp < 1.0 || sp > 2.0 {
			sp = 1.2
		}
		for i := 0; i < ninds; i++ {
			o.fitness[i] = 2.0 - sp + 2.0*(sp-1.0)*float64(ninds-i-1)/float64(ninds-1)
			sumfit += o.fitness[i]
		}
	} else {
		ovmin, ovmax := o.Pop[0].ObjValue, o.Pop[0].ObjValue
		for _, ind := range o.Pop {
			ovmin = min(ovmin, ind.ObjValue)
			ovmax = max(ovmax, ind.ObjValue)
		}
		if math.Abs(ovmax-ovmin) < 1e-14 {
			for i := 0; i < ninds; i++ {
				o.fitness[i] = float64(i) / float64(ninds-1)
				sumfit += o.fitness[i]
			}
		} else {
			for i, ind := range o.Pop {
				o.fitness[i] = (ovmax - ind.ObjValue) / (ovmax - ovmin)
				sumfit += o.fitness[i]
			}
		}
	}
	for i := 0; i < ninds; i++ {
		o.prob[i] = o.fitness[i] / sumfit
	}
	CumSum(o.cumprob, o.prob)

	// selection
	if o.Roulette {
		RouletteSelect(o.selinds, o.cumprob, nil)
	} else {
		SUSselect(o.selinds, o.cumprob, -1)
	}
	FilterPairs(o.A, o.B, o.selinds)

	// reproduction
	h := ninds / 2
	for i := 0; i < ninds/2; i++ {
		Crossover(o.BkpPop[i], o.BkpPop[h+i], o.Pop[o.A[i]], o.Pop[o.B[i]], o.CxNcuts, o.CxCuts, o.CxProbs, o.CxFuncs)
		Mutation(o.BkpPop[i], o.MtNchanges, o.MtProbs, o.MtExtra, o.MtFuncs)
		Mutation(o.BkpPop[h+i], o.MtNchanges, o.MtProbs, o.MtExtra, o.MtFuncs)
	}

	// compute objective values
	for _, ind := range o.BkpPop {
		o.ObjFunc(ind, 0, nil)
	}

	// sort
	o.BkpPop.Sort()

	// elitism
	if o.Elitism {
		if o.Pop[0].ObjValue < o.BkpPop[0].ObjValue {
			o.Pop[0].CopyInto(o.BkpPop[ninds-1])
			o.BkpPop.Sort()
		}
	}

	// swap populations
	o.Pop, o.BkpPop = o.BkpPop, o.Pop
}