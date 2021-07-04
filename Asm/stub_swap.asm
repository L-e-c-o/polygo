section .text
    global _start
_start:

jmp short addr

len:
    pop esi
    xor eax, eax
    xor ebx, ebx
    xor ecx, ecx
    mov cl, 0x0                     ; shellcode len (counter)

loop:
    mov byte al, [esi + ecx - 1]
    mov byte bl, [esi + ecx - 2]
    mov byte [esi + ecx - 1], bl
    mov byte [esi + ecx - 2], al
    sub cl, 2
    jnz loop
    jmp short shellcode

addr:
    call len                       ; push shellcode addr on the stack

shellcode: