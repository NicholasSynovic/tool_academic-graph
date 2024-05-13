from concurrent.futures import ProcessPoolExecutor, ThreadPoolExecutor
from json import loads as jsonLoad
from os import listdir
from pathlib import Path
from typing import Iterable, List

import click
from progress.bar import Bar
from pyfs import isDirectory, resolvePath


def readJSONLines(f: Path) -> List[dict]:
    with open(file=f, mode="r") as jsonFile:
        lines: List[str] = jsonFile.readlines()
        jsonFile.close()

    with Bar("Reading JSON lines...", max=len(lines)) as bar:

        def _read(line: str) -> dict:
            data: dict = jsonLoad(s=line)
            bar.next()
            return data

        with ProcessPoolExecutor() as executor:
            results: Iterable[dict] = executor.map(_read, lines)

    return list(results)


@click.command()
@click.option(
    "-i",
    "--input",
    "inputDirectory",
    required=True,
    type=Path,
    help="Path to a directory containing an OpenAlex database dump",
)
def main(inputDirectory: Path) -> None:
    assert isDirectory(path=inputDirectory)
    directory: Path = resolvePath(path=inputDirectory)

    files: List[Path] = [Path(directory, file) for file in listdir(path=directory)]

    for file in files:
        readJSONLines(f=file)


if __name__ == "__main__":
    main()
