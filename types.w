% types.w --- gweave 조판 힌트 모음. (advent.w가 림보에서 @i로 끌어온다.)
% 이 web의 const 이름들을 @d로 등록해 CWEB의 매크로처럼 타자체로 조판한다.
% gtangle은 @d를 완전히 무시하므로 생성되는 Go 코드에는 아무 영향이 없다.
% (gweave가 const 선언을 자동 분류하려 하지만, `N motion = iota` 뒤에
% 이름만 잇는 연속 줄은 알아채지 못하므로 여기서 손수 일러 준다.
% @d는 한 줄에 하나씩만 쓴다 --- 이름 뒤에 오는 것은 값으로 보고 버린다.)

% 타입 이름들을 타입처럼(굵게) 조판하도록 등록한다.
@s Reader int
@s Writer int
@s Game int
@s wordtype int
@s vocabEntry int
@s motion int
@s object int
@s action int
@s location int
@s instruction int
@s inputResult int
@s objResult int
@s smResult int
@s paResult int
@s actResult int
@s minorResult int
@s saveState int

% 경계·낱개 상수들
@d minInCave
@d minLowerLoc
@d minForcedLoc
@d maxLoc
@d maxSpec
@d lighted
@d oil
@d liquid
@d caveHint
@d minTreasure
@d maxObj
@d nHints
@d pitchDarkMsg
@d nd
@d chestLoc
@d messageLoc
@d maxPirateLoc
@d maxDeaths
@d maxScore
@d highestClass
@d saveFile
