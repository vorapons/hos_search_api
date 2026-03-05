-- ============================================================
-- SAMPLE DATA
-- ============================================================


-- ============================================================
-- HOSPITALS (5)
-- ============================================================
INSERT INTO hospitals (id, name) VALUES
  ('BKH01', 'Bangkok Hospital'),
  ('BKH02', 'Bangkok Hospital Samui'),
  ('TBR01', 'Thonburi Hospital'),
  ('VIM01', 'Vimut Hospital'),
  ('SMJ01', 'Samitivej Srinakarin');


-- ============================================================
-- PATIENTS — Thai (20)
-- Both Thai and English first+last name; national_id; other fields random
-- ============================================================
INSERT INTO patients (
  hospital_id,
  first_name_th, last_name_th,
  first_name_en, last_name_en,
  national_id, date_of_birth, gender, phone_number, email, patient_hn
) VALUES
  -- Bangkok Hospital (4)
  ('BKH01', 'สมชาย',    'ใจดี',      'Somchai',   'Jaidee',     '1100100012341', '1985-03-15', 'male',   '0812345001', 'somchai.j@email.com',    'BKH-0001'),
  ('BKH01', 'กนกวรรณ',  'ทองดี',     'Kanokwan',  'Thongdee',   '2101200023452', '1992-07-22', 'female', '0823456002', 'kanokwan.t@email.com',   'BKH-0002'),
  ('BKH01', 'วิชัย',    'มีสุข',     'Wichai',    'Meesuk',     '3200300034563', '1978-11-05', 'male',   '0834567003', NULL,                     'BKH-0003'),
  ('BKH01', 'นิดา',     'แสงดาว',    'Nida',      'Saengdao',   '4300400045674', '2000-01-30', 'female', NULL,         'nida.s@email.com',       'BKH-0004'),

  -- Bangkok Hospital Samui (4)
  ('BKH02', 'ธีรพล',    'จันทร์เพ็ญ','Theerapol', 'Chanpen',    '1400500056785', '1990-06-18', 'male',   '0845678004', 'theerapol.c@email.com',  'BKS-0001'),
  ('BKH02', 'จิราพร',   'วงศ์ทอง',   'Jiraporn',  'Wongtong',   '2500600067896', '1988-09-12', 'female', '0856789005', 'jiraporn.w@email.com',   'BKS-0002'),
  ('BKH02', 'อนุชา',    'พงษ์ไพร',   'Anucha',    'Phongphai',  '3600700078907', '1975-04-25', 'male',   '0867890006', NULL,                     'BKS-0003'),
  ('BKH02', 'สุดา',     'ปิ่นแก้ว',  'Suda',      'Pinkaew',    '4700800089018', '1995-12-08', 'female', NULL,         'suda.p@email.com',       'BKS-0004'),

  -- Thonburi Hospital (4)
  ('TBR01', 'ชัยวัฒน์', 'บุญมี',     'Chaiwat',   'Boonmee',    '1800900090129', '1982-02-14', 'male',   '0878901007', 'chaiwat.b@email.com',    'TH-0001'),
  ('TBR01', 'รัตนา',    'ศรีสุข',    'Rattana',   'Srisuk',     '2900100001230', '1970-08-20', 'female', '0889012008', 'rattana.s@email.com',    'TH-0002'),
  ('TBR01', 'พงษ์ศักดิ์','แก้วมณี',  'Phongsak',  'Kaewmanee',  '3100200012341', '1993-05-07', 'male',   '0890123009', NULL,                     'TH-0003'),
  ('TBR01', 'พิมพ์ใจ',  'ดาวเรือง',  'Pimjai',    'Daowruang',  '4200300023452', '2002-10-15', 'female', NULL,         'pimjai.d@email.com',     'TH-0004'),

  -- Vimut Hospital (4)
  ('VIM01', 'สุรชัย',   'เพ็ชรดี',   'Surachai',  'Phetdee',    '1300400034563', '1965-07-03', 'male',   '0801234010', 'surachai.p@email.com',   'VM-0001'),
  ('VIM01', 'กัญญา',    'ลิ้มสุข',   'Kanya',     'Limsuk',     '2400500045674', '1998-03-19', 'female', '0812345011', 'kanya.l@email.com',      'VM-0002'),
  ('VIM01', 'นพดล',     'ศรีวิไล',   'Noppadon',  'Sriwilai',   '3500600056785', '1987-11-28', 'male',   '0823456012', NULL,                     'VM-0003'),
  ('VIM01', 'มานะ',     'สุขใจ',     'Mana',      'Sukjai',     '4600700067896', '1972-06-10', 'male',   NULL,         'mana.s@email.com',       'VM-0004'),

  -- Samitivej Srinakarin (4)
  ('SMJ01', 'อัญชลี',  'บุญสม',     'Anchalee',  'Boonsom',    '1700800078907', '1996-04-22', 'female', '0834567013', 'anchalee.b@email.com',   'SMJ-0001'),
  ('SMJ01', 'ธนวัฒน์', 'ใจเย็น',    'Thanawat',  'Jaiyen',     '2800900089018', '1981-09-05', 'male',   '0845678014', 'thanawat.j@email.com',   'SMJ-0002'),
  ('SMJ01', 'ประภา',    'วิเศษ',     'Prapha',    'Wiset',      '3900100090129', '1959-12-17', 'female', '0856789015', NULL,                     'SMJ-0003'),
  ('SMJ01', 'สมหญิง',  'รักดี',     'Somying',   'Rakdee',     '4100200001230', '2003-08-01', 'female', NULL,         'somying.r@email.com',    'SMJ-0004');


