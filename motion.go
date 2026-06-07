package main

import "fmt"

// 이동 처리. 큰 순환은 이동 동사 mot가 주어지고 그에 맞는 newloc을
// 계산하면 끝난다. (advent.w "Motions")

// handleSpecialMotion이 큰 순환에 돌려주는 신호.
type smResult int

const (
	smNormal  smResult = iota // 보통: oldloc 갱신 후 다음 위치 계산
	smStay                    // 제자리 (큰 순환 처음으로)
	smGoForIt                 // BACK: oldloc 갱신을 끝냈으니 곧장 다음 위치 계산
)

// handleSpecialMotion은 이동표를 직접 거치지 않는 이동들을 먼저 처리한다.
// (advent.w "Handle special motion words")
func (g *Game) handleSpecialMotion() smResult {
	g.newloc = g.loc // 기본은 제자리
	switch g.mot {
	case NOWHERE:
		return smStay
	case BACK:
		return g.tryGoBack()
	case LOOK:
		// 어둠을 못 본 척하고 긴 설명을 다시 보여준다(구덩이에 안 빠지게).
		if g.lookCount < 3 {
			g.lookCount++
			fmt.Fprintf(g.out, "미안하지만 더 자세히는 알려줄 수 없어.  네 위치의 긴 설명을\n"+
				"다시 보여줄게.\n")
		}
		g.wasDark = false
		g.visits[g.loc] = 0
		return smStay
	case CAVE:
		if g.loc < minInCave {
			fmt.Fprintf(g.out, "동굴이 어디 있는지 안 보이지만, 이 근처에선 어떤 시내도 지표면을\n"+
				"오래 흐르지 못해.  나라면 시내를 따라가 보겠어.\n")
		} else {
			fmt.Fprintf(g.out, "그걸 하려면 더 자세한 지시가 필요해.\n")
		}
		return smStay
	}
	return smNormal
}

// tryGoBack은 loc에서 oldloc(또는 oldloc이 강제이동이면 oldoldloc)으로
// 돌아가는 이동을 찾는다. (advent.w "Try to go back")
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

// determineNextLocation은 이동표를 해석해 newloc을 정한다.
// (advent.w "Determine the next location")
func (g *Game) determineNextLocation() {
	qEnd := caveStart[g.loc+1]
	q := caveStart[g.loc]
	for ; q < qEnd; q++ {
		if forcedMove(g.loc) || caveTravels[q].mot == g.mot {
			break
		}
	}
	if q == qEnd {
		g.reportInapplicableMotion()
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
			g.choosePloverPassage()
			return
		case pdrop:
			// 에메랄드를 떨어뜨리고(플로버 통로를 쓰게 만든다) 조건을 재평가한다.
			g.drop(EMERALD, g.loc)
			q = noGoodSkip(q)
			continue
		case troll:
			g.crossTrollBridge()
			return
		}
	}
}

// advanceCondition은 q의 조건이 만족될 때까지 같은 목적지/조건 묶음을
// 건너뛰며 전진한다. (advent.w "If the condition ... isn't satisfied")
//
// 조건값 cond의 인코딩 (advent.w:832-834):
//
//	cond==0       : 항상 참
//	0<cond<100    : cond% 확률로 참
//	cond==100     : 난쟁이만 빼고 항상 참
//	100<cond<=200 : 물체 (cond%100)을 들고 있어야 함
//	200<cond<=300 : 물체 (cond%100)이 현재 장소에 있어야 함
//	cond>300      : prop[cond%100] != (cond-300)/100 여야 함
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

// noGoodSkip은 q와 목적지·조건이 같은 명령들을 모두 건너뛴다. (advent.w no_good)
func noGoodSkip(q int) int {
	qq := q
	for q++; caveTravels[q].dest == caveTravels[qq].dest && caveTravels[q].cond == caveTravels[qq].cond; q++ {
	}
	return q
}

// reportInapplicableMotion은 여기서 쓸 수 없는 이동에 대해 사정을 설명한다.
// (advent.w "Report on inapplicable motion")
func (g *Game) reportInapplicableMotion() {
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
			fmt.Fprintf(g.out, "여기선 안과 밖을 모르겠어.  나침반 방향을 쓰거나 가려는 쪽에\n"+
				"있는 걸 말해 줘.")
		case FORWARD, L, R:
			fmt.Fprintf(g.out, "네가 어느 쪽을 보고 있는지 모르겠어.  나침반 방향이나 가까운 물건으로 말해 줘.")
		default:
			fmt.Fprintf(g.out, "그 방향으로는 갈 수 없어.")
		}
	default:
		fmt.Fprintf(g.out, "그 단어를 여기서 어떻게 써야 할지 모르겠어.")
	}
	fmt.Fprintf(g.out, "\n")
}

// choosePloverPassage: 에메랄드만(램프조차 안 됨) 플로버-알코브 통로를 지난다.
// (advent.w "Choose newloc via plover-alcove passage")
func (g *Game) choosePloverPassage() {
	if g.holding == 0 || (g.toting(EMERALD) && g.holding == 1) {
		g.newloc = alcove + proom - g.loc
	} else {
		fmt.Fprintf(g.out, "네가 든 뭔가가 굴을 너랑 같이 통과하지 못해.\n"+
			"소지품을 확인하고 뭔가 내려놓는 게 좋겠어.\n")
		g.newloc = g.loc
	}
}

// crossTrollBridge: 트롤 다리 건너기. 곰을 데리고 있으면 다리가 무너진다.
// (advent.w "Cross troll bridge if possible")
func (g *Game) crossTrollBridge() {
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
		"다리가 휘청해.  넌 필사적으로 붙잡을 곳을\n"+
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
}
