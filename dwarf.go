package main

import "fmt"

// 난쟁이와 해적. 동굴에는 다섯 난쟁이가 떠돌고, 해적(dwarf[0])이 보물을
// 노린다. (advent.w "Dwarf stuff")

const (
	nd           = 5     // 난쟁이 수
	chestLoc     = dead2 // 해적의 보물 상자가 숨는 곳
	messageLoc   = pony  // 해적을 본 뒤 단서가 놓이는 곳
	maxPirateLoc = dead2 // 해적이 갈 수 있는 가장 깊은 곳 (advent.w:1859)
)

// dloc 초기값: 처음엔 어떤 두 난쟁이도 인접하지 않는다. (advent.w:3704)
var dwarfStart = [nd + 1]location{chestLoc, hmk, wfiss, y2, like3, complex}

// dwarf는 난쟁이가 (해적 말고) 현재 장소에 있는지 본다. (advent.w:3712)
func (g *Game) dwarf() bool {
	if g.dflag < 2 {
		return false
	}
	for j := 1; j <= nd; j++ {
		if g.dloc[j] == g.loc {
			return true
		}
	}
	return false
}

// moveDwarves는 당신이 새 장소로 옮긴 직후 다른 이들을 움직인다.
// 해적이 못 가는 곳이거나 다음 이동이 강제면 건너뛴다. (advent.w "Possibly move dwarves and the pirate")
func (g *Game) moveDwarves() {
	if g.loc > maxPirateLoc || g.loc == limbo {
		return
	}
	switch {
	case g.dflag == 0:
		if g.loc >= minLowerLoc {
			g.dflag = 1
		}
	case g.dflag == 1:
		if g.loc >= minLowerLoc && g.pct(5) {
			g.advanceDflag2()
		}
	default:
		g.moveDwarvesAndPirate()
	}
}

// advanceDflag2: 단계 2에 이르면 난쟁이 0~2마리를 조용히 없애고, 한 명이
// 도끼를 던지고 떠난다. (advent.w "Advance dflag to 2")
func (g *Game) advanceDflag2() {
	g.dflag = 2
	for j := 0; j < 2; j++ {
		if g.pct(50) {
			g.dloc[1+g.ran(nd)] = limbo
		}
	}
	for j := 1; j <= nd; j++ {
		if g.dloc[j] == g.loc {
			g.dloc[j] = nugget
		}
		g.odloc[j] = g.dloc[j]
	}
	fmt.Fprintf(g.out, "작은 난쟁이 하나가 모퉁이를 돌다 널 보더니, 작은 도끼를\n"+
		"네게 던지고는, 욕을 내뱉으며 달아났어.  (도끼는 빗나갔어.)\n")
	g.drop(AXE, g.loc)
}

// throwAxeAtDwarf: 도끼로 난쟁이를 공격한다. 2/3 확률로 맞는다. (advent.w "Throw the axe at a dwarf")
func (g *Game) throwAxeAtDwarf() actResult {
	j := 1
	for ; j <= nd; j++ {
		if g.dloc[j] == g.loc {
			break
		}
	}
	if g.ran(3) < 2 {
		g.dloc[j] = limbo
		g.dseen[j] = false
		g.dkill++
		if g.dkill == 1 {
			fmt.Fprintf(g.out, "작은 난쟁이를 죽였어.  사체가 느끼한 검은 연기 구름 속으로\n"+
				"사라져.\n")
		} else {
			fmt.Fprintf(g.out, "작은 난쟁이를 죽였어.\n")
		}
	} else {
		fmt.Fprintf(g.out, "작은 난쟁이를 공격하지만, 녀석이 잽싸게 피해.\n")
	}
	g.drop(AXE, g.loc)
	return g.stayPut()
}

