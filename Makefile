build:
	poetry build
	pip install dist/*.tar.gz

	(\
		cd extractors; \
		go build; \
	)

create-dev:
	rm -rf env
	python3.10 -m venv env
	( \
		. env/bin/activate; \
		pip install -r requirements.txt; \
		poetry install; \
		deactivate; \
	)
