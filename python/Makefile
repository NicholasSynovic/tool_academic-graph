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

	poetry build
	pip install dist/*.tar.gz


create-dev:
	rm -rf env bin
	python3.10 -m venv env
	( \
		. env/bin/activate; \
		pip install -r requirements.txt; \
		poetry install; \
		deactivate; \
	)
