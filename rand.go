package main

// 원본의 선형 합동 난수 생성기 (advent.w "Random numbers").
// rx는 Game에 보관되며, 시드를 고정하면 게임 전체가 결정론적이다.

// ran은 0 이상 rng 미만의 균등 정수를 돌려준다.
func (g *Game) ran(rng int) int {
	g.rx = (1021 * g.rx) & 0xfffff // 1021을 곱하고 2^20으로 나머지
	return (rng * g.rx) >> 20
}

// pct는 r 퍼센트의 확률로 true를 돌려준다. (원본 매크로 pct)
func (g *Game) pct(r int) bool { return g.ran(100) < r }
