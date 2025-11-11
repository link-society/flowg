import pytest

import subprocess
import json

from ..settings import has_test_suite


@pytest.mark.skipif(not has_test_suite("bench"), reason="Benchmark excluded")
def test_perf(iteration_count, flowg_admin_token):
    print("Running performance benchmark:\n")

    payload = json.dumps({"records": [
        {"message": "hello world"}
    ]})

    subprocess.check_call(
        f"""
        oha \
            -n {iteration_count} \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer {flowg_admin_token}" \
            "http://localhost:5080/api/v1/pipelines/default/logs/struct" \
            -m "POST" \
            -d '{payload}'
        """,
        shell=True,
    )
