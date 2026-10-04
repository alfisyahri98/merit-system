-- =====================================================================
-- MERIT SYSTEM PERSONEL POLRI - SCHEMA POSTGRESQL
-- Jalankan sekali di database kosong:  psql -d merit_db -f schema.sql
-- =====================================================================

BEGIN;

-- Dibutuhkan untuk EXCLUDE constraint (cegah jabatan definitif tumpang tindih)
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- =====================================================================
-- FUNGSI BANTU: auto-update kolom updated_at
-- =====================================================================
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =====================================================================
-- TABEL REFERENSI
-- =====================================================================

CREATE TABLE pangkat (
  id        SERIAL PRIMARY KEY,
  jenis     VARCHAR(10)  NOT NULL CHECK (jenis IN ('POLRI', 'PNS')),
  kode      VARCHAR(20)  NOT NULL UNIQUE,          -- BRIPTU, AKP, III/a
  nama      VARCHAR(100) NOT NULL,                 -- Brigadir Polisi Satu
  kelompok  VARCHAR(50)  NOT NULL,                 -- Bintara, Perwira Pertama, Golongan III
  urutan    INT          NOT NULL CHECK (urutan > 0),
  UNIQUE (jenis, urutan),                          -- urutan tidak boleh dobel per jenis
  UNIQUE (id, jenis)                               -- target composite FK dari personel
);

CREATE TABLE satker (
  id        SERIAL PRIMARY KEY,
  kode      VARCHAR(100) NOT NULL UNIQUE,
  nama      VARCHAR(200) NOT NULL,
  parent_id INT REFERENCES satker(id) ON DELETE RESTRICT,
  tingkat   VARCHAR(20)  NOT NULL CHECK (tingkat IN
            ('MABES','SATKER_MABES','POLDA','POLRES','POLSEK','DIREKTORAT','SUBDIREKTORAT',
             'BIRO','BIDANG','SUBBIDANG','BAGIAN','SUBBAGIAN','SATUAN','UNIT')),
  CHECK (parent_id IS NULL OR parent_id <> id)     -- tidak boleh jadi induk dirinya sendiri
);
CREATE INDEX idx_satker_parent ON satker(parent_id);

CREATE TABLE fungsi (
  id    SERIAL PRIMARY KEY,
  nama  VARCHAR(100) NOT NULL UNIQUE               -- Reskrim, Narkoba, Lantas, Samapta
);

CREATE TABLE nivelering (
  id      SERIAL PRIMARY KEY,
  kode    VARCHAR(10) NOT NULL UNIQUE,             -- IA, IB, IIA, IIIB, dst
  urutan  INT NOT NULL UNIQUE CHECK (urutan > 0)
);

CREATE TABLE master_pendidikan (
  id      SERIAL PRIMARY KEY,
  jenis   VARCHAR(20)  NOT NULL CHECK (jenis IN ('DIKPOL','DIKUM','DIKBANG','PELATIHAN')),
  kode    VARCHAR(50)  NOT NULL UNIQUE,            -- SESPIMMEN, BINTARA, S1, OPS_KOMPUTER
  nama    VARCHAR(200) NOT NULL,
  urutan  INT                                      -- jenjang: Sespimma < Sespimmen < Sespimti
);

-- =====================================================================
-- DATA INTI
-- =====================================================================

