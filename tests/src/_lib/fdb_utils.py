import pytest

from contextlib import contextmanager
from time import sleep

from . import docker_utils


# The FoundationDB server image. Must match the client library (libfdb_c) bundled
# in the FlowG image (see docker/flowg.dockerfile) and the pinned API version.
FDB_IMAGE = "foundationdb/foundationdb:7.3.77"

# Path of the cluster file inside the FoundationDB container. FlowG nodes are
# given a copy of this file to connect to the cluster.
CLUSTER_FILE = "/etc/foundationdb/fdb.cluster"


@contextmanager
def container(
    docker_client,
    *,
    name,
    network,
    report_dir,
):
    print(f"Pulling FoundationDB image: {FDB_IMAGE}")
    docker_client.images.pull(FDB_IMAGE)

    print(f"Creating container: {name}")
    container = docker_client.containers.create(
        image=FDB_IMAGE,
        name=name,
        environment={
            # "container" mode makes the coordinator advertise the container's
            # own IP in the cluster file, which is reachable by the other
            # containers on the same Docker network.
            "FDB_NETWORKING_MODE": "container",
            "FDB_PORT": "4500",
            "FDB_CLUSTER_FILE": CLUSTER_FILE,
        },
        network=network.name,
        hostname=name,
        detach=True,
    )
    print(f"Starting container: {name}")
    container.start()

    try:
        _wait_for_cluster_file(container)
        _initialize(container)

    except RuntimeError as err:
        docker_utils.teardown_container(container, report_dir)
        pytest.fail(f"{err}", pytrace=False)

    yield container

    docker_utils.teardown_container(container, report_dir)


def read_cluster_file(container):
    print(f"Reading FoundationDB cluster file: {container.name}")
    code, output = container.exec_run(["cat", CLUSTER_FILE])
    if code != 0:
        raise RuntimeError(
            f"Failed to read FoundationDB cluster file: {output!r}"
        )
    return output.decode().strip()


def _wait_for_cluster_file(container, attempts=60):
    print(f"Waiting for FoundationDB cluster file: {container.name}")
    for _ in range(attempts):
        code, _ = container.exec_run(["test", "-s", CLUSTER_FILE])
        if code == 0:
            return
        sleep(1)

    raise RuntimeError("FoundationDB cluster file was not created")


def _initialize(container, attempts=60):
    print(f"Initializing FoundationDB: {container.name}")
    for _ in range(attempts):
        code, output = container.exec_run(
            ["fdbcli", "-C", CLUSTER_FILE, "--exec", "status minimal", "--timeout", "5"]
        )
        if code == 0 and b"The database is available" in output:
            print(f"FoundationDB is available: {container.name}")
            return

        # Not configured yet: create a fresh single-node database.
        container.exec_run(
            ["fdbcli", "-C", CLUSTER_FILE, "--exec", "configure new single ssd", "--timeout", "20"]
        )
        sleep(1)

    raise RuntimeError("FoundationDB did not become available")
