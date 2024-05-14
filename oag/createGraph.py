from math import ceil
from pathlib import Path

import click
import pandas
from pandas import DataFrame
from progress.bar import Bar
from pyfs import isFile, resolvePath
from sqlalchemy import Engine

from oag.sqlite import createDBConnection
from oag.sqlite.db import DB


def addWorks(db: DB, dbConn: Engine, chunksize: int = 1000) -> None:
    worksCount: int = db.getWorkCount()

    with Bar("Adding works...", max=ceil(worksCount / chunksize)) as bar:
        df: DataFrame
        for df in pandas.read_sql_table(
            table_name="works",
            con=dbConn,
            chunksize=1000,
        ):
            bar.next()


@click.command()
@click.option(
    "-i",
    "--input",
    "inputFP",
    type=Path,
    required=True,
    help="Path to SQLite3 database",
)
def main(inputFP: Path) -> None:
    absInputFP: Path = resolvePath(path=inputFP)
    assert isFile(path=absInputFP)

    dbConn: Engine = createDBConnection(dbPath=absInputFP)
    db: DB = DB(dbConn=dbConn)
    # citesLargestID: int = db.getLargestCitesID()

    addWorks(db=db, dbConn=dbConn)


if __name__ == "__main__":
    main()
