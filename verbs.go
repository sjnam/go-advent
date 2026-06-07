package main

import "fmt"

// 싸움·먹이기·열고닫기·읽기·말하기 동사. (advent.w "The other actions")

// ---- KILL (attack) ----

func (g *Game) doKill() actResult {
	if g.obj == NOTHING {
		if r, redirect := g.uniqueAttack(); redirect {
			return r
		}
	}
	switch g.obj {
	case NOTHING:
		return g.rep("여기엔 공격할 게 없어.")
	case BIRD:
		return g.dispatchBird()
	case DRAGON:
		if g.prop[DRAGON] == 0 {
			return g.funStuffDragon()
		}
		return g.rep("맙소사, 그 불쌍한 녀석은 이미 죽었어!")
	case CLAM, OYSTER:
		return g.rep("껍데기가 아주 단단해서 공격이 안 통해.")
	case SNAKE:
		return g.rep("뱀을 공격하는 건 소용도 없고 아주 위험해.")
	case DWARF:
		if g.closed {
			return g.dwarvesUpset()
		}
		return g.rep("뭘로?  맨손으로?")
	case TROLL:
		return g.rep("트롤은 바위와 가까운 친척이라 살갗이\n" +
			"코뿔소 가죽처럼 질겨.  트롤은 네 공격을 가뿐히 막아내.")
	case BEAR:
		switch g.prop[BEAR] {
		case 0:
			return g.rep("뭘로?  맨손으로?  녀석의 곰 같은 손을 상대로?")
		case 3:
			return g.rep("맙소사, 그 불쌍한 녀석은 이미 죽었어!")
		default:
			return g.rep("곰이 어리둥절해해.  그냥 네 친구가 되고 싶을 뿐이야.")
		}
	default:
		return g.repDefault()
	}
}

// uniqueAttack: 공격할 대상이 하나로 정해지는지 본다. 둘 이상이면 되묻는다.
// (advent.w "See if there's a unique object to attack")
func (g *Game) uniqueAttack() (actResult, bool) {
	k := 0
	if g.dwarf() {
		k++
		g.obj = DWARF
	}
	if g.here(SNAKE) {
		k++
		g.obj = SNAKE
	}
	if g.isAtLoc(DRAGON) && g.prop[DRAGON] == 0 {
		k++
		g.obj = DRAGON
	}
	if g.isAtLoc(TROLL) {
		k++
		g.obj = TROLL
	}
	if g.here(BEAR) && g.prop[BEAR] == 0 {
		k++
		g.obj = BEAR
	}
	if k == 0 {
		if g.here(BIRD) && g.oldverb != TOSS {
			k++
			g.obj = BIRD
		}
		if g.here(CLAM) || g.here(OYSTER) {
			k++
			g.obj = CLAM
		}
	}
	if k > 1 {
		return g.getObjR(), true
	}
	return actResult{}, false
}

func (g *Game) dispatchBird() actResult {
	if g.closed {
		return g.rep("아, 그 불쌍한 새는 좀 내버려 둬.")
	}
	g.destroy(BIRD)
	g.prop[BIRD] = 0
	if g.place[SNAKE] == hmk {
		g.lostTreasures++
	}
	return g.rep("작은 새가 죽었어.  사체가 사라져.")
}

// funStuffDragon: 맨손으로 용을 공격하겠다고 우기면 용이 죽는다. (advent.w "Fun stuff for dragon")
func (g *Game) funStuffDragon() actResult {
	fmt.Fprintf(g.out, "뭘로?  맨손으로?\n")
	g.verb = ABSTAIN
	g.obj = NOTHING
	if !g.listen() {
		g.quitting = true
		return aDone()
	}
	if !(streq(g.word1, "예") || streq(g.word1, "응") || streq(g.word1, "y")) {
		// TODO: 정확히는 이 입력을 명령으로 재처리(goto pre_parse).
		return aDone()
	}
	fmt.Fprintf(g.out, "%s\n", objNote[objOffset[DRAGON]+1])
	g.prop[DRAGON] = 2 // 죽음
	g.prop[RUG] = 0
	objBase[RUG] = NOTHING // 이제 쓸 수 있는 보물
	objBase[DRAGON_] = DRAGON_
	g.destroy(DRAGON_)
	objBase[RUG_] = RUG_
	g.destroy(RUG_)
	for t := object(1); t <= maxObj; t++ {
		if g.place[t] == scan1 || g.place[t] == scan3 {
			g.move(t, scan2)
		}
	}
	g.loc = scan2
	return g.stayPut()
}

