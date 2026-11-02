prog = (vs (stmt | cmmt))*

stmt = datadef | instr | labeldef

datadef = hex

labeldef = '(\w+|>)'r

hex = '[0-9A-F]+\b'r

cmmt = ';[^\n]*'r

hs = '[ \t]+'r

vs = '\s+'r
