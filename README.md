# Academic Graph Generator

> Converts OpenAlex JSON documents to SQLite3 and Neo4J databases for analysis

## Table of Contents

- [Academic Graph Generator](#academic-graph-generator)
  - [Table of Contents](#table-of-contents)
  - [About](#about)
  - [How to Install](#how-to-install)
    - [Dependencies](#dependencies)
    - [Installation steps](#installation-steps)
  - [How to Run](#how-to-run)
    - [1. Extract Relevant Data From OpenAlex JSON Files](#1-extract-relevant-data-from-openalex-json-files)
    - [2. Create SQLite3 Database](#2-create-sqlite3-database)

## About

Academic tradition supports reusing and citing others published work in order to
further science. [OpenAlex](https://openalex.org) provides metadata of over 250
million documents within the academic network. This tool aims to take these
downloads and parse them into SQLite3 and Neo4J databases to be used as proper
datasources within future applications.

## How to Install

This project was tested on x86-64 Linux machines

### Dependencies

- `Python 3.10`
- `go 1.22.3`

### Installation steps

1. `make create-dev`
1. `source env/bin/activate`
1. `make build`

## How to Run

There are several applications contained within this project. The order in which
you should run the applications are as follows:

### 1. Extract Relevant Data From OpenAlex JSON Files

OpenAlex JSON files contains metadata that is irrelevant for the current scope
of the project. We thus need to extract and keep only the relevant information.
To do so, run the `extractors/extract_json` utility.

```shell
./extract_json -h

Usage of ./extract_json:
  -i string
        Path to OpenAlex "Works" JSON Lines file
  -o string
        Path to output JSON file
```

### 2. Create SQLite3 Database

By leveraging SQLite3, we can query data more efficiently than with standard
JSON. To convert the output JSON to a SQLite3 database, run
`oag/createDatabase.py`

```shell
python createDatabase.py --help

Usage: createDatabase.py [OPTIONS]

Options:
  -i, --input PATH   Path to JSON file  [required]
  -o, --output PATH  Path to output database  [required]
  --help             Show this message and exit.
```
