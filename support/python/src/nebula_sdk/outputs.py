"Outputs allow a step to write data that will be available for later steps"
import json
from typing import Any
from urllib.parse import quote

from requests import Session

from .util import JSONEncoder


class Outputs:
    "Use Outputs to write keys to the metadata service that later steps can retrieve"
    def __init__(self, client: Session) -> None:
        self._client = client

    def set(self, name: str, value: Any) -> None:
        """Set writes a value of a given name to the service

        Args:
            name: a string containing the key name to write
            value: the value of the key to be written. Can
                be any data type that can serialize to JSON
        """
        r = self._client.put(
            'http+api://api/outputs/{0}'.format(quote(name)),
            data=json.dumps(value, cls=JSONEncoder),
            headers={'content-type': 'application/json'},
        )
        r.raise_for_status()
