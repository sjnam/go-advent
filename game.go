package main

import (
	"bufio"
	"fmt"
	"io"
)

// Game은 원본 C 프로그램의 모든 전역 상태를 한곳에 모은 것이다.
// 단계가 진행되면서 필드가 계속 추가된다.
type Game struct {
	in  *bufio.Reader // 사용자 입력
	out io.Writer     // 게임 출력
	eof bool          // 입력이 끝났는가 (스크립트 재생/테스트용)

	rx int // 난수 생성기 상태 (rand.go 참고)

	word1, word2 string // listen()이 채우는 명령 단어 둘

	// 어휘 (vocab.go 참고)
	words      map[string]vocabEntry // 우리가 아는 단어들
	defaultMsg [30]string            // 동작이 특별 조건을 못 채웠을 때의 기본 메시지(action별)
	message    [13]string            // help, info 등 고정 메시지

	// 동굴 상태 (cave.go / cavedata.go 참고)
	visits                         [maxLoc + 1]int // 각 장소를 몇 번 방문했나
	oldoldloc, oldloc, loc, newloc location        // 최근/다음 위치

	// 객체 상태 (objects.go / objdata.go 참고)
	place   [maxObj + 1]location // 각 객체의 현재 위치
	prop    [maxObj + 1]int      // 각 객체의 현재 속성값
	link    [maxObj + 1]object   // 같은 장소의 다음 객체
	first   [maxLoc + 1]object   // 각 장소의 첫 객체
	holding int                  // 들고 있는 객체 수

	// 파싱/명령 상태 (parse.go 참고)
	mot      motion // 지금 지정된 이동
	verb     action // 지금 지정된 동작
	oldverb  action // 바뀌기 전 verb
	obj      object // 지금 지정된 객체
	oldobj   object // 이전 obj
	turns    int    // 명령을 읽은 횟수
	speakIdx int    // message_type 단어가 가리키는 메시지 인덱스

	// 상태 보고 상태
	wasDark       bool // 최근에 어두웠는가
	lookCount     int  // LOOK을 몇 번 했나
	interval      int  // BRIEF면 10000이 됨
	tally         int  // 아직 못 본 보물 수
	lostTreasures int  // 영영 못 볼 보물 수

	// 난쟁이/해적 상태 (dwarf.go 참고)
	dflag    int              // 난쟁이가 얼마나 화났나
	dkill    int              // 난쟁이를 몇 마리 죽였나
	dloc     [nd + 1]location // 각 난쟁이(0번은 해적)의 현재 위치
	odloc    [nd + 1]location // 직전 위치
	dseen    [nd + 1]bool     // 난쟁이가 나를 봤는가
	dtotal   int              // 같은 방에 있는 난쟁이 수
	attack   int              // 칼을 뽑을 틈이 있던 난쟁이 수
	stick    int              // 칼을 명중시킨 난쟁이 수
	knifeLoc int              // 칼이 언급된 곳, 없으면 -1

	// 폐쇄/죽음 상태 (closing.go / death.go 참고)
	clock1     int  // 모든 보물을 본 뒤 폐쇄까지의 카운트다운 (시작 15)
	clock2     int  // 폐쇄 시작 뒤 동굴이 닫힐 때까지 (시작 30)
	panicked   bool // 폐쇄 중 탈출을 시도해 당황했는가
	warned     bool // 램프 전력 부족 경고를 받았는가
	deathCount int  // 몇 번 죽었나
	bonus      int  // 마지막 퍼즐에서 얻는 추가 점수
	dying      bool // 이번에 죽었는가 (죽음 처리 대기)
	closed     bool // 동굴이 완전히 닫혔는가
	foobar     int  // fee-fie-foe-foo 주문 진행도

	// 힌트/환영 메시지 상태
	hinted    [nHints]bool
	hintCount [nHints]int // 이 힌트가 필요해 보인 지 얼마나 됐나
	limit     int         // 램프 수명 카운트다운

	quitting bool // 게임을 끝내야 하는가
	gaveUp   bool // 살아 있는데 스스로 그만뒀는가
}

// ---- 위치/액체/어둠 판정 (advent.w의 매크로들) ----

