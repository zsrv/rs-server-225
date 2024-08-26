package wordenc

func filterFragments(chars []rune) {
	for currentIndex := range len(chars) {
		numberIndex := indexOfNumber(chars, currentIndex)
		if numberIndex == -1 {
			return
		}

		isSymbolOrNotLowercaseAlpha := false
		for index := currentIndex; index >= 0 && index < numberIndex && !isSymbolOrNotLowercaseAlpha; index++ {
			if !isSymbol(chars[index]) && !isLowercaseAlpha(chars[index]) {
				isSymbolOrNotLowercaseAlpha = true
			}
		}

		startIndex := 1
		currentIndex = numberIndex

		value := 0
		for index := numberIndex; index < len(chars) && index < currentIndex; index++ {
			value = value*10 + int(chars[index]) - 48
		}

		if value <= 255 && currentIndex-numberIndex <= 8 {
			startIndex++
		} else {
			startIndex = 0
		}

		if startIndex == 4 {
			maskChars(numberIndex, currentIndex, chars)
			startIndex = 0
		}
		currentIndex = indexOfNonNumber(chars, currentIndex)
	}
}

func isBadFragment(wf *WordFilter, chars []rune) bool {
	if isNumericalChars(chars) {
		return true
	}

	value := getInteger(chars)
	fragments := wf.Fragments
	fragmentsLength := len(wf.Fragments)

	if value == fragments[0] || value == fragments[fragmentsLength-1] {
		return true
	}

	start := 0
	end := fragmentsLength - 1

	for start <= end {
		mid := ((start + end) / 2) | 0
		if value == fragments[mid] {
			return true
		} else if value < fragments[mid] {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}
	return false
}

func getInteger(chars []rune) int {
	if len(chars) > 6 {
		return 0
	}
	value := 0
	for i := range len(chars) {
		char := chars[len(chars)-i-1]
		if isLowercase(char) {
			value = value*38 + int(char) + 1 - 'a'
		} else if char == '\'' {
			value = value*38 + 27
		} else if isNumber(char) {
			value = value*38 + int(char) + 28 - '0'
		} else if char != '\u0000' {
			return 0
		}
	}
	return value
}

func indexOfNumber(chars []rune, offset int) int {
	for i := offset; i < len(chars) && i >= 0; i++ {
		if isNumber(chars[i]) {
			return i
		}
	}
	return -1
}

func indexOfNonNumber(chars []rune, offset int) int {
	for i := offset; i < len(chars) && i >= 0; i++ {
		if !isNumber(chars[i]) {
			return i
		}
	}
	return len(chars)
}
