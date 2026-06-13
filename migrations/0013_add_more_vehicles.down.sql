DELETE FROM vehicles
WHERE (brand, model, year) IN (
    ('Fiat',       'Mobi',      2025),
    ('Fiat',       'Argo',      2025),
    ('Hyundai',    'HB20',      2025),
    ('Volkswagen', 'Polo',      2025),
    ('Fiat',       'Pulse',     2025),
    ('Volkswagen', 'Nivus',     2025),
    ('Jeep',       'Renegade',  2025),
    ('Jeep',       'Compass',   2025),
    ('Fiat',       'Bravo',     2013),
    ('Hyundai',    'HB20',      2019),
    ('Peugeot',    '208',       2019),
    ('Volkswagen', 'Polo',      2020),
    ('Honda',      'HR-V',      2019),
    ('Jeep',       'Renegade',  2019),
    ('Jeep',       'Compass',   2018),
    ('Chevrolet',  'S10',       2015)
);
