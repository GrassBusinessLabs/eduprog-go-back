BEGIN;

ALTER TABLE eduprog.public.eduprogscheme ADD COLUMN discipline_id integer references discipline(id);

COMMIT;