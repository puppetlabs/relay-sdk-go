import pytest
from nebula_sdk.interface import Dynamic, Interface, UnresolvableException
from requests_mock import Adapter


@pytest.mark.parametrize(
    'test_input, expected',
    [
        (Dynamic.foo.bar, 'foo["bar"]'),
        (Dynamic.foo, 'foo'),
        (Dynamic.very.deep.nesting, 'very["deep"]["nesting"]'),
    ],
)
def test_dynamic(test_input: Dynamic, expected: str) -> bool:
    return str(test_input) == expected


class TestInterface:

    def test_all(self, requests_mock: Adapter) -> None:
        requests_mock.register_uri(
            'GET', 'http+api://api/spec',
            text=r'{"value": {"foo": "bar"}, "complete": true}',
        )
        assert Interface(api_url='http://api').get() == {'foo': 'bar'}

    def test_query(self, requests_mock: Adapter) -> None:
        requests_mock.register_uri(
            'GET', 'http+api://api/spec?q=foo',
            complete_qs=True,
            text=r'{"value": "bar", "complete": true}',
        )
        assert Interface(api_url='http://api').get(Dynamic.foo) \
            == 'bar'

    def test_incomplete(self, requests_mock: Adapter) -> None:
        requests_mock.register_uri(
            'GET', 'http+api://api/spec',
            text=r'''{
                "value": {"foo": {"$type": "Secret", "name": "foo"}},
                "complete": false}
            ''',
        )
        with pytest.raises(UnresolvableException):
            Interface(api_url='http://api').get()
