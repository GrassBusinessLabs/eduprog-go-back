CREATE TABLE IF NOT EXISTS public.competencies_matrix (
    eduprog_id integer NOT NULL,
    component_id integer NOT NULL,
    competency_id integer NOT NULL
);

ALTER TABLE public.competencies_matrix
    ADD FOREIGN KEY (component_id) REFERENCES public.eduprogcomp(id) ON DELETE CASCADE;

ALTER TABLE public.competencies_matrix
    ADD FOREIGN KEY (competency_id) REFERENCES public.eduprogcompetencies(id) ON DELETE CASCADE;

ALTER TABLE public.competencies_matrix
    ADD CONSTRAINT PK_competencies_matrix PRIMARY KEY (component_id,competency_id);