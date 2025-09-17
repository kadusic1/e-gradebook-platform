from tenant_helper import get_tenant_databases


def collect_behaviour_embeddings(cursor):
    """Collect behaviour embeddings from all tenant databases"""
    tenant_ids = get_tenant_databases(cursor)

    data_to_insert = []

    for tenant_id in tenant_ids:
        tenant_db_name = f"ednevnik_tenant_db_tenant_id_{tenant_id}"

        # Get behaviour data
        query = f"""
        SELECT pb.pupil_id, p.name, p.last_name, s.section_code, s.class_code,
               pb.behaviour, t.tenant_name, cs.course_name, s.archived, p.account_id
        FROM {tenant_db_name}.pupil_behaviour pb
        JOIN {tenant_db_name}.pupils p ON pb.pupil_id = p.id
        JOIN {tenant_db_name}.sections s ON pb.section_id = s.id
        JOIN ednevnik_workspace.tenant t ON s.tenant_id = t.id
        LEFT JOIN ednevnik_workspace.curriculum c ON s.curriculum_code = c.curriculum_code
        LEFT JOIN ednevnik_workspace.courses_secondary cs ON c.course_code = cs.course_code
        """

        cursor.execute(query)
        behaviours = cursor.fetchall()

        for behaviour in behaviours:
            (
                pupil_id,
                name,
                last_name,
                section_code,
                class_code,
                behaviour_val,
                tenant_name,
                course_name,
                archived,
                account_id,
            ) = behaviour

            section_name = f"{class_code}-{section_code}"
            archive_status = "Da" if archived else "Ne"
            archive_text = " (Odjeljenje je arhivirano)" if archived else ""

            # Map behaviour
            behaviour_map = {
                "primjerno": "Primjereno",
                "vrlodobro": "Vrlo dobro",
                "dobro": "Dobro",
                "zadovoljavajuće": "Zadovoljavajuće",
                "loše": "Loše",
            }
            behaviour_bosnian = behaviour_map.get(behaviour_val, "Primjereno")

            metadata = {
                "vrsta": "Vladanje učenika u određenom odjeljenju",
                "institucija": tenant_name,
                "odjeljenje": section_name,
                "ime_ucenika": name,
                "prezime_ucenik": last_name,
                "vladanje": behaviour_bosnian,
                "arhivirano": archive_status,
                # Filtering fields
                "source": "behaviour",
                "tenant_id": tenant_id,
                "account_id": account_id,
            }

            if course_name:
                metadata["smjer"] = course_name
                sentence = f"Vladanje učenika: {name} {last_name} ({section_name}{archive_text} - {tenant_name} smjer: {course_name}) je: {behaviour_bosnian}. Tip entiteta: vladanje."
            else:
                sentence = f"Vladanje učenika: {name} {last_name} ({section_name}{archive_text} - {tenant_name}) je: {behaviour_bosnian}. Tip entiteta: vladanje."

            data_to_insert.append(
                {
                    "metadata": metadata,
                    "sentence": sentence,
                }
            )

    return data_to_insert
