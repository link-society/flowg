import pytest

import multiprocessing

from tqdm import tqdm

from .utils import iteration, timeit
from ..settings import has_test_suite


@pytest.mark.skipif(not has_test_suite("bench"), reason="Benchmark excluded")
def test_perf(testdata, iteration_count, job_count, flowg_config, flowg_admin_token):
    print("Running performance benchmark:\n")

    with timeit(iteration_count):
        with multiprocessing.Pool(job_count) as pool:
            params = [
                (flowg_admin_token, testdata)
                for _ in range(iteration_count)
            ]
            iterable = tqdm(
                pool.imap_unordered(iteration, params),
                total=iteration_count,
                desc="Sending logs",
                unit="req",
            )

            for _ in iterable:
                pass
