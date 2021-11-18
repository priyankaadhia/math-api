package main

import (
	"log"
	"math"
	"sort"
)

type Calculator interface {
	Min(quantifier int) []float64
	Max(quantifier int) []float64
	Average() float64
	Median() float64
	Percentile(quantifier int) float64
}

type Calculation struct {
	Numbers []float64
}

func validateQuantifier(quantifier *int, numbersLength int) {
	if *quantifier > numbersLength {
		log.Printf("Quantifier can't be greater than length of numbers array")
		*quantifier = numbersLength
	}
}

func (c Calculation) Min(quantifier int) []float64 {
	validateQuantifier(&quantifier, len(c.Numbers))
	sort.Float64s(c.Numbers)
	return c.Numbers[:quantifier]
}

func (c Calculation) Max(quantifier int) []float64 {
	validateQuantifier(&quantifier, len(c.Numbers))
	sort.Float64s(c.Numbers)
	return c.Numbers[len(c.Numbers)-quantifier:]
}

func (c Calculation) Average() float64 {
	var sum float64 = 0
	for _, v := range c.Numbers {
		sum += v
	}
	return sum / float64(len(c.Numbers))
}

func (c Calculation) Median() float64 {
	sort.Float64s(c.Numbers)
	// handle even length of numbers array
	if len(c.Numbers)%2 == 0 {
		return (c.Numbers[len(c.Numbers)/2-1] + c.Numbers[len(c.Numbers)/2]) / 2
	}

	return c.Numbers[len(c.Numbers)/2]
}

func (c Calculation) Percentile(quantifier int) float64 {
	sort.Float64s(c.Numbers)
	index := int(math.Round((float64(quantifier) / 100) * float64(len(c.Numbers))))
	if index < 0 || index > len(c.Numbers) {
		log.Printf("Quantifier must be between 0 and 100")
		return 0
	}
	if index == 0 {
		return c.Numbers[index]
	}
	return c.Numbers[index-1]
}
