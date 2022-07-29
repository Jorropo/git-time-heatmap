package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	err := mainRet()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func mainRet() error {
	r, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("opening repo: %w", err)
	}

	itr, err := r.Log(&git.LogOptions{})
	if err != nil {
		return fmt.Errorf("getting log iterator: %w", err)
	}

	var boxes [7][24]uint

	err = itr.ForEach(func(c *object.Commit) error {
		t := c.Author.When
		boxes[int(t.Weekday())][t.Hour()] += 1
		return nil
	})
	if err != nil {
		return fmt.Errorf("logging: %w", err)
	}

	boxesStr := [8][24]string{{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}}
	maxHitSize := 2
	for i, hours := range boxes {
		for j, h := range hours {
			r := strconv.FormatUint(uint64(h), 10)
			if len(r) > maxHitSize {
				maxHitSize = len(r)
			}
			boxesStr[i+1][j] = r
		}
	}

	for i, hours := range boxesStr {
		for j, h := range hours {
			for len(h) < maxHitSize {
				h += " "
			}
			boxesStr[i][j] = h
		}
	}

	for i, hours := range boxesStr {
		var prefix string
		if i == 0 {
			prefix = "     "
		} else {
			prefix = time.Weekday(i - 1).String()[:3] + ": "
		}
		fmt.Println(prefix + strings.Join(hours[:], " | "))
	}

	return nil
}
