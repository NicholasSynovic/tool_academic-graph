from pathlib import Path
from typing import List

from pandas import DataFrame
from works import buildDataFrame, readFile

from __init__ import createJSON


def main() -> None:
    jlPath: Path = Path("../../data/json/part_000")
    strData: List[str] = readFile(jlFilePath=jlPath)
    jsonData: List[dict] = createJSON(data=strData)
    df: DataFrame = buildDataFrame(data=jsonData)

    print(df)


if __name__ == "__main__":
    main()
