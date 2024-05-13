from sqlalchemy import (
    Boolean,
    Column,
    Connection,
    DateTime,
    ForeignKeyConstraint,
    Integer,
    MetaData,
    PrimaryKeyConstraint,
    String,
    Table,
)


class Works:
    def __init__(self, dbConn: Connection) -> None:
        self.metadata: MetaData = MetaData()
        self.dbConn: Connection = dbConn

        self.tableSchema: Table = Table(
            "works",
            self.metadata,
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

        self.citesSchema: Table = Table(
            "relationship_cites",
            self.metadata,
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

        self.metadata.create_all(bind=self.dbConn)
