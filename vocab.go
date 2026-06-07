package main

import (
	"fmt"
	"strings"
)

// 우리가 아는 단어들. 한글판에서는 명령어를 한국어로 등록한다.
// 마법 주문(xyzzy, plugh, plover, fee-fie-foe-foo)만 영어로 둔다.

type vocabEntry struct {
	typ     wordtype // 단어의 종류
	meaning int      // motion/object/action 값, 또는 message 인덱스
}

// lookup은 단어를 찾는다. 한글판은 단어 전체로 매칭한다(원본의 앞 5글자
// 절단은 한글 UTF-8에서 글자를 깨뜨리므로 쓰지 않는다).
func (g *Game) lookup(w string) (vocabEntry, bool) {
	e, ok := g.words[w]
	return e, ok
}

// listen은 한두 단어짜리 명령을 읽어 word1, word2에 담는다.
// (advent.w "Low-level input"의 listen). EOF면 false.
func (g *Game) listen() bool {
	for {
		fmt.Fprintf(g.out, "* ")
		line, more := g.readLine()
		if !more {
			return false
		}
		fields := strings.Fields(line)
		switch {
		case len(fields) == 0:
			fmt.Fprintf(g.out, " 뭐라도 시켜 봐.\n")
		case len(fields) > 2:
			fmt.Fprintf(g.out, " 한두 단어로만 말해 줘.\n")
		default:
			g.word1 = strings.ToLower(fields[0])
			g.word2 = ""
			if len(fields) == 2 {
				g.word2 = strings.ToLower(fields[1])
			}
			return true
		}
	}
}

