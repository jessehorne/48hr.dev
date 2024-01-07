package mb

func ForN(s string, n int) string {
	end := ""
	for i, c := range s {
		if i < n {
			end += string(c)
		}
	}
	return end
}
