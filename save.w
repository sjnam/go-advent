@* 저장과 불러오기.
원작에 없던 단 하나의 새 기능이다. 1977년의 모험가는 전산실 단말
앞에서 한 판을 내리 끝내야 했지만(그래서 생산성이 그토록 떨어졌을
것이다), 우리는 \.{저장}이라고 치면 한 판을 파일에 접어 두고 내일
이어 할 수 있다.

구현은 게임의 설계가 거저 준 선물이다. 영속 상태가 죄다 |Game| 한
곳에 모여 있으니, 그걸 적었다 되읽으면 끝이다. 다만 \GO/의
\.{encoding/gob}은 비공개 필드를 직렬화하지 못하므로, 내보낸(대문자)
필드를 가진 |saveState|로 옮겨 담는다. 어휘와 메시지는 시작할 때
다시 만들면 되니 젖혀 두고, 턴 사이에만 뜻이 있는 임시 값(|mot|,
|verb|, |obj|$\ldots$)도 저장 시점(턴 경계)에는 무의미하니 뺀다.
난수 상태 |rx|는 꼭 담아야 한다---그래야 저장하고 불러와도
결정론이 이어진다. 전역인 |objBase|도 게임 중에 바뀌므로 함께 간다.

@<저장과 불러오기@>=
const saveFile = "advent.save" // 고정 저장 파일명

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

	@<난쟁이부터 힌트까지의 필드@>@;
}

@ 이어서 난쟁이부터 힌트까지, 나머지 장부들.

@<난쟁이부터 힌트까지의 필드@>=
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

@ 옮겨 담기와 되돌리기. 재미없는 필드 나르기의 연속이지만, 문학적
프로그램이라고 세금이 면제되지는 않는다.

@<저장과 불러오기@>=
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

@ @<저장과 불러오기@>=
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

@ 파일 입출력. |save|는 스냅숏을 gob으로 적고, |load|는 읽어 되돌린다.

@<저장과 불러오기@>=
func (g *Game) save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(g.snapshot())
}

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

@ 저장의 시험. 상태를 어질러 놓고 저장했다가 새 게임에 불러와,
어질러 둔 그대로인지 본다.

@(advent_test.go@>=
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
