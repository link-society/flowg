from time import sleep


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


def teardown_container(container, report_dir):
    print(f"Stopping container: {container.name}")
    container.stop()

    print(f"Writing logs: {container.name}")
    with open(report_dir / f"docker-{container.name}.log", "wb") as f:
        for data in container.logs(stream=True):
            f.write(data)

    print(f"Removing container: {container.name}")
    container.remove(force=True)
