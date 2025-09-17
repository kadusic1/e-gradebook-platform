from typing import List
from general_helper import get_difficulty_for_tenant


def get_tenant_databases(cursor) -> List[int]:
    """Get list of all tenant IDs from workspace database"""
    cursor.execute("SELECT id FROM ednevnik_workspace.tenant")
    return [row[0] for row in cursor.fetchall()]


def collect_tenant_embeddings(cursor):
    """Collect tenant institution embeddings"""
    # Get tenant data with average grades for difficulty calculation
    query = """
    SELECT t.id, t.tenant_name, t.tenant_city, t.tenant_type, 
           c.canton_name, t.phone, t.email, t.director_name,
           t.longitude, t.latitude, t.domain, t.specialization
    FROM ednevnik_workspace.tenant t
    LEFT JOIN ednevnik_workspace.cantons c ON t.canton_code = c.canton_code
    """

    cursor.execute(query)
    tenant_data = cursor.fetchall()

    data_to_insert = []

    for tenant in tenant_data:
        (
            tenant_id,
            name,
            city,
            school_type,
            canton,
            phone,
            email,
            director,
            longitude,
            latitude,
            domain,
            specialization,
        ) = tenant

        # Calculate average grades across all tenant databases for this tenant
        tenant_db_name = f"ednevnik_tenant_db_tenant_id_{tenant_id}"

        # Get average grade for this tenant
        grade_query = f"""
        SELECT AVG(grade) as avg_grade 
        FROM {tenant_db_name}.student_grades 
        WHERE type IN ('exam', 'oral', 'written_assignment') 
        AND grade IS NOT NULL
        """
        cursor.execute(grade_query)
        result = cursor.fetchone()
        avg_grade = result[0] if result and result[0] else 3.0

        difficulty_category = get_difficulty_for_tenant(avg_grade)

        # Map specialization
        specialization_map = {
            "regular": "obična",
            "religion": "vjerska",
            "musical": "muzička",
        }
        spec_bosnian = specialization_map.get(specialization, "obična")

        # Map school type
        type_map = {"primary": "Osnovna škola", "secondary": "Srednja škola"}
        type_bosnian = type_map.get(school_type, "Osnovna škola")

        metadata = {
            "vrsta": "Opšte informacije o instituciji",
            "ime": name,
            "grad": city or "",
            "tip": type_bosnian,
            "kanton": canton or "",
            "telefon": phone or "",
            "email": email or "",
            "direktor": director or "",
            "longitude": str(longitude) if longitude else "",
            "latitude": str(latitude) if latitude else "",
            "domena": domain or "",
            "specijalizacija": spec_bosnian,
            "težina": difficulty_category,
            # Filtering fields
            "source": "tenant",
            "tenant_id": tenant_id,
        }

        sentence = f"Institucija (škola) {name}: se nalazi u gradu {city or ''}, tip institucije (škole): {type_bosnian}, u kantonu {canton or ''}. Direktor institucije (škole) {name} je {director or ''}. Kontakt institucije (škole) {name}, telefon: {phone or ''}, email: {email or ''}. Web domena institucije (škole) {name}: {domain or ''}. Specijalizacija institucije (škole) {name}: {spec_bosnian}. Geografska širina institucije (škole) {name}: {latitude or ''}, geografska dužina institucije (škole) {name}: {longitude or ''}. Težina institucije (škole) {name} je {difficulty_category}. Tip entiteta: škola."

        data_to_insert.append(
            {
                "metadata": metadata,
                "sentence": sentence,
            }
        )

    return data_to_insert
