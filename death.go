package main

import "fmt"

// 죽음과 부활. 끈질긴 모험가는 최대 세 번까지 되살아날 수 있다.
// (advent.w "Death and resurrection")

const maxDeaths = 3 // 허용되는 죽음 횟수 (advent.w:4168)

// 죽음·부활 때 하는 말. 짝수 인덱스는 부활 권유, 홀수는 부활 시 메시지. (advent.w:4183)
var deathWishes = [2 * maxDeaths]string{
	"Oh dear, you seem to have gotten yourself killed.  I might be able to\n" +
		"help you out, but I've never really done this before.  Do you want me\n" +
		"to try to reincarnate you?",
	"All right.  But don't blame me if something goes wr......\n" +
		"                 --- POOF!! ---\n" +
		"You are engulfed in a cloud of orange smoke.  Coughing and gasping,\n" +
		"you emerge from the smoke and find....",
	"You clumsy oaf, you've done it again!  I don't know how long I can\n" +
		"keep this up.  Do you want me to try reincarnating you again?",
	"Okay, now where did I put my resurrection kit?....  >POOF!<\n" +
		"Everything disappears in a dense cloud of orange smoke.",
	"Now you've really done it!  I'm out of orange smoke!  You don't expect\n" +
		"me to do a decent reincarnation without any orange smoke, do you?",
	"Okay, if you're so smart, do it yourself!  I'm leaving!",
}

// handleDeath는 죽음을 처리한다. 부활하면 true, 게임을 끝내야 하면 false.
// 죽었을 때 newloc은 의미가 없고 oldloc은 당신을 죽인 곳이므로, 마지막으로
// 안전했던 oldoldloc을 본다. (advent.w "Deal with death and resurrection")
func (g *Game) handleDeath() bool {
	g.dying = false
	g.deathCount++
	if g.closing() {
		fmt.Fprintf(g.out, "It looks as though you're dead.  Well, seeing as how it's so close\n"+
			"to closing time anyway, let's just call it a day.\n")
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
	fmt.Fprintf(g.out, "The resulting ruckus has awakened the dwarves.  There are now several\n"+
		"threatening little dwarves in the room with you!  Most of them throw\n"+
		"knives at you!  All of them get you!\n")
	g.die()
	return aDone()
}
