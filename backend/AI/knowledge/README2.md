# Core Concepts
## Natural Language Processing (NLP) (https://www.ibm.com/think/topics/natural-language-processing)

**Natural language processing (NLP)** is a subfield of computer science and AI
that uses machine learning to help computers understand and communicate in human language.

### Three NLP Approaches

- **Rules-based NLP** – Early systems used if-then rules for specific prompts
  (e.g., Moviefone). They are limited and not scalable.

- **Statistical NLP** – Uses machine learning to classify text and voice,
  mapping words to vectors for analysis via statistical methods like regression.
  Early applications included spellcheckers.

- **Deep learning NLP** – Uses neural networks on massive datasets for more accurate results. Subtypes include:
  - **Seq2Seq models** (RNN-based) for machine translation.
  - **Transformer models** using tokenization and self-attention (e.g., BERT).
  - **Autoregressive models** trained to predict the next word (e.g., GPT, Llama).
  - **Foundation models** (e.g., IBM Granite) support content generation, named entity recognition, and retrieval-augmented generation.

### NLP Pipeline

- **Text preprocessing** – Tokenization, lowercasing, stop word removal, stemming/lemmatization, and cleaning remove unnecessary elements and standardize text.

- **Feature extraction** – Converts text into numerical representations using Bag of Words, TF-IDF, or embeddings like Word2Vec, GloVe, and contextual embeddings.

- **Text analysis** – Extracts meaning via POS tagging, named entity recognition, dependency parsing, sentiment analysis, topic modeling, and natural language understanding (NLU).

- **Model training** – Uses processed data to train models that learn patterns, make predictions, and generate outputs, refined through evaluation and fine-tuning.

---

## Transformer Architecture (https://medium.com/@amanatulla1606/transformer-architecture-explained-2c49e2257b4c)

The Transformer model revolutionized natural language processing (NLP) by
addressing limitations of previous sequence-to-sequence models. Unlike RNNs,
Transformers process entire sequences in parallel, enabling more efficient training
and better handling of long-range dependencies.

### Core Components of a Transformer

1. **Tokenization** *(Encoder & Decoder)*  
   Splits input text into smaller units, such as words or subwords, to facilitate processing.

2. **Embedding** *(Encoder & Decoder)*  
   Converts tokens into dense vectors that capture semantic information.

3. **Positional Encoding** *(Encoder & Decoder)*  
   Adds information about the position of tokens in the sequence, since Transformers lack inherent sequential processing.

4. **Transformer Blocks**  
   Stacked units that process sequences:
   - **Multi-Head Self-Attention**  
     - Encoder: Attends to all input tokens.  
     - Decoder: Masked attention to past tokens only.
   - **Encoder-Decoder Attention** *(Encoder-Decoder Transformer only)*  
     Allows the decoder to focus on relevant parts of the encoder output.  
   - **Feedforward Neural Networks**  
     Apply non-linear transformations to enhance learning of complex patterns.

5. **Softmax / Output Layer** *(Decoder only)*  
   Converts the decoder’s output into probabilities for each possible next token,
   enabling tasks like language modeling and text generation.

### Model Variants

There are several variants of the Transformer architecture, each tailored for specific tasks:

- **Encoder-Decoder Model** – Used for tasks like machine translation, where the encode
  processes the input sequence, and the decoder generates the output sequence.
- **Encoder-Only Model** – Employed for tasks like text classification and named entity recognition, where understanding the input sequence is paramount.
- **Decoder-Only Model** – Utilized for text generation tasks, where generating coherent and contextually appropriate text is the goal.

## Retrieval-Augmented Generation (RAG) (https://aws.amazon.com/what-is/retrieval-augmented-generation/)

**Retrieval-Augmented Generation (RAG)** enhances the capabilities of large language models (LLMs) by integrating external knowledge sources into the generation process. This approach allows LLMs to produce more accurate, context-aware, and up-to-date responses without the need for retraining the model.

### Why RAG Matters

While LLMs are powerful, they have limitations:

- **Static Knowledge**: LLMs are trained on data available up to a certain point,
  making them unaware of events or information that emerged afterward.

- **Potential Inaccuracies**: Without access to authoritative sources, LLMs might generate responses based on outdated or incorrect information.

RAG addresses these challenges by enabling LLMs to reference external, authoritative knowledge bases before generating responses, ensuring that the information is both current and accurate.

