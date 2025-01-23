package workers

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"net/http"
	"os"
)

func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return png.Decode(resp.Body)
}

func saveImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

func new6DigitID() string {
	num := rand.Intn(999999)
	str := fmt.Sprintf("%06d", num)
	ch1 := numToLetter(int(str[0]))
	ch2 := numToLetter(int(str[1]))
	return ch1 + ch2 + str[2:6]
}

// Returns the letter in the alphabet's position of the given number.
// Supports lowercase and uppercase, and wraps around if the number is out of bounds.
//
//	n = 1  ▶ numToLetter(1)  = "a"
//	n = 2  ▶ numToLetter(2)  = "b"
//	n = 26 ▶ numToLetter(26) = "z"
//	n = 27 ▶ numToLetter(27) = "A"
//	n = 52 ▶ numToLetter(52) = "Z"
//	n = 53 ▶ numToLetter(53) = "a"
func numToLetter(n int) string {
	if n < 1 {
		n = 1
	}
	if n > 52 {
		n = 52
	}

	const asciiLowerA = 96
	return string(rune(asciiLowerA + n))
}
