@* 게임의 진행.
게임 전체는 결국 두 겹의 순환이다. {\it 큰 순환}은 걸어 다니기다 ---
한 장소로 이동하고, 그곳을 묘사한다. 그 안의 {\it 작은 순환}은 그
자리에서 하기다 --- 명령을 받아 집고, 켜고, 읽고, 그러다 이동 명령이
나오면 큰 순환으로 돌아간다. 원작은 이 두 순환을 |main| 함수 하나에
|goto|로 촘촘히 짜 넣었고, 커누스는 서문에서 ``|goto|가 싫으면 읽지
말라''고 눙쳤다. 이 포팅은 그 흐름의 대부분을 값으로 풀었지만 ---
각 단계가 결과값을 돌려주고 호출자가 갈래를 튼다 --- 큰 순환의 뼈대
만큼은 |goto| 몇 개를 남겼다. 죽음에서 부활로, 부활에서 다시 묘사로
건너뛰는 흐름은 |goto|가 가장 정직한 표현이라서다.

시작하기 전에 살림 몇 가지. 입출력을 필드로 두어 시험이 문자열 버퍼를
꽂을 수 있게 했고, 나머지는 진행을 기록하는 장부들이다.

@<게임의 상태@>=
in  *bufio.Reader // 사용자 입력
out io.Writer     // 게임 출력
eof bool          // 입력이 끝났는가

wasDark       bool // 직전에 어두웠는가
lookCount     int  // ``봐''를 몇 번 했나
interval      int  // 긴 설명을 되풀이하는 주기 (``간략''하면 10000)
tally         int  // 아직 못 본 보물 수
lostTreasures int  // 영영 못 볼 보물 수

quitting bool // 게임을 끝내야 하는가
gaveUp   bool // 살아 있는데 스스로 그만뒀는가

@ 원작의 C~매크로들은 메서드가 됐다. 어두운가, 물이나 기름이 있는가.
어둠의 판정에 램프의 소지와 점등이 끼어드는 데 유의 --- 램프를 들고
있어도 꺼져 있으면 캄캄하다. |ok|는 사소한 동작에 대한 만능 응답
``알았어.''인데, 상수가 아니라 항상 현재 값을 참조한다: 어휘를 만들기
전에는 빈 문자열이라, 환영 인사에 ``아니''라고 답해도 조용히 넘어간다.

@<게임의 진행@>=
func (g *Game) dark() bool {
	return caveFlags[g.loc]&lighted == 0 && (g.prop[LAMP] == 0 || !g.here(LAMP))
}

func (g *Game) waterHere() bool    { return caveFlags[g.loc]&(liquid+oil) == liquid }
func (g *Game) oilHere() bool      { return caveFlags[g.loc]&(liquid+oil) == liquid+oil }
func (g *Game) noLiquidHere() bool { return caveFlags[g.loc]&liquid == 0 }

func (g *Game) ok() string { return g.defaultMsg[RELAX] }

@ 가장 낮은 층의 입력. 한 줄을 읽어 개행을 떼서 돌려주고, 입력이
바닥나면 |ok=false|로 알린다. 원작의 |fgets|에 해당한다.

@<게임의 진행@>=
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
	for len(line) > 0 && (line[len(line)-1] == '\n' || line[len(line)-1] == '\r') {
		line = line[:len(line)-1]
	}
	return line, true
}

@ 예/아니오를 묻는 |yes|. 질문 |q|를 던지고, 긍정이면 |y|를, 부정이면
|n|을 (비어 있지 않을 때) 들려준 뒤 진위를 돌려준다. 답의 판별은 첫
글자만 본다 --- \.{예}·\.{응}이면 긍정, \.{아니}·\.{안}이면 부정,
영어 사용자 습관을 존중해 \.{y}/\.{n}도 받아 준다. ``아마도''는 없다.

@<게임의 진행@>=
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

