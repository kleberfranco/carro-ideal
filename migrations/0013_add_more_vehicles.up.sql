-- Expande o catálogo com novos e seminovos populares no Brasil, cobrindo gaps
-- de preço e categoria: hatches econômicos, SUVs médios e picapes robustas.

INSERT INTO vehicles (
    category_id, brand, model, version, year, fuel_type, transmission,
    price_min, price_max, seats, trunk_capacity, consumption_city,
    consumption_highway, description, strengths, weaknesses, condition, match_profile
)
SELECT c.id, v.brand, v.model, v.version, v.year, v.fuel, v.transmission,
       v.price_min, v.price_max, v.seats, v.trunk, v.city, v.highway,
       v.description, v.strengths, v.weaknesses, v.condition, v.profile::jsonb
FROM vehicle_categories c
JOIN (VALUES
    -- -------------------------------------------------------------------------
    -- NOVOS
    -- -------------------------------------------------------------------------

    -- Hatches novos
    ('Hatch','Fiat','Mobi','Like 1.0',2025,'Flex','Manual',65000,72000,4,200,13.5,14.8,
     'Ultra compacto para cidade, o menor custo do mercado.',
     'Preço imbatível, facilidade de estacionar e baixíssimo consumo.',
     'Espaço reduzido e conforto básico.',
     'novo','{"novo":1,"budget_low":1,"urban":1,"compact":1,"efficiency":1,"manual":1}'),

    ('Hatch','Fiat','Argo','Drive 1.3',2025,'Flex','Manual',83000,98000,5,300,12.6,15.4,
     'Hatch popular com bom espaço e acabamento superior ao segmento.',
     'Porta-malas generoso para a classe, acabamento e segurança.',
     'Câmbio manual pode afastar quem prefere automático.',
     'novo','{"novo":1,"budget_mid":0.7,"urban":0.9,"compact":0.9,"efficiency":0.8,"family":0.5,"manual":1}'),

    ('Hatch','Hyundai','HB20','Comfort Plus 1.0',2025,'Flex','Manual',87000,102000,5,295,12.9,15.7,
     'Hatch equilibrado com ótima avaliação de segurança.',
     '5 estrelas Latin NCAP, conectividade e baixo custo de manutenção.',
     'Desempenho modesto em alta velocidade.',
     'novo','{"novo":1,"budget_mid":0.8,"urban":1,"compact":0.9,"efficiency":0.8,"family":0.5,"manual":1}'),

    ('Hatch','Volkswagen','Polo','Comfortline 1.0 TSI',2025,'Flex','Automático',105000,122000,5,300,12.2,14.3,
     'Hatch premium com motor turbo e câmbio automático de dupla embreagem.',
     'Dirigibilidade, acabamento europeu e câmbio automático eficiente.',
     'Custo de manutenção acima da média do segmento.',
     'novo','{"novo":1,"budget_mid":1,"urban":0.9,"mixed":0.7,"compact":0.8,"comfort":0.8,"efficiency":0.7,"automatic":1}'),

    -- SUVs novos
    ('SUV','Fiat','Pulse','Impetus 1.3 Turbo',2025,'Flex','Automático',118000,145000,5,370,11.8,13.6,
     'SUV compacto moderno com motor turbo e design arrojado.',
     'Espaço interno, motor turbo ágil e equipamentos de série.',
     'Desempenho off-road limitado; feito para asfalto.',
     'novo','{"novo":1,"budget_mid":1,"budget_high":0.5,"urban":0.8,"mixed":0.9,"family":0.8,"space":0.7,"comfort":0.7,"automatic":1}'),

    ('SUV','Volkswagen','Nivus','Highline 1.0 TSI',2025,'Flex','Automático',142000,162000,5,415,12.0,13.4,
     'SUV cupê compacto com apelo esportivo e bom espaço interno.',
     'Estilo, desempenho do motor turbo e tecnologia embarcada.',
     'Entrada traseira limitada pelo teto em arco.',
     'novo','{"novo":1,"budget_high":0.8,"urban":0.8,"mixed":0.9,"family":0.7,"comfort":0.8,"performance":0.6,"automatic":1}'),

    ('SUV','Jeep','Renegade','Longitude 1.3 Turbo',2025,'Flex','Automático',148000,170000,5,320,10.9,12.5,
     'SUV compacto com capacidade off-road e personalidade robusta.',
     'Capacidade fora de estrada, segurança e identidade de marca.',
     'Consumo mais alto que concorrentes e porta-malas pequeno.',
     'novo','{"novo":1,"budget_high":0.9,"mixed":0.9,"offroad":1,"family":0.8,"space":0.7,"performance":0.7,"automatic":1}'),

    ('SUV','Jeep','Compass','Longitude 1.3 Turbo',2025,'Flex','Automático',195000,230000,5,438,10.7,12.1,
     'SUV médio premium com tecnologia avançada e enorme apelo no Brasil.',
     'Conforto, equipamentos, espaço interno e status de marca.',
     'Preço alto e custo de manutenção elevado.',
     'novo','{"novo":1,"budget_high":1,"mixed":1,"family":1,"comfort":1,"space":0.9,"performance":0.8,"automatic":1}'),

    -- -------------------------------------------------------------------------
    -- SEMINOVOS
    -- -------------------------------------------------------------------------

    -- Hatches seminovos
    ('Hatch','Fiat','Bravo','Essence 1.8',2013,'Flex','Manual',38000,48000,5,280,9.1,12.2,
     'Hatch médio esportivo com motor potente por preço muito acessível.',
     'Desempenho, espaço e custo de aquisição baixíssimo.',
     'Idade avançada e custo de manutenção acima da média.',
     'seminovo','{"seminovo":1,"budget_low":1,"urban":0.8,"mixed":0.6,"compact":0.7,"performance":0.8,"manual":1}'),

    ('Hatch','Hyundai','HB20','Comfort 1.0',2019,'Flex','Manual',58000,68000,5,290,11.5,14.9,
     'Hatch moderno e confiável com boas avaliações de segurança.',
     'Confiabilidade, baixo custo e boa revenda.',
     'Desempenho modesto comparado a turboalimentados.',
     'seminovo','{"seminovo":1,"budget_low":1,"urban":1,"compact":0.9,"efficiency":0.8,"manual":1}'),

    ('Hatch','Peugeot','208','Active 1.6',2019,'Flex','Manual',65000,76000,5,285,10.4,13.8,
     'Hatch estiloso com design europeu diferenciado.',
     'Design, conforto e custo acessível para a marca.',
     'Custo de manutenção e peças acima da média nacional.',
     'seminovo','{"seminovo":1,"budget_low":1,"urban":0.9,"compact":0.9,"comfort":0.6,"efficiency":0.6,"manual":1}'),

    ('Hatch','Volkswagen','Polo','Comfortline 1.0 TSI',2020,'Flex','Automático',72000,85000,5,300,11.0,13.5,
     'Hatch premium seminovo com turbo e câmbio automático por preço de médio.',
     'Tecnologia, motor eficiente e câmbio de dupla embreagem.',
     'Histórico de manutenção pode variar conforme uso anterior.',
     'seminovo','{"seminovo":1,"budget_low":0.8,"urban":0.9,"compact":0.8,"comfort":0.7,"efficiency":0.7,"automatic":1}'),

    -- SUVs seminovos
    ('SUV','Honda','HR-V','EX 1.8',2019,'Flex','CVT',95000,112000,5,440,10.9,13.7,
     'SUV compacto premium com excelente espaço interno e confiabilidade Honda.',
     'Bancos traseiros rebatíveis (Magic Seat), conforto e durabilidade.',
     'Motor aspirado sem tanto vigor em ultrapassagens.',
     'seminovo','{"seminovo":1,"budget_mid":1,"urban":0.8,"mixed":0.9,"family":1,"space":0.8,"comfort":0.8,"automatic":1}'),

    ('SUV','Jeep','Renegade','Longitude 1.8',2019,'Flex','Automático',92000,108000,5,320,9.5,11.8,
     'SUV compacto com traço aventureiro e boa capacidade off-road.',
     'Capacidade fora de estrada, segurança e carisma de marca.',
     'Consumo e espaço de porta-malas abaixo dos rivais.',
     'seminovo','{"seminovo":1,"budget_mid":1,"mixed":0.9,"offroad":0.9,"family":0.8,"space":0.6,"performance":0.6,"automatic":1}'),

    ('SUV','Jeep','Compass','Longitude 2.0',2018,'Flex','Automático',105000,125000,5,438,9.0,11.4,
     'SUV médio que domina as vendas do segmento no Brasil, seminovo.',
     'Espaço, conforto, reputação de marca e revenda forte.',
     'Consumo e custo de manutenção acima da média.',
     'seminovo','{"seminovo":1,"budget_mid":1,"mixed":1,"family":1,"comfort":0.9,"space":0.9,"performance":0.7,"automatic":1}'),

    -- Picape seminova
    ('Picape','Chevrolet','S10','LTZ 2.8 Diesel',2015,'Diesel','Automático',115000,135000,5,0,7.5,11.2,
     'Picape robusta de grande capacidade com motor diesel duradouro.',
     'Torque, capacidade de carga, durabilidade e tração 4x4.',
     'Consumo alto na cidade e custo de revisão elevado.',
     'seminovo','{"seminovo":1,"budget_mid":1,"offroad":1,"cargo":1,"mixed":0.8,"performance":0.9,"automatic":1}')

) AS v(category,brand,model,version,year,fuel,transmission,price_min,price_max,seats,trunk,city,highway,description,strengths,weaknesses,condition,profile)
  ON c.name = v.category
WHERE NOT EXISTS (
    SELECT 1 FROM vehicles existing
    WHERE existing.brand = v.brand AND existing.model = v.model AND existing.year = v.year
);
