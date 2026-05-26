package mf

import v2k "github.com/blabtm/v2k/model"
import "github.com/blabtm/v2k/model/material"

#Material: material.#Material & {
	relativePermeability: float | *1.0
}

#Yoke: {
	material: #Material
}

#Coil: {
	material:       #Material
	currentAmp!:    >0.0
	numberOfTurns!: >0
}

#Config: v2k.#ServiceConfig & {
	yoke: material.#Material
	coils: {
		nbti!: #Coil
		nbsn!: #Coil
		comp!: #Coil
	}
}
