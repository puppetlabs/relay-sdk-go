from __future__ import annotations

import asyncio
import os
from typing import TYPE_CHECKING, Optional, Union, cast

from asgiref.wsgi import WsgiToAsgi
from hypercorn.asyncio.run import worker_serve
from hypercorn.config import Config, Sockets
from hypercorn.typing import ASGIFramework

from .util import SignalTerminationPolicy, TerminationPolicy, is_async_callable

if TYPE_CHECKING:
    from wsgiref.types import WSGIApplication

DEFAULT_TERMINATION_POLICY = SignalTerminationPolicy(timeout_sec=25)
DEFAULT_PORT = 8080


class PortParseError(Exception):
    pass


class WebhookServer:

    _app: ASGIFramework
    _config: Config
    _sockets: Sockets
    _termination_policy: TerminationPolicy
    _task: Optional[asyncio.Task[None]]

    def __init__(self, app: Union[ASGIFramework, 'WSGIApplication'], *,
                 termination_policy: Optional[TerminationPolicy] = None,
                 port: Optional[int] = None):
        if not is_async_callable(app):
            app = WsgiToAsgi(app)

        if termination_policy is None:
            termination_policy = DEFAULT_TERMINATION_POLICY

        if port is None:
            try:
                port = int(os.environ['PORT'])
            except KeyError:
                port = DEFAULT_PORT
            except ValueError:
                raise PortParseError()

        self._app = cast(ASGIFramework, app)
        self._config = Config()
        self._config.bind = [f'0.0.0.0:{port}']
        self._sockets = self._config.create_sockets()
        self._termination_policy = termination_policy
        self._task = None

    @property
    def port(self) -> int:
        return self._sockets.insecure_sockets[0].getsockname()[1]

    def listening(self) -> bool:
        return self._task is not None and not self._task.done()

    async def serve(self) -> None:
        if self._task is None:
            loop = asyncio.get_running_loop()
            shutdown_trigger = await self._termination_policy.attach()

            self._task = loop.create_task(worker_serve(
                self._app, self._config,
                sockets=self._sockets,
                shutdown_trigger=shutdown_trigger,
            ))

        await self._task

    async def _serve_wait(self) -> None:
        current_task = asyncio.current_task()
        assert current_task is not None

        await self.serve()
        await asyncio.gather(
            *filter(lambda t: t != current_task, asyncio.all_tasks()),
            return_exceptions=True,
        )

    def serve_forever(self) -> None:
        asyncio.run(self._serve_wait())
