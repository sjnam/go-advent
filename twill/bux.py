#!/usr/bin/env python3
# go-advent의 웹을 gtwill로 조판하기 전에, 저장소 꼭대기에 .bux 두 개를 쓴다.
# go-mmix의 twill/bux.py를 이 저장소에 맞게 줄인 것이다.
#
#   python3 twill/bux.py 웹이름        (저장소 꼭대기에서 부른다)
#
# system.bux: 같은 꾸러미 main에 들었지만 웹에서 나오지 않은 .go 파일들(C 원본에서
#   덤프한 cavedata.go와 objdata.go)이 꾸러미 수준으로 선언한 이름. 웹은 이것들을
#   정의 없이 쓰므로, 그대로 두면 미니 색인에 ???로 찍힌다. 다른 꾸러미는 표준
#   라이브러리뿐이라 go-mmix처럼 남의 .aux를 모을 일은 없다.
#
# 웹이름.bux: twill/hints.txt의 힌트. 서술의 |...|에는 나오지만 코드가 정의하지
#   않는 이름들(대개 커누스의 C 원본에 있던 이름)에 뜻을 준다.
#
# gtwill은 system.bux, 웹 자신의 정의, 웹이름.bux 차례로 읽으므로 웹 자신의 정의가
# 덤프한 데이터를 이기고, 힌트가 둘 다를 이긴다. 둘 다 생성물이라 커밋하지 않는다.
import glob
import os
import re
import subprocess
import sys
import tempfile

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
HINTS = os.path.join(ROOT, 'twill', 'hints.txt')
SECTION = re.compile(r'@[ *\t\n]')


def die(msg):
    sys.exit('twill/bux.py: ' + msg)


def lines_by_section(path, out=None, n=None):
    """웹의 줄마다 (절 번호, 줄)을 준다. @i로 끼운 파일의 절도 센다."""
    if out is None:
        out, n = [], [0]
    for ln in open(path, encoding='utf-8'):
        if ln.startswith('@i '):
            inc = ln[3:].strip().strip('"')
            lines_by_section(os.path.join(os.path.dirname(path), inc), out, n)
            continue
        if SECTION.match(ln):
            n[0] += 1
        out.append((n[0], ln))
    return out


def limbo(stem):
    """웹의 림보(첫 절 앞): 매크로와 @s를 typeset에 물려주려고."""
    w = open(os.path.join(ROOT, stem + '.w'), encoding='utf-8').read()
    m = re.search(r'(?m)^@[ *\t\n]', w)
    return w[:m.start()] if m else w


def title(stem):
    m = re.search(r'\\def\\title\{((?:[^{}]|\{[^{}]*\})*)\}', limbo(stem))
    return m.group(1) if m else stem


def untangled():
    """웹에서 나오지 않은 .go 파일들. gtangle의 산출물은 //line 지시로 시작한다."""
    out = []
    for f in sorted(glob.glob(os.path.join(ROOT, '*.go'))):
        with open(f, encoding='utf-8') as g:
            if not g.readline().startswith('//line '):
                out.append(f)
    return out


def data_decls():
    """(이름, 파일, 형) 목록. 덤프한 데이터는 꼴이 정해져 있다:
    `var 이름 = 형{', `var 이름 형', `var 이름 = 정수'."""
    out = []
    for f in untangled():
        for ln in open(f, encoding='utf-8'):
            m = (re.match(r'var (\w+) = (.+)\{$', ln) or
                 re.match(r'var (\w+) ([^=]+?)\s*$', ln))
            if m:
                out.append((m.group(1), os.path.basename(f), m.group(2)))
            elif m := re.match(r'var (\w+) = -?\d+\s*$', ln):
                out.append((m.group(1), os.path.basename(f), 'int'))
            elif re.match(r'(var|const|type|func) ', ln):
                die('%s: 이 꼴의 선언은 읽을 줄 모른다: %s' % (f, ln.strip()))
    return out


