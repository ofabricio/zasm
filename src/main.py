import sys
from lexer import tokenize
from parzer import parse
from codegen import generate


def assemble(input):
    tokens = tokenize(input)
    ast = parse(tokens)
    return generate(ast)


if __name__ == '__main__':
    input = ' '.join(sys.argv[1:]).replace('\\n', '\n')
    data = assemble(input)
    print(bytes(data).hex(' ').upper())
