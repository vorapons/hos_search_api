# 🧪 Service Layer Test Summary

📦 Package: `services_test`

---

## 👤 StaffService — `staff_service_test.go`

### 🔑 Login

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestLogin_Success` | ✅ positive | valid email + correct password | no error · JWT token returned |
| 2 | `TestLogin_EmptyEmail` | ❌ negative | empty email | `ErrInvalidInput` |
| 3 | `TestLogin_EmptyPassword` | ❌ negative | empty password | `ErrInvalidInput` |
| 4 | `TestLogin_NotEmailFormat` | ❌ negative | non-email string (`"notanemail"`) | `ErrUnauthorized` · repo IS called (no format short-circuit) |
| 5 | `TestLogin_UserNotFound` | ❌ negative | email not registered in DB | `ErrUnauthorized` |
| 6 | `TestLogin_WrongPassword` | ❌ negative | correct email + wrong password | `ErrUnauthorized` |
| 7 | `TestLogin_DBError` | ❌ negative | repo returns DB error | error propagated · not `ErrUnauthorized` |

> ⚠️ **Security note:** Tests 4 and 5 both return `ErrUnauthorized` intentionally — prevents user enumeration. A non-email format goes to the repo just like any other login attempt.

---

### 🏥 CreateStaff

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestCreateStaff_Success` | ✅ positive | valid email + strong password + known hospital | no error · JWT token returned |
| 2 | `TestCreateStaff_EmptyFields` | ❌ negative | any field empty (3 sub-cases) | `ErrInvalidInput` |
| 3 | `TestCreateStaff_InvalidEmail` | ❌ negative | malformed email (`"not-an-email"`) | `ErrInvalidInput` |
| 4 | `TestCreateStaff_WeakPassword` | ❌ negative | 5 weak password patterns (no upper/lower/digit/special, too short) | `ErrInvalidInput` |
| 5 | `TestCreateStaff_AlreadyExists` | ❌ negative | email already registered | `ErrStaffExists` |
| 6 | `TestCreateStaff_HospitalNotFound` | ❌ negative | hospital name not in DB | `ErrHospitalNotFound` |
| 7 | `TestCreateStaff_FindByEmailDBError` | ❌ negative | DB error on email duplicate check | error propagated |
| 8 | `TestCreateStaff_FindHospitalDBError` | ❌ negative | DB error on hospital lookup | error propagated |
| 9 | `TestCreateStaff_PasswordTooLong` | ❌ negative | password > 72 bytes (bcrypt limit) | error propagated |
| 10 | `TestCreateStaff_DBCreateError` | ❌ negative | DB error on staff insert | error propagated |

> ⚠️ **Note:** `TestCreateStaff_PasswordTooLong` — a 73-byte password passes `isStrongPassword()` but bcrypt silently truncates at 72 bytes, making the stored hash unreliable. The test confirms an error is returned.

---

### 🚪 Logout

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestLogout_PersistsToDBAndMemory` | ✅ positive | valid JWT token | no error · token blacklisted in DB and in-memory cache |
| 2 | `TestLogout_InvalidTokenFallsBackTo24h` | ❌ negative (edge) | unparseable token string | no error · falls back to 24h expiry |
| 3 | `TestLogout_TokenWithNoExpiryClaim` | ❌ negative (edge) | valid JWT but missing `exp` claim | no error · falls back to 24h expiry |
| 4 | `TestLogout_DBError` | ❌ negative | DB error when blacklisting | error returned |

---

### 🔍 IsTokenBlacklisted

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestIsTokenBlacklisted_FalseInitially` | ✅ positive | any token on fresh service | `false` |

---

### 📥 LoadBlacklist

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestLoadBlacklist_LoadsTokensIntoMemory` | ✅ positive | DB returns 2 tokens | no error · both tokens blacklisted in memory |
| 2 | `TestLoadBlacklist_DBError` | ❌ negative | DB error on load | error returned |

---

## 🩺 PatientService — `patient_service_test.go`

### 🔍 GetPatientByID

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestGetPatientByID_FoundByNationalID` | ✅ positive | valid national ID + hospital ID | no error · patient returned |
| 2 | `TestGetPatientByID_FoundByPassportID` | ✅ positive | valid passport ID + hospital ID | no error · patient returned |
| 3 | `TestGetPatientByID_NotFound` | ❌ negative | ID not in DB (repo returns nil) | `ErrNotFound` |
| 4 | `TestGetPatientByID_EmptyID` | ❌ negative | empty patient ID | `ErrInvalidInput` |
| 5 | `TestGetPatientByID_EmptyHospitalID` | ❌ negative | empty hospital ID | `ErrInvalidInput` |
| 6 | `TestGetPatientByID_DBError` | ❌ negative | repo returns DB error | error propagated · not `ErrNotFound` |

---

### 🔎 GetPatientByCondition

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestGetPatientByCondition_Success` | ✅ positive | `last_name: "Smith"` + hospital ID | no error · 2 patients returned |
| 2 | `TestGetPatientByCondition_EmptyResult` | ✅ positive | valid condition, no matches | no error · empty slice |
| 3 | `TestGetPatientByCondition_ByNationalID` | ✅ positive | `national_id` condition | no error · 1 patient returned |
| 4 | `TestGetPatientByCondition_ByDateOfBirth` | ✅ positive | `date_of_birth` condition | no error · 1 patient returned |
| 5 | `TestGetPatientByCondition_NoCondition` | ❌ negative | all fields nil | `ErrInvalidInput` |
| 6 | `TestGetPatientByCondition_EmptyHospitalID` | ❌ negative | empty hospital ID | `ErrInvalidInput` |
| 7 | `TestGetPatientByCondition_DBError` | ❌ negative | repo returns DB error | error propagated |

---

## 📊 Totals

| Service | Method | ✅ Positive | ❌ Negative | 🔢 Total |
|---------|--------|------------|------------|---------|
| 👤 StaffService | Login | 1 | 6 | 7 |
| | CreateStaff | 1 | 9 | 10 |
| | Logout | 1 | 3 | 4 |
| | IsTokenBlacklisted | 1 | 0 | 1 |
| | LoadBlacklist | 1 | 1 | 2 |
| 🩺 PatientService | GetPatientByID | 2 | 4 | 6 |
| | GetPatientByCondition | 4 | 3 | 7 |
| **Total** | | **11** | **26** | **37** |
