@* 객체.
동굴에는 예순일곱 가지 물건이 있다. 열쇠처럼 들고 다니는 것, 뱀처럼
꿈쩍 않는 것, 그리고 보물 열다섯 가지. 밑줄로 끝나는 이름(|GRATE_|,
|TROLL_| 같은)은 새 물건이 아니라 같은 물건의 ``두 번째 조각''이다 ---
격자문은 웅덩이 쪽과 동굴 쪽 두 장소에서 동시에 보여야 하니, 조각을
둘 두고 한 몸으로 묶는다. 여기서도 값의 순서는 신성하다: 같은 물건의
조각들은 연속한 값이어야 하고, |GOLD|부터 끝까지가 보물이라는 사실
자체가 코드 곳곳에서 쓰인다.

@<객체@>=
type object int

const (
	NOTHING object = iota
	KEYS
	LAMP
	GRATE
	GRATE_
	CAGE
	ROD
	ROD2
	TREADS
	TREADS_
	BIRD
	DOOR
	PILLOW
	SNAKE
	CRYSTAL
	CRYSTAL_
	TABLET
	CLAM
	OYSTER
	MAG
	DWARF
	KNIFE
	FOOD
	BOTTLE
	WATER
	OIL
	MIRROR
	MIRROR_
	PLANT
	PLANT2
	PLANT2_
	STALACTITE
	SHADOW
	SHADOW_
	AXE
	ART
	PIRATE
	DRAGON
	DRAGON_
	BRIDGE
	BRIDGE_
	TROLL
	TROLL_
	TROLL2
	TROLL2_
	BEAR
	MESSAGE
	GEYSER
	PONY
	BATTERIES
	MOSS
	GOLD
	DIAMONDS
	SILVER
	JEWELS
	COINS
	CHEST
	EGGS
	TRIDENT
	VASE
	EMERALD
	PYRAMID
	PEARL
	RUG
	RUG_
	SPICES
	CHAIN
)

@ @<객체@>=
const (
	minTreasure = GOLD
	maxObj      = CHAIN
)

func isTreasure(t object) bool { return t >= minTreasure }

@ 물건의 살림은 배열 넷과 계수기 하나로 꾸린다. |place|는 물건이 있는
장소(들고 있으면 |inhand|, 세상에 없으면 |limbo|), |prop|은 물건의
상태값(문이 열렸는지, 새가 새장에 있는지$\ldots$), 그리고 같은 장소의
물건들은 |first|와 |link|로 엮인 연결 리스트에 매달린다. 1970년대의
자료구조가 반세기를 건너 그대로 왔다 --- 리스트에 앞에서 끼우니, 나중에
떨어뜨린 물건이 먼저 눈에 띈다는 자연스러운 부수 효과까지 그대로다.

@<게임의 상태@>=
place   [maxObj + 1]location // 각 객체의 현재 위치
prop    [maxObj + 1]int      // 각 객체의 현재 속성값
link    [maxObj + 1]object   // 같은 장소의 다음 객체
first   [maxLoc + 1]object   // 각 장소의 첫 객체
holding int                  // 들고 있는 객체 수

@ 판정 둘: 들고 있는가, 여기 있는가.

@<객체@>=
func (g *Game) toting(t object) bool { return g.place[t] < 0 }
func (g *Game) here(t object) bool   { return g.toting(t) || g.place[t] == g.loc }

@ |drop|은 (리스트에 없던) 물건을 장소 |l|에 놓고, |carry|는 물건을
집어 들며 리스트에서 떼어 낸다. |carry|의 안쪽 |for|는 몸통 없이
돌며 리스트에서 |t|의 앞 원소를 찾는다 --- 커누스 세대의 포인터 춤이다.
|move|는 그 둘의 합성이고, |destroy|는 |limbo|로 보내는 완곡어법이다.

@<객체@>=
func (g *Game) drop(t object, l location) {
	if g.toting(t) {
		g.holding--
	}
	g.place[t] = l
	if l < 0 {
		g.holding++
	} else if l > 0 {
		g.link[t] = g.first[l]
		g.first[l] = t
	}
}

