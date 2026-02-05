-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE exercises (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    muscle_group TEXT NOT NULL,
    description TEXT
);

CREATE TABLE athletes(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    age INT,
    username TEXT NOT NULL UNIQUE,
    name TEXT,
    password_hash TEXT NOT NULL,
    current_cycle TEXT DEFAULT 'maintenance',
    created_at TIMESTAMP DEFAULT NOW(),
    role TEXT,
    email TEXT NOT NULL UNIQUE,
    gender TEXT
);



CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    total_time INT,
    grade_of_training INT,
    date_of_training TIMESTAMP DEFAULT NOW(),
    athlete_id UUID NOT NULL,
    FOREIGN KEY (athlete_id) REFERENCES athletes(id)
);

CREATE TABLE sets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exercise_id INT NOT NULL,
    workout_id UUID NOT NULL,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE,
    FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    weight FLOAT NOT NULL,
    set_order INT,
    reps INT NOT NULL,
    rpe INT
);

CREATE TABLE body_measurements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    athlete_id UUID NOT NULL,
    FOREIGN KEY (athlete_id) REFERENCES athletes(id) ON DELETE CASCADE,
    weight FLOAT NOT NULL,
    height FLOAT,
    chest_volume FLOAT,
    biceps_volume FLOAT,
    hip_volume FLOAT,
    waist_volume FLOAT,
    calves_volume FLOAT,
    forearm_volume FLOAT,
    gluteal_volume FLOAT,
    date_of_measuring TIMESTAMP DEFAULT NOW()
);


-- +goose Down
DROP TABLE body_measurements;
DROP TABLE sets;
DROP TABLE workouts;
DROP TABLE athletes;
DROP TABLE exercises;