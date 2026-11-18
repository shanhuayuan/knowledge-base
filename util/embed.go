package util

func Float64MatrixToFloat32(matrix [][]float64) [][]float32 {
	result := make([][]float32, len(matrix))

	for i, row := range matrix {
		result[i] = make([]float32, len(row))

		for j, value := range row {
			result[i][j] = float32(value)
		}
	}

	return result
}
