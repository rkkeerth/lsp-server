package main

import "fmt"

// Calculator provides basic arithmetic operations
type Calculator struct {
	result float64
}

// NewCalculator creates a new calculator instance
func NewCalculator() *Calculator {
	return &Calculator{result: 0}
}

// Add adds a number to the current result
func (c *Calculator) Add(n float64) float64 {
	c.result += n
	return c.result
}

// Subtract subtracts a number from the current result
func (c *Calculator) Subtract(n float64) float64 {
	c.result -= n
	return c.result
}

// Multiply multiplies the current result by a number
func (c *Calculator) Multiply(n float64) float64 {
	c.result *= n
	return c.result
}

// Divide divides the current result by a number
func (c *Calculator) Divide(n float64) float64 {
	if n != 0 {
		c.result /= n
	}
	return c.result
}

// Reset resets the calculator result to zero
func (c *Calculator) Reset() {
	c.result = 0
}

// GetResult returns the current result
func (c *Calculator) GetResult() float64 {
	return c.result
}

const (
	// MaxValue represents the maximum allowed value
	MaxValue = 1000000
	// MinValue represents the minimum allowed value
	MinValue = -1000000
)

var globalCalculator *Calculator

func main() {
	// TODO: Add more operations
	calc := NewCalculator()
	
	calc.Add(10)
	calc.Multiply(2)
	calc.Subtract(5)
	
	// FIXME: Handle division by zero properly
	result := calc.GetResult()
	
	fmt.Printf("Result: %.2f\n", result)
}