@ 힌트. 같은 곳에서 한참 헤매는 낌새가 보이면 내가 먼저 말을 건다 ---
``새를 잡으려는 거야?'' 물론 공짜는 아니어서, 힌트마다 점수가 깎인다.
슬롯 0은 환영 인사(게임 설명), 슬롯 1은 굴 껍데기 속 단서, 2부터 7까지가
장소에 걸린 진짜 힌트다. |hintThresh|는 힌트를 꺼내기까지 뜸을 들이는
턴 수 --- 미로 힌트(75턴)는 어지간히 헤매기 전엔 안 나온다.

@<게임의 상태@>=
hinted    [nHints]bool // 이 힌트를 이미 줬는가
hintCount [nHints]int  // 이 힌트가 필요해 보인 지 몇 턴째인가
limit     int          // 램프 수명 카운트다운

@ @<게임의 진행@>=
const nHints = 8

var hintCost = [nHints]int{5, 10, 2, 2, 2, 4, 5, 3}
var hintThresh = [nHints]int{0, 0, 4, 5, 8, 75, 25, 20}


@ 힌트를 꺼내기 전에 던지는 질문들. 0번은 환영 인사다.

@<게임의 진행@>=
var hintPrompt = [nHints]string{
	"어드벤처에 온 걸 환영해!! 설명을 들어볼래?",
	"흐음, 이건 단서 같은데, 읽으면 10점이 깎여.\n" +
		"그래도 읽어줄까?",
	"동굴 안으로 들어가려는 거야?",
	"새를 잡으려는 거야?",
	"뱀을 어떻게든 해보려는 거야?",
	"미로에서 빠져나가는 데 도움이 필요해?",
	"플로버 방 너머를 탐험하려는 거야?",
	"여기서 나가는 데 도움이 필요해?",
}


@ 그리고 힌트의 본문. 0번은 게임 설명으로, 마지막 줄의 crediting이 원작의 것 그대로다 --- ``버그는 전부 Don Knuth가 더했지.''

@<게임의 진행@>=
var hintText = [nHints]string{
	"이 근처 어딘가에 거대한 동굴이 있어. 많은 이들이 그곳에서 보물과\n" +
		"황금으로 부를 얻었지만, 들어간 사람 중 더러는 다시 보이지 않는다는\n" +
		"소문도 있지. 동굴에서는 마법이 통한다고 해. 내가 네 눈과 손이\n" +
		"되어줄게. 한두 단어짜리 명령으로 날 움직여 줘. 이를테면 \"북\",\n" +
		"\"열쇠 가져가\", \"램프 켜\" 처럼. 마법 주문(xyzzy, plugh, plover\n" +
		"같은 것들)은 영어 그대로 입력하면 돼. 막히면 \"도움말\"을 쳐서\n" +
		"일반적인 힌트를 받을 수 있어. 모험을 끝내는 법 등이 궁금하면\n" +
		"\"정보\"를 쳐 봐.\n" +
		"                        -  -  -\n" +
		"최초의 어드벤처 프로그램은 Willie Crowther가 만들었어.\n" +
		"지금 프로그램의 기능 대부분은 Don Woods가 더했고,\n" +
		"버그는 전부 Don Knuth가 더했지.",
	"이렇게 적혀 있어. \"이곳에는 묘한 기운이 감돌아서, 늘 알고 있던 단어\n" +
		"하나가 이제 새로운 효력을 낸다.\"",
	"창살은 아주 단단하고 강철 자물쇠가 채워져 있어. 열쇠 없이는 들어갈\n" +
		"수 없는데, 열쇠는 보이지 않네. 다른 곳을 찾아보는 게 좋겠어.",
	"지금은 뭔가가 새를 겁먹게 하고 있어서, 무슨 수를 써도 잡을 수 없어.\n" +
		"나중에 다시 해보는 게 좋겠어.",
	"뱀은 죽일 수도, 쫓아낼 수도, 피할 수도 없어. 지나갈 방법이 있긴\n" +
		"한데, 지금은 필요한 게 없네.",
	"물건을 떨어뜨려 두면 통로들이 덜 비슷해 보이게 만들 수 있어.",
	"그 구역은 구덩이에 빠질 걱정 없이 탐험하는 방법이 있어. 지금 가진\n" +
		"물건 중엔 당장 그 비밀을 푸는 데 쓸모 있는 게 없어.",
	"서쪽으로 가지 마.",
}

