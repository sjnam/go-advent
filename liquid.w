@* 액체 자산.
독자들은 이미 눈치챘겠지만, 이 게임에서 물과 기름은 이상한 물건이다.
병에 담기면 물건이고, 바닥에 흐르면 장소의 속성이며, 부으면 사라진다.
물건 노릇을 하는 것은 병 |BOTTLE| 하나이고, |WATER|와 |OIL|은 병의
상태(|prop[BOTTLE]|이 0이면 물, 2면 기름, 1이면 빈 병)를 따라다니는
그림자 객체다. 원작의 장 제목 ``Liquid assets''를 그대로 빌려 왔다 ---
이만한 말장난은 번역해도 죽지 않는다.

@<동작@>=
func (g *Game) bottleEmpty() bool { return g.prop[BOTTLE] == 1 || g.prop[BOTTLE] < 0 }

func (g *Game) objectInBottle() bool {
	return (g.obj == WATER && g.prop[BOTTLE] == 0) || (g.obj == OIL && g.prop[BOTTLE] == 2)
}

@ 마시기. 대상 없이 \.{마셔}라고 하면 주변의 물이나 병의 물을 찾고,
물 아닌 것을 마시겠다면 \.{먹어}의 핀잔으로 돌린다. 개울물을 그냥
마시는 경우는 |repDefault|가 처리한다---\.{마셔}의 기본 메시지가
바로 그 ``광물 맛이 강하게 나지만'' 하는 감상평이다.

@<동작@>=
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

@ 붓기. 무엇을 붓는지 정하고(병만 가리켰으면 내용물로 바꾼다), 병을
비운 뒤, 어디에 부었는지가 관건이다. 식물에 부으면 콩나무가 자라고,
녹슨 문에 부으면---물이냐 기름이냐가 갈림길이다.

@<동작@>=
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
		@<식물에 물을 준다@>@;
	}
	if g.loc == g.place[DOOR] {
		@<문의 경첩에 붓는다@>@;
	}
	return g.rep("병이 비고 바닥이 젖었어.")
}

@ 식물은 물을 줄 때마다 한 단계씩 자란다: 새싹, 12피트 콩나무, 구멍까지
닿는 거목---그리고 세 번째 물은 과유불급, 도로 새싹이 된다.
|prop[PLANT]|가 그 단계(0, 2, 4)이고, 위층 구멍으로 보이는 분신
|PLANT2|의 |prop|은 절반 값으로 따라간다. 기름을 주면 식물이 발끈한다.

@<식물에 물을 준다@>=
if g.obj != WATER {
	return g.rep("식물이 발끈해서 잎의 기름을 털어내며 묻네. \"물은?\"")
}
fmt.Fprintf(g.out, "%s\n", objNote[g.prop[PLANT]+1+objOffset[PLANT]])
g.prop[PLANT] += 2
if g.prop[PLANT] > 4 {
	g.prop[PLANT] = 0
}
g.prop[PLANT2] = g.prop[PLANT] >> 1
return g.stayPut()

@ 문의 경첩. 물은 녹을 더할 뿐이고 기름이라야 문이 열린다. 이 문 너머에
보석 박힌 삼지창이 기다린다.

@<문의 경첩에 붓는다@>=
switch g.obj {
case WATER:
	g.prop[DOOR] = 0
	return g.rep("이제 경첩이 완전히 녹슬어서 꿈쩍도 안 해.")
case OIL:
	g.prop[DOOR] = 1
	return g.rep("기름이 경첩을 풀어줘서 이제 문이 열려.")
}
return aDone()

@ 채우기. 병을 채우려면 병이 여기 있고, 비어 있고, 장소에 액체가
흘러야 한다. 무엇이 채워지는지는 장소가 정한다---|caveFlags|의
기름 비트가 그대로 |prop[BOTTLE]|이 된다(기름 2, 물 0).

