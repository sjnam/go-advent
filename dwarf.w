@* 난쟁이와 해적.
여태 난쟁이 이야기를 두루뭉술하게만 했는데, 이제 털어놓을 때다. 동굴
깊은 곳에는 난쟁이 다섯이 떠돌고, 여섯째 자리(0번)에는 해적이 있다.
난쟁이는 칼을 던지고, 해적은 보물을 훔쳐 미로 깊숙한 상자에 숨긴다.
훔쳐 가는 순간은 부아가 치밀지만 사실 선물이다---그 상자 자체가
열다섯 번째 보물이라, 해적이 숨겨 주지 않으면 찾을 수도 없다.

@<난쟁이와 해적@>=
const (
	nd           = 5     // 난쟁이 수
	chestLoc     = dead2 // 해적의 보물 상자가 숨는 곳
	messageLoc   = pony  // 해적을 본 뒤 단서가 놓이는 곳
	maxPirateLoc = dead2 // 해적이 갈 수 있는 가장 깊은 곳
)

var dwarfStart = [nd + 1]location{chestLoc, hmk, wfiss, y2, like3, complex}

@ 그들의 장부. |dflag|는 난쟁이 사태의 단계다: 0은 아직 못 만남, 1은
깊은 동굴에 들어섬, 2는 첫 조우(도끼가 날아온다), 3부터는 칼 명중률이
올라가는 분노 단계다. |knifeLoc|은 칼이 마지막으로 ``언급''된 곳 ---
칼을 주우려는 시도에 한 번 퇴짜를 놓기 위한 것으로, $-1$이면 언급된
칼이 없다.

@<게임의 상태@>=
dflag    int              // 난쟁이가 얼마나 성났나
dkill    int              // 난쟁이를 몇 마리 죽였나
dloc     [nd + 1]location // 각 난쟁이(0번은 해적)의 현재 위치
odloc    [nd + 1]location // 직전 위치
dseen    [nd + 1]bool     // 난쟁이가 나를 봤는가
dtotal   int              // 같은 방에 있는 난쟁이 수
attack   int              // 칼을 뽑을 틈이 있던 난쟁이 수
stick    int              // 칼을 명중시킨 난쟁이 수
knifeLoc int              // 칼이 언급된 곳, 없으면 -1

@ 난쟁이가 (해적 말고) 지금 이 방에 있는가.

@<난쟁이와 해적@>=
func (g *Game) dwarf() bool {
	if g.dflag < 2 {
		return false
	}
	for j := 1; j <= nd; j++ {
		if g.dloc[j] == g.loc {
			return true
		}
	}
	return false
}

@ 네가 한 걸음 옮길 때마다 그들도 움직인다. 해적이 못 가는 깊이거나
림보라면 아무도 안 움직이고, |dflag|의 단계에 따라 하는 일이 다르다.
1단계에서 2단계로 오르는 것은 깊은 동굴에서 5\%의 복권이 당첨될 때다.

@<난쟁이와 해적@>=
func (g *Game) moveDwarves() {
	if g.loc > maxPirateLoc || g.loc == limbo {
		return
	}
	switch {
	case g.dflag == 0:
		if g.loc >= minLowerLoc {
			g.dflag = 1
		}
	case g.dflag == 1:
		if g.loc >= minLowerLoc && g.pct(5) {
			@<dflag를 2로 올리며 첫 조우를 연출한다@>@;
		}
	default:
		@<난쟁이들과 해적을 실제로 움직인다@>@;
	}
}

@ 첫 조우. 난이도를 고르듯 난쟁이 0$\sim$2마리를 소리 없이 지우고
(50\% 동전을 두 번 던진다), 나와 같은 방의 난쟁이는 너깃 방으로
물린 다음, 한 녀석이 모퉁이에서 도끼를 던지고 달아난다. 도끼는
빗나가서 바닥에 떨어진다---이 도끼가 앞으로 너의 유일한 무기다.

@<dflag를 2로 올리며 첫 조우를 연출한다@>=
g.dflag = 2
for j := 0; j < 2; j++ {
	if g.pct(50) {
		g.dloc[1+g.ran(nd)] = limbo
	}
}
for j := 1; j <= nd; j++ {
	if g.dloc[j] == g.loc {
		g.dloc[j] = nugget
	}
	g.odloc[j] = g.dloc[j]
}
fmt.Fprintf(g.out, "작은 난쟁이 하나가 모퉁이를 돌다 널 보더니, 작은 도끼를\n"+
	"네게 던지고는, 욕을 내뱉으며 달아났어. (도끼는 빗나갔어.)\n")
