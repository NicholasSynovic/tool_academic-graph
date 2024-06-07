#!/bin/bash
source optparse.bash
optparse.define short=i long=input desc="Directory containing OpenAlex \"Works\" JSON documents" variable=inputPath
optparse.define short=s long=suffix desc="Suffix for output JSON files" variable=fileSuffix
source $( optparse.build )

absProgPath=$(readlink -f ../go/bin)
absInputPath=$(readlink -f $inputPath)

if [ ! -d $absInputPath ]; then
    echo "$absInputPath is a file. Expected a directory."
fi

cd ../go/bin
mkdir -p data

ls $absInputPath/part* | xargs -I % bash -c 'basename "$1"' _ % | xargs -I % bash -c '$2/ag-ce.bin -i $3/% -o data/ce_$4_%.json' _ % "$absProgPath" "$absInputPath" "$fileSuffix"
