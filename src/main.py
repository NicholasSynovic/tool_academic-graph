from pathlib import Path
from typing import List

import click
from pandas import DataFrame
from pyfs import isFile, resolvePath
from sqlalchemy import Connection

from src.json import createJSON
from src.json.works import buildDataFrame, readFile
from src.sqlite import createDBConnection, saveData
from src.sqlite.works import Works


def insertWorks(df: DataFrame, dbConn: Connection) -> None:
    pass


@click.command()
@click.option(
    "-i",
    "--input",
    "inputFP",
    type=Path,
    required=True,
    help="Path to OpenAlex JSON Lines file",
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
    assert absInputFP

    absOutputFP: Path = resolvePath(path=outputFP)

    dbConn: Connection = createDBConnection(dbPath=absOutputFP)
    works: Works = Works(dbConn=dbConn)  # TODO: Make this a generic SQLite obj

    strData: List[str] = readFile(jlFilePath=absInputFP)
    jsonData: List[dict] = createJSON(data=strData)

    df: DataFrame = buildDataFrame(data=jsonData)
    df.drop(columns="cites", inplace=True)
    df.drop(columns="venue_oa_id", inplace=True)

    saveData(df=df, table="works", dbConn=dbConn)


if __name__ == "__main__":
    main()
