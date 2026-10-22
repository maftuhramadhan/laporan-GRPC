
import grpc
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "pb"))

import user_pb2
import user_pb2_grpc


def run():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = user_pb2_grpc.UserServiceStub(channel)

        request = user_pb2.UserRequest(user_id=1)
        response = stub.GetUserProfile(request)

        print("=== HASIL DATA USER ===")
        print(f"Nama  : {response.name}")
        print(f"Email : {response.email}")


if __name__ == "__main__":
    run()
