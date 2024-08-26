package wordenc

func filterBadWords(wf *WordFilter, chars []rune) {
	for range 2 {
		for index := len(wf.BadWords) - 1; index >= 0; index-- {
			filterBadCombinations(wf, wf.BadCombinations[index], chars, wf.BadWords[index])
		}
	}
}

func filterBadCombinations(wf *WordFilter, badCombinations [][2]int, chars []rune, fragment []rune) {
	if len(fragment) > len(chars) {
		return
	}

	for start := 0; start <= len(chars)-len(fragment); start++ {
		end := start
		end, fragOff, isSymbolx, isEmulated, isNumeral := processBadCharacters(chars, fragment, end)

		if !(fragOff >= len(fragment) && (!isEmulated || !isNumeral)) {
			continue
		}

		currentChar := '\u0000'
		if end < len(chars) {
			currentChar = chars[end]
		}

		nextChar := '\u0000'
		if end+1 < len(chars) {
			nextChar = chars[end+1]
		}

		shouldFilter := true
		cur := 0

		if isSymbolx {
			badCurrent := false
			badNext := false
			if start-1 < 0 || (isSymbol(chars[start-1]) && chars[start-1] != '\'') {
				badCurrent = true
			}
			if end >= len(chars) || (isSymbol(chars[end]) && chars[end] != '\'') {
				badNext = true
			}
			if !badCurrent || !badNext {
				good := false
				cur = start - 2
				if badCurrent {
					cur = start
				}
				for !good && cur < end {
					if cur >= 0 && (!isSymbol(chars[cur]) || chars[cur] == '\'') {
						frag := make([]rune, 3)
						off := 0
						for off = 0; off < 3 &&
							cur+off < len(chars) &&
							(!isSymbol(chars[cur+off]) || chars[cur+off] == '\''); off++ {
							frag[off] = chars[cur+off]
						}
						valid := off != 0
						if off < 3 && cur-1 >= 0 && (!isSymbol(chars[cur-1]) || chars[cur-1] == '\'') {
							valid = false
						}
						if valid && !isBadFragment(wf, frag) {
							good = true
						}
					}
					cur++
				}
				if !good {
					shouldFilter = false
				}
			}
		} else {
			currentChar = ' '
			if start-1 >= 0 {
				currentChar = chars[start-1]
			}
			nextChar = ' '
			if end < len(chars) {
				nextChar = chars[end]
			}
			currentIndex := getIndex(currentChar)
			nextIndex := getIndex(nextChar)
			if badCombinations != nil && comboMatches(currentIndex, badCombinations, nextIndex) {
				shouldFilter = false
			}
		}

		if !shouldFilter {
			continue
		}

		numeralCount := 0
		alphaCount := 0
		for i := start; i < end; i++ {
			if isNumber(chars[i]) {
				numeralCount++
			} else if isAlpha(chars[i]) {
				alphaCount++
			}
		}
		if numeralCount <= alphaCount {
			maskChars(start, end, chars)
		}
	}
}

func processBadCharacters(chars []rune, fragment []rune, start int) (int, int, bool, bool, bool) {
	end := start
	fragOff := 0
	iterations := 0
	isSymbolx := false
	isEmulated := false
	isNumeral := false

	for end < len(chars) && !(isEmulated && isNumeral) {
		if end >= len(chars) || (isEmulated && isNumeral) {
			break
		}

		currentChar := chars[end]

		nextChar := '\u0000'
		if end+1 < len(chars) {
			nextChar = chars[end+1]
		}

		if fragOff < len(fragment) && getEmulatedSize(nextChar, fragment[fragOff], currentChar) > 0 {
			charLen := getEmulatedSize(nextChar, fragment[fragOff], currentChar)
			if charLen == 1 && isNumber(currentChar) {
				isEmulated = true
			}
			if charLen == 2 && (isNumber(currentChar) || isNumber(nextChar)) {
				isEmulated = true
			}
			end += charLen
			fragOff++
		} else {
			if fragOff == 0 {
				break
			}
			prevCharLen := getEmulatedSize(nextChar, fragment[fragOff-1], currentChar)
			if prevCharLen > 0 {
				end += prevCharLen
			} else {
				if fragOff >= len(fragment) || !isLowercaseAlpha(currentChar) {
					break
				}
				if isSymbol(currentChar) && currentChar != '\'' {
					isSymbolx = true
				}
				if isNumber(currentChar) {
					isNumeral = true
				}
				end++
				iterations++
				if (iterations*100)/(end-start) > 90 {
					break
				}
			}
		}
	}
	return end, fragOff, isSymbolx, isEmulated, isNumeral
}

