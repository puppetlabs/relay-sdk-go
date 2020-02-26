import os
from urllib.parse import urljoin, urlsplit, urlunsplit

from requests import Session
from requests.adapters import HTTPAdapter


class MetadataAPIAdapter(HTTPAdapter):

    def __init__(self, base_url=None):
        if base_url is None:
            base_url = os.environ['METADATA_API_URL']

        self._base_url = base_url

        super(MetadataAPIAdapter, self).__init__()

    def send(self, request, stream=False, timeout=None, verify=True, cert=None,
             proxies=None):
        (_, _, path, query, fragment) = urlsplit(request.url)
        request.prepare_url(
            urljoin(
                self._base_url,
                urlunsplit(('', '', path, query, fragment)),
            ),
            {},
        )

        return super(MetadataAPIAdapter, self).send(
            request,
            stream=stream,
            timeout=timeout,
            verify=verify,
            cert=cert,
            proxies=proxies,
        )


def new_session(api_url=None):
    sess = Session()
    sess.mount('http+api://', MetadataAPIAdapter(base_url=api_url))
    return sess
