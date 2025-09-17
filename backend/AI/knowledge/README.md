# Retrieval-Augmented Generation (RAG) Chatbots with Embeddings

## 1. Introduction

RAG chatbots combine **retrieval** from a database of knowledge with
**generation** by a large language model (LLM).
They allow chatbots to answer questions using both structured and unstructured data.

### 1.1. Why Embeddings and Vector Search are Useful

**Embeddings** are numerical representations of text or data in a high-dimensional space. Essentially, they are **vectors** — ordered lists of numbers — that capture the meaning of the data in a form that machines can process.

- **Vector Search** stores these embeddings in a database and lets the system find the most relevant information based on similarity rather than exact matches.
- This is crucial for RAG because user queries rarely match the source text exactly.
- Example: Searching for “John’s performance in math” will retrieve relevant
  teacher notes even if the notes say “John excelled in algebra” instead of using
  the exact words.

By representing data as vectors, the system can perform **semantic similarity searches**,
which go beyond keyword matching to find meaningfully related content.

---

## 2. Key Concepts

### 2.1. Vector Databases

In MariaDB, vector databases are supported natively through the **VECTOR** data type,
which allows storing embeddings and performing similarity searches directly in the database.

Key points for MariaDB:

- **Storing embeddings:** Use the `VECTOR` column type to store high-dimensional
  vectors generated from text or structured data.
- **Indexing methods:** MariaDB supports **HNSW (Hierarchical Navigable Small World)** indexing for efficient nearest neighbor search.
- **Distance functions:** You can use functions like **COSINE_DISTANCE** or **EUCLIDEAN_DISTANCE** to measure similarity between query embeddings and stored vectors.

Example usage:

```sql
CREATE TABLE products (
    name varchar(128),
    description varchar(2000),
    metadata JSON,
    embedding VECTOR(4) NOT NULL,  -- Make sure the dimensions match your AI model's output
    VECTOR INDEX (embedding) M=6 DISTANCE=euclidean)
ENGINE=InnoDB;
```

#### Comparing MariaDB to a dedicated vector database eg. Weaviate vector DB

| Feature                 | MariaDB Vector                        | Weaviate Vector DB                         |
| ----------------------- | ------------------------------------- | ------------------------------------------ |
| **Primary Use Case**    | Relational + vector search            | Pure vector/semantic search                |
| **Indexing**            | HNSW                                  | HNSW and other optimized indexes           |
| **Distance Metrics**    | COSINE\_DISTANCE, EUCLIDEAN\_DISTANCE | COSINE, Euclidean, dot-product             |
| **Data Model**          | Relational tables with VECTOR type    | Object-based with vector + metadata        |
| **APIs**                | SQL                                   | REST, GraphQL                              |
| **Integration with AI** | Needs external embedding generation   | Can integrate directly with AI models      |

If the application is already running on MariaDB, like eDnevnik, and you want to
add a RAG chatbot, using MariaDB Vector is the more practical choice.

#### 2.1.1. Choosing Index Parameters for _ednevnik_

When setting up a vector index in MariaDB, there are three critical parameters to decide on:
the **graph connectivity (M)**, the **distance function**, and the **vector dimension size**.

---

#### M = 6

- **Definition:**  
  `M` defines the maximum number of connections (edges) a node can have in each layer
  of the HNSW graph.
- **Why important:**  
  Higher `M` values improve recall but also increase memory usage. Lower values save memory but can reduce accuracy.
- **Choice for _ednevnik_:**  
  Since the system runs with **8 GB of RAM**, `M = 6` strikes a balance:
  - Keeps memory usage manageable
  - Provides sufficiently accurate similarity results for semantic text search
  - Ensures queries remain efficient at scale

---

#### Distance Function = Cosine

- **Options:**
  - `EUCLIDEAN_DISTANCE` → better for geometric similarity (e.g., image vectors)
  - `COSINE_DISTANCE` → better for semantic similarity (e.g., text embeddings)
- **Choice for _ednevnik_:**  
  Because we are embedding **Markdown text and descriptions**, cosine similarity
  is the most appropriate. It measures the _angle_ between vectors, making it robust
  for semantic meaning, regardless of vector magnitude.

---

### 2.2. Metadata / Filtering

