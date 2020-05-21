from typing import Container, Text, Tuple, Union

HTTPTimeout = Union[float, Tuple[float, float], Tuple[float, None]]
HTTPClientCertificate = Union[bytes, Text, Container[Union[bytes, Text]]]
