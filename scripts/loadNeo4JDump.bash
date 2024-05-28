#!/bin/bash

source optparse.bash
optparse.define short=i long=input desc="The Neo4J database to load" variable=db default="neo4j"
optparse.define short=d long=dump desc="Path to directory with the Neo4J dump to load" variable=dumpPath
optparse.define long=delete desc="Path to directory with the Neo4J dump to load" variable=dumpPath
source $( optparse.build )

sudo neo4j-admin database load $db --verbose --from-path $dumpPath
