package main

import "fmt"

// 액체(물·기름·병)와, 그것에 얽힌 집기/놓기/던지기. (advent.w "Liquid assets",
// "The other actions"의 DROP/TOSS)

// bottleEmpty: 병이 비었는가. (advent.w 매크로 bottle_empty)
func (g *Game) bottleEmpty() bool { return g.prop[BOTTLE] == 1 || g.prop[BOTTLE] < 0 }

// objectInBottle: 지금 obj가 병 안에 든 그 액체인가. (advent.w 매크로 object_in_bottle)
func (g *Game) objectInBottle() bool {
	return (g.obj == WATER && g.prop[BOTTLE] == 0) || (g.obj == OIL && g.prop[BOTTLE] == 2)
}

// ---- DRINK / POUR / FILL ----

func (g *Game) doDrink() actResult {
	if g.obj == NOTHING {
		if !g.waterHere() && !(g.here(BOTTLE) && g.prop[BOTTLE] == 0) {
			return g.getObjR()
		}
	} else if g.obj != WATER {
		return g.defTo(EAT)
	}
	if !(g.here(BOTTLE) && g.prop[BOTTLE] == 0) {
		return g.repDefault()
	}
	g.prop[BOTTLE] = 1
	g.place[WATER] = limbo
	return g.rep("The bottle of water is now empty.")
}

func (g *Game) doPour() actResult {
	if g.obj == NOTHING || g.obj == BOTTLE {
		switch g.prop[BOTTLE] {
		case 0:
			g.obj = WATER
		case 2:
			g.obj = OIL
		default:
			g.obj = NOTHING
		}
		if g.obj == NOTHING {
			return g.getObjR()
		}
	}
	if !g.toting(g.obj) {
		return g.repDefault()
	}
	if g.obj != WATER && g.obj != OIL {
		return g.rep("You can't pour that.")
	}
	g.prop[BOTTLE] = 1
	g.place[g.obj] = limbo
	if g.loc == g.place[PLANT] {
		return g.waterPlant()
	}
	if g.loc == g.place[DOOR] {
		return g.pourOnDoor()
	}
	return g.rep("Your bottle is empty and the ground is wet.")
}

func (g *Game) waterPlant() actResult {
	if g.obj != WATER {
		return g.rep("The plant indignantly shakes the oil off its leaves and asks, \"Water?\"")
	}
	fmt.Fprintf(g.out, "%s\n", objNote[g.prop[PLANT]+1+objOffset[PLANT]])
	g.prop[PLANT] += 2
	if g.prop[PLANT] > 4 {
		g.prop[PLANT] = 0
	}
	g.prop[PLANT2] = g.prop[PLANT] >> 1
	return g.stayPut()
}

func (g *Game) pourOnDoor() actResult {
	switch g.obj {
	case WATER:
		g.prop[DOOR] = 0
		return g.rep("The hinges are quite thoroughly rusted now and won't budge.")
	case OIL:
		g.prop[DOOR] = 1
		return g.rep("The oil has freed up the hinges so that the door will now open.")
	}
	return aDone()
}

func (g *Game) doFill() actResult {
	if g.obj == VASE {
		return g.fillVase()
	}
	if !g.here(BOTTLE) {
		if g.obj == NOTHING {
			return g.getObjR()
		}
		return g.repDefault()
	} else if g.obj != NOTHING && g.obj != BOTTLE {
		return g.repDefault()
	}
	if !g.bottleEmpty() {
		return g.rep("Your bottle is already full.")
	}
	if g.noLiquidHere() {
		return g.rep("There is nothing here with which to fill the bottle.")
	}
	g.prop[BOTTLE] = caveFlags[g.loc] & oil
	if g.toting(BOTTLE) {
		if g.prop[BOTTLE] != 0 {
			g.place[OIL] = inhand
		} else {
			g.place[WATER] = inhand
		}
	}
	if g.prop[BOTTLE] != 0 {
		fmt.Fprintf(g.out, "Your bottle is now full of oil.\n")
	} else {
		fmt.Fprintf(g.out, "Your bottle is now full of water.\n")
	}
	return aDone()
}

