package main

import (
	"fmt"
	"strings"
)

// 동작(verb) 디스패치. (advent.w "Simple verbs", "Liquid assets", "The other actions")

// performAction이 작은 순환에 돌려주는 결과.
type paResult int

const (
	paDone       paResult = iota // 처리 끝 (다음 명령으로)
	paNeedObject                 // 객체가 더 필요함 ("Take what?" 뒤 verb 유지 재입력)
	paTryMove                    // 이동으로 전환 (g.mot 설정됨)
	paCommence                   // 제자리에서 다시 묘사 (불 켜기 등)
)

// actResult는 동사 하나를 처리한 결과다. 원본의 continue/change_to/
// try_motion/goto transitive 같은 제어 흐름을 값으로 표현한다.
type actResult struct {
	kind         paResult // 기본 결과
	changeTo     action   // change_to: 이 동사로 다시 디스패치 (ABSTAIN이면 없음)
	toTransitive bool     // 자동사 갈래에서 타동사 갈래로 넘어감
}

func aDone() actResult           { return actResult{} }
func aNeed() actResult           { return actResult{kind: paNeedObject} }
func aTrans() actResult          { return actResult{toTransitive: true} }
func aChange(v action) actResult { return actResult{changeTo: v} }

// 자주 쓰는 결과 헬퍼들 (원본 매크로 report/default_to/try_motion/stay_put).
func (g *Game) rep(m string) actResult        { g.report(m); return aDone() }
func (g *Game) defTo(v action) actResult      { g.report(g.defaultMsg[v]); return aDone() }
func (g *Game) tryMotionR(m motion) actResult { g.mot = m; return actResult{kind: paTryMove} }
func (g *Game) stayPut() actResult            { return g.tryMotionR(NOWHERE) }

// repDefault는 현재 동사의 기본 메시지를 낸다(있으면). (advent.w report_default)
func (g *Game) repDefault() actResult {
	if g.defaultMsg[g.verb] != "" {
		g.report(g.defaultMsg[g.verb])
	}
	return aDone()
}

// getObjR은 "<Verb> what?"을 묻고 객체 재입력을 요청한다. (advent.w get_object)
func (g *Game) getObjR() actResult {
	w := g.word1
	if w != "" {
		w = strings.ToUpper(w[:1]) + w[1:]
	}
	fmt.Fprintf(g.out, "%s 뭘?\n", w)
	return aNeed()
}

// performAction은 verb(와 obj)에 따라 동작을 수행한다.
func (g *Game) performAction(transitive bool) paResult {
	for {
		var r actResult
		if transitive {
			r = g.doTransitive()
		} else {
			r = g.doIntransitive()
			if r.toTransitive {
				transitive = true
				continue
			}
		}
		if r.changeTo != ABSTAIN {
			g.oldverb = g.verb
			g.verb = r.changeTo
			transitive = true
			continue
		}
		return r.kind
	}
}

// doIntransitive는 객체 없이 들어온 동작을 처리한다. (advent.w intransitive)
func (g *Game) doIntransitive() actResult {
	switch g.verb {
	case GO, RELAX:
		return g.repDefault()
	case ON, OFF, POUR, FILL, DRINK, BLAST, KILL:
		return aTrans()
	case TAKE:
		if g.first[g.loc] == NOTHING || g.link[g.first[g.loc]] != NOTHING || g.dwarf() {
			return g.getObjR()
		}
		g.obj = g.first[g.loc]
		return aTrans()
	case EAT:
		if !g.here(FOOD) {
			return g.getObjR()
		}
		g.obj = FOOD
		return aTrans()
	case OPEN, CLOSE:
		if g.place[GRATE] == g.loc || g.place[GRATE_] == g.loc {
			g.obj = GRATE
		} else if g.place[DOOR] == g.loc {
			g.obj = DOOR
		} else if g.here(CLAM) {
			g.obj = CLAM
		} else if g.here(OYSTER) {
			g.obj = OYSTER
		}
		if g.here(CHAIN) {
			if g.obj != NOTHING {
				return g.getObjR()
			}
			g.obj = CHAIN
		}
		if g.obj != NOTHING {
			return aTrans()
		}
		return g.rep("여기엔 자물쇠 달린 게 없어!")
	case READ:
		if g.dark() {
			return g.getObjR()
		}
		if g.here(MAG) {
			g.obj = MAG
		}
		if g.here(TABLET) {
			if g.obj != NOTHING {
				return g.getObjR()
			}
			g.obj = TABLET
		}
		if g.here(MESSAGE) {
			if g.obj != NOTHING {
				return g.getObjR()
			}
			g.obj = MESSAGE
		}
		if g.closed && g.toting(OYSTER) {
			g.obj = OYSTER
		}
		if g.obj != NOTHING {
			return aTrans()
		}
		return g.getObjR()
	case INVENTORY:
		return g.doInventory()
	case BRIEF:
		g.interval = 10000
		g.lookCount = 3
		return g.rep("알았어, 이제부터 어떤 장소든 처음 왔을 때만 전부 설명할게.\n" +
			"전체 설명을 보려면 \"봐\"라고 해.")
	case SCORE:
		fmt.Fprintf(g.out, "지금 그만두면 너는 %d점을 얻어.\n"+
			"만점은 %d점이야.\n", g.score()-4, maxScore)
		if !g.yes("정말로 지금 그만둘래?", g.ok(), g.ok()) {
			return aDone()
		}
		g.gaveUp = true
		g.quitting = true
		return aDone()
	case QUIT:
		if !g.yes("정말 지금 그만두고 싶어?", g.ok(), g.ok()) {
			return aDone()
		}
		g.gaveUp = true
		g.quitting = true
		return aDone()
	case FEEFIE:
		k := 0
		for !streq(g.word1, incantation[k]) {
			k++
		}
		if g.foobar == -k {
			return g.proceedFoobar(k)
		}
		if g.foobar == 0 {
			return g.defTo(WAVE) // nada_sucede
		}
		return g.rep("왜 그래, 글도 못 읽어?  이제 처음부터 다시 하는 게 좋겠어.")
	default:
		return g.getObjR()
	}
}