// ---- FEED ----

func (g *Game) doFeed() actResult {
	switch g.obj {
	case BIRD:
		return g.rep("배가 안 고파 (그냥 피오르를 그리워하는 것뿐이야).  게다가 넌\n" +
			"새 모이도 없잖아.")
	case TROLL:
		return g.rep("탐식은 트롤의 악덕이 아니야.  하지만 탐욕은 맞지.")
	case DRAGON:
		if g.prop[DRAGON] != 0 {
			return g.rep(g.defaultMsg[EAT])
		}
	case SNAKE:
		if !g.closed && g.here(BIRD) {
			g.destroy(BIRD)
			g.prop[BIRD] = 0
			g.lostTreasures++
			return g.rep("뱀이 네 새를 삼켜 버렸어.")
		}
	case BEAR:
		if !g.here(FOOD) {
			if g.prop[BEAR] == 0 {
				break
			}
			if g.prop[BEAR] == 3 {
				g.verb = EAT
			}
			return g.repDefault()
		}
		g.destroy(FOOD)
		g.prop[BEAR] = 1
		g.prop[AXE] = 0
		objBase[AXE] = NOTHING // 도끼를 다시 들 수 있다
		return g.rep("곰이 네 음식을 게걸스레 먹어 치우더니, 그러고 나서 한결\n" +
			"누그러지고 심지어 꽤 친근해지기까지 해.")
	case DWARF:
		if !g.here(FOOD) {
			return g.repDefault()
		}
		g.dflag++
		return g.rep("이 바보야, 난쟁이는 석탄만 먹어!  이제 녀석을 제대로 화나게 했어!")
	default:
		return g.rep(g.defaultMsg[CALM])
	}
	return g.rep("여기엔 그게 먹고 싶어 할 만한 게 없어 (너라면 모를까).")
}

// ---- OPEN / CLOSE ----

func (g *Game) doOpenClose() actResult {
	switch g.obj {
	case OYSTER, CLAM:
		return g.openClam()
	case GRATE, CHAIN:
		if !g.here(KEYS) {
			return g.rep("열쇠가 없잖아!")
		}
		return g.openGrateChain()
	case KEYS:
		return g.rep("열쇠를 잠그거나 열 순 없어.")
	case CAGE:
		return g.rep("그건 자물쇠가 없어.")
	case DOOR:
		if g.prop[DOOR] != 0 {
			return g.defTo(RELAX)
		}
		return g.rep("문은 몹시 녹슬어서 열리지 않아.")
	default:
		return g.repDefault()
	}
}

func (g *Game) openGrateChain() actResult {
	if g.obj == CHAIN {
		return g.openCloseChain()
	}
	if g.closing() {
		g.panicClosing()
		return aDone()
	}
	k := g.prop[GRATE]
	if g.verb == OPEN {
		g.prop[GRATE] = 1
	} else {
		g.prop[GRATE] = 0
	}
	switch k + 2*g.prop[GRATE] {
	case 0:
		return g.rep("이미 잠겨 있었어.")
	case 1:
		return g.rep("이제 창살이 잠겼어.")
	case 2:
		return g.rep("이제 창살이 열렸어.")
	default: // 3
		return g.rep("이미 열려 있었어.")
	}
}

func (g *Game) openCloseChain() actResult {
	if g.verb == OPEN {
		return g.openChain()
	}
	if g.loc != barr {
		return g.rep("여기엔 사슬을 채울 만한 게 없어.")
	}
	if g.prop[CHAIN] != 0 {
		return g.rep("이미 잠겨 있었어.")
	}
	g.prop[CHAIN] = 2
	objBase[CHAIN] = CHAIN
	if g.toting(CHAIN) {
		g.drop(CHAIN, g.loc)
	}
	return g.rep("이제 사슬이 잠겼어.")
}