**Metadata** is structured information stored alongside embeddings to help narrow
down the search before performing vector similarity queries. This allows
for **more precise and efficient retrieval**.

Key points:

- **Structured information stored alongside embeddings:** Examples include pupil ID, class, subject, or date.
- **Used to filter relevant rows before vector search:** Instead of searching the entire table, you can first filter by metadata (e.g., only math grades or a specific pupil) and then perform the vector similarity search.
- **Examples of metadata fields:**
  - `entity_type`: “pupil”, “teacher”, “assignment”
  - `subject`: “math”, “history”
  - `user_id`: ID of the student or teacher
  - `category`: “grades”, “notes”, “attendance”

Example usage in MariaDB:

```sql
-- Find top 3 most similar embeddings for math pupils only
SELECT pupil_id, name
FROM pupil_embeddings
WHERE subject = 'math'
ORDER BY COSINE_DISTANCE(grade_vector, ?) ASC
LIMIT 3;
```

**Application-specific filtering:**

- In **eDnevnik**, **metadata filtering will use `account_type`** to control access:
  - For example, admins can access all embeddings, teachers can access only
    embeddings related to their assigned tenants, and pupils can access only
    embeddings associated with their own accounts.

---

### 2.3. Comparing MariaDB vector to Cloud-based Vector Solutions (Google Vertex AI Vector Search)

**Key Features:**
- **Integration:** Native integration with Google Cloud ecosystem
- **Indexing:** Uses Approximate Nearest Neighbor (ANN) algorithms
- **Scalability:** Handles billions of vectors with sub-second latency
- **Multimodal:** Supports text, image, and other data types

**Workflow Integration:**
```
Source Bucket (GCS) → Document Processing → Embedding Generation 
→ Vector Index → Search Endpoint → RAG Application
```

**Comparison with MariaDB Vector:**

| Step          | Vertex AI Workflow                        | MariaDB Approach                  |
|---------------|-------------------------------------------|-----------------------------------|
| **Storage**   | GCS buckets                               | MariaDB tables                    |
| **Embedding** | Managed (Vertex AI, multimodal)           | External, manual (text only)      |
| **Vector Index** | Vertex AI Vector Search  | MariaDB VECTOR + HNSW             |
| **Search API**| Managed endpoint                          | SQL queries                       |
| **Scale**     | Billions of vectors                       | Limited by DB hardware            |
| **Multimodal**| Yes (text, image, etc.)                   | No (text only, unless extended)   |
| **Cost**      | Pay-per-use                               | DB

### 2.3.1. Airtable vs MariaDB for RAG and Embeddings

**Airtable** is a cloud-based platform that combines spreadsheet-like views with database features and a user-friendly interface. While it is popular for rapid prototyping, collaboration, and lightweight data management, it has important differences compared to MariaDB when used for Retrieval-Augmented Generation (RAG) and embeddings.

#### Key Differences

| Feature                | Airtable                                 | MariaDB Vector                        |
|------------------------|------------------------------------------|---------------------------------------|
| **Data Model**         | Spreadsheet-style tables, no native vector type | Relational DB with native VECTOR type |
| **Embeddings Support** | No native support; must store as arrays/strings | Native VECTOR column, optimized for embeddings |
| **Similarity Search**  | Not supported; requires external processing | Built-in (HNSW, COSINE_DISTANCE, etc.) |
| **Indexing**           | Basic (single/multi-field), no vector index | Advanced (HNSW, vector indexes)       |
| **APIs**               | REST, GraphQL, GUI automations           | REST, SQL, connectors, programmatic APIs    |
| **Scalability**        | Good for small/medium datasets           | Scales to millions of rows, large vectors |
| **Automation**         | Built-in automations, scripting          | Requires external tools or triggers   |
| **Integration**        | Easy with Zapier, Make, etc.             | Integrates with backend apps, Python, etc. |
| **Cost**               | Subscription-based, per user/record      | Open source/self-hosted or enterprise |

#### When to Use Airtable

- **Prototyping**: Quickly build and share simple RAG demos or data annotation tools.
- **Collaboration**: Teams can edit, comment, and manage data easily.
- **No-code/Low-code**: Non-developers can manage data and workflows.

#### Limitations for RAG/Embeddings

