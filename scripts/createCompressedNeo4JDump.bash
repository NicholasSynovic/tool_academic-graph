#!/bin/bash

source optparse.bash
optparse.define short=i long=input desc="The Neo4J database to dump" variable=db default="neo4j"
source $( optparse.build )

sudo neo4j-admin database dump $db --to-path .

tar -czvf $db.dump.tar.gz $db.dump
