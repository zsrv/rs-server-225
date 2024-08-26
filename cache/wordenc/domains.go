package wordenc

func filterDomains(wf *WordFilter, chars []rune) {
	ampersand := make([]rune, len(chars))
	copy(ampersand, chars)

	period := make([]rune, len(chars))
	copy(period, chars)

	filterBadCombinations(wf, nil, ampersand, []rune{'(', 'a', ')'})
	filterBadCombinations(wf, nil, period, []rune{'d', 'o', 't'})

	for index := len(wf.Domains) - 1; index >= 0; index-- {
		filterDomain(period, ampersand, wf.Domains[index], chars)
	}
}

func getEmulatedDomainCharLen(nextChar rune, domainChar rune, currentChar rune) int {
	if domainChar == currentChar {
		return 1
	} else if domainChar == 'o' && currentChar == '0' {
		return 1
	} else if domainChar == 'o' && currentChar == '(' && nextChar == ')' {
		return 2
	} else if domainChar == 'c' && (currentChar == '(' || currentChar == '<' || currentChar == '[') {
		return 1
	} else if domainChar == 'e' && currentChar == '€' {
		return 1
	} else if domainChar == 's' && currentChar == '$' {
		return 1
	} else if domainChar == 'l' && currentChar == 'i' {
		return 1
	}
	return 0
}

func filterDomain(period []rune, ampersand []rune, domain []rune, chars []rune) {
	domainLength := len(domain)
	charsLength := len(chars)
	for index := 0; index <= charsLength-domainLength; index++ {
		matched, currentIndex := findMatchingDomain(index, domain, chars)
		if !matched {
			continue
		}

		ampersandStatus := prefixSymbolStatus(index, chars, 3, ampersand, []rune{'@'})
		periodStatus := suffixSymbolStatus(currentIndex-1, chars, 3, period, []rune{'.', ','})
		shouldFilter := ampersandStatus > 2 || periodStatus > 2
		if !shouldFilter {
			continue
		}
		maskChars(index, currentIndex, chars)
	}
}

func findMatchingDomain(startIndex int, domain []rune, chars []rune) (bool, int) {
	domainLength := len(domain)
	currentIndex := startIndex
	domainIndex := 0

	for currentIndex < len(chars) && domainIndex < domainLength {
		currentChar := chars[currentIndex]

		var nextChar rune
		if currentIndex+1 < len(chars) {
			nextChar = chars[currentIndex+1]
		} else {
			nextChar = '\u0000'
		}

		currentLength := getEmulatedDomainCharLen(nextChar, domain[domainIndex], currentChar)
		if currentLength > 0 {
			currentIndex += currentLength
			domainIndex++
		} else {
			if domainIndex == 0 {
				break
			}

			previousLength := getEmulatedDomainCharLen(nextChar, domain[domainIndex-1], currentChar)
			if previousLength > 0 {
				currentIndex += previousLength
				if domainIndex == 1 {
					startIndex++
				}
			} else {
				if domainIndex >= domainLength || !isSymbol(currentChar) {
					break
				}
				currentIndex++
			}
		}
	}
	return domainIndex >= domainLength, currentIndex
}
