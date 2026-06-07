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
	return g.rep("이제 물병이 비었어.")
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
		return g.rep("그건 부을 수 없어.")
	}
	g.prop[BOTTLE] = 1
	g.place[g.obj] = limbo
	if g.loc == g.place[PLANT] {
		return g.waterPlant()
	}
	if g.loc == g.place[DOOR] {
		return g.pourOnDoor()
	}
	return g.rep("병이 비고 바닥이 젖었어.")
}

func (g *Game) waterPlant() actResult {
	if g.obj != WATER {
		return g.rep("식물이 발끈해서 잎의 기름을 털어내며 묻네.  \"물은?\"")
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
		return g.rep("이제 경첩이 완전히 녹슬어서 꿈쩍도 안 해.")
	case OIL:
		g.prop[DOOR] = 1
		return g.rep("기름이 경첩을 풀어줘서 이제 문이 열려.")
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
		return g.rep("병은 이미 가득 찼어.")
	}
	if g.noLiquidHere() {
		return g.rep("여기엔 병을 채울 만한 게 없어.")
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
		fmt.Fprintf(g.out, "이제 병에 기름이 가득해.\n")
	} else {
		fmt.Fprintf(g.out, "이제 병에 물이 가득해.\n")
	}
	return aDone()
}

func (g *Game) fillVase() actResult {
	if g.noLiquidHere() {
		return g.rep("여기엔 꽃병을 채울 만한 게 없어.\n")
	}
	if !g.toting(VASE) {
		return g.defTo(DROP)
	}
	fmt.Fprintf(g.out, "급격한 온도 변화에 꽃병이 섬세하게 산산조각 났어.\n")
	return g.smashVase()
}

// ---- TAKE ----

func (g *Game) doTake() actResult {
	if g.toting(g.obj) {
		return g.repDefault() // 이미 들고 있음
	}
	if objBase[g.obj] != NOTHING { // 움직일 수 없는 물건
		if g.obj == CHAIN && g.prop[BEAR] != 0 {
			return g.rep("사슬은 아직 잠겨 있어.")
		}
		if g.obj == BEAR && g.prop[BEAR] == 1 {
			return g.rep("곰은 아직 벽에 사슬로 묶여 있어.")
		}
		if g.obj == PLANT && g.prop[PLANT] <= 0 {
			return g.rep("식물은 뿌리가 유난히 깊어서 뽑아낼 수가 없어.")
		}
		return g.rep("설마 진심은 아니겠지!")
	}
	if g.obj == WATER || g.obj == OIL {
		if r, done := g.takeLiquid(); done {
			return r
		}
	}
	if g.holding >= 7 {
		return g.rep("더는 못 들어.  먼저 뭔가 내려놔야 해.")
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
	return g.rep("그걸 담아 갈 게 없어."), true
}

// takeBird: 새는 막대를 들고 있으면 못 잡고, 새장이 있어야 데려간다.
// (advent.w "Check special cases for taking a bird")
func (g *Game) takeBird() (actResult, bool) {
	if g.toting(ROD) {
		return g.rep("네가 들어왔을 땐 새가 겁내지 않았는데, 다가가니까\n" +
			"불안해해서 잡을 수가 없어."), true
	}
	if g.toting(CAGE) {
		g.prop[BIRD] = 1
	} else {
		return g.rep("새를 잡을 순 있지만, 데려갈 순 없어."), true
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
		fmt.Fprintf(g.out, "작은 새가 초록 뱀을 공격해, 놀랍도록 격렬하게 퍼덕이며\n"+
			"뱀을 쫓아 버려.\n")
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
		return g.rep("작은 새가 초록 용을 공격해, 놀랍도록 격렬하게 퍼덕이다가\n" +
			"잿더미가 되도록 타 버려.  재가 바람에 흩날려 가."), true
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
	fmt.Fprintf(g.out, "곰이 트롤 쪽으로 어슬렁어슬렁 다가가자, 트롤은 깜짝 놀라 비명을\n"+
		"지르며 잽싸게 달아나.  곰은 곧 추격을 포기하고 어슬렁어슬렁 돌아와.\n")
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
		fmt.Fprintf(g.out, "도끼가 용의 두꺼운 비늘에 맞고 힘없이 튕겨 나가.\n")
	} else if g.isAtLoc(TROLL) {
		fmt.Fprintf(g.out, "트롤이 도끼를 날렵하게 받아 꼼꼼히 살펴보더니, 다시 던지며\n"+
			"이렇게 말해.  \"솜씨는 좋군, 한데 값어치가 영 모자라.\"\n")
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
	return g.rep("트롤이 네 보물을 낚아채 시야 밖으로 잽싸게 달아나.")
}

func (g *Game) throwAxeAtBear() actResult {
	g.drop(AXE, g.loc)
	g.prop[AXE] = 1
	objBase[AXE] = AXE
	if g.place[BEAR] == g.loc {
		g.move(BEAR, g.loc)
	}
	return g.rep("도끼가 빗나가 곰 옆에 떨어져서, 네가 닿을 수 없는 곳에 가 버렸어.")
}
