"""
Helper script to see current state of the databases using
a SQL file.
"""

import subprocess
import os

SQL_DIR = os.path.join(os.path.dirname(__file__), '../sql')
SQL_FILE = os.path.join(SQL_DIR, 'test_db_state.sql')

with open(SQL_FILE, 'rb') as sql_file:
    subprocess.run([
        "mysql",
        "-u", "root",
        "-p1234",
        "ednevnik_workspace"
    ], stdin=sql_file)