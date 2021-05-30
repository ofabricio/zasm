from assembler import Assembler
from lexer import tokenize
from parzer import parse


def test_all():
    def gen(inp): return Assembler().assemble_prettyhex(inp)

    assert_file('golden/lexer', lambda inp: tokenize(inp))
    assert_file('golden/parser-datadef', lambda inp: parse(tokenize(inp)))
    assert_file('golden/parser-directives', lambda inp: parse(tokenize(inp)))
    assert_file('golden/parser-instructions', lambda inp: parse(tokenize(inp)))
    assert_file('golden/codegen-datadef', gen)
    assert_file('golden/codegen-directive-align', gen)
    assert_file('golden/codegen-inst-mov', gen)


def assert_file(file, fn):
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
