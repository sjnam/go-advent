package main

// 원본 C 프로그램의 enum들을 그대로 옮긴 것이다.
// 값의 순서는 절대 바꾸면 안 된다 — 이동표/객체표가 이 정수값으로 색인된다.

// ---- 단어의 종류 (advent.w "The vocabulary") ----

type wordtype int

const (
	noType wordtype = iota
	motionType
	objectType
	actionType
	messageType
)

// ---- 이동 동사 (advent.w "Cave connections" 직전) ----

type motion int

const (
	N motion = iota
	S
	E
	W
	NE
	SE
	NW
	SW
	U
	D
	L
	R
	IN
	OUT
	FORWARD
	BACK
	OVER
	ACROSS
	UPSTREAM
	DOWNSTREAM
	ENTER
	CRAWL
	JUMP
	CLIMB
	LOOK
	CROSS
	ROAD
	WOODS
	VALLEY
	HOUSE
	GULLY
	STREAM
	DEPRESSION
	ENTRANCE
	CAVE
	ROCK
	SLAB
	BED
	PASSAGE
	CAVERN
	CANYON
	AWKWARD
	SECRET
	BEDQUILT
	RESERVOIR
	GIANT
	ORIENTAL
	SHELL
	BARREN
	BROKEN
	DEBRIS
	VIEW
	FORK
	PIT
	SLIT
	CRACK
	DOME
	HOLE
	WALL
	HALL
	ROOM
	FLOOR
	STAIRS
	STEPS
	COBBLES
	SURFACE
	DARK
	LOW
	OUTDOORS
	Y2
	XYZZY
	PLUGH
	PLOVER
	OFFICE
	NOWHERE
)

// ---- 사물 (advent.w "Data structures for objects") ----
// 밑줄로 끝나는 이름(GRATE_ 등)은 같은 물체의 "두 번째 상태"를 가리킨다.

type object int

const (
	NOTHING object = iota
	KEYS
	LAMP
	GRATE
	GRATE_
	CAGE
	ROD
	ROD2
	TREADS
	TREADS_
	BIRD
	DOOR
	PILLOW
	SNAKE
	CRYSTAL
	CRYSTAL_
	TABLET
	CLAM
	OYSTER
	MAG
	DWARF
	KNIFE
	FOOD
	BOTTLE
	WATER
	OIL
	MIRROR
	MIRROR_
	PLANT
	PLANT2
	PLANT2_
	STALACTITE
	SHADOW
	SHADOW_
	AXE
	ART
	PIRATE
	DRAGON
	DRAGON_
	BRIDGE
	BRIDGE_
	TROLL
	TROLL_
	TROLL2
	TROLL2_
	BEAR
	MESSAGE
	GEYSER
	PONY
	BATTERIES
	MOSS
	GOLD
	DIAMONDS
	SILVER
	JEWELS
	COINS
	CHEST
	EGGS
	TRIDENT
	VASE
	EMERALD
	PYRAMID
	PEARL
	RUG
	RUG_
	SPICES
	CHAIN
)

// 보물과 사물 범위에 대한 파생 상수 (advent.w:346-348)
const (
	minTreasure = GOLD
	maxObj      = CHAIN
)

func isTreasure(t object) bool { return t >= minTreasure }

// ---- 동작 동사 (advent.w "Vocabulary", action enum) ----

type action int

const (
	ABSTAIN action = iota
	TAKE
	DROP
	OPEN
	CLOSE
	ON
	OFF
	WAVE
	CALM
	GO
	RELAX
	POUR
	EAT
	DRINK
	RUB
	TOSS
	WAKE
	FEED
	FILL
	BREAK
	BLAST
	KILL
	SAY
	READ
	FEEFIE
	BRIEF
	FIND
	INVENTORY
	SCORE
	QUIT
)

// ---- 장소 (advent.w "Cave data") ----

type location int

const (
	inhand location = iota - 1 // 들고 있는 물체의 위치 코드 (-1)
	limbo                      // 0: 아무 데도 없음
	road
	hill
	house
	valley
	forest
	woods
	slit
	outside
	inside
	cobbles
	debris
	awk
	bird
	spit
	emist
	nugget
	efiss
	wfiss
	wmist
	like1
	like2
	like3
	like4
	like5
	like6
	like7
	like8
	like9
	like10
	like11
	like12
	like13
	like14
	brink
	elong
	wlong
	diff0
	diff1
	diff2
	diff3
	diff4
	diff5
	diff6
	diff7
	diff8
	diff9
	diff10
	pony
	cross
	hmk
	west
	south
	ns
	y2
	jumble
	windoe
	dirty
	clean
	wet
	dusty
	complex
	shell
	arch
	ragged
	sac
	ante
	witt
	bedquilt
	cheese
	soft
	e2pit
	w2pit
	epit
	wpit
	narrow
	giant
	block
	immense
	falls
	steep
	abovep
	sjunc
	tite
	low
	crawl
	window
	oriental
	misty
	alcove
	proom
	droom
	slab
	abover
	mirror
	res
	scan1
	scan2
	scan3
	secret
	wide
	tight
	tall
	boulders
	scorr
	swside
	dead0
	dead1
	dead2
	dead3
	dead4
	dead5
	dead6
	dead7
	dead8
	dead9
	dead10
	dead11
	neside
	corr
	fork
	warm
	view
	chamber
	lime
	fbarr
	barr
	neend
	swend
	crack
	neck
	lose
	cant
	climb
	check
	snaked
	thru
	duck
	sewer
	upnout
	didit
	ppass
	pdrop
	troll
)

// 장소 범위에 대한 파생 상수 (advent.w:704-708)
const (
	minInCave    = inside
	minLowerLoc  = emist
	minForcedLoc = crack
	maxLoc       = didit
	maxSpec      = troll
)

// forced_move: 도착 즉시 다른 곳으로 보내지는 더미 장소인가 (advent.w:2046)
func forcedMove(l location) bool { return l >= minForcedLoc }

// 만점 (advent.w:4250)
const maxScore = 350
