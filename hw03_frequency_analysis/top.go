package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type Pair struct {
	Value string
	Count int
}

func Top10(str string) []string {
	reg := regexp.MustCompile(`(?i)[,.а-яё-]+-?[а-яё]*[^, .!?\n\r]|[а-яё]`)
	// найти все слова
	matches := reg.FindAllString(str, -1)
	if len(matches) == 0 {
		return nil
	}
	// найти все повторения
	counts := make(map[string]int)
	for _, v := range matches {
		counts[strings.ToLower(v)]++
	}
	// отсортировать результат по количеству
	forSort := make([]Pair, len(counts))
	i := 0
	for k, v := range counts {
		forSort[i] = Pair{
			Value: k,
			Count: v,
		}
		i++
	}
	sort.SliceStable(forSort, func(i, j int) bool {
		if forSort[i].Count == forSort[j].Count {
			return strings.Compare(forSort[i].Value, forSort[j].Value) == -1
		}
		return forSort[i].Count > forSort[j].Count
	})
	// вернуть первые 10
	return firstTenStrings(forSort)
}

func firstTenStrings(list []Pair) []string {
	result := make([]string, len(list))
	for i := range list {
		result[i] = list[i].Value
	}
	if len(result) < 10 {
		return result
	}
	return result[0:10]
}
