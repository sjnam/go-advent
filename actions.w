@* 동작.
동사를 실제 행동으로 바꾸는 층이다. 원작에서 각 동사는 처리가 끝나면
매크로로 갈 곳을 알렸다 --- |report|(한마디 하고 다음 명령),
|change_to|(다른 동사로 변신해 재도전), |try_motion|(사실 이동이었음).
이 포팅은 그 매크로들을 값 |actResult|로 바꿨다. 동사 하나의 처리는
``무슨 일이 있었고 다음은 어디로''를 담은 값을 돌려주는 순수한 함수
호출이 되고, 갈래를 트는 일은 |performAction| 한 곳이 맡는다.

@<동작@>=
type paResult int

const (
	paDone       paResult = iota // 처리 끝 (다음 명령으로)
	paNeedObject                 // 객체가 더 필요함 (verb를 쥐고 재입력)
	paTryMove                    // 이동으로 전환 (mot 설정됨)
	paCommence                   // 제자리에서 다시 묘사 (불 켜기 등)
)

type actResult struct {
	kind         paResult // 기본 결과
	changeTo     action   // 이 동사로 다시 디스패치 (ABSTAIN이면 없음)
	toTransitive bool     // 자동사 갈래에서 타동사 갈래로 넘어감
}

func aDone() actResult           { return actResult{} }
func aNeed() actResult           { return actResult{kind: paNeedObject} }
func aTrans() actResult          { return actResult{toTransitive: true} }
func aChange(v action) actResult { return actResult{changeTo: v} }

@ 자주 쓰는 결과 도우미들. |rep|은 말하고 끝, |defTo|는 남의 기본
메시지를 빌려 말하고 끝, |tryMotionR|은 이동으로 전환, |stayPut|은
제자리 재묘사(NOWHERE 이동)다. |getObjR|은 ``$\langle$동사$\rangle$
뭘?''이라고 되물으며 객체 재입력을 청한다.

@<동작@>=
func (g *Game) rep(m string) actResult        { g.report(m); return aDone() }
func (g *Game) defTo(v action) actResult      { g.report(g.defaultMsg[v]); return aDone() }
func (g *Game) tryMotionR(m motion) actResult { g.mot = m; return actResult{kind: paTryMove} }
func (g *Game) stayPut() actResult            { return g.tryMotionR(NOWHERE) }

func (g *Game) repDefault() actResult {
	if g.defaultMsg[g.verb] != "" {
		g.report(g.defaultMsg[g.verb])
	}
	return aDone()
}

func (g *Game) getObjR() actResult {
	w := g.word1
	if w != "" {
		w = strings.ToUpper(w[:1]) + w[1:]
	}
	fmt.Fprintf(g.out, "%s 뭘?\n", w)
	return aNeed()
}

@ 갈래 트기. 자동사가 객체를 알아맞히면 타동사로 넘어가고
(|toTransitive|), 동사가 변신하면(|changeTo|) 그 동사의 타동사 처리로
다시 든다. \.{던져 음식}이 \.{먹여}로, \.{먹여}가 다시 \.{먹어}로
꼬리를 무는 식이다.

@<동작@>=
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

@ 객체 없이 들어온 동사. 상당수는 객체를 짐작할 수 있다 --- 방에 물건이
꼭 하나면 \.{가져가}의 대상은 그것이고, 음식이 보이면 \.{먹어}의 대상도
뻔하다. 짐작이 둘 이상으로 갈리면 되묻는다. 나머지는 각자의 사연대로.

@<동작@>=
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
		@<열거나 닫을 것을 짐작한다@>@;
	case READ:
		@<읽을 것을 짐작한다@>@;
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
	case SAVE:
		if err := g.save(saveFile); err != nil {
			return g.rep("저장하지 못했어.")
		}
		return g.rep("게임을 저장했어.")
	case RESTORE:
		if err := g.load(saveFile); err != nil {
			return g.rep("불러올 저장 파일이 없어.")
		}
		g.report("게임을 불러왔어.")
		return actResult{kind: paCommence} // 복원된 자리에서 다시 묘사
	case FEEFIE:
		@<주문의 진행을 살핀다@>@;
	default:
		return g.getObjR()
	}
}

@ 열고 닫기의 짐작. 격자문, 문, 대합, 굴 순으로 살피되, 사슬까지
보이면 후보가 둘이니 되묻는다.

@<열거나 닫을 것을 짐작한다@>=
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

@ 읽을 것의 짐작. 어두우면 짐작할 도리가 없다. 잡지·서판·쪽지가
후보이고, 폐쇄 후에 굴(oyster)을 들고 있다면 그 껍데기의 단서다.

@<읽을 것을 짐작한다@>=
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

@ 객체와 함께 들어온 동사는 각자의 처리로 흩어진다. 구현이 짧은 몇은
여기서 바로 해치운다 --- \.{문질러}는 램프면 기본 핀잔, 다른 것이면
희한하다는 반응이고, \.{깨워}는 폐쇄 후 난쟁이를 건드리는 단 한
경우에만 결과가 있다(치명적인 결과가).

