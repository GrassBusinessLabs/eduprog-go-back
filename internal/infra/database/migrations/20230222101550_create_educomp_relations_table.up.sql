CREATE TABLE IF NOT EXISTS public.educomp_relations (
    eduprog_id integer NOT NULL,
    base_comp_id integer NOT NULL,
    child_comp_id integer NOT NULL

);

ALTER TABLE public.educomp_relations ADD FOREIGN KEY (base_comp_id) REFERENCES public.eduprogcomp(id);

ALTER TABLE public.educomp_relations ADD FOREIGN KEY (child_comp_id) REFERENCES public.eduprogcomp(id);