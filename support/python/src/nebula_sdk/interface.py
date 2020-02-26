import json

from .client import new_session
from .outputs import Outputs
from .util import json_object_hook


class UnresolvableException(Exception):
    pass


class DynamicMetaclass(type):

    def __getattr__(self, name):
        return Dynamic(name)


class Dynamic(metaclass=DynamicMetaclass):
    """A query interface for inspecting a spec.

    This class allows arbitrary traversal that can be converted to a query to
    the metadata API.
    """

    def __init__(self, name, parent=None):
        self._name = name
        self._parent = parent

    def __getattr__(self, name):
        return Dynamic(name, parent=self)

    def __str__(self):
        if self._parent is None:
            return self._name

        return '{0}[{1}]'.format(self._parent, json.dumps(self._name))


class Interface(object):

    def __init__(self, api_url=None):
        self._client = new_session(api_url=api_url)

    def get(self, q=None):
        params = {}
        if q is not None:
            params['q'] = str(q)

        r = self._client.get('http+api://api/spec', params=params)
        r.raise_for_status()

        data = json.loads(r.text, object_hook=json_object_hook)
        if not data['complete']:
            raise UnresolvableException()

        return data['value']

    @property
    def outputs(self):
        return Outputs(self._client)
