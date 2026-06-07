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
// 힌트 비트는 giveHint가 caveHint(8)부터 두 배씩 순회한다:
// 8 동굴 진입, 16 새, 32 뱀, 64 미로, 128 어둠, 256 Witt's End.
const (
	lighted  = 1 // 어둡지 않은 장소
	oil      = 2 // 기름이 있음
	liquid   = 4 // 액체(물 또는 기름)가 있음
	caveHint = 8 // 힌트 비트의 시작(동굴 진입 힌트)
)

// remarkOf는 dest가 가리키는 "그 자리에 머물며 하는 말"을 돌려준다.
// dest가 maxSpec보다 크면 remark 색인이다. (advent.w:768)
func remarkOf(dest location) string {
	return caveRemarks[int(dest)-int(maxSpec)]
}
