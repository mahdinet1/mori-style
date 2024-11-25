import abc


class Config:
    DEVICE = "cpu"


class ModelInterface:
    def __init__(self):
        pass

    @abc.abstractmethod
    def get_text_features(self, text):
        pass

    @abc.abstractmethod
    def get_image_features(self, image):
        pass


class EmbedderSvc:
    def __init__(self, model: ModelInterface):
        self.model = model

    def embed_text(self, text):
        return self.model.get_text_features(text)

    def embed_image(self, image):
        return self.model.get_image_features(image)