@ @<객체@>=
func (g *Game) carry(t object) {
	l := g.place[t]
	if l >= limbo {
		g.place[t] = inhand
		g.holding++
		if l > limbo {
			var r, s object
			for r, s = 0, g.first[l]; s != t; r, s = s, g.link[s] {
			}
			if r == 0 {
				g.first[l] = g.link[s]
			} else {
				g.link[r] = g.link[s] // 리스트에서 t 제거
			}
		}
	}
}

func (g *Game) move(t object, l location) { g.carry(t); g.drop(t, l) }
func (g *Game) destroy(t object)          { g.move(t, limbo) }

@ 여러 조각으로 된 물건이 ``여기 있는가''는 조각을 훑어야 안다.
조각들의 대표를 |objBase|가 알려 주는데, 움직일 수 있는 보통 물건은
|NOTHING|, 붙박이는 제 대표 조각이다. 이 맵은 생성 파일
\.{objdata.go}에 있지만 game 중에 바뀌기도 한다 --- 용을 잡으면
양탄자가 붙박이에서 보물로 풀려나는 식이다.

@<객체@>=
func (g *Game) isAtLoc(t object) bool {
	if objBase[t] == NOTHING {
		return g.place[t] == g.loc
	}
	for tt := t; objBase[tt] == t; tt++ {
		if g.place[tt] == g.loc {
			return true
		}
	}
	return false
}

@ 물건들의 첫 배치 역시 검증된 C~원본에서 덤프한 초기표
(\.{objdata.go}의 |objInitPlace| 등)를 그대로 씌운다. 원작은 게임을
시작할 때마다 |new_obj| 호출 예순일곱 번으로 이 상태를 지었지만, 그
결과는 늘 같으니 결과만 저장해 둔 것이다.

@<객체@>=
func (g *Game) loadObjects() {
	g.place = objInitPlace
	g.prop = objInitProp
	g.first = objInitFirst
	g.link = objInitLink
	g.holding = objInitHolding
}

@ 물건의 텍스트는 사람 몫이다. 먼저 이름 --- 소지품 목록에 찍히는
문자열이다. 이름이 없는 물건(뱀, 용, 트롤$\ldots$)은 들고 다닐 수 없는
것들이라 목록에 나올 일도 없다.

@<객체@>=
var objName = map[object]string{
	KEYS:      "열쇠 꾸러미",
	LAMP:      "놋쇠 램프",
	CAGE:      "버들 새장",
	ROD:       "검은 막대",
	ROD2:      "검은 막대",
	BIRD:      "새장 속 작은 새",
	PILLOW:    "벨벳 베개",
	CLAM:      "거대한 대합 >끙!<",
	OYSTER:    "거대한 굴 >끄응!<",
	MAG:       "\"오늘의 동굴탐험가\"",
	FOOD:      "맛있는 음식",
	BOTTLE:    "작은 병",
	WATER:     "병 속의 물",
	OIL:       "병 속의 기름",
	AXE:       "난쟁이의 도끼",
	BATTERIES: "배터리",
	GOLD:      "커다란 금덩이",
	DIAMONDS:  "다이아몬드 여러 개",
	SILVER:    "은괴",
	JEWELS:    "값진 보석",
	COINS:     "희귀한 동전",
	CHEST:     "보물 상자",
	EGGS:      "황금 알",
	TRIDENT:   "보석 박힌 삼지창",
	VASE:      "명나라 꽃병",
	EMERALD:   "달걀만 한 에메랄드",
	PYRAMID:   "백금 피라미드",
	PEARL:     "반짝이는 진주",
	RUG:       "페르시아 양탄자",
	SPICES:    "희귀한 향신료",
	CHAIN:     "황금 사슬",
}

@ 그리고 묘사 --- 방에 들어섰을 때 물건이 제 상태에 따라 들려주는
문장들이다. 물건 |t|의 상태 |prop[t]|에 따른 묘사가
|objNote[objOffset[t]+prop[t]]|에 있다. 오프셋표는 \.{objdata.go}에서
오지만 문장은 여기 있다. 빈 문자열은 ``그 상태에서는 아무 말 없음''이다.