@<동작@>=
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

@ 소지품 목록. 여러 조각짜리 물건은 대표만 세고, 곰은 목록이 아니라
따로 한 문장을 받는다 --- 곰은 소지품이 아니라 일행이니까.

@<동작@>=
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

@ 거인의 주문 \.{fee fie foe foo}. 잭과 콩나무의 거인이 흥얼대던 그
노래다(\.{fum}까지 알아듣지만, \.{fum}은 언제나 헛소리 취급이다).
네 마디를 순서대로 외쳐야 하고, |foobar|가 진행도를 기억한다. 음수로
저장해 두는 것은 ``방금 그 단어까지 맞게 왔음''의 표시다.

@<게임의 상태@>=
foobar int // fee-fie-foe-foo 주문 진행도

@ @<주문의 진행을 살핀다@>=
k := 0
for !streq(g.word1, incantation[k]) {
	k++
}
if g.foobar == -k {
	return g.proceedFoobar(k)
}
if g.foobar == 0 {
	return g.defTo(WAVE) // 아무 일도 일어나지 않아
}
return g.rep("왜 그래, 글도 못 읽어?  이제 처음부터 다시 하는 게 좋겠어.")

@ 주문이 완성되면 황금 알이 거인의 방으로 돌아온다 --- 트롤에게
통행료로 던져 준 알을 이 주문으로 되찾는 것이 정석 공략이다. 알이
트롤과 함께 사라진 상태라면 트롤을 다리로 복귀시키는 뒷정리도 한다.
메시지는 어디서 외쳤는지에 따라 다르다: 거인의 방이면 알이 눈앞에
나타나고, 알이 있는 방이면 사라지는 광경을, 다른 곳이면 아무것도
못 본다.

@<동작@>=
var incantation = [...]string{"fee", "fie", "foe", "foo", "fum"}

func (g *Game) proceedFoobar(k int) actResult {
	g.foobar = k + 1
	if g.foobar != 4 {
		return g.defTo(RELAX)
	}
	g.foobar = 0
	if g.place[EGGS] == giant || (g.toting(EGGS) && g.loc == giant) {
		return g.defTo(WAVE)
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

@ 먹기. 음식은 먹으면 없어지고, 살아 있는 것들에게는 입맛이 떨어진다.

@<동작@>=
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

@ 흔들기. 딱 한 곳에서 마법이 일어난다: 갈라진 틈 가에서 검은 막대를
흔들면 수정 다리가 나타나고, 또 흔들면 사라진다. |prop[CRYSTAL]|을
$1$과 $0$ 사이에서 뒤집는 |1-prop|이 그 스위치다. 폐쇄 중에는 마법이
듣지 않는다.

@<동작@>=
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

@ 폭파. 게임의 마지막 수수께끼라 자세한 이야기는 아껴 두지만, 요점은
폐쇄 후 보관실에서 별 달린 검은 막대(사실은 다이너마이트)와 얼마나
멀리 떨어져 있느냐다. 막대를 곁에 두고 터뜨리면 함께 산화하고(25점),
북동 끝이면 그저 그런 결말(30점), 남서 끝까지 물러났으면 그랜드마스터의
결말(45점)이다. 보너스를 5로 나누면 마침 결말 메시지의 인덱스가 되는
것은 우연이 아니라 원작의 설계다.

@<동작@>=
func (g *Game) doBlast() actResult {
	if g.closed && g.prop[ROD2] >= 0 {
		g.bonus = 45
		if g.here(ROD2) {
			g.bonus = 25
		} else if g.loc == neend {
			g.bonus = 30
		}
		fmt.Fprintf(g.out, "%s\n", g.message[g.bonus/5])
		g.quitting = true
		return aDone()
	}
	return g.repDefault()
}

@ 찾기. 나는 네 눈과 손이지 레이더가 아니라서, 눈앞에 있는 것만
확인해 준다. 폐쇄 후에는 능청스러운 답 하나로 통일 --- 어차피 모든
게 이 근처 어딘가에 있으니까.

@<동작@>=
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

@ 부수기. 명나라 꽃병은 ``섬세하게 내동댕이''칠 수 있고, 거울은 폐쇄
후에만 손이 닿는데 깨는 순간 난쟁이들이 깨어난다. |smashVase|는 꽃병을
못 쓰는 붙박이로 강등시킨다 --- 깨진 꽃병 조각은 들고 다닐 수 없다.

@<동작@>=
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

func (g *Game) smashVase() actResult {
	g.prop[VASE] = 2
	objBase[VASE] = VASE // 더는 움직일 수 없다
	return aDone()
}

@ 램프 켜고 끄기. 켜는 순간 방금까지 어두웠다면 제자리 재묘사를
청한다 --- 불을 켰으니 방이 보여야지. 끄는 순간 어두워지면 칠흑
경고를 들려준다.

@<동작@>=
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
