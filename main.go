package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shopspring/decimal"
)

// Thai number words
var thaiNumbers = []string{"", "หนึ่ง", "สอง", "สาม", "สี่", "ห้า", "หก", "เจ็ด", "แปด", "เก้า"}

// Thai place values
var thaiPlaces = []string{"", "สิบ", "ร้อย", "พัน", "หมื่น", "แสน"}

// Convert an integer to Thai words
func convertIntToThaiText(number int) string {
	if number == 0 {
		return "ศูนย์"
	}

	numStr := fmt.Sprintf("%d", number)
	length := len(numStr)

	// Handling large numbers (millions, billions, etc.)
	millionIndex := length % 6
	if millionIndex == 0 {
		millionIndex = 6
	}

	chunks := []string{}
	for len(numStr) > 0 {
		if len(numStr) > millionIndex {
			chunks = append(chunks, numStr[:millionIndex])
			numStr = numStr[millionIndex:]
			millionIndex = 6
		} else {
			chunks = append(chunks, numStr)
			numStr = ""
		}
	}

	result := ""
	for i, chunk := range chunks {
		chunkNum := 0
		fmt.Sscanf(chunk, "%d", &chunkNum)
		if chunkNum == 0 {
			continue
		}

		chunkText := convertSmallIntToThaiText(chunkNum)
		if i < len(chunks)-1 {
			chunkText += "ล้าน"
		}

		result += chunkText
	}

	return result
}

// Convert small integers (under a million) to Thai words
func convertSmallIntToThaiText(number int) string {
	if number == 0 {
		return ""
	}

	result := ""
	numStr := fmt.Sprintf("%d", number)
	length := len(numStr)

	for i, digit := range numStr {
		n := int(digit - '0')
		position := length - i - 1

		if n == 0 {
			continue
		}

		if position == 1 && n == 1 {
			result += "สิบ"
		} else if position == 1 && n == 2 {
			result += "ยี่สิบ"
		} else if position == 0 && n == 1 && length > 1 {
			result += "เอ็ด"
		} else {
			result += thaiNumbers[n] + thaiPlaces[position]
		}
	}

	return result
}

// Convert a decimal number to Thai currency words
// Convert a decimal number to Thai currency words
func convertBahtToThaiText(amount decimal.Decimal) string {
	integerPart := amount.IntPart()
	decimalPart := amount.Sub(decimal.NewFromInt(integerPart)).Mul(decimal.NewFromInt(100)).IntPart()

	// Explicitly handle the case where the amount is 0
	if integerPart == 0 && decimalPart == 0 {
		return "ศูนย์บาทถ้วน"
	}

	result := ""
	if integerPart > 0 {
		result += convertIntToThaiText(int(integerPart)) + "บาท"
	}

	if decimalPart > 0 {
		result += convertSmallIntToThaiText(int(decimalPart)) + "สตางค์"
	} else if integerPart > 0 {
		result += "ถ้วน"
	}

	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter an amount (or type 'exit' to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "exit" {
			fmt.Println("Exiting program. Goodbye!")
			break
		}

		amount, err := decimal.NewFromString(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid number.")
			continue
		}

		fmt.Println("Thai currency words:", convertBahtToThaiText(amount))
	}
}
