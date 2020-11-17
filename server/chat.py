#!/usr/bin/env python3
import asyncio, websockets, os, logging, socket

IP = os.getenv("IP")
PORT = os.getenv("BOT_PORT")
LOG_IP = os.getenv("LOG_IP")
LOG_PORT = os.getenv("LOG_PORT")
logger = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
connectedToLogger = True

# Used to send the most important logs to the logger server
def LogMessage(msg: str):
    global logger
    global connectedToLogger
    if not connectedToLogger:
        logging.error(msg)
        return
    try: logger.send(msg.encode("utf-8"))
    except Exception as e: logging.error(e)


async def handler(ws, path):
    while True:
        try:
            msg = await ws.recv()
            if msg == "": break
            print(msg)
            await ws.send(f"Bot: {msg}")
        except Exception as e:
            if str(e) != "code = 1001 (going away), no reason": logging.info(e)
            break

if __name__ == "__main__":
    # Set the backup logger
    logging.basicConfig(level=logging.DEBUG, format="%(asctime)s: %(message)s")

    # Check to make sure the IP and PORT have been set
    if IP == "":
        log.error('Environ variable "IP" not set... using "129.119.172.61"')
        IP = "129.119.172.61"
    if PORT == "":
        log.error('Environ variable "PORT" not set... using "8008"')
        PORT = "8008"
    if LOG_IP == "":
        log.error('Environ variable "LOG_IP" not set... using "localhost"')
        LOG_IP = "localhost"
    if LOG_PORT == "":
        log.error('Environ variable "LOG_PORT" not set... using "7000"')
        LOG_PORT = "7000"

    # Connect to the log server
    logger = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try: logger.bind((LOG_IP, int(LOG_PORT)))
    except Exception as e:
        logging.error(e)
        connectedToLogger = False

    server = websockets.serve(handler, IP, int(PORT))
    asyncio.get_event_loop().run_until_complete(server)
    asyncio.get_event_loop().run_forever()