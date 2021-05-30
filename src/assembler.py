from lexer import tokenize
from metadata import registers
from parzer import Parser


class Assembler:

    def assemble(self, input):
        tokens = tokenize(input)
        ast = Parser().parse(tokens)
        return self.generate(ast)

    def assemble_prettyhex(self, input):
        return bytes(self.assemble(input)).hex(' ').upper()

    def generate(self, ast):
        if ast[0] == 'Program':
            return [byte for stmt in ast[1] for byte in self.generate(stmt)]
        if ast[0] == 'DD':
            return self.data_def(ast)
        if ast[0] == 'Num':
            return [int(ast[1], 16)]
        if ast[0] == 'BinOp':
            return self.binop(ast)
        if ast[0] == 'Inst':
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

    def binop(self, ast):
        a = self.generate(ast[1])[0]
        o = ast[2][0]
        b = self.generate(ast[3])[0]
        if o == 'AddOp':
            return [a+b]
        if o == 'SubOp':
            return [a-b]
        if o == 'MulOp':
            return [a*b]
        if o == 'DivOp':
            return [a/b]
        return []

    def data_def(self, ast):
        if ast[1] == 'db':
            return [b for stmt in ast[2] for data in self.generate(stmt) for b in self.to_bytes(data)]
        if ast[1] == 'dw':
            return [data for stmt in ast[2] for data in self.generate(stmt) + [0]]
        return []

    def to_bytes(self, x):
        signed = True if x < 0 else False
        length = 1
        a = abs(x)
        # There must be an equation to remove the if checks.
        if a > 255:
            length = 2
        if a > 255 * 255:
            length = 3
        if a > 255 * 255 * 255:
            length = 4
        return list(x.to_bytes(length, byteorder='big', signed=signed))

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
        return [opcode, mod+reg+rm]
