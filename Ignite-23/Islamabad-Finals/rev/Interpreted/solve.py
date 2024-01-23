#!/usr/bin/env python3
import string

with open("./enc.data", "rb") as f:
	data = f.read()

key = data[:128]
data = data[128:]

flag = b""
ret = 0

for i in range(len(data)):
	print(chr(key[i % 128] ^ data[i]), end='')
