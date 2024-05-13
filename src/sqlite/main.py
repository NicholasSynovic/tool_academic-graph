import pickle
from pathlib import Path
from typing import List

from pandas import DataFrame
from sqlalchemy import Connection

from src.json import createJSON
from src.json.works import buildDataFrame, readFile
from src.sqlite import createDBConnection, saveData
from src.sqlite.works import createSchema


def main() -> None:
    jlPath: Path = Path("../../data/json/part_000")
    dbPath: Path = Path("works_" + jlPath.name + ".db")

    dbConn: Connection = createDBConnection(dbPath=dbPath)
    createSchema(dbConn=dbConn)

    # strData: List[str] = readFile(jlFilePath=jlPath)
    # jsonData: List[dict] = createJSON(data=strData)

    # with open("jd.pickle", "wb") as pf:
    #     pickle.dump(jsonData, pf)
    #     pf.close()

    with open("jd.pickle", "rb") as pf:
        jsonData: List[dict] = pickle.load(pf)
        pf.close()

    df: DataFrame = buildDataFrame(data=jsonData)
    df.drop(columns="cites", inplace=True)
    df.drop(columns="venue_oa_id", inplace=True)

    print(df.columns)

    saveData(df=df, table="works", dbConn=dbConn)


if __name__ == "__main__":
    main()
