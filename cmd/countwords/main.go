// Courtesy of Ben Hoyt (https://github.com/benhoyt)
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fredrb/htable"
)

type S = htable.StringKey

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	counts := htable.New()
	var uniques []string
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		pvalue, ok := counts.Get(S(word))
		if !ok {
			counts.Set(S(word), 1)
			uniques = append(uniques, word)
		} else {
			counts.Set(S(word), pvalue.(int)+1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sort.Slice(uniques, func(i, j int) bool {
		v, _ := counts.Get(S(uniques[i]))
		iCount := v.(int)
		v, _ = counts.Get(S(uniques[j]))
		jCount := v.(int)
		return iCount > jCount
	})

	for _, word := range uniques {
		v, _ := counts.Get(S(word))
		count := v.(int)
		fmt.Println(word, count)
	}

	fmt.Fprintln(os.Stderr, "Hash table dump:")
	fmt.Fprintln(os.Stderr, "----------------")
	counts.Dump(os.Stderr)
}
