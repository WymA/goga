// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/gofem/fem"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func set_enabled_disabled(dom *fem.Domain, enabled []int) {
	cells := dom.Msh.Cells
	for cid, ena := range enabled {
		cells[cid].Disabled = false
		if ena == 0 {
			cells[cid].Disabled = true
		}
	}
}

func main() {

	// settings: upward movement
	/*
		enabled := []int{1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1}
		areas := []float64{34.00821693709039, 0.09, 9.968259379257471, 35, 15.831579773931853, 8.654957754397607, 11.965649280046627, 19.413683184371774, 7.6546849806620525, 5.387748841496445, 35, 29.504717529844843, 26.86909134752426, 35, 20.10785632243804, 3.446115518045177, 0.09, 35, 26.216339636590078, 9.542311851327804, 35, 0.09, 0.09, 28.407557441773008, 29.933108719044267, 10.922581748461933, 1.8067072461717366, 0.09, 14.7804274343423, 0.09, 11.730811122600027, 35, 35, 0.09, 35, 35}
			weight  = 10169.99460559597
			umax    = 0.03190476190476189
			smax    = 20.518951772912644
	*/

	// settings: upward movement
	//enabled := []int{1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1}
	//areas := []float64{29.673971565565495, 9.804876568305883, 12.722191736896143, 31.128370429558302, 12.28498763291402, 0.09, 1.0976024620675062, 29.054175458117097, 0.09, 12.074834078336714, 0.9518626611813701, 0.804189446111032, 9.620926416457582, 23.064311560926264, 4.903260570239974, 0.09, 14.345604382360431, 31.10942565747464, 0.09, 7.790214820299472, 7.491266591459749, 21.567320209602265, 4.905574834666627, 0.09, 9.065113525395782, 23.84052973418943, 6.7235867554969975, 3.6046158266920836, 24.589638797955896, 0.09, 31.780396612723077, 23.409598016209728, 3.50718429240112, 15.956651688597585, 35, 12.255743491145445}

	//enabled := []int{1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1}
	//areas := []float64{21.480878134095338, 0.09, 14.694965198034584, 24.824367367224532, 6.729854812525405, 0.09, 0.09, 18.170644951273943, 27.43068988046519, 16.340137823665955, 0.09, 35, 33.257655346869484, 0.09, 10.739467844959764, 1.2284583619296825, 0.09, 13.836693890116672, 3.223634459640667, 31.609509632805768, 2.9580912890580233, 0.09, 11.66346650529349, 11.839368679149583, 8.037665593492571, 18.4772618019285, 6.0722754499289335, 8.299339699920758, 18.092667282860184, 0.09, 3.95809930082411, 35, 24.98891088932943, 0.09, 20.001590440636104, 0.4232030075463411}

	//enabled := []int{1, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1}
	//areas := []float64{1.2035515116115514, 0.811277071473871, 32.37830152936991, 11.012589123896603, 29.419388200704844, 11.517528463414674, 26.842094480154707, 0.09, 7.545801867132738, 22.246892098984826, 33.64813536709853, 35, 18.79453561647245, 19.72091117582699, 24.417433685262015, 17.139485224780174, 14.64143052284774, 6.017622261768879, 18.627730008706013, 6.034380625351308, 15.909160991008125, 3.010643800045916, 35, 1.7855841010723912, 23.882989565364397, 4.179630598025799, 8.060183267136836, 27.61994718378331, 26.443620790772826, 35, 0.9889261275628931, 0.09, 22.110211729649148, 31.153765658657143, 19.868907703384732, 23.523896513200622}

	enabled := []int{1, 1, 1, 1, 0, 0, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1}
	areas := []float64{22.757099750148996, 28.986754816914374, 6.957451713927281, 8.528672093936208, 9.26287758087651, 2.4766961118094857, 19.144883540012557, 6.382281135196173, 10.365771948733226, 2.0524134188513186, 3.3776112861797856, 22.444923164834712, 0.09, 0.09, 0.09, 8.68418380122592, 0.09, 0.3395486752846164, 4.356831930853984, 12.792965016026955, 20.651430212430448, 4.881368992183173, 17.009172478115723, 14.806321101924194, 9.298936701386527, 5.820319254311902, 11.792969696093445, 12.323517103405779, 1.4343013440743113, 2.3392600723999366, 0.09, 11.352516128577138, 11.223982208350751, 23.98665707191376, 0.09, 0.09}

	// start simulation
	processing := fem.NewFEM("cantilever.sim", "", true, true, false, false, true, 0)

	// set enabled/disabled
	dom := processing.Domains[0]
	if true {
		set_enabled_disabled(dom, enabled)
	}

	// set stage
	err := processing.SetStage(0)
	if err != nil {
		io.PfRed("SetStage failed:\n%v", err)
		return
	}

	// set areas
	lwds := make(map[int]float64)
	if true {
		weight := 0.0
		for _, elem := range dom.Elems {
			ele := elem.(*fem.ElastRod)
			cid := ele.Cell.Id
			ele.Mdl.A = areas[cid]
			ele.Recompute(false)
			weight += ele.Mdl.Rho * ele.Mdl.A * ele.L
			lwds[cid] = 0.1 + ele.Mdl.A/10.0
		}
		io.Pforan("weight = %v\n", weight)
	}

	// mobility
	n := len(dom.Nodes)
	m := len(dom.Elems)
	d := len(dom.EssenBcs.Bcs)
	F := 2*n - m - d
	io.Pforan("mobility = %v\n", F)

	// plot
	msh := dom.Msh
	msh.Draw2d(true, lwds)
	plt.SaveD("/tmp/goga", "rods.eps")

	// run FE analysis
	err = processing.SolveOneStage(0, true)
	if err != nil {
		io.PfRed("Run failed:\n%v", err)
		return
	}

	// post-processing
	vid := msh.VertTag2verts[-4][0].Id
	nod := dom.Vid2node[vid]
	eqy := nod.GetEq("uy")
	uy := dom.Sol.Y[eqy]
	io.PfYel("%2d : uy = %g\n", vid, uy)
	smax := 0.0
	for _, elem := range dom.Elems {
		ele := elem.(*fem.ElastRod)
		sig := math.Abs(ele.CalcSig(dom.Sol))
		if sig > smax {
			smax = sig
		}
	}
	io.Pfred("smax = %v\n", smax)
}

//mats := dom.Sim.MatModels
//reg := dom.Reg
//dat := reg.Etag2data(-1)
//matname := dat.Mat
//mat := mats.Get(matname)
//io.Pforan("mat = %v\n", mat)
//io.Pforan("dat = %v\n", dat)

/*
	enabled = [1 0 1 1 0 1 1 0 1 0 1 0 1 0 1 1 0 1 1 1 1 1 1 1 1 0 1 0 1 1 1 1 0 0 0 1]
	areas   = [35.000000 24.038714 15.111653 14.992768 35.000000 35.000000 6.004577 35.000000 24.874932 35.000000 12.921652 16.298292 18.844584 35.000000 14.289417 31.588943 0.090000 22.623541 35.000000 35.000000 18.884884 20.930816 11.169120 17.494270 20.525592 5.581622 35.000000 7.399986 34.363928 5.274137 30.251305 28.175011 35.000000 0.664373 22.616175 10.328802]
	weight  = 11598.0725140786
	umax    = 0.09928756240366518
	smax    = 24.499761348024837
	errU    = 0
	errS    = 0
*/
