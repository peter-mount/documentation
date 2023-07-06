package edittf

import "fmt"

// PadLines ensures the slice is 25 entries long & each entry is 40 characters, padded by spaces as required
func PadLines(frame []string) []string {
	if len(frame) > 25 {
		frame = frame[:25]
	}

	if len(frame) < 25 {
		// Add row to start & end so that we have 25 rows
		frame = append([]string{""}, frame...)
		for len(frame) < 25 {
			frame = append(frame, "")
		}
	}

	// Ensure lines are 40 chars
	for r, s := range frame {
		frame[r] = (s + "                                        ")[:40]
	}

	return frame
}

// Encode a frame into edit.tf format
func Encode(frame []string) string {
	frame = PadLines(frame)

	b64 := make([]byte, 1167)

	for r := 0; r < 25; r++ {
		for c := 0; c < 40; c++ {
			for b := 0; b < 7; b++ {
				framebit := 7*((r*40)+c) + b
				b64bitoffset := framebit % 6
				b64charoffset := (framebit - b64bitoffset) / 6

				bitval := frame[r][c] & (1 << (6 - b))
				if bitval > 0 {
					bitval = 1
				}

				b64[b64charoffset] |= bitval << (5 - b64bitoffset)
			}
		}
	}

	fmt.Println("Encode")
	var url []byte
	url = append(url, []byte("https://edit.tf/#0:")...)
	for _, c := range b64 {
		url = append(url, alphabet[c])
	}

	return string(url)
}
