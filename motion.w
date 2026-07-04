@* 이동.
이동 동사 |mot|가 정해지면 큰 순환은 여기로 온다. 할 일은 이동표에서
알맞은 명령을 찾아 |newloc|을 정하는 것 --- 다만 몇몇 이동은 표를
거치지 않는 별종이라 먼저 골라낸다. |NOWHERE|는 제자리걸음이고,
\.{봐}(LOOK)는 어둠을 못 본 척하며 긴 설명을 다시 들려주고(어둠 속
구덩이 판정을 일부러 건너뛰어, 구경만 하다 죽는 일은 없게 한다),
\.{동굴}(CAVE)은 위치에 따라 힌트 또는 핀잔이다.

@<이동@>=
type smResult int

const (
	smNormal  smResult = iota // 보통: oldloc 갱신 후 다음 위치 계산
	smStay                    // 제자리
	smGoForIt                 // BACK: oldloc 갱신을 끝냈으니 곧장 계산
)

func (g *Game) handleSpecialMotion() smResult {
	g.newloc = g.loc // 기본은 제자리
	switch g.mot {
	case NOWHERE:
		return smStay
	case BACK:
		return g.tryGoBack()
	case LOOK:
		if g.lookCount < 3 {
			g.lookCount++
			fmt.Fprintf(g.out, "미안하지만 더 자세히는 알려줄 수 없어. 네 위치의 긴 설명을\n"+
				"다시 보여줄게.\n")
		}
		g.wasDark = false
		g.visits[g.loc] = 0
		return smStay
	case CAVE:
		if g.loc < minInCave {
			fmt.Fprintf(g.out, "동굴이 어디 있는지 안 보이지만, 이 근처에선 어떤 시내도 지표면을\n"+
				"오래 흐르지 못해. 나라면 시내를 따라가 보겠어.\n")
		} else {
			fmt.Fprintf(g.out, "그걸 하려면 더 자세한 지시가 필요해.\n")
		}
		return smStay
	}
	return smNormal
}

@ \.{돌아가}(BACK)는 직전 위치로 --- 직전 위치가 강제 이동 더미였다면
그 전 위치로 --- 가는 길을 이동표에서 거꾸로 찾는다. 곧장 가는 명령이
없으면, 강제 이동 더미를 거쳐 닿는 명령(|qq|)이라도 받아 준다. 그마저
없으면 여기서는 거기로 못 간다고 실토한다.

@<이동@>=
func (g *Game) tryGoBack() smResult {
	var l location
	if forcedMove(g.oldloc) {
		l = g.oldoldloc
	} else {
		l = g.oldloc
	}
	g.oldoldloc = g.oldloc
	g.oldloc = g.loc
	if l == g.loc {
		fmt.Fprintf(g.out, "미안한데, 네가 여기 어떻게 왔는지 이제 기억이 안 나.\n")
		return smStay
	}
	found, qq := -1, -1
	for q := caveStart[g.loc]; q < caveStart[g.loc+1]; q++ {
		ll := caveTravels[q].dest
		if ll == l {
			found = q
			break
		}
		if ll <= maxLoc && forcedMove(ll) && caveTravels[caveStart[ll]].dest == l {
			qq = q
		}
	}
	q := found
	if q < 0 {
		if qq < 0 {
			fmt.Fprintf(g.out, "여기서는 거기로 갈 수 없어.\n")
			return smStay
		}
		q = qq
	}
	g.mot = caveTravels[q].mot
	return smGoForIt
}

@ 본론인 이동표 해석. 현재 장소의 명령들(|caveStart[loc]|부터 다음 장소
직전까지) 가운데 |mot|에 맞는 첫 명령을 찾고 --- 강제 이동 장소면
무조건 첫 명령이다 --- 조건을 만족할 때까지 전진한 뒤 목적지를 읽는다.
목적지가 진짜 장소면 끝. remark면 한마디 듣고 제자리. 그 사이의 특수
목적지 셋은 이 게임의 명장면들이다: 플로버 통로(|ppass|), 에메랄드
자동 낙하(|pdrop|), 그리고 트롤 다리(|troll|).

