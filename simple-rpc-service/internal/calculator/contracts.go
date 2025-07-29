package calculator

type CalculatorContract interface {
	Add(a int32, b int32) int32
	Substract(a int32, b int32) int32
	Multiply(a int32, b int32) int32
	Divide(a int32, b int32) int32
}

type CalculatorService struct {
	CalculatorContract
}

func (m *CalculatorService) Add(a int32, b int32) int32 {
	return a + b
}
func (m *CalculatorService) Substract(a int32, b int32) int32 {
	return a - b
}
func (m *CalculatorService) Multiply(a int32, b int32) int32 {
	return a * b
}
func (m *CalculatorService) Divide(a int32, b int32) int32 {
	return a / b
}
