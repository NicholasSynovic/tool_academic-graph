from sqlalchemy import (
    Boolean,
    Column,
    DateTime,
    Engine,
    ForeignKeyConstraint,
    Integer,
    MetaData,
    PrimaryKeyConstraint,
    String,
    Table,
)


class DB:
    def __init__(self, dbConn: Engine) -> None:
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
            "relationship_cites",
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

        metadata.create_all(bind=dbConn)
