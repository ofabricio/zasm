
![stage](https://img.shields.io/static/v1?label=stage&message=draft&color=red)

An assembler with a clean syntax.

## Content

* [General](#general)
* [Comments](#comments)
* [Data](#data)
* [Directives](#directives)
* [Labels](#labels)
* [Numbers](#numbers)

## General

No comma between operands.

Instructions and registers must be lowercase.

Hexadecimal numbers must be uppercase.

Numbers default to hexadecimals.

## Comments

```
; single line comment
```

## Data

```asm
db 0
dw 1 2 3 4
db 'hello' A 0
db ?

times 200 - ($ - $$) db 0
```

Also any two-character hexadecimal value not associated with an instruction is considered a byte definition.

For example:

```asm
        00
    arr 01 00 02 00 03 00

; Is the same as:

        db 0
    arr dw 1 2 3

; And:

    mov eax ebx
    90
    mov ecx ebx

; Is the same as:

    mov eax ebx
    nop
    mov ecx ebx
```

## Numbers

```asm
1234        ; hex
ABCD        ; hex
1010b       ; binary
1234o       ; octal
1234d       ; decimal
```

## Labels

Labels are _global_, _local_ or _reference_ labels.

```asm
global_label
    mov eax ebx
.loc_label
    dec eax
    jnz .loc_label
    ret
```

#### Reference labels

This label allows forward or backward reference according to the arrow direction.

```asm
val_to_hex_str
    pusha
    mov di HEX_OUT + 5
    mov cx 04
  > mov ax bx
    and al 0F
    cmp al 09
    jle >           ; jump to "add al 30" below.
    add al 07
  > add al 30
    mov [di] al
    dec di
    shr bx 04
    dec cx
    jnz <<          ; jump to "mov ax bx" above. A single '<' would jump to "add al 30".
    popa
    ret

HEX_OUT db '0x0000' 0
```

## Directives

Directives start with `@`.

### @16 @32 @64

Generates code in 16, 32 or 64 bits. Defaults to `@64`.

### @org

Organizes offset.

```asm
@org 7C00
```

### @align

Aligns code or data to the specified boundary.

```asm
@align 4
```

### @include

Includes a file to assembly if its extension is `.asm`; otherwise, includes as a binary file.

```asm
@include 'example.asm'
```

### @struct

```asm
@struct Point {
    x db 0
    y db 0
}
```

### @print

> Or maybe @debug

Prints information about a piece of code at assembly time.
Good for optimizing things.

```asm
@print {
    mov eax ebx
    add eax 1
}
@print 'optimized' {
    mov eax ebx
    inc eax
}
```

The code above would print:

```sh
$ zasm example.asm

    size:       5 bytes

    89 D8       mov eax ebx
    83 C0 01    add eax 1

    --- optimized ---

    size:       4 bytes

    89 D8       mov eax ebx
    FF C0       inc eax
```

### @define

Defines a directive.

```asm
@define VALUE 123
```

### @ifndef

Includes code only if the parameter is not defined.

```asm
@ifndef __PRINT__
@define __PRINT__
    ; code...
@endif
```

### @guard

Shorter form of `@ifndef @define @endif` that wraps the whole file.

```asm
@guard __PRINT__

; code...
```

### @fn

This would be a macro definition that abstracts a function call.

```asm
    @call example(eax, ebx, ecx, out r)
    mov eax r
    ret

@fn example (keep eax, ebx, stack v, out r) {
    mov eax ebx
    add eax v
    mov r eax
}
```

The code above would generate this:

```asm
    sub  esp 4         ; "out" keyword reserves space for the returning value.
    push ecx           ; "stack" keyword uses the stack.
    call example
    add  esp 8
    mov  eax [esp-4]   ; [esp-4] is r from "out r"
    ret

example 
    push eax           ; "keep" keyword prevents trashing a register.
    mov  eax ebx
    add  eax [esp+4]   ; [esp+4] is v from "stack v"
    mov  [esp+8] eax   ; [esp+8] is r from "out r"
    pop  eax
    ret
```

| Keyword | Description |
| --- | --- |
| `keep` | Prevents trashing a register. |
| `stack` | Pushes the value before calling the function. |
| `out` | Reserves space on the stack for the returning value. If `out` uses a register the stack is not used. |
| `v`, `r` | These are just aliases to parameters of type `stack` and `out`. Can use any name. |

### @template

Defines templated structs.

> A template language would allow the definition of templates like the example below.
> Many well-known structs like GDT, IDT, etc. would be predefined making it
> easy to setup them.
>
> TODO: describe the template language.

```asm
@template 'GDT' code {
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
}
```

The definition above would become:

```asm
code dq 0x00209A0000000000
```
