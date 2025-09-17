USE ednevnik_workspace;

INSERT IGNORE INTO cantons (canton_code, canton_name, country)
VALUES
('ZDK', 'Zeničko-dobojski', 'BiH'),
('KS', 'Kanton Sarajevo', 'BiH'),
('TK', 'Tuzlanski', 'BiH'),
('USK', 'Unsko-Sanski', 'BiH'),
('BPK', 'Bosansko-podrinjski', 'BiH'),
('SBK', 'Srednjobosanski', 'BiH'),
('HNK', 'Hercegovačko-neretvanski', 'BiH'),
('ZHK', 'Zapadnohercegovački', 'BiH'),
('PK', 'Posavski', 'BiH'),
('K10', 'Kanton 10', 'BiH');

INSERT INTO classes (class_code)
VALUES
('I'),
('II'),
('III'),
('IV'),
('V'),
('VI'),
('VII'),
('VIII'),
('IX');

-- Variables to track class ids
SET @class_1 = (SELECT id FROM classes WHERE class_code = 'I');
SET @class_2 = (SELECT id FROM classes WHERE class_code = 'II');
SET @class_3 = (SELECT id FROM classes WHERE class_code = 'III');
SET @class_4 = (SELECT id FROM classes WHERE class_code = 'IV');
SET @class_5 = (SELECT id FROM classes WHERE class_code = 'V');
SET @class_6 = (SELECT id FROM classes WHERE class_code = 'VI');
SET @class_7 = (SELECT id FROM classes WHERE class_code = 'VII');
SET @class_8 = (SELECT id FROM classes WHERE class_code = 'VIII');
SET @class_9 = (SELECT id FROM classes WHERE class_code = 'IX');

INSERT INTO subjects (subject_name, subject_code)
VALUES
('Biologija', 'BIO'),
('Bosanski jezik i književnost', 'BJZ'),
('Hrvatski jezik i književnost', 'HJZ'),
('Srpski jezik i književnost', 'SJZ'),
('Demokratija i ljudska prava', 'DMK'),
('Njemački jezik', 'NJE'),
('Turski jezik', 'TUR'),
('Francuski jezik', 'FRA'),
('Arapski jezik', 'ARA'),
('Engleski jezik', 'ENG'),
('Fizika', 'FIZ'),
('Geografija', 'GEO'),
('Građansko obrazovanje', 'GOB'),
('Hemija', 'HEM'),
('Historija', 'HIS'),
('Informatika', 'INF'),
('Kultura življenja', 'KZ'),
('Likovna kultura', 'LK'),
('Muzička kultura', 'MK'),
('Matematika', 'MM'),
('Tjelesni i zdravstveni odgoj', 'TIZO'),
('Tehnička kultura', 'TK'),
('Vjeronauka', 'VJR');

-- Variables to track subject ids
SET @Biologija = (SELECT id FROM subjects WHERE subject_code = 'BIO');
SET @Bosanski = (SELECT id FROM subjects WHERE subject_code = 'BJZ');
SET @Hrvatski = (SELECT id FROM subjects WHERE subject_code = 'HJZ');
SET @Srpski = (SELECT id FROM subjects WHERE subject_code = 'SJZ');
SET @Demokratija = (SELECT id FROM subjects WHERE subject_code = 'DMK');
SET @Njemacki = (SELECT id FROM subjects WHERE subject_code = 'NJE');
SET @Turski = (SELECT id FROM subjects WHERE subject_code = 'TUR');
SET @Francuski = (SELECT id FROM subjects WHERE subject_code = 'FRA');
SET @Arapski = (SELECT id FROM subjects WHERE subject_code = 'ARA');
SET @Engleski = (SELECT id FROM subjects WHERE subject_code = 'ENG');
SET @Fizika = (SELECT id FROM subjects WHERE subject_code = 'FIZ');
SET @Geografija = (SELECT id FROM subjects WHERE subject_code = 'GEO');
SET @Gradansko = (SELECT id FROM subjects WHERE subject_code = 'GOB');
SET @Hemija = (SELECT id FROM subjects WHERE subject_code = 'HEM');
SET @Historija = (SELECT id FROM subjects WHERE subject_code = 'HIS');
SET @Informatika = (SELECT id FROM subjects WHERE subject_code = 'INF');
SET @Kultura = (SELECT id FROM subjects WHERE subject_code = 'KZ');
SET @Likovna = (SELECT id FROM subjects WHERE subject_code = 'LK');
SET @Muzicka = (SELECT id FROM subjects WHERE subject_code = 'MK');
SET @Matematika = (SELECT id FROM subjects WHERE subject_code = 'MM');
SET @Tjelesni = (SELECT id FROM subjects WHERE subject_code = 'TIZO');
SET @Tehnicka = (SELECT id FROM subjects WHERE subject_code = 'TK');
SET @Vjeronauka = (SELECT id FROM subjects WHERE subject_code = 'VJR');

INSERT INTO npp (npp_name, npp_code)
VALUES
('Bosanski jezik', 'BJZ'),
('Hrvatski jezik', 'HJZ'),
('Srpski jezik', 'SJZ'),
('EU-VET', 'EU-VET');

-- Variables to track NPP ids
SET @npp_bosanski = (SELECT id FROM npp WHERE npp_code = 'BJZ');
SET @npp_hrvatski = (SELECT id FROM npp WHERE npp_code = 'HJZ');
SET @npp_srpski = (SELECT id FROM npp WHERE npp_code = 'SJZ');
SET @npp_eu_vet = (SELECT id FROM npp WHERE npp_code = 'EU-VET');

