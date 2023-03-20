select <col1>,<col2>,count(*) as c from <table> group by <col1>,<col2> having c > 1;

CREATE UNIQUE INDEX CONCURRENTLY equipment_equip_id ON equipment (equip_id);

ALTER TABLE <table> DROP CONSTRAINT constraint_name;

ALTER TABLE equipment ADD CONSTRAINT unique_equip_id UNIQUE USING INDEX equipment_equip_id;
