import json
from typing import Any
from urllib.parse import quote

from requests import Session

from .util import JSONEncoder


class Outputs:

    def __init__(self, client: Session) -> None:
        self._client = client

    def set(self, name: str, value: Any) -> None:
        r = self._client.put(
            'http+api://api/outputs/{0}'.format(quote(name)),
            data=json.dumps(value, cls=JSONEncoder),
            headers={'content-type': 'application/json'},
        )
        r.raise_for_status()
