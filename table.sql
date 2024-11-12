
CREATE table students (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    date_of_birth DATE,
    gender TEXT,
    email TEXT UNIQUE,
    identity_number TEXT,
    phone_number TEXT,
    address TEXT,
    password TEXT,
    class TEXT,
    school_year TEXT,
    field_of_study TEXT
);
CREATE table courses (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    number_of_credit INT NOT NULL
);
CREATE table course_results (
    id SERIAL PRIMARY KEY,
    course_id TEXT REFERENCES courses(id),
    student_id TEXT REFERENCES students(id),
    semester_number INT NOT NULL,
    academic_year TEXT NOT NULL
);
CREATE table component_scores(
    id SERIAL PRIMARY KEY,
    course_result_id INT REFERENCES course_results(id),
    name TEXT,
    score_weight DOUBLE PRECISION,
    score DOUBLE PRECISION
);
INSERT INTO courses (id, name, number_of_credit) VALUES
('CSE101', 'Introduction to Computer Science', 3),
('MAT201', 'Calculus II', 4),
('PHY301', 'Physics I', 4),
('ENG101', 'English Literature', 2),
('HIS202', 'World History', 3);
INSERT INTO course_results (course_id, student_id, semester_number, academic_year) VALUES
('CSE101', '123457', 1, '2023-2024'),
('MAT201', '123457', 1, '2023-2024'),
('PHY301', '123457', 1, '2023-2024'),
('ENG101', '123457', 1, '2023-2024'),
('HIS202', '123457', 1, '2023-2024');

INSERT INTO component_scores (course_result_id, name, score_weight, score) VALUES
(1, 'Midterm Exam', 0.4, 8.0),
(1, 'Final Exam', 0.6, 9.0);

INSERT INTO component_scores (course_result_id, name, score_weight, score) VALUES
(2, 'Midterm Exam', 0.3, 7.0),
(2, 'Final Exam', 0.7, 8.0);

INSERT INTO component_scores (course_result_id, name, score_weight, score) VALUES
(3, 'Lab Work', 0.2, 8.0),
(3, 'Midterm Exam', 0.3, 7.0),
(3, 'Final Exam', 0.5, 8.0);

INSERT INTO component_scores (course_result_id, name, score_weight, score) VALUES
(4, 'Midterm Essay', 0.5, 7.0),
(4, 'Final Essay', 0.5, 8.0);

INSERT INTO component_scores (course_result_id, name, score_weight, score) VALUES
(5, 'Quiz', 0.2, 7.0),
(5, 'Midterm Exam', 0.4, 7.0),
(5, 'Final Exam', 0.4, 7.0);

DELETE from component_scores;

SELECT name, score_weight, score, course_result_id
FROM component_scores
WHERE course_result_id = ANY(ARRAY[1, 2, 3, 4, 5]);

