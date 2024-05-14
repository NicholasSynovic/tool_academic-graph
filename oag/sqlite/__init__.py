from pathlib import Path

from pandas import DataFrame
from sqlalchemy import Engine, create_engine


def createDBConnection(dbPath: Path) -> Engine:
    return create_engine(url=f"sqlite:///{dbPath}")


def saveData(df: DataFrame, table: str, dbConn: Engine) -> None:
    df.to_sql(
        name=table,
        con=dbConn,
        index=False,
        if_exists="append",
    )
