import json
from urllib.parse import quote

from .util import JSONEncoder


class Outputs(object):

    def __init__(self, client):
        self._client = client

    def set(self, name, value):
        r = self._client.put(
            'http+api://api/outputs/{0}'.format(quote(name)),
            data=json.dumps(value, cls=JSONEncoder),
            headers={'content-type': 'application/json'},
        )
        r.raise_for_status()
