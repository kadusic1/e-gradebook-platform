from tenant_helper import get_tenant_databases


def collect_grade_embeddings(cursor):
    """Collect grade embeddings from all tenant databases"""
    tenant_ids = get_tenant_databases(cursor)

    data_to_insert = []

    for tenant_id in tenant_ids:
        tenant_db_name = f"ednevnik_tenant_db_tenant_id_{tenant_id}"

        # Get grade data grouped by pupil, section, and subject
        query = f"""
        SELECT sg.pupil_id, p.name, p.last_name, s.section_code, s.class_code,
               subj.subject_name, t.tenant_name, cs.course_name,
               GROUP_CONCAT(sg.grade ORDER BY sg.grade_date) as grades,
               GROUP_CONCAT(sg.signature ORDER BY sg.grade_date) as signatures,
               s.archived, p.account_id
        FROM {tenant_db_name}.student_grades sg
        JOIN {tenant_db_name}.pupils p ON sg.pupil_id = p.id
        JOIN {tenant_db_name}.sections s ON sg.section_id = s.id
        JOIN ednevnik_workspace.subjects subj ON sg.subject_code = subj.subject_code
        JOIN ednevnik_workspace.tenant t ON s.tenant_id = t.id
        LEFT JOIN ednevnik_workspace.curriculum c ON s.curriculum_code = c.curriculum_code
        LEFT JOIN ednevnik_workspace.courses_secondary cs ON c.course_code = cs.course_code
        WHERE sg.type IN ('exam', 'oral', 'written_assignment') AND sg.grade IS NOT NULL
        GROUP BY sg.pupil_id, sg.section_id, sg.subject_code
        """

        cursor.execute(query)
        grade_groups = cursor.fetchall()

        for grade_group in grade_groups:
            (
                pupil_id,
                name,
                last_name,
                section_code,
                class_code,
                subject_name,
                tenant_name,
                course_name,
                grades_str,
                signatures_str,
                archived,
                account_id,
            ) = grade_group

            grades_list = [int(g) for g in grades_str.split(",") if g.strip()]
            signatures_list = [
                s.strip() for s in signatures_str.split(",") if s.strip()
            ]
            grades_with_signature = ", ".join(
                [f"{g} ({s})" for g, s in zip(grades_list, signatures_list)]
            )
            section_name = f"{class_code}-{section_code}"
            archive_status = "Da" if archived else "Ne"
            archive_text = " (Odjeljenje je arhivirano)" if archived else ""

            metadata = {
                "vrsta": "Ocjene učenika za određeni predmet u određenom odjeljenju",
                "institucija": tenant_name,
                "odjeljenje": section_name,
                "ime_ucenika": name,
                "prezime_ucenik": last_name,
                "predmet": subject_name,
                "ocjene": grades_list,
                "potpisi": signatures_list,
                "arhivirano": archive_status,
                # Filtering fields
                "source": "grade",
                "tenant_id": tenant_id,
                "account_id": account_id,
            }

            if course_name:
                metadata["smjer"] = course_name
                sentence = f"Ocjene: učenika {name} {last_name} iz ({section_name}{archive_text} - {tenant_name}) smjera {course_name} za predmet {subject_name} su: {grades_with_signature}. Tip entiteta: ocjena."
            else:
                sentence = f"Ocjene: učenika {name} {last_name} (odjeljenje: {section_name}{archive_text} - {tenant_name}) za predmet {subject_name} su: {grades_with_signature}. Tip entiteta: ocjena."

            data_to_insert.append(
                {
                    "metadata": metadata,
                    "sentence": sentence,
                }
            )

    return data_to_insert
