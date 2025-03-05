import pytest

from pathlib import Path
from shutil import rmtree

from tenacity import retry, stop_after_attempt, wait_fixed
import requests
import docker


@pytest.fixture(scope='module', autouse=True)
def log_separator(request):
    print("-" * 80)
    print("--- Running test:", request.node.path)
    yield
    print("-" * 80)


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

    print("Cleaning up old containers")
    for container in client.containers.list(all=True):
        if container.name.startswith("test-flowg-node"):
            container.remove(force=True)

    print("Cleaning up old volumes")
    for volume in client.volumes.list():
        if volume.name.startswith("test-flowg-node"):
            volume.remove(force=True)

    print("Cleaning up old networks")
    for network in client.networks.list():
        if network.name.startswith("test-flowg"):
            network.remove()

    yield client


@pytest.fixture(scope='module')
def flowg_network(docker_client):
    print("Creating network: test-flowg")
    network = docker_client.networks.create(name="test-flowg", driver="bridge")
    yield network
    print("Removing network: test-flowg")
    network.remove()


@pytest.fixture(scope='module')
def flowg_node0_volume(docker_client):
    print("Creating volume: test-flowg-node0")
    volume = docker_client.volumes.create(name="test-flowg-node0")
    yield volume
    print("Removing volume: test-flowg-node0")
    volume.remove()


@pytest.fixture(scope='module')
def flowg_node1_volume(docker_client):
    print("Creating volume: test-flowg-node1")
    volume = docker_client.volumes.create(name="test-flowg-node1")
    yield volume
    print("Removing volume: test-flowg-node1")
    volume.remove()


@pytest.fixture(scope='module')
def flowg_node2_volume(docker_client):
    print("Creating volume: test-flowg-node2")
    volume = docker_client.volumes.create(name="test-flowg-node2")
    yield volume
    print("Removing volume: test-flowg-node2")
    volume.remove()


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
    print("Creating container: test-flowg-node0")
    container = docker_client.containers.run(
        image=flowg_image,
        name="test-flowg-node0",
        environment={
            "FLOWG_SECRET_KEY": "s3cr3!",
            "FLOWG_CLUSTER_NODE_ID": "test-flowg-node0",
            "FLOWG_AUTH_DIR": "/data/auth",
            "FLOWG_CONFIG_DIR": "/data/config",
            "FLOWG_LOG_DIR": "/data/logs",
        },
        network=flowg_network.name,
        hostname="test-flowg-node0",
        ports={
            "5080/tcp": 5080,
            "9113/tcp": 9113,
            "5514/udp": 5514,
        },
        volumes={
            flowg_node0_volume.name: {"bind": "/data", "mode": "rw"}
        },
        detach=True,
    )
    print("Waiting for healthcheck: test-flowg-node0")
    wait_for_healthcheck("test-flowg-node0", "http://localhost:9113/health")

    yield

    print("Stopping container: test-flowg-node0")
    container.stop()

    print("Writing logs: test-flowg-node0")
    with open(report_dir / "docker-node0.log", "wb") as f:
        for data in container.logs(stream=True):
            f.write(data)

    print("Removing container: test-flowg-node0")
    container.remove(force=True)


@pytest.fixture(scope='module')
def flowg_node1_container(
    report_dir,
    docker_client,
    flowg_network,
    flowg_node1_volume,
    flowg_image,
):
    print("Creating container: test-flowg-node1")
    container = docker_client.containers.run(
        image=flowg_image,
        name="test-flowg-node1",
        environment={
            "FLOWG_SECRET_KEY": "s3cr3!",
            "FLOWG_CLUSTER_NODE_ID": "test-flowg-node1",
            "FLOWG_CLUSTER_JOIN_NODE_ID": "test-flowg-node0",
            "FLOWG_CLUSTER_JOIN_ENDPOINT": "http://test-flowg-node0:9113",
            "FLOWG_HTTP_BIND_ADDRESS": ":5081",
            "FLOWG_MGMT_BIND_ADDRESS": ":9114",
            "FLOWG_SYSLOG_BIND_ADDRESS": ":5515",
            "FLOWG_AUTH_DIR": "/data/auth",
            "FLOWG_CONFIG_DIR": "/data/config",
            "FLOWG_LOG_DIR": "/data/logs",
        },
        network=flowg_network.name,
        hostname="test-flowg-node1",
        ports={
            "5081/tcp": 5081,
            "9114/tcp": 9114,
            "5515/udp": 5515,
        },
        volumes={
            flowg_node1_volume.name: {"bind": "/data", "mode": "rw"}
        },
        detach=True,
    )
    print("Waiting for healthcheck: test-flowg-node1")
    wait_for_healthcheck("test-flowg-node1", "http://localhost:9114/health")

    yield

    print("Stopping container: test-flowg-node1")
    container.stop()

    print("Writing logs: test-flowg-node1")
    with open(report_dir / "docker-node1.log", "wb") as f:
        for data in container.logs(stream=True):
            f.write(data)

    print("Removing container: test-flowg-node1")
    container.remove(force=True)


