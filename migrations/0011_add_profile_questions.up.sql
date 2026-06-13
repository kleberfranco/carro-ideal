-- Amplia o questionário com 4 perguntas de perfil (passageiros, quilometragem,
-- combustível e prioridade de decisão). As opções mapeiam para dimensões já
-- existentes no match_profile dos veículos, mantendo o algoritmo de scoring
-- (fallback) coerente; o ChatGPT aproveita o texto integral das respostas.

INSERT INTO questions (question_text, weight, display_order)
SELECT v.question_text, v.weight, v.display_order
FROM (VALUES
    ('Quantas pessoas costumam andar no carro?', 1.100, 6),
    ('Quanto você dirige por mês?', 1.000, 7),
    ('Qual tipo de combustível você prefere?', 0.700, 8),
    ('O que mais pesa na sua decisão?', 1.200, 9)
) AS v(question_text, weight, display_order)
WHERE NOT EXISTS (
    SELECT 1 FROM questions q WHERE q.display_order = v.display_order
);

INSERT INTO answer_options (question_id, option_text, score_profile, display_order)
SELECT q.id, value.option_text, value.profile::jsonb, value.display_order
FROM questions q
JOIN (VALUES
    -- 6: Quantas pessoas costumam andar no carro?
    (6, 'Geralmente 1 ou 2 pessoas', '{"compact":1,"urban":0.4}', 1),
    (6, 'Uma família (3 a 5 pessoas)', '{"family":1,"space":0.6}', 2),
    (6, 'Preciso levar muita gente ou bagagem', '{"space":1,"cargo":0.5,"family":0.6}', 3),
    -- 7: Quanto você dirige por mês?
    (7, 'Pouco — uso leve na cidade', '{"urban":1,"compact":0.3}', 1),
    (7, 'Uso moderado — cidade e viagens ocasionais', '{"mixed":1,"comfort":0.4}', 2),
    (7, 'Rodo muito — economia é essencial', '{"efficiency":1,"comfort":0.4}', 3),
    -- 8: Qual tipo de combustível você prefere?
    (8, 'Flex (gasolina/etanol)', '{"efficiency":0.3}', 1),
    (8, 'Diesel — força e economia na estrada', '{"offroad":0.6,"cargo":0.6,"performance":0.4}', 2),
    (8, 'Híbrido ou elétrico — máxima economia', '{"efficiency":1,"comfort":0.4}', 3),
    (8, 'Sem preferência', '{"efficiency":0.3}', 4),
    -- 9: O que mais pesa na sua decisão?
    (9, 'Menor preço de compra', '{"budget_low":1,"efficiency":0.4}', 1),
    (9, 'Custo-benefício equilibrado', '{"budget_mid":1,"efficiency":0.5}', 2),
    (9, 'Conforto e tecnologia', '{"comfort":1,"automatic":0.5}', 3),
    (9, 'Segurança e robustez', '{"space":0.6,"performance":0.6,"offroad":0.5}', 4)
) AS value(question_order, option_text, profile, display_order)
  ON q.display_order = value.question_order
ON CONFLICT (question_id, display_order) DO NOTHING;
