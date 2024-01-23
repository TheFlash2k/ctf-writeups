with open("aaaa.txt") as f:
	lines = f.readlines()

d = []

for line in lines:
	print(line.split(','))