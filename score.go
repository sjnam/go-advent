package main

import "fmt"

// 점수 계산과 마무리. (advent.w "Scoring")

const highestClass = 8 // 최고 등급 인덱스 (advent.w:4392)

// 등급 경계 점수와 그에 대응하는 칭호. (advent.w:4406)
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

// score는 현재 점수를 계산한다. 보물은 깨지지 않고 건물에 두었을 때만
// 만점을 받지만, 보기만 해도 2점을 준다. (advent.w score)
func (g *Game) score() int {
	s := 2
	if g.dflag != 0 {
		s += 25 // 동굴 깊숙이 들어감
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
		s++ // Witt's End을 방문한 증거
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

// printScore는 점수와 등급을 출력하며 작별 인사를 한다. (advent.w "Print the score and say adieu")
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
