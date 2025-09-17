from langchain_mariadb.vectorstores import MariaDBStore


class CustomMariaDBStore(MariaDBStore):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._embedding_table_name = "embeddings"
        self._embedding_id_col_name = "id"
        self._embedding_emb_col_name = "embedding"
        self._embedding_meta_col_name = "metadata"
        self._collection_id = "00000000-0000-0000-0000-000000000001"
