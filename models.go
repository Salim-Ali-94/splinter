package main

import ( "strings"
		 "fmt"
		 "math"
		 "strconv" )


type Polynomial struct {

	Numerator 	   	 Expression
	Denominator    	 Expression
	TransferFunction string

}

type Expression struct {

	Terms 	  	   []Term
	Expansion 	   map[int64]float64
	Reduction 	   map[int64]float64
	Representation string

}

type Term struct {

	Coefficient float64
	Exponent	int64
	Variable 	string

}

func (e *Expression) sort() {

	orderedList, dc := findConstant(e.Terms)
	cycles := len(orderedList)

	if (cycles%2 == 0) {

		cycles = cycles / 2

	} else {

		cycles = (cycles - 1) / 2

	}

	for address := 0; address < cycles; address++ {

		front := address
		back := len(orderedList[front:]) - address - 1
		max := orderedList[front:back + 1][0]
		last := len(orderedList[front:back + 1]) - 1
		min := orderedList[front:back + 1][last]
		position := back

		for index, term := range orderedList[address:position + 1] {

			if (term.Exponent >= max.Exponent) {

				max = term
				front = index

			}

			if (term.Exponent <= min.Exponent) {

				min = term
				back = index

			}

		}

		orderedList[front] = orderedList[address]
		orderedList[back] = orderedList[position]
		orderedList[address] = max
		orderedList[position] = min

	}

	if (len(orderedList) < len(e.Terms)) {

		orderedList = append(orderedList, dc)

	}

	e.Terms = orderedList

}

func (p *Polynomial) multiply(q Polynomial) {

	if ((len(p.Numerator.Terms) > 0) &&
		(len(p.Denominator.Terms) > 0) &&
		(len(q.Numerator.Terms) > 0) &&
		(len(q.Denominator.Terms) > 0)) {

		s := p.Denominator.Terms[0].Variable
		expression := dotProduct(p.Numerator.Terms, q.Numerator.Terms)
		p.Numerator.transform(expression, s)
		p.Numerator.expand()
		p.Numerator.build()

		expression = dotProduct(p.Denominator.Terms, q.Denominator.Terms)
		p.Denominator.transform(expression, s)
		p.Denominator.expand()
		p.Denominator.build()
		p.build()

	}

}

func (e *Expression) build() {

	terms := ""
	e.sort()

	for index, term := range e.Terms {

		coefficient := strconv.FormatFloat(term.Coefficient, 'f', -1, 64)
		exponent := strconv.FormatInt(term.Exponent, 10)

		if (coefficient != "0") {

			if (((exponent == "0") && (coefficient == "1")) ||
				((exponent == "0") && (coefficient == "-1") && (index == 0)) ||
				((exponent == "0") && (coefficient != "-1") && (coefficient != "1")) ||
				((exponent == "0") && (term.Coefficient < 0) && (index == 0))) {

				terms += coefficient

			} else if (((exponent == "0") && (coefficient == "-1") && (index > 0)) ||
					   ((exponent == "0") && (term.Coefficient < 0) && (index > 0))) {

				terms += strings.ReplaceAll(coefficient, "-", "")

			} else if (((exponent == "1") && (coefficient == "1")) ||
					   ((exponent == "1") && (coefficient == "-1") && (index > 0))) {

				terms += term.Variable

			} else if ((exponent == "1") && (coefficient == "-1") && (index == 0)) {

				terms += "-" + term.Variable

			} else if (((exponent == "-1") && (coefficient == "1")) ||
					   ((exponent == "-1") && (coefficient == "-1") && (index > 0)) ||
					   ((exponent != "-1") && (exponent != "1") && (exponent != "0") && (coefficient == "1")) ||
					   ((exponent != "-1") && (exponent != "1") && (exponent != "0") && (coefficient == "-1") && (index > 0))) {

				exponent = getSuperscript(exponent)
				terms += term.Variable + exponent

			} else if (((exponent == "-1") && (coefficient == "-1") && (index == 0)) ||
					   ((exponent != "-1") && (exponent != "1") && (exponent != "0") && (coefficient == "-1") && (index == 0))) {

				exponent = getSuperscript(exponent)
				terms += "-" + term.Variable + exponent

			} else if ((exponent == "1") && (term.Coefficient < 0) && (index > 0)) {

				terms += strings.ReplaceAll(coefficient, "-", "") + term.Variable

			} else if (((exponent == "1") && (coefficient != "-1") && (coefficient != "1")) ||
					   ((exponent == "1") && (term.Coefficient < 0) && (index == 0))) {

				terms += coefficient + term.Variable

			} else if (((exponent == "-1") && (term.Coefficient < 0) && (index > 0)) ||
					   ((exponent != "-1") && (exponent != "1") && (exponent != "0") && (term.Coefficient < 0) && (index > 0))) {

				exponent = getSuperscript(exponent)
				terms += strings.ReplaceAll(coefficient, "-", "") + term.Variable + exponent

			} else {

				exponent = getSuperscript(exponent)
				terms += coefficient + term.Variable + exponent

			}

			if (index < len(e.Terms) - 1) {

				if (e.Terms[index + 1].Coefficient < 0) {

					terms += " - "

				} else {

					terms += " + "

				}

			}

		}

	}

	e.Representation = terms

}

