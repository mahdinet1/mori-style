import grpc
import embedder_pb2_grpc
import embedder_pb2


class UnaryClient(object):
    """
    Client for gRPC functionality
    """

    def __init__(self):
        self.host = 'localhost'
        self.server_port = 50051

        # instantiate a channel
        self.channel = grpc.insecure_channel(
            '{}:{}'.format(self.host, self.server_port))

        # bind the client and the server
        self.stub = embedder_pb2_grpc.EmbedderStub(self.channel)

    def get_vector(self, message):
        """
        Client function to call the rpc for GetServerResponse
        """
        message = embedder_pb2.MessageRequest(query=message)
        print(f'{message}')
        response = self.stub.ReturnEmbedder(message)
        return response.vector


if __name__ == '__main__':
    client = UnaryClient()
    result = client.get_vector(message="Hello Server you there?")
    print(f'{result}')
