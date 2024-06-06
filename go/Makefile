build:
	( \
		mkdir bin >> /dev/null; \
		cd graphML_generator; \
		go build -o ../bin/graphML-generator.bin; \
	)

	( \
		mkdir bin >> /dev/null; \
		cd oa_extractor; \
		go build -o ../bin/oa-extractor.bin; \
	)

create-dev:
	go work sync
