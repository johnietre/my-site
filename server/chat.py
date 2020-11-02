#!/usr/bin/env python3
import socket, select

IP, PORT = "localhost", 7001

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.bind((IP, PORT))
s.listen(5)

def handle(conn):
    while True:
        msg = conn.recv(512).decode("utf-8")
        if msg == "":
            conn.close()
            return
        elif "hello" in msg.lower(): conn.send("Wassup".encode("utf-8"))
        else: conn.send("Huh?!?".encode("utf-8"))

# Handle using select
while True:
    conn, addr = s.accept()
    handle(conn)