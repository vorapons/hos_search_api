-- ============================================================
-- Hospital REST API Database Schema
-- PostgreSQL
-- ============================================================

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";


-- ============================================================
-- TRIGGER FUNCTION: auto-update updated_at on every UPDATE
-- ============================================================
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = clock_timestamp();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ============================================================
-- TABLE: hospitals
-- ============================================================
CREATE TABLE hospitals (
  id          VARCHAR(5) PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,

  created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  deleted_at  TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TRIGGER set_updated_at
  BEFORE UPDATE ON hospitals
  FOR EACH ROW EXECUTE FUNCTION update_updated_at();


-- ============================================================
-- TABLE: staff
-- ============================================================
CREATE TABLE staff (
  id           SERIAL PRIMARY KEY,
  hospital_id  VARCHAR(5) NOT NULL REFERENCES hospitals(id),

  email        VARCHAR(255) NOT NULL UNIQUE,
  password     VARCHAR(255) NOT NULL,

  created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  deleted_at   TIMESTAMP WITH TIME ZONE DEFAULT NULL,

  CONSTRAINT chk_staff_email_format
    CHECK (email ~* '^[^@]+@[^@]+\.[^@]+$')
);

CREATE TRIGGER set_updated_at
  BEFORE UPDATE ON staff
  FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Index for fast login lookup
CREATE INDEX idx_staff_email       ON staff(email);
CREATE INDEX idx_staff_hospital_id ON staff(hospital_id);


-- ============================================================
-- TABLE: patients
-- ============================================================
CREATE TABLE patients (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hospital_id  VARCHAR(5) NOT NULL REFERENCES hospitals(id),

  -- Thai name (split)
  first_name_th   VARCHAR(100),
  middle_name_th  VARCHAR(100),
  last_name_th    VARCHAR(100),

  -- English name (split)
  first_name_en   VARCHAR(100),
  middle_name_en  VARCHAR(100),
  last_name_en    VARCHAR(100),

  -- Identity
  national_id  VARCHAR(13)  UNIQUE,   -- Thai national ID (nullable)
  passport_id  VARCHAR(20)  UNIQUE,   -- Passport (nullable)
  patient_hn   VARCHAR(20),           -- Hospital number (unique per hospital)

  -- Demographics
  date_of_birth  DATE,
  gender         VARCHAR(10),         -- e.g. 'male', 'female', 'other'
  phone_number   VARCHAR(20),
  email          VARCHAR(255),

  created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  deleted_at  TIMESTAMP WITH TIME ZONE DEFAULT NULL,

  -- at least first+last name in Thai OR English must exist
  CONSTRAINT chk_patient_name
    CHECK (
      (first_name_th IS NOT NULL AND last_name_th IS NOT NULL)
      OR (first_name_en IS NOT NULL AND last_name_en IS NOT NULL)
    ),

  -- at least one identity document must exist
  CONSTRAINT chk_patient_identity
    CHECK (national_id IS NOT NULL OR passport_id IS NOT NULL),

  -- HN is unique within a hospital
  CONSTRAINT uq_patient_hn UNIQUE (hospital_id, patient_hn)
);

CREATE TRIGGER set_updated_at
  BEFORE UPDATE ON patients
  FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Indexes for search performance
CREATE INDEX idx_patients_hospital_id ON patients(hospital_id);
CREATE INDEX idx_patients_national_id ON patients(national_id);
CREATE INDEX idx_patients_passport_id ON patients(passport_id);
CREATE INDEX idx_patients_patient_hn  ON patients(hospital_id, patient_hn);
CREATE INDEX idx_patients_first_th    ON patients(first_name_th);
CREATE INDEX idx_patients_last_th     ON patients(last_name_th);
CREATE INDEX idx_patients_first_en    ON patients(first_name_en);
CREATE INDEX idx_patients_last_en     ON patients(last_name_en);
CREATE INDEX idx_patients_phone       ON patients(phone_number);
CREATE INDEX idx_patients_email       ON patients(email);


-- ============================================================
-- TABLE: blacklisted_tokens
-- ============================================================
CREATE TABLE blacklisted_tokens (
  token      TEXT                     PRIMARY KEY,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_blacklisted_tokens_expires_at ON blacklisted_tokens(expires_at);

-- Auto-cleanup: delete tokens older than 7 days on every INSERT
CREATE OR REPLACE FUNCTION cleanup_old_blacklisted_tokens()
RETURNS TRIGGER AS $$
BEGIN
  DELETE FROM blacklisted_tokens WHERE created_at < NOW() - INTERVAL '7 days';
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_cleanup_blacklisted_tokens
  AFTER INSERT ON blacklisted_tokens
  FOR EACH STATEMENT EXECUTE FUNCTION cleanup_old_blacklisted_tokens();


-- ============================================================
-- VIEW: active_patients (excludes soft-deleted)
-- ============================================================
CREATE VIEW active_patients AS
  SELECT * FROM patients
  WHERE deleted_at IS NULL;

-- ============================================================
-- VIEW: active_staff (excludes soft-deleted)
-- ============================================================
CREATE VIEW active_staff AS
  SELECT * FROM staff
  WHERE deleted_at IS NULL;