// dark는 지금 어두운지 본다. (advent.w:2680)
func (g *Game) dark() bool {
	return caveFlags[g.loc]&lighted == 0 && (g.prop[LAMP] == 0 || !g.here(LAMP))
}

// waterHere/oilHere/noLiquidHere: 현재 장소의 액체 상태 (advent.w:2494-2496)
func (g *Game) waterHere() bool    { return caveFlags[g.loc]&(liquid+oil) == liquid }
func (g *Game) oilHere() bool      { return caveFlags[g.loc]&(liquid+oil) == liquid+oil }
func (g *Game) noLiquidHere() bool { return caveFlags[g.loc]&liquid == 0 }

// ok는 사소한 동작에 대한 기본 응답("OK.")이다. 원본의 매크로
// `default_msg[RELAX]`와 같이 항상 현재 값을 참조한다 — 어휘 빌드 전에는
// 빈 문자열이므로 환영 메시지에서 "no"라고 답해도 아무것도 출력되지 않는다.
func (g *Game) ok() string { return g.defaultMsg[RELAX] }

// readLine은 한 줄을 읽어 돌려준다(개행 제외). EOF면 ok=false.
// 원본의 fgets에 해당한다.
func (g *Game) readLine() (string, bool) {
	if g.eof {
		return "", false
	}
	line, err := g.in.ReadString('\n')
	if len(line) == 0 && err != nil {
		g.eof = true
		return "", false
	}
	if err != nil {
		g.eof = true // 마지막 줄에 개행이 없는 경우
	}
	// 끝의 \n 또는 \r\n 제거
	for len(line) > 0 && (line[len(line)-1] == '\n' || line[len(line)-1] == '\r') {
		line = line[:len(line)-1]
	}
	return line, true
}

// yes는 질문 q를 출력하고 예/아니오 답을 기다린다. 답이 긍정이면 y를,
// 부정이면 n을 (비어 있지 않을 때) 출력하고 그 진위를 돌려준다.
// (advent.w "Low-level input"의 yes 서브루틴)
func (g *Game) yes(q, y, n string) bool {
	for {
		fmt.Fprintf(g.out, "%s\n** ", q)
		line, more := g.readLine()
		if !more {
			return false // EOF: 부정으로 간주
		}
		switch affirmative(line) {
		case 1:
			if y != "" {
				fmt.Fprintf(g.out, "%s\n", y)
			}
			return true
		case 0:
			if n != "" {
				fmt.Fprintf(g.out, "%s\n", n)
			}
			return false
		default:
			fmt.Fprintf(g.out, " \"예\" 아니면 \"아니\"로 대답해 줘.\n")
		}
	}
}

// affirmative는 답을 긍정(1)/부정(0)/모호(-1)로 가린다.
// 한글 "예·응·ㅇ" 또는 영문 y는 긍정, "아니·안·ㄴ" 또는 영문 n은 부정.
func affirmative(s string) int {
	r := []rune(s)
	if len(r) == 0 {
		return -1
	}
	switch r[0] {
	case '예', '응', 'ㅇ', 'y', 'Y':
		return 1
	case '아', '안', 'ㄴ', 'n', 'N':
		return 0
	}
	return -1
}

// offer는 환영 메시지(j==0)나 힌트(j>=2)를 제안한다.
// (advent.w "Scoring"의 offer 서브루틴)
func (g *Game) offer(j int) {
	if j > 1 {
		if !g.yes(hintPrompt[j], " 힌트를 줄 준비가 됐어,", g.ok()) {
			return
		}
		fmt.Fprintf(g.out, " 하지만 %d점이 깎일 거야.  ", hintCost[j])
		g.hinted[j] = g.yes("힌트를 원해?", hintText[j], g.ok())
	} else {
		g.hinted[j] = g.yes(hintPrompt[j], hintText[j], g.ok())
	}
	if g.hinted[j] && g.limit > 30 {
		g.limit += 30 * hintCost[j]
	}
}

