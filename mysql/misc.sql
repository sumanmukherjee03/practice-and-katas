INSERT INTO audit_log
SELECT *
FROM audit_log_tmp
ON CONFLICT DO NOTHING;

select count(*) from audit_log where date >= NOW() - INTERVAL '8 Hours';

SELECT current_database() AS database, :'HOST' AS hostname;

DO $$ -- for the purpose of local variable declaration
DECLARE clientOrderId varchar := '<someid>';
BEGIN
  UPDATE <table_name>
SET
<column_name> = <value>
WHERE client_order_id = clientOrderId;
END;
$$;
