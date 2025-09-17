from script_util import script_setup
import subprocess
import random


def main():
    workspace_db_connection_cmd = script_setup()

    # Tenant choice
    print("Chose a tenant by ID")
    tenant_find_command = workspace_db_connection_cmd + [
        "-e",
        "SELECT id, tenant_name FROM tenant;",
    ]
    print("-" * 40)
    result = subprocess.run(
        tenant_find_command, capture_output=True, text=True, encoding="utf-8"
    )
    if not result.stdout.strip():
        print("[ERROR] No tenant found.")
        return
    print(result.stdout.strip())
    print("-" * 40)

    selected_tenant_id = input("Enter tenant ID: ").strip()

    check_tenant_command = workspace_db_connection_cmd + [
        "-e",
        f"SELECT id FROM tenant WHERE id = {selected_tenant_id};",
    ]
    result = subprocess.run(
        check_tenant_command, capture_output=True, text=True, encoding="utf-8"
    )
    if result.returncode != 0 or not result.stdout.strip():
        print(f'[ERROR] Tenant ID "{selected_tenant_id}" is not valid.')
        return

    tenant_db_connection_cmd = workspace_db_connection_cmd.copy()
    tenant_db_connection_cmd[-1] = f"ednevnik_tenant_db_tenant_id_{selected_tenant_id}"

    # Section choice
    print("Chose a section by ID")
    print("-" * 40)
    section_find_command = tenant_db_connection_cmd + [
        "-e",
        "SELECT id, section_code, class_code, year, curriculum_code FROM sections\
        WHERE archived = 0",
    ]
    result = subprocess.run(
        section_find_command, capture_output=True, text=True, encoding="utf-8"
    )
    if not result.stdout.strip():
        print("[ERROR] No section found.")
        return
    print(result.stdout.strip())
    print("-" * 40)

    selected_section_id = input("Enter section ID: ").strip()
    check_section_command = tenant_db_connection_cmd + [
        "-e",
        f"SELECT id FROM sections WHERE id = {selected_section_id};",
    ]
    result = subprocess.run(
        check_section_command, capture_output=True, text=True, encoding="utf-8"
    )
    if result.returncode != 0 or not result.stdout.strip():
        print(f'[ERROR] Section ID "{selected_section_id}" is not valid.')
        return

    # Pupil choice
    print("Chose a pupil by ID")
    print("-" * 40)
    pupil_find_command = tenant_db_connection_cmd + [
        "-e",
        f"SELECT p.id, p.name, p.last_name FROM pupils p\
        JOIN pupils_sections ps ON ps.pupil_id = p.id\
        WHERE ps.section_id = {selected_section_id} AND ps.is_active = 1;",
    ]
    result = subprocess.run(
        pupil_find_command, capture_output=True, text=True, encoding="utf-8"
    )
    if not result.stdout.strip():
        print("[ERROR] No pupil found.")
        return
    print(result.stdout.strip())
    print("-" * 40)

    selected_pupil_id = input("Enter pupil ID: ").strip()
    check_pupil_command = tenant_db_connection_cmd + [
        "-e",
        f"SELECT p.id FROM pupils p\
        JOIN pupils_sections ps ON ps.pupil_id = p.id\
        WHERE ps.section_id = {selected_section_id} AND p.id = {selected_pupil_id}\
        AND ps.is_active = 1;",
    ]
    result = subprocess.run(
        check_pupil_command, capture_output=True, text=True, encoding="utf-8"
    )
    if result.returncode != 0 or not result.stdout.strip():
        print(f'[ERROR] Pupil ID "{selected_pupil_id}" is not valid.')
        return

    # Drop all existing grades
    grade_delete_query = tenant_db_connection_cmd + [
        "-e",
        f"DELETE FROM student_grades WHERE pupil_id = {selected_pupil_id}\
        AND section_id = {selected_section_id};",
    ]
    result = subprocess.run(grade_delete_query)
    if result.returncode != 0:
        print("[ERROR] Failed to delete existing grades.")
        return
    print("[INFO] Existing grades deleted successfully.")

    # Get all subject codes for section
    subjects_get_command = tenant_db_connection_cmd + [
        "-e",
        f"SELECT DISTINCT cs.subject_code FROM sections sec\
        JOIN ednevnik_workspace.curriculum_subjects cs ON\
        cs.curriculum_code = sec.curriculum_code\
        WHERE sec.id = {selected_section_id};",
    ]

    result = subprocess.run(
        subjects_get_command, capture_output=True, text=True, encoding="utf-8"
    )
    if result.returncode != 0 or not result.stdout.strip():
        print(f"[ERROR] Unable to retrieve subjects for section")
        return
    # Construct subject code array
    lines = result.stdout.strip().split("\n")
    subject_codes = [line.strip() for line in lines[1:] if line.strip()]

    insert_grade_query = (
        "INSERT INTO student_grades "
        "(type, pupil_id, section_id, subject_code, grade, semester_code, signature) VALUES "
    )

    for subject_code in subject_codes:
        first_semester_grade = random.randint(2, 5)
        second_semester_grade = random.randint(2, 5)

        insert_grade_query += f"('exam', {selected_pupil_id}, {selected_section_id}, '{subject_code}', {first_semester_grade}, '1POL', 'Script User'),"
        insert_grade_query += f"('final', {selected_pupil_id}, {selected_section_id}, '{subject_code}', {first_semester_grade}, '1POL', 'Script User'),"
        insert_grade_query += f"('exam', {selected_pupil_id}, {selected_section_id}, '{subject_code}', {second_semester_grade}, '2POL', 'Script User'),"
        insert_grade_query += f"('final', {selected_pupil_id}, {selected_section_id}, '{subject_code}', {second_semester_grade}, '2POL', 'Script User'),"

    insert_grade_query = insert_grade_query[:-1] + ";"

    insert_grade_command = tenant_db_connection_cmd + ["-e", insert_grade_query]

    result = subprocess.run(insert_grade_command)
    if result.returncode != 0:
        print("[ERROR] Failed to insert grades.")
        return
    print("[INFO] Grades inserted successfully.")


if __name__ == "__main__":
    main()
