// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"encoding/json"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// ConfParams is an auxiliary structure to hold configuration parameters for setting the GA up
type ConfParams struct {

	// initialisation
	Nisl   int     // number of islands
	Ninds  int     // number of individuals: population size
	Nbases int     // number of bases in chromosome
	Grid   bool    // generate individuals based on grid
	Noise  float64 // apply noise when generate based on grid (if Noise > 0)

	// time control
	Tf    int // number of generations
	Dtout int // increment of time for output
	Dtmig int // increment of time for migration
	Dtreg int // increment of time for regeneration

	// regeneration
	RegIni    bool    // regenerate initial population
	RegTol    float64 // tolerance for ρ to activate regeneration
	RegNmax   int     // max number of generations. use -1 for unlimited number of generations
	RegBest   bool    // enforce that regeneration is always based on based individual, regardless the population is homogeneous or not
	RegPct    float64 // percentage of individuals to be regenerated
	RegMmin   float64 // multiplier to decrease reference value; e.g. 0.1
	RegMmax   float64 // multiplier to increase reference value; e.g. 10.0
	UseStdDev bool    // use standard deviation (σ) instead of average deviation in Stat

	// selection and reproduction
	Pc    float64 // probability of crossover
	Pm    float64 // probability of mutation
	Elite bool    // use elitism
	Rws   bool    // use Roulette-Wheel selection method
	Rnk   bool    // ranking
	RnkSp float64 // selective pressure for ranking

	// diversity
	StatOorSkip bool // skip oor individuals from statistics

	// output
	Json      bool   // output results as .json files; not tables
	DirOut    string // directory to save output files. "" means "/tmp/goga"
	FnKey     string // filename key for output files. "" means no output files
	DoPlot    bool   // plot results
	ShowBases bool   // show also bases when printing results (if any)

	// crossover
	CxNcuts   map[string]int         // crossover number of cuts for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxCuts    map[string][]int       // crossover specific cuts for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxProbs   map[string]float64     // crossover probabilities for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxFuncs   map[string]interface{} // crossover functions for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxIntFunc CxIntFunc_t            // crossover function
	CxFltFunc CxFltFunc_t            // crossover function
	CxStrFunc CxStrFunc_t            // crossover function
	CxKeyFunc CxKeyFunc_t            // crossover function
	CxBytFunc CxBytFunc_t            // crossover function
	CxFunFunc CxFunFunc_t            // crossover function

	// mutation
	MtNchanges map[string]int         // mutation number of changes for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	MtProbs    map[string]float64     // mutation probabilities for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	MtExtra    map[string]interface{} // mutation extra parameters for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	MtIntFunc  MtIntFunc_t            // mutation function
	MtFltFunc  MtFltFunc_t            // mutation function
	MtStrFunc  MtStrFunc_t            // mutation function
	MtKeyFunc  MtKeyFunc_t            // mutation function
	MtBytFunc  MtBytFunc_t            // mutation function
	MtFunFunc  MtFunFunc_t            // mutation function
}

// SetDefault sets default parameters
func (o *ConfParams) SetDefault() {

	// initialisation
	o.Nisl = 1
	o.Ninds = 20
	o.Nbases = 10
	o.Grid = true
	o.Noise = 0.2

	// time control
	o.Tf = 100
	o.Dtout = 10
	o.Dtmig = 30
	o.Dtreg = 60

	// regeneration
	o.RegIni = false
	o.RegTol = 1e-3
	o.RegNmax = -1
	o.RegBest = false
	o.RegPct = 0.3
	o.RegMmin = 0.1
	o.RegMmax = 10.0
	o.UseStdDev = false

	// selection and reproduction
	o.Pc = 0.8
	o.Pm = 0.01
	o.Elite = true
	o.Rws = false
	o.Rnk = true
	o.RnkSp = 1.2

	// diversity
	o.StatOorSkip = false

	// output
	o.Json = false
	o.DirOut = "/tmp/goga"
	o.FnKey = ""
	o.DoPlot = false
	o.ShowBases = true
}

// NewConfParams returns a new ConfParams structure, with default values set
func NewConfParams() *ConfParams {
	var o ConfParams
	o.SetDefault()
	return &o
}

// ReadConfParams reads configuration parameters from JSON file
func ReadConfParams(filenamepath string) *ConfParams {

	// new params
	var o ConfParams
	o.SetDefault()

	// read file
	b, err := io.ReadFile(filenamepath)
	if err != nil {
		chk.Panic("cannot read parameters file %q", filenamepath)
	}

	// decode
	err = json.Unmarshal(b, &o)
	if err != nil {
		chk.Panic("cannot unmarshal parameters file %q", filenamepath)
	}

	// set probabilities
	pc, pm := o.Pc, o.Pm
	o.CxProbs = map[string]float64{"int": pc, "flt": pc, "str": pc, "key": pc, "byt": pc, "fun": pc}
	o.MtProbs = map[string]float64{"int": pm, "flt": pm, "str": pm, "key": pm, "byt": pm, "fun": pm}

	// results
	return &o
}