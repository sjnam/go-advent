package main

import (
	"bufio"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveRestore(t *testing.T) {
	g := &Game{in: bufio.NewReader(strings.NewReader("")), out: io.Discard, rx: 42}
	g.buildVocabulary()
	g.loadObjects()
	// 상태를 좀 바꿔 둔다
	g.loc = hmk
	g.turns = 17
	g.prop[LAMP] = 1
	g.carry(KEYS)
	g.dflag = 2
	g.clock1 = 5

	path := filepath.Join(t.TempDir(), "t.save")
	if err := g.save(path); err != nil {
		t.Fatal(err)
	}
	// 다른 게임에 불러오기
	g2 := &Game{in: bufio.NewReader(strings.NewReader("")), out: io.Discard}
	g2.buildVocabulary()
	g2.loadObjects()
	if err := g2.load(path); err != nil {
		t.Fatal(err)
	}
	if g2.loc != hmk || g2.turns != 17 || g2.prop[LAMP] != 1 ||
		!g2.toting(KEYS) || g2.dflag != 2 || g2.clock1 != 5 || g2.rx != g.rx {
		t.Errorf("복원 상태 불일치: loc=%d turns=%d lamp=%d keys=%v dflag=%d clock1=%d",
			g2.loc, g2.turns, g2.prop[LAMP], g2.toting(KEYS), g2.dflag, g2.clock1)
	}
}
