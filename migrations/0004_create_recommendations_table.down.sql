CREATE TABLE recommendations (
                                 id SERIAL PRIMARY KEY,
                                 user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                 carro_id INT NOT NULL REFERENCES cars(id),
                                 score DECIMAL(5,2), -- relevância calculada pela IA
                                 criado_em TIMESTAMP DEFAULT NOW()
);