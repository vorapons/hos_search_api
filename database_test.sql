-- ============================================================
-- DATABASE TEST SCRIPT
-- Run with:
--   docker exec -i pt_search_hos_db psql -U postgres -d pt_search_hos < database_test.sql
-- ============================================================

\echo ''
\echo '======================================================'
\echo ' DATABASE TESTS'
\echo '======================================================'

-- ============================================================
-- HELPER: track pass/fail
-- ============================================================
CREATE TEMP TABLE test_results (
  test_name  TEXT,
  status     TEXT,  -- 'PASS' or 'FAIL'
  detail     TEXT
);

CREATE OR REPLACE FUNCTION assert(p_name TEXT, p_condition BOOLEAN, p_detail TEXT DEFAULT '')
RETURNS VOID AS $$
BEGIN
  IF p_condition THEN
    INSERT INTO test_results VALUES (p_name, 'PASS', p_detail);
  ELSE
    INSERT INTO test_results VALUES (p_name, 'FAIL', p_detail);
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ============================================================
-- 1. TABLES EXIST
-- ============================================================
\echo ''
\echo '-- 1. Tables exist'

SELECT assert('hospitals table exists',
  EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hospitals'));

SELECT assert('staff table exists',
  EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'staff'));

SELECT assert('patients table exists',
  EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'patients'));


-- ============================================================
-- 2. VIEWS EXIST
-- ============================================================
\echo '-- 2. Views exist'

SELECT assert('active_patients view exists',
  EXISTS (SELECT 1 FROM information_schema.views WHERE table_name = 'active_patients'));

SELECT assert('active_staff view exists',
  EXISTS (SELECT 1 FROM information_schema.views WHERE table_name = 'active_staff'));


-- ============================================================
-- 3. COLUMNS EXIST — patients
-- ============================================================
\echo '-- 3. Patient columns exist'

