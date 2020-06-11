Relay Python SDK
================

This is the Python SDK for use with the internal `Relay <https://relay.sh>`_ service APIs.
The SDK requires Python 3.8.

It is intended for use by integration authors who are building containers to run 
inside the service. For running workflows and interacting with the user-facing
service APIs, use the `Relay CLI <https://github.com/puppetlabs/relay/>`_.

The API documentation is auto-generated from the source code. Here are some
higher-level examples that show the main SDK classes that integration authors
will interact with.


Installation
------------

The SDK is available to install via pip:

.. code-block:: console

  pip --no-cache-dir install "https://packages.nebula.puppet.net/sdk/support/python/v1/nebula_sdk-1-py3-none-any.whl"

If you use the `relaysh/core:latest-python <https://hub.docker.com/r/relaysh/core/tags>`_ container image as your base
image, it'll be pre-installed.

Usage
-----
The main purpose of the SDK is to provide helpers for interacting with Relay's
metadata service. Each container that runs in Relay has access to this service,
which allows the container to read and write key-value data, emit events, and
generate logs.

Accessing Data from the step spec
---------------------------------

The `Interface class <./reference.html#module-nebula_sdk.interface>`_ is the primary way to interact with the service.
Import it and instantiate an object, then call methods on that object to access metadata,
which comes from the ``spec`` section of the step and global Connection information.
The ``Dynamic`` class provides syntactic sugar for getting data like connection credentials, 
workflow-specific parameters, and secrets. It represents nested data structures as dot-separated
method accessors.

.. code-block:: python

  from nebula_sdk import Interface, Dynamic as D

  relay = Interface()
  azuresecret = relay.get(D.azure.connection.secret) # using Dynamic
  azureclient = relay.get('azure["connection"]["clientID"]') # same as above
  secret = relay.get(D.mysecret)
  relay.outputs.set("outputkey","This will be the value of outputkey")

Webhook Triggers
----------------

The `WebhookServer class <./reference.html#module-nebula_sdk.webhook>`_ provides a
helper that sets up a webserver to handle incoming requests for Trigger actions. 

This example, from the `Docker Hub integration <https://github.com/relay-integrations/relay-dockerhub/>`_, makes use of
the Interface class to access the ``events.emit`` method, which will cause
the workflow associated with this trigger to be run with the inline mapping
of workflow parameters to values extracted from the webhook payload.

The WebhookServer class can run any WSGI_ or ASGI_ application passed to it. The
integrations the Relay team develops internally use the Quart_ web app framework.

.. _WSGI: https://www.python.org/dev/peps/pep-3333/
.. _ASGI: https://asgi.readthedocs.io/en/latest/specs/main.html
.. _Quart: https://pgjones.gitlab.io/quart/index.html

.. code-block:: python

  from nebula_sdk import Interface, WebhookServer
  from quart import Quart, request, jsonify, make_response

  relay = Interface()
  app = Quart('image-pushed')

  @app.route('/', methods=['POST'])
  async def handler():
      event_payload = await request.get_json()

      if event_payload is None:
          return await make_response(jsonify(message='not a valid Docker Hub event'), 400)

      pd = event_payload['push_data']
      rd = event_payload['repository']

      relay.events.emit({
          'pushedAt': pd['pushed_at'],
          'pusher': pd['pusher'],
          'tag': pd['tag'],
          'name': rd['repo_name']
      })

      return await make_response(jsonify(message='success'), 200)


  if __name__ == '__main__':
      WebhookServer(app).serve_forever()

