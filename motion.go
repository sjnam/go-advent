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
			fmt.Fprintf(g.out, "Sorry, but I am not allowed to give more detail.  I will repeat the\n"+
				"long description of your location.\n")
		}
		g.wasDark = false
		g.visits[g.loc] = 0
		return smStay
	case CAVE:
		if g.loc < minInCave {
			fmt.Fprintf(g.out, "I can't see where the cave is, but hereabouts no stream can run on\n"+
				"the surface for long.  I would try the stream.\n")
		} else {
			fmt.Fprintf(g.out, "I need more detailed instructions to do that.\n")
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
		fmt.Fprintf(g.out, "Sorry, but I no longer seem to remember how you got here.\n")
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
			fmt.Fprintf(g.out, "You can't get there from here.\n")
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
		fmt.Fprintf(g.out, "Which way?")
	case g.mot == XYZZY || g.mot == PLUGH:
		fmt.Fprintf(g.out, "%s", g.defaultMsg[WAVE])
	case g.verb == FIND || g.verb == INVENTORY:
		fmt.Fprintf(g.out, "%s", g.defaultMsg[FIND])
	case g.mot <= FORWARD:
		switch g.mot {
		case IN, OUT:
			fmt.Fprintf(g.out, "I don't know in from out here.  Use compass points or name something\n"+
				"in the general direction you want to go.")
		case FORWARD, L, R:
			fmt.Fprintf(g.out, "I am unsure how you are facing.  Use compass points or nearby objects.")
		default:
			fmt.Fprintf(g.out, "There is no way to go in that direction.")
		}
	default:
		fmt.Fprintf(g.out, "I don't know how to apply that word here.")
	}
	fmt.Fprintf(g.out, "\n")
}

// choosePloverPassage: 에메랄드만(램프조차 안 됨) 플로버-알코브 통로를 지난다.
// (advent.w "Choose newloc via plover-alcove passage")
func (g *Game) choosePloverPassage() {
	if g.holding == 0 || (g.toting(EMERALD) && g.holding == 1) {
		g.newloc = alcove + proom - g.loc
	} else {
		fmt.Fprintf(g.out, "Something you're carrying won't fit through the tunnel with you.\n"+
			"You'd best take inventory and drop something.\n")
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
	fmt.Fprintf(g.out, "Just as you reach the other side, the bridge buckles beneath the\n"+
		"weight of the bear, who was still following you around.  You\n"+
		"scrabble desperately for support, but as the bridge collapses you\n"+
		"stumble back and fall into the chasm.\n")
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
