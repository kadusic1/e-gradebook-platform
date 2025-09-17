-- Use python create_workspace_db.py to run this script
-- With: python create_workspace_db.py

-- Izbrisi wp bazu ako postoji
DROP DATABASE IF EXISTS ednevnik_workspace;

-- Kreira samo workspace bazu i tabele prema db_design.md

-- 1. WORKSPACE DB
-- Kreiranje nove workspace baze
SELECT '[LOG] Starting workspace database creation...' AS info;

CREATE DATABASE ednevnik_workspace CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
SELECT '[LOG] Created database ednevnik_workspace.' AS info;

-- Prebacivanje na workspace bazu
USE ednevnik_workspace;
SELECT '[LOG] Using workspace database.' AS info;

DELIMITER //
CREATE PROCEDURE DropLocalUsers()
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE user_name VARCHAR(255);
    DECLARE cur CURSOR FOR
        SELECT user
        FROM mysql.user
        WHERE host = 'localhost'
        AND user NOT IN ('root', 'mariadb.sys');
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO user_name;
        IF done THEN
            LEAVE read_loop;
        END IF;
        SET @sql = CONCAT('DROP USER IF EXISTS ''', user_name, '''@''localhost'';');
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END LOOP;
    CLOSE cur;
END//
DELIMITER ;

SELECT '[LOG] Dropping local users except root and mariadb.sys...' AS info;
CALL DropLocalUsers();
SELECT '[LOG] Local users dropped.' AS info;
DROP PROCEDURE DropLocalUsers;

-- Kreiranje tabele kantona
SELECT '[LOG] Creating core tables...' AS info;
CREATE TABLE cantons (
    canton_code VARCHAR(10) PRIMARY KEY,
    canton_name VARCHAR(50) NOT NULL,
    country VARCHAR(30) DEFAULT 'BiH',
    minister_name VARCHAR(100),
    ministry_name VARCHAR(100),
    ministry_phone VARCHAR(20)
);
-- PRIMJER: { 'canton_code': '01', 'canton_name': 'Sarajevski', 'country': 'BiH' }

-- Kreiranje tabele razreda
CREATE TABLE classes (
    class_code VARCHAR(10) PRIMARY KEY
);
-- PRIMJER: { 'class_code': 'VI' }

-- Kreiranje tabele predmeta
CREATE TABLE subjects (
    subject_code VARCHAR(15) PRIMARY KEY,
    subject_name VARCHAR(100) NOT NULL
);
-- PRIMJER: { 'subject_code': 'MM', 'subject_name': 'Matematika' }
CREATE INDEX idx_subject_name ON subjects (subject_name);

CREATE TABLE accounts (
    id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_by_teacher_id INT, -- NULL if the account is created via registration
    account_type ENUM('root', 'tenant_admin', 'teacher', 'pupil', 'parent') NOT NULL, -- staf is for teachers, admin, etc.
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Kreiranje tabele teachers
CREATE TABLE teachers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    phone VARCHAR(20) UNIQUE,
    contractions VARCHAR(50),
    title VARCHAR(50),
    account_id INT,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
-- PRIMJER: { 'id': 1, 'name': 'Ivana', 'last_name': 'Ivić', 'email': 'ivana.ivic@skole.ba', 'password': 'hash', 'phone': '061333333', 'tenant_id': 1, 'role': 'teacher' }

CREATE TABLE tenant (
    id INT PRIMARY KEY AUTO_INCREMENT,
    tenant_name VARCHAR(200) NOT NULL,
    tenant_city VARCHAR(100),
    tenant_type ENUM('primary', 'secondary') NOT NULL,
    canton_code VARCHAR(10),
    address VARCHAR(200),
    phone VARCHAR(20) UNIQUE,
    email VARCHAR(100),
    director_name VARCHAR(100),
    longitude DECIMAL(9, 6),
    latitude DECIMAL(9, 6),
    AI_enabled BOOLEAN DEFAULT FALSE,
    tenant_admin_id INT,
    domain VARCHAR(100) UNIQUE,
    color_config ENUM ('0', '1', '2', '3'),
    teacher_display ENUM ('card', 'table') DEFAULT 'card',
    teacher_invite_display ENUM ('card', 'table') DEFAULT 'card',
    pupil_display ENUM ('card', 'table') DEFAULT 'card',
    pupil_invite_display ENUM ('card', 'table') DEFAULT 'card',
    section_display ENUM ('card', 'table') DEFAULT 'card',
    curriculum_display ENUM ('card', 'table') DEFAULT 'card',
    semester_display ENUM ('card', 'table') DEFAULT 'card',
    lesson_display ENUM ('card', 'table') DEFAULT 'card',
    absence_display ENUM ('card', 'table') DEFAULT 'card',
    classroom_display ENUM ('card', 'table') DEFAULT 'card',
    specialization ENUM ('regular', 'religion', 'musical') DEFAULT 'regular',
    FOREIGN KEY (canton_code) REFERENCES cantons(canton_code),
    FOREIGN KEY (tenant_admin_id) REFERENCES teachers(id) ON DELETE CASCADE
);
-- PRIMJER: { 'id': 1, 'tenant_name': 'OŠ Skender Kulenović', 'tenant_type': 'primary', 'canton_code': '01', 'address': 'Adresa 1', 'phone': '033123456', 'email': 'os.skender@skole.ba', 'director_name': 'Neko Nekić' }

CREATE INDEX idx_tenant_email ON tenant (email);
CREATE INDEX idx_tenant_domain ON tenant (domain);

ALTER TABLE accounts
    ADD CONSTRAINT fk_created_by_teacher
    FOREIGN KEY (created_by_teacher_id) REFERENCES teachers(id) ON DELETE SET NULL;

CREATE TABLE teacher_tenant (
    teacher_id INT,
    tenant_id INT,
    PRIMARY KEY(teacher_id, tenant_id),
    FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
);

-- Kreiranje tabele nastavnih planova i programa (NPP)
CREATE TABLE npp (
    npp_code VARCHAR(20) PRIMARY KEY,
    npp_name VARCHAR(100) NOT NULL
);

-- Kreiranje tabele kurseva
CREATE TABLE courses_secondary (
    course_code VARCHAR(20) PRIMARY KEY,
    course_name VARCHAR(200) NOT NULL,
    course_duration ENUM('III', 'IV'),
    course_icon VARCHAR(30),
    course_type ENUM(
        'regular',
        'language',
        'natural_science',
        'it',
        'math',
        'mechanical',
        'electrical',
        'construction',
        'chemical',
        'agricultural',
        'art',
        'musical',
        'sports',
        'woodworking',
        'diet',
        'geodetic',
        'transport',
        'textile',
        'dental',
        'health',
        'vet',
        'ecological',
        'skin',
        'cater',
        'turism',
        'seller',
        'economic',
        'service',
        'other'
    ),
    is_defitiant BOOLEAN,
    has_scholarship BOOLEAN
);

-- Kreiranje tabele kurikuluma
CREATE TABLE curriculum (
    curriculum_code VARCHAR(30) PRIMARY KEY,
    curriculum_name VARCHAR(250) NOT NULL,
    class_code VARCHAR(10) NOT NULL,
    npp_code VARCHAR(20) NOT NULL,
    course_code VARCHAR(20) NULL, -- NULL ako nije srednja škola
    canton_code VARCHAR(10) NOT NULL,
    tenant_type ENUM('primary', 'secondary') NOT NULL,
    final_curriculum BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (class_code) REFERENCES classes(class_code),
    FOREIGN KEY (course_code) REFERENCES courses_secondary(course_code),
    FOREIGN KEY (npp_code) REFERENCES npp(npp_code),
    FOREIGN KEY (canton_code) REFERENCES cantons(canton_code)
);
-- PRIMJERI (Osnovna škola):
-- { 'curriculum_code': 'bos_primary_6', 'curriculum_name': 'Test', 'class_code': 'VI', 'npp': 'BOS', 'course_id': NULL, 'canton_code': 'ZDK' }
-- { 'curriculum_code': bos_primary_7, 'curriculum_name': 'Test', 'class_code': 'VII', 'npp': 'BOS', 'course_id': NULL, 'canton_code': 'ZDK' }
-- PRIMJERI (Srednja škola):
-- { 'curriculum_code': bos_secondary_RTIA_1, 'curriculum_name': 'Test', 'class_code': 'I', 'npp': 'BOS', 'course_id': 'RTIA', 'canton_code': 'ZDK' }
-- { 'curriculum_code': bos_secondary_RTIA_2, 'curriculum_name': 'Test', 'class_code': 'II', 'npp': 'BOS', 'course_id': 'RTIA', 'canton_code': 'ZDK' }

-- Kreiranje tabele predmeta po kurikulumu
CREATE TABLE curriculum_subjects (
    curriculum_code VARCHAR(30),
    subject_code VARCHAR(15),
    PRIMARY KEY(curriculum_code, subject_code),
    FOREIGN KEY (curriculum_code) REFERENCES curriculum(curriculum_code),
    FOREIGN KEY (subject_code) REFERENCES subjects(subject_code)
);

-- PRIMJERI:
-- { 'curriculum_code': 'bos_primary_6', 'subject_code': 'MM' }
-- { 'curriculum_code': 'bos_primary_6', 'subject_code': 'BJZ' }
-- { 'curriculum_code': 'bos_primary_7', 'subject_code': 'MM' }

-- Tabela semestar - sve moguće vrste semestara
CREATE TABLE semester(
    semester_code VARCHAR(10) PRIMARY KEY,
    semester_name VARCHAR(30),
    -- Field that indicates the order of the semester in curriculum
    -- For example primary school progress level for bosanski NPP can be either
    -- 1 OR 2
    progress_level INT
);

CREATE TABLE npp_semester (
    npp_code VARCHAR(20),
    semester_code VARCHAR(10),
    start_date DATE,
    end_date DATE,
    PRIMARY KEY(npp_code, semester_code),
    FOREIGN KEY (npp_code) REFERENCES npp(npp_code),
    FOREIGN KEY (semester_code) REFERENCES semester(semester_code)
);

CREATE TABLE tenant_semester (
    tenant_id INT,
    semester_code VARCHAR(10),
    start_date DATE,
    end_date DATE,
    npp_code VARCHAR(20),
    PRIMARY KEY(tenant_id, semester_code, npp_code),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (semester_code) REFERENCES semester(semester_code),
    FOREIGN KEY (npp_code) REFERENCES npp(npp_code)
);

CREATE TABLE curriculum_tenant (
    tenant_id INT,
    curriculum_code VARCHAR(30),
    PRIMARY KEY(tenant_id, curriculum_code),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (curriculum_code) REFERENCES curriculum(curriculum_code)
);

CREATE TABLE pupil_global (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    jmbg VARCHAR(13) UNIQUE,
    gender ENUM('M', 'F'),
    address VARCHAR(200),
    guardian_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) UNIQUE,
    guardian_number VARCHAR(20),
    date_of_birth DATE,
    religion ENUM ('Islam', 'Catholic', 'Orthodox', 'Jewish', 'Other', 'NotAttendingReligion'),
    child_of_martyr BOOLEAN,
    father_name VARCHAR(50),
    mother_name VARCHAR(50),
    parents_rvi BOOLEAN,
    living_condition ENUM('both_parents', 'one_parent', 'another_family_or_alone',
    'institution_for_children_without_parents'),
    student_dorm BOOLEAN,
    refugee BOOLEAN,
    returnee_from_abroad BOOLEAN,
    place_of_birth VARCHAR(100) NOT NULL,
    country_of_birth VARCHAR(100),
    country_of_living VARCHAR(100),
    citizenship VARCHAR(100),
    ethnicity VARCHAR(100),
    father_occupation ENUM('PhD', 'MR', 'VSS', 'VŠS', 'SSS', 'KV', 'OS', 'NoOccupation'),
    mother_occupation ENUM('PhD', 'MR', 'VSS', 'VŠS', 'SSS', 'KV', 'OS', 'NoOccupation'),
    has_no_parents BOOLEAN,
    extra_information VARCHAR(250),
    child_alone VARCHAR(250),
    is_commuter BOOLEAN,
    commuting_type ENUM('Walking', 'Bike', 'Car', 'Bus', 'Train', 'NotTraveling'),
    distance_to_school_km ENUM('<=5km', '5km - 10km', '10km - 25km', '>25km'),
    has_hifz BOOLEAN,
    special_honors BOOLEAN,
    account_id INT,
    parent_access_code VARCHAR(36) NOT NULL DEFAULT (UUID()),
    -- The eupis_link_id is used to associate the pupil with their EUPIS account
    -- It is NULL if the pupil is not linked to an EUPIS account
    eupis_link_id INT,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
CREATE INDEX idx_pupil_global_last_name_name ON pupil_global (last_name, name);

CREATE TABLE pupil_tenant (
    pupil_id INT,
    tenant_id INT,
    available_for_enrollment BOOLEAN DEFAULT FALSE,
    transferred_to_enrollment BOOLEAN DEFAULT FALSE,
    PRIMARY KEY(pupil_id, tenant_id),
    FOREIGN KEY (pupil_id) REFERENCES pupil_global(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
);

-- Pending Accounts Table with UUID verification token
CREATE TABLE pending_accounts (
    id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    account_type ENUM('root', 'tenant_admin', 'teacher', 'pupil', 'parent') NOT NULL,
    verification_token CHAR(36) NOT NULL UNIQUE DEFAULT (UUID()),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL 1 DAY)
);

-- Pending Teachers Table
CREATE TABLE pending_teachers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    phone VARCHAR(20) UNIQUE,
    contractions VARCHAR(50),
    title VARCHAR(50),
    account_id INT NOT NULL,
    FOREIGN KEY (account_id) REFERENCES pending_accounts(id) ON DELETE CASCADE
);

-- Pending Pupils Table
CREATE TABLE pending_pupil_global (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    jmbg VARCHAR(13) UNIQUE,
    gender ENUM('M', 'F'),
    address VARCHAR(200),
    guardian_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) UNIQUE,
    guardian_number VARCHAR(20),
    date_of_birth DATE,
    religion ENUM('Islam', 'Catholic', 'Orthodox', 'Jewish', 'Other', 'NotAttendingReligion'),
    account_id INT NOT NULL,
    FOREIGN KEY (account_id) REFERENCES pending_accounts(id) ON DELETE CASCADE
);

CREATE TABLE invite_index (
    id INT PRIMARY KEY AUTO_INCREMENT,
    invite_id INT NOT NULL,
    account_id INT NOT NULL,
    tenant_id INT NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
);

CREATE INDEX idx_global_invite_id_account_tenant ON invite_index (invite_id, account_id, tenant_id);

CREATE TABLE global_domains (
    domain VARCHAR(100) PRIMARY KEY
);

CREATE TABLE primary_school_final_grades (
    pupil_id INT,
    tenant_id INT,
    subject_code VARCHAR(15),
    class_code VARCHAR(10),
    grade INT CHECK (grade BETWEEN 1 AND 5),
    school_specialization ENUM ('regular', 'religion', 'musical'),
    PRIMARY KEY(pupil_id, class_code, school_specialization, tenant_id, subject_code),
    FOREIGN KEY (pupil_id) REFERENCES pupil_global(id) ON DELETE CASCADE,
    FOREIGN KEY (subject_code) REFERENCES subjects(subject_code),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (class_code) REFERENCES classes(class_code)
);

CREATE TABLE primary_school_behaviour_grades (
    pupil_id INT,
    tenant_id INT,
    class_code VARCHAR(10),
    behaviour ENUM('primjerno', 'vrlodobro', 'dobro', 'zadovoljavajuće', 'loše') NOT NULL DEFAULT 'primjerno',
    school_specialization ENUM ('regular', 'religion', 'musical'),
    PRIMARY KEY(pupil_id, class_code, school_specialization, tenant_id),
    FOREIGN KEY (pupil_id) REFERENCES pupil_global(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (class_code) REFERENCES classes(class_code)
);

CREATE TABLE high_school_final_grades (
    pupil_id INT,
    tenant_id INT,
    subject_code VARCHAR(15),
    class_code VARCHAR(10),
    grade INT CHECK (grade BETWEEN 1 AND 5),
    school_specialization ENUM ('regular', 'religion', 'musical'),
    PRIMARY KEY(pupil_id, class_code, school_specialization, tenant_id, subject_code),
    FOREIGN KEY (pupil_id) REFERENCES pupil_global(id) ON DELETE CASCADE,
    FOREIGN KEY (subject_code) REFERENCES subjects(subject_code),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (class_code) REFERENCES classes(class_code)
);

CREATE TABLE high_school_behaviour_grades (
    pupil_id INT,
    tenant_id INT,
    class_code VARCHAR(10),
    behaviour ENUM('primjerno', 'vrlodobro', 'dobro', 'zadovoljavajuće', 'loše') NOT NULL DEFAULT 'primjerno',
    school_specialization ENUM ('regular', 'religion', 'musical'),
    PRIMARY KEY(pupil_id, class_code, school_specialization, tenant_id),
    FOREIGN KEY (pupil_id) REFERENCES pupil_global(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (class_code) REFERENCES classes(class_code)
);

CREATE TABLE embeddings (
    id INT PRIMARY KEY AUTO_INCREMENT,
    metadata JSON,
    content TEXT,
    -- collection_id is used to group embeddings, it is needed for langchain to work
    -- but we wont use it for anything specific, so we set it to a default value
    collection_id CHAR(36) NOT NULL DEFAULT '00000000-0000-0000-0000-000000000001',
    -- Vector size of LaBSE embeddings is 768
    embedding VECTOR(768) NOT NULL,
    VECTOR INDEX (embedding) M=6 DISTANCE=cosine
);

-- Event to automatically delete expired pending accounts every hour
DELIMITER $$
CREATE EVENT cleanup_expired_pending_accounts
ON SCHEDULE EVERY 1 HOUR
DO
BEGIN
    DELETE FROM pending_accounts WHERE expires_at < NOW();
END$$
DELIMITER ;

-- Enable event scheduler if not already enabled
SET GLOBAL event_scheduler = ON;


SELECT '[LOG] Creating users and granting privileges...' AS info;

SELECT '[LOG] Dropping user eacon if exists...' AS info;
DROP USER IF EXISTS 'eacon'@'localhost';

SELECT '[LOG] Creating user eacon...' AS info;
CREATE USER 'eacon'@'localhost' IDENTIFIED VIA mysql_native_password
USING PASSWORD('test1234');

SELECT '[LOG] Granting all privileges to user eacon...' AS info;
GRANT ALL PRIVILEGES ON *.* TO 'eacon'@'localhost' WITH GRANT OPTION;
GRANT CREATE USER ON *.* TO 'eacon'@'localhost' WITH GRANT OPTION;
GRANT RELOAD ON *.* TO 'eacon'@'localhost' WITH GRANT OPTION;


SELECT '[LOG] Dropping user tenant_admin if exists...' AS info;
DROP USER IF EXISTS 'tenant_admin'@'localhost';

SELECT '[LOG] Creating user tenant_admin...' AS info;
CREATE USER 'tenant_admin'@'localhost';

SELECT '[LOG] Granting tenant admin privileges...' AS info;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.accounts TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.teachers TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.teacher_tenant TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.tenant_semester TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.curriculum_tenant TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.pupil_global TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.pupil_tenant TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.pending_accounts TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.pending_teachers TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.pending_pupil_global TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.invite_index TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.primary_school_final_grades TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.primary_school_behaviour_grades TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.high_school_final_grades TO 'tenant_admin'@'localhost' WITH GRANT OPTION;
GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.high_school_behaviour_grades TO 'tenant_admin'@'localhost' WITH GRANT OPTION;


SELECT '[LOG] Dropping user teacher if exists...' AS info;
DROP USER IF EXISTS 'teacher'@'localhost';

SELECT '[LOG] Creating user teacher...' AS info;
CREATE USER 'teacher'@'localhost';

SELECT '[LOG] Granting teacher privileges...' AS info;
GRANT DELETE ON ednevnik_workspace.pupil_tenant TO 'teacher'@'localhost' WITH GRANT OPTION;
GRANT INSERT, DELETE ON ednevnik_workspace.invite_index TO 'teacher'@'localhost' WITH GRANT OPTION;
GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.primary_school_final_grades TO 'teacher'@'localhost' WITH GRANT OPTION;
GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.primary_school_behaviour_grades TO 'teacher'@'localhost' WITH GRANT OPTION;
GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.high_school_final_grades TO 'teacher'@'localhost' WITH GRANT OPTION;
GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.high_school_behaviour_grades TO 'teacher'@'localhost' WITH GRANT OPTION;


SELECT '[LOG] Dropping user pupil if exists...' AS info;
DROP USER IF EXISTS 'pupil'@'localhost';

SELECT '[LOG] Creating user pupil...' AS info;
CREATE USER 'pupil'@'localhost';

SELECT '[LOG] Granting pupil privileges...' AS info;
GRANT SELECT on ednevnik_workspace.* TO 'pupil'@'localhost' WITH GRANT OPTION;


SELECT '[LOG] Dropping user service_reader if exists...' AS info;
DROP USER IF EXISTS 'service_reader'@'localhost';

SELECT '[LOG] Creating user service_reader...' AS info;
CREATE USER 'service_reader'@'localhost';

SELECT '[LOG] Granting service_reader privileges...' AS info;
GRANT SELECT ON ednevnik_workspace.* TO 'service_reader'@'localhost' WITH GRANT OPTION;
GRANT DELETE ON ednevnik_workspace.pupil_tenant TO 'service_reader'@'localhost' WITH GRANT OPTION;
GRANT INSERT, DELETE ON ednevnik_workspace.invite_index TO 'service_reader'@'localhost' WITH GRANT OPTION;
GRANT TRIGGER ON *.* TO 'service_reader'@'localhost' WITH GRANT OPTION;

SELECT '[LOG] Granting service DB user workspace privileges...' AS info;
GRANT INSERT ON ednevnik_workspace.pupil_tenant TO 'service_reader'@'localhost' WITH GRANT OPTION;
GRANT INSERT ON ednevnik_workspace.teacher_tenant TO 'service_reader'@'localhost' WITH GRANT OPTION;
GRANT INSERT, UPDATE ON ednevnik_workspace.accounts TO 'service_reader'@'localhost' WITH GRANT OPTION;
GRANT UPDATE ON ednevnik_workspace.pupil_global TO 'service_reader'@'localhost' WITH GRANT OPTION;
GRANT UPDATE ON ednevnik_workspace.teachers TO 'service_reader'@'localhost' WITH GRANT OPTION;


FLUSH PRIVILEGES;
SELECT '[LOG] Users and privileges created.' AS info;


-- Trigger on teachers table to prevent SQL injection on INSERT
SELECT '[LOG] Creating triggers for security...' AS info;
DELIMITER $$

CREATE DEFINER='service_reader'@'localhost' TRIGGER prevent_suspicious_email_insert
BEFORE INSERT ON accounts
FOR EACH ROW
BEGIN
  IF NEW.email NOT LIKE '%@%.%' OR
     INSTR(NEW.email, ' ') > 0 OR
     INSTR(LOWER(NEW.email), 'select') > 0 OR
     INSTR(LOWER(NEW.email), 'insert') > 0 OR
     INSTR(LOWER(NEW.email), 'update') > 0 OR
     INSTR(LOWER(NEW.email), 'delete') > 0 OR
     INSTR(LOWER(NEW.email), 'drop') > 0 OR
     INSTR(LOWER(NEW.email), 'create') > 0 OR
     INSTR(LOWER(NEW.email), 'alter') > 0 OR
     INSTR(LOWER(NEW.email), 'union') > 0 OR
     INSTR(LOWER(NEW.email), ' or ') > 0 OR
     INSTR(LOWER(NEW.email), ' and ') > 0 OR
     INSTR(LOWER(NEW.email), 'where') > 0 OR
     INSTR(LOWER(NEW.email), 'from') > 0 OR
     INSTR(LOWER(NEW.email), 'into') > 0 OR
     INSTR(LOWER(NEW.email), 'values') > 0 OR
     INSTR(LOWER(NEW.email), '--') > 0 OR
     INSTR(NEW.email, '/*') > 0 OR
     INSTR(NEW.email, '*/') > 0 OR
     INSTR(NEW.email, "'") > 0 OR
     INSTR(NEW.email, '"') > 0 OR
     INSTR(NEW.email, ';') > 0 OR
     INSTR(NEW.email, '\\') > 0 OR
     INSTR(LOWER(NEW.email), 'exec') > 0 OR
     INSTR(LOWER(NEW.email), 'execute') > 0
  THEN
    SIGNAL SQLSTATE '45000'
    SET MESSAGE_TEXT = 'Suspicious email address detected (insert)!';
  END IF;
END$$

CREATE DEFINER='service_reader'@'localhost' TRIGGER prevent_suspicious_email_update
BEFORE UPDATE ON accounts
FOR EACH ROW
BEGIN
  IF NEW.email NOT LIKE '%@%.%' OR
     INSTR(NEW.email, ' ') > 0 OR
     INSTR(LOWER(NEW.email), 'select') > 0 OR
     INSTR(LOWER(NEW.email), 'insert') > 0 OR
     INSTR(LOWER(NEW.email), 'update') > 0 OR
     INSTR(LOWER(NEW.email), 'delete') > 0 OR
     INSTR(LOWER(NEW.email), 'drop') > 0 OR
     INSTR(LOWER(NEW.email), 'create') > 0 OR
     INSTR(LOWER(NEW.email), 'alter') > 0 OR
     INSTR(LOWER(NEW.email), 'union') > 0 OR
     INSTR(LOWER(NEW.email), ' or ') > 0 OR
     INSTR(LOWER(NEW.email), ' and ') > 0 OR
     INSTR(LOWER(NEW.email), 'where') > 0 OR
     INSTR(LOWER(NEW.email), 'from') > 0 OR
     INSTR(LOWER(NEW.email), 'into') > 0 OR
     INSTR(LOWER(NEW.email), 'values') > 0 OR
     INSTR(LOWER(NEW.email), '--') > 0 OR
     INSTR(NEW.email, '/*') > 0 OR
     INSTR(NEW.email, '*/') > 0 OR
     INSTR(NEW.email, "'") > 0 OR
     INSTR(NEW.email, '"') > 0 OR
     INSTR(NEW.email, ';') > 0 OR
     INSTR(NEW.email, '\\') > 0 OR
     INSTR(LOWER(NEW.email), 'exec') > 0 OR
     INSTR(LOWER(NEW.email), 'execute') > 0
  THEN
    SIGNAL SQLSTATE '45000'
    SET MESSAGE_TEXT = 'Suspicious email address detected (update)!';
  END IF;
END$$

DELIMITER ;
SELECT '[LOG] Triggers created. Workspace DB setup complete.' AS info;

FLUSH PRIVILEGES;
