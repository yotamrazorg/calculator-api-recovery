"""Functional tests for the Go calculator API — verifies behavior parity with Python source."""
import re
import requests

BASE_URL = "http://localhost:8000"


# --- Health Check ---

def test_health():
    resp = requests.get(f"{BASE_URL}/health")
    assert resp.status_code == 200
    data = resp.json()
    assert data["status"] == "ok"
    assert data["version"] == "0.1.0"


# --- Arithmetic Endpoints ---

def test_add():
    resp = requests.post(f"{BASE_URL}/add", json={"a": 5, "b": 3})
    assert resp.status_code == 200
    assert resp.json() == {"result": 8.0}


def test_subtract():
    resp = requests.post(f"{BASE_URL}/subtract", json={"a": 10, "b": 4})
    assert resp.status_code == 200
    assert resp.json() == {"result": 6.0}


def test_multiply():
    resp = requests.post(f"{BASE_URL}/multiply", json={"a": 3, "b": 7})
    assert resp.status_code == 200
    assert resp.json() == {"result": 21.0}


def test_divide():
    resp = requests.post(f"{BASE_URL}/divide", json={"a": 15, "b": 3})
    assert resp.status_code == 200
    assert resp.json() == {"result": 5.0}


def test_divide_by_zero():
    resp = requests.post(f"{BASE_URL}/divide", json={"a": 5, "b": 0})
    assert resp.status_code == 400
    data = resp.json()
    assert data["detail"] == "Cannot divide by zero"


def test_add_with_floats():
    resp = requests.post(f"{BASE_URL}/add", json={"a": 1.5, "b": 2.5})
    assert resp.status_code == 200
    assert resp.json() == {"result": 4.0}


# --- POST /calculations ---

def test_create_calculation_add():
    resp = requests.post(f"{BASE_URL}/calculations", json={"operation": "add", "a": 5, "b": 3})
    assert resp.status_code == 201
    data = resp.json()
    assert data["operation"] == "add"
    assert data["a"] == 5.0
    assert data["b"] == 3.0
    assert data["result"] == 8.0
    assert "id" in data
    assert isinstance(data["id"], int)
    assert data["id"] >= 1
    # Verify created_at format matches Python (no trailing Z)
    assert "created_at" in data
    assert not data["created_at"].endswith("Z"), f"created_at should not end with Z: {data['created_at']}"
    # Match ISO 8601 without timezone
    assert re.match(r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?$", data["created_at"]), \
        f"created_at format mismatch: {data['created_at']}"


def test_create_calculation_all_operations():
    tests = [
        ("add", 10, 5, 15.0),
        ("sub", 10, 5, 5.0),
        ("mul", 10, 5, 50.0),
        ("div", 10, 5, 2.0),
    ]
    for op, a, b, expected in tests:
        resp = requests.post(f"{BASE_URL}/calculations", json={"operation": op, "a": a, "b": b})
        assert resp.status_code == 201, f"op={op}: status={resp.status_code}"
        assert resp.json()["result"] == expected, f"op={op}: result={resp.json()['result']}"


def test_create_calculation_unknown_operation():
    resp = requests.post(f"{BASE_URL}/calculations", json={"operation": "foo", "a": 5, "b": 3})
    assert resp.status_code == 400
    data = resp.json()
    assert data["detail"] == "Unknown operation: foo. Use: ['add', 'sub', 'mul', 'div']"


def test_create_calculation_divide_by_zero():
    resp = requests.post(f"{BASE_URL}/calculations", json={"operation": "div", "a": 1, "b": 0})
    assert resp.status_code == 400
    data = resp.json()
    assert data["detail"] == "Cannot divide by zero"


# --- GET /calculations ---

def test_list_calculations():
    # Create a fresh pair to ensure the list is non-empty
    requests.post(f"{BASE_URL}/calculations", json={"operation": "add", "a": 100, "b": 200})
    requests.post(f"{BASE_URL}/calculations", json={"operation": "sub", "a": 50, "b": 25})

    resp = requests.get(f"{BASE_URL}/calculations")
    assert resp.status_code == 200
    data = resp.json()
    assert isinstance(data, list)
    assert len(data) >= 2

    # Verify ordered by created_at descending
    for i in range(len(data) - 1):
        assert data[i]["created_at"] >= data[i + 1]["created_at"], \
            f"Not ordered desc: {data[i]['created_at']} < {data[i+1]['created_at']}"

    # Verify each item has all expected fields
    for item in data:
        assert "id" in item
        assert "operation" in item
        assert "a" in item
        assert "b" in item
        assert "result" in item
        assert "created_at" in item


# --- GET /calculations/{id} ---

def test_get_calculation_by_id():
    # Create a calculation first
    create_resp = requests.post(f"{BASE_URL}/calculations", json={"operation": "mul", "a": 4, "b": 5})
    assert create_resp.status_code == 201
    calc_id = create_resp.json()["id"]

    resp = requests.get(f"{BASE_URL}/calculations/{calc_id}")
    assert resp.status_code == 200
    data = resp.json()
    assert data["id"] == calc_id
    assert data["operation"] == "mul"
    assert data["a"] == 4.0
    assert data["b"] == 5.0
    assert data["result"] == 20.0


def test_get_calculation_not_found():
    resp = requests.get(f"{BASE_URL}/calculations/99999")
    assert resp.status_code == 404
    data = resp.json()
    assert data["detail"] == "Calculation not found"


# --- DELETE /calculations/{id} ---

def test_delete_calculation():
    # Create a calculation first
    create_resp = requests.post(f"{BASE_URL}/calculations", json={"operation": "add", "a": 1, "b": 1})
    assert create_resp.status_code == 201
    calc_id = create_resp.json()["id"]

    # Delete it
    resp = requests.delete(f"{BASE_URL}/calculations/{calc_id}")
    assert resp.status_code == 204
    assert resp.text == "" or resp.content == b""

    # Verify it's gone
    resp = requests.get(f"{BASE_URL}/calculations/{calc_id}")
    assert resp.status_code == 404


def test_delete_calculation_not_found():
    resp = requests.delete(f"{BASE_URL}/calculations/99999")
    assert resp.status_code == 404
    data = resp.json()
    assert data["detail"] == "Calculation not found"
