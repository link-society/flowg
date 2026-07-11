import pytest

from pathlib import Path
from shutil import rmtree
import os

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
def flowg_image():
    img = os.getenv("FLOWG_TEST_DOCKER_IMAGE_NAME", "linksociety/flowg:latest")
    print(f"Using Flowg Docker image: {img}")
    return img


@pytest.fixture(scope='module')
def flowg_node0_container(
    report_dir,
    docker_client,
    flowg_network,
    flowg_node0_volume,
    flowg_image,
):
    with flowg_utils.badgerdb_container(
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
def flowg_server(
    flowg_node0_container,
):
    yield


@pytest.fixture(scope='module')
def flowg_admin_token(flowg_server):
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