@<이동@>=
func (g *Game) determineNextLocation() {
	qEnd := caveStart[g.loc+1]
	q := caveStart[g.loc]
	for ; q < qEnd; q++ {
		if forcedMove(g.loc) || caveTravels[q].mot == g.mot {
			break
		}
	}
	if q == qEnd {
		@<쓸 수 없는 이동에 대해 해명한다@>@;
		return
	}
	for {
		q = g.advanceCondition(q)
		g.newloc = caveTravels[q].dest
		if g.newloc <= maxLoc {
			return
		}
		if g.newloc > maxSpec {
			fmt.Fprintf(g.out, "%s\n", remarkOf(g.newloc))
			g.newloc = g.loc
			return
		}
		switch g.newloc {
		case ppass:
			@<가능하면 플로버 통로를 지난다@>@;
			return
		case pdrop:
			g.drop(EMERALD, g.loc) // 에메랄드를 떨구고 조건을 재평가
			q = noGoodSkip(q)
			continue
		case troll:
			@<가능하면 트롤 다리를 건넌다@>@;
			return
		}
	}
}

@ 조건값 |cond|의 인코딩은 이렇다. 포트란 시절의 정수 하나에 우겨 넣은
살림이 그대로 왔다:

$$\vbox{\halign{\indent\hfil#\hfil&\quad#\hfil\cr
|cond==0|&        항상 참\cr
$0<|cond|<100$&   |cond|\%의 확률로 참\cr
|cond==100|&      난쟁이만 빼고 항상 참\cr
$100<|cond|\le200$& 물체 |cond%100|을 들고 있어야 함\cr
$200<|cond|\le300$& 물체 |cond%100|이 현재 장소에 있어야 함\cr
$|cond|>300$&     |prop[cond%100] != (cond-300)/100| 이어야 함\cr}}$$

조건이 어긋나면 |noGoodSkip|이 목적지·조건이 같은 명령 무리를 통째로
건너뛰고, 만족하는 명령이 나올 때까지 되풀이한다. 표가 올바르다면
마지막 명령은 반드시 잡히게 짜여 있다.

@<이동@>=
func (g *Game) advanceCondition(q int) int {
	for {
		j := caveTravels[q].cond
		satisfied := false
		switch {
		case j > 300:
			satisfied = g.prop[object(j%100)] != (j-300)/100
		case j <= 100:
			satisfied = j == 0 || g.pct(j)
		default:
			satisfied = g.toting(object(j%100)) || (j >= 200 && g.isAtLoc(object(j%100)))
		}
		if satisfied {
			return q
		}
		q = noGoodSkip(q)
	}
}

func noGoodSkip(q int) int {
	qq := q
	for q++; caveTravels[q].dest == caveTravels[qq].dest && caveTravels[q].cond == caveTravels[qq].cond; q++ {
	}
	return q
}

@ 이 장소에 그 이동이 아예 없으면 사정을 설명한다. 대꾸는 이동 동사의
결에 맞춘다 --- 주문을 엉뚱한 데서 외치면 ``아무 일도 일어나지
않아''(WAVE의 기본 메시지를 빌려 쓴다), 방향 감각이 필요한 동사면
나침반을 권한다.