// Run은 게임을 처음부터 끝까지 진행한다.
// (advent.w main + "Launching the program")
func (g *Game) Run() {
	// 대부분의 초기화는 환영 메시지를 읽는 동안 이뤄진다.
	g.offer(0) // 환영 메시지와, 원하면 게임 설명
	if g.hinted[0] {
		g.limit = 1000
	} else {
		g.limit = 330
	}
	g.buildVocabulary()
	g.loadObjects()
	g.interval = 5
	g.tally = 15
	g.knifeLoc = -1
	g.clock1 = 15
	g.clock2 = 30
	g.dloc = dwarfStart
	g.oldoldloc, g.oldloc, g.loc, g.newloc = road, road, road, road

	g.simulate()
	g.printScore() // 점수를 매기고 작별 인사
}

// simulate는 게임의 큰 순환(major cycle: 이동+상태보고)과
// 그 안의 작은 순환(minor cycle: 입력+동작)을 돈다.
// (advent.w "The main control loop")
func (g *Game) simulate() {
	for {
		g.checkInterference() // 폐쇄/난쟁이가 이동을 막는지
		g.loc = g.newloc      // 실제로 이동
		g.moveDwarves()       // 난쟁이/해적 이동
		if g.dying {
			goto death
		}

	commence:
		if g.reportState() { // forced move 등으로 곧장 다시 이동해야 하면 true
			goto tryMove
		}
		if g.dying {
			goto death
		}
		if g.quitting {
			return
		}
		switch g.minorCycle() {
		case mcEnd:
			return
		case mcDeath:
			goto death
		case mcTryMove:
			goto tryMove
		case mcCommence:
			goto commence // 불을 켜는 등으로 제자리에서 다시 묘사
		}

	tryMove:
		switch g.handleSpecialMotion() {
		case smStay:
			continue // 큰 순환 처음으로 (제자리에서 다시 상태 보고)
		case smGoForIt:
			// BACK: oldloc 갱신을 이미 끝냈으니 곧장 다음 위치 계산
		case smNormal:
			g.oldoldloc = g.oldloc
			g.oldloc = g.loc
		}
		g.determineNextLocation()
		if g.dying {
			goto death
		}
		continue

	death:
		if !g.handleDeath() { // 부활 실패면 게임 종료
			return
		}
		goto commence // 부활: 새 장소에서 상태 다시 보고
	}
}

// minorCycle이 큰 순환에 돌려주는 결과.
type minorResult int

const (
	mcEnd      minorResult = iota // 게임 종료(EOF 또는 quit)
	mcTryMove                     // 이동 요청 — 큰 순환의 이동 처리로
	mcCommence                    // 제자리에서 상태만 다시 보고 (불 켜기 등)
	mcDeath                       // 죽음 — 큰 순환의 death 처리로
)

// minorCycle은 작은 순환을 돈다: 명령을 받아 동작을 수행하고, 이동
// 명령이 나오면 큰 순환으로 돌려보낸다. (advent.w 안쪽 while 루프)
func (g *Game) minorCycle() minorResult {
	keep := false // "Take what?" 뒤처럼 verb를 유지해야 하는가
	for {
		var res paResult
		switch g.getUserInput(keep) {
		case inputEOF:
			return mcEnd
		case inputMotion:
			return mcTryMove
		case inputTransitive:
			res = g.performAction(true)
		case inputIntransitive:
			res = g.performAction(false)
		case inputSpeak:
			g.report(g.message[g.speakIdx])
			res = paDone
		}
		if g.quitting {
			return mcEnd
		}
		if g.dying {
			return mcDeath
		}
		switch res {
		case paTryMove:
			return mcTryMove
		case paCommence:
			return mcCommence
		case paNeedObject:
			keep = true
		default:
			keep = false
		}
	}
}

// report는 메시지를 출력한다(원본 매크로 report). 작은 순환은 호출 측에서
// 자연히 다음 입력으로 넘어가며 끝난다.
func (g *Game) report(msg string) {
	if msg != "" {
		fmt.Fprintf(g.out, "%s\n", msg)
	}
}

