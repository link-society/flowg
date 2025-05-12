import pytest

from contextlib import contextmanager
import requests

from . import docker_utils


@contextmanager
def network(docker_client, *, name):
    print(f"Creating network: {name}")
    network = docker_client.networks.create(name=name, driver="bridge")
    yield network
    print(f"Removing network: {name}")
    network.remove()


@contextmanager
def volume(docker_client, *, name):
    print(f"Creating volume: {name}")
    volume = docker_client.volumes.create(name=name)
    yield volume
    print(f"Removing volume: {name}")
    volume.remove()


@contextmanager
def container(
    docker_client,
    *,
    name,
    network,
    volume,
    image,
    environment,
    ports,
    report_dir,
):
    env = {
        "FLOWG_SECRET_KEY": "s3cr3!",
        "FLOWG_CLUSTER_NODE_ID": name,
        "FLOWG_AUTH_DIR": "/data/auth",
        "FLOWG_CONFIG_DIR": "/data/config",
        "FLOWG_LOG_DIR": "/data/logs",
    }
    env.update(environment)

    print(f"Creating container: {name}")
    container = docker_client.containers.run(
        image=image,
        name=name,
        environment=env,
        network=network.name,
        hostname=name,
        ports=ports,
        volumes={
            volume.name: {"bind": "/data", "mode": "rw"}
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


def create_token(*, username, password):
    print(f"Creating {username} token")
    resp = requests.post(
        "http://localhost:5080/api/v1/auth/login",
        json={"username": username, "password": password},
    )
    resp.raise_for_status()
    data = resp.json()
    admin_jwt = data["token"]

    resp = requests.post(
        "http://localhost:5080/api/v1/token",
        headers={"Authorization": f"Bearer {admin_jwt}"},
    )
    resp.raise_for_status()
    data = resp.json()
    return data["token"]


def create_user(*, token, username, password):
    print(f"Creating {username} user")
    resp = requests.put(
        f"http://localhost:5080/api/v1/users/{username}",
        headers={"Authorization": f"Bearer {token}"},
        json={"password": password, "roles": []},
    )
    resp.raise_for_status()
