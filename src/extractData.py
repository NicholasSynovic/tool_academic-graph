import pickle
from json import JSONDecodeError, loads
from pathlib import Path
from subprocess import DEVNULL, PIPE, CompletedProcess, run
from typing import List

import click
from pandas import DataFrame
from progress.bar import Bar
from pyfs import isFile, resolvePath


def readFile(jlFilePath: Path) -> List[str]:
    print(f"Extracting data from {jlFilePath}...")
    process: CompletedProcess = run(
        [
            "jq",
            "-c",
            "{doi: .doi, title: .title, oa_id: .id, paratext: .is_paratext, retracted: .is_retracted, venue_oa_id: .primary_location.source.id, published: .publication_date, cites_oa_id: .referenced_works, oa_type: .type, cf_type: .type_crossref}",
            jlFilePath,
        ],
        stderr=DEVNULL,
        stdout=PIPE,
    )
    outputString: str = process.stdout.decode()
    print(f"Extracted data from {jlFilePath}")
    return outputString.split(sep="\n")


def createJSON(data: List[str]) -> List[dict]:
    jsonObjects: List[dict] = []

    with Bar("Converting JSON strings into JSON objects...", max=len(data)) as bar:
        datum: str
        for datum in data:
            try:
                json: dict = loads(s=datum)
            except JSONDecodeError:
                bar.next()
                continue

            jsonObjects.append(json)
            bar.next()

    return jsonObjects


def buildDataFrame(data: List[dict]) -> DataFrame:
    df: DataFrame = DataFrame(data=data)

    with open("df.pickle", "wb") as pf:
        pickle.dump(obj=df, file=pf)
        pf.close()

    df["doi"] = df["doi"].str.replace("https://doi.org/", "")
    return df


@click.command()
@click.option(
    "-i",
    "--input",
    "inputJSONLinesFilePath",
    type=Path,
    required=True,
    help="Path to JSON Lines file from OpenAlex",
)
def main(inputJSONLinesFilePath: Path) -> None:
    absJSONLinesFilePath: Path = resolvePath(path=inputJSONLinesFilePath)
    if isFile(path=absJSONLinesFilePath) == False:
        quit(1)

    strData: List[str] = readFile(jlFilePath=absJSONLinesFilePath)
    jsonData: List[dict] = createJSON(data=strData)
    df: DataFrame = buildDataFrame(data=jsonData)
    print(df)


if __name__ == "__main__":
    main()