@<객체@>=
var objNote = []string{
	"바닥에 페르시아 양탄자가 펼쳐져 있어!",
	"용이 페르시아 양탄자 위에 큰대자로 뻗어 있어!!",
	"트롤은 어디에도 보이지 않아.",
	"건장한 트롤이 다리 옆에 서서, 건너려면 보물을 하나 던지라고\n" +
		"으름장을 놔.",
	"트롤이 다리 밑에서 나와 네 길을 막아.",
	"",
	"낡은 나무다리가 협곡을 가로질러 안개 속으로 사라져.  다리에 걸린\n" +
		"표지판엔 \"멈춰!  트롤에게 통행료를 내라!\"라고 적혀 있어.",
	"협곡 바닥에 부서진 다리(와 죽은 곰)가 보여.",
	"거대하고 사나운 초록 용이 길을 막아!",
	"축하해!  맨손으로 용을 무찔렀어!  (믿기지 않지?)",
	"거대한 초록 용의 시체가 한쪽에 널브러져 있어.",
	"흐릿한 형상이 네 주의를 끌려는 것 같아.",
	"",
	"12피트짜리 콩나무 꼭대기가 서쪽 구덩이에서 삐져나와 있어.",
	"서쪽 구덩이에서 거대한 콩나무가 구멍까지 자라 있어.",
	"",
	"이제 수정 다리가 틈을 가로질러 놓여 있어.",
	"수정 다리가 사라졌어!",
	"거친 돌계단이 구덩이 아래로 이어져.",
	"거친 돌계단이 돔 위로 이어져.",
	"창살이 잠겨 있어.",
	"창살이 열려 있어.",
	"바닥에 황금 사슬이 무더기로 놓여 있어!",
	"곰이 황금 사슬로 벽에 묶여 있어!",
	"황금 사슬이 벽에 잠겨 있어!",
	"여기 희귀한 향신료가 있어!",
	"한쪽에 반짝이는 진주가 놓여 있어!",
	"여기 한 변이 8인치인 백금 피라미드가 있어!",
	"여기 물떼새 알만 한 에메랄드가 있어!",
	"여기 섬세하고 값진 명나라 꽃병이 있어!",
	"이제 꽃병이 벨벳 베개 위에 살포시 놓여 있어.",
	"바닥엔 쓸모없는 도자기 파편이 널려 있어.",
	"명나라 꽃병이 섬세하게 쨍그랑 떨어져.",
	"여기 보석이 박힌 삼지창이 있어!",
	"여기 황금 알이 가득한 큰 둥지가 있어!",
	"황금 알 둥지가 사라졌어!",
	"끝!",
	"해적의 보물 상자가 여기 있어!",
	"여기 동전이 잔뜩 있어!",
	"여기 값진 보석이 있어!",
	"여기 은괴가 있어!",
	"여기 다이아몬드가 있어!",
	"여기 크고 반짝이는 금덩이가 있어!",
	"",
	"여기 새 배터리가 있어.",
	"다 쓴 배터리 몇 개가 근처에 버려져 있어.",
	"여기 거대한 자판기가 있어.  거기 적힌 설명은 이래.\n" +
		"\"새 배터리를 받으려면 여기 동전을 넣으시오.\"",
	"",
	"먼지 위에 화려한 필체로 휘갈긴 쪽지가 있어.  이렇게 적혀 있어.\n" +
		"\"여긴 해적이 보물 상자를 두는 미로가 아니다.\"",
	"방 저쪽 끝에서 사나운 동굴 곰이 널 노려보고 있어!",
	"온순한 동굴 곰이 한구석에 얌전히 앉아 있어.",
	"만족스러운 표정의 곰이 근처를 어슬렁거려.",
	"",
	"",
	"",
	"여기 작은 도끼가 있어.",
	"곰 옆에 작은 도끼가 놓여 있어.",
	"",
	"구덩이 속 아주 작은 식물이 \"물, 물, ...\" 하고 중얼거려.",
	"식물이 몇 초간 맹렬하게 자라나.",
	"12피트짜리 콩나무가 구덩이에서 솟아올라 \"물!!  물!!\" 하고 고함쳐.",
	"식물이 폭발하듯 자라나 구덩이 바닥을 거의 채워.",
	"거대한 콩나무가 구멍까지 쭉 자라 있어.",
	"물을 너무 줬어!  식물이 쪼그라들어!  이, 이런...",
	"",
	"여기 물이 든 병이 있어.",
	"여기 빈 병이 있어.",
	"여기 기름이 든 병이 있어.",
	"여기 음식이 있어.",
	"여기 \"오늘의 동굴탐험가\" 잡지 최신호 몇 권이 있어.",
	"여기 껍데기를 꽉 다문 거대한 굴이 있어.",
	"흥미로운걸.  굴 밑면에 뭔가 적혀 있는 것 같아.",
	"여기 껍데기를 꽉 다문 거대한 대합이 있어.",
	"벽에 박힌 육중한 돌판엔 이렇게 적혀 있어.\n" +
		"\"어둠의 방에 빛을 가져온 걸 축하한다!\"",
	"거대하고 사나운 초록 뱀이 길을 막아!",
	"",
	"작은 벨벳 베개가 바닥에 놓여 있어.",
	"북쪽 길은 육중하고 녹슨 철문으로 막혀 있어.",
	"북쪽 길은 육중하고 녹슨 철문을 통해 이어져.",
	"쾌활한 작은 새가 여기 앉아 노래하고 있어.",
	"새장 안에 작은 새가 있어.",
	"끝에 녹슨 자국이 있는 3피트짜리 검은 막대가 근처에 놓여 있어.",
	"끝에 녹슨 별이 있는 3피트짜리 검은 막대가 근처에 놓여 있어.",
	"작은 버들 새장이 근처에 버려져 있어.",
	"근처에 반짝이는 놋쇠 램프가 있어.",
	"근처에 램프가 빛나고 있어.",
	"바닥에 열쇠 몇 개가 있어.",
}

