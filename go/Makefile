build:
	( \
		mkdir bin >> /dev/null; \
		cd cmd/citation_extraction; \
		go build -o ../../bin/ag-ce.bin; \
	)

	( \
		mkdir bin >> /dev/null; \
		cd cmd/document_extraction; \
		go build -o ../../bin/ag-de.bin; \
	)

	( \
		mkdir bin >> /dev/null; \
		cd cmd/oaid_doi_mapper; \
		go build -o ../../bin/ag-odm.bin; \
	)

create-dev:
	go work sync
