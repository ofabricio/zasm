import re
from metadata import registers, instructions, data_def, symbol_table


def tokenize(input):
    tokens = []

    pos = 0
    while len(input := input[pos:]) > 0:

        if pos := match(input, re_nl):
            type = 'NL'

        elif pos := match(input, re_whitespace):
            type = 'Whitespace'

        elif pos := match(input, re_comment):
            type = 'Comment'

        elif pos := match(input, re_alphanum):
            type = 'Ident'

        elif pos := match(input, re_symbols):
            type = 'Symbol'

        else:
            raise Exception(f'Unexpected token {input[0]}')

        if type in ['Whitespace', 'Comment']:
            continue

        token = input[:pos]
        tokens.append(typify(type, token))

    return tokens


def match(str, pattern):
    if m := re.match(pattern, str):
        return m.span(0)[1]
    return None


def typify(type, token):
    if token in symbol_table:
        type = symbol_table[token]
    elif token in registers:
        type = 'Reg'
    elif token in data_def:
        type = 'DD'
    elif token in instructions:
        type = 'Inst'
    elif match(token, re_num):
        type = 'Num'
    return type, token


re_nl = re.compile(r'\n')
re_whitespace = re.compile(r'\s')
re_alphanum = re.compile(r'[0-9a-zA-Z_]+')
re_comment = re.compile(r'\;.*')
re_symbols = re.compile(r'[()\[\]{}|&?*+\-<>$,.@#=]')
re_num = re.compile(r'^(([0-9A-F][0-9A-F])+|[0-9]+d)$')
