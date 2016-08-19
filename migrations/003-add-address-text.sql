-- Adds the memory.address_text column to capture the original 
-- location/address entered by the user.
alter table memory add column address_text text;

commit;