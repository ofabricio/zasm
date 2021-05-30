import sys
from assembler import Assembler


if __name__ == '__main__':
    input = ' '.join(sys.argv[1:]).replace('\\n', '\n')
    data = Assembler().assemble_prettyhex(input)
    print(data)
