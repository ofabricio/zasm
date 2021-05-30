import re

instruction = re.compile(r'^(mov)$')
registers32 = re.compile(r'^(eax|ebx)$')


def parse(tokens):
    return Parser().parse(tokens)


class Parser:

    def parse(self, tokens):
        self.pos = 0
        self.tokens = tokens + [('EOF',)]
        return ('Program', self.block('EOF'))

    def block(self, end):
        stmts = []
        while not self.match(end):
            self.skip('NL')
            if s := self.statement():
                stmts.append(s)
            else:
                stmts.append(('Unknown', self.cur()))
                self.advance()
        return stmts

    def statement(self):
        return self.instruction() or self.data_def()

    def instruction(self):
        if i := self.match('Inst'):
            if a := self.match('Reg'):
                if b := self.match('Reg'):
                    return i + (a,) + (b,)
        return ()

    def data_def(self):
        dd = self.match('DD') or ('DD', 'db')
        ns = self.collect(self.expr)
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
        if x := self.expr_start():
            return ('Expr', x)
        return ()

    def expr_start(self):
        return self.bin_op(self.term, self.term_op, self.expr_start)

    def term(self):
        return self.bin_op(self.factor, self.factor_op, self.term)

    def factor(self):
        if self.match('LPar'):
            if exp := self.expr_start():
                if self.match('RPar'):
                    return exp
        if o := self.term_op():
            return ('BinOp', ('Num', '00'), o, self.factor())
        return self.match('Num')

    def term_op(self):
        return self.match('AddOp') or self.match('SubOp')

    def factor_op(self):
        return self.match('MulOp') or self.match('DivOp')

    def bin_op(self, fnA, fnOp, fnB):
        a = fnA()
        if o := fnOp():
            if b := fnB():
                return ('BinOp', a, o, b)
        return a

    def error(self, msg):
        raise Exception(msg)

    def match(self, type):
        curtype = self.cur()[0]
        if curtype != type:
            return ()
        return self.advance()

    def cur(self):
        return self.tokens[self.pos]

    def advance(self):
        cur = self.cur()
        self.pos = self.pos + 1
        return cur