func (p *Polynomial) build() {

	p.TransferFunction = p.Numerator.Representation
	
	constant := (p.Denominator.Terms[0].Exponent == 0) &&
				(len(p.Denominator.Terms) == 1)

	if !constant {

		lengthNumerator := len(p.Numerator.Representation)
		lengthDenominator := len(p.Denominator.Representation)
		lengthLine := max(lengthNumerator, lengthDenominator)
		line := ""

		for index := 0; index < lengthLine + 4; index++ {

			line += "-"

		}

		centered := p.Numerator.Representation

		if (lengthLine == lengthNumerator) {

			centered = p.Denominator.Representation

		}

		delta := len(line) - len(centered)
		empty := math.Floor(float64(delta / 2))

		for index := 0.0; index < empty; index++ {

			centered = " " + centered

		}

		if (lengthLine == lengthNumerator) {

			p.TransferFunction = fmt.Sprintf("  %s\n%s\n%s", p.Numerator.Representation, line, centered)

		} else {

			p.TransferFunction = fmt.Sprintf("%s\n%s\n  %s", centered, line, p.Denominator.Representation)

		}

	}

}

func (p *Polynomial) representation(label ...string) {

	function := "H"

	if (len(label) > 0) {

		function = label[0]

	}

	variable := p.Numerator.Terms[0].Variable
	transferFunction := fmt.Sprintf("\n%s(%s) = ", function, variable)
	prefix := ""
	suffix := "-\n"

	for index := 0; index < len(transferFunction) + 1; index++ {

		prefix += " "
		suffix += " "

	}

	formatted := strings.Replace(p.TransferFunction, "\n", transferFunction, 1)
	formatted = strings.Replace(formatted, "-\n", suffix, 1)
	formatted = prefix + formatted
	fmt.Printf("\n%s\n", formatted)

}

func (e *Expression) representation(label ...string) {

	function := "H"
	variable := e.Terms[0].Variable

	if (len(label) > 0) {

		function = label[0]

	}

	fmt.Printf("\n%s(%s) = %s\n", function, variable, e.Representation)

}

func (e *Expression) expand() {

	if (len(e.Terms) > 0) {

		maxPower := e.Terms[0].Exponent
		minPower := e.Terms[0].Exponent

		for _, term := range e.Terms {

			if (term.Exponent > maxPower) {

				maxPower = term.Exponent

			}

			if (term.Exponent < minPower) {

				minPower = term.Exponent

			}

		}

		for index := minPower; index <= maxPower; index++ {

			_, exists := e.Expansion[index]

			if !exists {

				e.Expansion[index] = 0.0

			}

		}

	}

}

