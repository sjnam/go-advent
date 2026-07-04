@* 동굴 폐쇄.
보물 열다섯을 모두 찾아낸 순간부터 --- 건물로 다 나르지 못했어도 ---
동굴의 시계가 움직이기 시작한다. |clock1|이 다 되면 음산한 폐쇄 예고가
울리고, 이어 |clock2|가 다 되면 동굴이 정말로 닫히며 최후의 퍼즐이
열린다. 우즈의 원안에는 대학 전산실의 황금 시간대에 동굴이 문을 닫는
장치까지 있었다는데(매카시가 AI 연구실의 생산성 급락을 보고 우겼다는
설이 유력하다), 커누스가 CWEB판에서 들어냈다. 우리의 폐쇄는 오직
게임 안의 사건이다.

@<게임의 상태@>=
clock1   int  // 모든 보물을 본 뒤 폐쇄 예고까지 (시작 15)
clock2   int  // 예고 뒤 동굴이 닫힐 때까지 (시작 30)
panicked bool // 폐쇄 중 탈출을 시도해 당황했는가
warned   bool // 램프 전력 부족 경고를 받았는가
closed   bool // 동굴이 완전히 닫혔는가

@ 시계는 매 턴 명령 해석 초입에서 돈다. |clock1|은 보물을 다 본 채
깊은 동굴에 있을 때만 줄어든다 --- 지상에서 밍기적거리는 시간은 폐쇄를
재촉하지 않는다. 동굴이 닫히는 턴에는 |true|를 돌려주어 제자리 이동으로
빠져나가게 한다.

@<동굴 폐쇄@>=
func (g *Game) checkClocks() bool {
	if g.tally == 0 && g.loc >= minLowerLoc && g.loc != y2 {
		g.clock1--
	}
	if g.clock1 == 0 {
		@<폐쇄가 임박했다고 경고한다@>@;
		return false
	}
	if g.clock1 < 0 {
		g.clock2--
	}
	if g.clock2 == 0 {
		@<동굴을 닫는다@>@;
		return true
	}
	@<램프를 점검한다@>@;
	return false
}

@ 폐쇄 예고와 함께 동굴은 정리 수순에 든다. 격자문이 잠기고, 수정
다리가 사라지고, 난쟁이와 해적이 모두 물러가고, 트롤 다리는 온전한
모습으로 복원된다. 사슬과 도끼는 자유로워진다 --- 어차피 이제 아무도
지키지 않으니.

@<폐쇄가 임박했다고 경고한다@>=
fmt.Fprintf(g.out, "동굴 전체에 울려 퍼지는 음산한 목소리가 말해.  \"곧 동굴이\n"+
	"닫힙니다.  모든 모험가는 즉시 본부를 통해 나가십시오.\"\n")
g.clock1 = -1
g.prop[GRATE] = 0
g.prop[CRYSTAL] = 0
for j := 0; j <= nd; j++ {
	g.dseen[j] = false
	g.dloc[j] = limbo
}
g.destroy(TROLL)
g.destroy(TROLL_)
g.move(TROLL2, swside)
g.move(TROLL2_, neside)
g.move(BRIDGE, swside)
g.move(BRIDGE_, neside)
if g.prop[BEAR] != 3 {
	g.destroy(BEAR)
}
g.prop[CHAIN] = 0
objBase[CHAIN] = NOTHING
g.prop[AXE] = 0
objBase[AXE] = NOTHING

@ 폐쇄 중에 밖으로 나가려 들면 녹음된 목소리가 타이르고, 처음
당황했을 때는 남은 시간을 15턴으로 줄여 준다 --- 어차피 못 나갈 거,
빨리 끝나기나 하라는 배려다.

@<동굴 폐쇄@>=
func (g *Game) panicClosing() {
	if !g.panicked {
		g.clock2 = 15
		g.panicked = true
	}
	fmt.Fprintf(g.out, "수상한 녹음된 목소리가 끼익 살아나며 알려:\n"+
		"\"이 출구는 닫혔습니다.  본부를 통해 나가 주십시오.\"\n")
}

@ 마침내 동굴이 닫히면, 너는 번쩍이는 빛(과 작은 주황색 연기 한 줄기)
속에서 최후의 보관실로 옮겨진다. 북동 끝에는 여행의 도구들이, 남서
끝에는 동굴의 옛 주민들이 --- 뱀과 새와 격자문까지 --- 박물관처럼
진열된다. 음수 |prop|은 ``내려놨다 집기 전에는 묘사하지 않음''의
표지다. 들고 있던 것은 모두 사라진다. 빈손으로 왔으니 빈손으로 가는
것이 동굴의 법이지만, 퇴장 보너스 10점은 챙겨 준다.

