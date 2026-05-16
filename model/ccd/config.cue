package ccd

#AxisConfig: {
	beta!:       >=0.0 | *0.0
	dispersion!: >=0.0 | *0.0
}

#AxesConfig: {
	x!: #AxisConfig
	z!: #AxisConfig
}

#CamConfig: {
	axes!: #AxesConfig
}

#Config: {
	cams!: {
		[string]: #CamConfig
	}
}