@<동작@>=
func (g *Game) doFill() actResult {
	if g.obj == VASE {
		@<꽃병을 채워 본다@>@;
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

@ 명나라 꽃병에 물을 받겠다는 발상은 도자기의 물리학이 응징한다.
급격한 온도 변화---산산조각.

@<꽃병을 채워 본다@>=
if g.noLiquidHere() {
	return g.rep("여기엔 꽃병을 채울 만한 게 없어.\n")
}
if !g.toting(VASE) {
	return g.defTo(DROP)
}
fmt.Fprintf(g.out, "급격한 온도 변화에 꽃병이 섬세하게 산산조각 났어.\n")
return g.smashVase()

@ 집기. 붙박이는 못 집고(사연이 있는 셋---잠긴 사슬, 묶인 곰, 뿌리
깊은 식물---은 사연을 말해 준다), 액체는 병이 필요하고, 일곱 개가
한계다. 새와 새장은 세트로 움직인다: 새장에 든 새를 집으면 새장이
따라오고, 새 든 새장을 집어도 마찬가지다. |BIRD+CAGE-obj|는 ``나머지
한쪽''을 꺼내는 산수다.

@<동작@>=
func (g *Game) doTake() actResult {
	if g.toting(g.obj) {
		return g.repDefault() // 이미 들고 있음
	}
	@<붙박이를 집으려는 시도@>@;
	if g.obj == WATER || g.obj == OIL {
		if r, done := g.takeLiquid(); done {
			return r
		}
	}
	if g.holding >= 7 {
		return g.rep("더는 못 들어. 먼저 뭔가 내려놔야 해.")
	}
	if g.obj == BIRD && g.prop[BIRD] == 0 {
		if r, done := g.takeBird(); done {
			return r
		}
	}
	@<새장 세트와 병을 챙긴다@>@;
	return g.defTo(RELAX) // 알았어, 집었어
}

@ 새와 새장은 세트로 움직이고(|BIRD+CAGE-obj|는 ``나머지 한쪽''을 꺼내는 산수다), 액체가 든 병을 집으면 그림자 객체도 손으로 따라온다.

@<새장 세트와 병을 챙긴다@>=
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

@ 붙박이 셋에게는 사연을 말해 준다---잠긴 사슬, 묶인 곰, 뿌리 깊은 식물.

@<붙박이를 집으려는 시도@>=
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

@ 액체 집기의 특수 사정. 그 액체가 든 병이 여기 있으면 병을 집는
것으로 바꾸고(|done=false|: 계속 진행), 빈 병만 들고 있으면 \.{채워}로
변신하고, 병조차 없으면 담아 갈 그릇이 없다.

@<동작@>=
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

@ 새 집기의 특수 사정. 검은 막대를 들고 있으면 새가 겁먹어 못 잡고
---이 게임에서 가장 유명한 함정이다. 막대를 내려놓고 오면 새는 순순히
잡히지만, 데려가려면 새장이 있어야 한다.

@<동작@>=
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

@ 놓기. 단순해 보이지만 이 게임 퍼즐의 절반이 여기서 풀린다. 동전을
자판기 앞에 놓으면 배터리가 나오고, 새를 뱀 앞에 놓으면 뱀을 쫓고,
꽃병은 베개 위가 아니면 깨지고, 곰을 트롤 앞에 풀면 트롤이 달아난다.
지역 변수 |k|는 특수 처리가 이미 제 할 말을 했으니 마지막의
``알았어.''를 누르라는 표지다.

@<동작@>=
func (g *Game) doDrop() actResult {
	if g.obj == ROD && g.toting(ROD2) && !g.toting(ROD) {
		g.obj = ROD2
	}
	if !g.toting(g.obj) {
		return g.repDefault()
	}
	k := false // ``알았어'' 메시지를 누를지
	if g.obj == COINS && g.here(PONY) {
		@<동전을 자판기에 넣는다@>@;
	}
	if g.obj == BIRD {
		if r, done := g.dropBird(&k); done {
			return r
		}
	}
	if g.obj == VASE && g.loc != soft {
		@<꽃병을 바닥에 놓는다@>@;
	}
	if g.obj == BEAR && g.isAtLoc(TROLL) {
		@<곰으로 트롤을 쫓는다@>@;
	}
	@<병을 놓으면 액체도 함께 간다@>@;
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

@ 배터리 자판기는 이 동굴의 수수께끼 중 가장 현대적인 물건이다.
크라우더가 다니던 BBN의 복도 자판기가 미로 속으로 이사 온 것인지도
모른다. 동전은 사라지고(환불 없음), 새 배터리가 떨어진다.

@<동전을 자판기에 넣는다@>=
g.destroy(COINS)
g.drop(BATTERIES, g.loc)
g.prop[BATTERIES] = 0
return g.rep(objNote[objOffset[BATTERIES]])

@ 새를 놓는 순간은 이 게임의 명장면 둘을 품고 있다. 뱀 앞에 놓으면
새가 뱀을 쫓아내 산속 왕의 홀 너머가 열린다---첫 관문의 해답이다.
그러나 용 앞에 놓으면 새는 잿더미가 된다. 원작 힌트의 표현을 빌리면,
새가 통하는 상대와 안 통하는 상대를 가리는 것도 공부다. 폐쇄 후의
소동은 언제나처럼 난쟁이들을 깨운다.

@<동작@>=
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
			"잿더미가 되도록 타 버려. 재가 바람에 흩날려 가."), true
	}
	return actResult{}, false
}

