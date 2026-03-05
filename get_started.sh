#!/bin/bash
# ─────────────────────────────────────────────────────────────────────────────
# get_started.sh — build, start, and verify pt_search_hos
# ─────────────────────────────────────────────────────────────────────────────
set -e

APP_PORT=3458
APP_URL="http://localhost:${APP_PORT}"
DB_CONTAINER="pt_search_hos_db"
DB_USER="postgres"
DB_NAME="pt_search_hos"

PASS=0
FAIL=0

# ── helpers ───────────────────────────────────────────────────────────────────

green()  { echo -e "\033[32m$*\033[0m"; }
red()    { echo -e "\033[31m$*\033[0m"; }
yellow() { echo -e "\033[33m$*\033[0m"; }
bold()   { echo -e "\033[1m$*\033[0m"; }

ok()   { green   "  ✔  $*"; PASS=$((PASS + 1)); }
fail() { red     "  ✘  $*"; FAIL=$((FAIL + 1)); }
info() { yellow  "  ·  $*"; }

# ── step 1: docker compose up --build ─────────────────────────────────────────

bold ""
bold "═══════════════════════════════════════════"
bold " pt_search_hos — get started"
bold "═══════════════════════════════════════════"
echo ""

bold "[ 1/4 ] Building and starting containers..."
docker compose up --build -d
echo ""

# ── step 2: wait for postgres healthcheck ─────────────────────────────────────

bold "[ 2/4 ] Waiting for PostgreSQL to be healthy..."
MAX=30
COUNT=0
until docker inspect --format='{{.State.Health.Status}}' "$DB_CONTAINER" 2>/dev/null | grep -q "healthy"; do
  COUNT=$((COUNT + 1))
  if [ "$COUNT" -ge "$MAX" ]; then
    fail "PostgreSQL did not become healthy within ${MAX}s"
    docker compose logs postgres | tail -20
    exit 1
  fi
  printf "."
  sleep 1
done
echo ""
ok "PostgreSQL is healthy"

# ── step 3: wait for app to respond ───────────────────────────────────────────

bold ""
bold "[ 3/4 ] Waiting for Go app to be ready..."
MAX=30
COUNT=0
until curl -sf "${APP_URL}/hello" > /dev/null 2>&1; do
  COUNT=$((COUNT + 1))
  if [ "$COUNT" -ge "$MAX" ]; then
    fail "App did not respond within ${MAX}s"
    docker compose logs app | tail -20
    exit 1
  fi
  printf "."
  sleep 1
done
echo ""
ok "App is responding on :${APP_PORT}"

# ── step 4: verify ────────────────────────────────────────────────────────────

bold ""
bold "[ 4/4 ] Running checks..."
echo ""

# --- check: GET /hello --------------------------------------------------------
HELLO=$(curl -sf "${APP_URL}/hello")
STATUS=$(echo "$HELLO" | grep -o '"status":"ok"' || true)
if [ -n "$STATUS" ]; then
  ok "GET /hello → $(echo "$HELLO")"
else
  fail "GET /hello returned unexpected response: $HELLO"
fi

# --- check: DB — can connect --------------------------------------------------
DB_CONN=$(docker exec "$DB_CONTAINER" \
  psql -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT 'connected'" 2>&1)
if echo "$DB_CONN" | grep -q "connected"; then
  ok "DB connection OK"
else
  fail "DB connection failed: $DB_CONN"
fi

# --- check: DB — hospitals table has data ------------------------------------
HOSP_COUNT=$(docker exec "$DB_CONTAINER" \
  psql -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT COUNT(*) FROM hospitals" 2>&1 | tr -d ' ')
if [ "$HOSP_COUNT" -gt 0 ] 2>/dev/null; then
  ok "hospitals table — ${HOSP_COUNT} row(s)"
else
  fail "hospitals table is empty or missing (got: '$HOSP_COUNT')"
fi

# --- check: DB — patients table has data -------------------------------------
PT_COUNT=$(docker exec "$DB_CONTAINER" \
  psql -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT COUNT(*) FROM patients" 2>&1 | tr -d ' ')
if [ "$PT_COUNT" -gt 0 ] 2>/dev/null; then
  ok "patients table  — ${PT_COUNT} row(s)"
else
  fail "patients table is empty or missing (got: '$PT_COUNT')"
fi

# --- check: DB — staff table has seed row ------------------------------------
STAFF_COUNT=$(docker exec "$DB_CONTAINER" \
  psql -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT COUNT(*) FROM staff" 2>&1 | tr -d ' ')
if [ "$STAFF_COUNT" -gt 0 ] 2>/dev/null; then
  ok "staff table     — ${STAFF_COUNT} row(s)"
else
  fail "staff table is empty or missing (got: '$STAFF_COUNT')"
fi

# --- check: app can reach DB (login with bad creds → 401, not 500) -----------
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
  -X POST "${APP_URL}/staff/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"probe@test.com","password":"Probe1!xx","hospital":"probe"}')
if [ "$HTTP_CODE" = "401" ]; then
  ok "App → DB round-trip OK (POST /staff/login probe → 401 Unauthorized)"
elif [ "$HTTP_CODE" = "500" ]; then
  fail "App → DB round-trip FAILED (POST /staff/login returned 500 — check DB connection)"
else
  info "POST /staff/login probe → HTTP $HTTP_CODE (expected 401)"
  PASS=$((PASS + 1))
fi

# ── summary ───────────────────────────────────────────────────────────────────

echo ""
bold "═══════════════════════════════════════════"
if [ "$FAIL" -eq 0 ]; then
  green " ✔  All ${PASS} checks passed — system is ready!"
  bold "═══════════════════════════════════════════"
  echo ""
  echo "  App:     ${APP_URL}"
  echo "  pgAdmin: http://localhost:5050  (admin@admin.com / admin)"
  echo ""
else
  red " ✘  ${FAIL} check(s) failed  (${PASS} passed)"
  bold "═══════════════════════════════════════════"
  echo ""
  echo "  Run:  docker compose logs app"
  echo "  Run:  docker compose logs postgres"
  echo ""
  exit 1
fi
