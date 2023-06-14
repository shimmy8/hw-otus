package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inp string) (string, error) {
	var prevStr string
	var outpBuilder strings.Builder

	for _, val := range inp {
		if '0' <= val && val <= '9' { // current rune is a number
			if prevStr == "" || ("0" <= prevStr && prevStr <= "9") { // prev string is empty or a number
				return "", ErrInvalidString
			}

			repeats, _ := strconv.Atoi(string(val))
			if repeats > 0 { // write only repeats > 0
				outpBuilder.WriteString(strings.Repeat(prevStr, repeats))
			}
		} else if "0" > prevStr || prevStr > "9" { // current rune and prevStr are not numbers
			outpBuilder.WriteString(prevStr)
		}

		prevStr = string(val)
	}

	if "0" > prevStr || prevStr > "9" { // write last symbol if not a number
		outpBuilder.WriteString(prevStr)
	}

	return outpBuilder.String(), nil
}