func (g *Game) openChain() actResult {
	if g.prop[CHAIN] == 0 {
		return g.rep("이미 열려 있었어.")
	}
	if g.prop[BEAR] == 0 {
		return g.rep("곰을 지나쳐 사슬을 풀 방법이 없어, 뭐\n" +
			"어쩌면 다행이지만.")
	}
	g.prop[CHAIN] = 0
	objBase[CHAIN] = NOTHING // 사슬이 풀렸다
	if g.prop[BEAR] == 3 {
		objBase[BEAR] = BEAR
	} else {
		g.prop[BEAR] = 2
		objBase[BEAR] = NOTHING
	}
	return g.rep("이제 사슬이 풀렸어.")
}

func (g *Game) openClam() actResult {
	name := "굴"
	if g.obj == CLAM {
		name = "대합"
	}
	if g.verb == CLOSE {
		return g.rep("뭐라고?")
	}
	if !g.toting(TRIDENT) {
		fmt.Fprintf(g.out, "%s을 열 만큼 강한 게 없어", name)
		return g.rep(".")
	}
	if g.toting(g.obj) {
		fmt.Fprintf(g.out, "%s을 열기 전에 내려놓는 게 좋겠어.  ", name)
		if g.obj == CLAM {
			return g.rep(">끙!<")
		}
		return g.rep(">으드득!<")
	}
	if g.obj == CLAM {
		g.destroy(CLAM)
		g.drop(OYSTER, g.loc)
		g.drop(PEARL, sac)
		return g.rep("반짝이는 진주가 대합에서 굴러 나와 떼구루루 굴러가.  세상에,\n" +
			"이거 진짜 굴이었나 봐.  (난 조개 종류를 구분하는 데\n" +
			"영 소질이 없었어.)  뭐가 됐든, 이제 다시 딱 닫혀 버렸어.")
	}
	return g.rep("굴이 삐걱 열리는데, 안에는 굴 말고 아무것도 없어.\n" +
		"곧바로 다시 딱 닫혀.")
}

// ---- READ ----

func (g *Game) doRead() actResult {
	if g.dark() {
		return g.cantSeeIt()
	}
	switch g.obj {
	case MAG:
		return g.rep("안타깝지만 잡지는 난쟁이 말로 쓰여 있어.")
	case TABLET:
		return g.rep("\"어둠의 방에 빛을 가져온 걸 축하한다!\"")
	case MESSAGE:
		return g.rep("\"여긴 해적이 보물 상자를 숨기는 미로가 아니다.\"")
	case OYSTER:
		if g.hinted[1] {
			if g.toting(OYSTER) {
				return g.rep("전에 했던 말이랑 똑같은 소리야.")
			}
		} else if g.closed && g.toting(OYSTER) {
			g.offer(1)
			return aDone()
		}
		return g.repDefault()
	default:
		return g.repDefault()
	}
}

// cantSeeIt: 보이지 않는 객체를 가리켰을 때. (advent.w cant_see_it)
func (g *Game) cantSeeIt() actResult {
	if (g.verb == FIND || g.verb == INVENTORY) && g.word2 == "" {
		return aChange(g.verb)
	}
	fmt.Fprintf(g.out, "여기 %s 같은 건 안 보여.\n", g.word1)
	return aDone()
}

// ---- SAY ----

func (g *Game) doSay() actResult {
	if g.word2 != "" {
		g.word1 = g.word2
	}
	if e, found := g.lookup(g.word1); found {
		magic := (e.typ == actionType && action(e.meaning) == FEEFIE) ||
			(e.typ == motionType && (motion(e.meaning) == XYZZY ||
				motion(e.meaning) == PLUGH || motion(e.meaning) == PLOVER))
		if magic {
			g.word2 = ""
			g.obj = NOTHING
			if e.typ == motionType {
				return g.tryMotionR(motion(e.meaning))
			}
			g.verb = action(e.meaning) // FEEFIE
			return aChange(g.verb)
		}
	}
	return g.rep(fmt.Sprintf("그래, \"%s\".", g.word1))
}