-- ============================================================
-- PATIENTS — Japanese (10)
-- English first+last name only (no Thai); passport_id; other fields random
-- ============================================================
INSERT INTO patients (
  hospital_id,
  first_name_en, last_name_en,
  passport_id, date_of_birth, gender, phone_number, email, patient_hn
) VALUES
  ('BKH01', 'Yuki',    'Tanaka',     'JP10234567', '1990-03-14', 'female', '+819011112001', 'yuki.tanaka@email.jp',   'BKH-0005'),
  ('BKH01', 'Hiroshi', 'Yamamoto',   'JP20345678', '1985-11-20', 'male',   '+819022223002', 'hiroshi.y@email.jp',     'BKH-0006'),
  ('BKH02', 'Akiko',   'Suzuki',     'JP30456789', '1993-07-08', 'female', '+819033334003', 'akiko.suzuki@email.jp',  'BKS-0005'),
  ('BKH02', 'Kenji',   'Watanabe',   'JP40567890', '1978-04-25', 'male',   NULL,            'kenji.w@email.jp',       'BKS-0006'),
  ('TBR01', 'Miho',    'Sato',       'JP50678901', '2001-09-15', 'female', '+819055556005', NULL,                      'TH-0005'),
  ('TBR01', 'Takeshi', 'Ito',        'JP60789012', '1969-01-30', 'male',   '+819066667006', 'takeshi.ito@email.jp',   'TH-0006'),
  ('VIM01', 'Yumi',    'Kato',       'JP70890123', '1997-06-12', 'female', NULL,            'yumi.kato@email.jp',     'VM-0005'),
  ('VIM01', 'Ryota',   'Nakamura',   'JP80901234', '1983-12-05', 'male',   '+819088889008', 'ryota.n@email.jp',       'VM-0006'),
  ('SMJ01', 'Haruki',  'Kobayashi',  'JP90012345', '1975-08-22', 'male',   '+819099990009', 'haruki.k@email.jp',      'SMJ-0005'),
  ('SMJ01', 'Sakura',  'Abe',        'JP00123456', '2000-02-14', 'female', '+819000001010', 'sakura.abe@email.jp',    'SMJ-0006');


