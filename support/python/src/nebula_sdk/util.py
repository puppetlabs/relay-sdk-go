import base64
import datetime
import functools
import json
from typing import Any, Callable, Mapping, Union


def json_object_hook(dct: Mapping[str, Any]) -> Any:
    if '$encoding' in dct:
        try:
            decoder: Callable[[str], str] = {
                'base64': base64.standard_b64decode,
                '': lambda data: data,
            }[dct['$encoding']]

            return decoder(dct['data'])
        except KeyError:
            # Either dct does not contain data or has an encoding that we can't
            # handle.
            pass
    return dct


class JSONEncoder(json.JSONEncoder):

    @functools.singledispatchmethod
    def default(self, obj: Any) -> Any:  # type: ignore[override]
        try:
            it = iter(obj)
        except TypeError:
            pass
        else:
            return list(it)

        return super(JSONEncoder, self).default(obj)

    @default.register
    def _datetime(self, obj: datetime.datetime) -> str:
        return obj.isoformat()

    @default.register
    def _bytes(self, obj: bytes) -> Union[str, Mapping[str, Any]]:
        try:
            return obj.decode('utf-8')
        except UnicodeDecodeError:
            return {
                '$encoding': 'base64',
                'data': base64.standard_b64encode(obj),
            }
