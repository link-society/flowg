import os
import time
from pathlib import Path
from shutil import rmtree

import pytest
from opentelemetry.proto.collector.logs.v1.logs_service_pb2 import ExportLogsServiceRequest
from opentelemetry.proto.common.v1.common_pb2 import AnyValue, KeyValue
from opentelemetry.proto.resource.v1.resource_pb2 import Resource

import boto3
from .._lib import docker_utils


@pytest.fixture(scope='module')
def floci_aws_image():
    img = os.getenv("FLOCI_AWS_TEST_DOCKER_IMAGE_NAME", "floci/floci:latest")
    print(f"Using floci Docker image: {img}")
    return img

@pytest.fixture(scope='module')
def floci_gcp_image():
    img = os.getenv("FLOCI_GCP_TEST_DOCKER_IMAGE_NAME", "floci/floci-gcp:latest")
    print(f"Using floci Docker image: {img}")
    return img

@pytest.fixture(scope='module')
def floci_aws_container(
        docker_client,
        flowg_network,
        report_dir,
        floci_aws_image
):
    name = "floci-aws-test-server"

    print(f"Creating Container: {name}")
    container = docker_client.containers.run(
        image=floci_aws_image,
        name=name,
        network=flowg_network.name,
        hostname=name,
        ports={
            "4566/tcp": 4566
        },
        detach=True,
    )

    yield

    docker_utils.teardown_container(container, report_dir)

@pytest.fixture(scope='module')
def floci_gcp_container(
        docker_client,
        flowg_network,
        report_dir,
        floci_gcp_image
):
    name = "floci-gcp-test-server"

    print(f"Creating Container: {name}")
    container = docker_client.containers.run(
        image=floci_gcp_image,
        name=name,
        network=flowg_network.name,
        hostname=name,
        ports={
            "4588/tcp": 4588
        },
        detach=True,
    )

    yield

    docker_utils.teardown_container(container, report_dir)


@pytest.fixture(scope="module")
def cloudwatch_log_stream(floci_aws_container):
    print("Creating CloudWatch log stream")
    client = boto3.client(
        "logs",
        endpoint_url='http://localhost:4566',
        region_name='us-east-1',
        aws_access_key_id='test',
        aws_secret_access_key='test'
    )

    client.create_log_group(logGroupName="flowg")
    client.create_log_stream(logGroupName="flowg", logStreamName="logs")


@pytest.fixture(scope="module")
def cache_dir():
    cache_dir = Path.cwd() / "cache" / "api"
    rmtree(cache_dir, ignore_errors=True)
    cache_dir.mkdir(parents=True)
    (cache_dir / "backup").mkdir()
    yield cache_dir
    rmtree(cache_dir, ignore_errors=True)


@pytest.fixture(scope="module")
def otlp_pb(cache_dir):
    # build request
    req = ExportLogsServiceRequest()
    rl = req.resource_logs.add()

    # resource attributes
    res = Resource()
    res.attributes.add(key="service.name", value=AnyValue(string_value="my-service"))
    rl.resource.CopyFrom(res)

    # one scopeLogs + one logRecord
    sl = rl.scope_logs.add()
    lr = sl.log_records.add()
    lr.time_unix_nano = int(time.time() * 1e9)
    lr.severity_number = 9
    lr.severity_text = "Info"
    lr.body.string_value = "Hello from Python"

    ingest_otlp = cache_dir / "ingest_otlp.pb"

    ingest_otlp.write_bytes(req.SerializeToString())

    yield

    ingest_otlp.unlink(missing_ok=True)
