import anvil
import re

region = anvil.Region.from_file("./W0rLD/region/r.0.-1.mca")

__i = []
y = 17
for y in range(y, y+1):
	for x in range(32):
		try:
			chunk = anvil.Chunk.from_region(region, x, y)
			for te in chunk.tile_entities:
				if str(te["id"]) == "minecraft:sign":
					_i = int(str(te["x"]))
					# print(chr(_i), te["Text1"], end='')
					__i.append(chr(_i))
			# print()
		except anvil.errors.ChunkNotFound as e:
					continue
		except: continue

print(__i)

# Flag{734501LMNGCFY_nr}