INSERT INTO courses_secondary (course_code, course_name, course_duration)
VALUES
('AT', 'Arhitektonski tehničar', 'IV'),
('ELEL', 'Elektrotehničar elektronike', 'IV'),
('ENERG', 'Elektrotehničar energetike', 'IV'),
('GEO', 'Geodetski tehničar (geometar)', 'IV'),
('MASTK', 'Mašinski tehničar za kompjutersko upravljanje mašinama', 'IV'),
('METT', 'Metalurški tehničar', 'IV'),
('RTIA', 'Elektrotehničar računarske tehnike i automatike', 'IV'),
('ST', 'Tehničar drumskog saobraćaja', 'IV');

-- Variables to track course ids
SET @AT = (SELECT id FROM courses_secondary WHERE course_code = 'AT');
SET @ELEL = (SELECT id FROM courses_secondary WHERE course_code = 'ELEL');
SET @ENERG = (SELECT id FROM courses_secondary WHERE course_code = 'ENERG');
SET @GEO = (SELECT id FROM courses_secondary WHERE course_code = 'GEO');
SET @MASTK = (SELECT id FROM courses_secondary WHERE course_code = 'MASTK');
SET @METT = (SELECT id FROM courses_secondary WHERE course_code = 'METT');
SET @RTIA = (SELECT id FROM courses_secondary WHERE course_code = 'RTIA');
SET @ST = (SELECT id FROM courses_secondary WHERE course_code = 'ST');

-- Kurikulumi za osnovnu školu po NPPovima
INSERT INTO curriculum (class_id, npp, canton_code, tenant_type)
VALUES
-- Bosanski npp
(@class_1, @npp_bosanski, 'ZDK', 'primary'),
(@class_2, @npp_bosanski, 'ZDK', 'primary'),
(@class_3, @npp_bosanski, 'ZDK', 'primary'),
(@class_4, @npp_bosanski, 'ZDK', 'primary'),
(@class_5, @npp_bosanski, 'ZDK', 'primary'),
(@class_6, @npp_bosanski, 'ZDK', 'primary'),
(@class_7, @npp_bosanski, 'ZDK', 'primary'),
(@class_8, @npp_bosanski, 'ZDK', 'primary'),
(@class_9, @npp_bosanski, 'ZDK', 'primary');
-- Hrvatski npp
(@class_1, @npp_hrvatski, 'ZDK', 'primary'),
(@class_2, @npp_hrvatski, 'ZDK', 'primary'),
(@class_3, @npp_hrvatski, 'ZDK', 'primary'),
(@class_4, @npp_hrvatski, 'ZDK', 'primary'),
(@class_5, @npp_hrvatski, 'ZDK', 'primary'),
(@class_6, @npp_hrvatski, 'ZDK', 'primary'),
(@class_7, @npp_hrvatski, 'ZDK', 'primary'),
(@class_8, @npp_hrvatski, 'ZDK', 'primary'),
(@class_9, @npp_hrvatski, 'ZDK', 'primary');

-- Kurikulumi za srednju školu smjer RTIA (npp bosanski jezik)
INSERT INTO curriculum (class_id, npp, course_id, canton_code, tenant_type)
VALUES
(@class_1, @npp_bosanski, @RTIA, 'ZDK', 'secondary'),
(@class_2, @npp_bosanski, @RTIA, 'ZDK', 'secondary'),
(@class_3, @npp_bosanski, @RTIA, 'ZDK', 'secondary'),
(@class_4, @npp_bosanski, @RTIA, 'ZDK', 'secondary'),

-- Variables to track curriculum ids
SET @curriculum_class_1_primary_school_BOS_ZDK = (
    SELECT id FROM curriculum WHERE class_id = @class_1 AND npp = @npp_bosanski AND tenant_type = 'primary'
);
SET @curriculum_class_2_primary_school_BOS_ZDK = (
    SELECT id FROM curriculum WHERE class_id = @class_2 AND npp = @npp_bosanski AND tenant_type = 'primary'
);
SET @curriculum_class_3_primary_school_BOS_ZDK = (
    SELECT id FROM curriculum WHERE class_id = @class_3 AND npp = @npp_bosanski AND tenant_type = 'primary'
);
SET @curriculum_class_4_primary_school_BOS_ZDK = (
    SELECT id FROM curriculum WHERE class_id = @class_4 AND npp = @npp_bosanski AND tenant_type = 'primary'
);
SET @curriculum_class_5_primary_school_BOS_ZDK = (
    SELECT id FROM curriculum WHERE class_id = @class_5 AND npp = @npp_bosanski AND tenant_type = 'primary'
);
SET @curriculum_class_6_primary_school_BOS_ZDK = (
    SELECT id FROM curriculum WHERE class_id = @class_6 AND npp = @npp_bosanski AND tenant_type = 'primary'
);
SET @curriculum_class_7_primary_school_BOS_ZDK = (
    SELECT id FROM curriculum WHERE class_id = @class_7 AND npp = @npp_bosanski AND tenant_type = 'primary'
);
SET @curriculum_class_8_primary_school_BOS_ZDK = (
    SELECT id FROM curriculum WHERE class_id = @class_8 AND npp = @npp_bosanski AND tenant_type = 'primary'
);
SET @curriculum_class_9_primary_school_BOS_ZDK = (
    SELECT id FROM curriculum WHERE class_id = @class_9 AND npp = @npp_bosanski AND tenant_type = 'primary'
);
