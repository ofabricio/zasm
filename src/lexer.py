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

        elif pos := match(input, re_num_hex):
            type = 'Num'

        elif pos := match(input, re_num_dec):
            type = 'Num'

        elif pos := match(input, re_alphanum):
            type = 'Ident'

        elif pos := match(input, re_symbols):
            type = 'Symbol'

        else:
            raise Exception(f'Unexpected token {input[0]}')

        if type in ['Whitespace', 'Comment']:
            continue

        token = input[:pos]
        tokens.append((typify(token) or type, token))

    return tokens


def match(str, pattern):
    if m := re.match(pattern, str):
        return m.span(0)[1]
    return None


def typify(token):
    if token in symbol_table:
        return symbol_table[token]
    if token in registers:
        return 'Reg'
    if token in data_def:
        return 'DD'
    if token in instructions:
        return 'Inst'
    return None


re_nl = re.compile(r'\n')
re_whitespace = re.compile(r'\s')
re_alphanum = re.compile(r'[0-9a-zA-Z_]+')
re_comment = re.compile(r'\;.*')
re_symbols = re.compile(r'[()\[\]{}|&?*+\-<>$,.@#=]')
re_num_hex = re.compile(r'[0-9A-F][0-9A-F]+\b')
re_num_dec = re.compile(r'[0-9]+d\b')
