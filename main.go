package main

import (
	"bufio"
	"flag"
	"os"
	"time"
)

func main() {
	// --seed를 주면 난수가 고정되어 게임이 완전히 결정론적이 된다(테스트/비교용).
	// 주지 않으면 원본처럼 현재 시각을 시드로 쓴다.
	seed := flag.Int64("seed", -1, "fixed RNG seed (default: time-based, like the original)")
	flag.Parse()

	g := &Game{
		in:  bufio.NewReader(os.Stdin),
		out: os.Stdout,
	}
	if *seed < 0 {
		g.rx = (int(time.Now().Unix()) & 0xfffff) | 1
	} else {
		g.rx = (int(*seed) & 0xfffff) | 1
	}

	g.Run()
}