// buildVocabulary는 우리가 아는 모든 단어와 기본 메시지를 채운다.
func (g *Game) buildVocabulary() {
	g.words = make(map[string]vocabEntry)
	mw := func(w string, m motion) { g.words[w] = vocabEntry{motionType, int(m)} }
	ow := func(w string, o object) { g.words[w] = vocabEntry{objectType, int(o)} }
	aw := func(w string, a action) { g.words[w] = vocabEntry{actionType, int(a)} }

	// ---- 이동 단어 ----
	mw("북", N)
	mw("n", N)
	mw("남", S)
	mw("s", S)
	mw("동", E)
	mw("e", E)
	mw("서", W)
	mw("w", W)
	mw("북동", NE)
	mw("ne", NE)
	mw("남동", SE)
	mw("se", SE)
	mw("북서", NW)
	mw("nw", NW)
	mw("남서", SW)
	mw("sw", SW)
	mw("위", U)
	mw("위로", U)
	mw("u", U)
	mw("올라가", U)
	mw("아래", D)
	mw("아래로", D)
	mw("d", D)
	mw("내려가", D)
	mw("왼쪽", L)
	mw("오른쪽", R)
	mw("안", IN)
	mw("안으로", IN)
	mw("나가", OUT)
	mw("밖", OUT)
	mw("밖으로", OUT)
	mw("앞으로", FORWARD)
	mw("계속", FORWARD)
	mw("뒤로", BACK)
	mw("돌아가", BACK)
	mw("너머", OVER)
	mw("가로질러", ACROSS)
	mw("상류", UPSTREAM)
	mw("하류", DOWNSTREAM)
	mw("들어가", ENTER)
	mw("기어가", CRAWL)
	mw("뛰어", JUMP)
	mw("점프", JUMP)
	mw("기어올라", CLIMB)
	mw("등반", CLIMB)
	mw("봐", LOOK)
	mw("둘러봐", LOOK)
	mw("살펴봐", LOOK)
	mw("건너", CROSS)
	mw("길", ROAD)
	mw("언덕", ROAD)
	mw("숲", WOODS)
	mw("계곡", VALLEY)
	mw("집", HOUSE)
	mw("건물", HOUSE)
	mw("도랑", GULLY)
	mw("개울", STREAM)
	mw("시내", STREAM)
	mw("웅덩이", DEPRESSION)
	mw("입구", ENTRANCE)
	mw("동굴", CAVE)
	mw("바위", ROCK)
	mw("석판", SLAB)
	mw("침대", BED)
	mw("통로", PASSAGE)
	mw("터널", PASSAGE)
	mw("동굴방", CAVERN)
	mw("협곡", CANYON)
	mw("어색한", AWKWARD)
	mw("비밀", SECRET)
	mw("베드퀼트", BEDQUILT)
	mw("저수지", RESERVOIR)
	mw("거인", GIANT)
	mw("동양식", ORIENTAL)
	mw("조개방", SHELL)
	mw("황량한", BARREN)
	mw("부서진", BROKEN)
	mw("잔해", DEBRIS)
	mw("전망", VIEW)
	mw("갈림길", FORK)
	mw("구덩이", PIT)
	mw("틈", SLIT)
	mw("균열", CRACK)
	mw("돔", DOME)
	mw("구멍", HOLE)
	mw("벽", WALL)
	mw("홀", HALL)
	mw("방", ROOM)
	mw("바닥", FLOOR)
	mw("계단", STAIRS)
	mw("층계", STEPS)
	mw("자갈", COBBLES)
	mw("지표", SURFACE)
	mw("어둠", DARK)
	mw("낮은", LOW)
	mw("바깥", OUTDOORS)
	mw("y2", Y2)
	mw("xyzzy", XYZZY)
	mw("plugh", PLUGH)
	mw("plover", PLOVER)
	mw("사무실", OFFICE)
	mw("본부", OFFICE)
	mw("아무데도", NOWHERE)

	// ---- 사물 단어 ----
	ow("열쇠", KEYS)
	ow("램프", LAMP)
	ow("등불", LAMP)
	ow("창살", GRATE)
	ow("격자", GRATE)
	ow("새장", CAGE)
	ow("막대", ROD)
	ow("새", BIRD)
	ow("문", DOOR)
	ow("베개", PILLOW)
	ow("쿠션", PILLOW)
	ow("뱀", SNAKE)
	ow("수정", CRYSTAL)
	ow("서판", TABLET)
	ow("대합", CLAM)
	ow("굴", OYSTER)
	ow("잡지", MAG)
	ow("난쟁이", DWARF)
	ow("칼", KNIFE)
	ow("음식", FOOD)
	ow("병", BOTTLE)
	ow("물", WATER)
	ow("기름", OIL)
	ow("거울", MIRROR)
	ow("식물", PLANT)
	ow("콩", PLANT)
	ow("종유석", STALACTITE)
	ow("그림자", SHADOW)
	ow("형상", SHADOW)
	ow("도끼", AXE)
	ow("그림", ART)
	ow("해적", PIRATE)
	ow("용", DRAGON)
	ow("낭떠러지", BRIDGE)
	ow("트롤", TROLL)
	ow("곰", BEAR)
	ow("쪽지", MESSAGE)
	ow("화산", GEYSER)
	ow("간헐천", GEYSER)
	ow("자판기", PONY)
	ow("배터리", BATTERIES)
	ow("건전지", BATTERIES)
	ow("이끼", MOSS)
	ow("금", GOLD)
	ow("금괴", GOLD)
	ow("다이아몬드", DIAMONDS)
	ow("은", SILVER)
	ow("은괴", SILVER)
	ow("보석", JEWELS)
	ow("동전", COINS)
	ow("상자", CHEST)
	ow("보물상자", CHEST)
	ow("알", EGGS)
	ow("둥지", EGGS)
	ow("삼지창", TRIDENT)
	ow("꽃병", VASE)
	ow("도자기", VASE)
	ow("에메랄드", EMERALD)
	ow("피라미드", PYRAMID)
	ow("백금", PYRAMID)
	ow("진주", PEARL)
	ow("양탄자", RUG)
	ow("카펫", RUG)
	ow("향신료", SPICES)
	ow("사슬", CHAIN)

	// ---- 동작 단어 + 기본 메시지 ----
	aw("가져가", TAKE)
	aw("집어", TAKE)
	aw("잡아", TAKE)
	aw("주워", TAKE)
	g.defaultMsg[TAKE] = "이미 들고 있잖아!"
	aw("버려", DROP)
	aw("내려놔", DROP)
	aw("떨어뜨려", DROP)
	g.defaultMsg[DROP] = "그건 들고 있지도 않아!"
	aw("열어", OPEN)
	g.defaultMsg[OPEN] = "그런 걸 어떻게 잠그고 여는지 모르겠어."
	aw("닫아", CLOSE)
	aw("잠가", CLOSE)
	g.defaultMsg[CLOSE] = g.defaultMsg[OPEN]
	aw("켜", ON)
	g.defaultMsg[ON] = "불을 밝힐 게 없어."
	aw("꺼", OFF)
	g.defaultMsg[OFF] = g.defaultMsg[ON]
	aw("흔들어", WAVE)
	g.defaultMsg[WAVE] = "아무 일도 일어나지 않아."
	aw("진정", CALM)
	aw("달래", CALM)
	g.defaultMsg[CALM] = "그래, 해보자.  어떻게 하는지 설명해 줄래?"
	aw("가", GO)
	aw("이동", GO)
	aw("걸어", GO)
	g.defaultMsg[GO] = "어디로?"
	aw("가만", RELAX)
	g.defaultMsg[RELAX] = "알았어."
	aw("부어", POUR)
	aw("따라", POUR)
	g.defaultMsg[POUR] = g.defaultMsg[DROP]
	aw("먹어", EAT)
	g.defaultMsg[EAT] = "말도 안 되는 소리!"
	aw("마셔", DRINK)
	g.defaultMsg[DRINK] = "개울물을 한 모금 마셨어.  물에서 광물 맛이 강하게 나지만,\n" +
		"그리 나쁘진 않아.  엄청나게 차가워."
	aw("문질러", RUB)
	aw("비벼", RUB)
	g.defaultMsg[RUB] = "전기 램프를 문질러 봐야 딱히 신통할 게 없어.  어차피\n" +
		"별일도 안 일어나고."
	aw("던져", TOSS)
	g.defaultMsg[TOSS] = g.defaultMsg[DROP]
	aw("깨워", WAKE)
	g.defaultMsg[WAKE] = g.defaultMsg[EAT]
	aw("먹여", FEED)
	g.defaultMsg[FEED] = "여기엔 먹일 게 없어."
	aw("채워", FILL)
	g.defaultMsg[FILL] = "그건 채울 수 없어."
	aw("부숴", BREAK)
	aw("깨", BREAK)
	g.defaultMsg[BREAK] = "그건 네 힘으로 할 수 있는 일이 아니야."
	aw("폭파", BLAST)
	aw("터뜨려", BLAST)
	g.defaultMsg[BLAST] = "폭파하려면 다이너마이트가 필요해."
	aw("죽여", KILL)
	aw("공격", KILL)
	aw("때려", KILL)
	g.defaultMsg[KILL] = g.defaultMsg[EAT]
	aw("말해", SAY)
	aw("외쳐", SAY)
	aw("read", READ)
	aw("읽어", READ)
	g.defaultMsg[READ] = "미안하지만 무슨 말인지 모르겠어."
	aw("fee", FEEFIE)
	aw("fie", FEEFIE)
	aw("foe", FEEFIE)
	aw("foo", FEEFIE)
	aw("fum", FEEFIE)
	g.defaultMsg[FEEFIE] = "어떻게 하는지 모르겠어."
	aw("간략", BRIEF)
	g.defaultMsg[BRIEF] = "뭘 말이야?"
	aw("찾아", FIND)
	aw("어디", FIND)
	g.defaultMsg[FIND] = "네가 돌아다니고 만지는 것만 알려줄 수 있어.  멀리 있는 게\n" +
		"어디 있는지는 말해줄 수 없어."
	aw("소지품", INVENTORY)
	aw("목록", INVENTORY)
	g.defaultMsg[INVENTORY] = g.defaultMsg[FIND]
	aw("점수", SCORE)
	g.defaultMsg[SCORE] = "응?"
	aw("그만", QUIT)
	aw("종료", QUIT)
	aw("끝내", QUIT)
	g.defaultMsg[QUIT] = g.defaultMsg[SCORE]

	// ---- 메시지 단어 (도움말, 정보 등) ----
	g.buildMessageWords()
}

