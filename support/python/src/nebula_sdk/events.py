"Generating events for the service to act on"
import json
from typing import Any, Mapping

from requests import Session

from .util import JSONEncoder


class Events:
    """Class for generating event payloads to the service API"""

    def __init__(self, client: Session) -> None:
        self._client = client

    def emit(self, data: Mapping[str, Any]) -> None:
        """Create and send an event to the service.

        Use this from a Trigger handler to start the workflow
        associated with the trigger.

        Accepts a Mapping of data fields that will bind
        names of workflow parameters to values from the
        request payload that came in to the Trigger handler"""
        r = self._client.post(
            'http+api://api/events',
            data=json.dumps({'data': data}, cls=JSONEncoder),
            headers={'content-type': 'application/json'},
        )
        r.raise_for_status()
