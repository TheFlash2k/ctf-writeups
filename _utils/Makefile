CC   := gcc
IN   := main.c
OUT  := main
BASE := $(CC) $(IN) -o $(OUT) -w

all:
	# compile with all protections
	$(BASE); $(MAKE) checksec
no:
	# compile with no protections
	$(BASE) -fno-stack-protector -no-pie -zexecstack -Wl,-z,norelro; $(MAKE) checksec
nocanary:
	# compile without a canary
	$(BASE) -fno-stack-protector; $(MAKE) checksec
nopie:
	# compile without PIE (also sets partial relro)
	$(BASE) -no-pie; $(MAKE) checksec
nonx:
	# compile without nx
	$(BASE) -zexecstack; $(MAKE) checksec
nopiecanary:
	# compile without pie and canary
	$(BASE) -fno-stack-protector -no-pie; $(MAKE) checksec
nocanarynx:
	# compile without canary and nx
	$(BASE) -fno-stack-protector -zexecstack; $(MAKE) checksec

checksec:
	checksec $(OUT)