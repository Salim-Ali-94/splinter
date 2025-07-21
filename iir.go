package main

import ( "strings"
		 "math" )


func designFilter(config Specs) Polynomial {

	domain, configuration,
	response, approximation,
	order, ripplePassband,
	rippleStopband, attenuationPassband,
	attenuationStopband, cutOffFrequency,
	lowerPassbandEdgeFrequency, upperPassbandEdgeFrequency,
	lowerStopbandEdgeFrequency, upperStopbandEdgeFrequency,
	bandwidth, centerFrequency, samplingPeriod = parseDesign(config)

	if (order == 0) {

		order = calculateOrder(approximation, response,
							   ripplePassband, rippleStopband,
							   attenuationPassband, attenuationStopband,
							   cutOffFrequency, lowerPassbandEdgeFrequency,
							   upperPassbandEdgeFrequency, lowerStopbandEdgeFrequency,
							   upperStopbandEdgeFrequency, bandwidth,
							   centerFrequency, samplingPeriod)

	}

	ltf = analogueLowPassFilterPrototype(approximation, order, epsilonPass, epsilonStop)
	// laplaceTransform
	tf := bilinearTransform(approximation, response, order)
	return tf

}

func analogueLowPassFilterPrototype(approximation string, order uint16, epsilonPass float64, epsilonStop float64) Polynomial {

	polynomial := constructPolynomial()

	if (approximation == "butterworth") {

		polynomial = butterworthTransferFunction(order)

	} else if contains(chebyshev1, approximation) {

		polynomial = chebyshevTransferFunction(order, epsilonPass)

	} else if contains(chebyshev2, approximation) {

		polynomial = inverseChebyshevTransferFunction(order, epsilonStop)

	} else if contains(cauer, approximation) {

		polynomial = ellipticTransferFunction(order, epsilonPass, epsilonStop)

	}

	return polynomial

}

func butterworthTransferFunction(order uint16) Polynomial {

	tf := constructPolynomial()
	initial := 1

	if (order%2 != 0) {

		term := map[string]interface{}{

			"variable": "s",
			"numerator": map[int64]float64{ 0: 1.0 },
			"denominator": map[int64]float64{ 0: 1.0, 1: 1.0 },

		}

		tf = constructPolynomial(term)
		initial = 2

	}

	for index := initial; index <= order - 1; index++ {

		angle := float64(index)*math.Pi / 6.0
		c := 2.0*math.Cos(angle)

		term := map[string]interface{}{

			"variable": "s",
			"numerator": map[int64]float64{ 0: 1.0 },
			"denominator": map[int64]float64{ 0: 1.0, 1: c, 2: 1.0 },

		}

		polynomial := constructPolynomial(term)
		tf.multiply(polynomial)

	}

	return tf

}

func chebyshevTransferFunction(order uint16, epsilonPass float64) Polynomial {

	tf := constructPolynomial()
	initial := 1
	d := math.Asinh(1.0 / epsilonPass) / float64(order)
	alpha := math.Sinh(d)
	beta := math.Cosh(d)

	if (order%2 != 0) {

		term := map[string]interface{}{

			"variable": "s"
			"numerator": map[int64]float64{ 0: 1.0 },
			"denominator": map[int64]float64{ 0: alpha, 1: 1.0 },

		}

		tf = constructPolynomial(term)
		initial = 2

	}

	for index := initial; index <= order - 1; index++ {

		angle := float64(2*index + 1)*math.Pi / 2.0*float64(order)
		a := -alpha*math.Sin(angle)
		b := beta*math.Cos(angle)
		c := math.Pow(a, 2) + math.Pow(b, 2)
		m := -2.0*a

		term := map[string]interface{}{

			"variable": "s",
			"numerator": map[int64]float64{ 0: 1.0 },
			"denominator": map[int64]float64{ 0: c, 1: m, 2: 1.0 },

		}

		polynomial := constructPolynomial(term)
		tf.multiply(polynomial)

	}

	return tf

}

