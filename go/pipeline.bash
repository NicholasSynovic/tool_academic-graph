#!/bin/bash
source optparse.bash

optparse.define short=a long=author-input desc="A directory of OpenAlex Authors files to read" variable=authorInput

optparse.define short=b long=work-input desc="A directory of OpenAlex Works files to read" variable=workInput
source $( optparse.build )

ai=$(realpath $authorInput)
ao=$(realpath ./authors.json)
aro=$(realpath ./authorship-relationships.json)

wi=$(realpath $workInput)
wo=$(realpath ./works.json)
wIndex=$(realpath ./works-index.json)
cro=$(realpath ./cites-relationships.json)

aoProg=$(realpath ./cmd/authorObjects/author-objects.go)
arProg=$(realpath ./cmd/authorshipRelationships/authorship-relationships.go)

woProg=$(realpath ./cmd/workObjects/work-objects.go)
croProg=$(realpath ./cmd/citesRelationships/cites-relationships.go)

wiProg=$(realpath ./cmd/workIndex/work-index.go)

# Author programs
go run $aoProg -i $ai -o $ao
go run $arProg -i $wi -o $aro

# Works programs
go run $woProg -i $wi -o $wo
go run $croProg -i $wi -o $cro

# Indices
go run $wiProg -i $wi -o $wIndex
