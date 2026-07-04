ALTER TABLE ledger_entries ADD COLUMN IF NOT EXISTS label TEXT;

UPDATE ledger_entries
SET label = note
WHERE label IS NULL
  AND type = 'accrual'
  AND note IS NOT NULL;
