import pytest

from pathlib import Path
from shutil import rmtree

@pytest.fixture(scope="module")
def cache_dir():
    cache_dir = Path.cwd() / "cache" / "api"
    rmtree(cache_dir, ignore_errors=True)
    cache_dir.mkdir(parents=True)
    (cache_dir / "backup").mkdir()
    yield cache_dir
    rmtree(cache_dir, ignore_errors=True)


@pytest.fixture(scope='module')
def consul_container(
    report_dir,
    docker_client,
    flowg_network
):
    
    consul_image = "hashicorp/consul"

    print("Pulling the Consul image")
    try:
        docker_client.images.pull(consul_image)
    except Exception as e:
            print(f"An unexpected error occurred while pulling the image: {e}")

    print("Creating Container: consul-container")
    container = docker_client.containers.run(
        image=consul_image,
        name="consul-container",
        network=flowg_network.name,
        hostname="consul-container",
        ports={
            "8500/tcp": 8500,
        },
        detach=True,
    )

    yield

    print("Stopping container: consul-container")
    container.stop()

    print("Writing logs: consul-container")
    with open(report_dir / "consul-container.log", "wb") as f:
        for data in container.logs(stream=True):
            f.write(data)

    print("Removing container: consul-contaier")
    container.remove(force=True)

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
            "CONSUL_URL": "http://consul-container:8500",
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
            "CONSUL_URL": "http://consul-container:8500",
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
            "CONSUL_URL": "http://consul-container:8500",
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
def flowg_cluster(consul_container, flowg_node0_container, flowg_node1_container, flowg_node2_container):
    yield

def wait_for_healthcheck(container):
    while True:
        container.reload()

        if container.health == "healthy":
            break

        elif container.health == "unhealthy":
            raise RuntimeError(f"Node {container.name} was not healthy")

        sleep(1)
