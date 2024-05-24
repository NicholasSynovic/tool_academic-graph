#!/bin/bash

source optparse.bash
optparse.define short=db long=database desc="The SQLite3 database to process" variable=db
optparse.define short=o long=output desc="The output GraphML file" variable=output default=graph.gml
source $( optparse.build )

./../bin/graphML-generator.bin -i $db -o $output

tar -czvf $output.tar.gz $output
