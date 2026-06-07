package main

import "fmt"

// 점수 계산과 마무리. (advent.w "Scoring")

const highestClass = 8 // 최고 등급 인덱스 (advent.w:4392)

// 등급 경계 점수와 그에 대응하는 칭호. (advent.w:4406)
var classScore = [...]int{35, 100, 130, 200, 250, 300, 330, 349, 9999}

var classMessage = [...]string{
	"You are obviously a rank amateur.  Better luck next time.",
	"Your score qualifies you as a novice class adventurer.",
	"You have achieved the rating \"Experienced Adventurer\".",
	"You may now consider yourself a \"Seasoned Adventurer\".",
	"You have reached \"Junior Master\" status.",
	"Your score puts you in Master Adventurer Class C.",
	"Your score puts you in Master Adventurer Class B.",
	"Your score puts you in Master Adventurer Class A.",
	"All of Adventuredom gives tribute to you, Adventure Grandmaster!",
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
	fmt.Fprintf(g.out, "You scored %d point%s out of a possible %d, using %d turn%s.\n",
		k, plural(k), maxScore, g.turns, plural(g.turns))
	j := 0
	for classScore[j] < k {
		j++
	}
	fmt.Fprintf(g.out, "%s\nTo achieve the next higher rating", classMessage[j])
	if j < highestClass {
		need := classScore[j] + 1 - k
		fmt.Fprintf(g.out, ", you need %d more point%s.\n", need, plural(need))
	} else {
		fmt.Fprintf(g.out, " would be a neat trick!\nCongratulations!!\n")
	}
}

// plural은 1이면 "", 아니면 "s"를 돌려준다(원본의 삼항식 처리).
func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
