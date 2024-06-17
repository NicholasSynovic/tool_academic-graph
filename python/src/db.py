from typing import List

from sqlalchemy import (
    Boolean,
    Column,
    CursorResult,
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

        worksTableSchema: Table = Table(
            "works",
            metadata,
            Column("oaid", String),
            Column("doi", String),
            Column("authorship_count", Integer),
            Column("cited_by_count", Integer),
            Column("concept_count", Integer),
            Column("cr_type", String),
            Column("created", DateTime),
            Column("distinct_country_count", Integer),
            Column("filepath", String),
            Column("grant_count", Integer),
            Column("institution_count", Integer),
            Column("is_paratext", Boolean),
            Column("is_retracted", Boolean),
            Column("keyword_count", Integer),
            Column("language", String),
            Column("license", String),
            Column("oa_type", String),
            Column("publication_location_count", Integer),
            Column("published", DateTime),
            Column("sustainable_development_goal_count", Integer),
            Column("title", String),
            Column("topic_count", Integer),
            Column("updated", DateTime),
            PrimaryKeyConstraint("oaid"),
        )

        metadata.create_all(bind=self.dbConn)

    def getLargestCitesID(self) -> int:
        sql: TextClause = text(text="SELECT id FROM cites ORDER BY id DESC LIMIT 1;")
        with self.dbConn.connect() as connection:
            try:
                result: Row = list(connection.execute(statement=sql))[0]
            except IndexError:
                return 0
        return result.tuple()[0]

    def getWorkCount(self) -> int:
        sql: TextClause = text(text="SELECT COUNT(oa_id) FROM works;")
        with self.dbConn.connect() as connection:
            result: CursorResult = connection.execute(statement=sql)
        return result.fetchone().tuple()[0]
