CREATE TABLE cars (
                      id SERIAL PRIMARY KEY,
                      nome TEXT NOT NULL,
                      marca TEXT NOT NULL,
                      ano INT NOT NULL,
                      preco DECIMAL(12,2) NOT NULL,
                      categoria TEXT NOT NULL, -- SUV, Sedan, Hatch...
                      consumo_urbano DECIMAL(5,2),
                      consumo_rodoviario DECIMAL(5,2),
                      criado_em TIMESTAMP DEFAULT NOW()
);