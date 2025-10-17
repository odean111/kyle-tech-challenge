-- Revert lothrop-backend:companies from pg

BEGIN;

DROP TABLE IF EXISTS companies;
DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;
