from json import JSONDecodeError, loads
from pathlib import Path
from subprocess import DEVNULL, PIPE, CompletedProcess, run
from typing import List

from pandas import DataFrame
from progress.bar import Bar


def readFile(jlFilePath: Path) -> List[str]:
    process: CompletedProcess = run(
        [
            "jq",
            "-c",
            "{doi: .doi, title: .title, oa_id: .id, paratext: .is_paratext, retracted: .is_retracted, venue_oa_id: .primary_location.source.id, published: .publication_date, oa_type: .type, cf_type: .type_crossref, cites: [.locations[] | .source.id]}",
            jlFilePath,
        ],
        stderr=DEVNULL,
        stdout=PIPE,
    )
    outputString: str = process.stdout.decode().strip()
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
    df["doi"] = df["doi"].str.replace("https://doi.org/", "")
    df["paratext"] = df["paratext"].astype(bool)
    df["retracted"] = df["retracted"].astype(bool)
    return df
