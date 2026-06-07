# 코드 읽기 안내 (Architecture)

이 문서는 **go-advent의 Go 코드를 직접 읽으며 이해하려는 사람**을 위한
안내서입니다. 원작은 Knuth가 CWEB(문학적 프로그래밍)으로 쓴 `advent.w`이고,
이 포팅은 그 동작을 충실히 옮기되 C의 전역변수·`goto`·매크로를 Go다운
형태로 재설계했습니다. 그 "재설계의 지도"를 그려 둡니다.

원작 자체에 대한 설명이 궁금하면 먼저 Knuth의 `advent.w` 서문을 읽어 보세요
(이 저장소엔 저작권상 포함하지 않습니다 — `README.md`의 출처 참고).

---

## 1. 한눈에 보는 그림

게임은 결국 **두 겹의 순환**입니다.

```text
큰 순환 (major) : 한 장소로 이동 → 그곳을 묘사
   └ 작은 순환 (minor) : 명령을 받아 동작 수행 …  (이동 명령이 나오면 큰 순환으로)
```

- **큰 순환**은 "걸어 다니기" — 장소를 옮기고 무엇이 보이는지 알려줍니다.
- **작은 순환**은 "그 자리에서 하기" — 집고, 켜고, 읽고… 결과를 알려줍니다.

이 구조가 `game.go`의 `simulate()`(큰 순환)와 `minorCycle()`(작은 순환)에
그대로 들어 있습니다. 여기만 이해하면 절반은 읽은 셈입니다.

---

## 2. 파일 지도

### 로직 (사람이 작성)

| 파일 | 역할 | 핵심 |
| --- | --- | --- |
| `main.go` | 진입점 | `--seed` 파싱, `Game` 생성, `Run()` |
| `game.go` | **심장부** | `Game` 구조체, `Run`/`simulate`/`minorCycle`, 상태 보고, 힌트 |
| `consts.go` | 상수 | location/object/motion/action/wordtype **enum** |
| `vocab.go` | 어휘 | 한글 단어 사전, `lookup`, `listen`, 기본 메시지 |
| `parse.go` | 명령 해석 | `getUserInput`, `makeObjMeaningful` |
| `motion.go` | 이동 | 특수 이동, 이동표 해석(`determineNextLocation`) |
| `actions.go` | 동작 디스패치 | `performAction`, `doIntransitive`/`doTransitive` |
| `liquid.go` | 동작 구현 ① | 물·기름·집기·놓기·던지기 |
| `verbs.go` | 동작 구현 ② | 싸움·먹이기·열고닫기·읽기·말하기 |
| `objects.go` | 객체 이동 | `carry`/`drop`/`here`/`toting` |
| `cave.go` | 동굴 타입 | `instruction`, 위치 속성 비트, `remarkOf` |
| `dwarf.go` | 난쟁이·해적 | 활성화·추격·공격·도둑질 |
| `closing.go` | 폐쇄·램프 | 시계, 램프 수명, 동굴 폐쇄 |
| `death.go` | 죽음·부활 | `handleDeath`, 부활 대사 |
| `score.go` | 점수 | `score`, `printScore`, 등급 |
| `rand.go` | 난수 | LCG `ran`, `pct` |
| `messages.go` | 긴 문자열 | 환영/안내문·힌트·공격 메시지 |

### 데이터

| 파일 | 성격 | 내용 |
| --- | --- | --- |
| `cavedata.go` | **생성**(DO NOT EDIT) | 장소 플래그·이동표(740개)·start 색인 |
| `objdata.go` | **생성**(DO NOT EDIT) | 객체 그룹·노트 오프셋·초기 위치/속성 |
| `cave_text.go` | 수동(한글) | 장소 설명·짧은 설명·remark |
| `obj_text.go` | 수동(한글) | 객체 이름·묘사 |

> **왜 나눴나**: 동굴의 *구조*(어디서 어디로 가는가)는 검증된 C 원본을 실행해
> 그대로 덤프했고(정확성 보장), 사람이 읽고 번역할 *텍스트*만 분리했습니다.
> 그래서 구조 파일은 손대지 않고, 번역은 `*_text.go`에서만 합니다.

---

## 3. 읽는 순서 (추천)

1. **`consts.go`** — 세계의 명사·동사 어휘(enum). 값의 순서가 이동표 색인과
   맞물리니 "순서가 곧 정보"입니다.
2. **`game.go`의 `Game` 구조체** — 게임의 모든 상태가 한곳에. 원작의 수많은
   전역변수가 여기로 모였습니다.
