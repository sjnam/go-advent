package main

import "fmt"

// 동굴 폐쇄와 램프 관리. 모든 보물을 찾으면 clock1이 똑딱이기 시작하고,
// 0이 되면 폐쇄가 시작되며 clock2가 흐른다. (advent.w "Closing the cave")

// checkClocks는 매 턴 시계와 램프를 점검한다. 동굴이 완전히 닫히면
// true(제자리로)를 돌려준다. (advent.w "Check the clocks and the lamp")
func (g *Game) checkClocks() bool {
	if g.tally == 0 && g.loc >= minLowerLoc && g.loc != y2 {
		g.clock1--
	}
	if g.clock1 == 0 {
		g.warnClosing()
		return false
	}
	if g.clock1 < 0 {
		g.clock2--
	}
	if g.clock2 == 0 {
		g.closeCave()
		return true // stay_put
	}
	g.checkLamp()
	return false
}

// zapLampIfElusive: 더는 보물을 찾을 수 없게 됐다면 램프 수명을 줄인다.
// (advent.w "Zap the lamp if the remaining treasures are too elusive")
func (g *Game) zapLampIfElusive() {
	if g.tally == g.lostTreasures && g.tally > 0 && g.limit > 35 {
		g.limit = 35
	}
}

// checkLamp는 램프 전력을 확인한다. (advent.w "Check the lamp")
func (g *Game) checkLamp() {
	if g.prop[LAMP] == 1 {
		g.limit--
	}
	switch {
	case g.limit <= 30 && g.here(BATTERIES) && g.prop[BATTERIES] == 0 && g.here(LAMP):
		g.replaceBatteries()
	case g.limit == 0:
		g.extinguishLamp()
	case g.limit < 0 && g.loc < minInCave:
		fmt.Fprintf(g.out, "여기 밖에서 어슬렁거려 봐야 별 의미 없고, 램프 없이는\n"+
			"동굴을 탐험할 수도 없어.  그냥 오늘은 여기까지 하자.\n")
		g.gaveUp = true
		g.quitting = true
	case g.limit <= 30 && !g.warned && g.here(LAMP):
		fmt.Fprintf(g.out, "램프가 어두워지고 있어")
		switch {
		case g.prop[BATTERIES] == 1:
			fmt.Fprintf(g.out, ", 게다가 여분 배터리도 없어.  슬슬\n"+
				"마무리하는 게 좋겠어.\n")
		case g.place[BATTERIES] == limbo:
			fmt.Fprintf(g.out, ".  슬슬 마무리하는 게 좋겠어, 새 배터리를\n"+
				"찾지 못한다면 말이야.  내 기억에 미로 어딘가에\n"+
				"자판기가 있었어.  동전을 좀 챙겨 가.\n")
		default:
			fmt.Fprintf(g.out, ".  그 배터리를 가지러 돌아가는 게 좋겠어.\n")
		}
		g.warned = true
	}
}

func (g *Game) replaceBatteries() {
	fmt.Fprintf(g.out, "램프가 어두워지고 있어.  내가 알아서 배터리를\n"+
		"갈아 끼울게.\n")
	g.prop[BATTERIES] = 1
	if g.toting(BATTERIES) {
		g.drop(BATTERIES, g.loc)
	}
	g.limit = 2500
}

func (g *Game) extinguishLamp() {
	g.limit = -1
	g.prop[LAMP] = 0
	if g.here(LAMP) {
		fmt.Fprintf(g.out, "램프 전력이 다 떨어졌어.")
	}
}

// warnClosing: 첫 경고. 격자문을 잠그고, 수정 다리를 부수고, 난쟁이와
// 해적을 모두 없애고, 트롤·곰을 치우고 폐쇄를 시작한다. (advent.w "Warn that the cave is closing")
func (g *Game) warnClosing() {
	fmt.Fprintf(g.out, "동굴 전체에 울려 퍼지는 음산한 목소리가 말해.  \"곧 동굴이\n"+
		"닫힙니다.  모든 모험가는 즉시 본부를 통해 나가십시오.\"\n")
	g.clock1 = -1
	g.prop[GRATE] = 0
	g.prop[CRYSTAL] = 0
	for j := 0; j <= nd; j++ {
		g.dseen[j] = false
		g.dloc[j] = limbo
	}
	g.destroy(TROLL)
	g.destroy(TROLL_)
	g.move(TROLL2, swside)
	g.move(TROLL2_, neside)
	g.move(BRIDGE, swside)
	g.move(BRIDGE_, neside)
	if g.prop[BEAR] != 3 {
		g.destroy(BEAR)
	}
	g.prop[CHAIN] = 0
	objBase[CHAIN] = NOTHING
	g.prop[AXE] = 0
	objBase[AXE] = NOTHING
}

// panicClosing: 폐쇄 중 밖으로 나가려 하면 몇 턴 더 준다. (advent.w "Panic at closing time")
func (g *Game) panicClosing() {
	if !g.panicked {
		g.clock2 = 15
		g.panicked = true
	}
	fmt.Fprintf(g.out, "수상한 녹음된 목소리가 끼익 살아나며 알려:\n"+
		"\"이 출구는 닫혔습니다.  본부를 통해 나가 주십시오.\"\n")
}

// closeCave: clock2가 0이 되면 최종 퍼즐의 보관실로 옮긴다. (advent.w "Close the cave")
func (g *Game) closeCave() {
	fmt.Fprintf(g.out, "음산한 목소리가 읊조려.  \"동굴은 이제 닫혔다.\"  메아리가 잦아들자,\n"+
		"눈부신 빛이 번쩍이고 (주황색 연기도 작게 한 줄기 피어올라).\n"+
		". . .    그러다 다시 초점이 잡히고, 둘러보니...\n")
	g.move(BOTTLE, neend)
	g.prop[BOTTLE] = -2
	g.move(PLANT, neend)
	g.prop[PLANT] = -1
	g.move(OYSTER, neend)
	g.prop[OYSTER] = -1
	g.move(LAMP, neend)
	g.prop[LAMP] = -1
	g.move(ROD, neend)
	g.prop[ROD] = -1
	g.move(DWARF, neend)
	g.prop[DWARF] = -1
	g.move(MIRROR, neend)
	g.prop[MIRROR] = -1
	g.loc, g.oldloc = neend, neend
	g.move(GRATE, swend) // prop[GRATE]는 여전히 0
	g.move(SNAKE, swend)
	g.prop[SNAKE] = -2
	g.move(BIRD, swend)
	g.prop[BIRD] = -2
	g.move(CAGE, swend)
	g.prop[CAGE] = -1
	g.move(ROD2, swend)
	g.prop[ROD2] = -1
	g.move(PILLOW, swend)
	g.prop[PILLOW] = -1
	g.move(MIRROR_, swend)
	g.place[WATER] = limbo // bugfix
	g.place[OIL] = limbo
	for j := object(1); j <= maxObj; j++ {
		if g.toting(j) {
			g.destroy(j)
		}
	}
	g.closed = true
	g.bonus = 10
}
