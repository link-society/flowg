import pytest
import os
import time
import requests
from pathlib import Path

from contextlib import contextmanager

from . import docker_utils


COREDNS_IMAGE = "coredns/coredns:latest"

@contextmanager
def container(docker_client, *, name, network, report_dir):
    print("Pulling the CoreDNS image")
    try:
        docker_client.images.pull(COREDNS_IMAGE)

    except Exception as err:
        pytest.fail(f"{err}", pytrace=False)

    local_corefile_path = Path(__file__).parent / "Corefile"
    container_corefile_path = "/etc/coredns"

    command = ["-conf", f"{container_corefile_path}/Corefile"]

    try:
        print(f"Creating Container: {name}")
        container = docker_client.containers.run(
            image=COREDNS_IMAGE,
            name=name,
            network=network.name,
            hostname=name,
            ports={
                "53/tcp": 54,
                "53/udp": 54,
                "8080/tcp":8080
            },
            volumes={
                str(local_corefile_path.parent): {
                    "bind": container_corefile_path,
                    "mode": "ro"
                }
            },
            command=command,
            detach=True,
        )

    except Exception as err:
        pytest.fail(f"{err}", pytrace=True)

    try:
        print(f"Waiting for healthcheck: {name}")
        wait_for_healthcheck(container)

    except RuntimeError as err:
        print(f"Container '{name}' healthcheck failed. Logs:")
        print(container.logs().decode('utf-8'))
        docker_utils.teardown_container(container, report_dir)
        pytest.fail(f"{err}", pytrace=False)

    yield

    docker_utils.teardown_container(container, report_dir)


def wait_for_healthcheck(container):
    attempts = 10
    interval = 1

    while True:
        try:
            resp = requests.get("http://localhost:8080/health")
            resp.raise_for_status()

            print(f"Health check for {container.name} successful!")
            return

        except Exception as err:
            if attempts == 0:
                raise TimeoutError(f"{container.name} not healthy: {err}")

            else:
                attempts -= 1
                time.sleep(interval)