from typing import List
from langchain_core.embeddings import Embeddings
from sentence_transformers import SentenceTransformer


class LaBSEEmbeddings(Embeddings):
    """LaBSE embeddings wrapper for LangChain"""

    def __init__(self):
        self.model_name = "LaBSE"
        print(f"Loading model {self.model_name}...")
        self.model = SentenceTransformer("LaBSE")
        self.dimension = self.model.get_sentence_embedding_dimension()
        print(f"Model {self.model_name} loaded with dimension {self.dimension}.")

    def embed_documents(self, texts):
        print(f"embed_documents called with: {texts}")

    def embed_query(self, text: str) -> List[float]:
        """
        Embed a single query

        Args:
            text: Query text to embed

        Returns:
            Embedding as list of floats
        """
        embedding = self.model.encode([text], convert_to_numpy=True)
        return embedding[0].tolist()
