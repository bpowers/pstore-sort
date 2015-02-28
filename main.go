package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Entry struct {
	Time float64
	Line string
}

func NewEntry(line string) (e *Entry) {
	e = &Entry{Line: line}

	s := strings.IndexByte(line, '[')
	if s < 0 {
		log.Printf("couldn't find [ in '%s'", line)
		return
	}
	s++ // skip the '['
	l := strings.IndexByte(line[s:], ']')
	if l < 0 {
		log.Printf("couldn't find ] in '%s'", line)
		return
	}
	t, err := strconv.ParseFloat(strings.TrimSpace(line[s:s+l]), 64)
	if err != nil {
		log.Printf("ParseFloat: %s", err)
		return
	}
	e.Time = t

	return e
}

type ByTime []*Entry

func (a ByTime) Len() int      { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool {
	return a[i].Time < a[j].Time
}

func shouldSkip(l string) bool {
	return strings.HasPrefix(l, "Oops#") ||
		strings.HasPrefix(l, "Panic#") || len(l) == 0
}

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("ReadAll(stdin): %s", err)
		os.Exit(1)
	}

	rawLines := strings.Split(string(buf), "\n")

	entries := make([]*Entry, 0, len(rawLines))
	for _, l := range rawLines {
		if shouldSkip(l) {
			continue
		}
		entries = append(entries, NewEntry(l))
	}

	sort.Sort(ByTime(entries))

	// sometimes we get dups.
	for i := 1; i < len(entries); i++ {
		if entries[i-1].Line == entries[i].Line {
			entries = append(entries[:i], entries[i+1:]...)
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for _, l := range entries {
		out.WriteString(l.Line)
		out.WriteByte('\n')
	}
}
