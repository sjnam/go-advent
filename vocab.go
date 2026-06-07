package main

import (
	"fmt"
	"strings"
)

// 원본은 단어를 직접 만든 해시 테이블에 넣었다(advent.w "The vocabulary").
// 해시 값 자체는 게임 출력에 영향이 없으므로, 여기서는 Go map으로 단순화한다.
// 의미가 같으면(같은 단어 → 같은 종류·뜻) 동작이 동일하다.

type vocabEntry struct {
	typ     wordtype // 단어의 종류
	meaning int      // motion/object/action 값, 또는 message 인덱스
}

// lookup은 단어를 찾는다. 원본처럼 앞 5글자만 본다(advent.w:195).
func (g *Game) lookup(w string) (vocabEntry, bool) {
	if len(w) > 5 {
		w = w[:5]
	}
	e, ok := g.words[w]
	return e, ok
}

// listen은 한두 단어짜리 명령을 읽어 word1, word2에 담는다.
// (advent.w "Low-level input"의 listen). EOF면 false.
func (g *Game) listen() bool {
	for {
		fmt.Fprintf(g.out, "* ")
		line, more := g.readLine()
		if !more {
			return false
		}
		fields := strings.Fields(line)
		switch {
		case len(fields) == 0:
			fmt.Fprintf(g.out, " Tell me to do something.\n")
		case len(fields) > 2:
			fmt.Fprintf(g.out, " Please stick to 1- and 2-word commands.\n")
		default:
			g.word1 = strings.ToLower(fields[0])
			g.word2 = ""
			if len(fields) == 2 {
				g.word2 = strings.ToLower(fields[1])
			}
			return true
		}
	}
}

