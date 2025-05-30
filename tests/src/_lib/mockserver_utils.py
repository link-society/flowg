import pytest

from contextlib import contextmanager
from time import sleep

import requests

from . import docker_utils


@contextmanager
def container(
    docker_client,
    *,
    name,
    network,
    config_dir,
    report_dir,
):
    print(f"Creating container: {name}")
    container = docker_client.containers.run(
        image="mockserver/mockserver:latest",
        name=name,
        network=network.name,
        hostname=name,
        environment={
            "MOCKSERVER_INITIALIZATION_JSON_PATH": "/config.json",
        },
        ports={
            "1080/tcp": 1080,
        },
        volumes={
            (config_dir / "mockserver-config.json").absolute().as_posix(): {
                "bind": "/config.json",
                "mode": "ro",
            }
        },
        detach=True,
    )

    try:
        print(f"Waiting for healthcheck: {name}")
        wait_for_healthcheck()

    except RuntimeError as err:
        docker_utils.teardown_container(container, report_dir)
        pytest.fail(f"{err}", pytrace=False)

    yield

    docker_utils.teardown_container(container, report_dir)


def wait_for_healthcheck():
    attempt = 0
    max_attempts = 10

    while True:
        try:
            response = requests.put(
                "http://localhost:1080/mockserver/status",
                timeout=1,
            )
            response.raise_for_status()

        except Exception:
            attempt += 1
            if attempt >= max_attempts:
                raise

            sleep(1)

        else:
            break
