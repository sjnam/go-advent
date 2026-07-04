@* 그 밖의 동작들.
싸우고, 먹이고, 열고 닫고, 읽고, 말하는 동사들이다. 동굴의 주민들과
얽히는 대목이라, 게임의 성격이 가장 잘 드러난다 --- 이 게임의 폭력은
대체로 소용없거나(뱀·트롤), 어리석거나(곰), 뜻밖의 방식으로만
통한다(용).

먼저 \.{죽여}. 대상 없이 외치면 공격할 만한 것을 세어 보고, 하나로
정해지면 그놈이고, 여럿이면 되묻는다.

@<동작@>=
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
		@<새를 처형한다@>@;
	case DRAGON:
		if g.prop[DRAGON] == 0 {
			@<용과 재미있는 일을 벌인다@>@;
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

@ 공격 대상 세기. 난쟁이·뱀·용·트롤·곰이 우선이고, 그들이 없을 때만
새와 조개가 후보에 든다 --- 새는 방금 던지기에서 넘어온 경우라면
빼 준다(도끼를 던졌지, 새를 노린 게 아니니까). 둘 이상이면 |redirect|로
``뭘?''을 되묻는다.

@<동작@>=
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

@ 새를 죽이겠다면 말리지는 않지만, 뱀이 아직 산속 왕의 홀에 버티고
있다면 방금 보물 하나를 영영 잃은 것이다 --- 뱀을 쫓을 수단이 새뿐이고,
그 너머의 보물이 잠기니까. 게임은 그런 말은 해 주지 않는다.

@<새를 처형한다@>=
if g.closed {
	return g.rep("아, 그 불쌍한 새는 좀 내버려 둬.")
}
g.destroy(BIRD)
g.prop[BIRD] = 0
if g.place[SNAKE] == hmk {
	g.lostTreasures++
}
return g.rep("작은 새가 죽었어.  사체가 사라져.")

@ 원작이 ``Fun stuff for dragon''이라 이름 붙인 절. 용을 공격하겠다면
게임은 능청스럽게 묻는다 --- ``뭘로? 맨손으로?'' 그런데 여기에
\.{예}라고 우기면, 정말로 맨손으로 용을 때려잡는다. 허세에 허세로
맞서는 자에게 게임이 바치는 경의다. 용이 죽으면 그 방(비밀 협곡)의
분신들이 하나로 합쳐지고, 깔고 있던 페르시아 양탄자가 비로소 보물이
된다.

@<용과 재미있는 일을 벌인다@>=
fmt.Fprintf(g.out, "뭘로?  맨손으로?\n")
g.verb = ABSTAIN
g.obj = NOTHING
if !g.listen() {
	g.quitting = true
	return aDone()
}
if !(streq(g.word1, "예") || streq(g.word1, "응") || streq(g.word1, "y")) {
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

@ 먹이기. 새는 배가 안 고프고(피오르를 그리워할 뿐이다 --- 몬티
파이선의 그 앵무새가 맞다), 트롤의 악덕은 탐식이 아니라 탐욕이며,
뱀에게 새를 먹이는 것은 보물 하나를 잃는 참사다. 곰에게 음식을 주면
온순해지고 도끼도 되찾는다. 난쟁이에게 먹이를 권하는 것은 --- 해 보면
안다.

@<동작@>=
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

@ 열고 닫기. 자물쇠가 있는 것은 격자문과 사슬뿐이고, 그나마 열쇠가
있어야 한다. 나머지는 사정 설명이다.

@<동작@>=
func (g *Game) doOpenClose() actResult {
	switch g.obj {
	case OYSTER, CLAM:
		@<조개를 열어 본다@>@;
	case GRATE, CHAIN:
		if !g.here(KEYS) {
			return g.rep("열쇠가 없잖아!")
		}
		@<창살이나 사슬을 열고 닫는다@>@;
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

@ 창살. 폐쇄 중이면 녹음된 목소리가 끼어들고, 아니면 전이표 하나로
네 가지 경우를 가린다: 이전 상태 |k|와 새 상태의 조합 |k+2*prop|이
0이면 이미 잠김, 1이면 방금 잠김, 2면 방금 열림, 3이면 이미 열림.

@<창살이나 사슬을 열고 닫는다@>=
if g.obj == CHAIN {
	@<사슬을 열고 닫는다@>@;
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

@ 사슬 잠그기는 텅 빈 방에서만 된다 --- 사슬 고리가 거기 박혀 있으니까.
잠그면 사슬은 도로 붙박이가 된다.

@<사슬을 열고 닫는다@>=
if g.verb == OPEN {
	@<사슬을 풀어 본다@>@;
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

@ 사슬을 풀려면 곰 문제부터 풀어야 한다. 사나운 곰을 지나쳐 자물쇠에
손을 뻗을 방법은 없다 --- ``뭐 어쩌면 다행이지만.'' 곰이 온순하면
사슬이 풀리며 보물이 되고, 곰도 (죽지 않았다면) 자유의 몸이 되어
데리고 다닐 수 있다.

@<사슬을 풀어 본다@>=
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

@ 조개(대합)를 열려면 삼지창을 들고 있어야 하고, 조개 자체는 내려놓고
열어야 한다. 열면 진주가 굴러 나와 아래 방의 부드러운 바닥까지
데구루루 --- 조개는 굴로 바뀐다(진주를 내놓았으니 이제 굴이다, 라는
것이 게임의 조개학이다). 이미 굴이면 열어도 헛수고다. 참고로 게임의
자칭 조개 감정 실력은 본문에서 스스로 실토한다.

@<조개를 열어 본다@>=
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

@ 읽기. 어두우면 글자가 안 보인다. 잡지는 난쟁이 말로 쓰여 있고(그래도
Witt's End에 배달되는 정기 간행물이다), 서판과 쪽지는 짧은 안내문이며,
폐쇄 후의 굴 껍데기에는 마지막 힌트가 적혀 있다 --- 값은 10점이다.

@<동작@>=
func (g *Game) doRead() actResult {
	if g.dark() {
		@<안 보인다고 답한다@>@;
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

@ @<안 보인다고 답한다@>=
if (g.verb == FIND || g.verb == INVENTORY) && g.word2 == "" {
	return aChange(g.verb)
}
fmt.Fprintf(g.out, "여기 %s 같은 건 안 보여.\n", g.word1)
return aDone()

@ 말하기. 대부분은 메아리다 --- ``그래, \.{안녕}.'' 하지만 말한 단어가
마법 주문이면 그 주문을 실제로 왼 것으로 친다. \.{말해 xyzzy}가
\.{xyzzy}와 같아야 마땅하니까.

@<동작@>=
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
