package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"
)

var index64 = [64]int{
	0, 1, 48, 2, 57, 49, 28, 3,
	61, 58, 50, 42, 38, 29, 17, 4,
	62, 55, 59, 36, 53, 51, 43, 22,
	45, 39, 33, 30, 24, 18, 12, 5,
	63, 47, 56, 27, 60, 41, 37, 16,
	54, 35, 52, 21, 44, 32, 23, 11,
	46, 26, 40, 15, 34, 20, 31, 10,
	25, 14, 19, 9, 13, 8, 7, 6,
}

func readDict(path string, maxLen int) (map[string]bool, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Either just use this:
	// (much shorter but also slower)
	/*
		   dict := make(map[string]bool)
		   b := bytes.Split(buf, []byte("\n"))
		   for i := range b {
			   dict[string(bytes.ToLower(b[i]))] = true
		   }
		   return dict, nil
	*/

	// Or the following:

	indexBitmap := make([]uint64, (len(buf)+63)/64)

	count := 0
	for i := 0; i < len(buf); i++ {
		// Find the next line break
		pos := bytes.IndexByte(buf[i:], '\n')
		if pos < 0 {
			break
		}
		if pos <= maxLen {
			if pos < 0 {
				break
			}
			count++
		}

		i += pos

		// Save pos in bitmap
		indexBitmap[i/64] |= 1 << (uint(i) % 64)
	}

	dict := make(map[string]bool, count)
	offset := 0
	for i, bm := range indexBitmap {
		for bm > 0 {
			// Get the least significant set bit (=1)
			lsb := bm & -bm

			// Remove the LSB from the bitmap
			bm ^= lsb

			// Get index from LSB with 64bit DeBruijn multiplication table
			// This could be done even faster using Assembly:
			// http://en.wikipedia.org/wiki/Find_first_set
			index := index64[(lsb*0x03f79d71b4cb0a89)>>58] + i*64

			if index-offset <= maxLen {
				dict[string(bytes.ToLower(buf[offset:index]))] = true
			}

			offset = index + 1
		}
	}
	return dict, nil
}

func buildPermSubseq(dict map[string]bool, str string) (p []string) {
	var charPermuted [256]bool
	var wg sync.WaitGroup
	var mu sync.Mutex
	for j := 0; j < 10; j++ {
		if !charPermuted[str[j]] {
			nStr := str[j : j+1]

			if dict[nStr] {
				mu.Lock()
				p = append(p, nStr)
				mu.Unlock()
			}

			postfix := str[:j] + str[j+1:]

			wg.Add(1)
			go func() {
				ret := permSubseq(dict, nStr, postfix, 1)
				mu.Lock()
				p = append(p, ret...)
				mu.Unlock()
				wg.Done()
			}()

			charPermuted[str[j]] = true
		}

	}
	wg.Wait()
	return p
}

func permSubseq(dict map[string]bool, prefix, str string, i int) (p []string) {
	var charPermuted [256]bool
	n := len(prefix) + len(str)
	for j := 0; j < len(str); j++ {
		if !charPermuted[str[j]] {
			nStr := prefix + str[j:j+1]

			if dict[nStr] {
				p = append(p, nStr)
			}

			if i+1 < n {
				p = append(p, permSubseq(dict, nStr, str[:j]+str[j+1:], i+1)...)
			}

			charPermuted[str[j]] = true
		}
	}
	return p
}

func main() {
	str := "racecardss"

	dict, err := readDict("/usr/share/dict/words", len(str))
	if err != nil {
		return
	}

	p := buildPermSubseq(dict, str)

	fmt.Println(p)
	fmt.Println(len(p))
}
