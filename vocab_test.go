package main

import (
	"bufio"
	"strings"
	"testing"
)

func newTestGame() *Game {
	g := &Game{in: bufio.NewReader(strings.NewReader("")), rx: 43}
	g.buildVocabulary()
	return g
}

func TestVocabLookup(t *testing.T) {
	g := newTestGame()
	cases := []struct {
		word string
		typ  wordtype
		mean int
		ok   bool
	}{
		{"take", actionType, int(TAKE), true},
		{"xyzzy", motionType, int(XYZZY), true},
		{"nugge", objectType, int(GOLD), true},
		{"northeast", motionType, int(NE), false}, // 5자 절단→"north"=N, NE 아님
		{"north", motionType, int(N), true},
		{"plugh", motionType, int(PLUGH), true},
		{"help", messageType, 1, true},
		{"abra", messageType, 0, true},
		{"swim", messageType, 12, true},
		{"tickl", noType, 0, false}, // 모르는 단어
		{"h2o", objectType, int(WATER), true},
	}
	for _, c := range cases {
		e, ok := g.lookup(c.word)
		if c.word == "northeast" {
			// "north"로 잘려 N으로 인식되는지 확인
			if !ok || e.typ != motionType || e.meaning != int(N) {
				t.Errorf("lookup(northeast) 절단 결과 = %+v,%v; want N", e, ok)
			}
			continue
		}
		if ok != c.ok {
			t.Errorf("lookup(%q) ok=%v, want %v", c.word, ok, c.ok)
			continue
		}
		if ok && (e.typ != c.typ || e.meaning != c.mean) {
			t.Errorf("lookup(%q) = {typ:%d mean:%d}, want {typ:%d mean:%d}", c.word, e.typ, e.meaning, c.typ, c.mean)
		}
	}
}

func TestDefaultMsgChaining(t *testing.T) {
	g := newTestGame()
	if g.defaultMsg[CLOSE] != g.defaultMsg[OPEN] || g.defaultMsg[OPEN] == "" {
		t.Error("CLOSE는 OPEN 메시지를 공유해야 함")
	}
	if g.defaultMsg[QUIT] != "Eh?" {
		t.Errorf("QUIT 기본메시지=%q, want SCORE의 Eh?", g.defaultMsg[QUIT])
	}
	if g.ok() != "OK." {
		t.Errorf("ok()=%q, want OK.", g.ok())
	}
}
