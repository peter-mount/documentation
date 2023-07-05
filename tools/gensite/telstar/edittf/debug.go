package edittf

func ToLines(s string) []string {
	var r []string
	for i := 0; i < len(s); i += 40 {
		r = append(r, s[i:i+40])
	}
	return r
}
