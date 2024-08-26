package wordenc

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/zsrv/rs-server-225/jagex2/io"
	"github.com/zsrv/rs-server-225/jagex2/packet"
)

type WordFilter struct {
	BadWords        [][]rune
	BadCombinations [][][2]int
	Domains         [][]rune
	Fragments       []int
	TLDs            [][]rune
	TLDTypes        []uint8
}

func LoadWordFilter(dir string) (*WordFilter, error) {
	jf, err := io.LoadJagfile(filepath.Join(dir, "client", "wordenc"))
	if err != nil {
		return nil, err
	}

	w := &WordFilter{}

	if err := w.readAll(jf); err != nil {
		return nil, err
	}

	return w, err
}

func (w *WordFilter) readAll(jf *io.Jagfile) error {
	fragments, err := jf.Read("fragmentsenc.txt")
	if err != nil {
		return err
	}

	bad, err := jf.Read("badenc.txt")
	if err != nil {
		return err
	}

	domain, err := jf.Read("domainenc.txt")
	if err != nil {
		return err
	}

	tld, err := jf.Read("tldlist.txt")
	if err != nil {
		return err
	}

	w.readBadWords(bad)
	w.readDomains(domain)
	w.readFragments(fragments)
	w.readTLD(tld)

	return nil
}

func (w *WordFilter) readTLD(buf *packet.Packet) {
	count := buf.G4()

	w.TLDs = make([][]rune, count)
	w.TLDTypes = make([]uint8, count)

	for i := range count {
		w.TLDTypes[i] = buf.G1()
		w.TLDs[i] = make([]rune, buf.G1())

		for j := range w.TLDs[i] {
			w.TLDs[i][j] = rune(buf.G1())
		}
	}
}

func (w *WordFilter) readBadWords(buf *packet.Packet) {
	count := buf.G4()

	badWords := make([][]rune, count)
	badCombinations := make([][][2]int, count)

	for i := range count {
		badWords[i] = make([]rune, buf.G1())
		for j := range badWords[i] {
			badWords[i][j] = rune(buf.G1())
		}

		combination := make([][2]int, buf.G1())
		for j := range combination {
			combination[j][0] = int(buf.G1())
			combination[j][1] = int(buf.G1())
		}
		if len(combination) > 0 {
			badCombinations[i] = combination
		}
	}

	w.BadWords = badWords
	w.BadCombinations = badCombinations
}

func (w *WordFilter) readDomains(buf *packet.Packet) {
	count := buf.G4()
	domains := make([][]rune, count)
	for i := range count {
		domains[i] = make([]rune, buf.G1())
		for j := range domains[i] {
			domains[i][j] = rune(buf.G1())
		}
	}
	w.Domains = domains
}

func (w *WordFilter) readFragments(buf *packet.Packet) {
	count := buf.G4()
	w.Fragments = make([]int, count)

	for i := range w.Fragments {
		w.Fragments[i] = int(buf.G2())
	}
}

func filterCharacters(chars []rune) {
	pos := 0
	for i := range chars {
		if isCharacterAllowed(chars[i]) {
			chars[pos] = chars[i]
		} else {
			chars[pos] = ' '
		}
		if pos == 0 || chars[pos] != ' ' || chars[pos-1] != ' ' {
			pos++
		}
	}
	for i := pos; i < len(chars); i++ {
		chars[i] = ' '
	}
}

func isCharacterAllowed(c rune) bool {
	return (c >= ' ' && c <= '\u007f') || c == ' ' || c == '\n' || c == '\t' || c == '£' || c == '€'
}

func filter(wf *WordFilter, input string) string {
	outputPre := []rune(input)
	filterCharacters(outputPre)
	trimmed := strings.TrimSpace(string(outputPre))
	lowercase := strings.ToLower(trimmed)
	output := []rune(lowercase)
	filterTLD(wf, output)
	filterBadWords(wf, output)
	filterDomains(wf, output)
	filterFragments(output)

	whitelist := []string{"cook", "cook's", "cooks", "seeks", "sheet"}

	for index := range whitelist {
		offset := -1
		for offset = strings.Index(lowercase[offset+1:], whitelist[index]); offset != -1; {
			whitelisted := []rune(whitelist[index])
			for charIndex := range whitelisted {
				output[charIndex+offset] = whitelisted[charIndex]
			}
		}
	}
	replaceUppercases(output, []rune(trimmed))
	formatUppercases(output)
	return strings.TrimSpace(string(output))
}

func isSymbol(r rune) bool {
	return !isAlpha(r) && !isNumber(r)
}

func isLowercaseAlpha(r rune) bool {
	if isLowercase(r) {
		return r == 'v' || r == 'x' || r == 'j' || r == 'q' || r == 'z'
	}
	return true
}

func isAlpha(r rune) bool {
	return isLowercase(r) || isUppercaseAlpha(r)
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLowercase(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isUppercaseAlpha(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isNumericalChars(s []rune) bool {
	for _, r := range s {
		if !isNumber(r) && r != '\u0000' {
			return false
		}
	}
	return true
}

func maskChars(offset int, length int, chars []rune) {
	for i := offset; i < length; i++ {
		chars[i] = '*'
	}
}

func maskedCountBackwards(chars []rune, offset int) int {
	count := 0
	for i := offset - 1; i >= 0 && isSymbol(chars[i]); i-- {
		if chars[i] == '*' {
			count++
		}
	}
	return count
}

func maskedCountForwards(chars []rune, offset int) int {
	count := 0
	for i := offset + 1; i < len(chars) && isSymbol(chars[i]); i++ {
		if chars[i] == '*' {
			count++
		}
	}
	return count
}

func maskedCharsStatus(chars []rune, filtered []rune, offset int, length int, prefix bool) int {
	var count int
	switch prefix {
	case true:
		count = maskedCountBackwards(filtered, offset)
	case false:
		count = maskedCountForwards(filtered, offset)
	}

	if count >= length {
		return 4
	}

	switch prefix {
	case true:
		if isSymbol(chars[offset-1]) {
			return 1
		}
	case false:
		if isSymbol(chars[offset+1]) {
			return 1
		}
	}

	return 0
}

func prefixSymbolStatus(offset int, chars []rune, length int, symbolChars []rune, symbols []rune) int {
	if offset == 0 {
		return 2
	}

	for i := offset - 1; i >= 0 && isSymbol(chars[i]); i-- {
		if slices.Contains(symbols, chars[i]) {
			return 3
		}
	}

	return maskedCharsStatus(chars, symbolChars, offset, length, true)
}

func suffixSymbolStatus(offset int, chars []rune, length int, symbolChars []rune, symbols []rune) int {
	if offset+1 == len(chars) {
		return 2
	}

	for i := offset + 1; i < len(chars) && isSymbol(chars[i]); i++ {
		if slices.Contains(symbols, chars[i]) {
			return 3
		}
	}

	return maskedCharsStatus(chars, symbolChars, offset, length, false)
}

func replaceUppercases(chars []rune, comparison []rune) {
	for i := range comparison {
		if chars[i] != '*' && isUppercaseAlpha(comparison[i]) {
			chars[i] = comparison[i]
		}
	}
}

func formatUppercases(chars []rune) {
	flagged := true
	for i := range chars {
		char := chars[i]
		if !isAlpha(char) {
			flagged = true
		} else if flagged {
			if isLowercase(char) {
				flagged = false
			}
		} else if isUppercaseAlpha(char) {
			chars[i] = char + 'a' - 65
		}
	}
}
