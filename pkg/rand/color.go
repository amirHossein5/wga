package rand

import "math/rand/v2"

func HexColor() string {
	allowed := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	return "#" + allowed[rand.IntN(len(allowed))] +
		allowed[rand.IntN(len(allowed))] +
		allowed[rand.IntN(len(allowed))] +
		allowed[rand.IntN(len(allowed))] +
		allowed[rand.IntN(len(allowed))] +
		allowed[rand.IntN(len(allowed))]
}
