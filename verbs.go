package main

import "fmt"

// 싸움·먹이기·열고닫기·읽기·말하기 동사. (advent.w "The other actions")

// ---- KILL (attack) ----

func (g *Game) doKill() actResult {
	if g.obj == NOTHING {
		if r, redirect := g.uniqueAttack(); redirect {
			return r
		}
	}
	switch g.obj {
	case NOTHING:
		return g.rep("There is nothing here to attack.")
	case BIRD:
		return g.dispatchBird()
	case DRAGON:
		if g.prop[DRAGON] == 0 {
			return g.funStuffDragon()
		}
		return g.rep("For crying out loud, the poor thing is already dead!")
	case CLAM, OYSTER:
		return g.rep("The shell is very strong and impervious to attack.")
	case SNAKE:
		return g.rep("Attacking the snake both doesn't work and is very dangerous.")
	case DWARF:
		if g.closed {
			return g.dwarvesUpset()
		}
		return g.rep("With what?  Your bare hands?")
	case TROLL:
		return g.rep("Trolls are close relatives with the rocks and have skin as tough as\n" +
			"a rhinoceros hide.  The troll fends off your blows effortlessly.")
	case BEAR:
		switch g.prop[BEAR] {
		case 0:
			return g.rep("With what?  Your bare hands?  Against HIS bear hands?")
		case 3:
			return g.rep("For crying out loud, the poor thing is already dead!")
		default:
			return g.rep("The bear is confused; he only wants to be your friend.")
		}
	default:
		return g.repDefault()
	}
}

// uniqueAttack: 공격할 대상이 하나로 정해지는지 본다. 둘 이상이면 되묻는다.
// (advent.w "See if there's a unique object to attack")
func (g *Game) uniqueAttack() (actResult, bool) {
	k := 0
	if g.dwarf() {
		k++
		g.obj = DWARF
	}
	if g.here(SNAKE) {
		k++
		g.obj = SNAKE
	}
	if g.isAtLoc(DRAGON) && g.prop[DRAGON] == 0 {
		k++
		g.obj = DRAGON
	}
	if g.isAtLoc(TROLL) {
		k++
		g.obj = TROLL
	}
	if g.here(BEAR) && g.prop[BEAR] == 0 {
		k++
		g.obj = BEAR
	}
	if k == 0 {
		if g.here(BIRD) && g.oldverb != TOSS {
			k++
			g.obj = BIRD
		}
		if g.here(CLAM) || g.here(OYSTER) {
			k++
			g.obj = CLAM
		}
	}
	if k > 1 {
		return g.getObjR(), true
	}
	return actResult{}, false
}

func (g *Game) dispatchBird() actResult {
	if g.closed {
		return g.rep("Oh, leave the poor unhappy bird alone.")
	}
	g.destroy(BIRD)
	g.prop[BIRD] = 0
	if g.place[SNAKE] == hmk {
		g.lostTreasures++
	}
	return g.rep("The little bird is now dead.  Its body disappears.")
}

// funStuffDragon: 맨손으로 용을 공격하겠다고 우기면 용이 죽는다. (advent.w "Fun stuff for dragon")
func (g *Game) funStuffDragon() actResult {
	fmt.Fprintf(g.out, "With what?  Your bare hands?\n")
	g.verb = ABSTAIN
	g.obj = NOTHING
	if !g.listen() {
		g.quitting = true
		return aDone()
	}
	if !(streq(g.word1, "yes") || streq(g.word1, "y")) {
		// TODO(11단계): 정확히는 이 입력을 명령으로 재처리(goto pre_parse).
		return aDone()
	}
	fmt.Fprintf(g.out, "%s\n", objNote[objOffset[DRAGON]+1])
	g.prop[DRAGON] = 2 // 죽음
	g.prop[RUG] = 0
	objBase[RUG] = NOTHING // 이제 쓸 수 있는 보물
	objBase[DRAGON_] = DRAGON_
	g.destroy(DRAGON_)
	objBase[RUG_] = RUG_
	g.destroy(RUG_)
	for t := object(1); t <= maxObj; t++ {
		if g.place[t] == scan1 || g.place[t] == scan3 {
			g.move(t, scan2)
		}
	}
	g.loc = scan2
	return g.stayPut()
}

// ---- FEED ----

func (g *Game) doFeed() actResult {
	switch g.obj {
	case BIRD:
		return g.rep("It's not hungry (it's merely pinin' for the fjords).  Besides, you\n" +
			"have no bird seed.")
	case TROLL:
		return g.rep("Gluttony is not one of the troll's vices.  Avarice, however, is.")
	case DRAGON:
		if g.prop[DRAGON] != 0 {
			return g.rep(g.defaultMsg[EAT])
		}
	case SNAKE:
		if !g.closed && g.here(BIRD) {
			g.destroy(BIRD)
			g.prop[BIRD] = 0
			g.lostTreasures++
			return g.rep("The snake has now devoured your bird.")
		}
	case BEAR:
		if !g.here(FOOD) {
			if g.prop[BEAR] == 0 {
				break
			}
			if g.prop[BEAR] == 3 {
				g.verb = EAT
			}
			return g.repDefault()
		}
		g.destroy(FOOD)
		g.prop[BEAR] = 1
		g.prop[AXE] = 0
		objBase[AXE] = NOTHING // 도끼를 다시 들 수 있다
		return g.rep("The bear eagerly wolfs down your food, after which he seems to calm\n" +
			"down considerably and even becomes rather friendly.")
	case DWARF:
		if !g.here(FOOD) {
			return g.repDefault()
		}
		g.dflag++
		return g.rep("You fool, dwarves eat only coal!  Now you've made him REALLY mad!")
	default:
		return g.rep(g.defaultMsg[CALM])
	}
	return g.rep("There's nothing here it wants to eat (except perhaps you).")
}

