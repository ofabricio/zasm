![stage](https://img.shields.io/static/v1?label=stage&message=draft&color=red)

An assembler.

## How to run

```sh
git clone --depth 1 https://github.com/ofabricio/zasm.git
cd zasm
go run cmd/zasm/main.go --src the_source_code.asm --out the_output.bin
```

## Implemented features

- [ ] Data definition keywords
- [x] Two-character hexadecimal value data definition
- [x] String literal data definition
- [x] Label definition (text and reference)
- [x] Comments
- [ ] Instructions
- [x] @dump
- [ ] @org
- [ ] @align
- [ ] @include
- [ ] @times
- [ ] @macro
- [ ] @if @define @ifndef @guard
- [ ] @16 @32 @64
- [ ] @fn
- [ ] @template

## General guidelines

- No comma between operands.
- Instructions and registers must be lowercase.
- Hexadecimal numbers must be uppercase.
- Numbers default to hexadecimal.

## Comments

```asm
; single line comment
```

## Data definition

To define data, use the data definition keywords.

```asm
d1 'hello' 00    ; define a string of bytes
d1  01           ; define 1 byte     (8 bits)
d2  01           ; define 2 bytes   (16 bits)
d4  01           ; define 4 bytes   (32 bits)
d6  01           ; define 6 bytes   (48 bits)
d8  01           ; define 8 bytes   (64 bits)
d10 01           ; define 10 bytes  (80 bits)
d16 01           ; define 16 bytes (128 bits)
d32 01           ; define 32 bytes (256 bits)
d64 01           ; define 64 bytes (512 bits)
```

Also, any two-character hexadecimal value not associated with an instruction is a byte definition.

```asm
    00
    01 00 02 00 03 00

; the above is the same as:

    d1 00
    d2 01 02 03
```

## Numbers

```asm
1234        ; hex
ABCD        ; hex
1234o       ; octal
1010b       ; binary
1234d       ; decimal
```

## Labels

There are two types of labels: text and reference labels.

#### Text labels

A text label is any word that is not a keyword or a hexadecimal value.

```asm
a_label
    xor eax eax
    ret
```

#### Reference labels

A reference label allows forward or backward references based on the direction and number of arrows.

```asm
val_to_hex_str
    pusha
    mov di HEX_OUT + 05
    mov cx 04
  > mov ax bx
    and al 0F
    cmp al 09
    jle >           ; jump to "add al 30" below
    add al 07
  > add al 30
    mov [di] al
    dec di
    shr bx 04
    dec cx
    jnz <<          ; jump to "mov ax bx" above. A single '<' would jump to "add al 30"
    popa
    ret

HEX_OUT d1 '0x0000' 00
```

## Directives

Directives start with `@`.

### @16 @32 @64

Generates code in 16, 32 or 64 bits. Defaults to `@64`.

### @org

Sets the origin offset.

```asm
@org 7C00
```

### @align

Aligns code or data to the specified boundary.

```asm
@align 04
```

### @include

Includes an `.asm` file for assembly; otherwise, it is included as a binary file.

```asm
@include 'example.asm'
```

### @dump

Prints information about a block of code at assembly time.

```asm
@dump
    mov eax ebx
    add eax 01
@end
@dump
    mov eax ebx
    inc eax
@end
```

The code above would print:

```sh
$ zasm example.asm

    --- 5 bytes ---

    00000000  89 D8       mov eax ebx
    00000002  83 C0 01    add eax 01

    --- 4 bytes ---

    00000005  89 D8       mov eax ebx
    00000007  FF C0       inc eax
```

### @define

```asm
@define VALUE 0123
```

### @if

```asm
@if VERSION == 0
    ; code...
@elseif VERSION == 1
    ; code...
@else
    ; code...
@end
```

### @ifndef

Includes code only if the parameter is not defined.

```asm
@ifndef __PRINT__
@define __PRINT__
    ; code...
@end
```

### @guard

A shorthand for `@ifndef @define @endif` that wraps the whole file.

> Maybe just make this the default behavior and @unguard to disable it.

```asm
@guard __PRINT__

; code...
```

### @macro

> Explicit macro usage with #

```asm
@macro pushes a b
    push a
    push b
@end

    #pushes eax ebx
    call hello
```

### @times

Repeats data n times.

```asm
; @times <n> <dd> <value>
@times 0200 - ($ - $$) d1 00
```

### @fn

> This would abstract a function call.

```asm
    #example(eax, ebx, ecx, out r)
    mov eax r
    ret

@fn example (save eax, ebx, stack v, out r)
    mov eax ebx
    add eax v
    mov r eax
@end
```

The code above would generate:

```asm
    sub  esp 04         ; "out" keyword reserves space for the return value
    push ecx            ; "stack" keyword uses the stack
    call example
    add  esp 08
    mov  eax [esp-04]   ; [esp-04] is r from "out r"
    ret

example 
    push eax            ; "save" keyword prevents trashing a register
    mov  eax ebx
    add  eax [esp+04]   ; [esp+04] is v from "stack v"
    mov  [esp+08] eax   ; [esp+08] is r from "out r"
    pop  eax
    ret
```

| Keyword | Description |
| --- | --- |
| `save` | Prevents trashing a register by pushing it to the stack. |
| `stack` | Pushes the value to the stack before calling the function. |
| `out` | Reserves space on the stack for the return value. If `out` uses a register, the stack is not used. |
| `v`, `r` | These are aliases. |

### @template

Defines data structure templates.

> A template language would allow the definition of templates like the example below.
> Many well-known structs like GDT, IDT, etc. would be predefined, making them easy to set up.

```asm
code @template 'GDT'
    Base        0
    Limit       0
    G           0
    DB          0
    L           1
    AVL         0
    P           1
    DPL         0
    S           0
    Code        1
    Conforming  0
    Readable    1
    Accessed    0
@end
```

The definition above would become:

```asm
code d8 00209A0000000000
```
