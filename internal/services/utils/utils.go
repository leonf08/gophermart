package utils

func IsNumber(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func LuhnValidate(s string) bool {
	sum := 0

	parity := len(s) % 2
	for i, r := range s {
		digit := int(r - '0')

		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return sum%10 == 0
}
