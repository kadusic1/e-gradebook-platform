from tenant_helper import get_tenant_databases
from general_helper import get_success_category


def collect_section_embeddings(cursor):
    """Collect section embeddings from all tenant databases"""
    tenant_ids = get_tenant_databases(cursor)

    data_to_insert = []

    for tenant_id in tenant_ids:
        tenant_db_name = f"ednevnik_tenant_db_tenant_id_{tenant_id}"

        # Get section data with curriculum and tenant info (including archived sections)
        query = f""" 
        SELECT s.id, s.section_code, s.class_code, s.year, 
               c.curriculum_name, t.tenant_name, s.archived 
        FROM {tenant_db_name}.sections s 
        JOIN ednevnik_workspace.curriculum c ON s.curriculum_code = c.curriculum_code 
        JOIN ednevnik_workspace.tenant t ON s.tenant_id = t.id 
        """

        cursor.execute(query)
        sections = cursor.fetchall()

        for section in sections:
            (
                section_id,
                section_code,
                class_code,
                year,
                curriculum_name,
                tenant_name,
                archived,
            ) = section

            # Calculate section average grade
            grade_query = f""" 
            SELECT AVG(grade) as avg_grade 
            FROM {tenant_db_name}.student_grades sg 
            WHERE sg.section_id = ? AND sg.type IN ('exam', 'oral', 'written_assignment') 
            AND sg.grade IS NOT NULL 
            """
            cursor.execute(grade_query, (section_id,))
            result = cursor.fetchone()
            avg_grade = result[0] if result and result[0] else 3.0

            success_category = get_success_category(avg_grade)
            section_name = f"{class_code}-{section_code}"

            # Determine archive status text
            archive_status = "Da" if archived else "Ne"
            archive_text = " (Odjeljenje je arhivirano)" if archived else ""

            metadata = {
                "vrsta": "Opšte informacije o odjeljenju",
                "ime": section_name,
                "institucija": tenant_name,
                "godina": year,
                "naziv_kurikuluma": curriculum_name,
                "uspjeh": success_category,
                "arhivirano": archive_status,
                # Filtering fields
                "source": "section",
                "tenant_id": tenant_id,
            }

            sentence = f"Odjeljenje: {section_name} - {tenant_name} u školskoj godini {year}, koristi kurikulum: {curriculum_name}. Uspjeh odjeljenja {section_name} je {success_category}{archive_text}. Tip entiteta: odjeljenje."

            data_to_insert.append(
                {
                    "metadata": metadata,
                    "sentence": sentence,
                }
            )

    return data_to_insert
