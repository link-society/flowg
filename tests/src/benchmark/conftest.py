import pytest

import os


@pytest.fixture()
def iteration_count():
    return int(os.getenv("BENCHMARK_ITERATIONS", 1_000_000))
