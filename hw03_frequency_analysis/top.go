package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type WordsCount struct {
	word string
	cnt  uint
}

func cleanWord(word string) string {
	// set word to lowercase and trim symbols
	word = strings.ToLower(word)
	return strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func Top10(input string) []string {
	const Limit = 10

	wordsCntMap := make(map[string]uint)
	// count words
	for _, word := range strings.Fields(input) {
		word = cleanWord(word)
		if word == "" {
			// skip words containing only symbols
			continue
		}
		cnt, inMap := wordsCntMap[word]
		if inMap {
			wordsCntMap[word] = cnt + 1
		} else {
			wordsCntMap[word] = 1
		}
	}

	// store counts to a slice of structs for further sorting
	wordsCntStructs := make([]WordsCount, 0, len(wordsCntMap))
	for word, cnt := range wordsCntMap {
		wordsCntStructs = append(wordsCntStructs, WordsCount{word, cnt})
	}
	// sort by cnt DESC and word ASC
	sort.Slice(wordsCntStructs, func(i, j int) bool {
		if wordsCntStructs[i].cnt != wordsCntStructs[j].cnt {
			return wordsCntStructs[i].cnt > wordsCntStructs[j].cnt
		}
		return wordsCntStructs[i].word < wordsCntStructs[j].word
	})

	// store top words to the result slice
	result := make([]string, 0, Limit)
	for ind, wordsCnt := range wordsCntStructs {
		result = append(result, wordsCnt.word)
		if ind == Limit-1 {
			break
		}
	}

	return result
}
