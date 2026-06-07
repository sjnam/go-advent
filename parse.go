package main

import "fmt"

// 사용자 명령을 읽고 해석하는 부분. (advent.w "The main control loop"의
// "Get user input" 및 관련 조각들)

// getUserInput이 작은 순환에 돌려주는 결과의 종류.
type inputResult int

const (
	inputEOF          inputResult = iota // 입력 끝
	inputMotion                          // 이동: g.mot 설정됨
	inputTransitive                      // 동작+객체: g.verb, g.obj 설정됨
	inputIntransitive                    // 동작만: g.verb 설정됨
	inputSpeak                           // 고정 메시지: g.speakIdx 설정됨
)

// tryMotionInput은 이동 m을 지정하고 작은 순환에 이동 요청을 돌려준다.
func (g *Game) tryMotionInput(m motion) inputResult {
	g.mot = m
	return inputMotion
}

// streq는 두 단어가 앞 5글자까지 같은지 본다. (advent.w 매크로 streq)
func streq(a, b string) bool {
	if len(a) > 5 {
		a = a[:5]
	}
	if len(b) > 5 {
		b = b[:5]
	}
	return a == b
}

// getUserInput은 명령 하나를 읽어 해석한다. 명령이 완성되면 그 종류를
// 돌려주고, 자질구레한 경우(모르는 단어, 되묻기 등)는 스스로 출력하고
// 다시 입력을 받는다.
// keep이 true면 직전 명령 상태(verb 등)를 유지한 채 입력만 다시 받는다
// ("Take what?" 뒤 객체만 받는 경우).
func (g *Game) getUserInput(keep bool) inputResult {
	if keep {
		goto cycle
	}
restart: // 작은 순환 재시작: 명령 상태를 비운다
	g.verb, g.oldverb = ABSTAIN, ABSTAIN
	g.oldobj = g.obj
	g.obj = NOTHING
cycle: // 명령 상태를 유지한 채 다시 입력만 받는다
	g.giveHint()
	g.wasDark = g.dark()
	g.ran(0) // 추격에 변화를 주려 난수기를 한 번 돌린다
	if g.knifeLoc > int(limbo) && g.knifeLoc != int(g.loc) {
		g.knifeLoc = int(limbo)
	}
	// 폐쇄 후, 들고 있는 prop<0 물건은 한 번 내려놨다 집어야 묘사된다. (advent.w:4092)
	if g.closed {
		if g.prop[OYSTER] < 0 && g.toting(OYSTER) {
			fmt.Fprintf(g.out, "%s\n", objNote[objOffset[OYSTER]+1])
		}
		for j := object(1); j <= maxObj; j++ {
			if g.toting(j) && g.prop[j] < 0 {
				g.prop[j] = -1 - g.prop[j]
			}
		}
	}
	if !g.listen() {
		return inputEOF
	}
	g.turns++

	// 특수 입력 처리: "say"에 두 단어를 주면 아무 말도 안 한다.
	if g.verb == SAY {
		if g.word2 != "" {
			g.verb = ABSTAIN
		} else {
			return inputTransitive // say <word1>
		}
	}
	if g.checkClocks() { // 동굴이 닫히면 제자리로(최종 퍼즐 장소에서 다시 묘사)
		return g.tryMotionInput(NOWHERE)
	}
	if g.quitting { // 램프가 다 닳아 포기
		return inputEOF
	}

	// "enter water/stream" 같은 경우 처리. ("enter"는 어휘상 이동이다)
	if streq(g.word1, "enter") {
		if streq(g.word2, "water") || streq(g.word2, "strea") {
			if g.waterHere() {
				g.report("Your feet are now wet.")
				goto restart
			}
			g.report(g.defaultMsg[GO]) // default_to(GO)
			goto restart
		} else if g.word2 != "" {
			goto shift
		}
	}

parse:
	// 친절한 안내: WEST를 자주 치면 W로 줄여 쓰라고 알려준다.
	if streq(g.word1, "west") {
		g.westCount++
		if g.westCount == 10 {
			fmt.Fprintf(g.out, " If you prefer, simply type W rather than WEST.\n")
		}
	}

	{
		e, found := g.lookup(g.word1)
		if !found {
			fmt.Fprintf(g.out, "Sorry, I don't know the word \"%s\".\n", g.word1)
			goto cycle
		}
		g.commandType = e.typ
		switch e.typ {
		case motionType:
			g.mot = motion(e.meaning)
			return inputMotion

		case objectType:
			g.obj = object(e.meaning)
			switch g.makeObjMeaningful() {
			case objMotion:
				return inputMotion
			case objReport:
				goto restart
			case objCantSee:
				// cant_see_it
				if (g.verb == FIND || g.verb == INVENTORY) && g.word2 == "" {
					return inputTransitive
				}
				fmt.Fprintf(g.out, "I see no %s here.\n", g.word1)
				goto restart
			}
			// objOK
			if g.word2 != "" {
				goto shift
			}
			if g.verb != ABSTAIN {
				return inputTransitive
			}
			fmt.Fprintf(g.out, "What do you want to do with the %s?\n", g.word1)
			goto cycle

		case actionType:
			g.verb = action(e.meaning)
			if g.verb == SAY {
				if g.word2 != "" {
					g.obj = object(g.word2[0]) // SAY는 말할 거리를 obj 첫 글자로
				} else {
					g.obj = NOTHING
				}
			} else if g.word2 != "" {
				goto shift
			}
			if g.obj != NOTHING {
				return inputTransitive
			}
			return inputIntransitive

		case messageType:
			g.speakIdx = e.meaning
			return inputSpeak
		}
	}

shift: // 둘째 단어를 첫째 자리로 옮겨 다시 해석한다
	g.word1 = g.word2
	g.word2 = ""
	goto parse
}

