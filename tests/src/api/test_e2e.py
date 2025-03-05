import pytest

from datetime import datetime, timedelta, timezone
import subprocess

from ..settings import has_test_suite


@pytest.mark.skipif(not has_test_suite("api"), reason="API test suite excluded")
def test_hurl(
    spec_dir,
    report_dir,
    cache_dir,
    flowg_admin_token,
    flowg_guest_token,
):
    print("Running Hurl test suite:\n")

    now = datetime.now(timezone.utc)
    timewindow_begin = (now - timedelta(minutes=5)).strftime("%Y-%m-%dT%H:%M:%SZ")
    timewindow_end = (now + timedelta(minutes=5)).strftime("%Y-%m-%dT%H:%M:%SZ")

    subprocess.check_call(
        f"""
        hurl \
            --file-root={cache_dir} \
            --variable admin_token={flowg_admin_token} \
            --variable guest_token={flowg_guest_token} \
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
