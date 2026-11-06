# 거대 동굴 탐험 문학적 프로그램 빌드용 Makefile.
#
#   make            # tangle + 문서 조판
#   make tangle     # advent.w -> advent.go, advent_test.go
#   make doc        # advent.pdf 조판 (한글이라 luatex)
#   make twill      # gtwill로 오른쪽 쪽마다 미니 색인을 단 책꼴 advent-twill.pdf 조판
#   make test       # go test ./...
#   make clean      # 생성물 삭제 (.w 원본과 데이터 .go는 남김)
#
# 매크로(gwebmac.tex, kotexgweb.tex)는 설치된 texmf 트리에서 자동으로 찾는다.

GO      ?= go
GTANGLE ?= gtangle
GWEAVE  ?= gweave
GTWILL  ?= gtwill
REFSORT ?= gtwill-refsort
LUATEX  ?= luatex
WARNS   := 'Overfull\|Underfull\|Error\|Missing\|Undefined'

WSRC := advent.w vocab.w cave.w objects.w game.w parse.w motion.w \
        actions.w liquid.w verbs.w rand.w dwarf.w closing.w death.w \
        score.w save.w launch.w

.PHONY: all build tangle doc twill test clean
.DEFAULT_GOAL := all

all: build doc

build: advent.go
	$(GO) build -o advent ./...

tangle: advent.go

advent.go: $(WSRC)
	$(GTANGLE) advent.w

doc: advent.pdf

advent.pdf: $(WSRC)
	$(GWEAVE) advent.w && luatex advent.tex </dev/null

# gtwill은 커누스의 CTWILL을 GWEB으로 옮긴 것이다: gtwill이 .tex와 .aux를 쓰고,
# TeX 첫 판이 미니 색인 항목을 .mref로 내놓으면 gtwill-refsort가 그것을 가나다순으로
# .smref에 갈무리하고, TeX 둘째 판이 그것을 읽어 미니 색인을 짠다. 묵은 .smref는
# 이제 없는 쪽을 가리키므로, 묵은 .pdf는 실패를 가리므로 첫 판 전에 지운다. gtwill도 gweave처럼 advent.tex를 쓰지만,
# 잡 이름을 advent-twill로 바꿔 조판하므로 .pdf·.log·.toc·.mref·.smref는
# make doc의 것과 섞이지 않는다(advent.tex는 make doc이 늘 새로 쓴다).
#
# gtwill 앞에서 twill/bux.py가 system.bux(cavedata.go와 objdata.go의 선언)와
# advent.bux(twill/hints.txt의 힌트)를 쓴다. 무엇을 어떻게 고르는지는 그 두
# 파일에 적혀 있다. luatex이 치명적 오류로 죽으면 로그에는 소문자 `error'만 남아
# WARNS에 걸리지 않으니, PDF가 생겼는지도 따로 본다. 뜻을 찾지 못한 미니 색인
# 항목은 ???로 찍히므로(\uninitialized) 그 수도 센다. 둘 다 0이어야 한다.
TWILL := $(LUATEX) --interaction=nonstopmode --jobname=advent-twill advent.tex >/dev/null

twill: advent-twill.pdf

advent-twill.pdf: $(WSRC) $(filter-out advent.go advent_test.go,$(wildcard *.go)) \
                  twill/bux.py twill/hints.txt
	python3 twill/bux.py advent
	$(GTWILL) advent.w >/dev/null
	@rm -f advent-twill.smref $@; $(TWILL); \
	  $(REFSORT) < advent-twill.mref > advent-twill.smref && $(TWILL); \
	  if [ -s $@ ]; then printf '%s: 조판 경고 %s개, ??? %s개\n' $@ \
	    `grep -ac $(WARNS) advent-twill.log` `grep -c '\\uninitialized' advent-twill.smref`; \
	  else echo "$@: 만들어지지 않았다(advent-twill.log를 볼 것)"; false; fi

test: tangle
	go test ./...

clean:
	rm -f advent.go advent_test.go advent go-advent
	rm -f advent.tex advent.log advent.toc advent.scn advent.idx advent.pdf
	rm -f advent.aux advent.bux system.bux advent-twill.*  # make twill