@ 객체의 시험. 보물의 초기 |prop|이 $-1$인지(처음 발견을 셈하는 표지다),
우물집에 보급품 네 가지가 놓였는지, |carry|/|drop|이 장부를 제대로
맞추는지 본다.

@(advent_test.go@>=
func TestObjectData(t *testing.T) {
	g := &Game{}
	g.loadObjects()
	// 보물은 prop=-1로 시작 (advent.w:2230)
	if g.prop[GOLD] != -1 {
		t.Errorf("prop[GOLD]=%d, want -1 (보물)", g.prop[GOLD])
	}
	if g.prop[KEYS] != 0 {
		t.Errorf("prop[KEYS]=%d, want 0 (비보물)", g.prop[KEYS])
	}
	// 집 안에 키/램프/음식/병이 있어야 함 (well house, house=3)
	g.loc = house
	for _, o := range []object{KEYS, LAMP, FOOD, BOTTLE} {
		if !g.here(o) {
			t.Errorf("%s가 house에 없음 (place=%d)", objName[o], g.place[o])
		}
	}
	// SNAKE는 움직이지 않는 그룹 객체 (base[SNAKE]==SNAKE)
	if objBase[SNAKE] != SNAKE {
		t.Errorf("base[SNAKE]=%d, want SNAKE", objBase[SNAKE])
	}
	// KEYS는 움직임 (base==NOTHING)
	if objBase[KEYS] != NOTHING {
		t.Errorf("base[KEYS]=%d, want NOTHING", objBase[KEYS])
	}
	// carry/drop 동작
	h0 := g.holding
	g.carry(KEYS)
	if !g.toting(KEYS) || g.holding != h0+1 {
		t.Errorf("carry(KEYS) 실패: toting=%v holding=%d", g.toting(KEYS), g.holding)
	}
	g.drop(KEYS, road)
	if g.toting(KEYS) || g.place[KEYS] != road {
		t.Errorf("drop(KEYS,road) 실패: place=%d", g.place[KEYS])
	}
}
