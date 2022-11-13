package main

import (
	"fmt"
	"time"
)

func LastSerial(path string) (str string, err error) {

	return "", nil
}

func findSerial(gool uint) (uint, error) {

	staps := 0
	base := uint(0)
	big := uint(1)
	smal := uint(0)

	for staps = 0; staps <= 300; staps++ {

		if base < gool {
			smal = base
			base = big
			big = base * 2
			if big >= (2305843009213693000) {
				// TODO find bigger then this number
				return 0, fmt.Errorf("big overflow at stap: %d\n", staps)
			}
		}

		if base > gool {
			big = base
			base = (smal + big) / 2
		}

		if base == gool {
			if base+1 > gool {
				println("gool", base)
				println("staps", staps)
				break
			}
		}
	}

	return base, nil
}

func genSerial() string {
	src := []string{"a", "b", "b", "c", "d", "e", "f", "j", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	res := ""

	for j := 0; j < 1000000; j++ {
		res = ""
		for i := 0; i < 5; i++ {
			res += string(src[i])
		}
	}

	fmt.Println("res : ", res)
	return res
}

func goLoop() {

	var mems = []string{"dog", "cat", "mouse", "koko", "bebe", "jojo", "haha", "jiji", "foo", "bar", "bax", "hik", "jik", "ors", "nos", "ren"}

	var bots = []string{"bot_1", "bot_2", "bot_3", "bot_4", "bot_5", "bot_6", "bot-7"}

	i := 0
	lbots := len(bots)
	for _, mem := range mems {
		go func(mem string, i int) {
			fmt.Printf("%s kik %s\n", bots[i], mem)
		}(mem, i)
		i++
		if lbots == i {
			i = 0
		}
	}
	time.Sleep(time.Millisecond * 10)
}
