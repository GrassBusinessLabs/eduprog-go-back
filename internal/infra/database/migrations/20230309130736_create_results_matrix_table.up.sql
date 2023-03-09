CREATE TABLE IF NOT EXISTS public.results_matrix (
    eduprog_id integer NOT NULL,
    component_id integer NOT NULL,
    eduprogresult_id integer NOT NULL
);

ALTER TABLE public.results_matrix
    ADD FOREIGN KEY (component_id) REFERENCES public.eduprogcomp(id) ON DELETE CASCADE;

ALTER TABLE public.results_matrix
    ADD FOREIGN KEY (eduprogresult_id) REFERENCES public.eduprogcompetencies(id) ON DELETE CASCADE;

ALTER TABLE public.results_matrix
    ADD CONSTRAINT PK_results_matrix PRIMARY KEY (component_id,eduprogresult_id);