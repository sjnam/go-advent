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
		{"가져가", actionType, int(TAKE), true},
		{"금괴", objectType, int(GOLD), true},
		{"북", motionType, int(N), true},
		{"북동", motionType, int(NE), true},       // 한글은 전체 단어 매칭(절단 없음)
		{"xyzzy", motionType, int(XYZZY), true}, // 주문은 영어 유지
		{"plugh", motionType, int(PLUGH), true},
		{"도움말", messageType, 1, true},
		{"수리수리", messageType, 0, true},
		{"수영", messageType, 12, true},
		{"틱탁", noType, 0, false}, // 모르는 단어
		{"물", objectType, int(WATER), true},
	}
	for _, c := range cases {
		e, ok := g.lookup(c.word)
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
	if g.defaultMsg[QUIT] != "응?" {
		t.Errorf("QUIT 기본메시지=%q, want SCORE의 \"응?\"", g.defaultMsg[QUIT])
	}
	if g.ok() != "알았어." {
		t.Errorf("ok()=%q, want \"알았어.\"", g.ok())
	}
}
