#!/bin/bash
source optparse.bash
optparse.define short=i long=input desc="Directory containing OpenAlex \"Works\" JSON documents" variable=inputPath
source $( optparse.build )

absProgPath=$(readlink -f ../go/bin)
absInputPath=$(readlink -f $inputPath)

if [ ! -d $absInputPath ]; then
    echo "$absInputPath is a file. Expected a directory."
fi

cd ../go/bin
mkdir -p data

ls $absInputPath/part* | xargs -I % basename % | parallel --bar -I {} $absProgPath/ag-ce.bin -i $absInputPath/{} -o data/ce_{}.json
