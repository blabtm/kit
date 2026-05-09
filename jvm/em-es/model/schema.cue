package em_es

#SolverConfig: {
	maxIterations:  >0
	maxEvaluations: >0
	absoluteCost:   >0.0
	relativeCost:   >0.0
}

#AxisSetup: {
	weight: >=0.0 & <=1.0 | *0.0
}

#AxesSetup: {
	x: #AxisSetup
	z: #AxisSetup
}

#CamSetup: {
	axes: #AxesSetup
}

#Config: {
	window: >0 | *3
	solver: #SolverConfig
	cams: {
		"1m1l": #CamSetup
		"1m1r": #CamSetup
		"1m2l": #CamSetup
		"1m2r": #CamSetup
		"2m1l": #CamSetup
		"2m1r": #CamSetup
		"2m2l": #CamSetup
		"2m2r": #CamSetup
		"3m1l": #CamSetup
		"3m1r": #CamSetup
		"3m2l": #CamSetup
		"3m2r": #CamSetup
		"4m1l": #CamSetup
		"4m1r": #CamSetup
		"4m2l": #CamSetup
		"4m2r": #CamSetup
	}
}

#Config