// moveDwarvesAndPirate: 살아 있는 각 난쟁이가 당신을 따라오거나 무작위로
// 움직인다. (advent.w "Move dwarves and the pirate")
func (g *Game) moveDwarvesAndPirate() {
	g.dtotal, g.attack, g.stick = 0, 0, 0
	for j := 0; j <= nd; j++ {
		if g.dloc[j] == limbo {
			continue
		}
		ploc := g.dwarfExits(j)
		if len(ploc) == 0 {
			ploc = []location{g.odloc[j]}
		}
		g.odloc[j] = g.dloc[j]
		g.dloc[j] = ploc[g.ran(len(ploc))] // 무작위 걸음
		g.dseen[j] = g.dloc[j] == g.loc || g.odloc[j] == g.loc ||
			(g.dseen[j] && g.loc >= minLowerLoc)
		if g.dseen[j] {
			g.dwarfFollow(j)
		}
	}
	if g.dtotal != 0 {
		g.dwarvesAttack()
	}
}

// dwarfExits: 난쟁이 j가 갈 수 있는 다음 칸 목록. 무작위 이동용이라
// scan1/2/3을 서로 다른 곳으로 친다. (advent.w "Make a table of all potential exits")
func (g *Game) dwarfExits(j int) []location {
	limit := minForcedLoc - 1
	if j == 0 {
		limit = maxPirateLoc
	}
	var ploc []location
	for q := caveStart[g.dloc[j]]; q < caveStart[g.dloc[j]+1]; q++ {
		nl := caveTravels[q].dest
		if nl >= minLowerLoc && nl != g.odloc[j] && nl != g.dloc[j] &&
			(len(ploc) == 0 || nl != ploc[len(ploc)-1]) && len(ploc) < 19 &&
			caveTravels[q].cond != 100 && nl <= limit {
			ploc = append(ploc, nl)
		}
	}
	return ploc
}

// dwarfFollow: 당신을 본 난쟁이(또는 해적)가 따라온다. (advent.w "Make dwarf j follow")
func (g *Game) dwarfFollow(j int) {
	g.dloc[j] = g.loc
	if j == 0 {
		g.pirateTrack()
		return
	}
	g.dtotal++
	if g.odloc[j] == g.dloc[j] {
		g.attack++
		if g.knifeLoc >= 0 {
			g.knifeLoc = int(g.loc)
		}
		if g.ran(1000) < 95*(g.dflag-2) {
			g.stick++
		}
	}
}

// dwarvesAttack: 위협하는 난쟁이들이 칼을 던진다. (advent.w "Make the threatening dwarves attack")
func (g *Game) dwarvesAttack() {
	if g.dtotal == 1 {
		fmt.Fprintf(g.out, "위협적인 작은 난쟁이가 하나")
	} else {
		fmt.Fprintf(g.out, "위협적인 작은 난쟁이가 %d마리", g.dtotal)
	}
	fmt.Fprintf(g.out, " 너랑 같은 방에 있어!\n")
	if g.attack == 0 {
		return
	}
	if g.dflag == 2 {
		g.dflag = 3
	}
	k := 0
	if g.attack == 1 {
		fmt.Fprintf(g.out, "날카롭고 고약한 칼 하나가 너에게 날아와")
	} else {
		k = 2
		fmt.Fprintf(g.out, " 그중 %d마리가 너에게 칼을 던져", g.attack)
	}
	fmt.Fprintf(g.out, " --- ")
	if g.stick <= 1 {
		fmt.Fprintf(g.out, "%s!\n", attackMsg[k+g.stick])
	} else {
		fmt.Fprintf(g.out, "그중 %d마리가 맞혔어!\n", g.stick)
	}
	if g.stick != 0 {
		g.oldoldloc = g.loc
		g.die()
	}
}

// tooEasy: 이 보물은 줍기 너무 쉬운가(피라미드는 플로버/다크룸에선 쉽다). (advent.w:3880)
func (g *Game) tooEasy(i object) bool {
	return i == PYRAMID && (g.loc == proom || g.loc == droom)
}

// pirateNotSpotted: 아직 해적을 못 봤는가. (advent.w:3879)
func (g *Game) pirateNotSpotted() bool { return g.place[MESSAGE] == limbo }

