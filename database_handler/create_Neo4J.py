import asyncio
from concurrent.futures import ThreadPoolExecutor
from math import ceil
from pathlib import Path
from typing import Generator

import click
import pandas
from pandas import DataFrame
from progress.bar import Bar
from pyfs import isFile, resolvePath
from sqlalchemy import Engine

from database_handler.graph.graph import Neo4J
from database_handler.sqlite import createDBConnection
from database_handler.sqlite.db import DB


def addWorks(
    db: DB,
    dbConn: Engine,
    neo4j: Neo4J,
    chunksize: int = 1000,
) -> None:
    worksCount: int = db.getWorkCount()

    with Bar("Adding work nodes...", max=ceil(worksCount / chunksize)) as bar:
        with ThreadPoolExecutor() as executor:

            def _run(df: DataFrame) -> None:
                bar.next()
                neo4j.addWorkNode(df=df)

            dfs: Generator = pandas.read_sql_table(
                table_name="works",
                con=dbConn,
                chunksize=chunksize,
            )

            executor.map(_run, dfs)


def addCites(
    db: DB,
    dbConn: Engine,
    neo4j: Neo4J,
    chunksize: int = 1000,
) -> None:
    citesCount: int = db.getLargestCitesID()

    with Bar(
        "Adding citation relationships...", max=ceil(citesCount / chunksize)
    ) as bar:
        with ThreadPoolExecutor() as executor:

            def _run(df: DataFrame) -> None:
                neo4j.addRelationship(df=df)
                bar.next()

            dfs: Generator = pandas.read_sql_table(
                table_name="cites",
                con=dbConn,
                chunksize=chunksize,
            )

            executor.map(_run, dfs)


def async_AddCites(
    db: DB,
    dbConn: Engine,
    neo4j: Neo4J,
) -> None:
    citesCount: int = db.getLargestCitesID()

    print("Reading works table...")
    worksDF: DataFrame = pandas.read_sql_table(
        table_name="works",
        con=dbConn,
    )
    print("Read works table")

    print("Reading cites table...")
    citesDF: DataFrame = pandas.read_sql_table(
        table_name="cites",
        con=dbConn,
    )
    print("Read cites table")

    df: DataFrame = citesDF[~citesDF["reference"].isin(worksDF["oa_id"])]
    print(citesDF.shape)
    print(df.shape)
    quit()

    with Bar("Adding citation relationships...", max=ceil(citesCount)) as bar:
        asyncio.run(neo4j.async_AddRelationships(df=df, bar=bar))


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
    neo4j: Neo4J = Neo4J(
        uri=neo4jURI,
        username=neo4jUsername,
        password=neo4jPassword,
    )

    # addWorks(
    #     db=db,
    #     dbConn=dbConn,
    #     neo4j=neo4j,
    # )

    # neo4j.createWorkNodeIndex()

    # addCites(db=db, dbConn=dbConn, neo4j=neo4j)

    async_AddCites(db=db, dbConn=dbConn, neo4j=neo4j)


if __name__ == "__main__":
    main()
