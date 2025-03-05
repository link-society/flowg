import pytest
import robot

from ..settings import has_test_suite


@pytest.mark.skipif(not has_test_suite("web"), reason="Web test suite excluded")
def test_robot(flowg_cluster, spec_dir, report_dir):
    print("Running Robot Framework:\n")

    rc = robot.run(spec_dir, outputdir=report_dir)
    assert rc == 0