// reportState는 현재 장소를 묘사한다. 방문 횟수에 따라 긴/짧은 설명을 고르고,
// 어두우면 그렇게 알린다. forced move 장소면 true를 돌려준다(곧장 이동).
// (advent.w "Report the current state")
func (g *Game) reportState() bool {
	if g.loc == limbo {
		g.die()
		return false
	}
	var p string
	if g.dark() && !forcedMove(g.loc) {
		if g.wasDark && g.pct(35) {
			// 어둠 속에서 구덩이에 떨어져 죽는다. (advent.w pitch_dark)
			fmt.Fprintf(g.out, "구덩이에 떨어져서 온몸의 뼈가 다 부러졌어!\n")
			g.oldoldloc = g.loc
			g.die()
			return false
		}
		p = pitchDarkMsg
	} else if caveShortDesc[g.loc] == "" || g.visits[g.loc]%g.interval == 0 {
		p = caveLongDesc[g.loc]
	} else {
		p = caveShortDesc[g.loc]
	}
	if g.toting(BEAR) {
		fmt.Fprintf(g.out, "아주 크고 온순한 곰이 널 따라오고 있어.\n")
	}
	if p != "" {
		fmt.Fprintf(g.out, "\n%s\n", p)
	}
	if forcedMove(g.loc) {
		return true
	}
	if g.loc == y2 && g.pct(25) && !g.closing() {
		fmt.Fprintf(g.out, "공허한 목소리가 \"PLUGH\"라고 말해.\n")
	}
	if !g.dark() {
		g.describeObjects()
	}
	return false
}

// describeObjects는 현재 장소에 있는 물건들을 묘사한다. 보물을 처음 보면
// 속성값을 초기화하고 남은 보물 수를 줄인다. (advent.w "Describe the objects")
func (g *Game) describeObjects() {
	g.visits[g.loc]++
	for t := g.first[g.loc]; t != NOTHING; t = g.link[t] {
		tt := objBase[t]
		if tt == NOTHING {
			tt = t
		}
		if g.prop[tt] < 0 { // 보물을 처음 발견
			if g.closed {
				continue // 폐쇄 후엔 자동 속성 변경 없음
			}
			if tt == RUG || tt == CHAIN {
				g.prop[tt] = 1
			} else {
				g.prop[tt] = 0
			}
			g.tally--
			g.zapLampIfElusive() // 후속(난쟁이/램프)
		}
		if tt == TREADS && g.toting(GOLD) {
			continue
		}
		idx := g.prop[tt] + objOffset[tt]
		if tt == TREADS && g.loc == emist {
			idx++
		}
		if p := objNote[idx]; p != "" {
			fmt.Fprintf(g.out, "%s\n", p)
		}
	}
}

// giveHint는 같은 곳에서 한참 헤매면 힌트를 제안한다.
// (advent.w "Check if a hint applies, and give it if requested")
func (g *Game) giveHint() {
	for j, k := 2, caveHint; j <= 7; j, k = j+1, k*2 {
		if g.hinted[j] {
			continue
		}
		if caveFlags[g.loc]&k == 0 {
			g.hintCount[j] = 0
			continue
		}
		g.hintCount[j]++
		if g.hintCount[j] < hintThresh[j] {
			continue
		}
		apply := false
		switch j {
		case 2: // 동굴에 들어가려는가
			apply = g.prop[GRATE] == 0 && !g.here(KEYS)
		case 3: // 새를 잡으려는가
			if g.here(BIRD) && g.oldobj == BIRD && g.toting(ROD) {
				apply = true
			} else {
				continue // hintCount 유지
			}
		case 4: // 뱀을 다루려는가
			apply = g.here(SNAKE) && !g.here(BIRD)
		case 5: // 미로에서 길을 잃었는가
			apply = g.first[g.loc] == NOTHING && g.first[g.oldloc] == NOTHING &&
				g.first[g.oldoldloc] == NOTHING && g.holding > 1
		case 6: // 플로버 방 너머를 탐험하려는가
			apply = g.prop[EMERALD] != -1 && g.prop[PYRAMID] == -1
		case 7: // 여기서 나가려는가 (Witt's End)
			apply = true
		}
		if apply {
			g.offer(j)
		}
		g.hintCount[j] = 0 // bypass
	}
}

// closing은 폐쇄가 시작됐는지 본다. (advent.w 매크로 closing)
func (g *Game) closing() bool { return g.clock1 < 0 }

// die는 죽음을 표시한다. 실제 처리(부활/종료)는 simulate의 death 분기에서.
func (g *Game) die() { g.dying = true }
