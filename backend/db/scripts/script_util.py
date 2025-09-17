import subprocess
import os
import argparse


def get_cmd(db_user, db_pass, is_root=True) -> list:
    """
    Returns the command to run the SQL files.
    """
    if os.name == "posix":
        if is_root:
            return ["sudo", "mariadb", "ednevnik_workspace"]

    if not db_user:
        if is_root:
            db_user = "root"
        else:
            raise ValueError("Database user must be provided.")

    return [
        "mariadb",
        f"-u{db_user}",
        f"-p{db_pass}",
        "ednevnik_workspace",
    ]


def check_mariadb_version(cmd, max_version=12) -> None:
    """
    Checks if the MariaDB version is compatible.
    GRANT TO PUBLIC must be 10.11 or higher.
    We want to use the latest version of MariaDB
    """
    try:
        version_cmd = cmd + ["-e", "SELECT VERSION();"]
        print(f"[INFO] Checking MariaDB version with command: {' '.join(version_cmd)}")
        result = subprocess.run(version_cmd, capture_output=True, text=True, check=True)
        print(f"[INFO] MariaDB version output: {result.stdout}")
        lines = result.stdout.strip().splitlines()
        if len(lines) < 2:
            raise RuntimeError("Could not parse MariaDB version output.")
        version_str = lines[1].split("-")[0]  # e.g. '10.6.22'
        major_version = int(version_str.split(".")[0])
        if major_version < max_version:
            raise RuntimeError(f"MariaDB version must be {max_version} or higher.")
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to check MariaDB version: {e}")


def parse_database_arguments():
    """Parse command line arguments for database connection configuration.

    Returns:
        argparse.Namespace: Parsed arguments containing database connection parameters

    Raises:
        NotImplementedError: If the operating system is not Windows or POSIX
    """
    parser = argparse.ArgumentParser(description="Database connection configuration")

    if not (os.name == "nt" or os.name == "posix"):
        raise NotImplementedError(
            "This script is only supported on Windows/POSIX systems."
        )

    parser.add_argument(
        "--password",
        "-p",
        required=False if os.name == "posix" else True,
        help="MySQL root password",
    )
    parser.add_argument(
        "--database",
        "-db",
        required=False,
        default="ednevnik_workspace",
        help="Workspace database name",
    )
    parser.add_argument(
        "--user", "-u", required=False, default="root", help="MySQL user"
    )

    return parser.parse_args()


def script_setup() -> list:
    """Sets up the script by parsing database arguments and checking the MariaDB version.

    Returns:
        list: The command to run the SQL files.
    """
    args = parse_database_arguments()
    cmd = get_cmd(args.user, args.password, is_root=True)
    print(f"[INFO] Command to execute SQL files: {' '.join(cmd)}")
    check_mariadb_version(cmd, max_version=11)
    return cmd
