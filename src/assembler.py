from lexer import tokenize
from metadata import registers
from parzer import Parser


class Assembler:

    def __init__(self):
        self.offset = 0

    def assemble(self, input):
        tokens = tokenize(input)
        ast = Parser().parse(tokens)
        return self.generate(ast)

    def assemble_prettyhex(self, input):
        return bytes(self.assemble(input)).hex(' ').upper()

    def generate(self, ast):
        if ast[0] == 'Program':
            return self.advance_offset(self.program(ast))
        if ast[0] == 'DD':
            return self.advance_offset(self.datadef(ast))
        if ast[0] == 'Directive':
            return self.advance_offset(self.directive(ast))
        if ast[0] == 'Inst':
            return self.advance_offset(self.instruction(ast))
        return []

    def expr(self, ast):
        if ast[0] == 'Num':
            return self.number(ast)
        if ast[0] == 'BinOp':
            return self.binop(ast)
        return 0

    def advance_offset(self, arr):
        self.offset = self.offset + len(arr)
        return arr

    def instruction(self, ast):
        inst = ast[1]
        if inst == 'mov':
            a = ast[2]
            b = ast[3]
            return self.inst_mov(a[1], b[1])
        if inst == 'add':
            a = ast[2]
            b = ast[3]
            return self.inst_add(a[1], b[1])
        return []

    def number(self, ast):
        return int(ast[1], 16)

    def program(self, ast):
        return [byte for stmt in ast[1] for byte in self.generate(stmt)]

    def directive(self, ast):
        name = ast[1]
        if name == 'align':
            return self.directive_align(ast)
        return []

    def directive_align(self, ast):
        val = self.expr(ast[2])
        val = min(val, 64)
        disp = val - (self.offset % val)
        if disp == val:
            disp = 0
        return [0 for _ in range(0, disp)]

    def binop(self, ast):
        a = self.expr(ast[1])
        o = ast[2][0]
        b = self.expr(ast[3])
        if o == 'AddOp':
            return a + b
        if o == 'SubOp':
            return a - b
        if o == 'MulOp':
            return a * b
        if o == 'DivOp':
            return a / b
        return 0

    def datadef(self, ast):
        if ast[1] == 'db':
            return [self.expr(stmt) for stmt in ast[2]]
        if ast[1] == 'dw':
            return [b for stmt in ast[2] for b in [self.expr(stmt), 0]]
        return []

    #
    # def to_bytes(self, x):
    #     signed = True if x < 0 else False
    #     length = 1
    #     a = abs(x)
    #     # There must be an equation to remove the if checks.
    #     if a > 255:
    #         length = 2
    #     if a > 255 * 255:
    #         length = 3
    #     if a > 255 * 255 * 255:
    #         length = 4
    #     return list(x.to_bytes(length, byteorder='big', signed=signed))

    def inst_mov(self, op1, op2):
        # Opcode | Instruction    | Op/En | 64-Bit Mode | Compat/Leg Mode | Description
        # 89 /r  | MOV r/m32,r32  | MR    | Valid       | Valid           | Move r32 to r/m32
        return self.opcode_r32_r32(0x89, op1, op2)

    def inst_add(self, op1, op2):
        # Opcode | Instruction    | Op/En | 64-Bit Mode | Compat/Leg Mode | Description
        # 01 /r  | ADD r/m32, r32 | MR    | Valid       | Valid           | Add r32 to r/m32.
        return self.opcode_r32_r32(0x01, op1, op2)

    def opcode_r32_r32(self, opcode, op1, op2):
        mod = registers[op1]['Mod'] << 6
        reg = registers[op2]['REG'] << 3
        rm = registers[op1]['RM'] << 0
        return [opcode, mod + reg + rm]
