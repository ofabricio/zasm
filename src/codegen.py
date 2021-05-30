from metadata import registers


def generate(ast):
    if ast[0] == 'Program':
        return [byte for stmt in ast[1] for byte in generate(stmt)]
    if ast[0] == 'DD':
        return gen_data_def(ast)
    if ast[0] == 'Expr':
        return generate(ast[1])
    if ast[0] == 'Num':
        return [int(ast[1], 16)]
    if ast[0] == 'BinOp':
        return [sum(generate(ast[1]) + generate(ast[3]))]
    if ast[0] == 'Inst':
        inst = ast[1]
        if inst == 'mov':
            a = ast[2]
            b = ast[3]
            return gen_mov(a[1], b[1])
        if inst == 'add':
            a = ast[2]
            b = ast[3]
            return gen_add(a[1], b[1])
    return []


def gen_data_def(ast):
    if ast[1] == 'db':
        return [data for stmt in ast[2] for data in generate(stmt)]
    if ast[1] == 'dw':
        return [data for stmt in ast[2] for data in generate(stmt) + [0]]
    return []


def gen_mov(op1, op2):
    # Opcode | Instruction    | Op/En | 64-Bit Mode | Compat/Leg Mode | Description
    # 89 /r  | MOV r/m32,r32  | MR    | Valid       | Valid           | Move r32 to r/m32
    return opcode_r32_r32(0x89, op1, op2)


def gen_add(op1, op2):
    # Opcode | Instruction    | Op/En | 64-Bit Mode | Compat/Leg Mode | Description
    # 01 /r  | ADD r/m32, r32 | MR    | Valid       | Valid           | Add r32 to r/m32.
    return opcode_r32_r32(0x01, op1, op2)


def opcode_r32_r32(opcode, op1, op2):
    mod = registers[op1]['Mod'] << 6
    reg = registers[op2]['REG'] << 3
    rm = registers[op1]['RM'] << 0
    return [opcode, mod+reg+rm]