// buildMessageWords는 고정 메시지를 출력하는 단어들을 등록한다.
// (advent.w의 message_type "Build the vocabulary" 섹션)
func (g *Game) buildMessageWords() {
	k := 0
	mw := func(w string) { g.words[w] = vocabEntry{messageType, k} } // 현재 메시지 슬롯을 가리킴
	nm := func(s string) { g.message[k] = s; k++ }                   // 메시지를 채우고 다음 슬롯으로

	mw("수리수리")
	mw("마수리")
	mw("열려라")
	mw("참깨")
	mw("아브라")
	mw("카다브라")
	mw("얍")
	nm("좋은 시도지만, 그건 낡아빠진 옛 마법 단어야.")
	mw("도움말")
	mw("도와줘")
	mw("?")
	nm("난 장소와 동작, 물건을 알아.  내 어휘 대부분은 장소를 가리키고\n" +
		"널 그곳으로 옮기는 데 써.  움직이려면 숲, 건물, 하류, 들어가,\n" +
		"동, 서, 북, 남, 위, 아래 같은 단어를 써 봐.  동굴에 숨겨진 검은\n" +
		"막대 같은 특별한 물건도 몇 개 알아.  이런 물건은 내가 아는 동작\n" +
		"단어로 다룰 수 있어.  보통 물건과 동작 단어를 (순서는 상관없이)\n" +
		"둘 다 줘야 하지만, 가끔은 동사만으로 물건을 짐작할 수도 있어.\n" +
		"어떤 물건은 동사를 함축하기도 해.  특히 \"소지품\"은 \"소지품\n" +
		"확인\"을 뜻해서, 네가 뭘 들고 있는지 목록을 보여줘.  물건엔 부수\n" +
		"효과가 있어.  예를 들어 막대는 새를 겁줘.  움직이기 힘들 땐 보통\n" +
		"단어를 몇 개 더 시도해 보면 돼.  물건을 다루는 데 자꾸 실패하면,\n" +
		"대개 네(또는 내!) 능력 밖의 일을 하려는 거니 완전히 다른 방법을\n" +
		"써 봐.  게임을 빨리 풀려면 가끔 한 단어로 먼 거리를 이동할 수\n" +
		"있어.  예를 들어 \"건물\"이라고 하면 (숲에서 길을 잃었을 때만\n" +
		"빼고) 보통 지상 어디서든 건물로 가.  또, 동굴 통로는 많이 휘어서,\n" +
		"어떤 방을 북쪽으로 나갔다고 다음 방에 남쪽으로 들어가리란 보장은\n" +
		"없어.  행운을 빌어!")
	mw("나무")
	mw("나무들")
	nm("숲의 나무는 큰 활엽수인 참나무와 단풍나무야, 가끔 소나무나\n" +
		"가문비나무 숲도 있고.  덤불도 꽤 우거졌는데, 주로 자작나무와\n" +
		"물푸레 묘목에 정체 모를 관목들이 섞여 있어.  이맘때면 잎 때문에\n" +
		"시야가 꽤 가리지만, 가문비나무와 산딸기 덤불을 피해 돌아가면\n" +
		"이동은 수월해.")
	mw("파")
	mw("발굴")
	nm("삽 없이 파는 건 영 비현실적이야.  삽이 있어도 진척은\n" +
		"거의 없을걸.")
	mw("길잃음")
	nm("나도 너만큼 헷갈려.")
	nm("요란한 폭발이 일어나고 넌 순식간에 방 벽에\n" +
		"흩뿌려져.")
	nm("요란한 폭발이 일어나고 저쪽 벽에 20피트 구멍이 뚫리면서, 뱀들이\n" +
		"잔해에 파묻혀.  녹은 용암의 강이 구멍으로 쏟아져 들어와, 너를\n" +
		"포함해 앞을 가로막는 모든 걸 파괴해!")
	mw("안개")
	nm("안개는 흰 수증기인데, 보통 물이고, 동굴에서 이따금 보여.  어디서든\n" +
		"나타날 수 있지만 흔히 물로 이어지는 깊은 구덩이의 신호야.")
	mw("젠장")
	nm("말조심해!")
	nm("요란한 폭발이 일어나고 저쪽 벽에 20피트 구멍이 뚫리면서, 난쟁이들이\n" +
		"잔해에 파묻혀.  넌 구멍을 통해 당당히 걸어 나가 본부에 다다라.\n" +
		"그곳에선 환호하는 다정한 요정 무리가 정복자 모험가인 너를 노을\n" +
		"속으로 데려가.")
	mw("멈춰")
	nm("\"멈춰\"라는 단어는 몰라.  그만두고 싶으면 \"그만\"을 써.")
	mw("정보")
	mw("안내")
	nm("모험을 일찍 끝내고 싶으면 \"그만\"이라고 해.  보물의 점수를 온전히\n" +
		"받으려면 건물에 안전하게 둬야 해. 그냥 찾기만 해도 부분 점수는 줘.\n" +
		"죽거나 그만두면 점수가 깎이는데, 죽는 쪽이 더 손해야.  동굴을\n" +
		"(조금이라도) 얼마나 탐험했는지에 따른 점수도 있어. 특히 동굴에\n" +
		"들어가는 것만으로 큰 보너스를 주지 (초보와 나머지를 구분하려고).\n" +
		"그리고 더 험난한 구간들을 지나왔는지 가리는 다른 방법들도 있어.\n" +
		"보물을 다 찾은 것 같으면, 그냥 한동안 더 탐험해 봐.  별일 없으면\n" +
		"아직 다 못 찾은 거야.  뭔가 흥미로운 일이 일어난다면, 보너스를\n" +
		"받고 있고 고수 구간에서 점수를 훨씬 더 딸 기회가 생겼다는 뜻이야.\n" +
		"네가 헤매는 것 같으면 가끔 힌트를 줄게.  그럴 땐 점수에 얼마나\n" +
		"영향을 줄지 미리 경고할게.  끝으로, 종이를 아끼려면 \"간략\"이라고\n" +
		"할 수 있어. 그러면 어떤 장소든 네가 직접 다시 보자고 하지 않는 한\n" +
		"전체 설명을 되풀이하지 않을게.")
	mw("수영")
	mw("헤엄")
	nm("어떻게 하는지 모르겠어.")
}
