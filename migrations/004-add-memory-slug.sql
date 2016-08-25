-- Adds the memory.slug column which will be used for the 
-- resource URI
alter table memory add column slug text;

update memory set slug = '';

commit;