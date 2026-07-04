@* 명령 해석.
네가 친 한 줄이 게임의 의도로 바뀌는 곳이다. 결과는 다섯 가지 중
하나다: 입력이 끝났거나, 이동이거나, 객체가 딸린 동작(타동사)이거나,
동작뿐(자동사)이거나, 정해진 대꾸(메시지)거나. 애매한 경우 --- 모르는
단어, ``열쇠로 뭘 하고 싶어?'' 같은 되물음 --- 는 이 안에서 스스로
해결하고 다시 입력을 받으므로, 밖에는 완성된 명령만 나간다.

명령을 담는 상태들. |oldverb|와 |oldobj|는 직전 명령의 잔상인데,
힌트 판정과 동사 바꿔치기에 쓰인다.

@<게임의 상태@>=
mot      motion // 지금 지정된 이동
verb     action // 지금 지정된 동작
oldverb  action // 바뀌기 전 verb
obj      object // 지금 지정된 객체
oldobj   object // 이전 obj
turns    int    // 명령을 읽은 횟수
speakIdx int    // 메시지 단어가 가리키는 인덱스

@ @<명령 해석@>=
type inputResult int

const (
	inputEOF          inputResult = iota // 입력 끝
	inputMotion                          // 이동: mot 설정됨
	inputTransitive                      // 동작+객체: verb, obj 설정됨
	inputIntransitive                    // 동작만: verb 설정됨
	inputSpeak                           // 고정 메시지: speakIdx 설정됨
)

func (g *Game) tryMotionInput(m motion) inputResult {
	g.mot = m
	return inputMotion
}

func streq(a, b string) bool { return a == b }

@ |getUserInput|은 원작에서 가장 |goto|가 무성하던 대목이라, 여기서도
label 넷이 남았다. |restart|는 명령 상태를 비우고 새로 시작,
|cycle|은 상태를 쥔 채 입력만 다시, |parse|는 |word1|의 해석,
|shift|는 둘째 단어를 첫째 자리로 밀어 재해석이다. 두 단어의 순서가
자유로운 비결이 이 |shift|다 --- \.{열쇠 가져가}는 객체를 먼저 만나
동사를 기다리고, \.{가져가 열쇠}는 동사를 먼저 만나 객체를 기다리는데,
어느 쪽이든 남은 단어를 밀어 넣으면 같은 곳에 닿는다.

인자 |keep|이 참이면 ``가져가 뭘?'' 뒤처럼 동사를 유지한 채 객체만
받는다.

@<명령 해석@>=
func (g *Game) getUserInput(keep bool) inputResult {
	if keep {
		goto cycle
	}
restart:
	g.verb, g.oldverb = ABSTAIN, ABSTAIN
	g.oldobj = g.obj
	g.obj = NOTHING
cycle:
	@<턴을 시작하는 집안일을 하고, 한 줄을 듣는다@>@;
	@<특수한 입력과 시계를 처리한다@>@;
parse:
	{
		e, found := g.lookup(g.word1)
		if !found {
			fmt.Fprintf(g.out, "\"%s\" 같은 단어는 모르겠어.\n", g.word1)
			goto cycle
		}
		switch e.typ {
		case motionType:
			g.mot = motion(e.meaning)
			return inputMotion
		case objectType:
			@<객체 단어를 처리한다@>@;
		case actionType:
			@<동작 단어를 처리한다@>@;
		case messageType:
			g.speakIdx = e.meaning
			return inputSpeak
		}
	}
shift:
	g.word1 = g.word2
	g.word2 = ""
	goto parse
}

@ 매 턴의 집안일: 힌트를 저울질하고, 어둠을 기억해 두고(다음 턴의
구덩이 판정용), 난수기를 공회전 한 번 시켜 난쟁이 추격에 변화를 주고,
다른 방에 두고 온 칼의 언급 권리를 지운다. 폐쇄 후에는 들고 있는
물건들의 음수 |prop|을 양수로 뒤집는다 --- 최종 보관실의 물건들은 한 번
내려놨다 집어야 묘사되게 한 원작의 잔꾀다. 그 뒤에 비로소 한 줄을 듣는다.