### How RAG Works

1. **Embedding Creation**: Documents are processed to create embeddings—numerical representations capturing the semantic content of the text.

2. **Storage in Vector Database**: These embeddings are stored in a vector database, allowing for efficient similarity searches.

3. **Query Processing**: When a user submits a query, it is transformed into an embedding and compared against the stored embeddings to retrieve the most relevant documents.

4. **Augmented Prompt Generation**: The retrieved documents are used to augment the original query, providing the LLM with context-specific information.

5. **Response Generation**: The augmented prompt is fed into the LLM, which generates a response informed by both its training data and the retrieved information.

### Benefits of RAG

- **Cost-Effective**: By leveraging existing LLMs and integrating external data sources, RAG reduces the need for expensive retraining processes.

- **Up-to-Date Information**: RAG enables LLMs to access the latest information, ensuring responses reflect current knowledge.

- **Enhanced Accuracy**: Access to authoritative sources improves the reliability of the generated content.

- **Transparency**: RAG systems can provide source citations, allowing users to verify the information and increasing trust in the responses.

## Embedding Models (https://medium.com/@nay1228/embedding-models-a-comprehensive-guide-for-beginners-to-experts-0cfc11d449f1)

Embedding models are pivotal in transforming complex data into lower-dimensional
representations, enabling machine learning algorithms to process and understand
intricate inputs like text, images, and graphs. These models have significantly
advanced fields such as natural language processing (NLP), computer vision, and recommendation systems.

### What Are Embeddings?

At the core of embedding models is the concept of **embedding**, which refers to representing high-dimensional data as vectors in a lower-dimensional space. This transformation is critical because it enables machine learning algorithms to process and understand complex inputs, such as words, sentences, images, and even graphs.

For example, words like “king” and “queen” can be represented as vectors that are close to each other in the vector space, reflecting their semantic similarity. This is in contrast to one-hot encoding, where each word is represented as a sparse binary vector, which fails to capture any semantic relationships.

Embeddings are learned through various machine learning techniques, and the resulting vectors can be used in downstream tasks such as classification, clustering, and recommendation.

### Types of Embeddings

- **Word Embeddings**: Represent individual words as vectors. Examples include Word2Vec, GloVe, and FastText. Word embeddings are foundational in NLP tasks and enable models to understand the semantic relationships between words.

  _Example_: In Word2Vec, a word like “king” might be represented by a vector such as `[0.5, 0.8, -0.1, 0.3, 0.9]`. The word “queen” would have a similar vector, showing their semantic similarity, while a word like “apple” would be more distant in the vector space.

- **Sentence Embeddings**: Capture the meaning of entire sentences or paragraphs. Models like the Universal Sentence Encoder and Sentence-BERT (SBERT) generate these embeddings by averaging or pooling word embeddings to create a single vector for the entire sentence.

  _Example_: In a sentence like “The cat sat on the mat,” sentence embeddings might map this entire sentence to a vector that encapsulates the overall meaning, rather than just the individual words.

- **Image Embeddings**: Used in computer vision to represent images as vectors. Convolutional neural networks (CNNs) often serve as the backbone for generating these embeddings, which can then be used for tasks such as image retrieval or classification.

  _Example_: An image of a cat might be mapped to a vector that is close to other images of cats in the embedding space, helping the model recognize similar objects across different images.

- **Graph Embeddings**: In graph embeddings, nodes or entire subgraphs are mapped to vectors that preserve the structural relationships within the graph. Models like DeepWalk and GraphSAGE are used to generate these embeddings.

  _Example_: In a social network graph, embeddings can be used to represent users as vectors, where the proximity of vectors indicates the strength of relationships between users.

### Key Embedding Models

- **Word2Vec** – Static Word Embeddings
- **GloVe (Global Vectors for Word Representation)** – Static Word Embeddings
- **FastText** – Static Word Embeddings
- **ELMo (Embeddings from Language Models)** – Contextual Word Embeddings
- **BERT (Bidirectional Encoder Representations from Transformers)** – Contextual Word Embeddings
- **Sentence-BERT (SBERT)** – Contextual Sentence Embeddings

### Comparing Embedding Models

References:

- https://dev.to/simplr_shcomparing-popular-embedding-models-choosing-the-right-one-for-your-use-case-43p1,