@ |offer|는 환영 인사(|j==0|)나 힌트(|j>=2|)를 제안한다. 힌트는 두 번
묻는다 --- 관심 있는지, 점수를 내고도 원하는지. 힌트를 받으면 램프
수명을 조금 보태 주는 인정도 베푼다.

@<게임의 진행@>=
func (g *Game) offer(j int) {
	if j > 1 {
		if !g.yes(hintPrompt[j], " 힌트를 줄 준비가 됐어,", g.ok()) {
			return
		}
		fmt.Fprintf(g.out, " 하지만 %d점이 깎일 거야. ", hintCost[j])
		g.hinted[j] = g.yes("힌트를 원해?", hintText[j], g.ok())
	} else {
		g.hinted[j] = g.yes(hintPrompt[j], hintText[j], g.ok())
	}
	if g.hinted[j] && g.limit > 30 {
		g.limit += 30 * hintCost[j]
	}
}

@ 힌트가 지금 어울리는지는 매 턴 |giveHint|가 가린다. 현재 장소에 힌트
비트가 걸려 있고, 그 상태로 |hintThresh|만큼 버텼고, 상황 조건까지
맞으면 제안한다. 새 힌트(|j==3|)만은 조건이 어긋나도 턴 수를 지우지
않는다 --- 막대를 든 채 새에게 다가가는 안쓰러운 시도는 이어지고 있는
셈이니까.

@<게임의 진행@>=
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
		@<힌트 j가 지금 어울리는지 apply에 담는다@>@;
		if apply {
			g.offer(j)
		}
		g.hintCount[j] = 0
	}
}

@ 조건들은 게임의 급소를 하나씩 짚는다: 열쇠 없이 격자문 앞(2), 막대를
든 채 새 앞(3), 새 없이 뱀 앞(4), 빈손도 아니면서 표식 없는 미로(5),
에메랄드는 봤는데 피라미드를 못 봄(6), 그리고 Witt's End(7)는 조건이
따로 없다 --- 거기 있다는 것 자체가 곤경이다.

@<힌트 j가 지금 어울리는지 apply에 담는다@>=
switch j {
case 2:
	apply = g.prop[GRATE] == 0 && !g.here(KEYS)
case 3:
	if g.here(BIRD) && g.oldobj == BIRD && g.toting(ROD) {
		apply = true
	} else {
		continue // hintCount 유지
	}
case 4:
	apply = g.here(SNAKE) && !g.here(BIRD)
case 5:
	apply = g.first[g.loc] == NOTHING && g.first[g.oldloc] == NOTHING &&
		g.first[g.oldoldloc] == NOTHING && g.holding > 1
case 6:
	apply = g.prop[EMERALD] != -1 && g.prop[PYRAMID] == -1
case 7:
	apply = true
}

@ 이제 게임 한 판. 초기화의 대부분이 환영 인사를 읽는 사이에 이뤄진다는
원작의 주석을 그대로 살렸다 --- 설명을 청한 신중한 모험가는 램프 수명
1000턴을, 곧장 뛰어든 성급한 모험가는 330턴을 받는다는 것이 이 게임
최초의 갈림길이다.