SELECT assert('patients.first_name_th exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'first_name_th'));
SELECT assert('patients.middle_name_th exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'middle_name_th'));
SELECT assert('patients.last_name_th exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'last_name_th'));
SELECT assert('patients.first_name_en exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'first_name_en'));
SELECT assert('patients.middle_name_en exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'middle_name_en'));
SELECT assert('patients.last_name_en exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'last_name_en'));
SELECT assert('patients.date_of_birth exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'date_of_birth'));
SELECT assert('patients.patient_hn exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'patient_hn'));
SELECT assert('patients.national_id exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'national_id'));
SELECT assert('patients.passport_id exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'passport_id'));
SELECT assert('patients.phone_number exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'phone_number'));
SELECT assert('patients.email exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'email'));
SELECT assert('patients.gender exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'gender'));
SELECT assert('patients.hospital_id exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'hospital_id'));
SELECT assert('patients.deleted_at exists (soft delete)',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'patients' AND column_name = 'deleted_at'));


-- ============================================================
-- 4. COLUMNS EXIST — staff
-- ============================================================
\echo '-- 4. Staff columns exist'

SELECT assert('staff.email exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'staff' AND column_name = 'email'));
SELECT assert('staff.password exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'staff' AND column_name = 'password'));
SELECT assert('staff.hospital_id exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'staff' AND column_name = 'hospital_id'));
SELECT assert('staff.created_at exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'staff' AND column_name = 'created_at'));
SELECT assert('staff.updated_at exists',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'staff' AND column_name = 'updated_at'));
SELECT assert('staff.deleted_at exists (soft delete)',
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'staff' AND column_name = 'deleted_at'));


-- ============================================================
-- 5. SAMPLE DATA COUNTS
-- ============================================================
\echo '-- 5. Sample data counts'

SELECT assert('5 hospitals loaded',
  (SELECT COUNT(*) FROM hospitals) = 5,
  'found: ' || (SELECT COUNT(*)::TEXT FROM hospitals));

SELECT assert('40 patients loaded',
  (SELECT COUNT(*) FROM patients) = 40,
  'found: ' || (SELECT COUNT(*)::TEXT FROM patients));

SELECT assert('20 Thai patients (have Thai names)',
  (SELECT COUNT(*) FROM patients WHERE first_name_th IS NOT NULL) = 20,
  'found: ' || (SELECT COUNT(*)::TEXT FROM patients WHERE first_name_th IS NOT NULL));

SELECT assert('20 foreign patients (no Thai names)',
  (SELECT COUNT(*) FROM patients WHERE first_name_th IS NULL) = 20,
  'found: ' || (SELECT COUNT(*)::TEXT FROM patients WHERE first_name_th IS NULL));

SELECT assert('3 patients have middle_name_en',
  (SELECT COUNT(*) FROM patients WHERE middle_name_en IS NOT NULL) = 3,
  'found: ' || (SELECT COUNT(*)::TEXT FROM patients WHERE middle_name_en IS NOT NULL));

SELECT assert('patients spread across all 5 hospitals',
  (SELECT COUNT(DISTINCT hospital_id) FROM patients) = 5);


-- ============================================================
-- 6. CHECK CONSTRAINT — patient name (must have first+last in Thai OR English)
-- ============================================================
\echo '-- 6. Constraints: patient name'

SELECT assert('chk_patient_name: reject no names at all', (
  SELECT COUNT(*) FROM (
    SELECT assert_fail FROM (VALUES (1)) v(assert_fail)
    WHERE NOT EXISTS (
      SELECT 1 FROM patients
      WHERE first_name_th IS NULL AND last_name_th IS NULL
        AND first_name_en IS NULL AND last_name_en IS NULL
    )
  ) t
) = 1);

-- Try to insert a patient with no names — should fail
DO $$
BEGIN
  BEGIN
    INSERT INTO patients (hospital_id, national_id)
    VALUES ('BKH01', '9999999999991');
    INSERT INTO test_results VALUES ('chk_patient_name: blocks insert with no names', 'FAIL', 'constraint did not fire');
  EXCEPTION WHEN check_violation THEN
    INSERT INTO test_results VALUES ('chk_patient_name: blocks insert with no names', 'PASS', '');
  END;
END;
$$;

-- Try to insert with only first_name_th but no last_name_th — should fail
DO $$
BEGIN
  BEGIN
    INSERT INTO patients (hospital_id, first_name_th, national_id)
    VALUES ('BKH01', 'ทดสอบ', '9999999999992');
    INSERT INTO test_results VALUES ('chk_patient_name: blocks first_name_th without last_name_th', 'FAIL', 'constraint did not fire');
  EXCEPTION WHEN check_violation THEN
    INSERT INTO test_results VALUES ('chk_patient_name: blocks first_name_th without last_name_th', 'PASS', '');
  END;
END;
$$;

-- Valid insert with Thai name — should succeed
DO $$
BEGIN
  INSERT INTO patients (hospital_id, first_name_th, last_name_th, national_id)
  VALUES ('BKH01', 'ทดสอบ', 'ระบบ', '9999999999993');
  INSERT INTO test_results VALUES ('chk_patient_name: allows Thai first+last name', 'PASS', '');
  DELETE FROM patients WHERE national_id = '9999999999993';
EXCEPTION WHEN OTHERS THEN
  INSERT INTO test_results VALUES ('chk_patient_name: allows Thai first+last name', 'FAIL', SQLERRM);
END;
$$;

-- Valid insert with English name — should succeed
DO $$
BEGIN
  INSERT INTO patients (hospital_id, first_name_en, last_name_en, national_id)
  VALUES ('BKH01', 'Test', 'User', '9999999999994');
  INSERT INTO test_results VALUES ('chk_patient_name: allows English first+last name', 'PASS', '');
  DELETE FROM patients WHERE national_id = '9999999999994';
EXCEPTION WHEN OTHERS THEN
  INSERT INTO test_results VALUES ('chk_patient_name: allows English first+last name', 'FAIL', SQLERRM);
END;
$$;


-- ============================================================
-- 7. CHECK CONSTRAINT — patient identity (national_id OR passport_id required)
-- ============================================================
\echo '-- 7. Constraints: patient identity'

-- Should fail — no identity document
DO $$
BEGIN
  BEGIN
    INSERT INTO patients (hospital_id, first_name_en, last_name_en)
    VALUES ('BKH01', 'No', 'Identity');
    INSERT INTO test_results VALUES ('chk_patient_identity: blocks insert with no identity', 'FAIL', 'constraint did not fire');
  EXCEPTION WHEN check_violation THEN
    INSERT INTO test_results VALUES ('chk_patient_identity: blocks insert with no identity', 'PASS', '');
  END;
END;
$$;


-- ============================================================
-- 8. UNIQUE CONSTRAINTS
-- ============================================================
\echo '-- 8. Unique constraints'

-- national_id must be globally unique
DO $$
DECLARE v_existing TEXT;
BEGIN
  SELECT national_id INTO v_existing FROM patients WHERE national_id IS NOT NULL LIMIT 1;
  BEGIN
    INSERT INTO patients (hospital_id, first_name_en, last_name_en, national_id)
    VALUES ('BKH02', 'Dup', 'Test', v_existing);
    INSERT INTO test_results VALUES ('unique: national_id is globally unique', 'FAIL', 'duplicate allowed');
  EXCEPTION WHEN unique_violation THEN
    INSERT INTO test_results VALUES ('unique: national_id is globally unique', 'PASS', '');
  END;
END;
$$;

-- passport_id must be globally unique
DO $$
DECLARE v_existing TEXT;
BEGIN
  SELECT passport_id INTO v_existing FROM patients WHERE passport_id IS NOT NULL LIMIT 1;
  BEGIN
    INSERT INTO patients (hospital_id, first_name_en, last_name_en, passport_id)
    VALUES ('BKH02', 'Dup', 'Test', v_existing);
    INSERT INTO test_results VALUES ('unique: passport_id is globally unique', 'FAIL', 'duplicate allowed');
  EXCEPTION WHEN unique_violation THEN
    INSERT INTO test_results VALUES ('unique: passport_id is globally unique', 'PASS', '');
  END;
END;
$$;

-- patient_hn must be unique per hospital (same HN in different hospital is OK)
DO $$
BEGIN
  BEGIN
    INSERT INTO patients (hospital_id, first_name_en, last_name_en, passport_id, patient_hn)
    VALUES ('BKH01', 'Dup', 'HN', 'XX99999999', 'BKH-0001');
    INSERT INTO test_results VALUES ('unique: patient_hn unique within hospital', 'FAIL', 'duplicate HN in same hospital allowed');
  EXCEPTION WHEN unique_violation THEN
    INSERT INTO test_results VALUES ('unique: patient_hn unique within hospital', 'PASS', '');
  END;
END;
$$;

DO $$
BEGIN
  BEGIN
    INSERT INTO patients (hospital_id, first_name_en, last_name_en, passport_id, patient_hn)
    VALUES ('BKH02', 'Same', 'HN', 'XY99999999', 'BKH-0001');
    INSERT INTO test_results VALUES ('unique: same patient_hn allowed in different hospital', 'PASS', '');
    DELETE FROM patients WHERE passport_id = 'XY99999999';
  EXCEPTION WHEN unique_violation THEN
    INSERT INTO test_results VALUES ('unique: same patient_hn allowed in different hospital', 'FAIL', 'rejected same HN in different hospital');
  END;
END;
$$;


-- ============================================================
-- 9. SOFT DELETE — deleted_at
-- ============================================================
\echo '-- 9. Soft delete & active_patients view'

DO $$
DECLARE v_id UUID;
BEGIN
  -- Insert a test patient
  INSERT INTO patients (hospital_id, first_name_en, last_name_en, passport_id)
  VALUES ('BKH01', 'Soft', 'Delete', 'SD99999999')
  RETURNING id INTO v_id;

  -- Should appear in active_patients
  PERFORM assert('soft delete: new patient appears in active_patients',
    EXISTS (SELECT 1 FROM active_patients WHERE id = v_id));

  -- Soft delete
  UPDATE patients SET deleted_at = NOW() WHERE id = v_id;

  -- Should NOT appear in active_patients
  PERFORM assert('soft delete: deleted patient hidden from active_patients',
    NOT EXISTS (SELECT 1 FROM active_patients WHERE id = v_id));

  -- Clean up
  DELETE FROM patients WHERE id = v_id;
END;
$$;


-- ============================================================
-- 10. TRIGGER — updated_at auto-updates
-- ============================================================
\echo '-- 10. Trigger: updated_at auto-updates'

DO $$
DECLARE
  v_id       UUID;
  v_before   TIMESTAMPTZ;
  v_after    TIMESTAMPTZ;
BEGIN
  INSERT INTO patients (hospital_id, first_name_en, last_name_en, passport_id)
  VALUES ('BKH01', 'Trigger', 'Test', 'TR99999999')
  RETURNING id, updated_at INTO v_id, v_before;

  PERFORM pg_sleep(0.01);
  UPDATE patients SET gender = 'other' WHERE id = v_id;
  SELECT updated_at INTO v_after FROM patients WHERE id = v_id;

  PERFORM assert('trigger: updated_at changes on UPDATE', v_after > v_before);

  DELETE FROM patients WHERE id = v_id;
END;
$$;


-- ============================================================
-- 11. STAFF — unique email constraint
-- ============================================================
\echo '-- 11. Staff: unique email'

DO $$
DECLARE v_id INT;
BEGIN
  -- Insert a test staff member
  INSERT INTO staff (hospital_id, email, password)
  VALUES ('BKH01', 'teststaff@test.com', 'hashed')
  RETURNING id INTO v_id;

  -- Try duplicate email — should fail
  BEGIN
    INSERT INTO staff (hospital_id, email, password)
    VALUES ('BKH02', 'teststaff@test.com', 'hashed');
    INSERT INTO test_results VALUES ('staff: email is unique', 'FAIL', 'duplicate email allowed');
  EXCEPTION WHEN unique_violation THEN
    INSERT INTO test_results VALUES ('staff: email is unique', 'PASS', '');
  END;

  DELETE FROM staff WHERE id = v_id;
END;
$$;


-- ============================================================
-- 12. STAFF — soft delete & active_staff view
-- ============================================================
\echo '-- 12. Staff: soft delete & active_staff view'

DO $$
DECLARE v_id INT;
BEGIN
  INSERT INTO staff (hospital_id, email, password)
  VALUES ('BKH01', 'softdelete@test.com', 'hashed')
  RETURNING id INTO v_id;

  PERFORM assert('staff soft delete: new staff appears in active_staff',
    EXISTS (SELECT 1 FROM active_staff WHERE id = v_id));

  UPDATE staff SET deleted_at = NOW() WHERE id = v_id;

  PERFORM assert('staff soft delete: deleted staff hidden from active_staff',
    NOT EXISTS (SELECT 1 FROM active_staff WHERE id = v_id));

  DELETE FROM staff WHERE id = v_id;
END;
$$;


-- ============================================================
-- 13. STAFF — trigger updated_at auto-updates
-- ============================================================
\echo '-- 13. Staff: trigger updated_at'

DO $$
DECLARE
  v_id      INT;
  v_before  TIMESTAMPTZ;
  v_after   TIMESTAMPTZ;
BEGIN
  INSERT INTO staff (hospital_id, email, password)
  VALUES ('BKH01', 'triggertest@test.com', 'hashed')
  RETURNING id, updated_at INTO v_id, v_before;

  PERFORM pg_sleep(0.01);
  UPDATE staff SET password = 'newhashedpassword' WHERE id = v_id;
  SELECT updated_at INTO v_after FROM staff WHERE id = v_id;

  PERFORM assert('staff trigger: updated_at changes on UPDATE', v_after > v_before);

  DELETE FROM staff WHERE id = v_id;
END;
$$;


-- ============================================================
-- 14. STAFF — hospital_id foreign key
-- ============================================================
\echo '-- 14. Staff: hospital_id foreign key'

DO $$
BEGIN
  BEGIN
    INSERT INTO staff (hospital_id, email, password)
    VALUES ('XXXXX', 'fktest@test.com', 'hashed');
    INSERT INTO test_results VALUES ('staff: rejects invalid hospital_id', 'FAIL', 'FK not enforced');
  EXCEPTION WHEN foreign_key_violation THEN
    INSERT INTO test_results VALUES ('staff: rejects invalid hospital_id', 'PASS', '');
  END;
END;
$$;


-- ============================================================
-- RESULTS SUMMARY
-- ============================================================
\echo ''
\echo '======================================================'
\echo ' RESULTS'
\echo '======================================================'

SELECT
  status,
  COUNT(*) AS count
FROM test_results
GROUP BY status
ORDER BY status;

\echo ''
\echo '-- Failed tests (if any):'
SELECT test_name, detail
FROM test_results
WHERE status = 'FAIL'
ORDER BY test_name;

\echo ''
\echo '-- All results:'
SELECT
  CASE status WHEN 'PASS' THEN '✓' ELSE '✗' END AS result,
  test_name,
  detail
FROM test_results
ORDER BY status DESC, test_name;

-- Cleanup
DROP FUNCTION assert(TEXT, BOOLEAN, TEXT);
DROP TABLE test_results;
