from concurrent.futures import ProcessPoolExecutor, ThreadPoolExecutor
from itertools import count
from pathlib import Path
from re import search
from typing import Hashable, List, Tuple
from xml.etree.ElementTree import Element, ElementTree

import click
import pandas
from lxml import etree
from pandas import DataFrame, Series
from progress.bar import Bar
from pyfs import isFile, resolvePath
from sqlalchemy import Engine

from oag.sqlite import createDBConnection
from oag.sqlite.db import DB


def buildXML(
    nodesDF: DataFrame,
    edgesDF: DataFrame,
) -> str:
    xmlns: str = f"http://graphml.graphdrawing.org/xmlns"

    rootNode: Element = etree.Element("graphml", nsmap={None: xmlns})

    graphNode: Element = etree.SubElement(rootNode, "graph")
    graphNode.set("id", "Graph")
    graphNode.set("edgedefault", "directed")

    nodeMapping: dict[str, str] = {}

    with Bar("Adding nodes...", max=nodesDF.shape[0]) as bar:
        with ThreadPoolExecutor() as executor:

            def _run(row: Tuple[Hashable, Series]) -> None:
                node: Element = etree.SubElement(graphNode, "node")
                oaID: Element = etree.SubElement(node, "data")

                nm: str = f"n{row[0]}"
                nodeMapping[row[1]["oa_id"]] = nm

                node.set("id", nm)
                oaID.set("key", "oa_id")
                oaID.text = row[1]["oa_id"]
                bar.next()

            executor.map(_run, nodesDF.iterrows())
            bar.finish()
            bar.update()

    with Bar("Adding edges...", max=edgesDF.shape[0]) as bar:
        row: Tuple[Hashable, Series]
        for row in edgesDF.iterrows():
            try:
                nodeMapping[row[1]["work"]]
                nodeMapping[row[1]["reference"]]
            except KeyError:
                bar.next()
                continue

            edge: Element = etree.SubElement(graphNode, "edge")
            edge.set(key="id", value=f"e{row[0]}")

            edge.set(
                key="source",
                value=nodeMapping[row[1]["work"]],
            )

            edge.set(
                key="target",
                value=nodeMapping[row[1]["reference"]],
            )

            bar.next()

    xmlStr: str = etree.tostring(rootNode, pretty_print=True).decode()
    return xmlStr


@click.command()
@click.option(
    "dbPath",
    "-i",
    "--input",
    nargs=1,
    type=Path,
    required=True,
    help="Path to SQLite3 database",
)
@click.option(
    "graphMLPath",
    "-o",
    "--output",
    nargs=1,
    type=Path,
    required=True,
    help="Path to store GraphML file",
)
def main(dbPath: Path, graphMLPath: Path) -> None:
    absInputFP: Path = resolvePath(path=dbPath)
    assert isFile(path=absInputFP)

    dbConn: Engine = createDBConnection(dbPath=absInputFP)
    db: DB = DB(dbConn=dbConn)

    print("Reading works table...")
    worksDF: DataFrame = pandas.read_sql_table(
        table_name="works",
        con=dbConn,
        columns=["oa_id"],
    )
    print("Read works table")

    print("Reading cites table...")
    citesDF: DataFrame = pandas.read_sql_table(table_name="cites", con=dbConn)
    print("Read cites table")

    xmlStr: str = buildXML(nodesDF=worksDF, edgesDF=citesDF)

    print("Writing file: ", graphMLPath)
    with open(file=graphMLPath, mode="w") as xmlFile:
        xmlFile.write(xmlStr)
        xmlFile.close()


if __name__ == "__main__":
    main()
