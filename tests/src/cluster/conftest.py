import pytest

from contextlib import ExitStack
from pathlib import Path
from shutil import rmtree

from .._lib import fdb_utils, flowg_utils


# Each node maps the FlowG HTTP port (5080) to a distinct host port so that the
# Hurl test suite can reach every node from the host.
CLUSTER_NODES = [
    ("test-flowg-node0", {"5080/tcp": 5080}),
    ("test-flowg-node1", {"5080/tcp": 5081}),
    ("test-flowg-node2", {"5080/tcp": 5082}),
]


@pytest.fixture(scope="module")
def cluster_conf_dir():
    conf_dir = Path.cwd() / "cache" / "cluster"
    rmtree(conf_dir, ignore_errors=True)
    conf_dir.mkdir(parents=True)
    yield conf_dir
    rmtree(conf_dir, ignore_errors=True)


@pytest.fixture(scope="module")
def fdb_container(docker_client, flowg_network, report_dir):
    with fdb_utils.container(
        docker_client,
        name="test-flowg-fdb",
        network=flowg_network,
        report_dir=report_dir,
    ) as container:
        yield container


@pytest.fixture(scope="module")
def fdb_cluster_file(fdb_container, cluster_conf_dir):
    content = fdb_utils.read_cluster_file(fdb_container)

    cluster_file = cluster_conf_dir / "fdb.cluster"
    cluster_file.write_text(content + "\n")
    cluster_file.chmod(0o644)

    yield cluster_file


@pytest.fixture(scope="module")
def flowg_cluster(
    docker_client,
    flowg_network,
    flowg_image,
    report_dir,
    cluster_conf_dir,
    fdb_cluster_file,
):
    with ExitStack() as stack:
        for name, ports in CLUSTER_NODES:
            stack.enter_context(
                flowg_utils.foundationdb_container(
                    docker_client,
                    name=name,
                    network=flowg_network,
                    conf_dir=cluster_conf_dir,
                    image=flowg_image,
                    ports=ports,
                    report_dir=report_dir,
                )
            )

        yield


@pytest.fixture(scope="module")
def flowg_admin_token(flowg_cluster):
    return flowg_utils.create_token(username="root", password="root")
