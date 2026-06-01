CREATE TABLE preferences (
                             id SERIAL PRIMARY KEY,
                             user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             idade_min INT,
                             idade_max INT,
                             valor_max DECIMAL(12,2),
                             tipo_preferencia TEXT, -- exemplo: "SUV", "sedan", "hatch"
                             criado_em TIMESTAMP DEFAULT NOW()
);