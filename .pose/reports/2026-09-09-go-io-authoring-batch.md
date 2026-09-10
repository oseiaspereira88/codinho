# Lote de autoria go-io — 2026-09-09

Spec: go-foundations-packs, permanece in-progress.

Adicionados dois desafios draft: decode-strict-config (combined) e
import-csv-inventory (functional_slice), 10 conceitos, 6 competências e
12 nodes. Autoria codex; sem revisão humana ou playtest declarados.
A consulta padrão de produção continua excluindo este pack não publicado.

Validação: pose validate --strict passou 22/22. Após ajustar critérios
intermediários para avaliação qualitativa, go run ./cmd/codinho catalog
validate --checks --json passou: 40/40 checks declarados verificados no
catálogo completo. Os dois desafios novos falham na fixture inicial e passam
na referência privada. Os testes incluem entrada fragmentada, JSON extra,
campos desconhecidos, CSV com aspas e UTF-8, duplicações, overflow de int,
erros de leitura/escrita e ausência de escrita diante de entrada inválida.
Docs-check declarou 14 documentos, zero erros/warnings de inventário.

Revisão técnica do lote: sem código executado durante carga de catálogo;
execução das fixtures restrita ao comando de validação autorizado. Sem rede,
secrets ou paths externos nas fixtures. O check final não é imposto aos
micropassos intermediários, que exigem evidência qualitativa e rubrica.
Não houve avaliação pedagógica independente nem playtest humano neste lote.

Pendências preservadas: quatro variantes de I/O, cinco trilhas fundamentais,
ampliação de conceitos/competências/nodes e publicação humana. O inventário
em docs/catalog/go-foundations.md registra 42 desafios fundamentais atuais
(33 A / 8 C / 1 F) e a divergência entre os sete atômicos de go-first-steps
e os seis planejados. Nenhum conteúdo foi removido ou recategorizado.
