#!/bin/bash

NAME="ret2win_armv5"
if [ "$1" == "DEBUG" ]; then
	echo "[*] Running GDB Server on PORT 7000"
	docker run -it --rm -p8000:8000 -p7000:7000 -e QEMU_GDB_PORT=7000 --name "$NAME" --hostname "$NAME" ropemporium-armv5:latest
else
	docker run -it --rm -p7000:7000 -p8000:8000 --name "$NAME" ropemporium-armv5:latest
fi
