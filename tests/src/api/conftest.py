import pytest

from pathlib import Path
from shutil import rmtree


@pytest.fixture(scope="module")
def cache_dir():
    cache_dir = Path.cwd() / "cache" / "api"
    rmtree(cache_dir, ignore_errors=True)
    cache_dir.mkdir(parents=True)
    (cache_dir / "backup").mkdir()
    yield cache_dir
    rmtree(cache_dir, ignore_errors=True)
