package main

import (
	"encoding/gob"
	"os"
)

// 게임 저장과 불러오기. 원작에는 없던 기능으로, 이 포팅에서 더했다.
//
// 게임의 영속 상태는 거의 다 Game 구조체에 모여 있어서, 그 상태만 파일에
// 적었다 되읽으면 그대로 저장/복원이 된다. 다만 비공개 필드는 encoding/gob이
// 직렬화하지 못하므로, 내보낸(대문자) 필드를 가진 saveState로 옮겨 담는다.
// 어휘·메시지(words/defaultMsg/message)는 시작 시 buildVocabulary로 다시
// 만들므로 저장하지 않는다. 턴 사이에만 의미 있는 임시 값(mot/verb/obj 등)도
// 저장 시점(턴 경계)에는 무의미하므로 제외한다.

const saveFile = "advent.save" // 고정 저장 파일명

// saveState는 한 판의 영속 상태를 담는 직렬화용 스냅샷이다.
type saveState struct {
	Rx int

	Oldoldloc, Oldloc, Loc, Newloc location
	Visits                         [maxLoc + 1]int

	Place   [maxObj + 1]location
	Prop    [maxObj + 1]int
	Link    [maxObj + 1]object
	First   [maxLoc + 1]object
	Holding int
	ObjBase map[object]object // 전역이지만 게임 중 바뀌므로 함께 저장

	Turns         int
	Interval      int
	Tally         int
	LostTreasures int

	Dflag    int
	Dkill    int
	Dloc     [nd + 1]location
	Odloc    [nd + 1]location
	Dseen    [nd + 1]bool
	KnifeLoc int

	Clock1, Clock2 int
	Panicked       bool
	Warned         bool
	DeathCount     int
	Bonus          int
	Closed         bool
	Foobar         int

	Hinted    [nHints]bool
	HintCount [nHints]int
	Limit     int
	GaveUp    bool
}

// snapshot은 현재 게임 상태를 saveState로 모은다.
func (g *Game) snapshot() saveState {
	return saveState{
		Rx:            g.rx,
		Oldoldloc:     g.oldoldloc,
		Oldloc:        g.oldloc,
		Loc:           g.loc,
		Newloc:        g.newloc,
		Visits:        g.visits,
		Place:         g.place,
		Prop:          g.prop,
		Link:          g.link,
		First:         g.first,
		Holding:       g.holding,
		ObjBase:       objBase,
		Turns:         g.turns,
		Interval:      g.interval,
		Tally:         g.tally,
		LostTreasures: g.lostTreasures,
		Dflag:         g.dflag,
		Dkill:         g.dkill,
		Dloc:          g.dloc,
		Odloc:         g.odloc,
		Dseen:         g.dseen,
		KnifeLoc:      g.knifeLoc,
		Clock1:        g.clock1,
		Clock2:        g.clock2,
		Panicked:      g.panicked,
		Warned:        g.warned,
		DeathCount:    g.deathCount,
		Bonus:         g.bonus,
		Closed:        g.closed,
		Foobar:        g.foobar,
		Hinted:        g.hinted,
		HintCount:     g.hintCount,
		Limit:         g.limit,
		GaveUp:        g.gaveUp,
	}
}

// restore는 saveState의 값을 현재 게임에 되돌려 놓는다.
func (g *Game) restore(s saveState) {
	g.rx = s.Rx
	g.oldoldloc, g.oldloc, g.loc, g.newloc = s.Oldoldloc, s.Oldloc, s.Loc, s.Newloc
	g.visits = s.Visits
	g.place, g.prop, g.link, g.first = s.Place, s.Prop, s.Link, s.First
	g.holding = s.Holding
	objBase = s.ObjBase
	g.turns = s.Turns
	g.interval, g.tally, g.lostTreasures = s.Interval, s.Tally, s.LostTreasures
	g.dflag, g.dkill = s.Dflag, s.Dkill
	g.dloc, g.odloc, g.dseen = s.Dloc, s.Odloc, s.Dseen
	g.knifeLoc = s.KnifeLoc
	g.clock1, g.clock2 = s.Clock1, s.Clock2
	g.panicked, g.warned = s.Panicked, s.Warned
	g.deathCount, g.bonus = s.DeathCount, s.Bonus
	g.closed, g.foobar = s.Closed, s.Foobar
	g.hinted, g.hintCount = s.Hinted, s.HintCount
	g.limit, g.gaveUp = s.Limit, s.GaveUp
}

// save는 현재 상태를 파일에 기록한다.
func (g *Game) save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(g.snapshot())
}

// load는 파일에서 상태를 읽어 게임에 되돌린다.
func (g *Game) load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var s saveState
	if err := gob.NewDecoder(f).Decode(&s); err != nil {
		return err
	}
	g.restore(s)
	return nil
}
