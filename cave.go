package main

// 동굴의 구조와 위치 속성. 실제 데이터(설명/플래그/이동표)는 cavedata.go에
// 들어 있으며, C reference를 그대로 덤프한 것이다(검증 완료).

// instruction은 이동표의 한 항목이다. "장소 loc에서 mot 방향으로 가려 할 때,
// 조건 cond가 맞으면 dest로 보낸다." (advent.w "Cave connections")
type instruction struct {
	mot  motion
	cond int
	dest location
}

// 위치 속성 비트. flags[loc]에 OR로 담긴다. (advent.w:796-804)
const (
	lighted   = 1   // 어둡지 않은 장소
	oil       = 2   // 기름이 있음
	liquid    = 4   // 액체(물 또는 기름)가 있음
	caveHint  = 8   // 동굴에 들어가는 힌트
	birdHint  = 16  // 새를 잡는 힌트
	snakeHint = 32  // 뱀을 다루는 힌트
	twistHint = 64  // 미로에서 길 잃는 힌트
	darkHint  = 128 // 어두운 방 힌트
	wittHint  = 256 // Witt's End 힌트
)

// 이동 조건(cond) 해석을 위한 헬퍼. (advent.w:832-834)
//
//	cond==0          : 항상 참
//	0<cond<100       : cond% 확률로 참
//	cond==100        : 난쟁이만 빼고 항상 참
//	holds(o)         : 물체 o를 들고 있어야 함
//	sees(o)          : 물체 o가 현재 장소에 있어야 함
//	notProp(o,k)     : prop[o] != k 여야 함
func holds(o object) int          { return 100 + int(o) }
func sees(o object) int           { return 200 + int(o) }
func notProp(o object, k int) int { return 300 + int(o) + 100*k }

// remarkOf는 dest가 가리키는 "그 자리에 머물며 하는 말"을 돌려준다.
// dest가 maxSpec보다 크면 remark 색인이다. (advent.w:768)
func remarkOf(dest location) string {
	return caveRemarks[int(dest)-int(maxSpec)]
}