- **No native vector search**: You cannot perform similarity search directly in Airtable. Embeddings must be exported to another system (e.g., Python, cloud vector DB) for retrieval.
- **Performance**: Not suitable for large-scale, high-performance RAG applications.
- **Data types**: Embeddings are stored as text/arrays, not as optimized vectors.

#### When to Use MariaDB

- **Production RAG systems**: Native vector support, fast similarity search, and scalability.
- **Integration**: Works well with backend services, APIs, and analytics.
- **Security and control**: More options for access control, backups, and compliance.

#### Example Workflow Comparison

- **Airtable**:  
  1. Store text and metadata in Airtable.  
  2. Generate embeddings externally (Python, cloud function).  
  3. Store embeddings as arrays/strings in Airtable.  
  4. For similarity search, export data to a script or another DB.

- **MariaDB**:  
  1. Store text, metadata, and embeddings (VECTOR) in MariaDB.  
  2. Use SQL queries for filtering and similarity search (COSINE_DISTANCE, HNSW).  
  3. Integrate directly with RAG pipelines (LangChain, custom apps).

**Summary:**  
Airtable is excellent for prototyping, collaboration, and simple data management, but lacks the native vector search and scalability needed for production RAG chatbots. MariaDB Vector is the better choice for robust, efficient, and scalable RAG and

### 2.4. System prompts

A system prompt is a set of instructions, rules, or context given to an AI model to guide its behavior and output. Unlike a user prompt, which is the direct question or request from the user, the system prompt sets the stage for the AI's persona and the parameters for its response.

```json
{
  "system_prompts": {
    "school_assistant": {
      "prompt": "You are an educational assistant...",
      "max_tokens": 150,
      "temperature": 0.7
    },
    "interview_coach": {
      "prompt": "You are a professional interview coach...",
      "max_tokens": 200,
      "temperature": 0.5
    }
  }
}
```

We'd handle system prompts with MariaDB by storing them in a dedicated table.
This approach allows us to manage and retrieve the prompts dynamically, rather
than hard-coding them into our application.

## 3. Data Modeling for Embeddings

Proper data modeling ensures embeddings are **accurate, useful, and easy to query**.

- **Entities to embed:** Choose what the chatbot will retrieve (e.g., users, documents, items).
- **Granularity:** Decide if each row represents a single entity or combined data.
  - **Single entity:** Simple to retrieve and update.
  - **Combined data:** Captures richer context but harder to maintain.
- **Content:** Select the text or structured data to embed (e.g., `"John, grade
      88, excellent in algebra"`).
- **Metadata:** Store structured info alongside embeddings for filtering and context (e.g., `{"subject": "math", "teacher_id": 42, "semester": "fall_2025"}`).

In **eDnevnik**:

- **Entities embedded:** Pupils, tenants, teachers and grades.
- **Granularity:** One row per entity (easier to retrieve and update).
- **Content:** Key fields combined into text embeddings (e.g., pupil name, tenant name, grades, notes).
- **Metadata:** Includes `entity_type`, `user_id`, `tenant_id`, `account_type`,
  and other fields for filtering and access control.

---

## 4. Choosing an embedding model and LLM

When building multilingual applications, selecting the right embedding model is
crucial for achieving optimal semantic understanding and search performance.
Below are four popular multilingual embedding models that represent different
trade-offs between performance, resource requirements, and language coverage.

| Model                                      | File Size (approx) | Model Parameters    | Vector Size | Other Important Details                                                |
|--------------------------------------------|--------------------|---------------------|-------------|----------------------------------------------------------------------|
| sentence-transformers/LaBSE                 | ~1.23 GB           | ~192 million        | 768         | Multilingual, supports 109 languages, BERT-based with Tanh activation |
| BAAI/bge-multilingual-gemma2                | ~1.11 GB           | ~340 million        | 768         | Large parameter count, multilingual                                   |
| sentence-transformers/distiluse-base-multilingual-cased-v2 | ~300 MB           | ~82 million         | 512         | DistilBERT-based, smaller & faster, multilingual                      |
| intfloat/multilingual-e5-large              | ~2.15 GB           | ~770 million        | 1024        | Large model, high dimensional vector, state-of-the-art multilingual  |

