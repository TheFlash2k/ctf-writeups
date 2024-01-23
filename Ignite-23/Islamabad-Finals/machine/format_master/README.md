# Format Master

This was probably the only challenge in the entire qualifiers that I had fun while solving.

---

We were only given the `privesc.sh` that we could run as root using sudo. We needed craft a simple `a.cfile`.

The original script has something like this as follows:

```bash
curl http://1337.x/a.cfile
```
Now, after enumerating, I came to know that the `/etc/passwd` file was user-writable, so I wrote my own server's IP and set the domain to be 1337.x.

This gave me control of the file. Now, all I had to do was somehow bypass the filters in the provided bash script by crafting such a perfect cfile that could execute the commands for me.

---

