# go-advent

Donald Knuth가 CWEB으로 작성한 **Colossal Cave Adventure**(`advent.w`)를 Go로
옮긴 것입니다. 원작은 Will Crowther(1975–76)가 만들고 Don Woods(1977)가 크게
확장한 FORTRAN 게임 *Adventure 1.0*이며, Knuth가 1998년에 CWEB으로 다시 썼습니다.

이 포팅의 목표는 단 하나입니다 — **원본과 한 바이트도 다르지 않게 동작할 것.**
동시에 전역변수·goto·매크로로 짜인 C 구조를 Go다운 형태로 재설계했습니다.

```sh
$ go run .
Welcome to Adventure!!  Would you like instructions?
** no
You are standing at the end of a road before a small brick building.
Around you is a forest.  A small stream flows out of the building and
down a gully.
* _
```

## 빌드와 실행

Go 1.26 이상이 필요합니다.

```sh
go build -o advent .   # 빌드
./advent               # 플레이 (매번 다른 모험 — 현재 시각을 난수 시드로)

# 또는 곧장 실행
go run .
```

`--seed` 플래그로 난수를 고정할 수 있습니다(테스트·비교용). 시드를 고정하면
게임 전체가 완전히 결정론적이 됩니다.

```sh
go run . --seed=42
```

### 명령 안내

한두 단어짜리 명령을 입력합니다. 앞 5글자만 인식하므로 `NORTHEAST`는 `NE`로
줄여 써야 `NORTH`와 구분됩니다.

- 이동: `north`(`n`), `south`, `east`, `west`, `up`, `down`, `in`, `out`, `back` …
- 마법 주문: `xyzzy`, `plugh`, `plover`
- 동작: `take`, `drop`, `open`, `on`, `off`, `read`, `eat`, `fill`, `pour` …
- 기타: `look`, `inventory`, `score`, `help`, `info`, `quit`

## 설계

원본 C는 "다상태 시스템"이라 goto가 많습니다(Knuth 본인도 서문에서 "goto를
싫어하면 읽지 말라"고 합니다). 이 포팅은 그 제어 흐름을 값으로 풀어
**goto를 메인 루프의 본질적인 분기 몇 곳에만** 남겼습니다.

- **모든 가변 상태는 하나의 `Game` 구조체**로 모았습니다. 전역변수가 없습니다.
- **C 매크로는 메서드로**: `dark()`, `toting()`, `here()`, `closing()`,
  `waterHere()`, `forcedMove()` 등.
- **제어 흐름은 명시적인 값으로**:
  - 메인 루프는 `simulate()`(큰 순환: 이동+상태보고)와 `minorCycle()`(작은
    순환: 입력+동작)으로 분리.
  - 파서는 `getUserInput()`이 `inputResult`(motion/transitive/intransitive/
    speak/eof)로 의도를 돌려줍니다.
  - 동작은 `actResult`로 원본의 `report`/`change_to`/`try_motion`/`continue`를
    값으로 표현합니다.
  - 이동·작은순환 신호는 `smResult`/`minorResult`/`paResult`.
- **역할별 파일 분리**:

  | 파일 | 내용 |
  | --- | --- |
  | `main.go` | 플래그 파싱, 게임 시작 |
  | `game.go` | `Game` 구조체, 메인 루프, 상태 보고 |
  | `consts.go` | location/object/motion/action/wordtype enum |
  | `vocab.go` | 어휘(약 300단어), lookup, listen |
  | `cave.go` | 동굴 구조·플래그·이동표 정의 |
  | `objects.go` | 객체 이동(carry/drop/here) |
  | `parse.go` | 명령 파싱 |
  | `motion.go` | 이동 처리, 다음 위치 계산 |
  | `actions.go` `liquid.go` `verbs.go` | 동사 디스패치와 구현 |
  | `dwarf.go` | 난쟁이·해적 AI |
  | `closing.go` `death.go` `score.go` | 폐쇄, 죽음·부활, 점수 |
  | `rand.go` `messages.go` | 난수 생성기, 긴 문자열 |

### 생성된 데이터

방대한 게임 데이터(이동 명령 740개, 장소 144곳, 객체 67개와 노트)는 손으로
옮기면 오타가 나기 쉽습니다. 그래서 **검증된 C 원본을 실행해 초기화된
테이블을 Go 소스로 덤프**했습니다.

- `cavedata.go` — 장소 설명·플래그·이동표·remark
- `objdata.go` — 객체 이름·그룹·노트·초기 위치/속성

두 파일은 `// Code generated ...; DO NOT EDIT.`로 표시되어 있으며 손으로 고치지
않습니다.

## 정확성 검증

원본은 단순한 선형 합동 난수 생성기를 쓰고 시드는 단 하나뿐입니다. 따라서
**시드를 고정하면 게임 전체가 결정론적**이고, C 원본과 Go의 출력을 그대로
비교(diff)할 수 있습니다.

- **golden 테스트** (`testdata/*.in`, `*.golden`) — 시드를 고정한 C 원본의
  출력을 정답으로 두고, Go 출력이 한 바이트도 다르지 않은지 확인합니다.
  걷기·보물수집·물/기름·격자문·난쟁이·죽음과 부활·긴 공략 시나리오를 담았습니다.

  ```sh
  go test ./...
  ```

- **무작위 fuzzing** — 무작위로 만든 긴 명령 시퀀스를 C 원본과 Go에 똑같이
  넣어, 어떤 경로를 타든 게임 출력이 일치함을 확인했습니다(난수 시퀀스, 난쟁이
  추격, 죽음 시점, 점수까지 동일).

검증 기준이 되는 C 원본은 CWEB 도구로 만듭니다(시드를 상수로 고정한 사본).

```sh
ctangle advent.w                 # advent.w → advent.c
# advent.c의 시드 한 줄을 상수로 바꾼 사본을 만들고
gcc -w -o advent-fixed advent-fixed.c
```

## 라이선스 / 출처

원작 `advent.w`의 저작권은 Don Woods와 Don Knuth에게 있습니다
(© 1998, all rights reserved). 이 저장소의 Go 코드는 그 동작을 충실히 옮긴
번역물입니다. 원본과 관련 정보: <http://www.rickadams.org/adventure/>
