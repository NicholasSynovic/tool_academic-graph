from itertools import product
from pathlib import Path
from typing import List, Tuple

import click
import pandas
from pandas import DataFrame
from progress.bar import Bar
from pyfs import isFile, resolvePath
from sqlalchemy import Engine

from oag import CITATIONS_JQ_FORMAT
from oag.json_lines.json_lines import JSONLines
from oag.sqlite import createDBConnection, saveData
from oag.sqlite.db import DB


def insertCitations(df: DataFrame, dbConn: Engine) -> None:
    data: List[dict[str, List[str]]] = []

    df["oa_id"] = df["oa_id"].str.replace(pat="https://openalex.org/", repl="")
    df["cites"] = df["cites"].apply(
        lambda values: [
            value.replace("https://openalex.org/", "")
            for value in values
            if value is not None
        ]
    )

    with Bar("Extracting citation relationships...", max=df.shape[0]) as bar:
        row: Tuple[int, str, List[str]]
        for row in df.itertuples():
            document: List[str] = [row[1]]
            cites: List[str] = row[2]
            pairs: List[dict[str, str]] = [
                {"work": p[0], "reference": p[1]} for p in product(document, cites)
            ]
            data.append(DataFrame(data=pairs))
            bar.next()

    foo: DataFrame = pandas.concat(objs=data, ignore_index=True)

    print(foo)

    saveData(df=foo, table="cites", dbConn=dbConn, includeIndex=True)


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