// buildVocabulary는 우리가 아는 모든 단어와 기본 메시지를 채운다.
// (advent.w의 세 "Build the vocabulary" 섹션)
func (g *Game) buildVocabulary() {
	g.words = make(map[string]vocabEntry)
	mw := func(w string, m motion) { g.words[w] = vocabEntry{motionType, int(m)} }
	ow := func(w string, o object) { g.words[w] = vocabEntry{objectType, int(o)} }
	aw := func(w string, a action) { g.words[w] = vocabEntry{actionType, int(a)} }

	// ---- 이동 단어 ----
	mw("north", N)
	mw("n", N)
	mw("south", S)
	mw("s", S)
	mw("east", E)
	mw("e", E)
	mw("west", W)
	mw("w", W)
	mw("ne", NE)
	mw("se", SE)
	mw("nw", NW)
	mw("sw", SW)
	mw("upwar", U)
	mw("up", U)
	mw("u", U)
	mw("above", U)
	mw("ascen", U)
	mw("downw", D)
	mw("down", D)
	mw("d", D)
	mw("desce", D)
	mw("left", L)
	mw("right", R)
	mw("inwar", IN)
	mw("insid", IN)
	mw("in", IN)
	mw("out", OUT)
	mw("outsi", OUT)
	mw("exit", OUT)
	mw("leave", OUT)
	mw("forwa", FORWARD)
	mw("conti", FORWARD)
	mw("onwar", FORWARD)
	mw("back", BACK)
	mw("retur", BACK)
	mw("retre", BACK)
	mw("over", OVER)
	mw("acros", ACROSS)
	mw("upstr", UPSTREAM)
	mw("downs", DOWNSTREAM)
	mw("enter", ENTER)
	mw("crawl", CRAWL)
	mw("jump", JUMP)
	mw("climb", CLIMB)
	mw("look", LOOK)
	mw("exami", LOOK)
	mw("touch", LOOK)
	mw("descr", LOOK)
	mw("cross", CROSS)
	mw("road", ROAD)
	mw("hill", ROAD)
	mw("fores", WOODS)
	mw("valle", VALLEY)
	mw("build", HOUSE)
	mw("house", HOUSE)
	mw("gully", GULLY)
	mw("strea", STREAM)
	mw("depre", DEPRESSION)
	mw("entra", ENTRANCE)
	mw("cave", CAVE)
	mw("rock", ROCK)
	mw("slab", SLAB)
	mw("slabr", SLAB)
	mw("bed", BED)
	mw("passa", PASSAGE)
	mw("tunne", PASSAGE)
	mw("caver", CAVERN)
	mw("canyo", CANYON)
	mw("awkwa", AWKWARD)
	mw("secre", SECRET)
	mw("bedqu", BEDQUILT)
	mw("reser", RESERVOIR)
	mw("giant", GIANT)
	mw("orien", ORIENTAL)
	mw("shell", SHELL)
	mw("barre", BARREN)
	mw("broke", BROKEN)
	mw("debri", DEBRIS)
	mw("view", VIEW)
	mw("fork", FORK)
	mw("pit", PIT)
	mw("slit", SLIT)
	mw("crack", CRACK)
	mw("dome", DOME)
	mw("hole", HOLE)
	mw("wall", WALL)
	mw("hall", HALL)
	mw("room", ROOM)
	mw("floor", FLOOR)
	mw("stair", STAIRS)
	mw("steps", STEPS)
	mw("cobbl", COBBLES)
	mw("surfa", SURFACE)
	mw("dark", DARK)
	mw("low", LOW)
	mw("outdo", OUTDOORS)
	mw("y2", Y2)
	mw("xyzzy", XYZZY)
	mw("plugh", PLUGH)
	mw("plove", PLOVER)
	mw("main", OFFICE)
	mw("offic", OFFICE)
	mw("null", NOWHERE)
	mw("nowhe", NOWHERE)

	// ---- 사물 단어 ----
	ow("key", KEYS)
	ow("keys", KEYS)
	ow("lamp", LAMP)
	ow("lante", LAMP)
	ow("headl", LAMP)
	ow("grate", GRATE)
	ow("cage", CAGE)
	ow("rod", ROD)
	ow("bird", BIRD)
	ow("door", DOOR)
	ow("pillo", PILLOW)
	ow("velve", PILLOW)
	ow("snake", SNAKE)
	ow("fissu", CRYSTAL)
	ow("table", TABLET)
	ow("clam", CLAM)
	ow("oyste", OYSTER)
	ow("magaz", MAG)
	ow("issue", MAG)
	ow("spelu", MAG)
	ow("\"spel", MAG)
	ow("dwarf", DWARF)
	ow("dwarv", DWARF)
	ow("knife", KNIFE)
	ow("knive", KNIFE)
	ow("food", FOOD)
	ow("ratio", FOOD)
	ow("bottl", BOTTLE)
	ow("jar", BOTTLE)
	ow("water", WATER)
	ow("h2o", WATER)
	ow("oil", OIL)
	ow("mirro", MIRROR)
	ow("plant", PLANT)
	ow("beans", PLANT)
	ow("stala", STALACTITE)
	ow("shado", SHADOW)
	ow("figur", SHADOW)
	ow("axe", AXE)
	ow("drawi", ART)
	ow("pirat", PIRATE)
	ow("drago", DRAGON)
	ow("chasm", BRIDGE)
	ow("troll", TROLL)
	ow("bear", BEAR)
	ow("messa", MESSAGE)
	ow("volca", GEYSER)
	ow("geyse", GEYSER)
	ow("vendi", PONY)
	ow("machi", PONY)
	ow("batte", BATTERIES)
	ow("moss", MOSS)
	ow("carpe", MOSS)
	ow("gold", GOLD)
	ow("nugge", GOLD)
	ow("diamo", DIAMONDS)
	ow("silve", SILVER)
	ow("bars", SILVER)
	ow("jewel", JEWELS)
	ow("coins", COINS)
	ow("chest", CHEST)
	ow("box", CHEST)
	ow("treas", CHEST)
	ow("eggs", EGGS)
	ow("egg", EGGS)
	ow("nest", EGGS)
	ow("tride", TRIDENT)
	ow("ming", VASE)
	ow("vase", VASE)
	ow("shard", VASE)
	ow("potte", VASE)
	ow("emera", EMERALD)
	ow("plati", PYRAMID)
	ow("pyram", PYRAMID)
	ow("pearl", PEARL)
	ow("persi", RUG)
	ow("rug", RUG)
	ow("spice", SPICES)
	ow("chain", CHAIN)

	// ---- 동작 단어 + 기본 메시지 ----
	aw("take", TAKE)
	aw("carry", TAKE)
	aw("keep", TAKE)
	aw("catch", TAKE)
	aw("captu", TAKE)
	aw("steal", TAKE)
	aw("get", TAKE)
	aw("tote", TAKE)
	g.defaultMsg[TAKE] = "You are already carrying it!"
	aw("drop", DROP)
	aw("relea", DROP)
	aw("free", DROP)
	aw("disca", DROP)
	aw("dump", DROP)
	g.defaultMsg[DROP] = "You aren't carrying it!"
	aw("open", OPEN)
	aw("unloc", OPEN)
	g.defaultMsg[OPEN] = "I don't know how to lock or unlock such a thing."
	aw("close", CLOSE)
	aw("lock", CLOSE)
	g.defaultMsg[CLOSE] = g.defaultMsg[OPEN]
	aw("light", ON)
	aw("on", ON)
	g.defaultMsg[ON] = "You have no source of light."
	aw("extin", OFF)
	aw("off", OFF)
	g.defaultMsg[OFF] = g.defaultMsg[ON]
	aw("wave", WAVE)
	aw("shake", WAVE)
	aw("swing", WAVE)
	g.defaultMsg[WAVE] = "Nothing happens."
	aw("calm", CALM)
	aw("placa", CALM)
	aw("tame", CALM)
	g.defaultMsg[CALM] = "I'm game.  Would you care to explain how?"
	aw("walk", GO)
	aw("run", GO)
	aw("trave", GO)
	aw("go", GO)
	aw("proce", GO)
	aw("explo", GO)
	aw("goto", GO)
	aw("follo", GO)
	aw("turn", GO)
	g.defaultMsg[GO] = "Where?"
	aw("nothi", RELAX)
	g.defaultMsg[RELAX] = "OK."
	aw("pour", POUR)
	g.defaultMsg[POUR] = g.defaultMsg[DROP]
	aw("eat", EAT)
	aw("devou", EAT)
	g.defaultMsg[EAT] = "Don't be ridiculous!"
	aw("drink", DRINK)
	g.defaultMsg[DRINK] = "You have taken a drink from the stream.  The water tastes strongly of\n" +
		"minerals, but is not unpleasant.  It is extremely cold."
	aw("rub", RUB)
	g.defaultMsg[RUB] = "Rubbing the electric lamp is not particularly rewarding.  Anyway,\n" +
		"nothing exciting happens."
	aw("throw", TOSS)
	aw("toss", TOSS)
	g.defaultMsg[TOSS] = g.defaultMsg[DROP]
	aw("wake", WAKE)
	aw("distu", WAKE)
	g.defaultMsg[WAKE] = g.defaultMsg[EAT]
	aw("feed", FEED)
	g.defaultMsg[FEED] = "There is nothing here to eat."
	aw("fill", FILL)
	g.defaultMsg[FILL] = "You can't fill that."
	aw("break", BREAK)
	aw("smash", BREAK)
	aw("shatt", BREAK)
	g.defaultMsg[BREAK] = "It is beyond your power to do that."
	aw("blast", BLAST)
	aw("deton", BLAST)
	aw("ignit", BLAST)
	aw("blowu", BLAST)
	g.defaultMsg[BLAST] = "Blasting requires dynamite."
	aw("attac", KILL)
	aw("kill", KILL)
	aw("fight", KILL)
	aw("hit", KILL)
	aw("strik", KILL)
	aw("slay", KILL)
	g.defaultMsg[KILL] = g.defaultMsg[EAT]
	aw("say", SAY)
	aw("chant", SAY)
	aw("sing", SAY)
	aw("utter", SAY)
	aw("mumbl", SAY)
	aw("read", READ)
	aw("perus", READ)
	g.defaultMsg[READ] = "I'm afraid I don't understand."
	aw("fee", FEEFIE)
	aw("fie", FEEFIE)
	aw("foe", FEEFIE)
	aw("foo", FEEFIE)
	aw("fum", FEEFIE)
	g.defaultMsg[FEEFIE] = "I don't know how."
	aw("brief", BRIEF)
	g.defaultMsg[BRIEF] = "On what?"
	aw("find", FIND)
	aw("where", FIND)
	g.defaultMsg[FIND] = "I can only tell you what you see as you move about and manipulate\n" +
		"things.  I cannot tell you where remote things are."
	aw("inven", INVENTORY)
	g.defaultMsg[INVENTORY] = g.defaultMsg[FIND]
	aw("score", SCORE)
	g.defaultMsg[SCORE] = "Eh?"
	aw("quit", QUIT)
	g.defaultMsg[QUIT] = g.defaultMsg[SCORE]

	// ---- 메시지 단어 (help, info 등) ----
	g.buildMessageWords()
}

