"""
Helper script to create (also drops the DBs so be careful) databases using
an SQL file.
"""

import subprocess
import os
from script_util import script_setup

SQL_DIR = os.path.join(os.path.dirname(__file__), "../sql")
STATIC_DIR = os.path.join(os.path.dirname(__file__), "../../static_data")
DROP_TENANT_SQL_FILE = os.path.join(SQL_DIR, "drop_tenant_databases.sql")
CREATE_WORKSPACE_SQL_FILE = os.path.join(SQL_DIR, "create_workspace_db.sql")
START_DATA_SQL_FILE = os.path.join(SQL_DIR, "sample_start_data.sql")


def main():
    cmd = script_setup()

    for sql_command in [
        DROP_TENANT_SQL_FILE,
        CREATE_WORKSPACE_SQL_FILE,
        START_DATA_SQL_FILE,
    ]:
        print(f"[INFO] Executing SQL file: {sql_command}")
        with open(sql_command, "rb") as sql_file:
            subprocess.run(cmd, stdin=sql_file)
        print(f"[INFO] Finished executing: {sql_command}")

    print(f"[INFO] Running Go script in directory: {STATIC_DIR}")
    subprocess.run(["go", "run", "."], cwd=STATIC_DIR)
    print("[INFO] Go script execution finished.")


if __name__ == "__main__":
    main()
