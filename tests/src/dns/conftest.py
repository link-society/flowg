import pytest

import os

from pathlib import Path
from shutil import rmtree

from .._lib import flowg_utils, docker_utils


@pytest.fixture(scope="module")
def cache_dir():
    cache_dir = Path.cwd() / "cache" / "dns"
    rmtree(cache_dir, ignore_errors=True)
    cache_dir.mkdir(parents=True)
    (cache_dir / "backup").mkdir()
    yield cache_dir
    rmtree(cache_dir, ignore_errors=True)


@pytest.fixture(scope='module')
def dns_server_image():
    img = os.getenv("DNSSERVER_TEST_DOCKER_IMAGE_NAME", "linksociety/dns-server:localdev")
    print(f"Using DNS Server Docker image: {img}")
    return img


@pytest.fixture(scope='module')
def dns_server_container(
    docker_client,
    flowg_network,
    report_dir,
    dns_server_image
):
    name="test-flowg-dns-server"

    print(f"Creating Container: {name}")
    container = docker_client.containers.run(
        image=dns_server_image,
        name=name,
        network=flowg_network.name,
        hostname=name,
        ports={
            "53/udp": 5333,
            "8080/tcp": 8080,
        },
        detach=True,
    )

    yield

    docker_utils.teardown_container(container, report_dir)


@pytest.fixture(scope='module')
def dns_client_image():
    img = os.getenv("DNSCLIENT_TEST_DOCKER_IMAGE_NAME", "linksociety/dns-client:localdev")
    print(f"Using DNS Client Docker image: {img}")
    return img


@pytest.fixture(scope='module')
def flowg_node0_container(
    report_dir,
    docker_client,
    flowg_network,
    flowg_node0_volume,
    dns_client_image,
):
    with flowg_utils.container(
        docker_client,
        name="test-flowg-node0",
        network=flowg_network,
        volume=flowg_node0_volume,
        image=dns_client_image,
        environment={
            "FLOWG_CLUSTER_FORMATION_DNS_SERVER": "test-flowg-dns-server:53",
            "FLOWG_CLUSTER_FORMATION_DNS_MANAGEMENT_SERVER": "test-flowg-dns-server:8080",
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
    dns_client_image,
):
    with flowg_utils.container(
        docker_client,
        name="test-flowg-node1",
        network=flowg_network,
        volume=flowg_node1_volume,
        image=dns_client_image,
        environment={
            "FLOWG_HTTP_BIND_ADDRESS": ":5081",
            "FLOWG_MGMT_BIND_ADDRESS": ":9114",
            "FLOWG_SYSLOG_BIND_ADDRESS": ":5515",
            "FLOWG_CLUSTER_FORMATION_DNS_SERVER": "test-flowg-dns-server:53",
            "FLOWG_CLUSTER_FORMATION_DNS_MANAGEMENT_SERVER": "test-flowg-dns-server:8080",
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
    dns_client_image,
):
    with flowg_utils.container(
        docker_client,
        name="test-flowg-node2",
        network=flowg_network,
        volume=flowg_node2_volume,
        image=dns_client_image,
        environment={
            "FLOWG_HTTP_BIND_ADDRESS": ":5082",
            "FLOWG_MGMT_BIND_ADDRESS": ":9115",
            "FLOWG_SYSLOG_BIND_ADDRESS": ":5516",
            "FLOWG_CLUSTER_FORMATION_DNS_SERVER": "test-flowg-dns-server:53",
            "FLOWG_CLUSTER_FORMATION_DNS_MANAGEMENT_SERVER": "test-flowg-dns-server:8080",
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
    dns_server_container,
    flowg_node0_container,
    flowg_node1_container,
    flowg_node2_container,
):
    yield
