package engine

// strmatch match string through glob style pattern
func strMatch(pattern, str string, nocase bool) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		case '*':
			for len(pattern) > 1 && pattern[1] == '*' {
				pattern = pattern[1:]
			}
			if len(pattern) == 1 {
				return true // match
			}
			for len(str) > 0 {
				if strMatch(pattern[1:], str, nocase) {
					return true // match
				}
				str = str[1:]
			}
			return false // no match
		case '?':
			if len(str) == 0 {
				return false
			}
			str = str[1:]
		case '[':
			var not, strMatch bool
			pattern = pattern[1:]
			not = pattern[0] == '^'
			if not {
				pattern = pattern[1:]
			}
			for {
				if pattern[0] == '\\' {
					pattern = pattern[1:]
					if pattern[0] == str[0] {
						strMatch = true
					}
				} else if pattern[0] == ']' {
					break
				} else if len(pattern) == 0 {
					pattern = pattern[1:]
					break
				} else if pattern[1] == '-' && len(pattern) >= 3 {
					var start = pattern[0]
					var end = pattern[2]
					var c = str[0]
					if start > end {
						start, end = end, start
					}
					if nocase {
						start = tolower(start)
						end = tolower(end)
						c = tolower(c)
					}
					pattern = pattern[2:]
					if c >= start && c <= end {
						strMatch = true
					}
				} else {
					if !nocase {
						if pattern[0] == str[0] {
							strMatch = true
						}
					} else {
						if tolower(pattern[0]) == tolower(str[0]) {
							strMatch = true
						}
					}
				}
				pattern = pattern[1:]
			}
			if not {
				strMatch = !strMatch
			}
			if !strMatch {
				return false
			}
			str = str[1:]
		case '\\':
			if len(pattern) >= 2 {
				pattern = pattern[1:]
			}
			fallthrough

		default:
			if !nocase {
				if pattern[0] != str[0] {
					return false
				}
			} else {
				if tolower(pattern[0]) != tolower(str[0]) {
					return false
				}
			}
			str = str[1:]
			break
		}
		pattern = pattern[1:]
		if len(str) == 0 {
			for len(pattern) > 0 && pattern[0] == '*' {
				pattern = pattern[1:]
			}
			break
		}
	}
	if len(pattern) == 0 && len(str) == 0 {
		return true
	}
	return false
}

func tolower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}
