package main

import "testing"

func TestObjectData(t *testing.T) {
	g := &Game{}
	g.loadObjects()
	// 보물은 prop=-1로 시작 (advent.w:2230)
	if g.prop[GOLD] != -1 {
		t.Errorf("prop[GOLD]=%d, want -1 (보물)", g.prop[GOLD])
	}
	if g.prop[KEYS] != 0 {
		t.Errorf("prop[KEYS]=%d, want 0 (비보물)", g.prop[KEYS])
	}
	// 집 안에 키/램프/음식/병이 있어야 함 (well house, house=3)
	g.loc = house
	for _, o := range []object{KEYS, LAMP, FOOD, BOTTLE} {
		if !g.here(o) {
			t.Errorf("%s가 house에 없음 (place=%d)", objName[o], g.place[o])
		}
	}
	// SNAKE는 움직이지 않는 그룹 객체 (base[SNAKE]==SNAKE)
	if objBase[SNAKE] != SNAKE {
		t.Errorf("base[SNAKE]=%d, want SNAKE", objBase[SNAKE])
	}
	// KEYS는 움직임 (base==NOTHING)
	if objBase[KEYS] != NOTHING {
		t.Errorf("base[KEYS]=%d, want NOTHING", objBase[KEYS])
	}
	// carry/drop 동작
	h0 := g.holding
	g.carry(KEYS)
	if !g.toting(KEYS) || g.holding != h0+1 {
		t.Errorf("carry(KEYS) 실패: toting=%v holding=%d", g.toting(KEYS), g.holding)
	}
	g.drop(KEYS, road)
	if g.toting(KEYS) || g.place[KEYS] != road {
		t.Errorf("drop(KEYS,road) 실패: place=%d", g.place[KEYS])
	}
}
