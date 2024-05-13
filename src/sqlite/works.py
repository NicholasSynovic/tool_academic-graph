from sqlalchemy import Boolean, Column, Connection, DateTime, MetaData, String, Table


def createSchema(dbConn: Connection) -> None:
    metadata: MetaData = MetaData()

    schema: Table = Table(
        "works",
        metadata,
        Column("oa_id", String, primary_key=True),
        Column("doi", String),
        Column("title", String),
        Column("paratext", Boolean),
        Column("retracted", Boolean),
        Column("published", DateTime),
        Column("oa_type", String),
        Column("cf_type", String),
    )

    metadata.create_all(bind=dbConn)
