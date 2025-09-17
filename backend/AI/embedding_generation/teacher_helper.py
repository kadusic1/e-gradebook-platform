from tenant_helper import get_tenant_databases
from general_helper import get_criteria_for_teacher


def collect_teacher_embeddings(cursor):
    """Collect teacher embeddings from workspace database"""

    # Get teacher data with account type and tenant admin info
    query = """
    SELECT t.id, t.name, t.last_name, t.phone, t.contractions, t.title, t.account_id,
           a.email, a.account_type,
           GROUP_CONCAT(DISTINCT ten.tenant_name SEPARATOR ', ') as admin_institutions
    FROM ednevnik_workspace.teachers t
    LEFT JOIN ednevnik_workspace.accounts a ON t.account_id = a.id
    LEFT JOIN ednevnik_workspace.tenant ten ON t.id = ten.tenant_admin_id
    GROUP BY t.id, t.name, t.last_name, t.phone, t.contractions, t.title, t.account_id, a.email, a.account_type
    """

    cursor.execute(query)
    teachers = cursor.fetchall()

    data_to_insert = []

    for teacher in teachers:
        (
            teacher_id,
            name,
            last_name,
            phone,
            contractions,
            title,
            account_id,
            email,
            account_type,
            admin_institutions,
        ) = teacher

        # Calculate teacher's average grades across all tenants
        tenant_ids = get_tenant_databases(cursor)
        all_grades = []
        odjeljenja = []

        for tenant_id in tenant_ids:
            tenant_db_name = f"ednevnik_tenant_db_tenant_id_{tenant_id}"

            # Get grades for criteria calculation
            grade_query = f"""
            SELECT grade FROM {tenant_db_name}.student_grades 
            WHERE teacher_id = ? AND grade IS NOT NULL
            AND type IN ('exam', 'oral', 'written_assignment')
            """
            cursor.execute(grade_query, (teacher_id,))
            grades = [row[0] for row in cursor.fetchall()]
            all_grades.extend(grades)

            # Get sections (odjeljenja) where this teacher teaches
            sections_query = f"""
            SELECT DISTINCT s.section_code, s.class_code, ten.tenant_name, 
                   c.course_name, s.id as section_id, s.archived
            FROM {tenant_db_name}.teachers_sections_subjects tss
            JOIN {tenant_db_name}.sections s ON tss.section_id = s.id
            JOIN ednevnik_workspace.tenant ten ON s.tenant_id = ten.id
            LEFT JOIN ednevnik_workspace.curriculum cur ON s.curriculum_code = cur.curriculum_code
            LEFT JOIN ednevnik_workspace.courses_secondary c ON cur.course_code = c.course_code
            WHERE tss.teacher_id = ?
            """
            cursor.execute(sections_query, (teacher_id,))
            sections = cursor.fetchall()

            for section in sections:
                (
                    section_code,
                    class_code,
                    tenant_name,
                    course_name,
                    section_id,
                    archived,
                ) = section

                # Check if this teacher is homeroom teacher for this section
                homeroom_query = f"""
                SELECT COUNT(*) FROM {tenant_db_name}.homeroom_assignments 
                WHERE section_id = ? AND teacher_id = ?
                """
                cursor.execute(homeroom_query, (section_id, teacher_id))
                is_homeroom = cursor.fetchone()[0] > 0

                odjeljenje_data = {
                    "kod": f"{class_code}-{section_code}",
                    "institucija": tenant_name,
                    "razrednik": "Da" if is_homeroom else "Ne",
                    "arhivirano": "Da" if archived else "Ne",
                }
                # Only add course/smjer for secondary schools
                if course_name:
                    odjeljenje_data["smjer"] = course_name

                odjeljenja.append(odjeljenje_data)

        avg_grade = sum(all_grades) / len(all_grades) if all_grades else 3.0
        teacher_criteria = get_criteria_for_teacher(avg_grade)

        tenant_ids_query = """
        SELECT tenant_id FROM ednevnik_workspace.teacher_tenant
        WHERE teacher_id = ?
        """
        cursor.execute(tenant_ids_query, (teacher_id,))
        available_in_tenant_ids = [row[0] for row in cursor.fetchall()]

        metadata = {
            "vrsta": "Opšte informacije o Nastavnik/Profesoru",
            "ime": name,
            "prezime": last_name,
            "telefon": phone or "",
            "email": email or "",
            "oslovljavanje": contractions or "",
            "titula": title or "",
            "kriterij": teacher_criteria,
            "odjeljenja": odjeljenja,
            # Filtering fields
            "source": "teacher",
            "account_id": account_id,
            "available_in_tenant_ids": available_in_tenant_ids,
        }

        # Build sentence with sections information
        sections_text = ""
        archived_text = ""
        if odjeljenja:
            sections_parts = []
            archived_parts = []
            for odj in odjeljenja:
                section_info = f"{odj['kod']} - {odj['institucija']}"
                if "smjer" in odj:
                    section_info += f" smjera {odj['smjer']}"
                if odj.get("Razrednik") == "Da":
                    section_info += " (razrednik)"
                if odj.get("arhivirano") == "Da":
                    archived_parts.append(section_info)
                else:
                    sections_parts.append(section_info)

            if len(archived_parts) > 0:
                archived_text = f" Nastavnik/Profesor: {name} {last_name} je predavao slijedećim odjeljenjima: {', '.join(sections_parts)}."
            if len(sections_parts) > 0:
                sections_text = f" Nastavnik/Profesor: {name} {last_name} predaje u slijedećim odjeljenjima: {', '.join(sections_parts)}."

        # Add role information
        role_text = ""
        if account_type == "root":
            role_text = f" Nastavnik/Profesor: {name} {last_name} je superadmin."
        elif account_type == "tenant_admin" and admin_institutions:
            role_text = f" Nastavnik/Profesor: {name} {last_name} je tenant admin za instituciju {admin_institutions}."

        parts = [
            f"Nastavnik/Profesor: {name} {last_name}, titula {title or ''}, oslovljavanje {contractions or ''}, kontakt telefon: {phone or ''}, email: {email or ''}.",
            f"Kriterij nastavnika/profesora {name} {last_name} je {teacher_criteria}.",
            "Tip entiteta: nastavnik/profesor.",
            sections_text if sections_text else "",
            archived_text if archived_text else "",
            role_text if role_text else "",
        ]

        # Join only non-empty parts with spaces
        sentence = " ".join([p for p in parts if p])

        data_to_insert.append(
            {
                "metadata": metadata,
                "sentence": sentence,
            }
        )

    return data_to_insert
