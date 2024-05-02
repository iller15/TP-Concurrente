package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Punto representa un registro de consumo de electricidad.
type Punto struct {
	kWh   float64 // Kilovatios-hora consumidos en un mes.
	costo float64 // Costo en la moneda local.
}

// calcularRegresionLineal encuentra los parámetros m y b en la fórmula y = mx + b para el mejor ajuste de línea.
func calcularRegresionLineal(puntos []Punto) (m, b float64) {
	var sumX, sumY, sumXY, sumXX float64
	n := float64(len(puntos))

	for _, p := range puntos {
		sumX += p.kWh
		sumY += p.costo
		sumXY += p.kWh * p.costo
		sumXX += p.kWh * p.kWh
	}

	m = (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
	b = (sumY - m*sumX) / n

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

	// Realizar pruebas múltiples y calcular la media recortada
	numTests := 1000
	runtimes := make([]time.Duration, numTests)

	for j := 0; j < numTests; j++ {
		startTime := time.Now()
		for i := 0; i < 1000000; i++ {
			consumoEstimado := float64(350 + i)
			predecirCosto(consumoEstimado, m, b)
		}
		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		runtimes[j] = elapsedTime
	}

	// Ordenar los tiempos de ejecución
	sort.Slice(runtimes, func(i, j int) bool { return runtimes[i] < runtimes[j] })

	// Calcular la media recortada de los tiempos de ejecución
	trimmedRuntime := calculateTrimmedMean(runtimes)
	fmt.Printf("Media recortada de los tiempos de ejecución: %s\n", trimmedRuntime)
}

func calculateTrimmedMean(data []time.Duration) time.Duration {
	numTrim := 50
	if len(data) <= 2*numTrim {
		sum := time.Duration(0)
		for _, value := range data {
			sum += value
		}
		return sum / time.Duration(len(data))
	}
	trimmedData := data[numTrim : len(data)-numTrim]
	sum := time.Duration(0)
	for _, value := range trimmedData {
		sum += value
	}
	return sum / time.Duration(len(trimmedData))
}
