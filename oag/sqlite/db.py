from typing import List

from sqlalchemy import (
    Boolean,
    Column,
    DateTime,
    Engine,
    ForeignKeyConstraint,
    Integer,
    MetaData,
    PrimaryKeyConstraint,
    Row,
    String,
    Table,
    TextClause,
    text,
)


class DB:
    def __init__(self, dbConn: Engine) -> None:
        self.dbConn: Engine = dbConn
        metadata: MetaData = MetaData()

        worksSchema: Table = Table(
            "works",
            metadata,
            Column("oa_id", String),
            Column("doi", String),
            Column("title", String),
            Column("paratext", Boolean),
            Column("retracted", Boolean),
            Column("published", DateTime),
            Column("oa_type", String),
            Column("cf_type", String),
            PrimaryKeyConstraint("oa_id"),
        )

        citesSchema: Table = Table(
            "cites",
            metadata,
            Column("id", Integer),
            Column("work", String),
            Column("reference", String),
            PrimaryKeyConstraint("id"),
            ForeignKeyConstraint(
                columns=["work"],
                refcolumns=["works.oa_id"],
            ),
            ForeignKeyConstraint(
                columns=["reference"],
                refcolumns=["works.oa_id"],
            ),
        )

        metadata.create_all(bind=self.dbConn)

    def getLargestCitesID(self) -> int:
        sql: TextClause = text(text="SELECT id FROM cites ORDER BY id DESC LIMIT 1;")
        with self.dbConn.connect() as connection:
            try:
                result: Row = list(connection.execute(sql))[0]
            except IndexError:
                return 0
        return result.tuple()[0]