func inverseChebyshevTransferFunction(order uint16, epsilonStop float64) Polynomial {

	tf := constructPolynomial()
	initial := 1
	d := math.Asinh(1.0 / epsilonStop) / float64(order)
	alpha := math.Sinh(d)
	beta := math.Cosh(d)

	if (order%2 != 0) {

		c := 1.0 / alpha

		term := map[string]interface{}{

			"variable": "s",
			"numerator": map[int64]float64{ 0: 1.0 },
			"denominator": map[int64]float64{ 0: c, 1: 1.0 },

		}

		tf = constructPolynomial(term)
		initial = 2

	}

	for index := initial; index <= order - 1; index++ {

		angle := (2.0*float64(index) + 1.0)*math.Pi / 2.0*float64(order)
		zero := 1.0 / math.Cos(angle)
		zero = math.Pow(zero, 2)
		a := -alpha*math.Sin(angle)
		b := beta*math.Cos(angle)
		c := math.Pow(a, 2) + math.Pow(b, 2)
		p := -2.0*a / c
		m := 1.0 / c

		term := map[int64]float64{

			"variable": "s",
			"numerator": map[int64]float64{ 0: zero, 1: 1.0 },
			"denominator": map[int64]float64{ 0: m, 1: p, 2: 1.0 },

		}

		polynomial := constructPolynomial(term)
		tf.multiply(polynomial)

	}

	return tf

}

func ellipticTransferFunction(order uint16, epsilonPass float64, epsilonStop float64) Polynomial {

	tf := constructPolynomial()

	sn, cn, dn, dSn, dCn, dDn := landenLoop()

	return tf

}

func butterworthOrder(epsilonPass float64, epsilonStop float64, normalizedFrequency float64) uint16 {

	numerator := math.Log(epsilonStop / epsilonPass)
	denominator := 2.0*math.Log(normalizedFrequency)
	denominator = math.Abs(denominator)
	order := numerator / denominator
	order = math.Ceil(order)
	return uint64(order)

}

func chebyshevOrder(epsilonPass float64, epsilonStop float64, normalizedFrequency float64) uint64 {

	epsilon := epsilonStop / epsilonPass
	epsilon = math.Sqrt(epsilon)
	numerator := math.Acosh(epsilon)
	denominator := math.Acosh(normalizedFrequency)
	order := numerator / denominator
	order = math.Ceil(order)
	return uint64(order)

}

func ellipticOrder(epsilonPass float64, epsilonStop float64, normalizedFrequency float64, response string) uint64 {

	epsilon := epsilonPass / epsilonStop
	epsilon = math.Sqrt(epsilon)
	tau := 1.0 / normalizedFrequency

	E := math.Pow(epsilon, 2)
	E = math.Sqrt(1.0 - E)
	T := math.Pow(tau, 2)
	T = math.Sqrt(1.0 - T)

	precision := 1e-9
	a := arithmeticGeometricMean(epsilon, precision)
	b := arithmeticGeometricMean(tau, precision)
	c := arithmeticGeometricMean(E, precision)
	d := arithmeticGeometricMean(T, precision)

	numerator := b*c
	denominator := a*d
	order := numerator / denominator
	order = math.Ceil(order)
	return uint64(order)

}

func arithmeticGeometricMean(k float64, tolerance float64) float64 {

	a := 1.0
	b := math.Sqrt(1.0 - math.Pow(k, 2))
	epsilon := math.Abs(tolerance)
	error := 2.0*epsilon

	for (error > epsilon) {

		A := (a + b) / 2.0
		B := math.Sqrt(a*b)
		error = math.Abs(A - B)
		a = A
		b = B

	}

	M := 2.0*a
	K := math.Pi / M
	return K

}