// makeObjMeaningful이 돌려주는 신호.
type objResult int

const (
	objOK      objResult = iota // obj가 여기서 말이 됨
	objCantSee                  // 여기 그런 건 안 보임
	objMotion                   // 사실은 이동이었음 (g.mot 설정됨)
	objReport                   // 메시지를 이미 출력함 (작은 순환 재시작)
)

// makeObjMeaningful은 지정한 객체가 현재 장소에서 말이 되는지 확인한다.
// 물/기름처럼 병 안이나 장소 특성으로 존재하는 특수 경우도 처리한다.
// (advent.w "Make sure obj is meaningful at the current location")
func (g *Game) makeObjMeaningful() objResult {
	if g.toting(g.obj) || g.isAtLoc(g.obj) {
		return objOK
	}
	objectInBottle := (g.obj == WATER && g.prop[BOTTLE] == 0) ||
		(g.obj == OIL && g.prop[BOTTLE] == 2)
	switch g.obj {
	case GRATE:
		if r := g.grateAsMotion(); r != objOK {
			return r
		}
		return objCantSee
	case DWARF:
		if g.dflag >= 2 && g.dwarf() {
			return objOK
		}
		return objCantSee
	case PLANT:
		if g.isAtLoc(PLANT2) && g.prop[PLANT2] != 0 {
			g.obj = PLANT2
			return objOK
		}
		return objCantSee
	case KNIFE:
		if int(g.loc) != g.knifeLoc {
			return objCantSee
		}
		g.knifeLoc = -1
		g.report("The dwarves' knives vanish as they strike the walls of the cave.")
		return objReport
	case ROD:
		if !g.here(ROD2) {
			return objCantSee
		}
		g.obj = ROD2
		return objOK
	case WATER, OIL:
		if g.here(BOTTLE) && objectInBottle {
			return objOK
		}
		if (g.obj == WATER && g.waterHere()) || (g.obj == OIL && g.oilHere()) {
			return objOK
		}
		return objCantSee
	}
	return objCantSee
}

// grateAsMotion은 GRATE가 사실 이동 단어로 쓰였는지 본다(표면에서).
// (advent.w "If GRATE is actually a motion word, move to it")
func (g *Game) grateAsMotion() objResult {
	if g.loc < minLowerLoc {
		switch g.loc {
		case road, valley, slit:
			g.mot = DEPRESSION
			return objMotion
		case cobbles, debris, awk, bird, spit:
			g.mot = ENTRANCE
			return objMotion
		}
	}
	return objOK
}
