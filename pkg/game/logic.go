package game

func IsDivisibleBy3(n int) bool {
	sum := 0
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum%3 == 0
}

func NeededToMakeDivisible(n int) int {
	r := n % 3
	if r == 1 {
		return 2
	} else if r == 2 {
		return 1
	}
	return 0
}