// doTransitive는 객체와 함께 들어온 동작을 처리한다. (advent.w transitive)
func (g *Game) doTransitive() actResult {
	switch g.verb {
	case SAY:
		return g.doSay()
	case TAKE:
		return g.doTake()
	case DROP:
		return g.doDrop()
	case TOSS:
		return g.doToss()
	case KILL:
		return g.doKill()
	case FEED:
		return g.doFeed()
	case OPEN, CLOSE:
		return g.doOpenClose()
	case READ:
		return g.doRead()
	case EAT:
		return g.doEat()
	case WAVE:
		return g.doWave()
	case BLAST:
		return g.doBlast()
	case RUB:
		if g.obj == LAMP {
			return g.repDefault()
		}
		return g.rep("희한하네.  별다른 일은 안 일어나.")
	case FIND, INVENTORY:
		return g.doFind()
	case BREAK:
		return g.doBreak()
	case WAKE:
		if g.closed && g.obj == DWARF {
			fmt.Fprintf(g.out, "넌 가장 가까운 난쟁이를 쿡 찔러.  녀석은 투덜대며 깨어나, 널 한 번\n"+
				"노려보고는, 욕을 내뱉으며 도끼를 집어 들어.\n")
			return g.dwarvesUpset()
		}
		return g.repDefault()
	case ON:
		return g.doLampOn()
	case OFF:
		return g.doLampOff()
	case DRINK:
		return g.doDrink()
	case POUR:
		return g.doPour()
	case FILL:
		return g.doFill()
	default:
		return g.repDefault()
	}
}

// doInventory는 들고 있는 물건 목록을 보여준다. (advent.w INVENTORY)
func (g *Game) doInventory() actResult {
	found := false
	for t := object(1); t <= maxObj; t++ {
		if g.toting(t) && (objBase[t] == NOTHING || objBase[t] == t) && t != BEAR {
			if !found {
				found = true
				fmt.Fprintf(g.out, "지금 들고 있는 건 이래:\n")
			}
			fmt.Fprintf(g.out, " %s\n", objName[t])
		}
	}
	if g.toting(BEAR) {
		return g.rep("아주 크고 온순한 곰이 널 따라오고 있어.")
	}
	if !found {
		return g.rep("넌 아무것도 들고 있지 않아.")
	}
	return aDone()
}

// proceedFoobar는 주문(fee-fie-foe-foo)을 한 단계 진행한다. (advent.w "Proceed foobarically")
func (g *Game) proceedFoobar(k int) actResult {
	g.foobar = k + 1
	if g.foobar != 4 {
		return g.defTo(RELAX)
	}
	g.foobar = 0
	if g.place[EGGS] == giant || (g.toting(EGGS) && g.loc == giant) {
		return g.defTo(WAVE) // nada_sucede
	}
	if g.place[EGGS] == limbo && g.place[TROLL] == limbo && g.prop[TROLL] == 0 {
		g.prop[TROLL] = 1
	}
	kk := 2
	if g.loc == giant {
		kk = 0
	} else if g.here(EGGS) {
		kk = 1
	}
	g.move(EGGS, giant)
	return g.rep(objNote[objOffset[EGGS]+kk])
}

