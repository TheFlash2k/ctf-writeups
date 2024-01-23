from base64 import b64decode
from pprint import pprint
import re
from pwn import xor

def dec(__s):
	return b64decode(str(__s).encode()).decode('latin-1')

with open("all.txt", "r") as f:
	data = [i[:-1] for i in f.readlines()]

decoded = []
buf = ""

for i in data:
	dcd = dec(i)
	portions = re.findall('{(.*?)_[A-Za-z0-9]{4}_(.*?)}', dcd)[0]
	buf += portions[1][0]
	decoded.append(''.join(portions))

print(buf)
final_buf = ""
i = 0
for buf in data:
	final_buf += buf[i]
	i+=1

print(final_buf)