package mf

#Material: {
	relativePermeability!: >0.0
}

#Coil: {
	currentAmp!:    >0.0
	numberOfTurns!: >0
}

#Config: {
	yoke!: #Material
	coils: {
		nbti!: #Coil
		nbsn!: #Coil
		comp!: #Coil
	}
}
