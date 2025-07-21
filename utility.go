package main

import ( "strings"
		 "strconv" )


func contains(array []string, search string) bool {

	flag := false

	for _, entry := range array {

		if strings.Contains(search, entry) {

			flag = true
			break

		}

	}

	return flag

}

func getSuperscript(exponent string) string {

	power := ""

	if (exponent[0] == '-') {

		power = "⁻" + power

	}

	for _, digit := range exponent {

		number, err := strconv.Atoi(string(digit))

		if (err == nil) {

			power += superscript[number]

		}

	}

	return power

}