func landenIterations(a float64, b float64, tolerance float64) (float64, float64, float64,float64, float64, float64) {

	sn := 0.0
	cn := 0.0
	dn := 0.0
	dSn := 0.0
	dCn := 0.0
	dDn := 0.0
	u := a
	k := b
	epsilon := math.Abs(tolerance)
	error := 2.0*epsilon
	kVector := []float64{}
	pVector := []float64{}

	for (error > epsilon) {

		c := math.Pow(k, 2)
		c = math.Sqrt(1.0 - k)
		k := (1.0 - c) / (1.0 + c)
		u := u*(1.0 + c) / 2.0
		kVector = append([]float64{ k }, kVector...)
		pVector = append([]float64{ c }, pVector...)

	}

	sn = math.Sin(u)
	cn = math.cos(u)
	dn = 1.0

	for (index := 0; index < len(kVector); index++) {

		c := kVector[index + 1]
		d := pVector[index + 1]
		s := math.Pow(sn, 2)
		sn = sn*(1 + d) / (1 + c*s)
		cn = s*(1 - d) / (1 + c*s)

	}

	return sn, cn, dn, dSn, dCn, dDn

}

func calculateOrder(approximation string, response string,
					ripplePassband float64, rippleStopband float64,
					attenuationPassband float64, attenuationStopband float64,
					cutOffFrequency float64, lowerPassbandEdgeFrequency float64,
					upperPassbandEdgeFrequency float64, lowerStopbandEdgeFrequency float64,
					upperStopbandEdgeFrequency float64, bandwidth float64,
					centerFrequency float64, samplingPeriod float64) uint64 {

	order := 0.0
	epsilonPass := math.Pow(ripplePassband, 2)
	epsilonStop := math.Pow(rippleStopband, 2)
	normalizedFrequency := 0.0

	if ((attenuationPassband > 0.0) &&
		(ripplePassband == 0.0)) {

		epsilonPass = math.Pow(10, 0.1*attenuationPassband) - 1.0

	}

	if ((attenuationStopband > 0.0) &&
		(rippleStopband == 0.0)) {

		epsilonStop = math.Pow(10, attenuationStopband / 10.0) - 1.0

	}

	if ((response == "lpf") ||
		(response == "hpf")) {

		frequencyDesignPass := 2.0*math.Pi*cutOffFrequency
		frequencyWarpedPass := 2.0*math.Tan(frequencyDesignPass*samplingPeriod / 2.0) / samplingPeriod
		frequencyDesignStop := 2.0*math.Pi*cutOffFrequency
		frequencyWarpedStop := 2.0*math.Tan(frequencyDesignStop*samplingPeriod / 2.0) / samplingPeriod
	
		if (response == "lpf") {

			normalizedFrequency = frequencyWarpedStop / frequencyWarpedPass

		} else if (response == "hpf") {

			normalizedFrequency = frequencyWarpedPass / frequencyWarpedStop

		}

	} else if (contains(bsf, response) ||
			   (response == "bpf")) {

		frequencyDesignCenter := 2.0*math.Pi*centerFrequency
		frequencyWarpedCenter := 2.0*math.Tan(frequencyDesignCenter*samplingPeriod / 2.0) / samplingPeriod

		frequencyDesignPassLower := 2.0*math.Pi*lowerPassbandEdgeFrequency
		frequencyWarpedPassLower := 2.0*math.Tan(frequencyDesignPassLower*samplingPeriod / 2.0) / samplingPreiod
		frequencyAnaloguePassUpper := math.Pow(frequencyWarpedCenter, 2) / frequencyWarpedPassLower
		bandwidthWarpedPass := frequencyAnaloguePassUpper - frequencyWarpedPassLower

		frequencyDesignStopLower := 2.0*math.Pi*lowerStopbandEdgeFrequency
		frequencyWarpedStopLower := 2.0*math.Tan(frequencyDesignStopLower*samplingPeriod / 2.0) / samplingPreiod
		frequencyAnalogueStopUpper := math.Pow(frequencyWarpedCenter, 2) / frequencyWarpedStopLower
		bandwidthWarpedStop := frequencyAnalogueStopUpper - frequencyWarpedStopLower

		frequencyDesignPassUpper := 2.0*math.Pi*upperPassbandEdgeFrequency
		frequencyWarpedPassUpper := 2.0*math.Tan(frequencyDesignPassUpper*samplingPeriod / 2.0) / samplingPreiod
		frequencyAnaloguePassLower := math.Pow(frequencyWarpedCenter, 2) / frequencyWarpedPassUpper
		bandwidthWarpedPassband := frequencyAnaloguePassLower - frequencyWarpedPassUpper

		frequencyDesignStopUpper := 2.0*math.Pi*upperStopbandEdgeFrequency
		frequencyWarpedStopUpper := 2.0*math.Tan(frequencyDesignStopUpper*samplingPeriod / 2.0) / samplingPreiod
		frequencyAnalogueStopLower := math.Pow(frequencyWarpedCenter, 2) / frequencyWarpedStopUpper
		bandwidthWarpedStopband := frequencyAnalogueStopLower - frequencyWarpedStopUpper

		if (response == "bpf") {

			if (bandwidthWarpedPass < bandwidthWarpedPassband) {

				normalizedFrequency = bandwidthWarpedStop / bandwidthWarpedPass

			} else {

				normalizedFrequency = bandwidthWarpedStopband / bandwidthWarpedPassband

			}

		} else {

			if (bandwidthWarpedStop > bandwidthWarpedStopband) {

				normalizedFrequency = bandwidthWarpedPass / bandwidthWarpedStop

			} else {

				normalizedFrequency = bandwidthWarpedPassband / bandwidthWarpedStopband

			}

		}

	}

	if (approximation == "butterworth") {

		order = butterworthOrder(epsilonPass, epsilonStop, normalizedFrequency)

	} else if strings.Contains(approximation, "chebyshev") {

		order = chebyshevOrder(epsilonPass, epsilonStop, normalizedFrequency)

	} else if contains(cauer, approximation) {

		order = ellipticOrder(epsilonPass, epsilonStop, normalizedFrequency, response)

	}

	return order

}

