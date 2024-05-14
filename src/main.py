from pathlib import Path
from typing import List

import click
from pandas import DataFrame
from pyfs import isFile, resolvePath
from sqlalchemy import Connection

from src.json.utils import createJSON
from src.json.works import buildDataFrame, readFile
from src.sqlite import createDBConnection, saveData
from src.sqlite.db import DB


def insertWorks(df: DataFrame, dbConn: Connection) -> None:
    df.drop(columns=["cites", "venue_oa_id"], inplace=True)
    saveData(df=df, table="works", dbConn=dbConn)


def insertCites(df: DataFrame, dbConn: Connection) -> None:
    df: DataFrame = df[df["oa_id", "cites"]]
    print(df)


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

    dbConn: Connection = createDBConnection(dbPath=absOutputFP)
    db: DB = DB(dbConn=dbConn)  # TODO: Make this a generic SQLite obj
    db.createTables()

    strData: List[str] = readFile(jlFilePath=absInputFP)
    jsonData: List[dict] = createJSON(data=strData)

    df: DataFrame = buildDataFrame(data=jsonData)
    insertWorks(df=df, dbConn=dbConn)
    insertCites(df=df, dbConn=dbConn)


if __name__ == "__main__":
    main()
