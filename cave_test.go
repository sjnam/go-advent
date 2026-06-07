package main

import "testing"

func TestCaveData(t *testing.T) {
	// 출발지 road의 설명이 그 유명한 첫 문장인지
	want := "You are standing at the end of a road before a small brick building.\n" +
		"Around you is a forest.  A small stream flows out of the building and\n" +
		"down a gully."
	if caveLongDesc[road] != want {
		t.Errorf("road long_desc 불일치:\n%q", caveLongDesc[road])
	}
	if caveShortDesc[road] != "You're at end of road again." {
		t.Errorf("road short_desc=%q", caveShortDesc[road])
	}
	if caveFlags[road] != lighted+liquid {
		t.Errorf("road flags=%d, want %d", caveFlags[road], lighted+liquid)
	}
	// 이동표 개수와 경계 일관성
	if len(caveTravels) != 740 {
		t.Errorf("travels 개수=%d, want 740", len(caveTravels))
	}
	if caveStart[maxLoc+1] != len(caveTravels) {
		t.Errorf("start[maxLoc+1]=%d, want %d", caveStart[maxLoc+1], len(caveTravels))
	}
	// start는 limbo(0)만 -1이어야 하고 나머지는 단조 증가
	prev := -1
	for l := road; l <= maxLoc; l++ {
		s := caveStart[l]
		if s < 0 {
			t.Errorf("start[%d] = -1 (limbo 외엔 없어야 함)", l)
		}
		if s < prev {
			t.Errorf("start[%d]=%d < 이전 %d (단조 증가 위반)", l, s, prev)
		}
		prev = s
	}
	// road의 첫 명령: W -> hill (advent.w:863)
	first := caveTravels[caveStart[road]]
	if first.mot != W || first.dest != hill {
		t.Errorf("road 첫 명령 = %+v, want {W,0,hill}", first)
	}
}