// ---- OPEN / CLOSE ----

func (g *Game) doOpenClose() actResult {
	switch g.obj {
	case OYSTER, CLAM:
		return g.openClam()
	case GRATE, CHAIN:
		if !g.here(KEYS) {
			return g.rep("You have no keys!")
		}
		return g.openGrateChain()
	case KEYS:
		return g.rep("You can't lock or unlock the keys.")
	case CAGE:
		return g.rep("It has no lock.")
	case DOOR:
		if g.prop[DOOR] != 0 {
			return g.defTo(RELAX)
		}
		return g.rep("The door is extremely rusty and refuses to open.")
	default:
		return g.repDefault()
	}
}

func (g *Game) openGrateChain() actResult {
	if g.obj == CHAIN {
		return g.openCloseChain()
	}
	if g.closing() {
		g.panicClosing()
		return aDone()
	}
	k := g.prop[GRATE]
	if g.verb == OPEN {
		g.prop[GRATE] = 1
	} else {
		g.prop[GRATE] = 0
	}
	switch k + 2*g.prop[GRATE] {
	case 0:
		return g.rep("It was already locked.")
	case 1:
		return g.rep("The grate is now locked.")
	case 2:
		return g.rep("The grate is now unlocked.")
	default: // 3
		return g.rep("It was already unlocked.")
	}
}

func (g *Game) openCloseChain() actResult {
	if g.verb == OPEN {
		return g.openChain()
	}
	if g.loc != barr {
		return g.rep("There is nothing here to which the chain can be locked.")
	}
	if g.prop[CHAIN] != 0 {
		return g.rep("It was already locked.")
	}
	g.prop[CHAIN] = 2
	objBase[CHAIN] = CHAIN
	if g.toting(CHAIN) {
		g.drop(CHAIN, g.loc)
	}
	return g.rep("The chain is now locked.")
}

func (g *Game) openChain() actResult {
	if g.prop[CHAIN] == 0 {
		return g.rep("It was already unlocked.")
	}
	if g.prop[BEAR] == 0 {
		return g.rep("There is no way to get past the bear to unlock the chain, which is\n" +
			"probably just as well.")
	}
	g.prop[CHAIN] = 0
	objBase[CHAIN] = NOTHING // 사슬이 풀렸다
	if g.prop[BEAR] == 3 {
		objBase[BEAR] = BEAR
	} else {
		g.prop[BEAR] = 2
		objBase[BEAR] = NOTHING
	}
	return g.rep("The chain is now unlocked.")
}

func (g *Game) openClam() actResult {
	name := "oyster"
	if g.obj == CLAM {
		name = "clam"
	}
	if g.verb == CLOSE {
		return g.rep("What?")
	}
	if !g.toting(TRIDENT) {
		fmt.Fprintf(g.out, "You don't have anything strong enough to open the %s", name)
		return g.rep(".")
	}
	if g.toting(g.obj) {
		fmt.Fprintf(g.out, "I advise you to put down the %s before opening it.  ", name)
		if g.obj == CLAM {
			return g.rep(">STRAIN!<")
		}
		return g.rep(">WRENCH!<")
	}
	if g.obj == CLAM {
		g.destroy(CLAM)
		g.drop(OYSTER, g.loc)
		g.drop(PEARL, sac)
		return g.rep("A glistening pearl falls out of the clam and rolls away.  Goodness,\n" +
			"this must really be an oyster.  (I never was very good at identifying\n" +
			"bivalves.)  Whatever it is, it has now snapped shut again.")
	}
	return g.rep("The oyster creaks open, revealing nothing but oyster inside.\n" +
		"It promptly snaps shut again.")
}

// ---- READ ----

func (g *Game) doRead() actResult {
	if g.dark() {
		return g.cantSeeIt()
	}
	switch g.obj {
	case MAG:
		return g.rep("I'm afraid the magazine is written in dwarvish.")
	case TABLET:
		return g.rep("\"CONGRATULATIONS ON BRINGING LIGHT INTO THE DARK-ROOM!\"")
	case MESSAGE:
		return g.rep("\"This is not the maze where the pirate hides his treasure chest.\"")
	case OYSTER:
		if g.hinted[1] {
			if g.toting(OYSTER) {
				return g.rep("It says the same thing it did before.")
			}
		} else if g.closed && g.toting(OYSTER) {
			g.offer(1)
			return aDone()
		}
		return g.repDefault()
	default:
		return g.repDefault()
	}
}

// cantSeeIt: 보이지 않는 객체를 가리켰을 때. (advent.w cant_see_it)
func (g *Game) cantSeeIt() actResult {
	if (g.verb == FIND || g.verb == INVENTORY) && g.word2 == "" {
		return aChange(g.verb)
	}
	fmt.Fprintf(g.out, "I see no %s here.\n", g.word1)
	return aDone()
}

// ---- SAY ----

func (g *Game) doSay() actResult {
	if g.word2 != "" {
		g.word1 = g.word2
	}
	if e, found := g.lookup(g.word1); found {
		magic := (e.typ == actionType && action(e.meaning) == FEEFIE) ||
			(e.typ == motionType && (motion(e.meaning) == XYZZY ||
				motion(e.meaning) == PLUGH || motion(e.meaning) == PLOVER))
		if magic {
			g.word2 = ""
			g.obj = NOTHING
			if e.typ == motionType {
				return g.tryMotionR(motion(e.meaning))
			}
			g.verb = action(e.meaning) // FEEFIE
			return aChange(g.verb)
		}
	}
	return g.rep(fmt.Sprintf("Okay, \"%s\".", g.word1))
}
