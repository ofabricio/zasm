```
<block> =
    <statement>*

<statement> =
    <label>? ( <directive> | <instr> | <hex> ) <comment>?

<inst> =
    <mov> | <nop> | <db> | ...

<number> =
    <hex> | <decimal> | <octal> | <binary> | <float>

<directive> =
    '@' <alphanum>

<directive_usage> =
    '#' <alphanum_>

<label> =
    <ref_label>
    '.' <alphanum_>
    <alphanum_>

<expr> =
    <term>
    <term> '+' <expr>

<term> =
    <factor>
    <factor> '*' <term>

<factor> =
    <value>
    '(' <expr> ')'

<ref_label> =
    '>'

<ref_label_usage> =
    '>'+
    '<'+

<alphanum> =
    [a-ZA-Z0-9]+

<alphanum_> =
    [a-ZA-Z0-9_]+

<hex> =
    [0-9A-F]+

<hex> =
    [0-9A-F]+

<string> =
    '.+?'

/* Instructions */

<nop> =
    nop

/* Directives */

<define> =
    'define' <alphanum_> <anything>

<include> =
    'include' <string>

<if> =
    'if' <expr> <block> <elseif>* <else>? '@end'

<elseif> =
    'elseif' <expr> <block>

<else> =
    'else' <block>

```
