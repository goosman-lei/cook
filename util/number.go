package util

// slice must sort first
func Uniq_int(s []int) []int {
	var (
		sl   int = len(s)
		p, i int
	)
	if sl <= 1 {
		return s
	}

	for p, i = 0, 1; i < sl; i++ {
		if s[i] != s[p] {
			if i != p {
				s[p+1] = s[i]
			}
			p++
		}
	}
	return s[:p+1]
}
