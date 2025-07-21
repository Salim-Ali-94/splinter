package main

import ( "fmt"
		 // "math/cmplx"
		 "math"
		 "slices" )


func constructPolynomial(parameters ...map[string]interface{}) Polynomial {

	polynomial := Polynomial{}
    top, bottom, variable := extract(parameters...)
    numerator := unpackExpression(top, variable)
    denominator := unpackExpression(bottom, variable, true)
	polynomial.Numerator = numerator
	polynomial.Denominator = denominator
	polynomial.Numerator.build()
	polynomial.Denominator.build()
	polynomial.build()
	return polynomial

}

func dotProduct(p []Term, q []Term) map[int64]float64 {

	expression := map[int64]float64{}

	for _, P := range p {

		for _, Q := range q {

			coefficient := P.Coefficient*Q.Coefficient
			exponent := P.Exponent + Q.Exponent

			_, exists := expression[exponent]

			if exists {

				expression[exponent] += coeffient

			} else {

				expression[exponent] = coefficient

			}

		}

	}

	return expression

}

func findConstant(terms []Term) ([]Term, Term) {

	removedTerm := Term{}
	modifiedExpression := terms

	for index, term := range terms {

		if (term.Exponent == int64(0)) {

			removedTerm = term
			modifiedExpression = slices.Delete(modifiedExpression, index, index + 1)
			break

		}

	}

	return modifiedExpression, removedTerm

}

func unpackExpression(expressionLUT map[int64]float64, variable string, denominatorFlag ...bool) Expression {

	flag := false
	expression := Expression{}

	if (len(expressionLUT) == 0) {

		coefficient := 0.0

		if (len(denominatorFlag) > 0) {

			flag = denominatorFlag[0]

		}

		if flag {

			coefficient = 1.0

		}

		lut := map[int64]float64{

			0: coefficient,

		}

		expression.transform(lut, variable)

	} else {

		expression.transform(expressionLUT, variable)

	}

	expression.expand()
	return expression

}

func extract(parameters ...map[string]interface{}) (map[int64]float64, map[int64]float64, string) {

	variable := "s"
	numerator := map[int64]float64{}
	denominator := map[int64]float64{}

    if (len(parameters) > 0) {

    	lut := parameters[0]
    	x, exists := lut["variable"]

    	if exists {

    		variable = x.(string)

    	}

    	top, exists := lut["numerator"]

    	if exists {

    		numerator = top.(map[int64]float64)

    	}

    	bottom, exists := lut["denominator"]

    	if exists {

    		denominator = bottom.(map[int64]float64)

    	}

    }

    return numerator, denominator, variable

}
