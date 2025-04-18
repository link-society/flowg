import pytest

from pathlib import Path
from shutil import rmtree
import time

from opentelemetry.proto.collector.logs.v1.logs_service_pb2 import ExportLogsServiceRequest
from opentelemetry.proto.resource.v1.resource_pb2 import Resource
from opentelemetry.proto.common.v1.common_pb2 import AnyValue, KeyValue


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
    res.attributes.add(key="service.name", value=AnyValue(string_value="my‑service"))
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
