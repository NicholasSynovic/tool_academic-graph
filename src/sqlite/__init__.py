from pathlib import Path

from pandas import DataFrame
from sqlalchemy import Connection, create_engine


def createDBConnection(dbPath: Path) -> Connection:
    return create_engine(url=f"sqlite:///{dbPath}").connect()


def saveData(df: DataFrame, table: str, dbConn: Connection) -> None:
    df.to_sql(name=table, con=dbConn, if_exists="replace")