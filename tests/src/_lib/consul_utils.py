import pytest

from contextlib import contextmanager

from . import docker_utils


CONSUL_IMAGE = "hashicorp/consul:latest"


@contextmanager
def container(docker_client, *, name, network, report_dir):
    print("Pulling the Consul image")
    try:
        docker_client.images.pull(CONSUL_IMAGE)

    except Exception as err:
        pytest.fail(f"{err}", pytrace=False)

    print(f"Creating Container: {name}")
    container = docker_client.containers.run(
        image=CONSUL_IMAGE,
        name=name,
        network=network.name,
        hostname=name,
        ports={
            "8500/tcp": 8500,
        },
        healthcheck={
            "test": ["CMD-SHELL", "consul info | awk '/health_score/{if ($3 >=1) exit 1; else exit 0}'"],
            "interval": 5 * 1000 * 1000 * 1000,  # 5 seconds (in nanoseconds)
            "timeout": 3 * 1000 * 1000 * 1000,  # 60 seconds (in nanoseconds)
            "retries": 3,
            "start_period": 10 * 1000 * 1000 * 1000,  # 10 seconds (in nanoseconds)
        },
        detach=True,
    )

    try:
        print(f"Waiting for healthcheck: {name}")
        docker_utils.wait_for_healthcheck(container)

    except RuntimeError as err:
        docker_utils.teardown_container(container, report_dir)
        pytest.fail(f"{err}", pytrace=False)

    yield

    docker_utils.teardown_container(container, report_dir)