-- ============================================================
-- PATIENTS — USA (10)
-- English first+last name only (no Thai); passport_id; 3 have middle_name_en
-- ============================================================
INSERT INTO patients (
  hospital_id,
  first_name_en, middle_name_en, last_name_en,
  passport_id, date_of_birth, gender, phone_number, email, patient_hn
) VALUES
  -- 3 with middle name
  ('BKH01', 'John',     'Michael', 'Smith',     'US10234567', '1982-05-10', 'male',   '+12125550001', 'john.m.smith@email.com',    'BKH-0007'),
  ('BKH02', 'Emily',    'Rose',    'Johnson',   'US20345678', '1994-09-23', 'female', '+13235550002', 'emily.r.johnson@email.com', 'BKS-0007'),
  ('TBR01', 'Robert',   'Lee',     'Williams',  'US30456789', '1970-12-07', 'male',   '+14045550003', 'robert.l.w@email.com',      'TH-0007'),
  -- 7 without middle name
  ('BKH01', 'James',    NULL,      'Brown',     'US40567890', '1988-03-18', 'male',   '+12015550004', 'james.brown@email.com',     'BKH-0008'),
  ('BKH02', 'Jennifer', NULL,      'Davis',     'US50678901', '1976-07-30', 'female', NULL,           'jennifer.davis@email.com',  'BKS-0008'),
  ('TBR01', 'Michael',  NULL,      'Wilson',    'US60789012', '2002-11-14', 'male',   '+15035550006', NULL,                         'TH-0008'),
  ('VIM01', 'Linda',    NULL,      'Martinez',  'US70890123', '1965-04-05', 'female', '+16025550007', 'linda.martinez@email.com',  'VM-0007'),
  ('VIM01', 'William',  NULL,      'Anderson',  'US80901234', '1991-08-19', 'male',   '+17025550008', 'william.a@email.com',       'VM-0008'),
  ('SMJ01', 'Patricia', NULL,      'Taylor',    'US90012345', '1984-01-26', 'female', '+18025550009', 'patricia.t@email.com',      'SMJ-0007'),
  ('SMJ01', 'Charles',  NULL,      'Thomas',    'US00123456', '1958-06-11', 'male',   '+19025550010', 'charles.thomas@email.com',  'SMJ-0008');


-- ============================================================
-- PATIENTS — Thai extra (20)
-- Includes name variations useful for contains-search testing
-- e.g. สมช*, กาน*, วิ*, multiple similar surnames
-- ============================================================
INSERT INTO patients (
  hospital_id,
  first_name_th, last_name_th,
  first_name_en, last_name_en,
  national_id, date_of_birth, gender, phone_number, email, patient_hn
) VALUES
  -- BKH01 — "สมช*" prefix group
  ('BKH01', 'สมชาติ',   'วงศ์ใหญ่',  'Somchat',   'Wongyai',    '5101100011111', '1980-04-10', 'male',   '0811111101', 'somchat.w@email.com',     'BKH-0009'),
  ('BKH01', 'สมชาญ',    'ดีเลิศ',    'Somchan',   'Deelert',    '5202200022222', '1975-09-25', 'male',   '0822222202', 'somchan.d@email.com',     'BKH-0010'),
  ('BKH01', 'สมหมาย',   'พุทธา',     'Sommai',    'Puttha',     '5303300033333', '1992-01-17', 'male',   NULL,         'sommai.p@email.com',      'BKH-0011'),

  -- BKH01 — "กาน*" prefix group
  ('BKH01', 'กานดา',    'ศรีทอง',    'Kanda',     'Srithong',   '5404400044444', '1988-06-03', 'female', '0844444404', 'kanda.s@email.com',       'BKH-0012'),
  ('BKH01', 'กานต์',    'พลอยดี',    'Karn',      'Ploydee',    '5505500055555', '1995-11-14', 'female', '0855555505', NULL,                      'BKH-0013'),

  -- BKH02 — "วิ*" prefix group
  ('BKH02', 'วิภา',     'ทองแดง',    'Wipa',      'Thongdaeng', '5606600066666', '1983-03-28', 'female', '0866666606', 'wipa.t@email.com',        'BKS-0009'),
  ('BKH02', 'วิไล',     'จันทร์งาม', 'Wilai',     'Channgam',   '5707700077777', '1990-08-19', 'female', NULL,         'wilai.c@email.com',       'BKS-0010'),
  ('BKH02', 'วิรัตน์',  'สุขสม',     'Wirat',     'Sooksom',    '5808800088888', '1977-12-01', 'male',   '0888888808', NULL,                      'BKS-0011'),

  -- TBR01 — same last name group (แก้ว*)
  ('TBR01', 'จินตนา',   'แก้วประเสริฐ','Jintana',  'Kaewprasert','5909900099999', '1985-05-09', 'female', '0899999909', 'jintana.k@email.com',    'TH-0009'),
  ('TBR01', 'สุพจน์',   'แก้วใส',    'Supot',     'Kaewsai',    '6100100010001', '1971-10-22', 'male',   NULL,         'supot.k@email.com',       'TH-0010'),
  ('TBR01', 'ภาวินี',   'แก้วสม',    'Pawinee',   'Kaewsom',    '6201200021112', '1999-02-18', 'female', '0801010110', 'pawinee.k@email.com',     'TH-0011'),

  -- VIM01 — mixed genders, unique names
  ('VIM01', 'ปรีดา',    'ฉัตรมงคล',  'Prida',     'Chatmongkol','6302300032223', '1968-07-31', 'male',   '0812121212', 'prida.c@email.com',       'VM-0009'),
  ('VIM01', 'อาภรณ์',   'มั่นคง',    'Arporn',    'Mankhong',   '6403400043334', '1994-04-14', 'female', NULL,         'arporn.m@email.com',      'VM-0010'),
  ('VIM01', 'ทรงศักดิ์','ฉลาดดี',   'Songsak',   'Chaladee',   '6504500054445', '1980-01-06', 'male',   '0834343434', NULL,                      'VM-0011'),

  -- SMJ01 — คน*, common short first names
  ('SMJ01', 'คนึง',     'รุ่งเรือง',  'Kanueng',   'Rungruang',  '6605600065556', '1963-09-09', 'female', '0845454545', 'kanueng.r@email.com',     'SMJ-0009'),
  ('SMJ01', 'คนอง',     'จิตต์ดี',   'Kanong',    'Jitdee',     '6706700076667', '1977-03-21', 'male',   NULL,         'kanong.j@email.com',      'SMJ-0010'),
  ('SMJ01', 'ปิยะ',     'ใจงาม',     'Piya',      'Jaingam',    '6807800087778', '2001-07-15', 'male',   '0867676767', NULL,                      'SMJ-0011'),

  -- Both hospitals — have BOTH national_id AND passport_id (dual-citizen)
  ('BKH01', 'นาวา',     'สมบูรณ์',   'Nawa',      'Somboon',    '6908900098889', '1986-05-30', 'male',   '0878787878', 'nawa.s@email.com',        'BKH-0014'),
  ('TBR01', 'ดวงใจ',    'เพชรงาม',   'Duangjai',  'Petchngam',  '7100100009990', '1973-11-03', 'female', '0889898989', 'duangjai.p@email.com',    'TH-0012'),
  ('VIM01', 'อุไร',     'ทองสุข',    'Urai',      'Thongsuk',   '7201200010001', '1958-08-27', 'female', NULL,         'urai.t@email.com',        'VM-0012');

