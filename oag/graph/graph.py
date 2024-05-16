from string import Template
from typing import List, Tuple

from neo4j import AsyncGraphDatabase, Driver, GraphDatabase, ManagedTransaction
from pandas import DataFrame, Series
from progress.bar import Bar


class Neo4J:
    def __init__(self, uri: str, username: str, password: str) -> None:
        self.uri: str = uri
        self.auth: Tuple[str, str] = (username, password)
        self.driver: Driver = GraphDatabase.driver(
            uri=self.uri,
            auth=self.auth,
        )

    @staticmethod
    async def _async_ExecuteQuery(tx: ManagedTransaction, query: str) -> None:
        await tx.run(query=query)

    @staticmethod
    def _executeQuery(tx: ManagedTransaction, query: str) -> None:
        tx.run(query=query)

    async def async_AddRelationships(
        self, df: DataFrame, bar: Bar, database: str = "neo4j"
    ) -> None:
        # https://community.neo4j.com/t/creating-relationship-over-several-millions-of-nodes/24390/3
        queryTemplate: Template = Template(
            template=r'MATCH (n:Work) WHERE n.oa_id = "${node1}" WITH n MATCH (m:Work) WHERE m.oa_id = "${node2}" CREATE (n)-[r:Cites]->(m)',
        )

        async with AsyncGraphDatabase.driver(self.uri, auth=self.auth) as driver:
            async with driver.session(database=database) as session:
                datum: Series
                for _, datum in df.iterrows():
                    query: str = queryTemplate.substitute(
                        node1=datum["work"],
                        node2=datum["reference"],
                    )
                    # queries.append(query)
                    await session.execute_write(self._async_ExecuteQuery, query)
                    bar.next()

    def addWorkNode(self, df: DataFrame) -> None:
        queries: List[str] = []

        queryTemplate: Template = Template(
            template=r'(n${id}:Work {oa_id: "${oa_id}", doi: "${doi}", title: "${title}", type: "${type}"})',
        )

        datum: Series
        idx: int
        for idx, datum in df.iterrows():
            query: str = queryTemplate.substitute(
                id=idx,
                type=datum["oa_type"],
                oa_id=datum["oa_id"],
                doi=datum["doi"],
                title=datum["title"],
            )
            queries.append(query)

        with self.driver.session() as session:
            query: str = r"CREATE " + ", ".join(queries)
            session.execute_write(self._executeQuery, query)

    def createWorkNodeIndex(self) -> None:
        with self.driver.session() as session:
            query: str = f"CREATE TEXT INDEX workNodes FOR (n:Work) ON (n.oa_id)"
            session.execute_write(self._executeQuery, query)

    def addRelationship(self, df: DataFrame) -> None:
        # https://community.neo4j.com/t/creating-relationship-over-several-millions-of-nodes/24390/3
        queries: List[str] = []
        queryTemplate: Template = Template(
            template=r'MATCH (n:Work) WHERE n.oa_id = "${node1}" WITH n MATCH (m:Work) WHERE m.oa_id = "${node2}" CREATE (n)-[r:Cites]->(m)',
        )

        datum: Series
        for _, datum in df.iterrows():
            query: str = queryTemplate.substitute(
                node1=datum["work"],
                node2=datum["reference"],
            )
            queries.append(query)

        with self.driver.session() as session:
            query: str
            for query in queries:
                session.execute_write(self._executeQuery, query)
