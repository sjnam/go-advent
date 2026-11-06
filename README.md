# 거대 동굴 탐험

This is a fully Korean-localized version of *Colossal Cave Adventure*, which Donald Knuth ported from CWEB to GWEB.
The original game was *Adventure 1.0*, a FORTRAN game created by Will Crowther (1975–76) and
significantly expanded by Don Woods (1977), which Knuth rewrote in CWEB in 1998.

Donald Knuth가 CWEB으로 쓴 **Colossal Cave Adventure**를 GWEB로 옮긴 뒤,
**한국어로 완전 한글화**한 것입니다. 원작은 Will Crowther(1975–76)가 만들고
Don Woods(1977)가 크게 확장한 FORTRAN 게임 *Adventure 1.0*이며, Knuth가
1998년에 CWEB으로 다시 썼습니다.

출력(장소 설명·메시지)과 명령어 모두 한국어이며, 말투는 친근한 반말입니다.
마법 주문 `xyzzy`, `plugh`, `plover`, `fee`/`fie`/`foe`/`foo`/`fum`만 원작 그대로
영어로 입력합니다.

```sh
$ go run .
어드벤처에 온 걸 환영해!! 설명을 들어볼래?
** 아니
넌 작은 벽돌 건물 앞, 길 끝에 서 있어. 주위는 온통 숲이야. 건물에서
작은 시냇물이 흘러나와 도랑을 따라 내려가.
* _
```

## 빌드와 실행

이 프로그램은 Knuth의 원작처럼 **문학적 프로그램**(한글 GWEB)으로
작성되어 있습니다. 소스는 `advent.w`와 장별 `.w` 파일들뿐이고, Go 코드
(`advent.go`, `advent_test.go`)와 문서(`advent.pdf`)는 거기서
생성합니다 — 저장소에는 두지 않습니다. Go 1.26 이상과
[GWEB](https://github.com/sjnam/gweb)(`gtangle`/`gweave`)이 필요합니다.
`make twill`에는 GWEB 0.10.7 이상의 `gtwill`과 `gtwill-refsort`, `python3`도
필요합니다.

```sh
make build            # advent.w -> advent.go, advent_test.go -> advent
./advent              # 플레이 (매번 다른 모험 — 현재 시각을 난수 시드로)

make doc              # advent.pdf 조판 (gweave + luatex)
make twill            # 오른쪽 쪽마다 미니 색인을 단 책꼴 advent-twill.pdf (gtwill)
make test             # go test ./...
```

`--seed` 플래그로 난수를 고정할 수 있습니다(테스트·비교용). 시드를 고정하면
게임 전체가 완전히 결정론적이 됩니다.

```sh
go run . --seed=42
```

## 명령 안내

한두 단어짜리 한글 명령을 입력합니다.

- **이동**: `북` `남` `동` `서` `북동` `남서` `위` `아래` `안` `밖` `들어가`
  `나가` `돌아가` `봐` …
- **마법 주문(영어)**: `xyzzy` `plugh` `plover`
- **동작**: `가져가`(집어) `버려`(내려놔) `열어` `닫아` `켜` `꺼` `읽어`
  `먹어` `마셔` `채워` `부어` `던져` `흔들어` `깨워` `먹여` `공격` …
- **사물**: `열쇠` `램프` `막대` `새` `새장` `병` `물` `기름` `창살` `금`
  `보석` `에메랄드` `삼지창` …
- **기타**: `봐`(둘러보기) `소지품` `점수` `도움말` `정보` `간략` `그만`
- **저장/불러오기**: `저장` / `불러오기` (`advent.save` 파일에 한 판 보관)
- **예/아니오**: `예`(응) / `아니`(안)

두 단어는 순서가 자유롭습니다 — `가져가 열쇠`와 `열쇠 가져가` 둘 다 됩니다.

예시:

```text
램프 가져가      램프를 집는다
켜             램프를 켠다
창살 열어       (열쇠가 있으면) 창살을 연다
물 부어         병의 물을 붓는다
```

처음엔 직접 헤매 보는 게 제맛이지만, 막히거나 끝까지 깨고 싶다면
[공략 가이드](GUIDE.md)를 참고하세요.

## 설계

원작 C는 "다상태 시스템"이라 `goto`가 많습니다. 이 포팅은 그 제어 흐름을
값으로 풀어 `goto`를 메인 루프의 본질적인 분기 몇 곳에만 남겼습니다.

- 모든 가변 상태는 하나의 `Game` 구조체에 모음(전역변수 없음)
- C 매크로는 메서드로(`dark()`, `toting()`, `closing()` …)
- 제어 흐름은 명시적 값으로: `simulate()`/`minorCycle()`,
  파서의 `inputResult`, 동작의 `actResult`

코드를 제대로 읽는 길은 문학적 프로그램 그 자체입니다 — `make doc`으로
조판한 **`advent.pdf`**가 들어가며부터 찾아보기까지 이야기로 안내합니다.
장 구성과 원작 대응의 요약은 [ARCHITECTURE.md](ARCHITECTURE.md)에도
있습니다.

### 텍스트와 구조의 분리

방대한 동굴 데이터(이동 명령 740개, 장소 144곳, 객체 67개)는 검증된 C 원본을
실행해 덤프한 **구조**(`cavedata.go`, `objdata.go` — 생성 파일)와, `.w`
문서 안에서 사람이 관리하는 **한글 텍스트**(장소 설명은 동굴 장, 객체
묘사는 객체 장, 힌트는 게임 진행 장)로 나뉩니다. 이 덕에 동굴 구조의
정확성은 생성으로 보장하면서, 번역은 문학적 프로그램 본문에서만
이뤄집니다.

### 한글 파서

원작은 명령을 "앞 5글자"로 잘라 인식하지만, 한글은 UTF-8 멀티바이트라 바이트
절단 시 글자가 깨집니다. 그래서 `lookup`은 **단어 전체로 매칭**하고, 예/아니오는
첫 글자(`예`/`응`/`아`/`안`)로 가립니다.

## 정확성 검증

원작은 단순한 선형 합동 난수 생성기를 쓰고 시드가 하나뿐이라, **시드를 고정하면
게임이 완전히 결정론적**입니다. 한글판도 이 성질을 그대로 가지므로, 시드를
고정한 출력을 golden으로 캡처해 회귀를 막습니다.

```sh
go test ./...
```

`testdata/ko_*.in`(한글 입력)과 `*.golden`(시드 고정 출력)으로 걷기·보물수집·
물/기름·격자문·죽음 시나리오를 검증합니다. 어휘·enum·동굴/객체 구조에 대한
단위 테스트도 함께 돕니다.

> 참고: 영어 원판(C 원본과 byte-identical)은 git 이력에 보존되어 있어 언제든
> 복원·비교할 수 있습니다.

## 라이선스 / 출처

원작 `advent.w`의 저작권은 Don Woods와 Don Knuth에게 있습니다
(© 1998, all rights reserved). 이 저장소의 Go 코드와 한글 번역은 그 동작을
충실히 옮긴 번역물입니다. 원본과 관련 정보:
<http://www.rickadams.org/adventure/>
