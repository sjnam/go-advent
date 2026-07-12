@* 프로그램 띄우기.
프로그램은 완성됐다. 남은 일은 난수 씨앗을 뿌리고 게임을 여는 것뿐.
씨앗을 주지 않으면 원작처럼 현재 시각을 쓰고 --- 그래서 매번 다른
모험이 된다 --- \.{--seed}를 주면 그 값으로 고정된다. 끝에 비트 OR로 1을 붙이는 것은
씨앗을 홀수로 만드는 보험이다: 짝수 씨앗은 LCG의 하위 비트를 일찍
죽인다.

@<난수 씨앗을 정하고 게임을 꾸린다@>=
seed := flag.Int64("seed", -1, "고정 난수 시드 (기본값: 원작처럼 현재 시각 기반)")
flag.Parse()

g := &Game{
	in:  bufio.NewReader(os.Stdin),
	out: os.Stdout,
}
if *seed < 0 {
	g.rx = (int(time.Now().Unix()) & 0xfffff) | 1
} else {
	g.rx = (int(*seed) & 0xfffff) | 1
}

@ 끝으로 이 번역의 파수꾼, golden 시험이다. \.{testdata/*.in}의 한글
명령 대본을 시드 42로 재생해, 출력이 같은 이름의 \.{.golden}(시드를
고정한 C~원본에서 캡처)과 byte-identical인지 확인한다. 걷기, 보물
수집, 물과 기름, 격자문, 죽음 --- 대본마다 다른 길로 동굴을 헤집는다.
이 시험이 초록인 한, 두 번 옮겨진 이 동굴은 1977년의 그 동굴이다.

@(advent_test.go@>=
// runGame은 주어진 입력으로 게임을 시드 42로 돌리고 출력을 돌려준다.
// 시드 42는 시드를 고정한 C reference(advent-fixed)와 동일하다.
func runGame(input string) string {
	var out bytes.Buffer
	g := &Game{
		in:  bufio.NewReader(strings.NewReader(input)),
		out: &out,
		rx:  (42 & 0xfffff) | 1,
	}
	g.Run()
	return out.String()
}


@ 대본마다 다른 길로 동굴을 헤집는다 --- 걷기, 보물, 물과 기름, 격자문, 죽음.

@(advent_test.go@>=
func TestGolden(t *testing.T) {
	ins, err := filepath.Glob("testdata/*.in")
	if err != nil {
		t.Fatal(err)
	}
	if len(ins) == 0 {
		t.Skip("testdata/*.in 없음")
	}
	for _, in := range ins {
		in := in
		name := strings.TrimSuffix(filepath.Base(in), ".in")
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(in)
			if err != nil {
				t.Fatal(err)
			}
			got := runGame(string(input))
			goldenPath := strings.TrimSuffix(in, ".in") + ".golden"
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("golden 읽기 실패: %v", err)
			}
			if got != string(want) {
				t.Errorf("출력이 golden과 다름.\n--- got ---\n%s\n--- want ---\n%s", got, want)
			}
		})
	}
}
