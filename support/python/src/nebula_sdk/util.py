import asyncio
import base64
import datetime
import functools
import inspect
import json
import signal
from typing import (Any, Awaitable, Callable, Iterable, Mapping, Optional,
                    Protocol, Union)


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


def is_async_callable(obj: Any) -> bool:
    if not callable(obj):
        return False

    return (
        inspect.iscoroutinefunction(obj) or
        inspect.iscoroutinefunction(obj.__call__)
    )


class TerminationPolicy(Protocol):

    def apply(self) -> Optional[Callable[..., Awaitable[None]]]: ...


class NoTerminationPolicy(TerminationPolicy):

    def apply(self) -> Optional[Callable[..., Awaitable[None]]]:
        return None


class SoftTerminationPolicy(TerminationPolicy):

    _event: asyncio.Event
    _timeout_sec: Optional[float]

    def __init__(self, *, timeout_sec: Optional[float] = None):
        self._event = asyncio.Event()
        self._timeout_sec = timeout_sec

    async def terminate(self) -> None:
        self._event.set()

        termination_task = asyncio.current_task()
        other_tasks = asyncio.gather(*filter(
            lambda t: t != termination_task,
            asyncio.all_tasks(),
        ), return_exceptions=True)

        if self._timeout_sec is not None:
            try:
                await asyncio.wait_for(other_tasks, self._timeout_sec)
            except asyncio.TimeoutError:
                other_tasks.cancel()

        await other_tasks

    def apply(self) -> Optional[Callable[..., Awaitable[None]]]:
        async def wait() -> None:
            await self._event.wait()

        return wait


class SignalTerminationPolicy(TerminationPolicy):

    _signals: Iterable[signal.Signals]
    _delegate: SoftTerminationPolicy

    def __init__(self, *,
                 signals: Optional[Iterable[signal.Signals]] = None,
                 timeout_sec: Optional[float] = None):
        if signals is None:
            signals = [signal.SIGINT, signal.SIGTERM]

        self._signals = signals
        self._delegate = SoftTerminationPolicy(timeout_sec=timeout_sec)

    def apply(self) -> Optional[Callable[..., Awaitable[None]]]:
        loop = asyncio.get_event_loop()
        wait = self._delegate.apply()

        for sig in self._signals:
            loop.add_signal_handler(
                sig,
                lambda: loop.create_task(self._delegate.terminate()),
            )

        return wait
