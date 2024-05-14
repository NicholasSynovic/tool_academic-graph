from string import Template
from typing import List, Tuple

from neo4j import Driver, GraphDatabase, ManagedTransaction
from pandas import DataFrame, Series


class Neo4J:
    def __init__(self, uri: str, username: str, password: str) -> None:
        self.uri: str = uri
        self.auth: Tuple[str, str] = (username, password)
        self.driver: Driver = GraphDatabase.driver(
            uri=self.uri,
            auth=self.auth,
        )

    @staticmethod
    def _createNode(tx: ManagedTransaction, query: str) -> None:
        tx.run(query=query)

    def addNode(self, df: DataFrame) -> None:
        queries: List[str] = []
        queryTemplate: Template = Template(
            template=r'CREATE (n:${type} {oa_id: "${oa_id}", doi: "${doi}", title: "${title}"})',
        )

        datum: Series
        for _, datum in df.iterrows():
            query: str = queryTemplate.substitute(
                type=datum["oa_type"],
                oa_id=datum["oa_id"],
                doi=datum["doi"],
                title=datum["title"],
            )
            queries.append(query)

        with self.driver.session() as session:
            query: str
            for query in queries:
                session.execute_write(self._createNode, query)
