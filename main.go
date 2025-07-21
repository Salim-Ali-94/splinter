package main


func main() {

	config := map[string]interface{}{

		"variable": "s",

		"numerator": map[int64]float64{

			-1: -0.83,
			0: 3.72,
			2: -4.09,

		},

		"denominator": map[int64]float64{

			1: -1.23,
			2: 16.24,

		},

	}

	polynomial := constructPolynomial(config)
	polynomial.representation()

}