-- Add passport to dual-citizen rows above
INSERT INTO patients (
  hospital_id,
  first_name_th, last_name_th,
  first_name_en, last_name_en,
  national_id, passport_id, date_of_birth, gender, phone_number, email, patient_hn
) VALUES
  ('BKH02', 'ณัฐพล',   'วรรณโชติ',  'Nattaphon', 'Wannachot',  '7302300021112', 'TH10000001', '1991-06-07', 'male',   '0890909090', 'nattaphon.w@email.com',  'BKS-0012'),
  ('SMJ01', 'ชลธิชา',  'บุญรักษ์',  'Chonthicha','Boonrak',    '7403400032223', 'TH20000002', '1997-10-19', 'female', '0801020304', 'chonthicha.b@email.com', 'SMJ-0012');


-- ============================================================
-- PATIENTS — Japanese extra (10)
-- Includes "Hiro*" name group for contains-search testing
-- ============================================================
INSERT INTO patients (
  hospital_id,
  first_name_en, last_name_en,
  passport_id, date_of_birth, gender, phone_number, email, patient_hn
) VALUES
  -- "Hiro*" group — searching "Hiro" returns all 4
  ('BKH01', 'Hiroshi',  'Nakagawa',   'JP11111111', '1988-04-12', 'male',   '+819011001101', 'hiroshi.n@email.jp',    'BKH-0015'),
  ('BKH01', 'Hiroki',   'Fujimoto',   'JP22222222', '1993-07-25', 'male',   '+819022002202', 'hiroki.f@email.jp',     'BKH-0016'),
  ('BKH02', 'Hiroyuki', 'Ogawa',      'JP33333333', '1979-11-08', 'male',   '+819033003303', 'hiroyuki.o@email.jp',   'BKS-0013'),
  ('BKH02', 'Hiromi',   'Matsuda',    'JP44444444', '1996-02-14', 'female', NULL,            'hiromi.m@email.jp',     'BKS-0014'),

  -- Other Japanese names
  ('TBR01', 'Naomi',    'Hayashi',    'JP55555555', '1984-09-30', 'female', '+819055005505', 'naomi.h@email.jp',      'TH-0013'),
  ('TBR01', 'Keiko',    'Shimizu',    'JP66666666', '1971-03-17', 'female', '+819066006606', 'keiko.s@email.jp',      'TH-0014'),
  ('VIM01', 'Taro',     'Yamada',     'JP77777777', '2000-08-05', 'male',   NULL,            'taro.yamada@email.jp',  'VM-0013'),
  ('VIM01', 'Hanako',   'Inoue',      'JP88888888', '1987-12-21', 'female', '+819088008808', 'hanako.i@email.jp',     'VM-0014'),
  ('SMJ01', 'Kazuo',    'Saito',      'JP99999999', '1965-06-10', 'male',   '+819099009909', 'kazuo.s@email.jp',      'SMJ-0013'),
  ('SMJ01', 'Yoko',     'Tanaka',     'JP00000001', '1992-01-28', 'female', '+819000000010', 'yoko.tanaka@email.jp',  'SMJ-0014');


