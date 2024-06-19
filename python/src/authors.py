from math import ceil
from pathlib import Path
from typing import List

import click
import pandas
from pandas import DataFrame
from pandas.io.json._json import JsonReader
from progress.bar import Bar
from pyfs import isFile, resolvePath
from sqlalchemy import Engine
from src import createDBConnection, saveData
from src.db import DB
from src.shell import getJSONSize


@click.command()
@click.option(
    "-i",
    "--input",
    "inputPath",
    type=Path,
    help="Path to OpenAlex JSON Authors File generated from Go scripts",
    required=True,
)
@click.option(
    "-o",
    "--output",
    "outputPath",
    type=Path,
    help="Path to SQLite3 Database to store JSON Authors data",
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

    authorsSet: set = set([])

    engine: Engine = createDBConnection(dbPath=absOutputPath)
    db: DB = DB(dbConn=engine)

    print("Getting the total number of objects in the JSON file...")
    objCount: int = getJSONSize(fp=absInputPath)

    dfs: JsonReader[DataFrame] = pandas.read_json(
        path_or_buf=absInputPath,
        lines=True,
        chunksize=chunksize,
        compression="infer",
    )

    with Bar(
        "Iterating through DataFrames...", max=int(ceil(objCount / chunksize))
    ) as bar:
        for df in dfs:
            df.set_index(keys="id", inplace=True)

            # Remove duplicate ORCIDs in the DataFrame
            df = df[~df.duplicated(subset="orcid", keep=False)]

            # Make sure ORCIDs are not in the set of used DOIs
            df = df[~df["orcid"].isin(values=authorsSet)]

            # Add ORCIDs to orcidSet
            authorsSet.update(df["orcid"].to_list())

            saveData(df=df, table="authors", dbConn=engine)
            bar.next()


if __name__ == "__main__":
    main()
