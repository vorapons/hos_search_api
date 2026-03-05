# 🧪 Staff Handler Test Summary — `handler/gin/staff_test.go`

📦 Package: `ginhandler_test`

---

## 🔑 POST /staff/login

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestLogin_Success` | ✅ positive | valid email + correct password | 200 · `token` in body |
| 2 | `TestLogin_BadBody` | ❌ negative | malformed JSON | 400 |
| 3 | `TestLogin_InvalidInput` | ❌ negative | empty email + empty password | 400 · `INVALID_INPUT` |
| 4 | `TestLogin_Unauthorized` | ❌ negative | valid email + wrong password | 401 · `UNAUTHORIZED` |
| 5 | `TestLogin_InvalidEmailFormat` | ❌ negative | non-email string as login (e.g. `"notanemail"`) | 401 · `UNAUTHORIZED` |
| 6 | `TestLogin_UserNotFound` | ❌ negative | email not registered in DB | 401 · `UNAUTHORIZED` |
| 7 | `TestLogin_InternalError` | ❌ negative | service returns unexpected error | 500 |

---

## 🏥 POST /staff/create

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestCreate_Success` | ✅ positive | valid email + strong password + known hospital | 201 · `token` in body |
| 2 | `TestCreate_BadBody` | ❌ negative | malformed JSON | 400 |
| 3 | `TestCreate_InvalidInput` | ❌ negative | all fields empty | 400 · `INVALID_INPUT` |
| 4 | `TestCreate_InvalidEmailFormat` | ❌ negative | login is not a valid email (e.g. `"notanemail"`) | 400 · `INVALID_INPUT` |
| 5 | `TestCreate_WeakPassword` | ❌ negative | password too weak (e.g. `"weak"`) | 400 · `INVALID_INPUT` |
| 6 | `TestCreate_StaffExists` | ❌ negative | email already registered | 409 · `CONFLICT` |
| 7 | `TestCreate_HospitalNotFound` | ❌ negative | hospital name not in DB | 404 · `NOT_FOUND` |
| 8 | `TestCreate_HospitalNameWrongCase` | ❌ negative | hospital name wrong case (`"bangkok hospital"`) | 404 · `NOT_FOUND` |
| 9 | `TestCreate_InternalError` | ❌ negative | service returns unexpected error | 500 |

> ⚠️ **Note:** Hospital name lookup is case-sensitive (`=` in PostgreSQL). `"bangkok hospital"` ≠ `"Bangkok Hospital"`.

---

## 👋 GET /staff/hello  *(requires JWT)*

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestHello_Success` | ✅ positive | valid non-blacklisted JWT | 200 · `login`, `hospital`, `expires_at` in body |
| 2 | `TestHello_NoAuthHeader` | ❌ negative | no Authorization header | 401 |
| 3 | `TestHello_InvalidToken` | ❌ negative | malformed / bad JWT string | 401 |
| 4 | `TestHello_BlacklistedToken` | ❌ negative | valid JWT but token is blacklisted | 401 · `UNAUTHORIZED` |

---

## 🚪 GET /staff/logout  *(requires JWT)*

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestLogout_Success` | ✅ positive | valid non-blacklisted JWT | 200 · `"Logged out successfully"` |
| 2 | `TestLogout_NoAuthHeader` | ❌ negative | no Authorization header | 401 |
| 3 | `TestLogout_ServiceError` | ❌ negative | service fails to blacklist token | 500 |

---

## 📊 Totals

| Endpoint | ✅ Positive | ❌ Negative | 🔢 Total |
|----------|------------|------------|---------|
| 🔑 POST /staff/login | 1 | 6 | 7 |
| 🏥 POST /staff/create | 1 | 8 | 9 |
| 👋 GET /staff/hello | 1 | 3 | 4 |
| 🚪 GET /staff/logout | 1 | 2 | 3 |
| **Total** | **4** | **19** | **23** |
