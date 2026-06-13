INSERT INTO vehicle_categories (name, description) VALUES
    ('Hatch', 'Compacto e eficiente para uso urbano'),
    ('Sedan', 'Conforto e porta-malas para família e estrada'),
    ('SUV', 'Espaço interno e versatilidade'),
    ('Picape', 'Capacidade de carga e uso misto')
ON CONFLICT (name) DO NOTHING;

INSERT INTO questions (question_text, weight, display_order) VALUES
    ('Qual é a sua faixa de orçamento?', 1.5, 1),
    ('Onde você usará mais o carro?', 1.2, 2),
    ('Quanto espaço interno você precisa?', 1.0, 3),
    ('Qual característica é mais importante?', 1.3, 4),
    ('Qual câmbio você prefere?', 0.8, 5)
ON CONFLICT DO NOTHING;

INSERT INTO answer_options (question_id, option_text, score_profile, display_order)
SELECT q.id, value.option_text, value.profile::jsonb, value.display_order
FROM questions q
JOIN (VALUES
    (1, 'Até R$ 80 mil', '{"budget_low":1}', 1),
    (1, 'De R$ 80 mil a R$ 130 mil', '{"budget_mid":1}', 2),
    (1, 'Acima de R$ 130 mil', '{"budget_high":1}', 3),
    (2, 'Principalmente na cidade', '{"urban":1,"efficiency":0.6}', 1),
    (2, 'Cidade e estrada', '{"mixed":1,"comfort":0.5}', 2),
    (2, 'Estradas e terrenos irregulares', '{"offroad":1,"space":0.4}', 3),
    (3, 'Compacto, para até duas pessoas', '{"compact":1}', 1),
    (3, 'Confortável para uma família', '{"family":1,"space":0.6}', 2),
    (3, 'Muito espaço ou capacidade de carga', '{"space":1,"cargo":1}', 3),
    (4, 'Economia de combustível', '{"efficiency":1}', 1),
    (4, 'Conforto e acabamento', '{"comfort":1}', 2),
    (4, 'Desempenho e robustez', '{"performance":1,"offroad":0.4}', 3),
    (5, 'Manual', '{"manual":1}', 1),
    (5, 'Automático ou CVT', '{"automatic":1}', 2),
    (5, 'Sem preferência', '{"manual":0.5,"automatic":0.5}', 3)
) AS value(question_order, option_text, profile, display_order)
  ON q.display_order = value.question_order
ON CONFLICT (question_id, display_order) DO NOTHING;

INSERT INTO vehicles (
    category_id, brand, model, version, year, fuel_type, transmission,
    price_min, price_max, seats, trunk_capacity, consumption_city,
    consumption_highway, description, strengths, weaknesses, match_profile
)
SELECT c.id, v.brand, v.model, v.version, v.year, v.fuel, v.transmission,
       v.price_min, v.price_max, v.seats, v.trunk, v.city, v.highway,
       v.description, v.strengths, v.weaknesses, v.profile::jsonb
FROM vehicle_categories c
JOIN (VALUES
    ('Hatch','Renault','Kwid','Zen 1.0',2025,'Flex','Manual',78000,82000,5,290,15.3,15.7,'Compacto econômico para a cidade.','Baixo consumo e manutenção acessível.','Acabamento simples e desempenho modesto.','{"budget_low":1,"urban":1,"compact":1,"efficiency":1,"manual":1}'),
    ('Hatch','Chevrolet','Onix','LT 1.0',2025,'Flex','Manual',90000,100000,5,275,13.8,16.9,'Hatch equilibrado para uso diário.','Segurança, conectividade e eficiência.','Porta-malas menor que sedãs.','{"budget_mid":0.8,"urban":1,"mixed":0.5,"compact":0.8,"efficiency":0.9,"manual":1}'),
    ('Hatch','Honda','City Hatchback','EXL',2025,'Flex','CVT',130000,145000,5,268,12.7,13.9,'Hatch confortável e confiável.','Conforto, confiabilidade e bom acabamento.','Preço elevado para a categoria.','{"budget_high":0.7,"urban":0.8,"mixed":0.8,"family":0.7,"comfort":0.9,"automatic":1}'),
    ('Sedan','Toyota','Corolla','GLi 2.0',2025,'Flex','CVT',155000,175000,5,470,11.9,14.5,'Sedã médio confortável para família e estrada.','Confiabilidade, conforto e revenda.','Preço de aquisição elevado.','{"budget_high":1,"mixed":1,"family":1,"comfort":1,"efficiency":0.6,"automatic":1}'),
    ('Sedan','Nissan','Versa','Advance',2025,'Flex','CVT',115000,130000,5,482,11.5,15.0,'Sedã espaçoso com bom porta-malas.','Espaço interno, conforto e consumo.','Desempenho apenas adequado.','{"budget_mid":1,"mixed":0.9,"family":1,"space":0.8,"comfort":0.8,"efficiency":0.7,"automatic":1}'),
    ('SUV','Volkswagen','T-Cross','Comfortline',2025,'Flex','Automático',155000,175000,5,373,11.2,13.5,'SUV compacto versátil e seguro.','Desempenho, segurança e espaço.','Acabamento rígido em algumas áreas.','{"budget_high":1,"urban":0.6,"mixed":1,"family":1,"space":0.8,"performance":0.8,"automatic":1}'),
    ('SUV','Hyundai','Creta','Comfort',2025,'Flex','Automático',145000,165000,5,422,11.9,12.6,'SUV confortável para família.','Conforto, espaço e garantia.','Consumo rodoviário mediano.','{"budget_high":0.9,"mixed":1,"family":1,"space":0.9,"comfort":1,"automatic":1}'),
    ('SUV','Renault','Duster','Intense',2025,'Flex','CVT',135000,150000,5,475,10.5,11.5,'SUV robusto e espaçoso.','Porta-malas, robustez e vão livre.','Acabamento simples e consumo.','{"budget_high":0.7,"offroad":0.9,"family":0.9,"space":1,"cargo":0.7,"performance":0.7,"automatic":1}'),
    ('Picape','Fiat','Strada','Freedom Cabine Dupla',2025,'Flex','Manual',120000,135000,5,844,11.6,13.3,'Picape compacta para trabalho e lazer.','Capacidade de carga e baixo custo operacional.','Conforto traseiro limitado.','{"budget_mid":0.9,"mixed":0.7,"offroad":0.7,"cargo":1,"performance":0.6,"manual":1}'),
    ('Picape','Ford','Maverick','Lariat Hybrid',2025,'Híbrido','Automático',235000,250000,5,943,15.7,13.6,'Picape híbrida confortável e potente.','Desempenho, economia urbana e caçamba.','Preço elevado e dimensões grandes.','{"budget_high":1,"urban":0.7,"mixed":0.9,"cargo":1,"performance":1,"efficiency":0.8,"comfort":0.8,"automatic":1}')
) AS v(category,brand,model,version,year,fuel,transmission,price_min,price_max,seats,trunk,city,highway,description,strengths,weaknesses,profile)
  ON c.name = v.category
WHERE NOT EXISTS (
    SELECT 1 FROM vehicles existing
    WHERE existing.brand = v.brand AND existing.model = v.model AND existing.year = v.year
);