CREATE TABLE personel (
  id             BIGSERIAL PRIMARY KEY,
  jenis          VARCHAR(10)  NOT NULL CHECK (jenis IN ('POLRI', 'PNS')),
  nrp_nip        VARCHAR(18)  NOT NULL UNIQUE,     -- VARCHAR: NRP bisa diawali angka 0
  nama           VARCHAR(150) NOT NULL CHECK (btrim(nama) <> ''),
  pangkat_id     INT          NOT NULL,
  satker_id      INT          NOT NULL REFERENCES satker(id) ON DELETE RESTRICT,
  tempat_lahir   VARCHAR(100),
  tanggal_lahir  DATE         NOT NULL CHECK (tanggal_lahir >= DATE '1940-01-01'),
  status         VARCHAR(20)  NOT NULL DEFAULT 'AKTIF'
                 CHECK (status IN ('AKTIF','PENSIUN','MENINGGAL','DIBERHENTIKAN','MUTASI_KELUAR')),
  created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
  deleted_at     TIMESTAMPTZ,

  -- Format nomor identitas tergantung jenis personel
  CONSTRAINT chk_nrp_nip_format CHECK (
    (jenis = 'POLRI' AND nrp_nip ~ '^[0-9]{8}$') OR
    (jenis = 'PNS'   AND nrp_nip ~ '^[0-9]{18}$')
  ),

  -- Composite FK: personel POLRI hanya boleh pakai pangkat POLRI, PNS pakai golongan PNS
  CONSTRAINT fk_personel_pangkat_jenis
    FOREIGN KEY (pangkat_id, jenis) REFERENCES pangkat(id, jenis)
);
CREATE INDEX idx_personel_satker ON personel(satker_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_personel_nama   ON personel(lower(nama));
CREATE TRIGGER trg_personel_updated BEFORE UPDATE ON personel
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Pangkat Jenderal Polisi hanya untuk Kapolri: maksimal 1 personel aktif
CREATE OR REPLACE FUNCTION cek_jenderal_tunggal() RETURNS trigger AS $$
BEGIN
  IF NEW.status = 'AKTIF' AND NEW.deleted_at IS NULL
     AND EXISTS (SELECT 1 FROM pangkat WHERE id = NEW.pangkat_id AND jenis = 'POLRI' AND kode = 'JENDERAL')
     AND EXISTS (SELECT 1 FROM personel p JOIN pangkat pg ON pg.id = p.pangkat_id
                 WHERE pg.kode = 'JENDERAL' AND p.status = 'AKTIF' AND p.deleted_at IS NULL
                   AND p.id <> NEW.id) THEN
    RAISE EXCEPTION 'Sudah ada personel aktif berpangkat Jenderal Polisi (Kapolri)';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_personel_jenderal BEFORE INSERT OR UPDATE ON personel
  FOR EACH ROW EXECUTE FUNCTION cek_jenderal_tunggal();

CREATE TABLE riwayat_jabatan (
  id              BIGSERIAL PRIMARY KEY,
  personel_id     BIGINT       NOT NULL REFERENCES personel(id) ON DELETE RESTRICT,
  nama_jabatan    VARCHAR(300) NOT NULL CHECK (btrim(nama_jabatan) <> ''),
  satker_id       INT          NOT NULL REFERENCES satker(id) ON DELETE RESTRICT,
  fungsi_id       INT          REFERENCES fungsi(id),
  nivelering_id   INT          REFERENCES nivelering(id),
  tmt_mulai       DATE         NOT NULL,
  tmt_selesai     DATE,                            -- NULL = jabatan saat ini
  status_jabatan  VARCHAR(20)  NOT NULL DEFAULT 'DEFINITIF'
                  CHECK (status_jabatan IN ('DEFINITIF','PLT','PLH')),
  nomor_skep      VARCHAR(100),
  keterangan      TEXT,
  created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
  deleted_at      TIMESTAMPTZ,

  CONSTRAINT chk_tmt_urut CHECK (tmt_selesai IS NULL OR tmt_selesai >= tmt_mulai),

  -- Jabatan DEFINITIF satu personel tidak boleh tumpang tindih waktunya.
  -- (PLT/PLH boleh merangkap, jadi dikecualikan)
  CONSTRAINT excl_jabatan_definitif_overlap EXCLUDE USING gist (
    personel_id WITH =,
    daterange(tmt_mulai, tmt_selesai, '[]') WITH &&
  ) WHERE (status_jabatan = 'DEFINITIF' AND deleted_at IS NULL)
);
-- Maksimal 1 jabatan DEFINITIF aktif (tmt_selesai NULL) per personel
CREATE UNIQUE INDEX uq_jabatan_aktif_per_personel
  ON riwayat_jabatan(personel_id)
  WHERE tmt_selesai IS NULL AND status_jabatan = 'DEFINITIF' AND deleted_at IS NULL;
CREATE INDEX idx_riwayat_jabatan_kronologis ON riwayat_jabatan(personel_id, tmt_mulai);
CREATE TRIGGER trg_riwayat_jabatan_updated BEFORE UPDATE ON riwayat_jabatan
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE riwayat_pangkat (
  id           BIGSERIAL PRIMARY KEY,
  personel_id  BIGINT      NOT NULL REFERENCES personel(id) ON DELETE RESTRICT,
  pangkat_id   INT         NOT NULL REFERENCES pangkat(id),
  tmt          DATE        NOT NULL,
  nomor_skep   VARCHAR(100),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (personel_id, pangkat_id),                -- pangkat yang sama tidak dicatat dua kali
  UNIQUE (personel_id, tmt)                        -- tidak ada dua kenaikan di tanggal sama
);
CREATE INDEX idx_riwayat_pangkat_kronologis ON riwayat_pangkat(personel_id, tmt);

CREATE TABLE kualifikasi (
  id                    BIGSERIAL PRIMARY KEY,
  personel_id           BIGINT       NOT NULL REFERENCES personel(id) ON DELETE RESTRICT,
  master_pendidikan_id  INT          NOT NULL REFERENCES master_pendidikan(id),
  institusi             VARCHAR(200),
  jurusan               VARCHAR(150),
  tahun_lulus           INT          NOT NULL CHECK (tahun_lulus BETWEEN 1950 AND 2100),
  nomor_dokumen         VARCHAR(100),
  keterangan            TEXT,
  created_at            TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at            TIMESTAMPTZ  NOT NULL DEFAULT now(),
  deleted_at            TIMESTAMPTZ
);
CREATE INDEX idx_kualifikasi_personel ON kualifikasi(personel_id, tahun_lulus);
CREATE TRIGGER trg_kualifikasi_updated BEFORE UPDATE ON kualifikasi
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- =====================================================================
-- AKSES
-- =====================================================================

CREATE TABLE users (
  id             SERIAL PRIMARY KEY,
  username       VARCHAR(50)  NOT NULL UNIQUE CHECK (username ~ '^[a-z0-9_.]{4,50}$'),
  password_hash  VARCHAR(255) NOT NULL,            -- bcrypt, BUKAN password asli
  role           VARCHAR(20)  NOT NULL CHECK (role IN ('ADMIN_SSDM','OPERATOR')),
  satker_id      INT REFERENCES satker(id),
  is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
  created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
  -- Operator wajib terikat ke satker; admin boleh NULL (akses semua)
  CONSTRAINT chk_operator_punya_satker CHECK (role = 'ADMIN_SSDM' OR satker_id IS NOT NULL)
);

CREATE TABLE api_clients (
  id               SERIAL PRIMARY KEY,
  nama_aplikasi    VARCHAR(150) NOT NULL,
  client_id        VARCHAR(64)  NOT NULL UNIQUE,
  secret_hash      VARCHAR(255) NOT NULL,
  scopes           TEXT[]       NOT NULL DEFAULT '{}',
  akses_satker_id  INT REFERENCES satker(id),      -- NULL = semua satker
  is_active        BOOLEAN      NOT NULL DEFAULT TRUE,
  created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
  -- Scope hanya boleh dari daftar yang dikenal sistem
  CONSTRAINT chk_scopes_valid CHECK (scopes <@ ARRAY[
    'personel:read','personel:write','personel:delete',
    'jabatan:read','jabatan:write','jabatan:delete',
    'pangkat:read','pangkat:write',
    'kualifikasi:read','kualifikasi:write','kualifikasi:delete'
  ]::TEXT[])
);

-- =====================================================================
-- FUNGSI: semua satker turunan (untuk scope akses operator / api client)
-- Contoh: SELECT * FROM satker_turunan(2);  -- Polda Metro + semua di bawahnya
-- =====================================================================
CREATE OR REPLACE FUNCTION satker_turunan(root_id INT)
RETURNS TABLE (id INT) AS $$
  WITH RECURSIVE pohon AS (
    SELECT s.id FROM satker s WHERE s.id = root_id
    UNION ALL
    SELECT s.id FROM satker s JOIN pohon p ON s.parent_id = p.id
  )
  SELECT pohon.id FROM pohon;
$$ LANGUAGE sql STABLE;

-- =====================================================================
-- VIEW: jabatan saat ini per personel (memudahkan query profil & list)
-- =====================================================================
CREATE VIEW v_jabatan_saat_ini AS
SELECT rj.personel_id, rj.id AS riwayat_jabatan_id, rj.nama_jabatan,
       rj.satker_id, rj.tmt_mulai,
       age(CURRENT_DATE, rj.tmt_mulai) AS lama_jabatan
FROM riwayat_jabatan rj
WHERE rj.tmt_selesai IS NULL
  AND rj.status_jabatan = 'DEFINITIF'
  AND rj.deleted_at IS NULL;

COMMIT;
