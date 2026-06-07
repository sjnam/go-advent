package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runGame은 주어진 입력으로 게임을 시드 42로 돌리고 출력을 돌려준다.
// 시드 42는 시드를 고정한 C reference(advent-fixed)와 동일하다.
func runGame(input string) string {
	var out bytes.Buffer
	g := &Game{
		in:  bufio.NewReader(strings.NewReader(input)),
		out: &out,
		rx:  (42 & 0xfffff) | 1,
	}
	g.Run()
	return out.String()
}

// TestGolden은 testdata/*.in 각각을 게임에 넣고, 출력이 같은 이름의
// .golden(시드 고정 C reference로 캡처)과 byte-identical인지 확인한다.
// -update 플래그로 golden을 갱신할 수 있다.
func TestGolden(t *testing.T) {
	ins, err := filepath.Glob("testdata/*.in")
	if err != nil {
		t.Fatal(err)
	}
	if len(ins) == 0 {
		t.Skip("testdata/*.in 없음")
	}
	for _, in := range ins {
		in := in
		name := strings.TrimSuffix(filepath.Base(in), ".in")
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(in)
			if err != nil {
				t.Fatal(err)
			}
			got := runGame(string(input))
			goldenPath := strings.TrimSuffix(in, ".in") + ".golden"
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("golden 읽기 실패: %v", err)
			}
			if got != string(want) {
				t.Errorf("출력이 golden과 다름.\n--- got ---\n%s\n--- want ---\n%s", got, want)
			}
		})
	}
}
