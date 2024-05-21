from pathlib import Path

import click
import pandas
from pandas import DataFrame
from pyfs import isFile, resolvePath
from sqlalchemy import Engine

from oag.sqlite import createDBConnection
from oag.sqlite.db import DB


@click.command()
@click.option(
    "-i",
    "--input",
    "inputFP",
    type=Path,
    required=True,
    help="Path to JSON file",
)
@click.option(
    "-o",
    "--output",
    "outputFP",
    type=Path,
    required=True,
    help="Path to output database",
)
def main(inputFP: Path, outputFP: Path) -> None:
    absInputFP: Path = resolvePath(path=inputFP)
    assert isFile(absInputFP)

    absOutputFP: Path = resolvePath(path=outputFP)

    dbConn: Engine = createDBConnection(dbPath=absOutputFP)
    db: DB = DB(dbConn=dbConn)
    del db

    print(f"Reading: ", absInputFP.name)
    df: DataFrame = pandas.read_json(path_or_buf=absInputFP)
    df.rename(
        columns={
            "Is_Paratext": "paratext",
            "Is_Retracted": "retracted",
            "Date_Published": "published",
        },
        inplace=True,
    )
    df.columns = df.columns.str.lower()

    print(f"Writing to: ", absOutputFP.name)
    df.to_sql(name="works", con=dbConn, if_exists="append", index=False)


if __name__ == "__main__":
    main()