- https://openai.com/index/new-embedding-models-and-api-updates/

- https://huggingface.co/sentence-transformers/LaBSE

- https://docs.pinecone.io/models/mistral-embed

| Model                  | Company  | Contextual? | Deployment   | Cost per 1M tokens         | Multilingual Support | Vector size |
| ---------------------- | -------- | ----------- | ------------ | -------------------------- | -------------------- | ----------- |
| text-embedding-ada-002 | OpenAI   | ✅ Yes      | API-based    | $0.10                      | ✅ Yes               | 1536        |
| LaBSE                  | Google   | ✅ Yes      | Local or API | Free                       | ✅ Yes               | 768         |
| FastText               | Facebook | ❌ No       | Local        | Free                       | ✅ Yes               | 300         |
| mistral-embed          | Mistral  | ✅ Yes      | Local or API | Offers a limited free tier | ✅ Yes               | 1024        |
| N/A                    | TrustLLM | N/A         | N/A          | N/A                        | N/A                  | N/A         |

For eDnevnik, **LaBSE** will be used because it is free, has an ideal (medium) vector size, and supports multiple languages.

## Vector databases (https://www.ibm.com/think/topics/vector-database)

A **vector database** is a specialized system designed to store, manage, and index high-dimensional vector data—numerical representations of complex data such as text, images, audio, and video. Unlike traditional relational databases that organize data in rows and columns, vector databases represent data points as vectors in a continuous vector space. This structure enables efficient similarity searches, making them ideal for applications like semantic search, recommendation systems, and retrieval-augmented generation (RAG).

### Key Features

- **High-Dimensional Indexing**: Vector databases index vectors with many dimensions, allowing for nuanced representation of complex data.
- **Similarity Search**: Support similarity search techniques, such as k-nearest neighbor (k-NN) search, to find data points similar to a query vector.
- **Scalability**: Designed to handle large-scale datasets, vector databases can efficiently process millions or even billions of vectors.
- **Integration with AI Models**: Often integrate with AI models to generate embeddings and perform advanced analytics.

