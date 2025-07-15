import pytest

from pathlib import Path
from shutil import rmtree

from .._lib import flowg_utils, consul_utils


@pytest.fixture(scope="module")
def cache_dir():
    cache_dir = Path.cwd() / "cache" / "consul"
    rmtree(cache_dir, ignore_errors=True)
    cache_dir.mkdir(parents=True)
    (cache_dir / "backup").mkdir()
    yield cache_dir
    rmtree(cache_dir, ignore_errors=True)


@pytest.fixture(scope='module')
def consul_container(
    docker_client,
    flowg_network,
    report_dir,
):
    with consul_utils.container(
        docker_client,
        name="test-flowg-consul",
        network=flowg_network,
        report_dir=report_dir,
    ):
        yield


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
        environment={
            "FLOWG_CLUSTER_FORMATION_STRATEGY": "consul",
            "FLOWG_CLUSTER_FORMATION_CONSUL_URL": "http://test-flowg-consul:8500",
        },
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
            "FLOWG_HTTP_BIND_ADDRESS": ":5081",
            "FLOWG_MGMT_BIND_ADDRESS": ":9114",
            "FLOWG_SYSLOG_BIND_ADDRESS": ":5515",
            "FLOWG_CLUSTER_FORMATION_STRATEGY": "consul",
            "FLOWG_CLUSTER_FORMATION_CONSUL_URL": "http://test-flowg-consul:8500",
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
            "FLOWG_HTTP_BIND_ADDRESS": ":5082",
            "FLOWG_MGMT_BIND_ADDRESS": ":9115",
            "FLOWG_SYSLOG_BIND_ADDRESS": ":5516",
            "FLOWG_CLUSTER_FORMATION_STRATEGY": "consul",
            "FLOWG_CLUSTER_FORMATION_CONSUL_URL": "http://test-flowg-consul:8500",
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
    consul_container,
    flowg_node0_container,
    flowg_node1_container,
    flowg_node2_container,
):
    yield
