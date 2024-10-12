import json
import socket

class DB:
    def __init__(self, api_token):
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.socket.connect(('localhost', 8080))
        message = {
            "action": "login",
            "value": api_token
        }
        self.socket.sendall(json.dumps(message).encode())
        response = self.socket.recv(1024).decode()
        try:
            response_data = json.loads(response)
            if response_data.get("state") != 200:
                print(f"Login Fail > {response_data.get('error')}")
                self.socket.close()
        except json.JSONDecodeError as e:
            print(f"Invalid response: {e}")
            self.socket.close()
    def read(self, collection):
        try:
            message = {
                "action": f"read:{collection}",
            }
            self.socket.sendall(json.dumps(message).encode())
            response = self.socket.recv(1024).decode()
            return json.loads(json.loads(response).get("data"))
        except Exception as e:
            self.socket.close()
            print(f"Error: {e}")
    def write(self, collection, data):
        try:
            message = {
                "action": f"write:{collection}",
                "value": json.dumps(data)
            }
            self.socket.sendall(json.dumps(message).encode())
            response = self.socket.recv(1024).decode()
            return json.loads(response)
        except Exception as e:
            self.socket.close()
            print(f"Error: {e}")