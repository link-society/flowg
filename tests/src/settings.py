import os


FLOWG_TEST_SUITES = [
    item.strip()
    for item in os.getenv("FLOWG_TEST_SUITES", "*").split(",")
]


def has_test_suite(name):
    return name in FLOWG_TEST_SUITES or "*" in FLOWG_TEST_SUITES