def typeset(stem, exprs):
    """Go 원문들을 gweave로 조판해 \\PB{...}들을 얻는다. 이 웹이 정의한 형은
    굵게 나오도록 웹의 림보 뒤에 @s를 단다."""
    if not exprs:
        return []
    types = set()
    for f in glob.glob(os.path.join(ROOT, '*.w')):
        types |= set(re.findall(r'(?m)^\s*type\s+([A-Za-z_]\w*)',
                                open(f, encoding='utf-8').read()))
    tmp = tempfile.mkdtemp()
    with open(os.path.join(tmp, 'x.w'), 'w', encoding='utf-8') as f:
        f.write(limbo(stem))
        f.writelines('@s %s int\n' % t for t in sorted(types))
        f.write('@* x.\n')
        f.writelines('|%s|\\par\n' % e for e in exprs)
    subprocess.run(['gweave', 'x.w'], cwd=tmp, check=True, capture_output=True)
    tex = open(os.path.join(tmp, 'x.tex'), encoding='utf-8').read()
    out, i = [], 0
    while (i := tex.find('\\PB{', i)) >= 0:
        j, depth = i + 3, 0
        while True:
            depth += {'{': 1, '}': -1}.get(tex[j], 0)
            if depth == 0:
                break
            j += 1
        out.append(tex[i:j + 1])
        i = j + 1
    if len(out) != len(exprs):
        die('gweave가 식 %d개에서 \\PB를 %d개 냈다' % (len(exprs), len(out)))
    return out


def hints():
    rows = []
    for ln in open(HINTS, encoding='utf-8'):
        ln = ln.strip()
        if ln and not ln.startswith('#'):
            rows.append(ln.split(None, 2))
    return rows


def in_bars(name, line):
    """|...| 안에(코드 주석의 |...| 포함) 이름이 낱말로 나오는가."""
    pat = re.compile(r'(?<![\w])' + re.escape(name) + r'(?![\w])')
    return any(pat.search(span) for span in line.split('|')[1::2])


def location(stem, name, loc, secs):
    """힌트의 자리를 gtwill의 꼴로. §만 적으면 그 이름이 처음 나오는 절이고,
    §12처럼 적으면 그 절에 정말 나오는지 검사한다."""
    if loc.startswith('"'):
        return loc
    if not loc.startswith('§'):
        die('%s: 자리 %s를 모른다' % (name, loc))
    num = loc[1:]
    hits = [n for n, l in secs if n > 0 and in_bars(name, l)]
    if not hits:
        die('%s가 %s의 |...|에 나오지 않는다' % (name, stem))
    if num and int(num) not in hits:
        die('%s가 §%s에 나오지 않는다(나오는 절: %s)' % (name, num, hits))
    return '{%s}%s' % (title(stem), num or hits[0])


def main():
    if len(sys.argv) != 2:
        die('쓰는 법: bux.py 웹이름')
    stem = sys.argv[1]
    secs = lines_by_section(os.path.join(ROOT, stem + '.w'))
    data, rows = data_decls(), hints()
    exprs = [t for _, _, t in data] + [m[1:] for _, _, m in rows if m[0] in '=:']
    pbs = iter(typeset(stem, exprs))
    head = '%% twill/bux.py가 만든 것이다. 고치려면 %s를 고친다.\n'
    with open(os.path.join(ROOT, 'system.bux'), 'w', encoding='utf-8') as f:
        f.write(head % '덤프한 .go')
        for name, file, _ in data:
            f.write('@$%s "%s" %s@>\n' % (name, file, next(pbs)))
    with open(os.path.join(ROOT, stem + '.bux'), 'w', encoding='utf-8') as f:
        f.write(head % 'twill/hints.txt')
        for name, loc, mean in rows:
            tex = '\\zip' if mean == '-' else ('=' if mean[0] == '=' else '') + next(pbs)
            f.write('@$%s %s %s@>\n' % (name, location(stem, name, loc, secs), tex))


main()
