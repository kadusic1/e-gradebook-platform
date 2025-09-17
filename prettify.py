# This is a helper script to prettify JavaScript files in the frontend directory
# based on the output of `git status`.
import subprocess
import os

def main():
    # Run git status command
    result = subprocess.run(
        ['git', 'status'],
        capture_output=True,
        text=True,
        cwd=os.getcwd()  # Run in current directory
    )
    text = result.stdout
    files_to_prettify = []

    # Iterate over each line in the output
    for line in text.splitlines():
        # Iterate over each word in the line
        previousWord = None
        for word in line.split():
            # See if previous word does not contain deleted
            if (
                (
                    previousWord is None or
                    "deleted:" not in previousWord
                )
                and "frontend/" in word
                # And file ends with js or jsx
                and (word.endswith('.js') or word.endswith('.jsx'))
            ):
                files_to_prettify.append(word)
            previousWord = word

    print("Prettifying files from git status!")
    print("-"*40)

    # Run prettier on each file in one command
    command = f"npx prettier --write {' '.join(files_to_prettify)}"
    print(f"Running: {command}")
    subprocess.run(
        command,
        shell=True,
        check=True,
        cwd=os.getcwd()
    )

if __name__ == "__main__":
    main()
