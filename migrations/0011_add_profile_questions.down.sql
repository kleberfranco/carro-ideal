-- Remove as 4 perguntas de perfil adicionadas na 0011.
-- O ON DELETE CASCADE em answer_options remove as opções automaticamente.
DELETE FROM questions
WHERE display_order IN (6, 7, 8, 9)
  AND question_text IN (
      'Quantas pessoas costumam andar no carro?',
      'Quanto você dirige por mês?',
      'Qual tipo de combustível você prefere?',
      'O que mais pesa na sua decisão?'
  );
