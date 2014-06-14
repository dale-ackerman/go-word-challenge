package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	combos = make(map[string]bool)
)

func mapwords() map[string]bool {
	wmap := make(map[string]bool)
	dfile, err := ioutil.ReadFile("/usr/share/dict/words")
	if err != nil {
		fmt.Println(err)
	}
	dwords := strings.Split(string(dfile), "\n")
	for _, v := range dwords {
		wmap[strings.ToLower(v)] = true
	}
	return wmap
}

func permute(prefix, str string, le int) {
	n := len(str)
	if n == 0 || len(prefix) == le {
		combos[prefix] = true
	} else {
		for i := 0; i < n; i++ {
			permute(prefix+string(str[i]), str[0:i]+str[i+1:n], le)
		}
	}
}

func main() {
	counter := 0
	for i := 0; i < len(os.Args[1]); i++ {
		permute("", os.Args[1], i)
	}
	wmap := mapwords()
	for k, _ := range combos {
		if wmap[k] {
			// fmt.Println(k)
			counter++
		}
	}
	fmt.Println(counter)
}