-- ============================================================
-- PATIENTS — European & Korean (15)
-- Smith* surname group + Korean names
-- ============================================================
INSERT INTO patients (
  hospital_id,
  first_name_en, middle_name_en, last_name_en,
  passport_id, date_of_birth, gender, phone_number, email, patient_hn
) VALUES
  -- "Smith*" surname group — searching "Smith" returns all 4
  ('BKH01', 'David',    NULL,       'Smith',      'GB10000001', '1980-06-15', 'male',   '+441234560001', 'david.smith@email.uk',      'BKH-0017'),
  ('BKH02', 'Sarah',    'Jane',     'Smith',      'GB20000002', '1992-03-08', 'female', '+441234560002', 'sarah.j.smith@email.uk',    'BKS-0015'),
  ('TBR01', 'Thomas',   NULL,       'Smithson',   'GB30000003', '1975-11-22', 'male',   NULL,            'thomas.smithson@email.uk',  'TH-0015'),
  ('VIM01', 'Oliver',   'James',    'Blacksmith', 'GB40000004', '1988-09-14', 'male',   '+441234560004', 'oliver.blacksmith@email.uk','VM-0015'),

  -- Other European
  ('BKH01', 'Marie',    'Claire',   'Dupont',     'FR10000001', '1985-07-04', 'female', '+33612340001',  'marie.dupont@email.fr',     'BKH-0018'),
  ('BKH02', 'Hans',     NULL,       'Müller',     'DE10000001', '1970-02-19', 'male',   '+49170340001',  'hans.muller@email.de',      'BKS-0016'),
  ('TBR01', 'Sofia',    NULL,       'Rossi',      'IT10000001', '1995-10-11', 'female', '+39334560001',  'sofia.rossi@email.it',      'TH-0016'),
  ('VIM01', 'Lucas',    'André',    'Bernard',    'FR20000002', '1982-04-27', 'male',   NULL,            'lucas.bernard@email.fr',    'VM-0016'),
  ('SMJ01', 'Emma',     NULL,       'Wilson',     'AU10000001', '1998-08-03', 'female', '+61412340001',  'emma.wilson@email.au',      'SMJ-0015'),
  ('SMJ01', 'Liam',     'Patrick',  'O''Brien',   'IE10000001', '1979-12-16', 'male',   '+35387340001',  'liam.obrien@email.ie',      'SMJ-0016'),

  -- Korean
  ('BKH01', 'Minjun',   NULL,       'Kim',        'KR10000001', '1991-05-20', 'male',   '+821012340001', 'minjun.kim@email.kr',       'BKH-0019'),
  ('BKH02', 'Jisoo',    NULL,       'Park',       'KR20000002', '1997-09-17', 'female', '+821012340002', 'jisoo.park@email.kr',       'BKS-0017'),
  ('TBR01', 'Hyunwoo',  NULL,       'Lee',        'KR30000003', '1986-03-05', 'male',   NULL,            'hyunwoo.lee@email.kr',      'TH-0017'),
  ('VIM01', 'Sooyeon',  NULL,       'Choi',       'KR40000004', '2001-07-29', 'female', '+821012340004', 'sooyeon.choi@email.kr',     'VM-0017'),
  ('SMJ01', 'Junho',    NULL,       'Jung',       'KR50000005', '1974-01-13', 'male',   '+821012340005', 'junho.jung@email.kr',       'SMJ-0017');


-- ============================================================
-- STAFF (1) — seed admin account  ← must come after hospitals (FK)
-- password: Admin1!xx  (bcrypt cost 10)
-- ============================================================
INSERT INTO staff (hospital_id, email, password, created_at, updated_at) VALUES
  ('BKH01', 'admin@bangkokhospital.com', '$2a$10$8uM7WKMoblegRY4WCjaa0OUuaA35NYNcBlkatHZUsK5/bXk6LS1zC', NOW(), NOW());