@<동굴을 닫는다@>=
fmt.Fprintf(g.out, "음산한 목소리가 읊조려.  \"동굴은 이제 닫혔다.\"  메아리가 잦아들자,\n"+
	"눈부신 빛이 번쩍이고 (주황색 연기도 작게 한 줄기 피어올라).\n"+
	". . .    그러다 다시 초점이 잡히고, 둘러보니...\n")
g.move(BOTTLE, neend)
g.prop[BOTTLE] = -2
g.move(PLANT, neend)
g.prop[PLANT] = -1
g.move(OYSTER, neend)
g.prop[OYSTER] = -1
g.move(LAMP, neend)
g.prop[LAMP] = -1
g.move(ROD, neend)
g.prop[ROD] = -1
g.move(DWARF, neend)
g.prop[DWARF] = -1
g.move(MIRROR, neend)
g.prop[MIRROR] = -1
g.loc, g.oldloc = neend, neend
g.move(GRATE, swend) // prop[GRATE]는 여전히 0
g.move(SNAKE, swend)
g.prop[SNAKE] = -2
g.move(BIRD, swend)
g.prop[BIRD] = -2
g.move(CAGE, swend)
g.prop[CAGE] = -1
g.move(ROD2, swend)
g.prop[ROD2] = -1
g.move(PILLOW, swend)
g.prop[PILLOW] = -1
g.move(MIRROR_, swend)
g.place[WATER] = limbo // bugfix
g.place[OIL] = limbo
for j := object(1); j <= maxObj; j++ {
	if g.toting(j) {
		g.destroy(j)
	}
}
g.closed = true
g.bonus = 10

@ 램프도 제 시계를 산다. 켜져 있으면 매 턴 |limit|가 준다. 30턴이
남으면 경고가 나오는데, 마침 새 배터리가 곁에 있으면 내가 알아서
갈아 끼워 준다 --- 이 게임에서 가장 친절한 자동화다. 경고문은
배터리의 행방에 따라 세 갈래다: 이미 써 버렸거나, 아직 자판기에
있거나(동전을 챙기라는 힌트까지 얹어 준다), 어딘가 두고 왔거나.

@<램프를 점검한다@>=
if g.prop[LAMP] == 1 {
	g.limit--
}
switch {
case g.limit <= 30 && g.here(BATTERIES) && g.prop[BATTERIES] == 0 && g.here(LAMP):
	fmt.Fprintf(g.out, "램프가 어두워지고 있어.  내가 알아서 배터리를\n"+
		"갈아 끼울게.\n")
	g.prop[BATTERIES] = 1
	if g.toting(BATTERIES) {
		g.drop(BATTERIES, g.loc)
	}
	g.limit = 2500
case g.limit == 0:
	g.limit = -1
	g.prop[LAMP] = 0
	if g.here(LAMP) {
		fmt.Fprintf(g.out, "램프 전력이 다 떨어졌어.")
	}
case g.limit < 0 && g.loc < minInCave:
	fmt.Fprintf(g.out, "여기 밖에서 어슬렁거려 봐야 별 의미 없고, 램프 없이는\n"+
		"동굴을 탐험할 수도 없어.  그냥 오늘은 여기까지 하자.\n")
	g.gaveUp = true
	g.quitting = true
case g.limit <= 30 && !g.warned && g.here(LAMP):
	fmt.Fprintf(g.out, "램프가 어두워지고 있어")
	switch {
	case g.prop[BATTERIES] == 1:
		fmt.Fprintf(g.out, ", 게다가 여분 배터리도 없어.  슬슬\n"+
			"마무리하는 게 좋겠어.\n")
	case g.place[BATTERIES] == limbo:
		fmt.Fprintf(g.out, ".  슬슬 마무리하는 게 좋겠어, 새 배터리를\n"+
			"찾지 못한다면 말이야.  내 기억에 미로 어딘가에\n"+
			"자판기가 있었어.  동전을 좀 챙겨 가.\n")
	default:
		fmt.Fprintf(g.out, ".  그 배터리를 가지러 돌아가는 게 좋겠어.\n")
	}
	g.warned = true
}

@ 하나 더. 남은 보물이 모두 ``영영 못 볼 보물''이 된 순간 --- 새를
죽여 뱀을 못 쫓게 됐다든가 --- 게임은 램프 수명을 35턴으로 끊어
빨리 끝나게 돕는다. 잔인함이 아니라 자비다.

@<남은 보물이 너무 꼭꼭 숨었으면 램프를 재촉한다@>=
if g.tally == g.lostTreasures && g.tally > 0 && g.limit > 35 {
	g.limit = 35
}