For this application, we want a balance between performance, resource usage,
and language coverage. **LaBSE** is the best choice as it provides strong multilingual support,
good semantic quality, and moderate resource

#### Testing and Validation Strategy

Before making a final decision, consider implementing a benchmark test with your specific data:

1. **Create representative datasets** in your target languages
2. **Measure semantic similarity accuracy** using your application's typical queries
3. **Benchmark inference speed** in your deployment environment
4. **Evaluate memory and storage requirements** for your infrastructure
5. **Test edge cases** with domain-specific terminology and less common language constructs

Remember that the "best" model is the one that meets your specific requirements while staying within your resource constraints and performance expectations.


## 5. Embedding Generation

Generating embeddings is a critical step in building a RAG system, as it transforms
your raw data into a format that can be **searched semantically**.

- **Choosing what text to embed:**
  - Identify the most relevant text or fields that capture the meaning of the entity.
  - What text to embed in **eDnevnik**?:
    - For example a pupil record in eDnevnik, combine pupil name, gender, age,
      religion, enrolled classes, archived claseses

- **Embedding Update Strategy:**
  - **Periodic jobs:** Run at regular intervals (cron jobs, Celery tasks) to regenerate embeddings for the dataset.
    - **Pros:**
      - Simple to implement and maintain
      - Ensures all embeddings are refreshed consistently
      - Useful for large datasets where real-time updates are not critical
    - **Cons:**
      - May introduce latency between data changes and updated embeddings

  - **On-insert/update:** Generate embeddings immediately when a new row is added or updated.
    - **Pros:**
      - Keeps embeddings up-to-date in real time
      - Ideal for applications that require immediate retrieval accuracy
    - **Cons:**
      - More complex to implement (triggers, hooks, or event-driven architecture)
      - May add latency to the insert/update operation
      - Could increase load on the embedding service if updates are frequent

---

## 6. Storing Embeddings

- **VECTOR column for embeddings**  
  Stores the numeric representation of the text/data.

- **Index for fast similarity search**  
  Use specialized indexes (HNSW for MariaDB) for efficient nearest-neighbor queries.

### Pros and Cons of Single Table vs Multiple Tables for Embeddings

#### Single Table for Embeddings

**Pros:**

- ✅ **Simplicity:** All embeddings are in one place, easier to query and maintain.
- ✅ **Unified indexing:** You can create a single vector index for all embeddings, which can speed up similarity searches.
- ✅ **Easier analytics:** Aggregating statistics across all embeddings is straightforward.
- ✅ **Lower schema complexity:** No need to manage multiple tables, relationships, or joins.

**Cons:**

- ❌ **Scalability issues:** If you have embeddings from very different data sources,
  the table can become huge and harder to manage.
- ❌ **Heterogeneous data handling:** You might need extra columns to differentiate
  embeddings by type/source.
- ❌ **Locking/Write contention:** High write frequency from multiple sources
  can create contention in a single table.
- ❌ **Potential performance hit:** Searching across all embeddings may be slower if some queries only need a subset of data.

---

#### Multiple Tables for Embeddings

**Pros:**

- ✅ **Better organization:** Each table can represent a specific data type (e.g., students, teachers, books), making it clear and manageable.
- ✅ **Improved query performance:** Searches can target only the relevant table, reducing search space.
- ✅ **Flexible schema per table:** Different types of embeddings may have different metadata without wasting columns.
- ✅ **Reduced contention:** Writes and updates are separated by table, improving concurrency.

**Cons:**

- ❌ **More complex management:** Maintaining multiple tables, indexes, and migrations is harder.
- ❌ **Cross-table queries are harder:** Aggregating or comparing embeddings from multiple tables may require extra logic.
- ❌ **Duplicate infrastructure:** Each table may need its own indexing or maintenance operations, increasing storage or computation overhead.
- ❌ **Fragmentation:** Similar embeddings may be split across tables, complicating global analytics.
- ❌ **High RAM usage:** If multiple tables are indexed for vector search (e.g., HNSW),
  it can consume large amounts of memory for each index.

---

**Rule of thumb:**

- Use **single table** if your embeddings are homogeneous and you want simplicity
  and easy global search.
- Use **multiple tables** if your embeddings come from very different
  sources/types or require different metadata and frequent isolated queries.

