import re


def parse(tokens):
    return Parser().parse(tokens)


class Parser:

    def __init__(self):
        self.pos = 0
        self.tokens = []

    def parse(self, tokens):
        self.pos = 0
        self.tokens = tokens + [('EOF',)]
        return 'Program', self.block('EOF')

    def block(self, end):
        stmts = []
        while not self.match(end):
            self.skip('NL')
            if s := self.statement():
                stmts.append(s)
            else:
                stmts.append(('Unknown', self.cur()))
                break
        return stmts

    def statement(self):
        return self.directive() or self.instruction() or self.data_def()

    def directive(self):
        if self.match('At') and (ident := self.match('Ident')):
            if (ident[1] == 'align') and (n := self.match('Num')):
                return 'Directive', ident[1], n
        return ()

    def instruction(self):
        if (i := self.match('Inst')) and (a := self.match('Reg')) and (b := self.match('Reg')):
            return i + (a,) + (b,)
        return ()

    def data_def(self):
        dd = self.match('DD') or ('DD', 'db')
        ns = self.collect(lambda: self.match('Num'))
        return dd + (ns,) if ns else ()

    def skip(self, type):
        while self.match(type):
            pass

    def collect(self, fn):
        xs = []
        while x := fn():
            xs.append(x)
        return xs

    def expr(self):
        return self.bin_op(self.term, self.term_op, self.expr)

    def term(self):
        return self.bin_op(self.factor, self.factor_op, self.term)

    def factor(self):
        if self.match('LPar') and (exp := self.expr()) and self.match('RPar'):
            return exp
        if o := self.term_op():
            return 'BinOp', ('Num', '00'), o, self.factor()
        return self.match('Num')

    def term_op(self):
        return self.match('AddOp') or self.match('SubOp')

    def factor_op(self):
        return self.match('MulOp') or self.match('DivOp')

    def bin_op(self, fnA, fnOp, fnB):
        a = fnA()
        if (o := fnOp()) and (b := fnB()):
            return 'BinOp', a, o, b
        return a

    def error(self, msg):
        raise Exception(msg)

    def match(self, type):
        cur = self.cur()
        if cur[0] == type:
            return self.advance()
        return ()

    def cur(self):
        return self.tokens[self.pos]

    def advance(self):
        cur = self.cur()
        self.pos = self.pos + 1
        return cur