func (e *Expression) transform(expression map[int64]float64, variable string) {

	e.Reduction = map[int64]float64{}
	e.Expansion = map[int64]float64{}
	e.Terms = []Term

	for exponent, coefficient := range expression {

		e.Expansion[exponent] = coefficient
		e.Reduction[exponent] = coefficient

		term := Term{

			Variable: 	 variable,
			Coefficient: coefficient,
			Exponent: 	 exponent,

		}

		e.Terms = append(e.Terms, term)

	}

}

func (e *Expression) function(s float64) float64 {

	y := 0.0

	for _, term := range e.Terms {

		y += term.Coefficient*math.Pow(s, float64(term.Exponent))

	}

	return y

}

func (p *Polynomial) evaluate(s float64) float64 {

	numerator := p.Numerator.function(s)
	denominator := p.Denominator.function(s)
	zero := float64(0)
	quotient := zero

	if (denominator != zero) {

		quotient = numerator / denominator

	}

	return quotient

}


// config
type Specs struct {

	Domain					   Domain
	Response				   Response
	Approximation			   Approximation
	Configuration			   Configuration
	PassbandRipple			   *float64
	StopbandRipple			   *float64
	PassbandAttenuation		   *float64
	StopbandAttenuation		   *float64
	CutoffFrequency			   *float64
	LowerPassbandEdgeFrequency *float64
	UpperPassbandEdgeFrequency *float64
	LowerStopbandEdgeFrequency *float64
	UpperStopbandEdgeFrequency *float64
	Bandwidth				   *float64
	CenterFrequency			   *float64
	TransitionWidth			   *float64
	SamplingFrequency		   *float64
	Order					   *uint16

}

type Domain string

const (

	Analogue Domain = "analogue"
	Digital  Domain = "digital"

)

func (d Domain) exists() bool {

	switch d {

		case Analogue, Digital:

			return true

		default:

			return false

	}

}

type Response string

const (

	LPF   Response = "lpf"
	HPF   Response = "hpf"
	BPF   Response = "bpf"
	BSF   Response = "bsf"

	BRF   Response = "brf"
	Notch Response = "notch"

)

func (r Response) exists() bool {

	switch r {

		case LPF, HPF, BPF, BSF:

			return true

		case BRF, Notch:

			return true

		default:

			return false

	}

}

type Approximation string

const (

	Butterworth 	 Approximation = "butterworth"
	Chebyshev		 Approximation = "chebyshev"
	InverseChebyshev Approximation = "inverse chebyshev"
	Elliptic		 Approximation = "elliptic"

	Cauer			 Approximation = "cauer"
	Zolotarev		 Approximation = "zolotarev"

	ChebyshevType1	 Approximation = "chebyshev type 1"
	ChebyshevTypeI	 Approximation = "chebyshev type i"
	ChebyshevType2	 Approximation = "chebyshev type 2"
	ChebyshevTypeII	 Approximation = "chebyshev type ii"

	Bessel			 Approximation = "bessel"
	Thiran			 Approximation = "thiran"

)

func (a Approximation) exists() bool {

	switch a {

		case Butterworth, Chebyshev, InverseChebyshev, Elliptic:

			return true

		case ChebyshevType1, ChebyshevTypeI, ChebyshevType2, ChebyshevTypeII:

			return true

		case Cauer, Zolotarev:

			return true

		case Bessel, Thiran:

			return true

		default:

			return false

	}

}

type Configuration string

const (

	IIR     Configuration = "iir"
	FIR     Configuration = "fir"
	Active  Configuration = "active"
	Passive Configuration = "passive"

)

func (c Configuration) exists() bool {

	switch c {

		case IIR, FIR:

			return true

		case Active, Passive:

			return true

		default:

			return false

	}

}