// buildMessageWords는 고정 메시지를 출력하는 단어들을 등록한다.
// (advent.w의 message_type "Build the vocabulary" 섹션)
func (g *Game) buildMessageWords() {
	k := 0
	mw := func(w string) { g.words[w] = vocabEntry{messageType, k} } // 현재 메시지 슬롯을 가리킴
	nm := func(s string) { g.message[k] = s; k++ }                   // 메시지를 채우고 다음 슬롯으로

	mw("abra")
	mw("abrac")
	mw("opens")
	mw("sesam")
	mw("shaza")
	mw("hocus")
	mw("pocus")
	nm("Good try, but that is an old worn-out magic word.")
	mw("help")
	mw("?")
	nm("I know of places, actions, and things.  Most of my vocabulary\n" +
		"describes places and is used to move you there.  To move, try words\n" +
		"like forest, building, downstream, enter, east, west, north, south,\n" +
		"up, or down.  I know about a few special objects, like a black rod\n" +
		"hidden in the cave.  These objects can be manipulated using some of\n" +
		"the action words that I know.  Usually you will need to give both the\n" +
		"object and action words (in either order), but sometimes I can infer\n" +
		"the object from the verb alone.  Some objects also imply verbs; in\n" +
		"particular, \"inventory\" implies \"take inventory\", which causes me to\n" +
		"give you a list of what you're carrying.  The objects have side\n" +
		"effects; for instance, the rod scares the bird.  Usually people having\n" +
		"trouble moving just need to try a few more words.  Usually people\n" +
		"trying unsuccessfully to manipulate an object are attempting something\n" +
		"beyond their (or my!) capabilities and should try a completely\n" +
		"different tack.  To speed the game you can sometimes move long\n" +
		"distances with a single word.  For example, \"building\" usually gets\n" +
		"you to the building from anywhere above ground except when lost in the\n" +
		"forest.  Also, note that cave passages turn a lot, and that leaving a\n" +
		"room to the north does not guarantee entering the next from the south.\n" +
		"Good luck!")
	mw("tree")
	mw("trees")
	nm("The trees of the forest are large hardwood oak and maple, with an\n" +
		"occasional grove of pine or spruce.  There is quite a bit of under-\n" +
		"growth, largely birch and ash saplings plus nondescript bushes of\n" +
		"various sorts.  This time of year visibility is quite restricted by\n" +
		"all the leaves, but travel is quite easy if you detour around the\n" +
		"spruce and berry bushes.")
	mw("dig")
	mw("excav")
	nm("Digging without a shovel is quite impractical.  Even with a shovel\n" +
		"progress is unlikely.")
	mw("lost")
	nm("I'm as confused as you are.")
	nm("There is a loud explosion and you are suddenly splashed across the\n" +
		"walls of the room.")
	nm("There is a loud explosion and a twenty-foot hole appears in the far\n" +
		"wall, burying the snakes in the rubble.  A river of molten lava pours\n" +
		"in through the hole, destroying everything in its path, including you!")
	mw("mist")
	nm("Mist is a white vapor, usually water, seen from time to time in\n" +
		"caverns.  It can be found anywhere but is frequently a sign of a deep\n" +
		"pit leading down to water.")
	mw("fuck")
	nm("Watch it!")
	nm("There is a loud explosion, and a twenty-foot hole appears in the far\n" +
		"wall, burying the dwarves in the rubble.  You march through the hole\n" +
		"and find yourself in the main office, where a cheering band of\n" +
		"friendly elves carry the conquering adventurer off into the sunset.")
	mw("stop")
	nm("I don't know the word \"stop\".  Use \"quit\" if you want to give up.")
	mw("info")
	mw("infor")
	nm("If you want to end your adventure early, say \"quit\".  To get full\n" +
		"credit for a treasure, you must have left it safely in the building,\n" +
		"though you get partial credit just for locating it.  You lose points\n" +
		"for getting killed, or for quitting, though the former costs you more.\n" +
		"There are also points based on how much (if any) of the cave you've\n" +
		"managed to explore; in particular, there is a large bonus just for\n" +
		"getting in (to distinguish the beginners from the rest of the pack),\n" +
		"and there are other ways to determine whether you've been through some\n" +
		"of the more harrowing sections.  If you think you've found all the\n" +
		"treasures, just keep exploring for a while.  If nothing interesting\n" +
		"happens, you haven't found them all yet.  If something interesting\n" +
		"DOES happen, it means you're getting a bonus and have an opportunity\n" +
		"to garner many more points in the master's section.\n" +
		"I may occasionally offer hints if you seem to be having trouble.\n" +
		"If I do, I'll warn you in advance how much it will affect your score\n" +
		"to accept the hints.  Finally, to save paper, you may specify \"brief\",\n" +
		"which tells me never to repeat the full description of a place\n" +
		"unless you explicitly ask me to.")
	mw("swim")
	nm("I don't know how.")
}
