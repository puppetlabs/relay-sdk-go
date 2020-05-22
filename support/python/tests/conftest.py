import asyncio
from concurrent import futures
from typing import Any, Dict, Iterable

import pytest


def event_loop_exception_handler(loop: asyncio.AbstractEventLoop,
                                 context: Dict[str, Any]) -> None:
    # Ignore dropped tasks in test.
    if loop.is_closed() and isinstance(context.get('task'), asyncio.Task):
        return

    loop.default_exception_handler(context)


@pytest.fixture
def event_loop() -> Iterable[asyncio.AbstractEventLoop]:
    loop = asyncio.new_event_loop()
    loop.set_exception_handler(event_loop_exception_handler)
    yield loop
    loop.close()


@pytest.fixture
def thread_pool_executor() -> Iterable[futures.Executor]:
    with futures.ThreadPoolExecutor() as pool:
        yield pool
