with open('got_it.txt') as f:
	lines = f.readlines()

data = ["103"]

for i in data:
	x = "/1/0/3"
	for l in range(len(lines)):
		if x in lines[l]:
			print(lines[l+2])