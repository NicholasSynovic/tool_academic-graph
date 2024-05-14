from pathlib import Path
from subprocess import DEVNULL, PIPE, CompletedProcess, run
from typing import List

import pandas
from pandas import DataFrame


def readFile(jlFilePath: Path) -> List[str]:
    process: CompletedProcess = run(
        [
            "jq",
            "-c",
            "{doi: .doi, title: .title, oa_id: .id, paratext: .is_paratext, retracted: .is_retracted, venue_oa_id: .primary_location.source.id, published: .publication_date, oa_type: .type, cf_type: .type_crossref, cites: .referenced_works[]}",
            jlFilePath,
        ],
        stderr=DEVNULL,
        stdout=PIPE,
    )
    outputString: str = process.stdout.decode().strip()
    return outputString.split(sep="\n")


def buildDataFrame(data: List[dict]) -> DataFrame:
    df: DataFrame = DataFrame(data=data)
    df["doi"] = df["doi"].str.replace("https://doi.org/", "")
    df["oa_id"] = df["oa_id"].str.replace("https://openalex.org/", "")
    df["paratext"] = df["paratext"].astype(bool)
    df["retracted"] = df["retracted"].astype(bool)
    df["published"] = pandas.to_datetime(
        df["published"],
        format="%Y-%m-%d",
        errors="coerce",
    )
    df.dropna(subset="published", inplace=True)

    return df
