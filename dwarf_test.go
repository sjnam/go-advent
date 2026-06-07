package main

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestDwarfActivation(t *testing.T) {
	g := &Game{in: bufio.NewReader(strings.NewReader("")), out: io.Discard, rx: 42}
	g.buildVocabulary()
	g.loadObjects()
	g.dloc = dwarfStart
	g.knifeLoc = -1
	g.loc = hmk // Hall of Mountain King(깊은 동굴)
	for i := 0; i < 300 && g.dflag < 2; i++ {
		g.moveDwarves()
	}
	if g.dflag < 1 {
		t.Fatal("동굴 깊은 곳에서도 dflag가 0 (난쟁이 비활성)")
	}
	t.Logf("dflag=%d 도달", g.dflag)
	// 같은 방의 난쟁이를 dwarf()가 알아채는지
	g.dflag = 2
	g.dloc[1] = hmk
	if !g.dwarf() {
		t.Error("같은 방의 난쟁이를 dwarf()가 못 알아챔")
	}
}
