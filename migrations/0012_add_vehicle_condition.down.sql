-- Reverte a introdução da condição (novo/seminovo).
DELETE FROM questions WHERE display_order = 10
  AND question_text = 'Você prefere carro novo ou seminovo?';

DELETE FROM vehicles WHERE condition = 'seminovo';

UPDATE vehicles SET match_profile = match_profile - 'novo';

ALTER TABLE vehicles DROP COLUMN IF EXISTS condition;
