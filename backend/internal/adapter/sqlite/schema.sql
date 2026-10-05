CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name_kk TEXT NOT NULL,
    name_ru TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    kazakh TEXT NOT NULL,
    russian TEXT NOT NULL,
    transcription TEXT NOT NULL DEFAULT '',
    example_kk TEXT NOT NULL DEFAULT '',
    example_ru TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_words_category ON words(category_id);

CREATE TABLE IF NOT EXISTS quizzes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL,
    finished_at TEXT
);

CREATE TABLE IF NOT EXISTS quiz_questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    quiz_id INTEGER NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    word_id INTEGER NOT NULL,
    prompt TEXT NOT NULL,
    options_json TEXT NOT NULL,
    correct_index INTEGER NOT NULL,
    selected_index INTEGER,
    correct INTEGER
);

CREATE INDEX IF NOT EXISTS idx_quiz_questions_quiz ON quiz_questions(quiz_id);

CREATE TABLE IF NOT EXISTS progress (
    category_id INTEGER PRIMARY KEY REFERENCES categories(id) ON DELETE CASCADE,
    total_answers INTEGER NOT NULL DEFAULT 0,
    correct_answers INTEGER NOT NULL DEFAULT 0,
    last_studied_at TEXT
);

CREATE TABLE IF NOT EXISTS situations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL UNIQUE,
    title_kk TEXT NOT NULL,
    title_ru TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    position INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS situation_words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    situation_id INTEGER NOT NULL REFERENCES situations(id) ON DELETE CASCADE,
    kazakh TEXT NOT NULL,
    russian TEXT NOT NULL,
    transcription TEXT NOT NULL DEFAULT '',
    position INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_situation_words_situation ON situation_words(situation_id, position);

CREATE TABLE IF NOT EXISTS units (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    situation_id INTEGER NOT NULL REFERENCES situations(id) ON DELETE CASCADE,
    title_kk TEXT NOT NULL,
    title_ru TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    position INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_units_situation ON units(situation_id);

CREATE TABLE IF NOT EXISTS lessons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    title_kk TEXT NOT NULL,
    title_ru TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    position INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_lessons_unit ON lessons(unit_id, position);

CREATE TABLE IF NOT EXISTS exercises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    kind TEXT NOT NULL,
    prompt TEXT NOT NULL,
    answer TEXT NOT NULL,
    options_json TEXT NOT NULL DEFAULT '[]',
    position INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_exercises_lesson ON exercises(lesson_id, position);
