#!/usr/bin/env python3
import asyncio, websockets, os, logging, socket

DEFAULT_IP = "192.168.1.125"
DEFAULT_PORT = "8008"
DEFAULT_LOG_IP = "localhost"
DEFAULT_LOG_PORT = "7000"

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
        logging.error(f"Bot not connected: {msg}")
        return
    try: logger.send(msg.encode("utf-8"))
    except Exception as e: logging.error(f"Bot sending err: {e}")


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
    if not IP:
        logging.error(f'Environ variable "IP" not set... using "{DEFAULT_IP}"')
        IP = DEFAULT_IP
    if not PORT:
        logging.error(f'Environ variable "PORT" not set... using "{DEFAULT_PORT}"')
        PORT = DEFAULT_PORT
    if not LOG_IP:
        logging.error(f'Environ variable "LOG_IP" not set... using "{DEFAULT_LOG_IP}"')
        LOG_IP = DEFAULT_LOG_IP
    if not LOG_PORT:
        logging.error(f'Environ variable "LOG_PORT" not set... using "{DEFAULT_LOG_PORT}"')
        LOG_PORT = DEFAULT_LOG_PORT

    # Connect to the log server
    print(IP, PORT, LOG_IP, LOG_PORT)
    logger = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try: logger.connect((LOG_IP, int(LOG_PORT)))
    except Exception as e:
        logging.error(f"Error connecting bot to logger: {e}")
        connectedToLogger = False

    server = websockets.serve(handler, IP, int(PORT))
    asyncio.get_event_loop().run_until_complete(server)
    asyncio.get_event_loop().run_forever()