g.drop(AXE, g.loc)

@ 본격 추격전. 살아 있는 각 난쟁이는 갈 수 있는 이웃 칸 목록을 만들어
무작위로 한 칸 옮긴다. 갈 곳이 없으면 왔던 곳으로 되돌아간다. 나를
봤는지(|dseen|)는 새 위치나 옛 위치가 나와 겹치는지, 그리고 깊은
동굴에서 이미 나를 쫓고 있었는지로 정한다. 본 녀석들은 따라붙고,
방에 모인 위협의 머릿수만큼 칼부림이 벌어진다.

@<난쟁이들과 해적을 실제로 움직인다@>=
g.dtotal, g.attack, g.stick = 0, 0, 0
for j := 0; j <= nd; j++ {
	if g.dloc[j] == limbo {
		continue
	}
	ploc := g.dwarfExits(j)
	if len(ploc) == 0 {
		ploc = []location{g.odloc[j]}
	}
	g.odloc[j] = g.dloc[j]
	g.dloc[j] = ploc[g.ran(len(ploc))] // 무작위 걸음
	g.dseen[j] = g.dloc[j] == g.loc || g.odloc[j] == g.loc ||
		(g.dseen[j] && g.loc >= minLowerLoc)
	if g.dseen[j] {
		g.dwarfFollow(j)
	}
}
if g.dtotal != 0 {
	@<위협하는 난쟁이들이 칼을 던진다@>@;
}

@ 난쟁이 |j|가 갈 만한 곳의 목록. 이동표를 그대로 쓰되 깊은 동굴만
허락하고, 왔던 곳으로 즉시 되밟기와 제자리걸음은 빼고, 난쟁이 금지
조건(|cond==100|)을 지키고, 해적(|j==0|)은 제 구역 밖으로 내보내지
않는다. 목록 길이 19 제한은 원작의 고정 배열 크기가 남긴 화석이다.

@<난쟁이와 해적@>=
func (g *Game) dwarfExits(j int) []location {
	limit := minForcedLoc - 1
	if j == 0 {
		limit = maxPirateLoc
	}
	var ploc []location
	for q := caveStart[g.dloc[j]]; q < caveStart[g.dloc[j]+1]; q++ {
		nl := caveTravels[q].dest
		if nl >= minLowerLoc && nl != g.odloc[j] && nl != g.dloc[j] &&
			(len(ploc) == 0 || nl != ploc[len(ploc)-1]) && len(ploc) < 19 &&
			caveTravels[q].cond != 100 && nl <= limit {
			ploc = append(ploc, nl)
		}
	}
	return ploc
}

@ 나를 본 난쟁이는 내 방까지 따라온다. 해적이면 보물 추적으로 넘어가고,
난쟁이면 위협 인원에 든다. 지난 턴에도 같은 방이었다면 칼을 뽑을 틈이
있었던 것이고(|attack|), 분노 단계가 높을수록 명중까지 한다(|stick|) ---
확률 $95(|dflag|-2)/1000$.

@<난쟁이와 해적@>=
func (g *Game) dwarfFollow(j int) {
	g.dloc[j] = g.loc
	if j == 0 {
		g.pirateTrack()
		return
	}
	g.dtotal++
	if g.odloc[j] == g.dloc[j] {
		g.attack++
		if g.knifeLoc >= 0 {
			g.knifeLoc = int(g.loc)
		}
		if g.ran(1000) < 95*(g.dflag-2) {
			g.stick++
		}
	}
}

@ 칼부림의 중계. 머릿수를 알리고, 던진 칼과 맞은 칼을 센다. 하나라도
맞았으면---죽음이다. 문장 조립이 잘게 쪼개진 것은 단수·복수와
명중 수의 조합을 원작의 출력 그대로 재현하기 위해서다.