func getEmulatedSize(nextChar rune, fragment rune, currentChar rune) int {
	if fragment == currentChar {
		return 1
	}
	if fragment >= 'a' && fragment <= 'm' {
		switch fragment {
		case 'a':
			if currentChar != '4' && currentChar != '@' && currentChar != '^' {
				if currentChar == '/' && nextChar == '\\' {
					return 2
				}
				return 0
			}
			return 1
		case 'b':
			if currentChar != '6' && currentChar != '8' {
				if currentChar == '1' && nextChar == '3' {
					return 2
				}
				return 0
			}
			return 1
		case 'c':
			if currentChar != '(' && currentChar != '<' && currentChar != '{' && currentChar != '[' {
				return 0
			}
			return 1
		case 'd':
			if currentChar == '[' && nextChar == ')' {
				return 2
			}
			return 0
		case 'e':
			if currentChar != '3' && currentChar != '€' {
				return 0
			}
			return 1
		case 'f':
			if currentChar == 'p' && nextChar == 'h' {
				return 2
			}
			if currentChar == '£' {
				return 1
			}
			return 0
		case 'g':
			if currentChar != '9' && currentChar != '6' {
				return 0
			}
			return 1
		case 'h':
			if currentChar == '#' {
				return 1
			}
			return 0
		case 'i':
			if currentChar != 'y' && currentChar != 'l' && currentChar != 'j' &&
				currentChar != '1' && currentChar != '!' && currentChar != ':' &&
				currentChar != ';' && currentChar != '|' {
				return 0
			}
			return 1
		case 'j':
			return 0
		case 'k':
			return 0
		case 'l':
			if currentChar != '1' && currentChar != '|' && currentChar != 'i' {
				return 0
			}
			return 1
		case 'm':
			return 0
		}
	}
	if fragment >= 'n' && fragment <= 'z' {
		switch fragment {
		case 'n':
			return 0
		case 'o':
			if currentChar != '0' && currentChar != '*' {
				if (currentChar != '(' || nextChar != ')') &&
					(currentChar != '[' || nextChar != ']') &&
					(currentChar != '{' || nextChar != '}') &&
					(currentChar != '<' || nextChar != '>') {
					return 0
				}
				return 2
			}
			return 1
		case 'p':
			return 0
		case 'q':
			return 0
		case 'r':
			return 0
		case 's':
			if currentChar != '5' && currentChar != 'z' && currentChar != '$' && currentChar != '2' {
				return 0
			}
			return 1
		case 't':
			if currentChar != '7' && currentChar != '+' {
				return 0
			}
			return 1
		case 'u':
			if currentChar == 'v' {
				return 1
			}
			if (currentChar != '\\' || nextChar != '/') &&
				(currentChar != '\\' || nextChar != '|') &&
				(currentChar != '|' || nextChar != '/') {
				return 0
			}
			return 2
		case 'v':
			if (currentChar != '\\' || nextChar != '/') &&
				(currentChar != '\\' || nextChar != '|') &&
				(currentChar != '|' || nextChar != '/') {
				return 0
			}
			return 2
		case 'w':
			if currentChar == 'v' && nextChar == 'v' {
				return 2
			}
			return 0
		case 'x':
			if (currentChar != ')' || nextChar != '(') &&
				(currentChar != '}' || nextChar != '{') &&
				(currentChar != ']' || nextChar != '[') &&
				(currentChar != '>' || nextChar != '<') {
				return 0
			}
			return 2
		case 'y':
			return 0
		case 'z':
			return 0
		}
	}
	if fragment >= '0' && fragment <= '9' {
		switch fragment {
		case '0':
			if currentChar == 'o' || currentChar == 'O' {
				return 1
			} else if (currentChar != '(' || nextChar != ')') &&
				(currentChar != '{' || nextChar != '}') &&
				(currentChar != '[' || nextChar != ']') {
				return 0
			} else {
				return 2
			}
		case '1':
			if currentChar == 'l' {
				return 1
			} else {
				return 0
			}
		default:
			return 0
		}
	}
	if fragment == ',' {
		if currentChar == '.' {
			return 1
		} else {
			return 0
		}
	}
	if fragment == '.' {
		if currentChar == ',' {
			return 1
		} else {
			return 0
		}
	}
	if fragment == '!' {
		if currentChar == 'i' {
			return 1
		} else {
			return 0
		}
	}
	return 0
}

func comboMatches(currentIndex int, combos [][2]int, nextIndex int) bool {
	start := 0
	end := len(combos) - 1

	for start <= end {
		mid := ((start + end) / 2) | 0
		if combos[mid][0] == currentIndex && combos[mid][1] == nextIndex {
			return true
		} else if currentIndex < combos[mid][0] || (currentIndex == combos[mid][0] && nextIndex < combos[mid][1]) {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}
	return false
}

func getIndex(char rune) int {
	if isLowercase(char) {
		return int(char - 1 - 'a')
	} else if char == '\'' {
		return 28
	} else if isNumber(char) {
		return int(char + 29 - '0')
	}
	return 27
}