@ 꽃병은 벨벳 베개 위에 놓아야 산다. 맨바닥이면---명나라의 장인에게
잠시 묵념.

@<꽃병을 바닥에 놓는다@>=
if g.place[PILLOW] == g.loc {
	g.prop[VASE] = 0
} else {
	g.prop[VASE] = 2
}
fmt.Fprintf(g.out, "%s\n", objNote[objOffset[VASE]+1+g.prop[VASE]])
k = true
if g.prop[VASE] != 0 {
	objBase[VASE] = VASE
}

@ 순해진 곰을 트롤 앞에 풀어놓으면 트롤이 비명을 지르며 달아난다.
통행료로 던져 준 보물을 되찾는 길이 이것이다.

@<곰으로 트롤을 쫓는다@>=
fmt.Fprintf(g.out, "곰이 트롤 쪽으로 어슬렁어슬렁 다가가자, 트롤은 깜짝 놀라 비명을\n"+
	"지르며 잽싸게 달아나. 곰은 곧 추격을 포기하고 어슬렁어슬렁 돌아와.\n")
k = true
g.destroy(TROLL)
g.destroy(TROLL_)
g.drop(TROLL2, swside)
g.drop(TROLL2_, neside)
g.prop[TROLL] = 2
g.move(BRIDGE, swside)
g.move(BRIDGE_, neside)

@ 병을 놓으면 안에 든 액체도 함께 바닥에 놓인다---물리적으로야
당연하지만, 그림자 객체의 장부는 손으로 맞춰야 한다.

@<병을 놓으면 액체도 함께 간다@>=
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

@ 던지기. 보물을 트롤에게 던지면 통행료가 되고, 음식을 곰에게 던지면
먹이가 되고(동사가 \.{먹여}로 변신한다), 도끼가 아닌 것을 던지는 건
그냥 놓는 것과 같다(``던져 봐야 여기가 동굴이지''라는 태도다). 도끼는
난쟁이·용·트롤·곰 순으로 표적을 찾고, 아무도 없으면 \.{죽여}로
변신해 본다.

@<동작@>=
func (g *Game) doToss() actResult {
	if g.obj == ROD && g.toting(ROD2) && !g.toting(ROD) {
		g.obj = ROD2
	}
	if !g.toting(g.obj) {
		return g.repDefault()
	}
	if isTreasure(g.obj) && g.isAtLoc(TROLL) {
		@<트롤에게 보물을 바친다@>@;
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
			"이렇게 말해. \"솜씨는 좋군, 한데 값어치가 영 모자라.\"\n")
	} else if g.here(BEAR) && g.prop[BEAR] == 0 {
		@<도끼를 곰에게 던진다@>@;
	} else {
		g.obj = NOTHING
		return aChange(KILL)
	}
	g.drop(AXE, g.loc)
	return g.stayPut()
}

@ 트롤은 보물이면 종류를 가리지 않고 낚아채 사라진다. 정석은 황금 알을
바치는 것---어차피 \.{fee fie foe foo}로 되찾을 수 있으니, 트롤만
모르는 환불 제도다.

@<트롤에게 보물을 바친다@>=
g.drop(g.obj, limbo)
g.destroy(TROLL)
g.destroy(TROLL_)
g.drop(TROLL2, swside)
g.drop(TROLL2_, neside)
g.move(BRIDGE, swside)
g.move(BRIDGE_, neside)
return g.rep("트롤이 네 보물을 낚아채 시야 밖으로 잽싸게 달아나.")

@ 사나운 곰에게 도끼를 던지는 것은 나쁜 생각이다. 도끼는 빗나가 곰
옆에 떨어지고, 곰이 지키는 한 주울 수 없는 붙박이가 된다.

@<도끼를 곰에게 던진다@>=
g.drop(AXE, g.loc)
g.prop[AXE] = 1
objBase[AXE] = AXE
if g.place[BEAR] == g.loc {
	g.move(BEAR, g.loc)
}
return g.rep("도끼가 빗나가 곰 옆에 떨어져서, 네가 닿을 수 없는 곳에 가 버렸어.")
