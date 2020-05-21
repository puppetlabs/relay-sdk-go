from __future__ import annotations

import asyncio
import base64
import concurrent.futures
import datetime
import functools
import inspect
import json
import signal
import weakref
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


TerminationEvent = Callable[[], Awaitable[None]]


class TerminationPolicy(Protocol):

    async def attach(self) -> Optional[TerminationEvent]: ...


class NoTerminationPolicy(TerminationPolicy):

    async def attach(self) -> Optional[TerminationEvent]:
        return None


class SoftTerminationPolicy(TerminationPolicy):

    _tasks: weakref.WeakKeyDictionary[asyncio.Task[Any], asyncio.Event]
    _timeout_sec: Optional[float]

    def __init__(self, *, timeout_sec: Optional[float] = None):
        self._tasks = weakref.WeakKeyDictionary()
        self._timeout_sec = timeout_sec

    async def terminate_task(self, task: asyncio.Task[Any]) -> None:
        event = self._tasks.get(task)
        if event is not None:
            event.set()

        if task.done():
            return

        if self._timeout_sec is not None:
            try:
                await asyncio.wait_for(task, self._timeout_sec)
            except asyncio.TimeoutError:
                task.cancel()

        await task

    def terminate_all(self) -> None:
        futs = [
            asyncio.run_coroutine_threadsafe(
                self.terminate_task(t),
                t.get_loop(),
            ) for t in self._tasks
        ]
        concurrent.futures.wait(futs)

    async def attach(self) -> Optional[TerminationEvent]:
        task = asyncio.current_task()
        assert task is not None

        try:
            event = self._tasks[task]
        except KeyError:
            event = asyncio.Event()
            self._tasks[task] = event

        async def wait() -> None:
            await event.wait()

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

    async def attach(self) -> Optional[TerminationEvent]:
        loop = asyncio.get_running_loop()
        task = asyncio.current_task()
        assert task is not None

        term = self._delegate.terminate_task(task)
        for sig in self._signals:
            loop.add_signal_handler(sig, lambda: loop.create_task(term))

        return await self._delegate.attach()
