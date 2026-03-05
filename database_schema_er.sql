-- ============================================================
-- ER DIAGRAM SCHEMA
-- Tables and relations only (no triggers, indexes, views)
-- ============================================================

CREATE TABLE hospitals (
  id          VARCHAR(5)  PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL,
  updated_at  TIMESTAMPTZ NOT NULL,
  deleted_at  TIMESTAMPTZ
);

CREATE TABLE staff (
  id           SERIAL      PRIMARY KEY,
  hospital_id  VARCHAR(5)  NOT NULL REFERENCES hospitals(id),
  email        VARCHAR(255) NOT NULL UNIQUE,
  password     VARCHAR(255) NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL,
  updated_at   TIMESTAMPTZ NOT NULL,
  deleted_at   TIMESTAMPTZ
);

CREATE TABLE patients (
  id              UUID        PRIMARY KEY,
  hospital_id     VARCHAR(5)  NOT NULL REFERENCES hospitals(id),

  first_name_th   VARCHAR(100),
  middle_name_th  VARCHAR(100),
  last_name_th    VARCHAR(100),

  first_name_en   VARCHAR(100),
  middle_name_en  VARCHAR(100),
  last_name_en    VARCHAR(100),

  national_id     VARCHAR(13)  UNIQUE,
  passport_id     VARCHAR(20)  UNIQUE,
  patient_hn      VARCHAR(20),

  date_of_birth   DATE,
  gender          VARCHAR(10),
  phone_number    VARCHAR(20),
  email           VARCHAR(255),

  created_at      TIMESTAMPTZ NOT NULL,
  updated_at      TIMESTAMPTZ NOT NULL,
  deleted_at      TIMESTAMPTZ,

  CONSTRAINT chk_patient_name
    CHECK (
      (first_name_th IS NOT NULL AND last_name_th IS NOT NULL)
      OR (first_name_en IS NOT NULL AND last_name_en IS NOT NULL)
    ),
  CONSTRAINT chk_patient_identity
    CHECK (national_id IS NOT NULL OR passport_id IS NOT NULL),
  CONSTRAINT uq_patient_hn
    UNIQUE (hospital_id, patient_hn)
);

CREATE TABLE blacklisted_tokens (
  token      TEXT        PRIMARY KEY,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);