A **vector database** is a specialized system designed to store, manage, and index high-dimensional vector data—numerical representations of complex data such as text, images, audio, and video. Unlike traditional relational databases that organize data in rows and columns, vector databases represent data points as vectors in a continuous vector space. This structure enables efficient similarity searches, making them ideal for applications like semantic search, recommendation systems, and retrieval-augmented generation (RAG) ([IBM](https://www.ibm.com/think/topics/vector-database)).

### Use Cases

- **Semantic Search**: Enhance search capabilities by retrieving information based on meaning rather than exact keyword matches.
- **Recommendation Systems**: Provide personalized recommendations by finding similar items based on user preferences.
- **Retrieval-Augmented Generation (RAG)**: Improve the accuracy and relevance of AI-generated content by grounding it in trusted data sources.
- **Anomaly Detection**: Identify outliers or unusual patterns in data by analyzing vector representations.

### Comparing Vector databases

Before comparing the platforms and databases, it is important to clarify some key terms used in the table:

- **Most Common Index Types** – The algorithm used to efficiently search for similar
vectors in high-dimensional space. Examples include:
  - **HNSW (Hierarchical Navigable Small World)**:  
    Fast approximate nearest neighbor search that builds a layered graph.  
    Large-scale semantic search, recommendation systems, real-time search applications.

  - **IVFFLAT (Inverted File Index + Flat Structure)**:  
    Functions similarly to IVF by clustering vectors, but keeps vectors in each cluster unsummarized and uncompressed, stored in a “flat” list. During search, a brute-force scan is done within selected clusters.

  - **IVF (Inverted File)**:  
    Partitions the vector space into clusters and searches only relevant clusters.  
    Medium to large-scale vector search, offline or batch processing scenarios, when indexing can be precomputed.

  - **Flat**:
    A flat index is the simplest type of index. It stores all vectors in a single list, and searches through all of them to find the nearest neighbors.
    This is extremely memory-efficient, but does not scale well at all, as the search time grows linearly with the number of vectors.

- **Most Common Distance / Similarity Functions**

- **Cosine similarity** – Measures the angle between vectors (ignores magnitude).
  Use case: Best for text embeddings, where semantic similarity matters.
- **Euclidean distance (L2)** – Measures straight-line distance between vectors.
  Use case: Best for image embeddings, where absolute vector differences indicate similarity.
- **Dot product** –  Measures vector alignment including magnitude. Higher
  values = more similar + stronger signals.
  Use caes: Search ranking - boosts results that are both relevant AND popular,
  since popular content has higher magnitude embeddings.


Below is a comparison of popular vector databases and platforms based on these criteria:

References:

- https://aws.amazon.com/rds/mariadb/
- https://mariadb.com/docs/server/reference/sql-structure/vectors/vector-overview
- https://docs.weaviate.io/academy/py/vector_index/overview
- https://docs.weaviate.io/weaviate/config-refs/distances
- https://cloud.google.com/alloydb/docs/ai/measure-vector-query-recall
- https://cloud.google.com/spanner/docs/choose-vector-distance-function
- https://github.com/pgvector/pgvector
- https://cloud.google.com/kubernetes-engine/docs/tutorials/deploy-pgvector
- https://ydb.tech/docs/en/dev/vector-indexes#types
- https://ydb.tech/docs/en/yql/reference/udf/list/knn#functions

| Platform / Database | Company / Provider | Deployment      | Common Index Types Supported | Euclidean | Cosine | Dot Product | Hybrid Queries | Open Source / Proprietary |
| ------------------- | ------------------ | --------------- | --------------------------- | --------- | ------ | ----------- | -------------- | ------------------------- |
| MariaDB Vector      | MariaDB            | Cloud / On-prem | HNSW                        | ✓         | ✓      | ✗           | Yes            | Open Source               |
| Weaviate            | SeMI Technologies  | Cloud / On-prem | HNSW / Flat       | ✓         | ✓      | ✓          | Yes            | Open Source               |
| Google Vertex AI    | Google             | Cloud           | HNSW / IVF / IVFFLAT | ✓         | ✓      | ✓           | Yes            | Proprietary               |
| PGVector            | PostgreSQL         | Cloud / On-prem | IVF / HNSW                  | ✓         | ✓      | ✗           | Yes            | Open Source               |
| YandexDB | Yandex | Cloud / On-prem | Test mode | ✓ | ✓ | ✓ | Yes | Proprietary |

##  LangChain (https://www.ibm.com/think/topics/langchain)

LangChain is an open-source framework designed to facilitate the development of
applications powered by large language models (LLMs). It provides modular
abstractions and components that enable developers to build data-aware and agentic
applications by connecting LLMs to external data sources and allowing them to
interact with their environment.

**Key Features**:

- **Modular Components**: LangChain offers a collection of components such as
prompt templates, chains, and agents that can be combined to create complex workflows.
- **Data Awareness**: It enables LLMs to access and process external data sources,
enhancing their capabilities beyond static knowledge.
- **Agentic Behavior**: LangChain supports the creation of agents that can make
decisions, perform actions, and interact with external systems autonomously.

### RAG with LangChain (http://medium.com/@ai-data-drive/hands-on-langchain-rag-dd639d0576f6)

**Workflow**:

1. **Document Loading**: Utilize LangChain's document loaders to ingest various
data formats such as PDFs, text files, and CSVs.

2. **Document Splitting**: Employ `CharacterTextSplitter` to divide large
documents into smaller, manageable chunks, ensuring they fit within the
LLM's input constraints.

3. **Embedding**: Convert the document chunks into vector embeddings using
models like SBERT`.

4. **Vector Store**: Store the embeddings in a vector database, such as Weaviate,
to facilitate efficient similarity searches.

5. **Retriever, Prompt Construction, and Generation**: In LangChain, these steps are typically combined into a single abstraction, such as the `RetrievalQA` chain, which:  
   - Retrieves the most relevant document chunks from the vector store,  
   - Combines them with the user's query to form the prompt, and  
   - Passes the prompt to the LLM to generate a contextually enriched response.

   Combining these steps simplifies the workflow, reduces boilerplate code,
   ensures consistency in prompt construction, and makes it easier to build and maintain RAG pipelines.

*Note*: Steps 1 and 2 are only necessary if we are extracting knowledge from external
files like PDFs. If we are not using files and instead extracting knowledge from a database
or other structured sources, we can directly form our embeddings using a script to
process the data before storing them in a vector database.

# AI Chatbot implementation
## Generating and Storing Embeddings