// doEat은 음식이나 생물을 먹으려는 시도를 처리한다. (advent.w EAT)
func (g *Game) doEat() actResult {
	switch g.obj {
	case FOOD:
		g.destroy(FOOD)
		return g.rep("고마워, 맛있었어!")
	case BIRD, SNAKE, CLAM, OYSTER, DWARF, DRAGON, TROLL, BEAR:
		return g.rep("갑자기 입맛이 뚝 떨어졌어.")
	default:
		return g.repDefault()
	}
}

// doWave는 무언가를 흔드는 동작이다. 갈라진 틈에서 막대를 흔들면 수정 다리가 생긴다.
// (advent.w WAVE)
func (g *Game) doWave() actResult {
	if g.obj != ROD || (g.loc != efiss && g.loc != wfiss) || !g.toting(g.obj) || g.closing() {
		if g.toting(g.obj) || (g.obj == ROD && g.toting(ROD2)) {
			return g.repDefault()
		}
		return g.defTo(DROP)
	}
	g.prop[CRYSTAL] = 1 - g.prop[CRYSTAL]
	return g.rep(objNote[objOffset[CRYSTAL]+2-g.prop[CRYSTAL]])
}

// doBlast는 폭파 시도다. 폐쇄 후 다이너마이트가 있어야 의미가 있다. (advent.w BLAST)
func (g *Game) doBlast() actResult {
	if g.closed && g.prop[ROD2] >= 0 {
		bonus := 45
		if g.here(ROD2) {
			bonus = 25
		} else if g.loc == neend {
			bonus = 30
		}
		fmt.Fprintf(g.out, "%s\n", g.message[bonus/5])
		g.quitting = true
		return aDone()
	}
	return g.repDefault()
}

// doFind는 보이지 않는 물건을 찾으려 할 때 안내한다. (advent.w FIND/INVENTORY)
func (g *Game) doFind() actResult {
	if g.toting(g.obj) {
		return g.defTo(TAKE)
	}
	if g.closed {
		return g.rep("네가 뭘 찾든 이 근처 어딘가에 있을 거야.")
	}
	objectInBottle := (g.obj == WATER && g.prop[BOTTLE] == 0) || (g.obj == OIL && g.prop[BOTTLE] == 2)
	if g.isAtLoc(g.obj) || (objectInBottle && g.place[BOTTLE] == g.loc) ||
		(g.obj == WATER && g.waterHere()) || (g.obj == OIL && g.oilHere()) ||
		(g.obj == DWARF && g.dwarf()) {
		return g.rep("네가 찾는 건 바로 여기 너랑 같이 있는 것 같은데.")
	}
	return g.repDefault()
}

// doBreak는 무언가를 부수는 시도다. 꽃병이나(폐쇄 후) 거울만 효과가 있다. (advent.w BREAK)
func (g *Game) doBreak() actResult {
	if g.obj == VASE && g.prop[VASE] == 0 {
		if g.toting(VASE) {
			g.drop(VASE, g.loc)
		}
		fmt.Fprintf(g.out, "넌 꽃병을 들어 바닥에 섬세하게 내동댕이쳤어.\n")
		return g.smashVase()
	}
	if g.obj != MIRROR {
		return g.repDefault()
	}
	if g.closed {
		fmt.Fprintf(g.out, "넌 거울을 쾅 후려쳐, 그러자 거울이 산산이\n"+
			"수많은 작은 조각으로 부서져.")
		return g.dwarvesUpset()
	}
	return g.rep("너무 높이 있어서 닿을 수가 없어.")
}

// smashVase는 꽃병을 깨뜨려 못 쓰게 만든다. (advent.w smash)
func (g *Game) smashVase() actResult {
	g.prop[VASE] = 2
	objBase[VASE] = VASE // 더는 움직일 수 없다
	return aDone()
}

// doLampOn/doLampOff: 램프 켜기/끄기. (advent.w ON/OFF)
func (g *Game) doLampOn() actResult {
	if !g.here(LAMP) {
		return g.repDefault()
	}
	if g.limit < 0 {
		return g.rep("램프 전력이 다 떨어졌어.")
	}
	g.prop[LAMP] = 1
	fmt.Fprintf(g.out, "이제 램프가 켜졌어.\n")
	if g.wasDark {
		return actResult{kind: paCommence}
	}
	return aDone()
}

func (g *Game) doLampOff() actResult {
	if !g.here(LAMP) {
		return g.repDefault()
	}
	g.prop[LAMP] = 0
	fmt.Fprintf(g.out, "이제 램프가 꺼졌어.\n")
	if g.dark() {
		fmt.Fprintf(g.out, "%s\n", pitchDarkMsg)
	}
	return aDone()
}