@pytest.fixture(scope='module')
def flowg_node2_container(
    report_dir,
    docker_client,
    flowg_network,
    flowg_node2_volume,
    flowg_image,
):
    print("Creating container: test-flowg-node2")
    container = docker_client.containers.run(
        image=flowg_image,
        name="test-flowg-node2",
        environment={
            "FLOWG_SECRET_KEY": "s3cr3!",
            "FLOWG_CLUSTER_NODE_ID": "test-flowg-node2",
            "FLOWG_CLUSTER_JOIN_NODE_ID": "test-flowg-node1",
            "FLOWG_CLUSTER_JOIN_ENDPOINT": "http://test-flowg-node1:9114",
            "FLOWG_HTTP_BIND_ADDRESS": ":5082",
            "FLOWG_MGMT_BIND_ADDRESS": ":9115",
            "FLOWG_SYSLOG_BIND_ADDRESS": ":5516",
            "FLOWG_AUTH_DIR": "/data/auth",
            "FLOWG_CONFIG_DIR": "/data/config",
            "FLOWG_LOG_DIR": "/data/logs",
        },
        network=flowg_network.name,
        hostname="test-flowg-node2",
        ports={
            "5082/tcp": 5082,
            "9115/tcp": 9115,
            "5516/udp": 5516,
        },
        volumes={
            flowg_node2_volume.name: {"bind": "/data", "mode": "rw"}
        },
        detach=True,
    )
    print("Waiting for healthcheck: test-flowg-node2")
    wait_for_healthcheck("test-flowg-node2", "http://localhost:9115/health")

    yield

    print("Stopping container: test-flowg-node2")
    container.stop()

    print("Writing logs: test-flowg-node2")
    with open(report_dir / "docker-node2.log", "wb") as f:
        for data in container.logs(stream=True):
            f.write(data)

    print("Removing container: test-flowg-node2")
    container.remove(force=True)


@pytest.fixture(scope='module')
def flowg_cluster(flowg_node0_container, flowg_node1_container, flowg_node2_container):
    yield


@pytest.fixture(scope='module')
def flowg_admin_token(flowg_cluster):
    print("Creating admin token")
    resp = requests.post(
        "http://localhost:5080/api/v1/auth/login",
        json={"username": "root", "password": "root"},
    )
    resp.raise_for_status()
    data = resp.json()
    admin_jwt = data["token"]

    resp = requests.post(
        "http://localhost:5080/api/v1/token",
        headers={"Authorization": f"Bearer jwt:{admin_jwt}"},
    )
    resp.raise_for_status()
    data = resp.json()
    return data["token"]


@pytest.fixture(scope='module')
def flowg_guest_token(flowg_admin_token):
    print("Creating guest token")
    resp = requests.put(
        "http://localhost:5080/api/v1/users/guest",
        headers={"Authorization": f"Bearer pat:{flowg_admin_token}"},
        json={"password": "guest", "roles": []},
    )
    resp.raise_for_status()

    resp = requests.post(
        "http://localhost:5080/api/v1/auth/login",
        json={"username": "guest", "password": "guest"},
    )
    resp.raise_for_status()
    data = resp.json()
    guest_jwt = data["token"]

    resp = requests.post(
        "http://localhost:5080/api/v1/token",
        headers={"Authorization": f"Bearer jwt:{guest_jwt}"},
    )
    resp.raise_for_status()
    data = resp.json()
    return data["token"]


def wait_for_healthcheck(nodename, endpoint):
    @retry(
        stop=stop_after_attempt(5),
        wait=wait_fixed(1),
    )
    def impl():
        resp = requests.get(endpoint)
        resp.raise_for_status()

    try:
        impl()

    except Exception:
        pytest.fail(f"Node {nodename} was not healthy", pytrace=False)


def pytest_report_teststatus(report, config):
    return report.outcome, "", report.outcome.upper()