func (g *Game) fillVase() actResult {
	if g.noLiquidHere() {
		return g.rep("There is nothing here with which to fill the vase.\n")
	}
	if !g.toting(VASE) {
		return g.defTo(DROP)
	}
	fmt.Fprintf(g.out, "The sudden change in temperature has delicately shattered the vase.\n")
	return g.smashVase()
}

// ---- TAKE ----

func (g *Game) doTake() actResult {
	if g.toting(g.obj) {
		return g.repDefault() // 이미 들고 있음
	}
	if objBase[g.obj] != NOTHING { // 움직일 수 없는 물건
		if g.obj == CHAIN && g.prop[BEAR] != 0 {
			return g.rep("The chain is still locked.")
		}
		if g.obj == BEAR && g.prop[BEAR] == 1 {
			return g.rep("The bear is still chained to the wall.")
		}
		if g.obj == PLANT && g.prop[PLANT] <= 0 {
			return g.rep("The plant has exceptionally deep roots and cannot be pulled free.")
		}
		return g.rep("You can't be serious!")
	}
	if g.obj == WATER || g.obj == OIL {
		if r, done := g.takeLiquid(); done {
			return r
		}
	}
	if g.holding >= 7 {
		return g.rep("You can't carry anything more.  You'll have to drop something first.")
	}
	if g.obj == BIRD && g.prop[BIRD] == 0 {
		if r, done := g.takeBird(); done {
			return r
		}
	}
	if g.obj == BIRD || (g.obj == CAGE && g.prop[BIRD] != 0) {
		g.carry(BIRD + CAGE - g.obj)
	}
	g.carry(g.obj)
	if g.obj == BOTTLE && !g.bottleEmpty() {
		if g.prop[BOTTLE] != 0 {
			g.place[OIL] = inhand
		} else {
			g.place[WATER] = inhand
		}
	}
	return g.defTo(RELAX) // OK, 집었다
}

// takeLiquid: 액체를 집으려면 병이 있어야 한다. (advent.w "Check special cases for taking a liquid")
// done=false면 obj=BOTTLE로 바꾼 뒤 계속 집는다.
func (g *Game) takeLiquid() (actResult, bool) {
	if g.here(BOTTLE) && g.objectInBottle() {
		g.obj = BOTTLE
		return actResult{}, false
	}
	g.obj = BOTTLE
	if g.toting(BOTTLE) {
		return aChange(FILL), true
	}
	return g.rep("You have nothing in which to carry it."), true
}

// takeBird: 새는 막대를 들고 있으면 못 잡고, 새장이 있어야 데려간다.
// (advent.w "Check special cases for taking a bird")
func (g *Game) takeBird() (actResult, bool) {
	if g.toting(ROD) {
		return g.rep("The bird was unafraid when you entered, but as you approach it becomes\n" +
			"disturbed and you cannot catch it."), true
	}
	if g.toting(CAGE) {
		g.prop[BIRD] = 1
	} else {
		return g.rep("You can catch the bird, but you cannot carry it."), true
	}
	return actResult{}, false
}

// ---- DROP ----

func (g *Game) doDrop() actResult {
	if g.obj == ROD && g.toting(ROD2) && !g.toting(ROD) {
		g.obj = ROD2
	}
	if !g.toting(g.obj) {
		return g.repDefault()
	}
	k := false // "OK" 메시지를 누를지
	if g.obj == COINS && g.here(PONY) {
		return g.putCoins()
	}
	if g.obj == BIRD {
		if r, done := g.dropBird(&k); done {
			return r
		}
	}
	if g.obj == VASE && g.loc != soft {
		g.dropVase(&k)
	}
	if g.obj == BEAR && g.isAtLoc(TROLL) {
		g.chaseTroll(&k)
	}
	g.dropLiquid()
	if g.obj == BIRD {
		g.prop[BIRD] = 0
	} else if g.obj == CAGE && g.prop[BIRD] != 0 {
		g.drop(BIRD, g.loc)
	}
	g.drop(g.obj, g.loc)
	if k {
		return aDone()
	}
	return g.defTo(RELAX)
}

func (g *Game) putCoins() actResult {
	g.destroy(COINS)
	g.drop(BATTERIES, g.loc)
	g.prop[BATTERIES] = 0
	return g.rep(objNote[objOffset[BATTERIES]])
}

