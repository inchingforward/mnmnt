-- Adds the memory.edit_uuid column to allow for editing
-- approved memories through the site.
alter table memory add column edit_uuid uuid;

commit;