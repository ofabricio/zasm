```
<statement> =
    <label>? ( <directive> | <instr> | <hex> ) <comment>?

<inst> =
    <mov> | <xor> | <db> | ...

<number> =
    <hex> | <decimal> | <octal> | <binary> | <float>

<directive> =
    '@' <alphanum>

<directive_usege> =
    '#' <alphanum>

<label> =
    <ref_label> | '.' <alphanum_> | <alphanum_>

<ref_label> =
    '>'

<alphanum> =
    [a-ZA-Z0-9]+

<alphanum_> =
    [a-ZA-Z0-9_]+

<hex> =
    [0-9A-F]+

<hex> =
    [0-9A-F]+

<string> =
    ' .*? '

```