@<쓸 수 없는 이동에 대해 해명한다@>=
switch {
case g.mot == CRAWL:
	fmt.Fprintf(g.out, "어느 쪽으로?")
case g.mot == XYZZY || g.mot == PLUGH:
	fmt.Fprintf(g.out, "%s", g.defaultMsg[WAVE])
case g.verb == FIND || g.verb == INVENTORY:
	fmt.Fprintf(g.out, "%s", g.defaultMsg[FIND])
case g.mot <= FORWARD:
	switch g.mot {
	case IN, OUT:
		fmt.Fprintf(g.out, "여기선 안과 밖을 모르겠어. 나침반 방향을 쓰거나 가려는 쪽에\n"+
			"있는 걸 말해 줘.")
	case FORWARD, L, R:
		fmt.Fprintf(g.out, "네가 어느 쪽을 보고 있는지 모르겠어. 나침반 방향이나 가까운 물건으로 말해 줘.")
	default:
		fmt.Fprintf(g.out, "그 방향으로는 갈 수 없어.")
	}
default:
	fmt.Fprintf(g.out, "그 단어를 여기서 어떻게 써야 할지 모르겠어.")
}
fmt.Fprintf(g.out, "\n")

@ 플로버 통로는 이 게임에서 가장 짓궂은 설계다. 알코브와 플로버 방을
잇는 이 좁은 굴은 빈손이거나 에메랄드 {\it 하나만} 들었을 때만
통과된다 --- 램프조차 안 된다. 에메랄드를 꺼내 오라고 만든 길이니까.
목적지 계산 |alcove+proom-loc|은 두 끝을 오가는 상투적인 산수다.

@<가능하면 플로버 통로를 지난다@>=
if g.holding == 0 || (g.toting(EMERALD) && g.holding == 1) {
	g.newloc = alcove + proom - g.loc
} else {
	fmt.Fprintf(g.out, "네가 든 뭔가가 굴을 너랑 같이 통과하지 못해.\n"+
		"소지품을 확인하고 뭔가 내려놓는 게 좋겠어.\n")
	g.newloc = g.loc
}

@ 트롤 다리. 통행료를 이미 냈으면(|prop[TROLL]==0|) 건너고, 트롤이
버티고 있으면(|prop[TROLL]==1|) 트롤이 다리 밑에서 나와 앞을 막는다.
건너는 데 성공해도 곰을 데리고 있다면 --- 다리가 곰의 무게를 견디지
못한다. 곰과 함께 협곡으로 떨어지며, 다리 건너에 남겨 둔 향신료와
사슬은 영영 잃은 보물이 된다. 이 죽음의 부활 지점만은 특별히 다리를
건넌 자리다: 건너다 죽었으니 건넌 곳에서 깨어나는 것이 순리라던
원작의 처사다.

@<가능하면 트롤 다리를 건넌다@>=
if g.prop[TROLL] == 1 { // 트롤이 길을 막음
	g.move(TROLL, swside)
	g.move(TROLL_, neside)
	g.prop[TROLL] = 0
	g.destroy(TROLL2)
	g.destroy(TROLL2_)
	g.move(BRIDGE, swside)
	g.move(BRIDGE_, neside)
	fmt.Fprintf(g.out, "%s\n", objNote[objOffset[TROLL]+1])
	g.newloc = g.loc
	return
}
g.newloc = neside + swside - g.loc // 건넌다
if g.prop[TROLL] == 0 {
	g.prop[TROLL] = 1
}
if !g.toting(BEAR) {
	return
}
fmt.Fprintf(g.out, "막 건너편에 닿는 순간, 여태 널 따라다니던 곰의 무게에\n"+
	"다리가 휘청해. 넌 필사적으로 붙잡을 곳을\n"+
	"더듬지만, 다리가 무너지면서 뒤로\n"+
	"휘청이다 협곡으로 떨어져.\n")
g.prop[BRIDGE] = 1
g.prop[TROLL] = 2
g.drop(BEAR, g.newloc)
objBase[BEAR] = BEAR
g.prop[BEAR] = 3 // 곰이 죽었다
if g.prop[SPICES] < 0 && g.place[SPICES] >= neside {
	g.lostTreasures++
}
if g.prop[CHAIN] < 0 && g.place[CHAIN] >= neside {
	g.lostTreasures++
}
g.oldoldloc = g.newloc // 되살아나면 다리를 건넌 상태
g.die()