@<위협하는 난쟁이들이 칼을 던진다@>=
if g.dtotal == 1 {
	fmt.Fprintf(g.out, "위협적인 작은 난쟁이가 하나")
} else {
	fmt.Fprintf(g.out, "위협적인 작은 난쟁이가 %d마리", g.dtotal)
}
fmt.Fprintf(g.out, " 너랑 같은 방에 있어!\n")
if g.attack == 0 {
	return
}
if g.dflag == 2 {
	g.dflag = 3
}
k := 0
if g.attack == 1 {
	fmt.Fprintf(g.out, "날카롭고 고약한 칼 하나가 너에게 날아와")
} else {
	k = 2
	fmt.Fprintf(g.out, " 그중 %d마리가 너에게 칼을 던져", g.attack)
}
fmt.Fprintf(g.out, "---")
if g.stick <= 1 {
	fmt.Fprintf(g.out, "%s!\n", attackMsg[k+g.stick])
} else {
	fmt.Fprintf(g.out, "그중 %d마리가 맞혔어!\n", g.stick)
}
if g.stick != 0 {
	g.oldoldloc = g.loc
	g.die()
}

@ @<난쟁이와 해적@>=
var attackMsg = [...]string{
	"빗나갔어", "맞았어",
	"아무도 못 맞혔어", "하나가 맞혔어",
}

@ 도끼 반격은 던지기 동사에서 넘어온다. 명중률은 2/3. 죽은 난쟁이는
느끼한 검은 연기가 되어 사라지고---첫 마리만 그 묘사를 듣는다,
두 번째부터는 너도 이골이 났을 테니---도끼는 도로 바닥에 떨어진다.

@<난쟁이와 해적@>=
func (g *Game) throwAxeAtDwarf() actResult {
	j := 1
	for ; j <= nd; j++ {
		if g.dloc[j] == g.loc {
			break
		}
	}
	if g.ran(3) < 2 {
		g.dloc[j] = limbo
		g.dseen[j] = false
		g.dkill++
		if g.dkill == 1 {
			fmt.Fprintf(g.out, "작은 난쟁이를 죽였어. 사체가 느끼한 검은 연기 구름 속으로\n"+
				"사라져.\n")
		} else {
			fmt.Fprintf(g.out, "작은 난쟁이를 죽였어.\n")
		}
	} else {
		fmt.Fprintf(g.out, "작은 난쟁이를 공격하지만, 녀석이 잽싸게 피해.\n")
	}
	g.drop(AXE, g.loc)
	return g.stayPut()
}

@ 이제 해적. 해적의 관심사는 오직 보물이다. 다만 플로버 방과 어둠의
방의 피라미드는 줍기가 너무 쉬워 훔칠 가치도 없다고 여긴다---해적
나름의 직업 윤리다. 해적을 봤다는 증거는 미로에 놓이는 쪽지(|MESSAGE|)
인데, 그것이 아직 림보에 있으면 못 본 것이다.

@<난쟁이와 해적@>=
func (g *Game) tooEasy(i object) bool {
	return i == PYRAMID && (g.loc == proom || g.loc == droom)
}

func (g *Game) pirateNotSpotted() bool { return g.place[MESSAGE] == limbo }

@ 해적의 추적. 내가 훔칠 만한 보물을 들고 있으면($k<0$) 노획하고,
모든 보물을 이미 본 참인데 이 방에 보물이 없고 램프도 켜져 있다면
---즉 상자를 찾을 때가 됐다면---모습을 드러내 준다. 그도 아니면
20\%의 확률로 등 뒤에서 바스락거리기만 한다. 동굴에서 가장 등골이
서늘한 소리다.

@<난쟁이와 해적@>=
func (g *Game) pirateTrack() {
	if g.loc == maxPirateLoc || g.prop[CHEST] >= 0 {
		return
	}
	k := 0
	for i := minTreasure; i <= maxObj; i++ {
		if !g.tooEasy(i) && g.toting(i) {
			k = -1
			break
		}
		if g.here(i) {
			k = 1
		}
	}
	switch {
	case k < 0:
		@<해적이 노획물을 챙긴다@>@;
	case g.tally == g.lostTreasures+1 && k == 0 && g.pirateNotSpotted() &&
		g.prop[LAMP] != 0 && g.here(LAMP):
		@<해적이 모습을 드러낸다@>@;
	case g.odloc[0] != g.dloc[0] && g.pct(20):
		fmt.Fprintf(g.out, "네 뒤 어둠 속에서 희미하게 바스락거리는 소리가 들려.\n")
	}
}

@ 노획. 이 방에서 챙길 수 있는 보물을 모두---바닥에 놓인 것까지
집어서---상자 은닉처로 나르고, 해적 자신도 그리로 사라진다.
첫 강도라면 상자와 쪽지도 이때 배치된다.

