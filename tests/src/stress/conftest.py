import pytest

import json

import requests

from .._lib import flowg_utils


@pytest.fixture()
def flowg_stress_container(limits, docker_client, flowg_network, flowg_image, report_dir):
    name = "test-flowg-stress"

    with flowg_utils.volume(docker_client, name=name) as volume:
        with flowg_utils.badgerdb_container(
            docker_client,
            name=name,
            network=flowg_network,
            volume=volume,
            image=flowg_image,
            environment={
                "FLOWG_BADGER_AUTH_MEMTABLE_SIZE": str(limits.memtable_size),
                "FLOWG_BADGER_LOG_MEMTABLE_SIZE": str(limits.memtable_size),
                "FLOWG_BADGER_CONFIG_MEMTABLE_SIZE": str(limits.memtable_size),
            },
            ports={
                "5080/tcp": 5080,
                "5514/udp": 5514,
                "9113/tcp": 9113,
            },
            report_dir=report_dir,
            nano_cpus=limits.nano_cpus,
            mem_limit=limits.memory,
        ):
            yield docker_client.containers.get(name)


@pytest.fixture()
def flowg_stress_token(flowg_stress_container):
    return flowg_utils.create_token(username="root", password="root")


@pytest.fixture()
def flowg_stress_config(flowg_stress_token, spec_dir):
    headers = {"Authorization": f"Bearer {flowg_stress_token}"}

    with open(spec_dir / "config" / "transformers" / "json.vrl") as f:
        resp = requests.put(
            "http://localhost:5080/api/v1/transformers/json",
            headers=headers,
            json={"script": f.read()},
        )
        resp.raise_for_status()

    with open(spec_dir / "config" / "pipelines" / "transform.json") as f:
        resp = requests.put(
            "http://localhost:5080/api/v1/pipelines/transform",
            headers=headers,
            json={"flow": json.load(f)},
        )
        resp.raise_for_status()
