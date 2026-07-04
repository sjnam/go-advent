@* 죽음과 부활.
끈질긴 모험가만이 죽는 법도 배운다. 구덩이에 떨어지거나, 난쟁이의
칼에 맞거나, 곰과 함께 다리에서 추락하거나. 다행히 이 동굴의 관리인은
주황색 연기를 피우는 부활 마법을 부린다 --- 다만 재고가 세 번뿐이고,
회를 거듭할수록 말투가 험해진다. 짝수 인덱스가 부활 권유, 홀수가
부활 장면이다.

@<죽음과 부활@>=
const maxDeaths = 3 // 허용되는 죽음 횟수

var deathWishes = [2 * maxDeaths]string{
	"이런, 너 그만 죽어 버린 것 같네. 내가 도와줄 수\n" +
		"있을지도 모르겠지만, 사실 이런 건 해본 적이 없어. 널\n" +
		"환생시켜 볼까?",
	"좋아. 근데 뭔가 잘못돼도 날 탓하진 마......\n" +
		"                 --- 펑!! ---\n" +
		"넌 주황색 연기 구름에 휩싸여. 콜록대고 헐떡이며\n" +
		"연기 속에서 빠져나와 보니....",
	"이 칠칠맞은 녀석, 또 저질렀구나! 내가 이걸 얼마나 오래\n" +
		"해줄 수 있을지 모르겠어. 널 또 환생시켜 볼까?",
	"자, 내 부활 키트를 어디 뒀더라?.... >펑!<\n" +
		"모든 것이 짙은 주황색 연기 구름 속으로 사라져.",
	"이번엔 진짜 사고 쳤어! 주황색 연기가 다 떨어졌다고! 설마\n" +
		"주황색 연기도 없이 제대로 된 환생을 바라는 건 아니겠지?",
	"그래, 그렇게 잘났으면 직접 해 봐! 난 갈래!",
}

@ 죽음의 장부는 두 칸이면 된다.

@<게임의 상태@>=
deathCount int  // 몇 번 죽었나
dying      bool // 이번에 죽었는가 (죽음 처리 대기)

@ 죽음의 처리. 폐쇄 시각이 코앞이면 부활도 생략이다. 부활하면 들고
있던 것을 죽은 자리 근처(|oldoldloc|)에 떨군다 --- 거꾸로 도는 것은
새를 새장보다 먼저 떨궈 새가 새장 밖으로 나오게 하려는 원작의 세심함
이다. 램프만은 특별히 길 끝에 놓아 준다. 시체 회수 여행을 최소한
불 없이 시작하지는 않도록. 너 자신은 건물 안에서 깨어난다.

@<죽음과 부활@>=
func (g *Game) handleDeath() bool {
	g.dying = false
	g.deathCount++
	if g.closing() {
		fmt.Fprintf(g.out, "넌 죽은 것 같아. 뭐, 어차피 폐쇄 시각이 코앞이니\n"+
			"그냥 오늘은 여기까지 하자.\n")
		return false
	}
	if !g.yes(deathWishes[2*g.deathCount-2], deathWishes[2*g.deathCount-1], g.ok()) ||
		g.deathCount == maxDeaths {
		return false
	}
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

@ 폐쇄 후의 보관실에서 난쟁이를 건드리면 --- 깨우든, 치든, 거울을
깨서 소음을 내든 --- 떼로 일어나 칼을 던진다. 전부 명중이다. 최후의
퍼즐 방에서는 얌전히 폭파 지점이나 고민하라는 뜻이다.

@<죽음과 부활@>=
func (g *Game) dwarvesUpset() actResult {
	fmt.Fprintf(g.out, "그 소동에 난쟁이들이 깨어났어. 이제 위협적인 작은 난쟁이\n"+
		"여럿이 너랑 같은 방에 있어! 대부분이 너에게 칼을\n"+
		"던져! 전부 다 맞혔어!\n")
	g.die()
	return aDone()
}