func (g *Game) dropBird(k *bool) (actResult, bool) {
	if g.here(SNAKE) {
		fmt.Fprintf(g.out, "The little bird attacks the green snake, and in an astounding flurry\n"+
			"drives the snake away.\n")
		*k = true
		if g.closed {
			return g.dwarvesUpset(), true
		}
		g.destroy(SNAKE)
		g.prop[SNAKE] = 1
		return actResult{}, false
	}
	if g.isAtLoc(DRAGON) && g.prop[DRAGON] == 0 {
		g.destroy(BIRD)
		g.prop[BIRD] = 0
		if g.place[SNAKE] == hmk {
			g.lostTreasures++
		}
		return g.rep("The little bird attacks the green dragon, and in an astounding flurry\n" +
			"gets burnt to a cinder.  The ashes blow away."), true
	}
	return actResult{}, false
}

func (g *Game) dropVase(k *bool) {
	if g.place[PILLOW] == g.loc {
		g.prop[VASE] = 0
	} else {
		g.prop[VASE] = 2
	}
	fmt.Fprintf(g.out, "%s\n", objNote[objOffset[VASE]+1+g.prop[VASE]])
	*k = true
	if g.prop[VASE] != 0 {
		objBase[VASE] = VASE
	}
}

func (g *Game) chaseTroll(k *bool) {
	fmt.Fprintf(g.out, "The bear lumbers toward the troll, who lets out a startled shriek and\n"+
		"scurries away.  The bear soon gives up the pursuit and wanders back.\n")
	*k = true
	g.destroy(TROLL)
	g.destroy(TROLL_)
	g.drop(TROLL2, swside)
	g.drop(TROLL2_, neside)
	g.prop[TROLL] = 2
	g.move(BRIDGE, swside)
	g.move(BRIDGE_, neside)
}

// dropLiquid: 병을 놓으면 안의 액체도 함께 떨어진다. (advent.w "Check special cases for dropping a liquid")
func (g *Game) dropLiquid() {
	if g.objectInBottle() {
		g.obj = BOTTLE
	}
	if g.obj == BOTTLE && !g.bottleEmpty() {
		if g.prop[BOTTLE] != 0 {
			g.place[OIL] = limbo
		} else {
			g.place[WATER] = limbo
		}
	}
}

// ---- TOSS (throw) ----

func (g *Game) doToss() actResult {
	if g.obj == ROD && g.toting(ROD2) && !g.toting(ROD) {
		g.obj = ROD2
	}
	if !g.toting(g.obj) {
		return g.repDefault()
	}
	if isTreasure(g.obj) && g.isAtLoc(TROLL) {
		return g.snarfTreasure()
	}
	if g.obj == FOOD && g.here(BEAR) {
		g.obj = BEAR
		return aChange(FEED)
	}
	if g.obj != AXE {
		return aChange(DROP)
	}
	if g.dwarf() {
		return g.throwAxeAtDwarf()
	}
	if g.isAtLoc(DRAGON) && g.prop[DRAGON] == 0 {
		fmt.Fprintf(g.out, "The axe bounces harmlessly off the dragon's thick scales.\n")
	} else if g.isAtLoc(TROLL) {
		fmt.Fprintf(g.out, "The troll deftly catches the axe, examines it carefully, and tosses it\n"+
			"back, declaring, \"Good workmanship, but it's not valuable enough.\"\n")
	} else if g.here(BEAR) && g.prop[BEAR] == 0 {
		return g.throwAxeAtBear()
	} else {
		g.obj = NOTHING
		return aChange(KILL)
	}
	g.drop(AXE, g.loc)
	return g.stayPut()
}

func (g *Game) snarfTreasure() actResult {
	g.drop(g.obj, limbo)
	g.destroy(TROLL)
	g.destroy(TROLL_)
	g.drop(TROLL2, swside)
	g.drop(TROLL2_, neside)
	g.move(BRIDGE, swside)
	g.move(BRIDGE_, neside)
	return g.rep("The troll catches your treasure and scurries away out of sight.")
}

func (g *Game) throwAxeAtBear() actResult {
	g.drop(AXE, g.loc)
	g.prop[AXE] = 1
	objBase[AXE] = AXE
	if g.place[BEAR] == g.loc {
		g.move(BEAR, g.loc)
	}
	return g.rep("The axe misses and lands near the bear where you can't get at it.")
}
