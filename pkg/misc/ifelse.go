package misc

// IfElse 模拟三元操作符
func IfElse(condition bool, positiveVal, negativeVal interface{}) interface{} {
	if condition {
		return positiveVal
	}

	return negativeVal
}
