import subprocess
from pathlib import Path
from subprocess import CompletedProcess


def getJSONSize(fp: Path) -> int:
    command: str = f"cat {fp} | tail -n 1 | jq '.id'"

    try:
        result: CompletedProcess = subprocess.run(
            command,
            shell=True,
            check=True,
            stdout=subprocess.PIPE,
            text=True,
        )

        return int(result.stdout.strip()) + 1

    except subprocess.CalledProcessError as e:
        print(f"An error occurred: {e}")
        quit(1)
    except ValueError as e:
        print(f"Failed to convert output to integer: {e}")
        quit(1)
