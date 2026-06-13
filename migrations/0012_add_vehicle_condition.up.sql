-- Introduz a noção de condição (novo/seminovo) nos veículos, popula o catálogo
-- com seminovos (carros maiores e mais antigos em faixas de preço menores) e
-- adiciona a pergunta correspondente ao questionário. A dimensão "novo"/
-- "seminovo" no match_profile mantém o algoritmo de scoring coerente; o ChatGPT
-- usa a coluna condition diretamente no prompt.

-- 1. Coluna de condição. Veículos existentes ficam como 'novo' por padrão.
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS condition VARCHAR(20) NOT NULL DEFAULT 'novo';

-- 2. Marca os veículos atuais (todos novos) com a dimensão "novo" para o scoring.
UPDATE vehicles SET match_profile = match_profile || '{"novo":1}'::jsonb
WHERE condition = 'novo';

-- 3. Catálogo de seminovos — carros grandes/confortáveis em preços acessíveis.
INSERT INTO vehicles (
    category_id, brand, model, version, year, fuel_type, transmission,
    price_min, price_max, seats, trunk_capacity, consumption_city,
    consumption_highway, description, strengths, weaknesses, condition, match_profile
)
SELECT c.id, v.brand, v.model, v.version, v.year, v.fuel, v.transmission,
       v.price_min, v.price_max, v.seats, v.trunk, v.city, v.highway,
       v.description, v.strengths, v.weaknesses, 'seminovo', v.profile::jsonb
FROM vehicle_categories c
JOIN (VALUES
    ('Sedan','Renault','Logan','Expression 1.6',2018,'Flex','Manual',48000,55000,5,510,9.6,12.8,'Sedã espaçoso e econômico por preço de hatch.','Porta-malas enorme, baixo custo e espaço interno.','Acabamento simples e desempenho modesto.','{"seminovo":1,"budget_low":1,"family":0.7,"space":0.7,"efficiency":0.6,"mixed":0.5,"manual":1}'),
    ('Sedan','Chevrolet','Cobalt','LTZ 1.8',2016,'Flex','Automático',52000,60000,5,563,8.9,12.0,'Sedã grande e confortável com câmbio automático.','Espaço traseiro, porta-malas e conforto.','Consumo urbano elevado.','{"seminovo":1,"budget_low":1,"family":0.8,"space":0.8,"comfort":0.6,"mixed":0.5,"automatic":1}'),
    ('Sedan','Toyota','Corolla','XEi 2.0',2015,'Flex','CVT',72000,82000,5,470,9.4,13.1,'Sedã médio confiável e confortável, seminovo.','Confiabilidade, conforto e revenda forte.','Preço ainda alto para um seminovo.','{"seminovo":1,"budget_low":0.8,"family":1,"comfort":1,"mixed":1,"efficiency":0.5,"automatic":1}'),
    ('Sedan','Honda','Civic','LXR 2.0',2014,'Flex','Automático',62000,72000,5,440,8.6,12.4,'Sedã esportivo e confortável por bom preço.','Desempenho, conforto e durabilidade.','Manutenção mais cara que rivais.','{"seminovo":1,"budget_low":1,"comfort":0.9,"mixed":0.8,"performance":0.7,"automatic":1}'),
    ('SUV','Renault','Duster','Dynamique 2.0',2015,'Flex','Manual',55000,65000,5,475,8.5,11.3,'SUV robusto e espaçoso a preço acessível.','Vão livre, porta-malas e robustez.','Acabamento simples e consumo.','{"seminovo":1,"budget_low":1,"offroad":0.9,"space":1,"family":0.8,"cargo":0.7,"manual":1}'),
    ('SUV','Hyundai','ix35','GLS 2.0',2016,'Flex','Automático',70000,82000,5,465,7.8,10.9,'SUV confortável e bem equipado, seminovo.','Conforto, espaço e equipamentos.','Consumo elevado na cidade.','{"seminovo":1,"budget_low":0.8,"comfort":0.9,"family":1,"space":0.9,"automatic":1}'),
    ('Picape','Fiat','Toro','Freedom 1.8',2017,'Flex','Automático',78000,90000,5,937,8.8,11.4,'Picape versátil para trabalho e família.','Caçamba, conforto de SUV e câmbio automático.','Motor 1.8 aspirado modesto na estrada.','{"seminovo":1,"budget_mid":0.9,"cargo":1,"mixed":0.8,"offroad":0.6,"performance":0.7,"automatic":1}')
) AS v(category,brand,model,version,year,fuel,transmission,price_min,price_max,seats,trunk,city,highway,description,strengths,weaknesses,profile)
  ON c.name = v.category
WHERE NOT EXISTS (
    SELECT 1 FROM vehicles existing
    WHERE existing.brand = v.brand AND existing.model = v.model AND existing.year = v.year
);

-- 4. Pergunta de condição (display_order 10) + opções.
INSERT INTO questions (question_text, weight, display_order)
SELECT 'Você prefere carro novo ou seminovo?', 1.000, 10
WHERE NOT EXISTS (SELECT 1 FROM questions WHERE display_order = 10);

INSERT INTO answer_options (question_id, option_text, score_profile, display_order)
SELECT q.id, value.option_text, value.profile::jsonb, value.display_order
FROM questions q
JOIN (VALUES
    (10, 'Novo (0 km)', '{"novo":1}', 1),
    (10, 'Seminovo (melhor preço)', '{"seminovo":1}', 2),
    (10, 'Tanto faz', '{"novo":0.5,"seminovo":0.5}', 3)
) AS value(question_order, option_text, profile, display_order)
  ON q.display_order = value.question_order
ON CONFLICT (question_id, display_order) DO NOTHING;
