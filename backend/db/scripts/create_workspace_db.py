"""
Helper script to create (also drops the DBs so be careful) databases using
an SQL file.
"""

import subprocess
import os

SQL_DIR = os.path.join(os.path.dirname(__file__), '../sql')
SQL_FILE = os.path.join(SQL_DIR, 'create_workspace_db.sql')

with open(SQL_FILE, 'rb') as sql_file:
    subprocess.run([
        "mysql",
        "-u", "root",
        "-p1234",
        "ednevnik_workspace"
    ], stdin=sql_file)
