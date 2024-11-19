package rand_test

import (
	"strings"
	"testing"

	"github.com/amirhossein5/wgo/pkg/rand"
)

func TestGeneratesRandomColors(t *testing.T) {
	randomColor := rand.HexColor()

	if randomColor[0] != '#' || len(randomColor) != 7 {
		t.Fatalf("invalid random hex color given %v", randomColor)
	}

	hexLetters := strings.Split(randomColor[1:], "")

	for _, letter := range hexLetters {
		isValid := (letter >= "0" && letter <= "9") || (letter >= "a" && letter <= "f")
		if !isValid {
			t.Fatalf("invalid random hex color given %v", randomColor)
		}
	}
}
