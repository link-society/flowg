import pytest

from datetime import datetime, timedelta, timezone
import subprocess

from ..settings import has_test_suite


@pytest.mark.skipif(not has_test_suite("cluster"), reason="Cluster test suite excluded")
def test_hurl(
    spec_dir,
    report_dir,
    flowg_cluster,
    flowg_admin_token,
):
    print("Running Hurl test suite:\n")

    now = datetime.now(timezone.utc)
    timewindow_begin = (now - timedelta(minutes=5)).strftime("%Y-%m-%dT%H:%M:%SZ")
    timewindow_end = (now + timedelta(minutes=5)).strftime("%Y-%m-%dT%H:%M:%SZ")

    subprocess.check_call(
        f"""
        hurl \
            --file-root={spec_dir} \
            --variable admin_token={flowg_admin_token} \
            --variable node0_url=http://localhost:5080 \
            --variable node1_url=http://localhost:5081 \
            --variable node2_url=http://localhost:5082 \
            --variable timewindow_begin={timewindow_begin} \
            --variable timewindow_end={timewindow_end} \
            --error-format long \
            --report-html {report_dir / 'html'} \
            --report-junit {report_dir / 'junit.xml'} \
            --jobs 1 \
            --test \
            {spec_dir}
        """,
        shell=True,
    )