**For the Ednevnik system, we will use a single table** for embeddings to enable
better global search, maintain simplicity across the system, and avoid excessive
RAM usage from multiple vector indexes.

---

## 7. Querying Embeddings

### 7.1 Filtering Using Metadata Before Similarity Search

Before performing a vector similarity search, you can filter rows based on metadata.
This reduces the search space and improves performance.

**Example:**

```sql
SELECT *
FROM embeddings
WHERE data_type = 'student_grades'
  AND school_id = 42;
```

### 7.2 Vector Similarity Search Syntax (MariaDB Example)

MariaDB supports similarity search on `VECTOR` columns using distance functions like `COSINE_DISTANCE` or `L2_DISTANCE`.

**Example:**

```sql
SELECT *, COSINE_DISTANCE(embedding, ?) AS distance
FROM embeddings
WHERE COSINE_DISTANCE(embedding, ?) < 0.3
ORDER BY distance;
```

- `?` is the query embedding vector
- `WHERE` clause filters results to only highly similar vectors
- Lower distance values indicate higher similarity

### 7.3 Retrieving Similarity-Based Results

Instead of limiting to a fixed number of results, use similarity thresholds to retrieve only semantically relevant embeddings.

**Example:**

```sql
SELECT *, COSINE_DISTANCE(embedding, ?) AS distance
FROM embeddings
WHERE COSINE_DISTANCE(embedding, ?) < 0.25
ORDER BY distance;
```

This retrieves all embeddings with cosine distance less than 0.25 (highly similar) to the query vector.

**Choosing Similarity Thresholds:**

- Depends on your embedding model and data characteristics
- For precise matches, use strict thresholds (e.g., < 0.15)
- For broader semantic retrieval, use more lenient thresholds (e.g., < 0.4)
- Test with representative queries to find optimal thresholds

**Hybrid Approach (Recommended):**

```sql
SELECT *, COSINE_DISTANCE(embedding, ?) AS distance
FROM embeddings
WHERE COSINE_DISTANCE(embedding, ?) < 0.3
ORDER BY distance
LIMIT 20; -- Safety limit to prevent excessive results
```

### 7.4 Similarity-Based Retrieval for Ednevnik

**Purpose of Similarity Thresholds:**
Similarity thresholds ensure the retriever only sends semantically relevant context to the LLM, improving answer quality while reducing noise from irrelevant information.

**Dynamic Context Coverage:**
Each embedding in Ednevnik covers comprehensive information about students per
subject, including grades, teacher comments, and performance data. Using
similarity thresholds (e.g., cosine distance < 0.3) allows the system to:

- Retrieve variable amounts of context based on query specificity
- Automatically adapt to different question types
- Avoid padding with irrelevant information

**Balance Between Relevance and Completeness:**

- **Strict thresholds (< 0.2):** High precision but may miss relevant context
  for complex queries
- **Lenient thresholds (< 0.4):** Better coverage but may include some less
  relevant information
- **Optimal range (0.25-0.35):** Balances semantic relevance with comprehensive
  context coverage

**Adaptive Retrieval Strategy:**

```sql
-- Primary query with strict threshold
SELECT *, COSINE_DISTANCE(embedding, ?) AS distance
FROM embeddings
WHERE COSINE_DISTANCE(embedding, ?) < 0.25
ORDER BY distance
LIMIT 15;

-- Fallback with relaxed threshold if insufficient results
-- (Implement in application logic)
```

**Performance Considerations:**

- Similarity-based retrieval scales naturally with query complexity
- Simple questions retrieve fewer, more focused results
- Complex questions automatically gather broader context
- Maximum limits prevent excessive context that could overwhelm the LLM

✅ **Conclusion:**
For Ednevnik's RAG chatbot, **similarity thresholds (cosine distance < 0.3)** with
a safety limit ensure the LLM receives semantically relevant context while
maintaining system efficiency. This approach provides more intelligent,
adaptive retrieval compared to fixed result counts.

---

## 8. Aggregating Context

When building a RAG (Retrieval-Augmented Generation) system, retrieving individual embedding rows is often not enough.  
To provide high-quality answers, we need to **aggregate multiple rows into a single context block** that can be passed to the LLM.

