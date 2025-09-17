
CREATE TABLE pupils (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    jmbg VARCHAR(13) UNIQUE,
    gender ENUM('M', 'F'),
    address VARCHAR(200),
    guardian_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) UNIQUE,
    guardian_number VARCHAR(20) UNIQUE,
    date_of_birth DATE,
    religion ENUM ('Islam', 'Catholic', 'Orthodox', 'Jewish', 'Other', 'NotAttendingReligion'),
    place_of_birth VARCHAR(100) NOT NULL,
    account_id INT,
    FOREIGN KEY (id) REFERENCES ednevnik_workspace.pupil_global(id),
    FOREIGN KEY (account_id) REFERENCES ednevnik_workspace.accounts(id)
);
CREATE INDEX idx_pupil_last_name_name ON pupils (last_name, name);

CREATE TABLE sections (
    id INT PRIMARY KEY AUTO_INCREMENT,
    section_code VARCHAR(10) NOT NULL,
    class_code VARCHAR(10),
    year VARCHAR(30) NOT NULL,
    tenant_id INT,
    curriculum_code VARCHAR(30),
    archived BOOLEAN DEFAULT FALSE,
    CONSTRAINT unique_section_class_year UNIQUE (section_code, class_code, year),
    CONSTRAINT check_section_year CHECK (
        year LIKE '____/____' AND
        LENGTH(year) = 9 AND
        SUBSTRING(year, 1, 4) REGEXP '^[0-9]{4}$' AND
        SUBSTRING(year, 6, 4) REGEXP '^[0-9]{4}$' AND
        CAST(SUBSTRING(year, 6, 4) AS UNSIGNED) = CAST(SUBSTRING(year, 1, 4) AS UNSIGNED) + 1
    ),
    FOREIGN KEY (class_code) REFERENCES ednevnik_workspace.classes(class_code),
    FOREIGN KEY (tenant_id) REFERENCES ednevnik_workspace.tenant(id),
    FOREIGN KEY (curriculum_code) REFERENCES ednevnik_workspace.curriculum(curriculum_code)
);

CREATE TABLE homeroom_assignments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    section_id INT NOT NULL,
    teacher_id INT NOT NULL,
    UNIQUE (section_id),
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE,
    FOREIGN KEY (teacher_id) REFERENCES ednevnik_workspace.teachers(id)
);
CREATE INDEX idx_homeroom_teacher_section ON homeroom_assignments (teacher_id, section_id);

CREATE TABLE pupils_sections (
    pupil_id INT,
    section_id INT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY(pupil_id, section_id),
    FOREIGN KEY (pupil_id) REFERENCES pupils(id) ON DELETE CASCADE,
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
) WITH SYSTEM VERSIONING;

CREATE TABLE pupils_sections_invite (
    id INT PRIMARY KEY AUTO_INCREMENT,
    pupil_id INT,
    section_id INT,
    invite_date DATE DEFAULT CURRENT_DATE,
    status ENUM('pending', 'accepted', 'declined') DEFAULT 'pending',
    FOREIGN KEY (pupil_id) REFERENCES ednevnik_workspace.pupil_global(id),
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
);
CREATE INDEX idx_psi_section_pupil_status ON pupils_sections_invite (
    section_id, pupil_id, status
);
CREATE INDEX idx_psi_section_status ON pupils_sections_invite (
    section_id, status
);


CREATE TABLE teachers_sections_invite (
    id INT PRIMARY KEY AUTO_INCREMENT,
    teacher_id INT,
    section_id INT,
    invite_date DATE DEFAULT CURRENT_DATE,
    status ENUM('pending', 'accepted', 'declined') DEFAULT 'pending',
    homeroom_teacher BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (teacher_id) REFERENCES ednevnik_workspace.teachers(id),
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
);
CREATE INDEX idx_tsi_id_homeroom ON teachers_sections_invite (id, homeroom_teacher);

CREATE TABLE teachers_sections_invite_subjects (
    invite_id INT,
    subject_code VARCHAR(15),
    PRIMARY KEY(invite_id, subject_code),
    FOREIGN KEY (invite_id) REFERENCES teachers_sections_invite(id) ON DELETE CASCADE,
    FOREIGN KEY (subject_code) REFERENCES ednevnik_workspace.subjects(subject_code)
);

