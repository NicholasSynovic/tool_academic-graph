from pathlib import Path
from typing import List

import click
import pandas
from pandas import DataFrame
from pyfs import isFile, resolvePath
from sqlalchemy import Engine

from oag import WORKS_JQ_FORMAT
from oag.json_lines.json_lines import JSONLines
from oag.sqlite import createDBConnection, saveData
from oag.sqlite.db import DB


def insertWorks(df: DataFrame, dbConn: Engine) -> None:
    df["oa_id"] = df["oa_id"].str.replace(pat="https://openalex.org/", repl="")
    df["doi"] = df["doi"].str.replace(pat="https://doi.org/", repl="")
    df["title"] = df["title"].str.title()
    df["paratext"] = df["paratext"].astype(dtype=bool)
    df["retracted"] = df["retracted"].astype(dtype=bool)
    df["published"] = pandas.to_datetime(
        arg=df["published"],
        errors="coerce",
        format="%Y-%m-%d",
    )
    foo: DataFrame = df.dropna(subset=["published"])
    saveData(df=foo, table="works", dbConn=dbConn)


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
    df: DataFrame = JL.read(jqFormat=WORKS_JQ_FORMAT)

    insertWorks(df=df, dbConn=dbConn)


if __name__ == "__main__":
    main()
