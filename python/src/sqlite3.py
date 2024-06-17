"""
Steps

5. Get the length of the JSON file from JQ
6. Read in JSON file as chunks into Pandas DFS
7. Pass each chunk into SQLite3 db
"""

from math import ceil
from pathlib import Path

import click
import pandas
from pandas import DataFrame
from pandas.io.json._json import JsonReader
from progress.bar import Bar
from pyfs import isFile, resolvePath
from src.shell import getJSONSize


@click.command()
@click.option(
    "-i",
    "--input",
    "inputPath",
    type=Path,
    help="Path to OpenAlex JSON Works File generated from Go scripts",
    required=True,
)
@click.option(
    "-o",
    "--output",
    "outputPath",
    type=Path,
    help="Path to SQLite3 Database to store JSON works data",
    required=True,
)
@click.option(
    "--chunksize",
    "chunksize",
    type=int,
    help="Chunksize while iterating through JSON files",
    required=False,
    default=10000,
)
def main(inputPath: Path, outputPath: Path, chunksize: int) -> None:
    absInputPath: Path = resolvePath(path=inputPath)
    absOutputPath: Path = resolvePath(path=outputPath)

    assert isFile(path=absInputPath)
    assert isFile(path=absOutputPath) == False

    print("Getting the total number of objects in the JSON file...")
    objCount: int = getJSONSize(fp=absInputPath)

    dfs: JsonReader[DataFrame] = pandas.read_json(
        path_or_buf=absInputPath,
        lines=True,
        chunksize=chunksize,
    )

    with Bar(
        "Iterating through DataFrames...", max=int(ceil(objCount / chunksize))
    ) as bar:
        for df in dfs:
            df.set_index(keys="id", inplace=True)
            print(df)
            quit()
            bar.next()


if __name__ == "__main__":
    main()
