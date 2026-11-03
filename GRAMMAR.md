prog = stmts

stmts = ( stmt | cmt )*

stmt = datadef | labeldef | directive

directive = dirdump

dirdump = '@dump' int? stmts '@end'

datadef = hex

labeldef = '(\w+|>)'r

hex = '[0-9A-F]+\b'r

str = '".*"'r

cmt = ';[^\n]*'r

int = '\d+'r

hsp = '[ \t]+'r

vsp = '\s+'r
