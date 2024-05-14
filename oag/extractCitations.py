from pathlib import Path
from typing import List

import click
import pandas
from pandas import DataFrame
from pyfs import isFile, resolvePath
from sqlalchemy import Engine

from oag import CITATIONS_JQ_FORMAT
from oag.json_lines.json_lines import JSONLines
from oag.sqlite import createDBConnection, saveData
from oag.sqlite.db import DB


def insertCitations(df: DataFrame, dbConn: Engine) -> None:
    df["oa_id"] = df["oa_id"].str.replace(pat="https://openalex.org/", repl="")
    df["cites"] = df["cites"].apply(
        lambda values: [
            value.replace("https://openalex.org/", "")
            for value in values
            if value is not None
        ]
    )

    saveData(df=df, table="cites", dbConn=dbConn)


@click.command()
@click.option(
    "-i",
    "--input",
    "inputFP",
    type=Path,
    required=True,
    help='Path to OpenAlex JSON Lines "Works" file',
)
@click.option(
    "-o",
    "--output",
    "outputFP",
    type=Path,
    required=True,
    help="Path to SQLite3 file",
)
def main(inputFP: Path, outputFP: Path) -> None:
    absInputFP: Path = resolvePath(path=inputFP)
    assert isFile(path=absInputFP)

    absOutputFP: Path = resolvePath(path=outputFP)

    dbConn: Engine = createDBConnection(dbPath=absOutputFP)
    DB(dbConn=dbConn)

    JL: JSONLines = JSONLines(jlFilePath=absInputFP)
    df: DataFrame = JL.read(jqFormat=CITATIONS_JQ_FORMAT)

    insertCitations(df=df, dbConn=dbConn)


if __name__ == "__main__":
    main()
