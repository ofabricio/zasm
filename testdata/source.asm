@dump

; data definition

        00                      ; 1 byte
        01 02                   ; 2 bytes
        01 02 03
        0123456789ABCDEF        ; 8 bytes

; label definition

      a 01                      ; text label for a byte
     ab 01
    a_b 01
      > 01                      ; reference label

    @dump 8
        0102030405060708 1020304050607080
    @end
        0123456789ABCDEF

    @dump
        01 02
    @end

@end
