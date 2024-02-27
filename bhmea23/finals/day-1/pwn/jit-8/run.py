#!/usr/local/bin/python
import tempfile
import os

if __name__ == '__main__':
    code = bytes.fromhex(input("CHIP8 Code (hex): "))
    assert len(code) < 100, "Too long!"

    with tempfile.NamedTemporaryFile('wb') as f:
        f.write(code)
        f.flush()

        os.system(f"./chip8 {f.name}")

