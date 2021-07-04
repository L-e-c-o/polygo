section .text
    global _start
_start:

jmp short addr

len:
    pop edi
    xor ecx, ecx
    mov cl, 0x0                   ; shellcode len (counter)
loop:
    sub byte [edi + ecx - 1], 0x0 ; random byte used to obfuscate
    jnz loop
    jmp short shellcode

addr:
    call len                      ; push shellcode address on the stack

shellcode:
	