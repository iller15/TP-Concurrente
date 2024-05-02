package main

import (
	"fmt"
	"sync"
	"time"
)

type Punto struct {
	kWh   float64
	costo float64
}

func calcularRegresionLineal(puntos []Punto) (m, b float64) {
	var sumX, sumY, sumXY, sumXX float64
	n := float64(len(puntos))

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(4)

	go func() {
		defer wg.Done()
		subSumX := 0.0
		for _, p := range puntos {
			subSumX += p.kWh
		}
		mu.Lock()
		sumX += subSumX
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		subSumY := 0.0
		for _, p := range puntos {
			subSumY += p.costo
		}
		mu.Lock()
		sumY += subSumY
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		subSumXY := 0.0
		for _, p := range puntos {
			subSumXY += p.kWh * p.costo
		}
		mu.Lock()
		sumXY += subSumXY
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		subSumXX := 0.0
		for _, p := range puntos {
			subSumXX += p.kWh * p.kWh
		}
		mu.Lock()
		sumXX += subSumXX
		mu.Unlock()
	}()

	wg.Wait()

	promX := sumX / n
	promY := sumY / n
	m = (sumXY - n*promX*promY) / (sumXX - n*promX*promX)
	b = promY - m*promX

	return m, b
}

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

	start := time.Now()

	for j := 0; j < 1000; j++ {
		for i := 0; i < 1000000; i++ {
			consumoEstimado := float64(350 + i)
			predecirCosto(consumoEstimado, m, b)
		}
		
	}

	duration := time.Since(start)
	fmt.Printf("Tiempo de ejecución: %s\n", duration)
}
