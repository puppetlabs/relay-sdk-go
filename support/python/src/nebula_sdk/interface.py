from __future__ import annotations

import json
from typing import Any, Optional, Union

from .client import new_session
from .events import Events
from .outputs import Outputs
from .util import json_object_hook


class UnresolvableException(Exception):
    pass


class DynamicMetaclass(type):

    def __getattr__(self, name: str) -> Dynamic:
        return Dynamic(name)


class Dynamic(metaclass=DynamicMetaclass):
    """A query interface for inspecting a spec.

    This class allows arbitrary traversal that can be converted to a query to
    the metadata API.
    """

    def __init__(self, name: str, parent: Optional[Dynamic] = None) -> None:
        self._name = name
        self._parent = parent

    def __getattr__(self, name: str) -> Dynamic:
        return Dynamic(name, parent=self)

    def __str__(self) -> str:
        if self._parent is None:
            return self._name

        return '{0}[{1}]'.format(self._parent, json.dumps(self._name))


class Interface:

    def __init__(self, api_url: Optional[str] = None):
        self._client = new_session(api_url=api_url)

    def get(self, q: Optional[Union[Dynamic, str]] = None) -> Any:
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
    def events(self) -> Events:
        return Events(self._client)

    @property
    def outputs(self) -> Outputs:
        return Outputs(self._client)
