package main

import "testing"

func TestCaveData(t *testing.T) {
	// 출발지 road의 설명이 그 유명한 첫 문장(한글)인지
	want := "넌 작은 벽돌 건물 앞, 길 끝에 서 있어.  주위는 온통 숲이야.  건물에서\n" +
		"작은 시냇물이 흘러나와 도랑을 따라 내려가."
	if caveLongDesc[road] != want {
		t.Errorf("road long_desc 불일치:\n%q", caveLongDesc[road])
	}
	if caveShortDesc[road] != "다시 길 끝에 왔어." {
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
