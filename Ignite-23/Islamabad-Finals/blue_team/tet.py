def xor_decode(data, key, length):
    decoded_data = bytearray(length)
    for i in range(length):
        decoded_data[i] = data[i] ^ key
    return decoded_data

def fifth_final():
    local_68 = [3, 0x2b, 0x30, 0xd, 0x34, 0x27, 0x30, 0x24, 0x2e, 0x2d, 0x35, 0x24, 0x72, 0x23, 0x76, 0x74, 0x71, 0]
    local_10 = 0x42
    local_14 = 0x12

    decoded_data = xor_decode(local_68, local_10, local_14)
    flag = ''.join([chr(byte) for byte in decoded_data])

    with open("flag.txt", "w") as file:
        file.write(flag)

if __name__ == "__main__":
    fifth_final()
