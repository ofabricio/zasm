prog = stmts

stmts = ( stmt | cmt )*

stmt = datadef | labeldef | directive

directive = dirdump | dirast

dirdump = '@dump' int? stmts '@end'
dirast  = '@ast' stmts '@end'

datadef = hex

labeldef = '(\w+|>)'r

hex = '[0-9A-F]+\b'r

str = '".*"'r

cmt = ';[^\n]*'r

int = '\d+'r

hsp = '[ \t]+'r

vsp = '\s+'r
