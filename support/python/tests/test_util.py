import asyncio
import datetime
import json
import signal
from typing import Any

import pytest
from nebula_sdk.util import (JSONEncoder, SignalTerminationPolicy,
                             is_async_callable, json_object_hook)


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
def test_json_encoding(test_input: Any, expected: Any) -> None:
    assert json.dumps(test_input, cls=JSONEncoder, sort_keys=True) == expected


def test_json_decoding() -> None:
    assert json.loads(
        r'{"$encoding": "base64", "data": "kA=="}',
        object_hook=json_object_hook,
    ) == b'\x90'


def test_is_async_callable() -> None:
    class AsyncCallable:

        async def __call__(self) -> None:
            pass

    class Callable:

        def __call__(self) -> None:
            pass

    async def async_callable() -> None:
        pass

    def sync_callable() -> None:
        pass

    assert is_async_callable(AsyncCallable())
    assert is_async_callable(async_callable)
    assert not is_async_callable(Callable)
    assert not is_async_callable(sync_callable)
    assert not is_async_callable(2)


class TestSignalTerminationPolicy:

    @pytest.mark.asyncio
    async def test_indefinite(self) -> None:
        pol = SignalTerminationPolicy(signals=[signal.SIGWINCH])

        ready_event = asyncio.Event()
        completion_event = asyncio.Event()

        async def run() -> None:
            waiter = pol.apply()
            assert waiter is not None
            ready_event.set()

            await waiter()
            completion_event.set()

        async def raise_signal() -> None:
            await ready_event.wait()
            signal.raise_signal(signal.SIGWINCH)

        await asyncio.gather(run(), raise_signal())
        await asyncio.wait_for(completion_event.wait(), 60)

    @pytest.mark.asyncio
    async def test_timeout(self) -> None:
        pol = SignalTerminationPolicy(
            signals=[signal.SIGWINCH],
            timeout_sec=0,  # Immediate timeout.
        )

        ready_event = asyncio.Event()

        async def run() -> None:
            waiter = pol.apply()
            assert waiter is not None
            ready_event.set()

            await waiter()

            # Instead of completing we'll force the policy to terminate us
            # early.
            await asyncio.sleep(30)

        async def raise_signal() -> None:
            await ready_event.wait()
            signal.raise_signal(signal.SIGWINCH)

        with pytest.raises(asyncio.CancelledError):
            await asyncio.gather(run(), raise_signal())
