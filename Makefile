# 거대 동굴 탐험 문학적 프로그램 빌드용 Makefile.
#
#   make            # tangle + 문서 조판
#   make tangle     # advent.w -> advent.go, advent_test.go
#   make doc        # advent.pdf 조판 (한글이라 luatex)
#   make test       # go test ./...
#   make clean      # 생성물 삭제 (.w 원본과 데이터 .go는 남김)
#
# 매크로(gwebmac.tex, kotexgweb.tex)는 설치된 texmf 트리에서 자동으로 찾는다.

GO      ?= go
GTANGLE ?= gtangle
GWEAVE  ?= gweave

WSRC := advent.w vocab.w cave.w objects.w game.w parse.w motion.w \
        actions.w liquid.w verbs.w rand.w dwarf.w closing.w death.w \
        score.w save.w launch.w

.PHONY: all build tangle doc test clean
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

test: tangle
	go test ./...

clean:
	rm -f advent.go advent_test.go advent go-advent
	rm -f advent.tex advent.log advent.toc advent.scn advent.idx advent.pdf
