with open("hexed.txt") as f:
	data = [i[:-1] for i in f.readlines()]

for item in data:
	print(chr(int(item, 16)), end='')
print()
