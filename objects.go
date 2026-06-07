package main

// 객체의 이동과 위치 관리. 불변 데이터(이름/그룹/노트)는 objdata.go에 있고,
// 가변 상태(place/prop/link/first/holding)는 Game에 있다.
// (advent.w "Data structures for objects")

// toting은 객체 t를 들고 있는지 본다. (advent.w:2155)
func (g *Game) toting(t object) bool { return g.place[t] < 0 }

// here는 객체 t가 지금 여기 있는지(들고 있거나 현재 장소에) 본다. (advent.w:2493)
func (g *Game) here(t object) bool { return g.toting(t) || g.place[t] == g.loc }

// drop은 (리스트에 없던) 객체 t를 장소 l에 놓는다. (advent.w:2173)
func (g *Game) drop(t object, l location) {
	if g.toting(t) {
		g.holding--
	}
	g.place[t] = l
	if l < 0 {
		g.holding++
	} else if l > 0 {
		g.link[t] = g.first[l]
		g.first[l] = t
	}
}

// carry는 객체 t를 집어 든다(현재 리스트에서 떼어낸다). (advent.w:2193)
func (g *Game) carry(t object) {
	l := g.place[t]
	if l >= limbo {
		g.place[t] = inhand
		g.holding++
		if l > limbo {
			var r, s object
			for r, s = 0, g.first[l]; s != t; r, s = s, g.link[s] {
			}
			if r == 0 {
				g.first[l] = g.link[s]
			} else {
				g.link[r] = g.link[s] // 리스트에서 t 제거
			}
		}
	}
}

// move는 객체 t를 장소 l로 옮긴다. destroy는 limbo로 보낸다. (advent.w:2189)
func (g *Game) move(t object, l location) { g.carry(t); g.drop(t, l) }
func (g *Game) destroy(t object)          { g.move(t, limbo) }

// isAtLoc은 (여러 부분으로 된) 객체 t가 현재 장소(loc)에 있는지 본다.
// 같은 그룹의 객체는 enum 값이 연속이고 base가 같다는 점을 이용한다. (advent.w:2215)
func (g *Game) isAtLoc(t object) bool {
	if objBase[t] == NOTHING {
		return g.place[t] == g.loc
	}
	for tt := t; objBase[tt] == t; tt++ {
		if g.place[tt] == g.loc {
			return true
		}
	}
	return false
}

// loadObjects는 객체들의 초기 상태를 데이터에서 불러온다.
// (advent.w "Object data"의 new_obj 호출들을 덤프한 결과)
func (g *Game) loadObjects() {
	g.place = objInitPlace
	g.prop = objInitProp
	g.first = objInitFirst
	g.link = objInitLink
	g.holding = objInitHolding
}
