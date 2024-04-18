package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func validateCards(cardNumbers []string) {
	pattern := regexp.MustCompile(`^[456]\d{3}-?\d{4}-?\d{4}-?\d{4}$`)
	repetitionPattern := regexp.MustCompile(`(.)\\1{3}`)

	for _, cardNumber := range cardNumbers {
		cardNumber2 := strings.ReplaceAll(cardNumber, "-", "")
		if pattern.MatchString(cardNumber) && !repetitionPattern.MatchString(cardNumber2) {
			fmt.Println("Valid")
		} else {
			fmt.Println("Invalid")
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var N int
	fmt.Scanln(&N)

	cardNumbers := make([]string, N)
	for i := 0; i < N && scanner.Scan(); i++ {
		cardNumbers[i] = scanner.Text()
	}

	validateCards(cardNumbers)
}
