import pytest

from pathlib import Path
from shutil import rmtree
import docker

from ._lib import docker_utils, flowg_utils, mockserver_utils


@pytest.fixture(scope='module', autouse=True)
def log_separator(request):
    print("-" * 80)
    print("--- Running test:", request.node.path)
    yield
    print("-" * 80)


@pytest.fixture(scope="module")
def config_dir():
    config_dir = Path.cwd() / "config"
    yield config_dir


@pytest.fixture(scope="module")
def spec_dir(request):
    spec_dir = Path.cwd() / "specs" / request.node.parent.name
    yield spec_dir


@pytest.fixture(scope="module")
def report_dir(request):
    report_dir = Path.cwd() / "reports" / request.node.parent.name
    rmtree(report_dir, ignore_errors=True)
    report_dir.mkdir(parents=True)
    yield report_dir


@pytest.fixture(scope='module')
def docker_client():
    client = docker.from_env()
    docker_utils.cleanup(client)
    yield client


@pytest.fixture(scope='module')
def flowg_network(docker_client):
    with flowg_utils.network(docker_client, name="test-flowg") as network:
        yield network


@pytest.fixture(scope='module')
def flowg_node0_volume(docker_client):
    with flowg_utils.volume(docker_client, name="test-flowg-node0") as volume:
        yield volume


@pytest.fixture(scope='module')
def flowg_node1_volume(docker_client):
    with flowg_utils.volume(docker_client, name="test-flowg-node1") as volume:
        yield volume


@pytest.fixture(scope='module')
def flowg_node2_volume(docker_client):
    with flowg_utils.volume(docker_client, name="test-flowg-node2") as volume:
        yield volume


@pytest.fixture(scope='module')
def flowg_image():
    return "linksociety/flowg:latest"


@pytest.fixture(scope='module')
def flowg_node0_container(
    report_dir,
    docker_client,
    flowg_network,
    flowg_node0_volume,
    flowg_image,
):
    with flowg_utils.container(
        docker_client,
        name="test-flowg-node0",
        network=flowg_network,
        volume=flowg_node0_volume,
        image=flowg_image,
        environment={},
        ports={
            "5080/tcp": 5080,
            "9113/tcp": 9113,
            "5514/udp": 5514,
        },
        report_dir=report_dir,
    ):
        yield


@pytest.fixture(scope='module')
def flowg_node1_container(
    report_dir,
    docker_client,
    flowg_network,
    flowg_node1_volume,
    flowg_image,
):
    with flowg_utils.container(
        docker_client,
        name="test-flowg-node1",
        network=flowg_network,
        volume=flowg_node1_volume,
        image=flowg_image,
        environment={
            "FLOWG_CLUSTER_JOIN_NODE_ID": "test-flowg-node0",
            "FLOWG_CLUSTER_JOIN_ENDPOINT": "http://test-flowg-node0:9113",
            "FLOWG_HTTP_BIND_ADDRESS": ":5081",
            "FLOWG_MGMT_BIND_ADDRESS": ":9114",
            "FLOWG_SYSLOG_BIND_ADDRESS": ":5515",
        },
        ports={
            "5081/tcp": 5081,
            "9114/tcp": 9114,
            "5515/udp": 5515,
        },
        report_dir=report_dir,
    ):
        yield


@pytest.fixture(scope='module')
def flowg_node2_container(
    report_dir,
    docker_client,
    flowg_network,
    flowg_node2_volume,
    flowg_image,
):
    with flowg_utils.container(
        docker_client,
        name="test-flowg-node2",
        network=flowg_network,
        volume=flowg_node2_volume,
        image=flowg_image,
        environment={
            "FLOWG_CLUSTER_JOIN_NODE_ID": "test-flowg-node1",
            "FLOWG_CLUSTER_JOIN_ENDPOINT": "http://test-flowg-node1:9114",
            "FLOWG_HTTP_BIND_ADDRESS": ":5082",
            "FLOWG_MGMT_BIND_ADDRESS": ":9115",
            "FLOWG_SYSLOG_BIND_ADDRESS": ":5516",
        },
        ports={
            "5082/tcp": 5082,
            "9115/tcp": 9115,
            "5516/udp": 5516,
        },
        report_dir=report_dir,
    ):
        yield


@pytest.fixture(scope='module')
def flowg_cluster(
    flowg_node0_container,
    flowg_node1_container,
    flowg_node2_container,
):
    yield


@pytest.fixture(scope='module')
def flowg_admin_token(flowg_cluster):
    return flowg_utils.create_token(username="root", password="root")


@pytest.fixture(scope='module')
def flowg_guest_token(flowg_admin_token):
    flowg_utils.create_user(
        token=flowg_admin_token,
        username="guest",
        password="guest",
    )

    return flowg_utils.create_token(
        username="guest",
        password="guest",
    )


@pytest.fixture(scope='module')
def mockserver_container(
    config_dir,
    report_dir,
    docker_client,
    flowg_network,
):
    with mockserver_utils.container(
        docker_client,
        name="test-flowg-mockserver",
        network=flowg_network,
        config_dir=config_dir,
        report_dir=report_dir,
    ):
        yield


def pytest_report_teststatus(report, config):
    return report.outcome, "", report.outcome.upper()