---

### 8.1. Combining Multiple Embedding Rows

- **Concatenate retrieved rows:**  
  Join the top-K retrieved embeddings (e.g., 5–20 rows) into a single text block.  
  Example: combine all grades and comments for a student across multiple subjects.
- **Preserve source attribution:**  
  Add metadata tags (e.g., `[Teacher: John | Subject: Math] Grade: 4`) so the LLM knows the origin of each piece.
- **Use separators:**  
  Clearly mark boundaries between rows, such as `---` or newlines, to avoid blending unrelated data.

---

### 8.2. Chunking or Truncating Context to Fit LLM Limits

- **Chunking:**  
  If the aggregated context is too large, split it into chunks (e.g., 500–1,000 tokens) and process each separately.
- **Truncation:**  
  Keep the most relevant top-K results and discard lower-ranked rows.
- **Summarization step:**  
  Use a smaller LLM or heuristic rules to summarize raw data before passing it
  to the main LLM.

---

### 8.3. Example of Aggregated Context

Suppose the embeddings table contains rows about students, subjects, and teacher comments.  
If a query asks _“How is Ana doing in school?”_, and we retrieve two rows, the
aggregated context might look like this:

```
[Student: Ana K. | Subject: Math | Teacher: John D.]
Grade: 5
Comment: Very good progress

[Student: Ana K. | Subject: English | Teacher: Mary P.]
Grade: 4
Comment: Needs improvement
```

This structured format makes it easy for the LLM to understand relationships while staying within token limits.

✅ **Best Practice:**

- Always structure aggregated context in a **consistent format**.
- Include only the **most relevant and recent data**.
- Use **metadata annotations** to help the LLM reason about sources and relationships.

---

## 8.4. Multimodal RAG Considerations (MariaDB Approach)

To support multimodal retrieval (text, images, and other data types) in the future using MariaDB, consider the following strategies:

- **Multimodal Embedding Models:**  
  Use models like CLIP or similar to generate embeddings for both text and images. Each content type (text, image, etc.) will have its own embedding vector.

- **Storing Multiple Embedding Types:**  
  Extend your MariaDB schema to include separate columns for each embedding type, for example:
  - `text_embedding VECTOR(...)`
  - `image_embedding VECTOR(...)`
  - Add a `content_type` or `modality` column to indicate the type of data (e.g., "text", "image", "audio").

- **Mixed Documents:**  
  When processing documents that contain both text and images, extract and embed each part separately. Store each embedding in the same row (if using a single table) or in separate rows/tables, and use metadata to link them (e.g., `document_id`, `section_id`, `content_type`).

- **Schema Design Options:**  
  - **Single Table:**  
    Add columns for each embedding type and use metadata to indicate which columns are populated for each row.
  - **Multiple Tables:**  
    Create separate tables for each modality (e.g., `text_embeddings`, `image_embeddings`) with shared metadata fields for linking and filtering.

- **Querying Multimodal Embeddings:**  
  When searching, select the appropriate embedding column and distance function based on the query type (text or image). Use metadata filters to narrow down results by modality.

- **Example Schema (Single Table):**
    ```sql
    CREATE TABLE multimodal_embeddings (
        id INT PRIMARY KEY,
        document_id INT,
        content_type VARCHAR(16), -- 'text', 'image', etc.
        text_embedding VECTOR(384),
        image_embedding VECTOR(512),
        metadata JSON,
        -- Add vector indexes as needed
        VECTOR INDEX (text_embedding) M=6 DISTANCE=cosine,
        VECTOR INDEX (image_embedding) M=6 DISTANCE=cosine
    );
    ```

- **Metadata for Linking:**  
  Use metadata fields (e.g., `document_id`, `section`, `caption`) to relate different modalities from the same source document.

**Summary:**  
By extending your schema and indexing strategy, MariaDB can support multimodal RAG in the future. This enables semantic search and retrieval across both text and images, with flexible filtering and aggregation using

## 9. Integrating LangChain in the RAG Pipeline

**LangChain** is a popular open-source framework for building applications with large language models (LLMs) and retrieval-augmented generation (RAG) workflows. It provides modular components for chaining together LLMs, retrievers, memory, and tools, making it easier to build, maintain, and extend RAG chatbots.

