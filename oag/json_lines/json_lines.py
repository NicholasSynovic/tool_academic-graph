from concurrent.futures import ThreadPoolExecutor
from json import loads
from pathlib import Path
from subprocess import DEVNULL, PIPE, CompletedProcess, run
from typing import List

import pandas
from pandas import DataFrame
from progress.bar import Bar


class JSONLines:
    def __init__(self, jlFilePath: Path) -> None:
        self.filepath: Path = jlFilePath
        assert jlFilePath.exists()

    def read(self, jqFormat: str) -> DataFrame:
        """
        read Read the JSON Lines file and store each object as a row within a Pandas DataFrame.

        :param jqFormat: A string representation of the data to be extracted.
            Follows jq output formatting.

            Example:
            "{doi: .doi, title: .title, oa_id: .id, paratext: .is_paratext, retracted: .is_retracted, venue_oa_id: .primary_location.source.id, published: .publication_date, oa_type: .type, cf_type: .type_crossref, cites: .referenced_works}"
        :type jqFormat: str
        :return: A DataFrame of the relevant data from the JSON Lines file
        :rtype: DataFrame
        """
        print(f"Reading {self.filepath.name}...")
        process: CompletedProcess = run(
            [
                "jq",
                "-c",
                jqFormat,
                self.filepath,
            ],
            stderr=DEVNULL,
            stdout=PIPE,
        )
        assert process.returncode == 0

        outputString: str = process.stdout.decode().strip()

        stringData: List[str] = outputString.split(sep="\n")

        with Bar(
            "Converting JSON strings to Python dicts...", max=len(stringData)
        ) as bar:
            with ThreadPoolExecutor() as executor:

                def _run(string: str) -> dict:
                    data: dict = loads(s=string)
                    bar.next()
                    return data

                dictData: List[dict] = list(executor.map(_run, stringData))

        return DataFrame(data=dictData)
