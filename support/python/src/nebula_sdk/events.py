import json
from typing import Any, Mapping

from requests import Session

from .util import JSONEncoder


class Events:

    def __init__(self, client: Session) -> None:
        self._client = client

    def emit(self, data: Mapping[str, Any]) -> None:
        r = self._client.post(
            'http+api://api/events',
            data=json.dumps({'data': data}, cls=JSONEncoder),
            headers={'content-type': 'application/json'},
        )
        r.raise_for_status()