---

### 9.1. Why Use LangChain?

- **Abstraction:** Simplifies integration of LLMs, vector stores, and prompt templates.
- **Flexibility:** Supports multiple LLM providers (OpenAI, HuggingFace, Gemini, etc.) and vector databases (MariaDB, Pinecone, Weaviate, etc.).
- **Composability:** Enables chaining of retrieval, context aggregation, and LLM calls.
- **Extensibility:** Easy to add custom logic, filters, or tools.

---

### 9.2. LangChain Components in Our Setup

For the Ednevnik RAG chatbot, we would use LangChain to orchestrate the following steps:

1. **Prompt Templates:** Define system and user prompts, including context formatting.
2. **Embeddings:** Use LangChain’s embedding wrappers to generate query and document embeddings (e.g., via HuggingFace or OpenAI).
3. **Vector Store Integration:** Connect LangChain to MariaDB Vector as the retrieval backend using a custom or community-supported vector store class.
4. **Retriever:** Use LangChain’s retriever abstraction to filter by metadata and perform similarity search.
5. **Context Aggregation:** Aggregate and format retrieved rows for the LLM.
6. **LLM Chain:** Pass the user query and aggregated context to the LLM for answer generation.
7. **Output Handling:** Return the LLM’s response to the user.

---

### 9.3. Example LangChain Workflow (Python Pseudocode)

```python
# Example: LangChain RAG workflow with memory (conversation history)
from langchain.embeddings import HuggingFaceEmbeddings
from langchain.vectorstores import MariaDBVectorStore  # Custom or community implementation
from langchain.llms import OpenAI
from langchain.prompts import PromptTemplate
from langchain.chains import ConversationalRetrievalChain
from langchain.memory import ConversationBufferMemory

# 1. Embedding model
embeddings = HuggingFaceEmbeddings(model_name="sentence-transformers/all-MiniLM-L6-v2")

# 2. Vector store (MariaDB Vector)
vector_store = MariaDBVectorStore(
    connection_params={...},
    embedding_function=embeddings,
    table_name="embeddings"
)

# 3. Retriever with metadata filtering
retriever = vector_store.as_retriever(
    search_type="similarity",
    search_kwargs={
        "filter": {"account_type": "teacher", "subject": "math"},
        "distance_threshold": 0.3,
        "k": 15
    }
)

# 4. Prompt template
prompt = PromptTemplate(
    template="""
    {system_prompt}
    Context:
    {context}
    Conversation history:
    {chat_history}
    User question: {question}
    """,
    input_variables=["system_prompt", "context", "chat_history", "question"]
)

# 5. LLM
llm = OpenAI(model="gpt-3.5-turbo", temperature=0.7)

# 6. Memory for conversation history
memory = ConversationBufferMemory(memory_key="chat_history", return_messages=True)

# 7. Conversational RetrievalQA chain
qa_chain = ConversationalRetrievalChain(
    retriever=retriever,
    llm=llm,
    memory=memory,
    prompt=prompt
)

# 8. Run the chain with conversation
response = qa_chain.run({
    "system_prompt": "You are an educational assistant...",
    "question": "Kako je Sara napredovala u matematici?",
})
print(response)
```

## 10. Building the RAG Chatbot

**Workflow Overview:**

1. Receive user query
2. Generate query embedding
3. Retrieve relevant embeddings from database
4. Aggregate multiple rows into context
5. Send query + context to LLM
6. Return LLM response to user


Za LLM tabela ogranicenja u smisulu finansija i tokena i ako imaju druga ogranicenje
LLAMA Bert OpenAI Gemini Mistral TrustLLM (EU)

Dijagram workflowa


https://mariadb.com/resources/blog/how-fast-is-mariadb-vector/
Staviti sta koje baze postoje i ova poredjenja

Staviti i vertex gdje je porednjenje mariadb i weviate, Yandex DB

Airtable

Fali sa odabir LLM

Zasto BERT ne moze generisati tekst a LLM moze

Dodati literaturu za AI i updateovati

Fali transformers

Opisati langchain driver

REY

https://github.com/ray-project/langchain-ray

https://python.langchain.com/docs/integrations/providers/ray_serve/

Hugging face vs LM studio