CREATE TABLE teachers_sections (
    teacher_id INT,
    section_id INT,
    PRIMARY KEY(teacher_id, section_id),
    FOREIGN KEY (teacher_id) REFERENCES ednevnik_workspace.teachers(id),
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
);

CREATE TABLE teachers_sections_subjects (
    section_id INT,
    subject_code VARCHAR(15),
    teacher_id INT,
    PRIMARY KEY(teacher_id, subject_code, section_id),
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE,
    FOREIGN KEY (subject_code) REFERENCES ednevnik_workspace.subjects(subject_code),
    FOREIGN KEY (teacher_id) REFERENCES ednevnik_workspace.teachers(id)
);
CREATE INDEX idx_tss_teacher_section ON teachers_sections_subjects (teacher_id, section_id);

CREATE TABLE student_grades (
    id INT PRIMARY KEY AUTO_INCREMENT,
    type ENUM('exam', 'oral', 'written_assignment', 'final'),
    pupil_id INT,
    section_id INT,
    subject_code VARCHAR(15),
    grade INT CHECK (grade BETWEEN 1 AND 5),
    grade_date DATE DEFAULT CURRENT_DATE,
    teacher_id INT,
    semester_code VARCHAR(10),
    signature VARCHAR(128),
    FOREIGN KEY (pupil_id) REFERENCES pupils(id) ON DELETE CASCADE,
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE,
    FOREIGN KEY (subject_code) REFERENCES ednevnik_workspace.subjects(subject_code),
    FOREIGN KEY (teacher_id) REFERENCES ednevnik_workspace.teachers(id),
    FOREIGN KEY (semester_code) REFERENCES ednevnik_workspace.semester(semester_code)
) WITH SYSTEM VERSIONING;

CREATE INDEX idx_grade_pupil_section_type_semester ON student_grades (
    pupil_id, section_id, type, semester_code
);
CREATE INDEX idx_grade_section_subject_semester_date ON student_grades (
    section_id, subject_code, semester_code, grade_date
);
CREATE INDEX idx_grade_section_type ON student_grades (
    section_id, type
);

CREATE TABLE pupil_behaviour (
    id INT PRIMARY KEY AUTO_INCREMENT,
    pupil_id INT NOT NULL,
    section_id INT NOT NULL,
    behaviour ENUM('primjerno', 'vrlodobro', 'dobro', 'zadovoljavajuće', 'loše') NOT NULL DEFAULT 'primjerno',
    semester_code VARCHAR(10),
    signature VARCHAR(128),
    FOREIGN KEY (pupil_id, section_id) REFERENCES pupils_sections(pupil_id, section_id) ON DELETE CASCADE,
    FOREIGN KEY (semester_code) REFERENCES ednevnik_workspace.semester(semester_code)
) WITH SYSTEM VERSIONING;

CREATE INDEX idx_behaviour_pupil_section_semester ON pupil_behaviour (
    pupil_id, section_id, semester_code
);

CREATE TABLE time_periods (
    id INT PRIMARY KEY AUTO_INCREMENT,
    section_id INT NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    batch_id VARCHAR(50) NOT NULL,
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
) WITH SYSTEM VERSIONING;
CREATE INDEX idx_time_period_start_time ON time_periods (start_time);

CREATE TABLE classroom (
    code VARCHAR(20) PRIMARY KEY,
    capacity INT NOT NULL,
    type VARCHAR(40)
);

CREATE TABLE schedule (
    id INT PRIMARY KEY AUTO_INCREMENT,
    section_id INT NOT NULL,
    time_period_id INT NOT NULL,
    subject_code VARCHAR(15) NOT NULL,
    weekday ENUM('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday') NOT NULL,
    classroom_code VARCHAR(40),
    type ENUM('regular', 'additional', 'remedial') NOT NULL DEFAULT 'regular',
    batch_id VARCHAR(50) NOT NULL,
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE,
    FOREIGN KEY (time_period_id) REFERENCES time_periods(id) ON DELETE CASCADE,
    FOREIGN KEY (subject_code) REFERENCES ednevnik_workspace.subjects(subject_code),
    FOREIGN KEY (classroom_code) REFERENCES classroom(code) ON DELETE SET NULL
) WITH SYSTEM VERSIONING;

