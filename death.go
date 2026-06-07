package main

import "fmt"

// 죽음과 부활. 끈질긴 모험가는 최대 세 번까지 되살아날 수 있다.
// (advent.w "Death and resurrection")

const maxDeaths = 3 // 허용되는 죽음 횟수 (advent.w:4168)

// 죽음·부활 때 하는 말. 짝수 인덱스는 부활 권유, 홀수는 부활 시 메시지. (advent.w:4183)
var deathWishes = [2 * maxDeaths]string{
	"이런, 너 그만 죽어 버린 것 같네.  내가 도와줄 수\n" +
		"있을지도 모르겠지만, 사실 이런 건 해본 적이 없어.  널\n" +
		"환생시켜 볼까?",
	"좋아.  근데 뭔가 잘못돼도 날 탓하진 마......\n" +
		"                 --- 펑!! ---\n" +
		"넌 주황색 연기 구름에 휩싸여.  콜록대고 헐떡이며\n" +
		"연기 속에서 빠져나와 보니....",
	"이 칠칠맞은 녀석, 또 저질렀구나!  내가 이걸 얼마나 오래\n" +
		"해줄 수 있을지 모르겠어.  널 또 환생시켜 볼까?",
	"자, 내 부활 키트를 어디 뒀더라?....  >펑!<\n" +
		"모든 것이 짙은 주황색 연기 구름 속으로 사라져.",
	"이번엔 진짜 사고 쳤어!  주황색 연기가 다 떨어졌다고!  설마\n" +
		"주황색 연기도 없이 제대로 된 환생을 바라는 건 아니겠지?",
	"그래, 그렇게 잘났으면 직접 해 봐!  난 갈래!",
}

// handleDeath는 죽음을 처리한다. 부활하면 true, 게임을 끝내야 하면 false.
// 죽었을 때 newloc은 의미가 없고 oldloc은 당신을 죽인 곳이므로, 마지막으로
// 안전했던 oldoldloc을 본다. (advent.w "Deal with death and resurrection")
func (g *Game) handleDeath() bool {
	g.dying = false
	g.deathCount++
	if g.closing() {
		fmt.Fprintf(g.out, "넌 죽은 것 같아.  뭐, 어차피 폐쇄 시각이 코앞이니\n"+
			"그냥 오늘은 여기까지 하자.\n")
		return false
	}
	if !g.yes(deathWishes[2*g.deathCount-2], deathWishes[2*g.deathCount-1], g.ok()) ||
		g.deathCount == maxDeaths {
		return false
	}
	// 부활: 들고 있던 것들을 oldoldloc에 떨군다(새는 새장보다 먼저 떨구려고
	// 거꾸로 돈다). 램프는 건물 밖에 두고, 당신은 건물 안에 둔다.
	if g.toting(LAMP) {
		g.prop[LAMP] = 0
	}
	g.place[WATER] = limbo // drop하면 안 되므로 직접
	g.place[OIL] = limbo
	for j := maxObj; j > 0; j-- {
		if g.toting(j) {
			dest := g.oldoldloc
			if j == LAMP {
				dest = road
			}
			g.drop(j, dest)
		}
	}
	g.loc, g.oldloc = house, house
	return true
}

// dwarvesUpset: 폐쇄 후 난쟁이를 깨우면 떼로 칼을 던져 당신을 죽인다.
// (advent.w dwarves_upset)
func (g *Game) dwarvesUpset() actResult {
	fmt.Fprintf(g.out, "그 소동에 난쟁이들이 깨어났어.  이제 위협적인 작은 난쟁이\n"+
		"여럿이 너랑 같은 방에 있어!  대부분이 너에게 칼을\n"+
		"던져!  전부 다 맞혔어!\n")
	g.die()
	return aDone()
}
