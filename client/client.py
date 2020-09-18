import socket

s = socket.socket()
s.connect(('127.0.0.1',8080))
while True:
    print(s.recv(1024).decode())
    qwe = input()
    if not qwe == '':
        if not qwe == 'exit':
            if qwe == 'w':
                s.send(b'2;0;0.9;90;False;\n')
            if qwe == 'a':
                s.send(b'2;-0.9;0;90;False;\n')
            if qwe == 's':
                s.send(b'2;0.9;0;90;False;\n')
            if qwe == 'd':
                s.send(b'2;0;-0.9;90;False;\n')
            if qwe == 'n':
                s.send(b'2;0;0;90;False;\n')
        else:
            break
s.close()