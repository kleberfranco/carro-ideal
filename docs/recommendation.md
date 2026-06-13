# Motor de recomendação

Cada opção do questionário contribui para dimensões como `urban`, `comfort`,
`efficiency` e `family`. A contribuição é multiplicada pelo peso da pergunta.

Para cada veículo, o motor percorre as dimensões do perfil, limita a aderência do
veículo a 1 e calcula:

`score = soma(peso_usuario * aderencia_veiculo) / soma(peso_usuario) * 100`

Os veículos são ordenados pelo score, com ID como desempate, e os dez primeiros
são persistidos. A resposta contém `reason` para leitura humana e
`matched_criteria` para futura geração de texto por IA. Sem dimensões compatíveis,
o score é zero e uma explicação geral é usada.

A complexidade é `O(V * D + V log V)`, onde `V` é o número de veículos e `D` o
número de dimensões. O catálogo ativo fica em cache TTL e o benchmark pode ser
executado com:

```sh
docker run --rm -v "$PWD":/app -w /app golang:1.22-alpine \
  go test -bench BenchmarkScoreVehicle -benchmem ./app/service
```
