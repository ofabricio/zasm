
![stage](https://img.shields.io/static/v1?label=stage&message=draft&color=red)

An assembler with a clean syntax.

## General

No comma between operands.

Instructions and registers must be lowercase.

Hexadecimal numbers must be uppercase and appear in pairs.

Numbers default to hexadecimals.

## Comments

```asm
; single line comment
```

## Data

To define data you can either use the data definition keywords or hexadecimal numbers. Examples:

```asm
db 00
dw 01 02 03
db 'hello' 0A 00
dd ?

times 0200 - ($ - $$) db 00
```

Also any two-character hexadecimal value not associated with an instruction is a byte definition.

For example:

```asm
        00
    arr 01 00 02 00 03 00

; The above is the same as:

        db 00
    arr dw 01 02 03
```

```asm
    mov eax 90
    90
    mov ecx ebx

; The above is the same as:

    mov eax 90
    nop
    mov ecx ebx
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
    mov di HEX_OUT + 05
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

HEX_OUT db '0x0000' 00
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
@align 04
```

### @include

Includes a file to assembly if its extension is `.asm`; otherwise, includes as a binary file.

```asm
@include 'example.asm'
```

### @struct

```asm
@struct Point {
    x db 00
    y db 00
}
```

### @print

> Or maybe @debug

Prints information about a piece of code at assembly time.
Good for optimizing things.

```asm
@print {
    mov eax ebx
    add eax 01
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
    83 C0 01    add eax 01

    --- optimized ---

    size:       4 bytes

    89 D8       mov eax ebx
    FF C0       inc eax
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

Shorter form of `@ifndef @define @endif` that wraps the whole file.

> Maybe just make this the default behavior and @unguard to disable it.

```asm
@guard __PRINT__

; code...
```

### @macro

> Explicit macro usage with #

```asm
@macro pushes a b {
    push a
    push b
}

    #pushes eax ebx
    call hello
```

### @fn

> This would abstract a function call.

```asm
    #example(eax, ebx, ecx, out r)
    mov eax r
    ret

@fn example (save eax, ebx, stack v, out r) {
    mov eax ebx
    add eax v
    mov r eax
}
```

The code above would generate this:

```asm
    sub  esp 04         ; "out" keyword reserves space for the returning value.
    push ecx            ; "stack" keyword uses the stack.
    call example
    add  esp 08
    mov  eax [esp-04]   ; [esp-04] is r from "out r"
    ret

example 
    push eax            ; "save" keyword prevents trashing a register.
    mov  eax ebx
    add  eax [esp+04]   ; [esp+04] is v from "stack v"
    mov  [esp+08] eax   ; [esp+08] is r from "out r"
    pop  eax
    ret
```

| Keyword | Description |
| --- | --- |
| `save` | Prevents trashing a register by pushing it to the stack. |
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
code dq 00209A0000000000
```