@<해적이 노획물을 챙긴다@>=
fmt.Fprintf(g.out, "네 뒤 그림자 속에서 수염 난 해적이 튀어나와! \"하르, 하르,\" 그가\n"+
	"낄낄대며 말해. \"이 노획물은 내가 몽땅 챙겨서 미로 깊숙이 내 상자에\n"+
	"숨겨 두지!\"  그러더니 네 보물을 낚아채 어둠 속으로 사라져.\n")
for i := minTreasure; i <= maxObj; i++ {
	if g.tooEasy(i) {
		continue
	}
	if objBase[i] == NOTHING && g.place[i] == g.loc {
		g.carry(i)
	}
	if g.toting(i) {
		g.drop(i, chestLoc)
	}
}
if g.pirateNotSpotted() {
	g.moveChest()
}
g.dloc[0], g.odloc[0] = chestLoc, chestLoc
g.dseen[0] = false

@ @<난쟁이와 해적@>=
func (g *Game) moveChest() {
	g.move(CHEST, chestLoc)
	g.move(MESSAGE, messageLoc)
}

@ 목격. 훔칠 것이 없는 해적은 머쓱하게 상자를 숨기러 달아나며, 그
덕에 상자가 세상에 나온다. 보물을 다 보고도 점수가 안 차는 이유를
궁금해할 무렵 일어나도록 조건이 짜여 있다.

@<해적이 모습을 드러낸다@>=
fmt.Fprintf(g.out, "네 뒤 어둠 속에서 희미하게 바스락거리는 소리가 들려. 그쪽으로\n"+
	"고개를 돌리자, 네 램프 불빛이 수염 난 해적을 비춰. 녀석은 큰 상자를\n"+
	"들고 있어. \"이런 제기랄!\" 그가 외쳐. \"들켜 버렸군! 얼른 미로로\n"+
	"가서 내 상자를 숨겨야겠어!\"  그러더니 어둠 속으로 사라져.\n")
g.moveChest()
g.dloc[0], g.odloc[0] = chestLoc, chestLoc
g.dseen[0] = false

@ 큰 순환의 첫걸음이 참조하는 조각이 마지막으로 남았다. 이동하려는
곳(|newloc|)이 정해진 뒤 실제로 옮기기 전에, 방해꾼이 있는지 본다.
폐쇄 중에 밖으로 나가려 하면 녹음된 목소리가 막고, 칼 든 난쟁이가
길목(내가 가려는 곳이 녀석이 방금 떠나온 곳)에 버티고 있으면 못
지나간다.

@<폐쇄나 난쟁이가 이동을 막으면 제자리로 돌린다@>=
if g.closing() && g.newloc < minInCave && g.newloc != limbo {
	g.panicClosing()
	g.newloc = g.loc
} else if g.newloc != g.loc && g.loc <= maxPirateLoc {
	for j := 1; j <= nd; j++ {
		if g.odloc[j] == g.newloc && g.dseen[j] {
			fmt.Fprintf(g.out, "큰 칼을 든 작은 난쟁이가 네 길을 막아.\n")
			g.newloc = g.loc
			break
		}
	}
}

@ 난쟁이의 시험. 깊은 동굴을 돌아다니면 |dflag|가 오르는지, 같은 방의
난쟁이를 알아채는지 본다.

@(advent_test.go@>=
func TestDwarfActivation(t *testing.T) {
	g := &Game{in: bufio.NewReader(strings.NewReader("")), out: io.Discard, rx: 42}
	g.buildVocabulary()
	g.loadObjects()
	g.dloc = dwarfStart
	g.knifeLoc = -1
	g.loc = hmk // Hall of Mountain King(깊은 동굴)
	for i := 0; i < 300 && g.dflag < 2; i++ {
		g.moveDwarves()
	}
	if g.dflag < 1 {
		t.Fatal("동굴 깊은 곳에서도 dflag가 0 (난쟁이 비활성)")
	}
	t.Logf("dflag=%d 도달", g.dflag)
	// 같은 방의 난쟁이를 dwarf()가 알아채는지
	g.dflag = 2
	g.dloc[1] = hmk
	if !g.dwarf() {
		t.Error("같은 방의 난쟁이를 dwarf()가 못 알아챔")
	}
}