// pirateTrack: 해적이 보물을 노리며 당신을 쫓는다. (advent.w "Make the pirate track you")
func (g *Game) pirateTrack() {
	if g.loc == maxPirateLoc || g.prop[CHEST] >= 0 {
		return
	}
	k := 0
	for i := minTreasure; i <= maxObj; i++ {
		if !g.tooEasy(i) && g.toting(i) {
			k = -1
			break
		}
		if g.here(i) {
			k = 1
		}
	}
	switch {
	case k < 0:
		g.takeBooty()
	case g.tally == g.lostTreasures+1 && k == 0 && g.pirateNotSpotted() &&
		g.prop[LAMP] != 0 && g.here(LAMP):
		g.pirateSpotted()
	case g.odloc[0] != g.dloc[0] && g.pct(20):
		fmt.Fprintf(g.out, "네 뒤 어둠 속에서 희미하게 바스락거리는 소리가 들려.\n")
	}
}

// takeBooty: 해적이 당신의 보물을 빼앗아 미로의 상자에 숨긴다. (advent.w "Take booty and hide it in the chest")
func (g *Game) takeBooty() {
	fmt.Fprintf(g.out, "네 뒤 그림자 속에서 수염 난 해적이 튀어나와!  \"하르, 하르,\" 그가\n"+
		"낄낄대며 말해. \"이 노획물은 내가 몽땅 챙겨서 미로 깊숙이 내 상자에\n"+
		"숨겨 두지!\"  그러더니 네 보물을 낚아채 어둠 속으로 사라져.\n")
	g.snatchTreasures()
	if g.pirateNotSpotted() {
		g.moveChest()
	}
	g.dloc[0], g.odloc[0] = chestLoc, chestLoc
	g.dseen[0] = false
}

// snatchTreasures: 여기서 챙길 수 있는 보물을 모두 빼앗는다. (advent.w "Snatch all treasures...")
func (g *Game) snatchTreasures() {
	for i := minTreasure; i <= maxObj; i++ {
		if g.tooEasy(i) {
			continue
		}
		if objBase[i] == NOTHING && g.place[i] == g.loc {
			g.carry(i)
		}
		if g.toting(i) {
			g.drop(i, chestLoc)
		}
	}
}

func (g *Game) moveChest() {
	g.move(CHEST, chestLoc)
	g.move(MESSAGE, messageLoc)
}

// pirateSpotted: 보물을 다 봤고 이 방엔 보물이 없을 때 해적을 목격한다. (advent.w "Let the pirate be spotted")
func (g *Game) pirateSpotted() {
	fmt.Fprintf(g.out, "네 뒤 어둠 속에서 희미하게 바스락거리는 소리가 들려.  그쪽으로\n"+
		"고개를 돌리자, 네 램프 불빛이 수염 난 해적을 비춰.  녀석은 큰 상자를\n"+
		"들고 있어.  \"이런 제기랄!\" 그가 외쳐. \"들켜 버렸군!  얼른 미로로\n"+
		"가서 내 상자를 숨겨야겠어!\"  그러더니 어둠 속으로 사라져.\n")
	g.moveChest()
	g.dloc[0], g.odloc[0] = chestLoc, chestLoc
	g.dseen[0] = false
}

// checkInterference: newloc으로 가려는데 폐쇄나 난쟁이가 막는지 본다.
// (advent.w "Check for interference with the proposed move")
func (g *Game) checkInterference() {
	if g.closing() && g.newloc < minInCave && g.newloc != limbo {
		g.panicClosing()
		g.newloc = g.loc
	} else if g.newloc != g.loc {
		g.dwarfBlock()
	}
}

// dwarfBlock: 칼 든 난쟁이가 가려는 길을 막으면 제자리에 머문다. (advent.w "Stay in loc if a dwarf is blocking...")
func (g *Game) dwarfBlock() {
	if g.loc <= maxPirateLoc {
		for j := 1; j <= nd; j++ {
			if g.odloc[j] == g.newloc && g.dseen[j] {
				fmt.Fprintf(g.out, "큰 칼을 든 작은 난쟁이가 네 길을 막아.\n")
				g.newloc = g.loc
				break
			}
		}
	}
}
