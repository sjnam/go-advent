@* 점수.
채점 기준은 이렇다. 참가상 2점에, 동굴 깊숙이 들어갔으면 25점.
보물은 보기만 해도 2점, 깨지지 않은 채 건물에 모셔 오면 등급대로
더 준다 --- 흔한 보물 10점, 해적의 상자 12점, 그 뒤의 귀물들은 14점.
쓰지 않은 목숨이 하나에 10점, 스스로 그만두지 않았으면 4점,
Witt's End에 잡지를 두고 왔으면 1점(원작의 주석은 이 점수의 근거를
``순전히 심술''이라고 밝힌다), 폐쇄 단계 도달 25점, 그리고 마지막
퍼즐의 보너스. 힌트 값은 여기서 제한다. 만점 350점.

@<점수@>=
const maxScore = 350

func (g *Game) score() int {
	s := 2
	if g.dflag != 0 {
		s += 25
	}
	for k := minTreasure; k <= maxObj; k++ {
		if g.prop[k] >= 0 {
			s += 2
			if g.place[k] == house && g.prop[k] == 0 {
				switch {
				case k < CHEST:
					s += 10
				case k == CHEST:
					s += 12
				default:
					s += 14
				}
			}
		}
	}
	s += 10 * (maxDeaths - g.deathCount)
	if !g.gaveUp {
		s += 4
	}
	if g.place[MAG] == witt {
		s++ // Witt's End 방문의 증거
	}
	if g.closing() {
		s += 25
	}
	s += g.bonus
	for j := 0; j < nHints; j++ {
		if g.hinted[j] {
			s -= hintCost[j]
		}
	}
	return s
}

@ 채점표 밖의 유일한 상태가 그 마지막 보너스다.

@<게임의 상태@>=
bonus int // 마지막 퍼즐에서 얻는 추가 점수

@ 점수는 아홉 등급으로 번역된다. 풋내기부터 어드벤처 그랜드마스터까지.
경계값 배열의 마지막 9999는 ``이보다 높은 등급은 없다''는 뜻의
울타리다.

@<점수@>=
const highestClass = 8

var classScore = [...]int{35, 100, 130, 200, 250, 300, 330, 349, 9999}

var classMessage = [...]string{
	"넌 누가 봐도 풋내기야.  다음엔 운이 따르길.",
	"이 점수면 초보 모험가라고 할 만해.",
	"\"숙련된 모험가\" 등급을 얻었어.",
	"이제 스스로를 \"노련한 모험가\"라 여겨도 돼.",
	"\"주니어 마스터\" 경지에 올랐어.",
	"이 점수면 마스터 모험가 C등급이야.",
	"이 점수면 마스터 모험가 B등급이야.",
	"이 점수면 마스터 모험가 A등급이야.",
	"온 모험계가 그대에게 경의를 표한다, 어드벤처 그랜드마스터여!",
}

@ 작별 인사. 점수와 등급을 알리고, 다음 등급까지 몇 점이 모자란지도
일러 준다 --- 다시 오라는 뜻이다.

@<점수@>=
func (g *Game) printScore() {
	k := g.score()
	fmt.Fprintf(g.out, "%d점 만점에 %d점을 얻었고, %d번 움직였어.\n",
		maxScore, k, g.turns)
	j := 0
	for classScore[j] < k {
		j++
	}
	fmt.Fprintf(g.out, "%s\n", classMessage[j])
	if j < highestClass {
		need := classScore[j] + 1 - k
		fmt.Fprintf(g.out, "다음 등급에 오르려면 %d점이 더 필요해.\n", need)
	} else {
		fmt.Fprintf(g.out, "더 높은 등급에 오르는 건 신기에 가까운 일이야!\n축하해!!\n")
	}
}
