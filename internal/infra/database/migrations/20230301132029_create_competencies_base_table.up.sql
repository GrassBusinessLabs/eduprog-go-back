CREATE TABLE IF NOT EXISTS public.competencies_base (
    id SERIAL PRIMARY KEY,
    "type" varchar(50) NOT NULL,
    code integer,
    definition varchar(500),
    specialty varchar(50),
    education_level varchar(100)
);

CREATE OR REPLACE FUNCTION competencies_base_education_level_trigger()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.education_level = 'ENTRY' THEN
        NEW.education_level := 'Початковий рівень (короткий цикл)';
    ELSIF NEW.education_level = 'FIRST' THEN
        NEW.education_level := 'Перший (бакалаврський) рівень';
    ELSIF NEW.education_level = 'SECOND' THEN
        NEW.education_level := 'Другий (магістерський) рівень';
    ELSIF NEW.education_level = 'THIRD' THEN
        NEW.education_level := 'Третій (освітньо-науковий/освітньо-творчий) рівень';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER competencies_base_education_level_trigger
    BEFORE INSERT ON competencies_base
    FOR EACH ROW
EXECUTE FUNCTION competencies_base_education_level_trigger();