@<게임의 진행@>=
func (g *Game) Run() {
	g.offer(0) // 환영 인사와, 원하면 게임 설명
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

@ 큰 순환의 뼈대. 남겨 둔 |goto|들이 여기 다 모여 있다. 순환의 꼭대기는
이동이다: 폐쇄나 난쟁이가 길을 막는지 보고, 실제로 옮기고, 난쟁이와
해적도 한 걸음 움직인다. |commence|는 도착지 묘사, |tryMove|는 다음
위치 계산, |death|는 죽음 처리다. 각 단계가 돌려주는 값이 다음 갈 곳을
말해 준다.

@<게임의 진행@>=
func (g *Game) simulate() {
	for {
		@<폐쇄나 난쟁이가 이동을 막으면 제자리로 돌린다@>@;
		g.loc = g.newloc // 실제로 이동
		g.moveDwarves()  // 난쟁이/해적 이동
		if g.dying {
			goto death
		}

	commence:
		@<도착지를 보고하고 작은 순환을 돈다@>@;
	tryMove:
		@<다음 위치를 계산한다@>@;
	death:
		if !g.handleDeath() { // 부활 실패면 게임 종료
			return
		}
		goto commence // 부활: 새 장소에서 상태 다시 보고
	}
}

@ |tryMove|: 특수 이동을 가리고, 보통 이동이면 위치 기억을 갱신한 뒤 이동표를 푼다.

@<다음 위치를 계산한다@>=
		switch g.handleSpecialMotion() {
		case smStay:
			continue // 제자리에서 다시 상태 보고
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


@ |commence|: 도착지를 묘사하고, 작은 순환이 돌려주는 값대로 갈래를 튼다.

@<도착지를 보고하고 작은 순환을 돈다@>=
		if g.reportState() { // 강제 이동 장소면 true: 곧장 다시 이동
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


@ 작은 순환. 명령을 받아 동작을 수행하고, 큰 순환으로 돌아가야 할 일이
생기면 그 사유를 값으로 돌려준다. |keep|은 ``가져가 뭘?''처럼 동사를
쥔 채 객체만 다시 묻는 상태다.

@<게임의 진행@>=
type minorResult int

const (
	mcEnd      minorResult = iota // 게임 종료 (EOF 또는 그만)
	mcTryMove                     // 이동 요청
	mcCommence                    // 제자리에서 상태만 다시 보고
	mcDeath                       // 죽음
)


@ 순환 자체는 단순하다: 입력을 받아 수행하고, 결과가 큰 순환의 소관이면 돌려보낸다.

@<게임의 진행@>=
func (g *Game) minorCycle() minorResult {
	keep := false
	for {
		var res paResult
		@<명령을 받아 수행하고 결과를 res에 담는다@>@;
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

func (g *Game) report(msg string) {
	if msg != "" {
		fmt.Fprintf(g.out, "%s\n", msg)
	}
}

@ 명령의 다섯 갈래. 이동이면 곧장 큰 순환으로, 동작이면 수행하고 그 결과를 살핀다.

@<명령을 받아 수행하고 결과를 res에 담는다@>=
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

@ 도착지의 묘사. 방문 횟수에 따라 긴 설명과 짧은 설명을 고르고, 곰을
데리고 다니면 그 사실도 알린다. 어둠 속을 계속 걸으면 35\%의 확률로
구덩이에 떨어져 죽는다 --- 이 게임이 플레이어에게 램프를 아끼라고
가르치는 방식이다. Y2 방에서 이따금 들려오는 공허한 목소리
``PLUGH''는 원작이 심어 둔 가장 유명한 힌트다.

@<게임의 진행@>=
const pitchDarkMsg = "이제 칠흑같이 어두워. 계속 가면 십중팔구 구덩이에 빠질 거야."

func (g *Game) reportState() bool {
	if g.loc == limbo {
		g.die()
		return false
	}
	var p string
	if g.dark() && !forcedMove(g.loc) {
		if g.wasDark && g.pct(35) {
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

@ 이어서 그 방의 물건들을 묘사한다. 보물을 처음 보는 순간이 여기다:
음수로 시작한 |prop|을 0(양탄자와 사슬은 1)으로 바꾸고 남은 보물 수
|tally|를 줄인다. 계단의 묘사는 금덩이를 들었는지, 어느 끝에 서 있는지에
따라 달라지는 잔재미도 있다.

@<게임의 진행@>=
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
			@<남은 보물이 너무 꼭꼭 숨었으면 램프를 재촉한다@>@;
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

@ 마지막으로 잔손 둘. 폐쇄가 시작됐는지의 판정과, 죽음의 표시.
실제 죽음 처리는 큰 순환의 |death| 분기가 맡는다.

@<게임의 진행@>=
func (g *Game) closing() bool { return g.clock1 < 0 }
func (g *Game) die()          { g.dying = true }
