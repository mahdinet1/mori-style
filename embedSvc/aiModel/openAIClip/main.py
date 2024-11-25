import torch
from transformers import CLIPProcessor, CLIPModel
from service.embedder.main import ModelInterface


class OpenAiCLIPModel(ModelInterface):
    def __init__(self):
        super().__init__()
        self.model = CLIPModel.from_pretrained("openai/clip-vit-base-patch32")
        self.processor = CLIPProcessor.from_pretrained("openai/clip-vit-base-patch32")

    def get_text_features(self, text):
        inputs = self.processor(text=[text], return_tensors="pt", padding=True)
        with torch.no_grad():
            text_embedding = self.model.get_text_features(**inputs)
        text_embedding = text_embedding.squeeze().tolist()
        return text_embedding

    def get_image_features(self, image):
        inputs = self.processor(images=image, return_tensors="pt")
        with torch.no_grad():
            image_embeddings = self.model.get_image_features(**inputs)
        return image_embeddings




