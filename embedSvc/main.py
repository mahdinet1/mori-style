from concurrent import futures
import logging
import cv2
import urllib
import numpy as np
import grpc
import embedder_pb2
import embedder_pb2_grpc
from service.embedder.main import EmbedderSvc
from aiModel.openAIClip.main import OpenAiCLIPModel


class Services(embedder_pb2_grpc.EmbedderServicer):
    def __init__(self):
        self.embeder_service = EmbedderSvc(model=OpenAiCLIPModel())

    def ReturnTextVector(self, request, context):
        query = request.query
        vector = self.embeder_service.embed_text(text=query)
        return embedder_pb2.TextToVectorReply(vector=vector)

    def ReturnImageVector(self, request, context):
        print("ReturnImageVector", len(request.image_url))
        images_url = request.image_url
        images = []

        for image_addr in images_url:
            try:
                req = urllib.request.urlopen(image_addr)
                arr = np.asarray(bytearray(req.read()), dtype=np.uint8)
                img = cv2.imdecode(arr, -1)
                images.append(img)
            except Exception as e:
                print(e)

        vectors = self.embeder_service.embed_image(images)
        # Build the response
        response = embedder_pb2.ImageVectorReply()
        for vector in vectors:
            # Create a new Vector message and populate its 'vector' field
            vector_message = embedder_pb2.Vector(vector=vector.tolist())
            response.vectors.append(vector_message)

        return response


def serve():
    port = "50051"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    embedder_pb2_grpc.add_EmbedderServicer_to_server(Services(), server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    print("Server started, listening on " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()