CREATE TABLE class_lesson (
    id INT PRIMARY KEY AUTO_INCREMENT,
    description VARCHAR(255),
    date DATE NOT NULL,
    period_number INT NOT NULL DEFAULT 1,
    section_id INT,
    subject_code VARCHAR(15),
    signature VARCHAR(128),
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE,
    FOREIGN KEY (subject_code) REFERENCES ednevnik_workspace.subjects(subject_code)
) WITH SYSTEM VERSIONING;

CREATE INDEX idx_class_lesson_date_period ON class_lesson (
    date, period_number
);

CREATE TABLE pupil_attendance (
    pupil_id INT,
    lesson_id INT,
    PRIMARY KEY(pupil_id, lesson_id),
    status ENUM('present', 'absent', 'unexcused', 'excused'),
    FOREIGN KEY (pupil_id) REFERENCES pupils(id) ON DELETE CASCADE,
    FOREIGN KEY (lesson_id) REFERENCES class_lesson(id) ON DELETE CASCADE
) WITH SYSTEM VERSIONING;

DELIMITER $$
CREATE DEFINER='service_reader'@'localhost' TRIGGER create_pupil_behaviour_after_pupil_section_insert
AFTER INSERT ON pupils_sections
FOR EACH ROW
BEGIN
    DECLARE signature_name VARCHAR(128);

    SELECT CONCAT(t.name, ' ', t.last_name)
    INTO signature_name
    FROM homeroom_assignments ha
    JOIN ednevnik_workspace.teachers t ON ha.teacher_id = t.id
    WHERE ha.section_id = NEW.section_id
    LIMIT 1;

    IF signature_name IS NULL THEN
        SELECT CONCAT(ta.name, ' ', ta.last_name)
        INTO signature_name
        FROM sections s
        JOIN ednevnik_workspace.tenant ten ON s.tenant_id = ten.id
        JOIN ednevnik_workspace.teachers ta ON ten.tenant_admin_id = ta.id
        WHERE s.id = NEW.section_id
        LIMIT 1;
    END IF;

    INSERT INTO pupil_behaviour (pupil_id, section_id, semester_code, signature)
    SELECT
        NEW.pupil_id,
        NEW.section_id,
        s.semester_code,
        signature_name
    FROM ednevnik_workspace.semester s;
END$$

DELIMITER ;

DELIMITER $$
CREATE DEFINER='service_reader'@'localhost' TRIGGER prevent_multiple_final_grades_insert
BEFORE INSERT ON student_grades
FOR EACH ROW
BEGIN
    DECLARE existing_count INT DEFAULT 0;
    IF NEW.type = 'final' THEN
        SELECT COUNT(*) INTO existing_count
        FROM student_grades
        WHERE pupil_id = NEW.pupil_id
          AND section_id = NEW.section_id
          AND subject_code = NEW.subject_code
          AND semester_code = NEW.semester_code
          AND type = 'final';

        IF existing_count > 0 THEN
            SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'A final grade already exists for this pupil-section-subject-semester combination';
        END IF;
    END IF;
END$$

DELIMITER ;

DELIMITER $$
CREATE DEFINER='service_reader'@'localhost' TRIGGER prevent_multiple_final_grades_update
BEFORE UPDATE ON student_grades
FOR EACH ROW
BEGIN
    DECLARE existing_count INT DEFAULT 0;

    IF NEW.type = 'final' AND OLD.type != 'final' THEN
        SELECT COUNT(*) INTO existing_count
        FROM student_grades
        WHERE pupil_id = NEW.pupil_id
          AND section_id = NEW.section_id
          AND subject_code = NEW.subject_code
          AND semester_code = NEW.semester_code
          AND type = 'final'
          AND id != NEW.id;

        IF existing_count > 0 THEN
            SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'A final grade already exists for this pupil-section-subject-semester combination';
        END IF;
    END IF;
END$$
DELIMITER ;

FLUSH PRIVILEGES;
