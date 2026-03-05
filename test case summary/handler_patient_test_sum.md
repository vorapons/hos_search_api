# 🧪 Patient Handler Test Summary — `handler/gin/patient_test.go`

📦 Package: `ginhandler_test`

> 🔐 All endpoints except `GET /hello` require a valid JWT (`Authorization: Bearer <token>`).
> The `hospital_id` claim from the JWT is used to scope all patient queries.

---

## 🔍 GET /patient/search/:id  *(requires JWT)*

> Searches for a single patient by either **national ID** or **passport ID**.

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestGetByID_FoundByNationalID` | ✅ positive | valid national ID (`1234567890123`) | 200 · `national_id` in body |
| 2 | `TestGetByID_FoundByPassportID` | ✅ positive | valid passport ID (`AB123456`) | 200 · `passport_id` in body |
| 3 | `TestGetByID_NotFound` | ❌ negative | ID not in DB | 404 · `NOT_FOUND` |
| 4 | `TestGetByID_NoAuthHeader` | ❌ negative | no Authorization header | 401 |
| 5 | `TestGetByID_InvalidInput` | ❌ negative | service rejects ID (handler branch coverage) | 400 · `INVALID_INPUT` |
| 6 | `TestGetByID_InternalError` | ❌ negative | service returns unexpected error | 500 |

---

## 🔎 POST /patient/search  *(requires JWT)*

> Searches for patients by one or more conditions (name, national ID, passport ID, DOB, etc.).

| # | Test Name | Type | Input | Expected |
|---|-----------|------|-------|----------|
| 1 | `TestSearch_Success` | ✅ positive | `last_name: "Smith"` → mixed Thai + Japanese patients | 200 · list of 2 · Thai fields + `national_id` present; foreign `passport_id` present |
| 2 | `TestSearch_EmptyResult` | ✅ positive | valid condition, no matches | 200 · empty list `[]` |
| 3 | `TestSearch_NoCondition` | ❌ negative | empty JSON body `{}` | 400 · `INVALID_INPUT` |
| 4 | `TestSearch_BadBody` | ❌ negative | non-object JSON (e.g. `"not-an-object"`) | 400 |
| 5 | `TestSearch_NoAuthHeader` | ❌ negative | no Authorization header | 401 |
| 6 | `TestSearch_InternalError` | ❌ negative | service returns unexpected error | 500 |

### 🗂️ Test Data in `TestSearch_Success`

| Patient | Nationality | ID Type | HN | Key Assertions |
|---------|-------------|---------|-----|----------------|
| Somchai Jaidee (สมชาย ใจดี) | 🇹🇭 Thai | `national_id` | `BKH-0001` | `first_name_th`, `last_name_th` present · `passport_id` is null |
| Yuki Tanaka | 🇯🇵 Japanese | `passport_id` | `BKH-0005` | `first_name_th`, `last_name_th` null · `national_id` is null |

---

## 📊 Totals

| Endpoint | ✅ Positive | ❌ Negative | 🔢 Total |
|----------|------------|------------|---------|
| 🔍 GET /patient/search/:id | 2 | 4 | 6 |
| 🔎 POST /patient/search | 2 | 4 | 6 |
| **Total** | **4** | **8** | **12** |
