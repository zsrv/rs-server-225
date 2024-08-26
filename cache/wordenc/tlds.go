package wordenc

func filterTLD(wf *WordFilter, chars []rune) {
	filteredDot := make([]rune, len(chars))
	copy(filteredDot, chars)
	filterBadCombinations(wf, nil, filteredDot, []rune{'d', 'o', 't'})

	filteredSlash := make([]rune, len(chars))
	copy(filteredSlash, chars)
	filterBadCombinations(wf, nil, filteredSlash, []rune{'s', 'l', 'a', 's', 'h'})

	for i := range wf.TLDs {
		filterTLDX(filteredSlash, wf.TLDTypes[i], chars, wf.TLDs[i], filteredDot)
	}
}

func filterTLDX(slash []rune, tldType uint8, chars []rune, tld []rune, period []rune) {
	if len(tld) > len(chars) {
		return
	}

	for index := 0; index <= len(chars)-len(tld); index++ {
		currentIndex, tldIndex := processTLDs(chars, tld, index)

		if tldIndex < len(tld) {
			continue
		}

		shouldFilter := false
		periodFilterStatus := prefixSymbolStatus(index, chars, 3, period, []rune{',', '.'})
		slashFilterStatus := suffixSymbolStatus(currentIndex-1, chars, 5, slash, []rune{'\\', '/'})
		if tldType == 1 && periodFilterStatus > 0 && slashFilterStatus > 0 {
			shouldFilter = true
		}
		if tldType == 2 && ((periodFilterStatus > 2 && slashFilterStatus > 0) || (periodFilterStatus > 0 && slashFilterStatus > 2)) {
			shouldFilter = true
		}
		if tldType == 3 && periodFilterStatus > 0 && slashFilterStatus > 2 {
			shouldFilter = true
		}
		if !shouldFilter {
			continue
		}

		startFilterIndex := index
		endFilterIndex := currentIndex - 1
		foundPeriod := false
		if periodFilterStatus > 2 {
			if periodFilterStatus == 4 {
				foundPeriod = false
				for periodIndex := index - 1; periodIndex >= 0; periodIndex-- {
					if foundPeriod {
						if period[periodIndex] != '*' {
							break
						}
						startFilterIndex = periodIndex
					} else if period[periodIndex] == '*' {
						startFilterIndex = periodIndex
						foundPeriod = true
					}
				}
			}
			foundPeriod = false
			for periodIndex := startFilterIndex - 1; periodIndex >= 0; periodIndex-- {
				if foundPeriod {
					if isSymbol(chars[periodIndex]) {
						break
					}
					startFilterIndex = periodIndex
				} else if !isSymbol(chars[periodIndex]) {
					foundPeriod = true
					startFilterIndex = periodIndex
				}
			}
		}
		if slashFilterStatus > 2 {
			if slashFilterStatus == 4 {
				foundPeriod = false
				for periodIndex := endFilterIndex + 1; periodIndex < len(chars); periodIndex++ {
					if foundPeriod {
						if slash[periodIndex] != '*' {
							break
						}
						endFilterIndex = periodIndex
					} else if slash[periodIndex] == '*' {
						endFilterIndex = periodIndex
						foundPeriod = true
					}
				}
			}
			foundPeriod = false
			for periodIndex := endFilterIndex + 1; periodIndex < len(chars); periodIndex++ {
				if foundPeriod {
					if isSymbol(chars[periodIndex]) {
						break
					}
					endFilterIndex = periodIndex
				} else if !isSymbol(chars[periodIndex]) {
					foundPeriod = true
					endFilterIndex = periodIndex
				}
			}
		}
		maskChars(startFilterIndex, endFilterIndex+1, chars)
	}
}

func processTLDs(chars []rune, tld []rune, currentIndex int) (int, int) {
	tldIndex := 0
	for currentIndex < len(chars) && tldIndex < len(tld) {
		currentChar := chars[currentIndex]

		nextChar := '\u0000'
		if currentIndex+1 < len(chars) {
			nextChar = chars[currentIndex+1]
		}

		currentLength := getEmulatedDomainCharLen(nextChar, tld[tldIndex], currentChar)
		if currentLength > 0 {
			currentIndex += currentLength
			tldIndex++
		} else {
			if tldIndex == 0 {
				break
			}
			previousLength := getEmulatedDomainCharLen(nextChar, tld[tldIndex-1], currentChar)
			if previousLength > 0 {
				currentIndex += previousLength
			} else {
				if !isSymbol(currentChar) {
					break
				}
				currentIndex++
			}
		}
	}
	return currentIndex, tldIndex
}
