from math import ceil
from pathlib import Path

import click
import pandas
from pandas import DataFrame
from progress.bar import Bar
from pyfs import isFile, resolvePath
from sqlalchemy import Engine

from oag.graph.graph import Neo4J
from oag.sqlite import createDBConnection
from oag.sqlite.db import DB


def addWorks(
    db: DB,
    dbConn: Engine,
    username: str,
    password: str,
    uri: str,
    chunksize: int = 1000,
) -> None:
    neo4j: Neo4J = Neo4J(uri=uri, username=username, password=password)

    worksCount: int = db.getWorkCount()

    with Bar("Adding works...", max=ceil(worksCount / chunksize)) as bar:
        df: DataFrame
        for df in pandas.read_sql_table(
            table_name="works",
            con=dbConn,
            chunksize=1000,
        ):
            neo4j.addNode(df=df)
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
@click.option(
    "-u",
    "--uri",
    "neo4jURI",
    type=str,
    required=False,
    show_default=True,
    default="bolt://localhost:7687",
    help="Bolt URI of Neo4J instance to connect to",
)
@click.option(
    "--username",
    "neo4jUsername",
    type=str,
    required=True,
    help="Username of Neo4J user",
)
@click.option(
    "--password",
    "neo4jPassword",
    type=str,
    required=True,
    help="Password of Neo4J user",
)
def main(
    inputFP: Path,
    neo4jUsername: str,
    neo4jPassword: str,
    neo4jURI: str = "bolt://localhost:7474",
) -> None:
    absInputFP: Path = resolvePath(path=inputFP)
    assert isFile(path=absInputFP)

    dbConn: Engine = createDBConnection(dbPath=absInputFP)
    db: DB = DB(dbConn=dbConn)
    # citesLargestID: int = db.getLargestCitesID()

    addWorks(
        db=db,
        dbConn=dbConn,
        username=neo4jUsername,
        password=neo4jPassword,
        uri=neo4jURI,
    )


if __name__ == "__main__":
    main()
