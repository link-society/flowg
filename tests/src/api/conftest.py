import os
import time
from pathlib import Path
from shutil import rmtree

import pytest
from opentelemetry.proto.collector.logs.v1.logs_service_pb2 import ExportLogsServiceRequest
from opentelemetry.proto.common.v1.common_pb2 import AnyValue, KeyValue
from opentelemetry.proto.resource.v1.resource_pb2 import Resource

from azure.core.pipeline.transport import RequestsTransport
from azure.core.credentials import AccessToken, TokenCredential
from azure.mgmt.monitor import MonitorManagementClient
from azure.mgmt.loganalytics import LogAnalyticsManagementClient

import urllib3

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
    print(f"Using floci-gcp Docker image: {img}")
    return img

@pytest.fixture(scope='module')
def floci_az_image():
    img = os.getenv("FLOCI_AZ_TEST_DOCKER_IMAGE_NAME", "floci/floci-az:latest")
    print(f"Using floci-az Docker image: {img}")
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


@pytest.fixture(scope='module')
def floci_az_container(
        docker_client,
        flowg_network,
        report_dir,
        floci_az_image
):
    name = "floci-az-test-server"

    print(f"Creating Container: {name}")
    container = docker_client.containers.run(
        image=floci_az_image,
        name=name,
        network=flowg_network.name,
        hostname=name,
        ports={
            "4577/tcp": 4577
        },
        environment={
            "FLOCI_AZ_TLS_ENABLED": "true"
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
def azuremonitor_setup_dcr(floci_az_container):
    base_url = "https://localhost:4577"
    subscription_id = "subscription_id"
    resource_group = "resource_group"
    dcr_name = "flowg"
    location = "westus"
    workspace_name = "workspace"

    print("Setting up Microsoft Azure Monitor environment")

    class StaticTokenCredential(TokenCredential):
        def __init__(self, token: str):
            self._token = token

        def get_token(self, *scopes, **kwargs):
            return AccessToken(self._token, 4102444800)  # 2100-01-01T00:00:00Z

    credential = StaticTokenCredential("TOKEN")
    transport = RequestsTransport(connection_verify=False)

    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

    client = MonitorManagementClient(
        credential,
        subscription_id,
        base_url=base_url,
        transport=transport
    )

    la = LogAnalyticsManagementClient(
        credential,
        subscription_id,
        base_url=base_url,
        transport=transport
    )

    workspace = la.workspaces.begin_create_or_update(
        resource_group_name=resource_group,
        workspace_name=workspace_name,
        parameters={
            "location": location,
        },
    ).result()

    dcr_definition = {
        "location": location,
        "properties": {
            "destinations": {
                "logAnalytics": [
                    {
                        "name": dcr_name,
                        "workspaceResourceId": workspace.id
                    }
                ]
            },
        }
    }

    dcr = client.data_collection_rules.create(
        resource_group_name=resource_group,
        data_collection_rule_name=dcr_name,
        body=dcr_definition,
    )

    return dcr.immutable_id

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
