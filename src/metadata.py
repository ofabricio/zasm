"""
Tokens in Lexer:
    ('Type', 'token')
Tokens in Parser:
    ('Type', 'token', ...)
"""

registers = {
    'eax': {'Mod': 3, 'REG': 0, 'RM': 0, 'bits': 32},
    'ecx': {'Mod': 3, 'REG': 1, 'RM': 1, 'bits': 32},
    'edx': {'Mod': 3, 'REG': 2, 'RM': 2, 'bits': 32},
    'ebx': {'Mod': 3, 'REG': 3, 'RM': 3, 'bits': 32},
    'esp': {'Mod': 3, 'REG': 4, 'RM': 4, 'bits': 32},
    'ebp': {'Mod': 3, 'REG': 5, 'RM': 5, 'bits': 32},
    'esi': {'Mod': 3, 'REG': 6, 'RM': 6, 'bits': 32},
    'edi': {'Mod': 3, 'REG': 7, 'RM': 7, 'bits': 32},
}

instructions = {
    'mov': {'operands': 2},
    'add': {'operands': 2},
}

data_def = {
    'db',
    'dw',
}

symbol_table = {
    '(': 'LPar',
    ')': 'RPar',
    '[': 'LBra',
    ']': 'RBra',
    '{': 'LCur',
    '}': 'RCur',

    '+': 'AddOp',
    '-': 'SubOp',
    '*': 'MulOp',
    '/': 'DivOp',

    '@': 'At',
}

"""
Reg = 'Reg'
Inst = 'Inst'
Ident = 'Ident'
Symbol = 'Symbol'
Comment = 'Comment'
Program = 'Program'  # Parser

NL = 'NL'
Whitespace = 'Whitespace'

LPar = 'LPar'
RPar = 'RPar'
LBra = 'LBra'
RBra = 'RBra'
LCur = 'LCur'
RCur = 'RCur'

Expr = 'Expr'   # Parser
BinOp = 'BinOp'  # Parser
AddOp = 'AddOp'
SubOp = 'SubOp'
MulOp = 'MulOp'
DivOp = 'DivOp'

EOF = 'EOF'  # Parser
"""
