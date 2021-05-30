from lexer import tokenize
from parzer import parse
from assembler import Assembler


def test_all():

    assertFile('golden/lexer', lambda inp: tokenize(inp))
    assertFile('golden/parser-datadef', lambda inp: parse(tokenize(inp)))
    assertFile('golden/parser-instructions', lambda inp: parse(tokenize(inp)))
    assertFile('golden/codegen-datadef', Assembler().assemble_prettyhex)
    assertFile('golden/codegen-inst-mov', Assembler().assemble_prettyhex)


def assertFile(file, fn):
    for inp, out in load_test_file(file):
        res = fn(inp)
        if str(res) != out:
            print(f'[Failed test {file}]')
            print(f'Input:\n{inp}')
            print(f'Exp:\n{out}')
            print(f'Got:\n{res}')
            print(f'---')


def load_test_file(path):
    tt = []  # [(in, out)]
    inp = ''
    file = open(path, 'r')
    for line in file.readlines():
        if ';;' in line:
            i, out = line.split(';;')
            tt.append(((inp + i).strip(), out.strip()))
            inp = ''
        else:
            inp += line
    file.close()
    return tt


if __name__ == "__main__":
    test_all()
