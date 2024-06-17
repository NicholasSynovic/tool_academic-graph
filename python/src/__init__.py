from pathlib import Path

from pandas import DataFrame
from sqlalchemy import Engine, create_engine


def createDBConnection(dbPath: Path) -> Engine:
    return create_engine(url=f"sqlite:///{dbPath}")


def saveData(
    df: DataFrame, table: str, dbConn: Engine, includeIndex: bool = False
) -> None:
    if includeIndex:
        df.to_sql(
            name=table,
            con=dbConn,
            index=True,
            index_label="id",
            if_exists="append",
        )
    else:
        df.to_sql(
            name=table,
            con=dbConn,
            index=False,
            if_exists="append",
        )
