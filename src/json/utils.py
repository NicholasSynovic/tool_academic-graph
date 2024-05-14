from json import JSONDecodeError, loads
from typing import List

from progress.bar import Bar


def createJSON(data: List[str]) -> List[dict]:
    jsonObjects: List[dict] = []

    with Bar("Converting JSON strings into JSON objects...", max=len(data)) as bar:
        datum: str
        for datum in data:
            try:
                json: dict = loads(s=datum)
            except JSONDecodeError:
                bar.next()
                continue

            jsonObjects.append(json)
            bar.next()

    return jsonObjects
