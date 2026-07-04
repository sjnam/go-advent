% types.w --- gweave 조판 힌트 모음. (advent.w가 림보에서 @i로 끌어온다.)
% 이 web의 const 이름들을 @d로 등록해 CWEB의 매크로처럼 타자체로 조판한다.
% gtangle은 @d를 완전히 무시하므로 생성되는 Go 코드에는 아무 영향이 없다.
% (gweave가 const 선언을 자동 분류하려 하지만, `N motion = iota` 뒤에
% 이름만 잇는 연속 줄은 알아채지 못하므로 여기서 손수 일러 준다.
% @d는 한 줄에 하나씩만 쓴다 --- 이름 뒤에 오는 것은 값으로 보고 버린다.)

% 타입 이름들을 타입처럼(굵게) 조판하도록 등록한다.
@s Reader int
@s Writer int
@s Game int
@s wordtype int
@s vocabEntry int
@s motion int
@s object int
@s action int
@s location int
@s instruction int
@s inputResult int
@s objResult int
@s smResult int
@s paResult int
@s actResult int
@s minorResult int
@s saveState int

% 단어의 종류 (wordtype)
@d noType
@d motionType
@d objectType
@d actionType
@d messageType

% 이동 동사 (motion)
@d N
@d S
@d E
@d W
@d NE
@d SE
@d NW
@d SW
@d U
@d D
@d L
@d R
@d IN
@d OUT
@d FORWARD
@d BACK
@d OVER
@d ACROSS
@d UPSTREAM
@d DOWNSTREAM
@d ENTER
@d CRAWL
@d JUMP
@d CLIMB
@d LOOK
@d CROSS
@d ROAD
@d WOODS
@d VALLEY
@d HOUSE
@d GULLY
@d STREAM
@d DEPRESSION
@d ENTRANCE
@d CAVE
@d ROCK
@d SLAB
@d BED
@d PASSAGE
@d CAVERN
@d CANYON
@d AWKWARD
@d SECRET
@d BEDQUILT
@d RESERVOIR
@d GIANT
@d ORIENTAL
@d SHELL
@d BARREN
@d BROKEN
@d DEBRIS
@d VIEW
@d FORK
@d PIT
@d SLIT
@d CRACK
@d DOME
@d HOLE
@d WALL
@d HALL
@d ROOM
@d FLOOR
@d STAIRS
@d STEPS
@d COBBLES
@d SURFACE
@d DARK
@d LOW
@d OUTDOORS
@d Y2
@d XYZZY
@d PLUGH
@d PLOVER
@d OFFICE
@d NOWHERE

% 동작 동사 (action)
@d ABSTAIN
@d TAKE
@d DROP
@d OPEN
@d CLOSE
@d ON
@d OFF
@d WAVE
@d CALM
@d GO
@d RELAX
@d POUR
@d EAT
@d DRINK
@d RUB
@d TOSS
@d WAKE
@d FEED
@d FILL
@d BREAK
@d BLAST
@d KILL
@d SAY
@d READ
@d FEEFIE
@d BRIEF
@d FIND
@d INVENTORY
@d SCORE
@d QUIT
@d SAVE
@d RESTORE

% 장소 (location)
@d inhand
@d limbo
@d road
@d hill
@d house
@d valley
@d forest
@d woods
@d slit
@d outside
@d inside
@d cobbles
@d debris
@d awk
@d bird
@d spit
@d emist
@d nugget
@d efiss
@d wfiss
@d wmist
@d like1
@d like2
@d like3
@d like4
@d like5
@d like6
@d like7
@d like8
@d like9
@d like10
@d like11
@d like12
@d like13
@d like14
@d brink
@d elong
@d wlong
@d diff0
@d diff1
@d diff2
@d diff3
@d diff4
@d diff5
@d diff6
@d diff7
@d diff8
@d diff9
@d diff10
@d pony
@d cross
@d hmk
@d west
@d south
@d ns
@d y2
@d jumble
@d windoe
@d dirty
@d clean
@d wet
@d dusty
@d complex
@d shell
@d arch
@d ragged
@d sac
@d ante
@d witt
@d bedquilt
@d cheese
@d soft
@d e2pit
@d w2pit
@d epit
@d wpit
@d narrow
@d giant
@d block
@d immense
@d falls
@d steep
@d abovep
@d sjunc
@d tite
@d low
@d crawl
@d window
@d oriental
@d misty
@d alcove
@d proom
@d droom
@d slab
@d abover
@d mirror
@d res
@d scan1
@d scan2
@d scan3
@d secret
@d wide
@d tight
@d tall
@d boulders
@d scorr
@d swside
@d dead0
@d dead1
@d dead2
@d dead3
@d dead4
@d dead5
@d dead6
@d dead7
@d dead8
@d dead9
@d dead10
@d dead11
@d neside
@d corr
@d fork
@d warm
@d view
@d chamber
@d lime
@d fbarr
@d barr
@d neend
@d swend
@d crack
@d neck
@d lose
@d cant
@d climb
@d check
@d snaked
@d thru
@d duck
@d sewer
@d upnout
@d didit
@d ppass
@d pdrop
@d troll

% 사물 (object)
@d NOTHING
@d KEYS
@d LAMP
@d GRATE
@d GRATE_
@d CAGE
@d ROD
@d ROD2
@d TREADS
@d TREADS_
@d BIRD
@d DOOR
@d PILLOW
@d SNAKE
@d CRYSTAL
@d CRYSTAL_
@d TABLET
@d CLAM
@d OYSTER
@d MAG
@d DWARF
@d KNIFE
@d FOOD
@d BOTTLE
@d WATER
@d OIL
@d MIRROR
@d MIRROR_
@d PLANT
@d PLANT2
@d PLANT2_
@d STALACTITE
@d SHADOW
@d SHADOW_
@d AXE
@d ART
@d PIRATE
@d DRAGON
@d DRAGON_
@d BRIDGE
@d BRIDGE_
@d TROLL
@d TROLL_
@d TROLL2
@d TROLL2_
@d BEAR
@d MESSAGE
@d GEYSER
@d PONY
@d BATTERIES
@d MOSS
@d GOLD
@d DIAMONDS
@d SILVER
@d JEWELS
@d COINS
@d CHEST
@d EGGS
@d TRIDENT
@d VASE
@d EMERALD
@d PYRAMID
@d PEARL
@d RUG
@d RUG_
@d SPICES
@d CHAIN

% 작은 순환의 결과 (minorResult)
@d mcEnd
@d mcTryMove
@d mcCommence
@d mcDeath

% 명령 해석의 결과 (inputResult)
@d inputEOF
@d inputMotion
@d inputTransitive
@d inputIntransitive
@d inputSpeak

% 객체 확인의 결과 (objResult)
@d objOK
@d objCantSee
@d objMotion
@d objReport

% 특수 이동의 결과 (smResult)
@d smNormal
@d smStay
@d smGoForIt

% 동작 수행의 결과 (paResult)
@d paDone
@d paNeedObject
@d paTryMove
@d paCommence

% 경계·낱개 상수들
@d minInCave
@d minLowerLoc
@d minForcedLoc
@d maxLoc
@d maxSpec
@d lighted
@d oil
@d liquid
@d caveHint
@d minTreasure
@d maxObj
@d nHints
@d pitchDarkMsg
@d nd
@d chestLoc
@d messageLoc
@d maxPirateLoc
@d maxDeaths
@d maxScore
@d highestClass
@d saveFile
