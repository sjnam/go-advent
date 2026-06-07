package main

// 게임이 출력하는 긴 문자열들. 원본의 char* 배열을 Go 문자열로 옮긴 것.
// 단계가 진행되면서 여기에 메시지가 계속 추가된다.

// 어두운 곳에서 더 나아가려 할 때의 경고 (advent.w:2711)
const pitchDarkMsg = "It is now pitch dark.  If you proceed you will most likely fall into a pit."

// 거인 방 벽에 적힌 주문 (advent.w:3425)
var incantation = [...]string{"fee", "fie", "foe", "foo", "fum"}

// 난쟁이 칼 공격 결과 메시지 (advent.w:3860)
var attackMsg = [...]string{"it misses", "it gets you",
	"none of them hit you", "one of them gets you"}

const nHints = 8

// 힌트(및 환영 메시지)의 비용. (advent.w hint_cost)
var hintCost = [nHints]int{5, 10, 2, 2, 2, 4, 5, 3}

// 힌트가 제공되기까지 기다리는 횟수. (advent.w hint_thresh)
var hintThresh = [nHints]int{0, 0, 4, 5, 8, 75, 25, 20}

// 각 힌트를 꺼내기 전에 던지는 질문. hintPrompt[0]은 환영 메시지.
var hintPrompt = [nHints]string{
	"Welcome to Adventure!!  Would you like instructions?",
	"Hmmm, this looks like a clue, which means it'll cost you 10 points to\n" +
		"read it.  Should I go ahead and read it anyway?",
	"Are you trying to get into the cave?",
	"Are you trying to catch the bird?",
	"Are you trying to deal somehow with the snake?",
	"Do you need help getting out of the maze?",
	"Are you trying to explore beyond the Plover Room?",
	"Do you need help getting out of here?",
}

// 각 힌트의 본문. hint[0]은 게임 설명(instructions).
var hintText = [nHints]string{
	"Somewhere nearby is Colossal Cave, where others have found fortunes in\n" +
		"treasure and gold, though it is rumored that some who enter are never\n" +
		"seen again.  Magic is said to work in the cave.  I will be your eyes\n" +
		"and hands.  Direct me with commands of one or two words.  I should\n" +
		"warn you that I look at only the first five letters of each word, so\n" +
		"you'll have to enter \"NORTHEAST\" as \"NE\" to distinguish it from\n" +
		"\"NORTH\".  Should you get stuck, type \"HELP\" for some general hints.\n" +
		"For information on how to end your adventure, etc., type \"INFO\".\n" +
		"                        -  -  -\n" +
		"The first adventure program was developed by Willie Crowther.\n" +
		"Most of the features of the current program were added by Don Woods;\n" +
		"all of its bugs were added by Don Knuth.",
	"It says, \"There is something strange about this place, such that one\n" +
		"of the words I've always known now has a new effect.\"",
	"The grate is very solid and has a hardened steel lock.  You cannot\n" +
		"enter without a key, and there are no keys in sight.  I would recommend\n" +
		"looking elsewhere for the keys.",
	"Something seems to be frightening the bird just now and you cannot\n" +
		"catch it no matter what you try.  Perhaps you might try later.",
	"You can't kill the snake, or drive it away, or avoid it, or anything\n" +
		"like that.  There is a way to get by, but you don't have the necessary\n" +
		"resources right now.",
	"You can make the passages look less alike by dropping things.",
	"There is a way to explore that region without having to worry about\n" +
		"falling into a pit.  None of the objects available is immediately\n" +
		"useful for discovering the secret.",
	"Don't go west.",
}
