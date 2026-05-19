package em_es

#SolverConfig: {
	maxIterations:  >0 | *500
	maxEvaluations: >0 | *500
	absoluteCost:   >0.0 | *1e-7
	relativeCost:   >0.0 | *1e-10
}

#AxisSetup: {
	weight: >=0.0 & <=1.0 | *0.0
}

#AxesSetup: {
	x!: #AxisSetup
	z!: #AxisSetup
}

#CamSetup: {
	axes!: #AxesSetup
}

#Config: {
	window: >0 | *3
	solver: #SolverConfig
	cams!: {
		[string]: #CamSetup
	}
}
