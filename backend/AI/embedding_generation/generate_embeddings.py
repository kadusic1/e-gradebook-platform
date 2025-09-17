import mariadb
import json
from tenant_helper import (
    collect_tenant_embeddings,
)
from section_helper import collect_section_embeddings
from teacher_helper import collect_teacher_embeddings
from pupil_helper import collect_pupil_embeddings
from grade_helper import collect_grade_embeddings
from behaviour_helper import collect_behaviour_embeddings
import json
from general_helper import safe_float_str


def main():
    # Database connection configuration
    conn = mariadb.connect(
        user="eacon",
        password="test1234",
        host="localhost",
        database="ednevnik_workspace",
        autocommit=False,
    )
    cursor = conn.cursor()

    # Clear existing embeddings
    cursor.execute("DELETE FROM embeddings")
    conn.commit()

    # Collect all embedding data
    all_data = []

    print("Collecting tenant embeddings...")
    all_data.extend(collect_tenant_embeddings(cursor))

    print("Collecting section embeddings...")
    all_data.extend(collect_section_embeddings(cursor))

    print("Collecting teacher embeddings...")
    all_data.extend(collect_teacher_embeddings(cursor))

    print("Collecting pupil embeddings...")
    all_data.extend(collect_pupil_embeddings(cursor))

    print("Collecting grade embeddings...")
    all_data.extend(collect_grade_embeddings(cursor))

    print("Collecting behaviour embeddings...")
    all_data.extend(collect_behaviour_embeddings(cursor))

    print(f"Total embeddings to insert: {len(all_data)}")

    print("Loading embedding model...")
    from sentence_transformers import SentenceTransformer

    model = SentenceTransformer("LaBSE")
    print("Model loaded successfully.")

    print("Generating and inserting embeddings in batches...")

    # Process in batches for both embedding generation and database insertion
    batch_size = 100
    total_batches = (len(all_data) + batch_size - 1) // batch_size

    for i in range(0, len(all_data), batch_size):
        batch = all_data[i : i + batch_size]
        batch_num = i // batch_size + 1

        print(f"Processing batch {batch_num}/{total_batches}...")

        # Generate embeddings for entire batch at once
        sentences = [entry["sentence"] for entry in batch]
        embeddings = model.encode(
            sentences, convert_to_numpy=True, show_progress_bar=True
        )

        # Prepare batch data for insertion
        batch_insert_data = []
        for j, entry in enumerate(batch):
            vec_text = "[" + ",".join(safe_float_str(x) for x in embeddings[j]) + "]"
            batch_insert_data.append(
                (
                    json.dumps(entry["metadata"], ensure_ascii=False),
                    entry["sentence"],
                    vec_text,
                )
            )

        # Batch insert into database
        cursor.executemany(
            "INSERT INTO embeddings (metadata, content, embedding) VALUES (?, ?, VEC_FROMTEXT(?))",
            batch_insert_data,
        )

        # Commit after each batch to avoid memory issues
        conn.commit()

        print(f"Completed batch {batch_num}/{total_batches}")

    print("All embeddings inserted successfully!")
    cursor.close()
    conn.close()


if __name__ == "__main__":
    main()
