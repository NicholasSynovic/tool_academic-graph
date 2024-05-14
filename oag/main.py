from pathlib import Path
from typing import List

import click
from oa_graph.json_lines.utils import createJSON
from oa_graph.json_lines.works import buildDataFrame, readFile
from oa_graph.sqlite import createDBConnection, saveData
from oa_graph.sqlite.db import DB
from pandas import DataFrame
from pyfs import isFile, resolvePath
from sqlalchemy import Engine


def insertWorks(df: DataFrame, dbConn: Engine) -> None:
    foo: DataFrame = df.drop(columns=["cites", "venue_oa_id"])
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

    print(f"Reading {absInputFP.name}...")
    strData: List[str] = readFile(jlFilePath=absInputFP)

    jsonData: List[dict] = createJSON(data=strData)

    df: DataFrame = buildDataFrame(data=jsonData)
    insertWorks(df=df, dbConn=dbConn)


if __name__ == "__main__":
    main()