func parseDesign(config Specs) (string, string,
								string, string,
								uint16, float64,
							 	float64, float64,
							 	float64, float64,
							 	float64, float64,
							 	float64, float64,
							 	float64, float64,
							 	float64) {

	domain := ""

	if config.Domain.exists() {

		domain = string(config.Domain)

	}

	response := ""

	if config.Response.exists() {

		response = string(config.Response)

	}

	approximation := ""

	if config.Approximation.exists() {

		approximation = string(config.Approximation)

	}

	configuration := ""

	if config.Configuration.exists() {

		configuration = string(config.Configuration)

	}

	order := 0

	if (config.Order != nil) {

		order = uint16(math.Abs(*config.Order))

	}

	ripplePassband := 0.0

	if (config.PassbandRipple != nil) {

		ripplePassband = *config.PassbandRipple

		if (ripplePassband < 0.0) {

			ripplePassband = math.Abs(ripplePassband)

		}

	}

	rippleStopband := 0.0

	if (config.StopbandRipple != nil) {

		rippleStopband = *config.StopbandRipple

		if (rippleStopband < 0.0) {

			rippleStopband = math.Abs(rippleStopband)

		}

	}

	attenuationPassband := 0.0

	if (config.PassbandAttenuation != nil) {

		attenuationPassband = *config.PassbandAttenuation

		if (attenuationPassband < 0.0) {

			attenuationPassband = math.Abs(attenuationPassband)

		}

	}

	attenuationStopband := 0.0

	if (config.StopbandAttenuation != nil) {

		attenuationStopband = *config.StopbandAttenuation

		if (attenuationStopband < 0.0) {

			attenuationStopband = math.Abs(attenuationStopband)

		}

	}

	cutOffFrequency := 0.0

	if (config.CutoffFrequency != nil) {

		cutOffFrequency = *config.CutoffFrequency

		if (cutOffFrequency < 0.0) {

			cutOffFrequency = math.Abs(cutOffFrequency)

		}

	}

	lowerPassbandEdgeFrequency := 0.0

	if (config.LowerPassbandEdgeFrequency != nil) {

		lowerPassbandEdgeFrequency = *config.LowerPassbandEdgeFrequency

		if (lowerPassbandEdgeFrequency < 0.0) {

			lowerPassbandEdgeFrequency = math.Abs(lowerPassbandEdgeFrequency)

		}

	}

	upperPassbandEdgeFrequency := 0.0

	if (config.UpperPassbandEdgeFrequency != nil) {

		upperPassbandEdgeFrequency = *config.UpperPassbandEdgeFrequency

		if (upperPassbandEdgeFrequency < 0.0) {

			upperPassbandEdgeFrequency = math.Abs(upperPassbandEdgeFrequency)

		}

	}

	if ((lowerPassbandEdgeFrequency > 0.0) &&
		(upperPassbandEdgeFrequency > 0.0)) {

		if (upperPassbandEdgeFrequency < lowerPassbandEdgeFrequency) {

			remember := upperPassbandEdgeFrequency
			upperPassbandEdgeFrequency = lowerPassbandEdgeFrequency
			lowerPassbandEdgeFrequency = remember

		}

	}

	lowerStopbandEdgeFrequency := 0.0

	if (config.LowerStopbandEdgeFrequency != nil) {

		lowerStopbandEdgeFrequency = *config.LowerStopbandEdgeFrequency

		if (lowerStopbandEdgeFrequency < 0.0) {

			lowerStopbandEdgeFrequency = math.Abs(lowerStopbandEdgeFrequency)

		}

	}

	upperStopbandEdgeFrequency := 0.0

	if (config.UpperStopbandEdgeFrequency != nil) {

		upperStopbandEdgeFrequency = *config.UpperStopbandEdgeFrequency

		if (upperStopbandEdgeFrequency < 0.0) {

			upperStopbandEdgeFrequency = math.Abs(upperStopbandEdgeFrequency)

		}

	}

	if ((lowerStopbandEdgeFrequency > 0.0) &&
		(upperStopbandEdgeFrequency > 0.0)) {

		if (lowerStopbandEdgeFrequency > upperStopbandEdgeFrequency) {

			remember := lowerStopbandEdgeFrequency
			lowerStopbandEdgeFrequency = upperStopbandEdgeFrequency
			upperStopbandEdgeFrequency = remember

		}

	}

	bandwidth := 0.0

	if (config.Bandwidth != nil) {

		bandwidth = *config.Bandwidth

		if (bandwidth < 0.0) {

			bandwidth = math.Abs(bandwidth)

		}

	} else if ((lowerPassbandEdgeFrequency > 0.0) &&
			   (upperPassbandEdgeFrequency > 0.0) &&
			   (response == "bpf")) {

		bandwidth = upperPassbandEdgeFrequency - lowerPassbandEdgeFrequency

	} else if ((lowerStopbandEdgeFrequency > 0.0) &&
			   (upperStopbandEdgeFrequency > 0.0)) {

		bandwidth = upperStopbandEdgeFrequency - lowerStopbandEdgeFrequency

	}

	centerFrequency := 0.0

	if (config.CenterFrequency != nil) {

		centerFrequency = *config.CenterFrequency

		if (centerFrequency < 0.0) {

			centerFrequency = math.Abs(centerFrequency)

		}

	} else if ((lowerPassbandEdgeFrequency > 0.0) &&
			   (upperPassbandEdgeFrequency > 0.0)) {

		centerFrequency = math.Sqrt(lowerPassbandEdgeFrequency*upperPassbandEdgeFrequency)

	} else if ((lowerStopbandEdgeFrequency > 0.0) &&
			   (upperStopbandEdgeFrequency > 0.0)) {

		centerFrequency = math.Sqrt(lowerStopbandEdgeFrequency*upperStopbandEdgeFrequency)

	}

	if ((bandwidth > 0.0) &&
		(centerFrequency > 0.0)) {

		if (lowerPassbandEdgeFrequency == 0.0) {

			lowerPassbandEdgeFrequency = centerFrequency - bandwidth / 2.0

		} else if (lowerStopbandEdgeFrequency == 0.0) {

			lowerStopbandEdgeFrequency = centerFrequency - bandwidth / 2.0

		}

		if (upperPassbandEdgeFrequency == 0.0) {

			upperPassbandEdgeFrequency = centerFrequency + bandwidth / 2.0

		} else if (upperStopbandEdgeFrequency == 0.0) {

			upperStopbandEdgeFrequency = centerFrequency + bandwidth / 2.0

		}

	}

	transitionWidth := 0.0

	if (config.TransitionWidth != nil) {

		transitionWidth = *config.TransitionWidth

		if (transitionWidth < 0.0) {

			transitionWidth = math.Abs(transitionWidth)

		}

	}

	if ((transitionWidth > 0.0) && (lowerPassbandEdgeFrequency == 0.0) ||
		(transitionWidth > 0.0) && (lowerStopbandEdgeFrequency == 0.0) ||
		(transitionWidth > 0.0) && (upperPassbandEdgeFrequency == 0.0) ||
		(transitionWidth > 0.0) && (upperStopbandEdgeFrequency == 0.0)) {

		if (lowerPassbandEdgeFrequency == 0.0) {

			lowerPassbandEdgeFrequency = lowerStopbandEdgeFrequency + transitionWidth

		} else if (lowerStopbandEdgeFrequency == 0.0) {

			lowerStopbandEdgeFrequency = lowerPassbandEdgeFrequency + transitionWidth

		}

		if (upperPassbandEdgeFrequency == 0.0) {

			upperPassbandEdgeFrequency = upperStopbandEdgeFrequency - transitionWidth

		} else if (upperStopbandEdgeFrequency == 0.0) {

			upperStopbandEdgeFrequency = upperPassbandEdgeFrequency - transitionWidth

		}

		if (lowerStopbandEdgeFrequency == 0.0) {

			lowerStopbandEdgeFrequency = lowerPassbandEdgeFrequency - transitionWidth

		} else if (lowerPassbandEdgeFrequency == 0.0) {

			lowerPassbandEdgeFrequency = lowerStopbandEdgeFrequency - transitionWidth

		}

		if (upperStopbandEdgeFrequency == 0.0) {

			upperStopbandEdgeFrequency = upperPassbandEdgeFrequency + transitionWidth

		} else if (upperPassbandEdgeFrequency == 0.0) {

			upperPassbandEdgeFrequency = upperStopbandEdgeFrequency + transitionWidth

		}

	}

	samplingFrequency := 0.0

	if (config.SamplingFrequency != nil) {

		samplingFrequency = *config.SamplingFrequency

		if (samplingFrequency < 0.0) {

			samplingFrequency = math.Abs(samplingFrequency)

		}

	}

	samplingPeriod := 0.0

	if (samplingFrequency != 0.0) {

		samplingPeriod = 1.0 / samplingFrequency

	}

	if (domain == "digital") {

		if (samplingPeriod == 0.0) {



		} else if (order == 0) {

			if ((ripplePassband == 0.0) ||
				(attenuationPassband == 0.0)) {

				if ((rippleStopband == 0.0) ||
				 	(attenuationStopband == 0.0)) {



				 	}

			} else if ((rippleStopband == 0.0) ||
					   (attenuationStopband == 0.0)) {

				if ((ripplePassband == 0.0) ||
				 	(attenuationPassband == 0.0)) {



				 	}

			}


		}

	} else {



	}

	return domain, response,
		   approximation, configuration,
		   order, ripplePassband,
		   rippleStopband, attenuationPassband,
		   attenuationStopband, cutOffFrequency,
		   lowerPassbandEdgeFrequency, upperPassbandEdgeFrequency,
		   lowerStopbandEdgeFrequency, upperStopbandEdgeFrequency,
		   bandwidth, centerFrequency, samplingPeriod

}