3. **`game.go`의 `simulate` → `minorCycle`** — 2겹 순환의 뼈대.
4. **`parse.go`의 `getUserInput`** — 입력 한 줄이 어떻게 의도로 바뀌는가.
5. **`motion.go`의 `determineNextLocation`** — 이동표를 어떻게 해석하는가.
6. **`actions.go`** — 동사가 어떻게 갈래를 타는가. 그다음 `liquid.go`/`verbs.go`로
   개별 동사 구현을 필요할 때 펼쳐 봅니다.
7. 나머지(`dwarf`/`closing`/`death`/`score`)는 독립적인 하위 시스템이라
   아무 때나 따로 읽어도 됩니다.

---

## 4. 상태: `Game` 구조체 (game.go)

원작은 전역변수로 상태를 들고 다녔지만, 여기서는 **모든 가변 상태를 하나의
`Game`에 모았습니다**. 입출력(`in`/`out`)도 필드라 테스트에서 문자열
버퍼로 바꿔 끼울 수 있습니다(그래서 golden 테스트가 가능).

대략의 묶음:

- 입출력·난수: `in`, `out`, `rx`
- 어휘: `words`, `defaultMsg`, `message`
- 동굴/위치: `loc`, `newloc`, `oldloc`, `oldoldloc`, `visits`
- 객체: `place`, `prop`, `link`, `first`, `holding`
- 파싱: `mot`, `verb`, `obj`, `turns` …
- 난쟁이/해적: `dflag`, `dloc`, `dseen` …
- 폐쇄/죽음: `clock1`, `clock2`, `dying`, `closed` …

원작의 C 매크로는 **메서드**가 됐습니다: `dark()`, `toting()`, `here()`,
`closing()`, `forcedMove()`, `pct()` 등. (예전엔 `holds()`/`sees()` 같은 빌더
매크로도 있었지만, 데이터를 생성으로 바꾸며 더는 필요 없어 정리했습니다.)

---

## 5. 핵심: 두 겹 순환과 제어 흐름 (game.go)

`simulate()`의 골격(주석 생략):

```go
for {
    g.checkInterference()   // 폐쇄·난쟁이가 이동을 막는지
    g.loc = g.newloc        // 실제로 이동
    g.moveDwarves()         // 난쟁이·해적도 움직임
    if g.dying { goto death }

commence:
    if g.reportState() { goto tryMove }   // 장소 묘사 (forced면 곧장 이동)
    switch g.minorCycle() {               // ← 작은 순환
    case mcEnd:      return
    case mcDeath:    goto death
    case mcTryMove:  goto tryMove
    case mcCommence: goto commence        // 불 켜기 등: 제자리 재묘사
    }

tryMove:
    switch g.handleSpecialMotion() { ... } // BACK/LOOK/NOWHERE
    g.determineNextLocation()              // 이동표로 newloc 결정
    if g.dying { goto death }
    continue

death:
    if !g.handleDeath() { return }         // 부활 실패면 종료
    goto commence                          // 부활: 새 장소에서 재묘사
}
```

