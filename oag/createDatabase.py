from pathlib import Path

import click
import pandas
from pandas import DataFrame
from pyfs import isFile, resolvePath
from sqlalchemy import Engine, text
from sqlalchemy.exc import IntegrityError

from oag.sqlite import createDBConnection
from oag.sqlite.db import DB


def readWrite(
    fp: Path,
    dbPath: Path,
    columns: dict[str, str],
    dbConn: Engine,
    dbTable: str,
    index: bool = False,
    updateDFIndex: bool = False,
) -> None:
    print(f"Reading: ", fp.name)
    df: DataFrame = pandas.read_json(path_or_buf=fp)
    df.rename(columns=columns, inplace=True)
    df.columns = df.columns.str.lower()

    if updateDFIndex:
        print(f"Updating index to support database table {dbTable}...")
        sqlQuery: str = f"SELECT id FROM {dbTable} ORDER BY id DESC LIMIT 1"
        offset: int
        try:
            offset = dbConn.connect().execute(text(text=sqlQuery)).fetchall()[0]
        except IndexError:
            offset = 0

        df.index = df.index + offset + 1

    print(f"Writing to: ", dbPath.name)

    try:
        df.to_sql(
            name=dbTable,
            con=dbConn,
            if_exists="append",
            index=index,
            index_label="id",
        )
    except IntegrityError:
        print("Cannot write to database: IntegrityError")
        quit(1)


@click.command()
@click.option(
    "-c",
    "--input-cites-file",
    "citesFP",
    type=Path,
    required=True,
    help="Path to Citations JSON file",
)
@click.option(
    "-o",
    "--output",
    "outputFP",
    type=Path,
    required=True,
    help="Path to output database",
)
@click.option(
    "-w",
    "--input-works-file",
    "worksFP",
    type=Path,
    required=True,
    help="Path to Works JSON file",
)
def main(worksFP: Path, citesFP: Path, outputFP: Path) -> None:
    absWorksFP: Path = resolvePath(path=worksFP)
    absCitessFP: Path = resolvePath(path=citesFP)
    assert isFile(absWorksFP)
    assert isFile(absCitessFP)

    absOutputFP: Path = resolvePath(path=outputFP)

    dbConn: Engine = createDBConnection(dbPath=absOutputFP)
    db: DB = DB(dbConn=dbConn)
    del db

    worksColumns: dict[str, str] = {
        "Is_Paratext": "paratext",
        "Is_Retracted": "retracted",
        "Date_Published": "published",
    }

    citesColumns: dict[str, str] = {
        "Work_OA_ID": "work",
        "Ref_OA_ID": "reference",
    }

    readWrite(
        fp=absWorksFP,
        dbPath=absOutputFP,
        columns=worksColumns,
        dbConn=dbConn,
        dbTable="works",
    )

    readWrite(
        fp=absCitessFP,
        dbPath=absOutputFP,
        columns=citesColumns,
        dbConn=dbConn,
        dbTable="cites",
        index=True,
        updateDFIndex=True,
    )


if __name__ == "__main__":
    main()
