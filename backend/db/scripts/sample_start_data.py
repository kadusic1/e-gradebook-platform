"""
Script to execute sample_start_data.sql
"""

import subprocess
import os

SQL_DIR = os.path.join(os.path.dirname(__file__), '../sql')
SQL_FILE = os.path.join(SQL_DIR, 'sample_start_data.sql')

command = [
    "mysql",
    "-u", "root",
    "-p1234",
    "ednevnik_workspace"
    "<", SQL_FILE
]

full_command = ' '.join(command)

subprocess.run(full_command, shell=True)