원작은 `main` 한 함수 안이 `goto`로 빽빽합니다(Knuth도 "goto를 싫어하면 읽지
말라"고 농담했죠). 이 포팅은 그 흐름을 **값으로** 풀어, `goto`를 위처럼 큰
순환의 본질적인 분기 몇 곳에만 남겼습니다. 그 "값"이 다음 여섯 타입입니다.

### 제어 흐름을 나타내는 타입들

| 타입 | 파일 | 누가 돌려주나 | 뜻 |
| --- | --- | --- | --- |
| `inputResult` | parse.go | `getUserInput` | EOF / 이동 / 타동사 / 자동사 / 고정메시지 |
| `objResult` | parse.go | `makeObjMeaningful` | 객체가 말이 됨 / 안 보임 / 사실은 이동 / 메시지 출력함 |
| `smResult` | motion.go | `handleSpecialMotion` | 보통 / 제자리 / (BACK)곧장 이동 |
| `paResult` | actions.go | `performAction` | 끝 / 객체 더 필요 / 이동 전환 / 제자리 재묘사 |
| `actResult` | actions.go | 개별 동사 | 메시지·`change_to`·`try_motion`을 값으로 표현 |
| `minorResult` | game.go | `minorCycle` | 종료 / 이동 / 재묘사 / 죽음 |

C의 관용구 → Go의 값 대응이 핵심입니다:

- `report(m)` (출력 후 continue) → 동사가 `g.rep(m)`을 돌려줌
- `change_to(v)` (다른 동사로 재디스패치) → `actResult{changeTo: v}`
- `try_motion(m)` (이동으로 전환) → `actResult{...}` → `paTryMove`
- `goto transitive/intransitive` → `inputResult`의 두 갈래
- `goto death` → `g.dying` 플래그 + `simulate`의 `death:` 분기

---

## 6. 한 턴의 생애 (데이터 흐름)

명령 한 줄이 처리되는 길을 따라가 봅시다.

```text
listen()            한 줄 읽어 word1/word2 (vocab.go)
   │
getUserInput()      (parse.go)
   │  lookup(word1) → 단어의 종류 판정
   ├─ 이동 단어        → inputMotion ─────────────┐
   ├─ 객체 단어        → makeObjMeaningful()       │
   │                     (여기 없는 것/사실은 이동 처리)
   ├─ 동작 단어        → inputTransitive/Intransitive
   └─ 메시지 단어      → inputSpeak (고정 메시지 출력)
   │
performAction()     (actions.go) 동작이면
   │  doIntransitive / doTransitive → 개별 동사(liquid/verbs)
   │  actResult 해석: 메시지 출력 / change_to / try_motion / 객체 되묻기
   │
   ▼ (이동이면)
handleSpecialMotion() → determineNextLocation()  (motion.go)
   │  이동표(caveTravels)에서 mot에 맞는 명령을 찾고
   │  조건(cond)을 advanceCondition으로 평가해 newloc 확정
   ▼
reportState()       (game.go) 새 장소 묘사 + describeObjects()
```

이동표 한 항목은 `instruction{mot, cond, dest}`(cave.go)입니다. 조건 `cond`의
인코딩(확률·소지·존재·속성)은 `motion.go`의 `advanceCondition` 주석에
정리해 두었습니다.

---

## 7. 어휘와 한글 파서 (vocab.go, parse.go)

- `buildVocabulary()`가 한글 단어 → `{종류, 의미}` 맵(`words`)을 채웁니다.
  명사·동사·이동어·메시지어가 모두 한 맵에 들어가고, 마법 주문만 영어입니다.
- 원작은 단어를 "앞 5글자"로 잘라 인식하지만, 한글은 UTF-8 멀티바이트라
  바이트로 자르면 글자가 깨집니다. 그래서 `lookup`은 **단어 전체**로
  매칭합니다.
- 예/아니오는 `game.go`의 `affirmative()`가 첫 글자(`예`/`응`/`아`/`안`/y/n)로
  가립니다.

---

## 8. 하위 시스템들

각자 독립적이라 따로 읽어도 됩니다.

- **objects.go** — 객체는 장소마다 연결 리스트(`first`/`link`)로 매달려
  있습니다. `carry`/`drop`이 그 리스트를 잇고 끊습니다.
- **dwarf.go** — `dflag`로 난쟁이 활성 단계를 올리고(0→무활동, 2→칼),
  무작위로 움직이며 추격·공격합니다. 해적(`dloc[0]`)은 보물을 훔쳐
  미로 상자에 숨깁니다.
- **closing.go** — 보물을 다 보면 `clock1`이 줄기 시작, 0이 되면 폐쇄 경고,
  `clock2`가 0이면 최종 보관실로 옮깁니다. 램프 수명(`limit`)도 여기서.
- **death.go** — 죽으면 `dying` 표시 → `simulate`의 `death:`가 받아
  `handleDeath`로 부활(최대 3회)/종료를 결정.
- **score.go** — 보물·생존·도달 단계로 점수를 매기고 등급을 출력.

---

## 9. 원작(advent.w)과의 대응

각 Go 파일·함수 주석에 대응하는 advent.w 섹션명을 적어 두었습니다
(예: `// (advent.w "The main control loop")`). 원작과 나란히 읽고 싶다면 그
주석을 길잡이로 삼으세요. 큰 대응은 이렇습니다.

| advent.w 섹션 | Go |
| --- | --- |
| The vocabulary | `vocab.go`, `consts.go` |
| Cave connections / data | `cavedata.go` + `cave.go` + `cave_text.go` |
| Data structures for objects | `objects.go` + `objdata.go` + `obj_text.go` |
| The main control loop | `game.go`(`simulate`/`minorCycle`), `parse.go` |
| Simple verbs / The other actions / Liquid assets | `actions.go`/`verbs.go`/`liquid.go` |
| Motions | `motion.go` |
| Dwarf stuff | `dwarf.go` |
| Closing the cave / Death / Scoring | `closing.go`/`death.go`/`score.go` |

---

## 10. 직접 확인하며 읽기

읽다가 "이게 정말 이렇게 도나?" 싶으면 시드를 고정해 돌려 보세요. 같은
입력엔 항상 같은 출력이 나옵니다.

```sh
printf '아니\n안\n가져가 램프\n켜\n그만\n예\n' | go run . --seed=42
go test ./...     # golden·단위 테스트
```

동굴/객체 데이터가 의심되면 `testdata/ko_*.golden`(시드 고정 캡처)과
대조하면 됩니다. 구조 데이터를 다시 만들고 싶을 때의 덤프 과정은
`README.md`의 검증 절에 적혀 있습니다.
