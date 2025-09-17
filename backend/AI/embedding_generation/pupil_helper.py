from tenant_helper import get_tenant_databases
from general_helper import get_success_category


def collect_pupil_embeddings(cursor):
    """Collect pupil embeddings from all tenant databases"""

    # Get pupil data
    query = """
    SELECT p.id, p.name, p.last_name, p.gender, p.address, p.guardian_name,
           p.phone_number, p.guardian_number, p.date_of_birth, p.place_of_birth,
           p.religion, p.account_id, a.email
    FROM ednevnik_workspace.pupil_global p
    LEFT JOIN ednevnik_workspace.accounts a ON p.account_id = a.id
    """

    cursor.execute(query)
    pupils = cursor.fetchall()

    data_to_insert = []
    tenant_ids = get_tenant_databases(cursor)

    for pupil in pupils:
        (
            pupil_id,
            name,
            last_name,
            gender,
            address,
            guardian_name,
            phone,
            guardian_phone,
            birth_date,
            birth_place,
            religion,
            account_id,
            email,
        ) = pupil

        # Calculate pupil's average grades across all tenants
        all_grades = []
        pupil_sections = []
        active_sections = []
        archived_sections = []

        for tenant_id in tenant_ids:
            tenant_db_name = f"ednevnik_tenant_db_tenant_id_{tenant_id}"

            # Check if pupil exists in this tenant
            check_query = f"SELECT 1 FROM {tenant_db_name}.pupils WHERE id = ?"
            cursor.execute(check_query, (pupil_id,))
            if not cursor.fetchone():
                continue

            # Get grades
            grade_query = f"""
            SELECT grade FROM {tenant_db_name}.student_grades sg
            WHERE sg.pupil_id = ? AND sg.grade IS NOT NULL
            AND sg.type IN ('exam', 'oral', 'written_assignment')
            """
            cursor.execute(grade_query, (pupil_id,))
            grades = [row[0] for row in cursor.fetchall()]
            all_grades.extend(grades)

            # Get sections and courses for this pupil with behavior and archive status
            section_query = f"""
            SELECT DISTINCT s.section_code, s.class_code, t.tenant_name, cs.course_name, s.archived,
                   pb.behaviour, ps.section_id
            FROM {tenant_db_name}.pupils_sections ps
            JOIN {tenant_db_name}.sections s ON ps.section_id = s.id
            JOIN ednevnik_workspace.tenant t ON s.tenant_id = t.id
            LEFT JOIN ednevnik_workspace.curriculum c ON s.curriculum_code = c.curriculum_code
            LEFT JOIN ednevnik_workspace.courses_secondary cs ON c.course_code = cs.course_code
            LEFT JOIN {tenant_db_name}.pupil_behaviour pb ON pb.pupil_id = ps.pupil_id AND pb.section_id = ps.section_id
            WHERE ps.pupil_id = ? AND ps.is_active = TRUE
            """
            cursor.execute(section_query, (pupil_id,))
            sections = cursor.fetchall()

            for section in sections:
                (
                    section_code,
                    class_code,
                    tenant_name,
                    course_name,
                    archived,
                    behavior,
                    section_id,
                ) = section

                # Map behavior
                behavior_map = {
                    "primjerno": "Primjerno",
                    "vrlodobro": "Vrlo dobro",
                    "dobro": "Dobro",
                    "zadovoljavajuće": "Zadovoljavajuće",
                    "loše": "Loše",
                }
                behavior_bosnian = (
                    behavior_map.get(behavior, "Primjerno") if behavior else "Primjerno"
                )

                section_info = {
                    "kod": f"{class_code}-{section_code}",
                    "institucija": tenant_name,
                    "ponašanje": behavior_bosnian,
                    "arhivirano": "Da" if archived else "Ne",
                }

                if course_name:
                    section_info["smjer"] = course_name

                pupil_sections.append(section_info)

                # Separate active and archived sections for sentences
                section_display = f"{class_code}-{section_code} ({tenant_name})"
                if course_name:
                    section_display += f" - {course_name}"

                if archived:
                    archived_sections.append(section_display)
                else:
                    active_sections.append(section_display)

        if not pupil_sections:  # Skip pupils not enrolled in any sections
            continue

        avg_grade = sum(all_grades) / len(all_grades) if all_grades else 3.0
        success_category = get_success_category(avg_grade)

        # Map gender and religion
        gender_map = {"M": "muškog", "F": "ženskog"}
        gender_bosnian = gender_map.get(gender, "muškog")

        religion_map = {
            "Islam": "Islamska",
            "Catholic": "Katolička",
            "Orthodox": "Pravoslavna",
            "Jewish": "Jevrejska",
            "Other": "Ostala",
            "NotAttendingReligion": "Ne pohađa vjeronauku",
        }
        religion_bosnian = religion_map.get(religion, "Islamska")

        tenant_ids_query = """
        SELECT tenant_id FROM ednevnik_workspace.pupil_tenant
        WHERE pupil_id = ?
        """
        cursor.execute(tenant_ids_query, (pupil_id,))
        available_in_tenant_ids = [row[0] for row in cursor.fetchall()]

        metadata = {
            "vrsta": "Opšte informacije o učeniku",
            "ime": name,
            "prezime": last_name,
            "spol": gender_bosnian.title(),
            "adresa": address or "",
            "staratelj": guardian_name,
            "email": email or "",
            "telefon": phone or "",
            "telefon_staratelja": guardian_phone or "",
            "datum_rodjenja": str(birth_date) if birth_date else "",
            "mjesto_rodjenja": birth_place or "",
            "vjeroispovijest": religion_bosnian,
            "uspjeh": success_category,
            "odjeljenja": pupil_sections,
            # Filtering fields
            "source": "pupil",
            "available_in_tenant_ids": available_in_tenant_ids,
            "account_id": account_id,
        }

        # Build base sentence
        sentence = f"Učenik: {name} {last_name}, {gender_bosnian} spola, rođen {birth_date or ''} u {birth_place or ''}, živi na adresi {address or ''}. Staratelj učenika {name} {last_name} je {guardian_name}, kontakt: {guardian_phone or ''}. Učenik {name} {last_name} ima telefon: {phone or ''} i email: {email or ''}. Vjeroispovijest učenika {name} {last_name} je {religion_bosnian}. Uspjeh učenika {name} {last_name} je {success_category}."

        # Add sentences for current and archived sections
        if active_sections:
            sentence += (
                f" Učenik pohađa slijedeća odjeljenja: {', '.join(active_sections)}."
            )

        if archived_sections:
            sentence += f" Učenik je pohađao slijedeća odjeljenja: {', '.join(archived_sections)}."

        sentence += " Tip entiteta: učenik."

        data_to_insert.append(
            {
                "metadata": metadata,
                "sentence": sentence,
            }
        )

    return data_to_insert
