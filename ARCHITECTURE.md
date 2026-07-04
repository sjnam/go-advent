# 코드 읽기 안내 (Architecture)

이 포팅은 Knuth의 원작 `advent.w`(CWEB)처럼 **문학적 프로그램**입니다.
사람이 쓰는 소스는 한글 GWEB `.w` 파일들이고, `advent.go`와
`advent_test.go`는 `gtangle`이 생성합니다. 제대로 읽는 길은 `make doc`으로
조판한 **`advent.pdf`** 한 권을 처음부터 읽는 것입니다 — 이 문서는 그
지도만 간단히 그려 둡니다.

## 장(章) 구성

`advent.w`가 뼈대(들어가며, 사용법, 패키지 로드맵)를 두고 `@i`로 장별
파일을 엮습니다. 장 순서는 Knuth 원작의 절 순서를 거의 따릅니다.

| .w 파일 | 장 | 원작 advent.w 절 |
| --- | --- | --- |
| `advent.w` | 들어가며 / 뼈대 | Introduction |
| `vocab.w` | 어휘 (단어 종류·사전·기본 메시지) | The vocabulary |
| `cave.w` | 동굴 (장소 enum·이동표·장소 텍스트) | Cave data / connections |
| `objects.w` | 객체 (연결 리스트·carry/drop·객체 텍스트) | Data structures for objects |
| `game.w` | 게임의 진행 (두 겹 순환·상태 보고·힌트) | Low-level input, The main control loop |
| `parse.w` | 명령 해석 | (main loop의 입력 부분) |
| `motion.w` | 이동 (이동표 해석·플로버 통로·트롤 다리) | Motions |
| `actions.w` | 동작 디스패치와 단순 동사 | Simple verbs |
| `liquid.w` | 액체 자산 (물·기름·집기·놓기·던지기) | Liquid assets |
| `verbs.w` | 그 밖의 동작들 (싸움·먹이기·열닫·읽기·말하기) | The other actions |
| `rand.w` | 난수 (LCG — 결정론의 뿌리) | Random numbers |
| `dwarf.w` | 난쟁이와 해적 | Dwarf stuff |
| `closing.w` | 동굴 폐쇄와 램프 | Closing the cave |
| `death.w` | 죽음과 부활 | Death and resurrection |
| `score.w` | 점수 | Scoring |
| `save.w` | 저장/불러오기 (이 포팅의 추가 기능) | — |
| `launch.w` | 프로그램 띄우기 + golden 시험 | Launching the program |

동굴/객체의 **구조 데이터**(`cavedata.go`, `objdata.go`)만은 문서에 싣지
않고 생성 파일로 곁에 둡니다 — 검증된 C 원본을 실행해 덤프한 것이라
기계가 정확성을 보증하고, 사람은 손대지 않습니다.

## 한눈에 보는 그림

게임은 결국 **두 겹의 순환**입니다.

```text
큰 순환 (major) : 한 장소로 이동 → 그곳을 묘사
   └ 작은 순환 (minor) : 명령을 받아 동작 수행 …  (이동 명령이 나오면 큰 순환으로)
```

이 구조가 `game.w`의 `simulate()`(큰 순환)와 `minorCycle()`(작은 순환)에
그대로 들어 있습니다. 원작의 빽빽한 `goto`는 대부분 **값으로** 풀었습니다
— 각 단계가 결과값을 돌려주고 호출자가 갈래를 틉니다:

| 타입 | 장 | 누가 돌려주나 | 뜻 |
| --- | --- | --- | --- |
| `inputResult` | parse.w | `getUserInput` | EOF / 이동 / 타동사 / 자동사 / 고정메시지 |
| `objResult` | parse.w | `makeObjMeaningful` | 말이 됨 / 안 보임 / 사실은 이동 / 출력했음 |
| `smResult` | motion.w | `handleSpecialMotion` | 보통 / 제자리 / (BACK) 곧장 이동 |
| `paResult` | actions.w | `performAction` | 끝 / 객체 더 필요 / 이동 전환 / 재묘사 |
| `actResult` | actions.w | 개별 동사 | report / change_to / try_motion을 값으로 |
| `minorResult` | game.w | `minorCycle` | 종료 / 이동 / 재묘사 / 죽음 |

모든 가변 상태는 `Game` 구조체 하나에 모이는데, 그 필드들은 각 장이 제
관심사를 논할 때 `@<게임의 상태@>`에 조금씩 보탭니다 — Knuth가 전역
변수를 도입 지점 근처에서 선언하던 방식 그대로입니다.

## 직접 확인하며 읽기

원작은 시드 하나짜리 선형 합동 난수를 쓰므로, **시드를 고정하면 게임이
완전히 결정론적**입니다.

```sh
printf '아니\n안\n가져가 램프\n켜\n그만\n예\n' | go run . --seed=42
go test ./...     # golden(시드 42, C 원본과 byte-identical) + 단위 테스트
```

`.w`를 고쳤다면 `make tangle`(또는 `make test`)이 `advent.go`를 다시
만듭니다.
