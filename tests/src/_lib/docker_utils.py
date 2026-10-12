from time import monotonic, sleep

import requests


def cleanup(docker_client):
    print("Cleaning up old containers")
    for container in docker_client.containers.list(all=True):
        if container.name.startswith("test-flowg-"):
            container.remove(force=True)

    print("Cleaning up old volumes")
    for volume in docker_client.volumes.list():
        if volume.name.startswith("test-flowg-"):
            volume.remove(force=True)

    print("Cleaning up old networks")
    for network in docker_client.networks.list():
        if network.name.startswith("test-flowg"):
            network.remove()


def wait_for_healthcheck(container):
    while True:
        container.reload()

        if container.health == "healthy":
            break

        elif container.health == "unhealthy":
            raise RuntimeError(f"Node {container.name} was not healthy")

        sleep(1)


def wait_for_http(container, url, *, timeout=60, verify=True):
    """
    Wait until the service in the container answers HTTP requests.
    Any HTTP response (whatever its status code) means the server is up.
    """

    deadline = monotonic() + timeout

    while True:
        try:
            requests.get(url, timeout=1, verify=verify)

        except requests.RequestException as err:
            container.reload()
            if container.status != "running":
                raise RuntimeError(
                    f"Container {container.name} is not running ({container.status})"
                ) from err

            if monotonic() >= deadline:
                raise RuntimeError(
                    f"Container {container.name} not ready after {timeout}s: {err}"
                ) from err

            sleep(1)

        else:
            break


def teardown_container(container, report_dir):
    print(f"Stopping container: {container.name}")
    container.stop()

    print(f"Writing logs: {container.name}")
    with open(report_dir / f"docker-{container.name}.log", "wb") as f:
        for data in container.logs(stream=True):
            f.write(data)

    print(f"Removing container: {container.name}")
    container.remove(force=True)
