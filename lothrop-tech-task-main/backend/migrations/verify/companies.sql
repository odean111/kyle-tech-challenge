-- Verify lothrop-backend:companies on pg

BEGIN;

SELECT id, jurisdiction, company_name, company_address, nature_of_business, 
       number_of_directors, number_of_shareholders, sec_code, 
       date_created, date_updated
FROM companies
WHERE FALSE;

ROLLBACK;
