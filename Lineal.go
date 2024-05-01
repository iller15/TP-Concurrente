package main

import (
	"fmt"
)

// Punto representa un registro de consumo de electricidad.
type Punto struct {
	kWh   float64 // Kilovatios-hora consumidos en un mes.
	costo float64 // Costo en la moneda local.
}

// calcularRegresionLineal encuentra los parámetros m y b en la fórmula y = mx + b para el mejor ajuste de línea.
func calcularRegresionLineal(puntos []Punto) (m, b float64) {
	var sumX, sumY, sumXY, sumXX, promX, promY float64
	n := float64(len(puntos))

	for _, p := range puntos {
		sumX += p.kWh
		sumY += p.costo
		sumXY += p.kWh * p.costo
		sumXX += p.kWh * p.kWh
	}

	promX = sumX / n
	promY = sumY / n
	m = (sumXY - n*promX*promY) / (sumXX - n*promX*promX)
	b = promY - m*promX

	return m, b
}

// predecirCosto estima el costo basado en los kilovatios-hora proporcionados y los parámetros de la regresión lineal.
func predecirCosto(kWh, m, b float64) float64 {
	return m*kWh + b
}

func main() {
	historial := []Punto{
		{kWh: 150, costo: 30},
		{kWh: 200, costo: 40},
		{kWh: 250, costo: 50},
		{kWh: 300, costo: 60},
	}

	m, b := calcularRegresionLineal(historial)
	fmt.Printf("Modelo de predicción: Costo = %.2fkWh + %.2f\n", m, b)

	// Imagina que quieres estimar el costo para un consumo de 350 kWh.
	consumoEstimado := 350.0
	costoEstimado := predecirCosto(consumoEstimado, m, b)
	fmt.Printf("El costo estimado para un consumo de %.2f kWh es de $ %.2f\n", consumoEstimado, costoEstimado)
}
