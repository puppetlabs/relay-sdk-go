import datetime
import json

import pytest
from nebula_sdk.util import JSONEncoder, json_object_hook


@pytest.mark.parametrize(
    'test_input, expected',
    [
        (
            {'foo': b'\x90'},
            r'{"foo": {"$encoding": "base64", "data": "kA=="}}',
        ),
        (range(0, 5), r'[0, 1, 2, 3, 4]'),
        (
            datetime.datetime(2020, 1, 2, 3, 4, 5, 0, datetime.timezone.utc),
            r'"2020-01-02T03:04:05+00:00"',
        ),
    ],
)
def test_json_encoding(test_input, expected):
    assert json.dumps(test_input, cls=JSONEncoder, sort_keys=True) == expected


def test_json_decoding():
    assert json.loads(
        r'{"$encoding": "base64", "data": "kA=="}',
        object_hook=json_object_hook,
    ) == b'\x90'