@<턴을 시작하는 집안일을 하고, 한 줄을 듣는다@>=
g.giveHint()
g.wasDark = g.dark()
g.ran(0)
if g.knifeLoc > int(limbo) && g.knifeLoc != int(g.loc) {
	g.knifeLoc = int(limbo)
}
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

@ 몇 가지 특수 사정. \.{말해}에 두 단어를 주면 동사를 지워 둘째 단어가
그대로 해석되게 하고(``말해 북''은 그냥 북쪽으로 간다), 시계를 돌려
폐쇄가 닥쳤으면 제자리 이동으로 빠져나간다. \.{들어가 물}은 이동표에
없는 관용구라 여기서 발을 적셔 준다.

@<특수한 입력과 시계를 처리한다@>=
if g.verb == SAY {
	if g.word2 != "" {
		g.verb = ABSTAIN
	} else {
		return inputTransitive // 말해 <word1>
	}
}
if g.checkClocks() { // 동굴이 닫히면 제자리로
	return g.tryMotionInput(NOWHERE)
}
if g.quitting { // 램프가 다 닳아 포기
	return inputEOF
}
if streq(g.word1, "들어가") {
	if streq(g.word2, "물") || streq(g.word2, "개울") {
		if g.waterHere() {
			g.report("이제 발이 젖었어.")
			goto restart
		}
		g.report(g.defaultMsg[GO])
		goto restart
	} else if g.word2 != "" {
		goto shift
	}
}

@ 객체 단어. 먼저 그 객체가 이 자리에서 말이 되는지
|makeObjMeaningful|에게 묻는다. 말이 되면, 둘째 단어가 있으면 밀어
해석하고, 동사가 이미 있으면 타동사 명령이 완성이고, 동사가 없으면
``그걸로 뭘 하고 싶어?''라고 되묻는다.

@<객체 단어를 처리한다@>=
g.obj = object(e.meaning)
switch g.makeObjMeaningful() {
case objMotion:
	return inputMotion
case objReport:
	goto restart
case objCantSee:
	if (g.verb == FIND || g.verb == INVENTORY) && g.word2 == "" {
		return inputTransitive
	}
	fmt.Fprintf(g.out, "여기 %s 같은 건 안 보여.\n", g.word1)
	goto restart
}
if g.word2 != "" {
	goto shift
}
if g.verb != ABSTAIN {
	return inputTransitive
}
fmt.Fprintf(g.out, "%s(으)로 뭘 하고 싶어?\n", g.word1)
goto cycle

@ 동작 단어. \.{말해}만 특별하다: 말할 거리는 어차피 |doSay|가
|word1|/|word2|에서 다시 꺼내므로, |obj|에는 첫 바이트를 표지 삼아
욱여넣는다 --- 객체 아닌 것을 객체 자리에 앉히는 원작의 능청을 그대로
옮긴 것이다.

@<동작 단어를 처리한다@>=
g.verb = action(e.meaning)
if g.verb == SAY {
	if g.word2 != "" {
		g.obj = object(g.word2[0])
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

@ 객체가 이 자리에서 말이 되는가. 눈앞에 있거나 손에 있으면 당연히
되고, 나머지는 사정이 딸린 특수 경우들이다: 창살은 멀리서 부르면
사실 이동 명령이고(다음 절), 난쟁이는 성나 있어야 하며, 식물은 다 자란
쪽 조각으로 바꿔치고, 칼은 이 방에서 언급된 직후에만 실체가 있고
(만지려 들면 사라진다 --- 난쟁이의 칼은 주울 수 없어야 하니까), 검은
막대는 닮은꼴 |ROD2|로 통하고, 물과 기름은 병 속에 있거나 장소에
흐르고 있어야 한다.

@<명령 해석@>=
type objResult int

const (
	objOK      objResult = iota // obj가 여기서 말이 됨
	objCantSee                  // 여기 그런 건 안 보임
	objMotion                   // 사실은 이동이었음 (mot 설정됨)
	objReport                   // 메시지를 이미 출력함 (재시작)
)

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
		g.report("난쟁이들의 칼이 동굴 벽에 부딪히며 사라져.")
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

@ 지상에서 \.{창살}이라고 외치는 모험가는 십중팔구 창살로 {\it 가고}
싶은 것이다. 어디 있느냐에 따라 웅덩이행 또는 입구행 이동으로 바꿔 준다.

@<명령 해석@>=
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
