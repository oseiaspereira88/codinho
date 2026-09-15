---
title: Catálogo fundamental de Go
doc_type: reference
---

# Catálogo fundamental de Go — plano e auditoria de trilhas

Tarefas editoriais: `plan-concept-depth` (baseline de expansão),
`audit-tracks` (reconciliação dos percursos) e `audit-depth-remediation`
(reconciliação pós-remediação). Inventário auditado em 2026-09-12 UTC contra o
manifest atual, depois dos lotes conceituais e dos dois lotes de remediação
aprovados; a reconciliação de fechamento foi conferida em 2026-09-15 UTC.
Este documento é um plano de autoria e auditoria, não uma
declaração de publicação: a spec `go-foundations-packs` continua `in-progress`,
os sete packs continuam em rascunho e revisão humana/playtest continuam
pendentes.

## 1. Ledger reconciliado após a remediação

Os números abaixo vêm de `go run ./cmd/codinho catalog validate --json` e de
uma contagem independente das árvores `layers[*].macro_steps`, recursiva. O nó
de `layer` não entra na contagem. `Base` inclui o conteúdo histórico; entre
parênteses, quando necessário, aparece o crédito conservador usado para a
meta: o protótipo histórico `go-data.slice-filter-preserve-input` tem um
micro-node, mas esse node continua não creditado até sua disposição explícita.
O estado antes da remediação era 281 nodes brutos e 280 creditáveis; os dois
lotes aprovados chegaram a 305 brutos e 304 creditáveis. Uma correção posterior
removeu um micropasso redundante cuja assinatura já vinha na fixture; o estado
atual é 304 brutos e 303 creditáveis.

| Pack | Temas A–M | Entradas de base/variantes | Conceitos | Competências | Nodes de base M/Meso/Micro | Nodes de variantes M/Meso/Micro |
|---|---|---:|---:|---:|---:|---:|
| `go-first-steps` | A, B, com E/C introdutórios | 7 bases | 22 | 9 | 6/6/30 = 42 (crédito 6/6/29) | 0/0/0 |
| `go-core` | C, D, K | 10 bases | 15 | 10 | 5/5/38 = 48 | 0/0/0 |
| `go-data-text` | E, F | 10 bases | 13 | 8 | 12/12/40 = 64 | 0/0/0 |
| `go-type-design` | G, H, I | 8 bases | 12 | 6 | 7/7/34 = 48 | 0/0/0 |
| `go-errors` | J | 7 bases | 17 | 10 | 5/5/35 = 45 | 0/0/0 |
| `go-io` | L | 2 bases + 4 variantes | 14 | 11 | 5/5/19 = 29 | 4/4/16 = 24 |
| `go-testing` | M | 1 base + 4 variantes | 11 | 6 | 5/5/18 = 28 | 5/5/22 = 32 |
| **Sete packs fundamentais** | **A–M** | **45 bases + 8 variantes** | **104** | **60** | **45/45/214 = 304** (**crédito 303**) | **9/9/38 = 56** |
| `go-debugging` (fora do escopo) | protótipo | 1 protótipo | 1 | 4 | 1/1/4 = 6 | 0/0/0 |
| **Manifest completo** | — | **46 bases/protótipo + 8 variantes** | **105** | **64** | **46/46/218 = 310** | **9/9/38 = 56** |

Os sete packs contêm **45 entradas que não estão marcadas
`canonical: false`**; as oito entradas restantes são variantes de I/O/testes
marcadas `canonical: false`. A medição de 45 é inventário, não a meta
fundamental: uma dessas entradas é o protótipo histórico
`go-data.slice-filter-preserve-input`. Excluindo-o conceitualmente, sem
removê-lo do pack, a distribuição fundamental fica em **44 entradas: 32
atômicos, 10 combinados e 2 fatias funcionais**. O protótipo continua fora da
trilha `go-from-zero`, que usa o desafio canônico
`go-data-text.filter-without-mutating-input`.

O segundo protótipo, `go-debug.slice-off-by-one`, está no pack
`go-debugging`, fora dos sete packs fundamentais. Os IDs dos dois protótipos e
as sessões existentes permanecem preservados; esta reconciliação não renomeia,
remove nem migra esses registros. A classificação, eventual migração e
publicação dos protótipos, assim como o aceite humano da distribuição 32/10/2,
continuam pendentes; nenhum deles deve ser reclassificado ou excluído
silenciosamente para atingir a meta.

### Regra de crédito

Há três projeções que devem permanecer visíveis em toda auditoria:

1. **Escopo fundamental:** sete packs, 104 conceitos e 60 competências hoje;
   45 entradas não marcadas `canonical: false` permanecem inventariadas,
   sendo 44 fundamentais e um protótipo histórico, além das 8 variantes.
2. **Crédito de profundidade:** nodes de bases, sem as oito variantes e sem o
   micro-node do protótipo histórico até a decisão de disposição; o resultado
   atual é 304 brutos e 303 creditáveis.
3. **Manifesto:** inclui `go-debugging` e as variantes para fins de inventário,
   mas não transforma esses itens em conteúdo fundamental novo.

Os lotes conceituais já alcançaram **104 conceitos e 60 competências nos sete
packs**. As remediações adicionaram **+6/+6/+12** em macro/meso/micro,
movendo o escopo fundamental de **39/39/203 = 281 brutos** (**crédito 280**)
para **45/45/215 = 305 brutos** (**crédito 304**). Depois da remoção do
micropasso redundante, o estado atual é **45/45/214 = 304** (**crédito 303**);
a meta de 300 nodes creditáveis continua satisfeita. Com `go-debugging`
preservado, o manifesto tem 105 conceitos, 64 competências e 310 nodes de
base; os 56 nodes de variantes continuam reportados separadamente, totalizando
366 nodes no inventário global.

## 2. Mapa reconciliado dos 104 conceitos e lacunas residuais

O mapeamento usa um tema primário por conceito; alguns conceitos têm relação
secundária indicada na pendência. A enumeração soma exatamente os 104 conceitos
dos sete packs atuais, sem contar `go-debugging`. Os IDs que antes apareciam
como propostas foram reconciliados com os lotes conceituais concluídos.

| Tema | Pack(s) e conceitos reconciliados | Pendência residual |
|---|---|---|
| **A — Ambiente e tooling** | `go-first-steps`: `go-modules`, `go-version-directive`, `printf-verb-matching`, `go-toolchain-environment`, `gofmt-idempotent-source`, `go-build-package-selection` | Contexto efetivo, formatação idempotente e seleção de packages foram entregues. Workspaces, tags e cross-compile ficam para lote posterior, pois não há desafio fundamental correspondente. |
| **B — Declarações e tipos** | `go-first-steps`: `named-types`, `field-types`, `zero-values`, `short-variable-declaration`, `constant-declaration`, `iota-sequence`, `type-conversion`, `narrowing-conversion-overflow`, `named-return-values`, `variable-shadowing`, `block-scope`, `untyped-constant-defaulting`, `numeric-constant-representability` | O contraste entre default contextual, representabilidade em compilação e conversão em runtime está coberto; não foi criado sinônimo adicional. |
| **C — Controle de fluxo** | `go-first-steps`: `range-loop`; `go-core`: `defer-execution-order`, `labeled-break-scope`, `short-circuit-nil-guard`, `if-initializer-scope`, `for-clause-lifecycle`, `defer-argument-evaluation` | As três decisões ausentes do plano foram autoradas. Paralelismo e controle concorrente ficam fora deste catálogo fundamental. |
| **D — Funções e métodos** | `go-core`: `variadic-parameters`, `pointer-vs-value-receiver`, `panic-recover-boundary`, `closure-capture-lifetime`, `method-value-binding` | Captura independente e ligação de method value estão cobertas. `pointer-vs-value-receiver` continua reutilizado por relação em G; `panic-recover-boundary` também se relaciona com J, sem duplicação. |
| **E — Arrays, slices e maps** | `go-first-steps`: `slice-declaration`; `go-data-text`: `slice-backing-array-aliasing`, `map-comma-ok-lookup`, `set-membership-with-empty-struct-map`, `slice-preallocation-length-vs-capacity`, `array-value-copy`, `append-reallocation-boundary`, `map-delete-and-clear` | Cópia de array, realocação de `append` e remoção/limpeza de map foram entregues com critérios distintos. |
| **F — Texto e dados binários** | `go-data-text`: `utf8-rune-boundary`, `csv-field-quoting`, `strconv-error-handling`, `rune-aware-whitespace-preservation`, `strings-builder-output`, `regexp-submatch-contract` | Construção incremental e submatches de regexp completam a cobertura planejada sem transformar normalização em novo sinônimo. |
| **G — Structs e ponteiros** | `go-first-steps`: `structs`; `go-type-design`: `struct-shallow-copy-aliasing`, `nil-safe-pointer-receiver`, `struct-embedding-promotion`, `struct-tags-reflection`, `method-set-addressability` | Embedding/promotion, tags e addressability foram entregues. `pointer-vs-value-receiver` e `unexported-field-invariant` entram como relações cruzadas. |
| **H — Interfaces** | `go-type-design`: `typed-nil-interface-trap`, `safe-type-assertion`, `compile-time-interface-satisfaction`, `interface-composition-consumer` | A satisfação em compilação e a composição pelo consumidor mínimo completam a cobertura de interfaces. |
| **I — Generics** | `go-type-design`: `generic-comparable-constraint`, `generic-two-type-parameters`, `generic-underlying-type-constraint` | O contraste entre tipo definido e tipo subjacente com `~` foi entregue; não foi criado outro conceito de “comparável”. |
| **J — Erros** | `go-errors`: `error-wrapping-with-percent-w`, `errors-as-type-extraction`, `errors-join-aggregation`, `error-translation-at-boundary`, `error-classification-precedence`, `error-identity-not-message`, `error-chain-type-traversal`, `error-joined-tree-order`, `catalog-zero-value`, `catalog-validation-before-mutation`, `catalog-deterministic-list`, `catalog-value-snapshots`, `custom-error-type-contract`, `errors-is-custom-matcher`, `errors-as-custom-matcher`, `errors-unwrap-method-contract`, `nil-error-success-contract` | Tipos customizados e hooks `Is`/`As`/`Unwrap`, além do sucesso nil, foram entregues com exemplos e relações. |
| **K — Packages e APIs** | `go-core`: `unexported-field-invariant`, `constructor-returns-interface`, `functional-options-pattern`, `package-import-export-boundary` | A fronteira entre package produtor e consumidor externo foi entregue sem repetir o conceito de módulo. |
| **L — I/O e serialização** | `go-io`: `io-reader-stream`, `io-partial-reads`, `io-eof-boundary`, `io-json-decoder`, `io-json-strict-fields`, `io-csv-quoted-fields`, `io-writer-errors`, `io-validate-before-write`, `io-csv-record-width`, `io-json-encode`, `bufio-scanner-token-boundary`, `io-short-write-contract`, `io-limit-reader-boundary`, `json-null-versus-absent` | Buffer/token, short write, limite explícito de Reader e presença JSON foram entregues. |
| **M — Testes** | `go-testing`: `testing-table-cases`, `testing-named-subtests`, `testing-clock-injection`, `testing-counting-double`, `testing-time-boundaries`, `testing-instant-versus-zone`, `testing-input-isolation`, `testing-repeatability`, `testing-helper-context`, `testing-cleanup-lifecycle`, `testing-fuzz-invariant` | Helper, cleanup e invariantes de fuzz foram entregues; paralelismo, golden, benchmark e race detector continuam temas posteriores. |

## 3. Distribuição dos 34 conceitos autorados

Cada linha registra conteúdo entregue, não apenas um título: os conceitos têm
`content.example` em contexto diferente dos desafios, relação inicial no grafo,
critério observável e referência primária. Os exemplos são mini-módulos ou
testes Go executáveis com biblioteca padrão, Go mínimo 1.25, sem rede e sem
arquivos externos. Quando a decisão é de compilação ou tooling, o conteúdo
registra o comando determinístico de `go/types`, `go env`, `gofmt` ou `go build`.

### `go-first-steps` — 17 → 22 conceitos

| Conceito autorado (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `go-toolchain-environment` (`inspect-go-toolchain-context`) | Mini-módulo que registra `go env GOMOD GOOS GOARCH` e compara o build com `runtime.GOOS`/`runtime.GOARCH`. | `go-modules` — `relates_to` | O teste identifica o módulo e o alvo efetivos sem depender do diretório corrente ou alterar o ambiente. | [CMD], [MOD], [RUNTIME] |
| `gofmt-idempotent-source` | Arquivo Go deliberadamente desalinhado; `gofmt -d` após formatação deve produzir diff vazio e a segunda execução não pode mudar bytes. | `printf-verb-matching` — `relates_to` | Mede idempotência da ferramenta; não confunde formato com correção de verbos de `fmt`. | [EG], [CMD] |
| `go-build-package-selection` | Mini-módulo com library, `cmd/` e package interno; `go list ./...` e `go build ./...` devem enumerar e compilar os três contextos permitidos. | `go-modules` — `relates_to`; `package-import-export-boundary` — `deepens_into` | Verifica seleção/compilação de packages, não só a existência de `go.mod`. | [S], [CMD], [MOD] |
| `untyped-constant-defaulting` | Teste atribui a mesma constante não tipada a `int`, `float64`, `rune` e `string`, verificando tipo estático por compilação e valor observado. | `constant-declaration` — `relates_to`; `type-conversion` — `contrasts_with` | O resultado nasce do contexto de uso; não há conversão explícita para mascarar o default. | [S] |
| `numeric-constant-representability` | `go/types` verifica snippets: uma constante representável passa e uma que excede o tipo falha antes da execução. | `narrowing-conversion-overflow` — `contrasts_with`; `type-conversion` — `relates_to` | Distingue erro de compilação de wraparound em conversão executada. | [S], [TYPES] |

### `go-core` — 9 → 15 conceitos

| Conceito autorado (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `if-initializer-scope` | Teste de uma decisão com inicializador local e um snippet em `go/types` que falha ao usar o nome depois do `if`. | `block-scope` — `deepens_into`; `short-circuit-nil-guard` — `relates_to` | A variável do inicializador existe somente na condição e nos ramos do `if`. | [S] |
| `for-clause-lifecycle` | Tabela executa `init`, condição, corpo, `continue` e `post` com um trace esperado, sem usar `range`. | `range-loop` — `contrasts_with`; `labeled-break-scope` — `relates_to` | O trace prova a ordem do `post` inclusive quando há `continue`; não é apenas uma contagem de iterações. | [S] |
| `defer-argument-evaluation` | Função adia uma função com argumento escalar e outra closure; o trace verifica o valor capturado na inscrição e o valor lido na chamada. | `defer-execution-order` — `deepens_into` | Isola tempo de avaliação do argumento, sem repetir a ordem LIFO dos `defer` existentes. | [S] |
| `closure-capture-lifetime` (`control-closure-capture`) | Fábrica retorna duas closures independentes; chamadas intercaladas devem manter contadores separados e estado após o retorno da fábrica. | `block-scope` — `relates_to`; `named-return-values` — `relates_to` | Observa ambiente capturado e independência entre instâncias, não apenas uma função anônima que retorna valor. | [S], [EG] |
| `method-value-binding` | Teste cria method values antes e depois de alterar um receiver e compara o efeito para receiver por valor e por ponteiro. | `pointer-vs-value-receiver` — `deepens_into` | Verifica quando o receiver é ligado ao method value; não repete a escolha de receiver do desafio-base. | [S] |
| `package-import-export-boundary` (`define-package-api-boundary`) | Mini-módulo tem teste em package externo: somente identificadores exportados compilam, e o consumidor usa a menor interface pública necessária. | `constructor-returns-interface` — `relates_to`; `go-modules` — `relates_to` | A fronteira é observada por um consumidor de outro package; não é apenas campo não exportado no mesmo package. | [S], [MOD], [EG] |

### `go-data-text` — 8 → 13 conceitos

| Conceito autorado (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `array-value-copy` | Teste atribui `[3]int`, altera a cópia e contrasta com uma slice sobre o mesmo conteúdo. | `slice-backing-array-aliasing` — `contrasts_with`; `slice-declaration` — `relates_to` | A atribuição de array copia valores; a comparação não permite concluir que toda slice seja cópia. | [S] |
| `append-reallocation-boundary` | Slice com capacidade conhecida sofre `append` antes e depois de cruzar a capacidade; o teste observa quando a mutação deixa de compartilhar backing array. | `slice-backing-array-aliasing` — `deepens_into`; `slice-preallocation-length-vs-capacity` — `relates_to` | O critério separa aliasing enquanto há capacidade de realocação e independência depois dela. | [S] |
| `map-delete-and-clear` | Teste insere chaves, usa `delete` para uma só, `clear` para todas e consulta mapa nil sem panic. | `map-comma-ok-lookup` — `deepens_into`; `set-membership-with-empty-struct-map` — `relates_to` | Verifica remoção seletiva e limpeza total, não apenas presença/ausência por comma-ok. | [S], [MAPS] |
| `strings-builder-output` | `strings.Builder` recebe fragmentos, `String` é comparada com a concatenação esperada e o caso vazio é verificado. | `rune-aware-whitespace-preservation` — `relates_to`; `utf8-rune-boundary` — `relates_to` | O contrato é acumulação e saída, não normalização de espaços nem corte de runes. | [STR] |
| `regexp-submatch-contract` | Teste usa uma regex compilada em linhas de inventário e verifica match completo, submatch e ausência (`nil`) separadamente. | `strconv-error-handling` — `relates_to`; `rune-aware-whitespace-preservation` — `relates_to` | Distingue ausência de match e grupo vazio sem transformar regex em mais um conceito de parsing numérico. | [REGEXP] |

### `go-type-design` — 6 → 12 conceitos

| Conceito autorado (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `struct-embedding-promotion` | Struct externa embute uma struct de endereço; teste acessa campo/método promovido e registra a ambiguidade quando dois campos têm o mesmo nome. | `struct-shallow-copy-aliasing` — `relates_to`; `structs` — `deepens_into` | O foco é seleção/promoção de membros, não cópia independente de slices. | [S] |
| `struct-tags-reflection` | `reflect.Type.FieldByName` e `StructTag.Get` leem uma tag `json` de um registro de catálogo. | `structs` — `relates_to`; `field-types` — `relates_to` | A tag altera metadado de serialização observado por reflexão; não altera nome nem tipo do campo. | [S], [REFLECT] |
| `method-set-addressability` | Teste compila chamada de método de ponteiro em valor endereçável e usa `go/types` para rejeitar a mesma chamada em elemento não endereçável de map. | `pointer-vs-value-receiver` — `deepens_into`; `nil-safe-pointer-receiver` — `relates_to` | O critério é addressability/method set, não somente o comportamento de receiver nil. | [S], [TYPES] |
| `compile-time-interface-satisfaction` | Atribuição de uma implementação a uma interface compila; um snippet mutante sem método é rejeitado por `go/types`. | `safe-type-assertion` — `contrasts_with`; `typed-nil-interface-trap` — `relates_to` | A falha é descoberta na compilação, antes de qualquer type assertion ou valor nil. | [S], [TYPES] |
| `interface-composition-consumer` | Package consumidor aceita uma interface de um método; um fake mínimo passa no teste sem conhecer o tipo concreto do produtor. | `constructor-returns-interface` — `deepens_into`; `safe-type-assertion` — `relates_to` | Mede dependência no menor contrato do consumidor, não só satisfação acidental de uma interface grande. | [S], [EG] |
| `generic-underlying-type-constraint` | Função com constraint `~int` aceita `type Score int` e rejeita `string`; o teste instancia ambos e mantém o erro no caso inválido. | `generic-comparable-constraint` — `deepens_into`; `generic-two-type-parameters` — `relates_to` | O critério é aceitar tipos definidos pelo tipo subjacente, sem duplicar a constraint `comparable`. | [S] |

### `go-errors` — 12 → 17 conceitos

| Conceito autorado (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `custom-error-type-contract` | Tipo `FieldError` carrega campo e valor inválido; teste verifica `Error` e extrai os dados com `errors.As`. | `errors-as-type-extraction` — `deepens_into`; `error-identity-not-message` — `relates_to` | O contrato do tipo inclui dados estruturados; a mensagem não é usada como protocolo. | [ERRORS] |
| `errors-is-custom-matcher` | Erro wrapper implementa `Is` para uma categoria estável; duas instâncias com mensagens diferentes devem satisfazer `errors.Is` somente pela categoria declarada. | `error-identity-not-message` — `deepens_into`; `error-classification-precedence` — `relates_to` | A equivalência semântica é explícita no método `Is`, sem comparar `Error()`. | [ERRORS] |
| `errors-as-custom-matcher` | Wrapper implementa `As` para expor apenas um tipo público; o teste verifica match do tipo permitido e ausência de match de um tipo interno. | `errors-as-type-extraction` — `deepens_into`; `error-translation-at-boundary` — `relates_to` | O limite do tipo exposto é controlado pelo método `As`, não pela profundidade casual da cadeia. | [ERRORS] |
| `errors-unwrap-method-contract` | Erro customizado implementa `Unwrap`; `errors.Is` atravessa uma causa, enquanto `Unwrap` nil encerra a cadeia. | `error-chain-type-traversal` — `deepens_into`; `errors-join-aggregation` — `relates_to` | Testa contrato de uma causa linear e seu término, sem recontar a ordem de árvore de `Join`. | [ERRORS] |
| `nil-error-success-contract` | Função de validação retorna `(value, nil)` em sucesso e um erro não nil em falha; o teste inclui entrada vazia e sucesso parcial. | `error-translation-at-boundary` — `relates_to`; `errors-join-aggregation` — `relates_to` | O critério é a identidade nil do caminho de sucesso, não a aparência da mensagem nem o typed nil de interface. | [ERRORS], [S] |

### `go-io` — 10 → 14 conceitos

| Conceito autorado (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `bufio-scanner-token-boundary` | Scanner lê linhas curtas e uma entrada que excede o buffer configurado; o teste distingue tokens lidos de `ErrTooLong`. | `io-reader-stream` — `deepens_into`; `io-partial-reads` — `relates_to` | A fronteira é do token/buffer, não uma leitura parcial genérica ou EOF. | [BUFIO], [IO] |
| `io-short-write-contract` | Writer de teste aceita somente um prefixo; a operação verifica `n < len(p)` e preserva `io.ErrShortWrite`/erro declarado. | `io-writer-errors` — `deepens_into`; `io-validate-before-write` — `relates_to` | O caso cobre escrita parcial e sua sinalização, não somente Writer que rejeita tudo. | [IO] |
| `io-limit-reader-boundary` | `io.LimitReader` recebe mais bytes que o limite; o teste lê exatamente o prefixo e confirma que o limite termina a leitura observável. | `io-reader-stream` — `relates_to`; `io-eof-boundary` — `deepens_into` | O limite imposto pelo consumidor é separado do EOF originado pela fonte. | [IO] |
| `json-null-versus-absent` | Decodifica três documentos usando presença em `map[string]json.RawMessage` e valor opcional: campo ausente, `null` e zero explícito produzem estados distinguíveis. | `io-json-decoder` — `deepens_into`; `io-json-strict-fields` — `relates_to` | Verifica semântica de presença/nulidade, não apenas rejeição de campo desconhecido. | [JSON] |

### `go-testing` — 8 → 11 conceitos

| Conceito autorado (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `testing-helper-context` | Helper chama `t.Helper()` e falha intencionalmente; o teste verifica que a linha reportada aponta para o chamador e que o nome do caso permanece legível. | `testing-named-subtests` — `deepens_into`; `testing-repeatability` — `relates_to` | O benefício é atribuição de diagnóstico, não apenas nomear subtestes. | [TEST] |
| `testing-cleanup-lifecycle` (`control-test-cleanup-lifecycle`) | Subteste registra dois `t.Cleanup` e o pai registra outro; o trace esperado demonstra ordem LIFO e escopo de execução. | `testing-repeatability` — `relates_to`; `testing-counting-double` — `relates_to` | O critério é lifecycle garantido pelo framework, sem `defer` dentro da função sob teste. | [TEST] |
| `testing-fuzz-invariant` | Fuzz target usa uma propriedade simples de normalização; seed conhecido passa e um mutante que viola a propriedade falha. | `testing-table-cases` — `deepens_into`; `testing-repeatability` — `relates_to` | O teste declara uma propriedade para entradas geradas, não apenas uma lista maior de exemplos. | [TEST], [FUZZ] |

### Referências oficiais

As referências são primárias e devem ser mantidas no conteúdo dos conceitos,
sem copiar a solução para briefing ou pista. `[S]` é a [especificação da
linguagem Go](https://go.dev/ref/spec); `[EG]` é [Effective
Go](https://go.dev/doc/effective_go); `[MOD]` é a [referência de
go.mod](https://go.dev/doc/modules/gomod-ref); `[CMD]` é a documentação do
[comando go](https://pkg.go.dev/cmd/go); `[TYPES]` é
[go/types](https://pkg.go.dev/go/types); `[RUNTIME]` é
[runtime](https://pkg.go.dev/runtime); `[STR]` é
[strings](https://pkg.go.dev/strings); `[REGEXP]` é
[regexp](https://pkg.go.dev/regexp); `[REFLECT]` é
[reflect](https://pkg.go.dev/reflect); `[MAPS]` é
[maps](https://pkg.go.dev/maps); `[ERRORS]` é o pacote
[errors](https://pkg.go.dev/errors); `[IO]` é
[io](https://pkg.go.dev/io); `[BUFIO]` é
[bufio](https://pkg.go.dev/bufio); `[JSON]` é
[encoding/json](https://pkg.go.dev/encoding/json); `[TEST]` é
[testing](https://pkg.go.dev/testing); e `[FUZZ]` é o [tutorial oficial de
fuzzing](https://go.dev/doc/tutorial/fuzz).

## 4. Quatro competências fundamentais adicionadas

O baseline anterior era 56. Os lotes adicionaram exatamente as quatro
competências abaixo, chegando a 60 sem contar as quatro competências de
`go-debugging` e sem promover variantes. As competências existentes continuam
reutilizáveis quando já expressam o comportamento; novo conceito não implica
nova competência por sinônimo.

| Nova competência | Pack/tema | Evidência mínima | Por que não é duplicata |
|---|---|---|---|
| `inspect-go-toolchain-context` — inspecionar o contexto efetivo da toolchain | `go-first-steps`, A | Registrar `go env`/build e explicar módulo, alvo e resultado reproduzível. | `initialize-a-go-module` cria o módulo; esta competência diagnostica o contexto em que o comando realmente roda. |
| `control-closure-capture` — controlar estado capturado por closures | `go-core`, D | Duas closures produzidas pela mesma fábrica preservam estados independentes após o retorno. | Nenhuma competência atual cobre tempo de vida/captura; `use-variadic-parameters` e retornos nomeados tratam outra decisão. |
| `define-package-api-boundary` — definir uma fronteira pública mínima entre packages | `go-core`, K | Consumidor externo compila usando apenas identificadores exportados e o menor contrato de interface. | `hide-concrete-type-behind-constructor` esconde a implementação dentro de uma API; esta competência verifica a fronteira entre packages. |
| `control-test-cleanup-lifecycle` — controlar cleanup no lifecycle de um teste | `go-testing`, M | Dois `t.Cleanup` no subteste e um no pai deixam observável a ordem LIFO e o escopo de execução. | As competências atuais cobrem tempo, fronteiras, aliasing, doubles e nomes, mas não cleanup garantido pelo framework; `t.Helper` permanece coberto por `testing-name-diagnostic-cases`. |

Distribuição reconciliada: `go-first-steps` 9, `go-core` 10,
`go-data-text` 8, `go-type-design` 6, `go-errors` 10, `go-io` 11 e
`go-testing` 6, totalizando 60. `go-debugging` permanece em 4 e aparece
somente na projeção global.

## 5. Nodes e decisões pedagógicas reconciliados

O quadro abaixo mede somente árvores de bases, excluindo variantes e não
creditando o micro-node do protótipo histórico. A coluna anterior registra o
estado imediatamente antes dos dois lotes aprovados; a coluna posterior é o
estado carregado no manifest atual. As adições são números de nodes
`macro/meso/micro`; cada micro autorado manteve uma intenção, um alvo, uma
evidência e um critério observável.

| Pack | Antes da remediação M/Meso/Micro | Depois da remediação M/Meso/Micro | Adição | Decisões cobertas |
|---|---:|---:|---:|---|
| `go-first-steps` | 6/6/30 | 6/6/30 | 0/0/0 | Contexto da toolchain, formatação/build, default de constantes e representabilidade já entregues; o micro histórico de filtro não é creditado. |
| `go-core` | 0/0/33 | 5/5/38 | +5/+5/+5 | Escopo de inicializador, ciclo de `for`, avaliação de argumentos de `defer`, captura de closure e ligação de method value; a fronteira com package consumidor já estava coberta pelo conceito e exemplo existentes. |
| `go-data-text` | 12/12/40 | 12/12/40 | 0/0/0 | Cópia de array, realocação de `append`, remoção/limpeza de map, builder e contrato de submatch já entregues. |
| `go-type-design` | 7/7/34 | 7/7/34 | 0/0/0 | Embedding/promotion, tags refletidas, addressability/method sets, satisfação em compilação, interface mínima e constraint com `~` já entregues. |
| `go-errors` | 5/5/35 | 5/5/35 | 0/0/0 | Tipo de erro com dados, hooks customizados `Is`/`As`/`Unwrap` e sucesso nil já entregues. |
| `go-io` | 5/5/19 | 5/5/19 | 0/0/0 | Buffer/token, short write, limite explícito de Reader e ausência/`null`/zero em JSON já entregues; variantes não contam. |
| `go-testing` | 4/4/11 | 5/5/18 | +1/+1/+7 | Nomes de subtestes, diagnóstico de helper, dono e momento de cleanup, seeds de fronteira, entrada gerada e estabilidade da propriedade de fuzz. Variantes não contam. |
| **Total bruto** | **39/39/203 = 281** | **45/45/215 = 305** | **+6/+6/+12 = +24** | O micro-node do protótipo histórico permanece inventariado, mas não creditado. |
| **Total creditável** | **39/39/202 = 280** | **45/45/214 = 304** | **+6/+6/+12 = +24** | A meta de 300 e o alvo operacional de 304 estão satisfeitos; variantes e protótipos permanecem separados. |

Esse quadro preserva o snapshot da auditoria de profundidade em 2026-09-12.
Depois dele, a correção de `go-first-steps.clamp-int-to-byte` removeu o node
redundante de assinatura que já estava materializado na fixture. A leitura
corrente é 304 nodes brutos e 303 creditáveis; a meta de 300 permanece
satisfeita.

### Como os nodes serão autorados

- Um macro/meso pode ser uma decisão de orientação, mas cada micro será uma
  intenção única. Não usar “declarar e implementar” ou “testar e explicar” no
  mesmo micro; separar análise, alteração, verificação e reflexão quando forem
  decisões cognitivas distintas.
- Cada micro que exige avaliação positiva terá hints 1–6 realmente
  progressivos. Hints 1–5 não conterão `kind: solution`; nenhuma instrução,
  nome de fixture ou comentário revelará a implementação antes do momento.
- Cada critério aponta para `source_inspection`, `compile`, `test` ou comando
  de tooling executável. Reflexões perguntarão por que o comportamento ocorre,
  sobretudo em aliasing, escopo, nil, ordem e fronteiras.
- Relações novas usarão apenas `relates_to`, `deepens_into`,
  `contrasts_with` ou `applies_in` quando a direção for justificável. Relações
  de contraste serão recíprocas; `relation_refs` apontará para arestas que
  também existirão no grafo. Não criar um conceito só para preencher uma aresta.
- A projeção original de 108 nodes adicionais foi distribuída nas bases já
  existentes e nas duas fatias funcionais, sem novos desafios canônicos. Os
  dois lotes aprovados entregaram os 24 nodes que faltavam desde o baseline
  anterior, sem alterar a decisão histórica de 44 entradas fundamentais, cinco
  trilhas e oito variantes.

## 6. Auditoria das cinco trilhas (`audit-tracks`)

As cinco trilhas vivem em `packs/go-first-steps.yaml` versão 1.9.0 e reutilizam
desafios dos sete packs. Elas têm 47 posições e 37 IDs de desafio distintos;
repetição entre percursos é intencional e não altera a contagem do catálogo.
As decisões abaixo explicam o recorte e a ordem, sem transformar uma sequência
de navegação em promessa de nível ou publicação.

### `go-from-zero` — Go do zero (11)

Seleciona o primeiro contato com módulo, declarações, conversões, tooling,
escopo, coleções, funções, mapas e erros. A ordem vai do contexto que permite
compilar ao capstone de estado e erros: cada bloco acrescenta uma decisão que o
seguinte reutiliza. O protótipo histórico de filtro permanece preservado no
pack, mas foi substituído aqui pelo desafio canônico com fixture e check.

Ordem autorada:

`go-first-steps.declare-a-minimal-module` →
`go-first-steps.enumerate-weekdays-with-iota` →
`go-first-steps.convert-celsius-to-fahrenheit` →
`go-first-steps.clamp-int-to-byte` →
`go-first-steps.format-rate-limit-message` →
`go-first-steps.avoid-shadowing-named-returns` →
`go-data-text.filter-without-mutating-input` →
`go-core.sum-variadic-numbers` →
`go-data-text.lookup-map-value-with-comma-ok` →
`go-errors.wrap-sentinel-error-with-context` →
`go-errors.maintain-memory-catalog`.

### `go-oop-transition` — Transição de orientação a objetos (10)

Começa com estado mutável em receiver de ponteiro, fecha o invariante e reduz
a API por interface antes de tratar nil, cópia independente e assertion. O
desafio de opções funcionais amplia a composição; o relógio injetado traz
testabilidade; o catálogo em memória fecha o percurso com estado e erros.

Ordem autorada:

`go-core.increment-counter-with-pointer-receiver` →
`go-core.guard-invariant-with-unexported-field` →
`go-core.hide-counter-type-behind-interface` →
`go-type-design.select-validator-true-nil` →
`go-type-design.sum-list-with-nil-safe-receiver` →
`go-type-design.clone-inventory-independently` →
`go-type-design.first-circle-clone` →
`go-core.configure-server-with-functional-options` →
`go-testing.select-active-tokens-with-injected-clock` →
`go-errors.maintain-memory-catalog`.

### `go-practical-fluency` — Fluência prática (9)

Seleciona situações frequentes de código de aplicação. A sequência parte de
controle de fluxo e dados ordenados, passa por parsing e pela escalada
`errors.Is`/`errors.As`/classificação, e termina com I/O JSON/CSV antes do
catálogo funcional. Assim, os desafios integrados aparecem depois de seus
átomos de controle, dados e erros.

Ordem autorada:

`go-core.find-first-value-at-least` →
`go-data-text.dedupe-preserving-first-occurrence` →
`go-data-text.sum-valid-integers-safely` →
`go-errors.wrap-sentinel-error-with-context` →
`go-errors.extract-field-with-errors-as` →
`go-errors.classify-wrapped-errors` →
`go-io.decode-strict-config` →
`go-io.import-csv-inventory` →
`go-errors.maintain-memory-catalog`.

### `go-data-and-types` — Dados e tipos (10)

Avança dos invariantes de slice e capacidade para fronteiras de UTF-8 e
espaçamento, compõe a formatação CSV, depois retoma cópia de structs e chega a
transformações genéricas e receiver nil-safe. A seleção evita usar generics
antes de estabelecer ordem, alocação e isolamento de dados.

Ordem autorada:

`go-data-text.filter-without-mutating-input` →
`go-data-text.preallocate-slice-with-zero-length` →
`go-data-text.truncate-bytes-at-rune-boundary` →
`go-data-text.title-first-letters-preserving-spacing` →
`go-data-text.quote-csv-field-when-needed` →
`go-data-text.format-names-as-csv-fields` →
`go-type-design.clone-inventory-independently` →
`go-type-design.contains-generic-comparable` →
`go-type-design.map-generic-transform` →
`go-type-design.generic-node-values-nil-safe`.

### `go-testing-and-design` — Testes e design (7)

Coloca isolamento de dados e encapsulamento antes do desafio de teste
determinístico. Depois, usa identidade, agregação e classificação de erros
para mostrar que o design testável inclui contratos de falha, fechando com o
catálogo que combina invariantes, estado e visões isoladas.

Ordem autorada:

`go-type-design.clone-inventory-independently` →
`go-core.guard-invariant-with-unexported-field` →
`go-testing.select-active-tokens-with-injected-clock` →
`go-errors.wrap-sentinel-error-with-context` →
`go-errors.validate-user-joining-all-errors` →
`go-errors.classify-wrapped-errors` →
`go-errors.maintain-memory-catalog`.

### Prova de pré-requisitos e alcance

O teste de integração carrega a cópia do catálogo real com `curriculum.Load`,
compara a membership e a ordem exatas acima e chama `ResolvePath` para cada
percurso. Ele também verifica que cada referência de `prerequisites` existe e,
quando aponta para outro desafio, aparece antes dele na mesma trilha. Os
pré-requisitos conceituais continuam referências de conceito — não são
silenciosamente convertidos em desafios extras.

Em seguida, o mesmo teste inicia cada ID via `session_start` no servidor MCP
real sobre stdio, consulta `session_get` e `instruction_get` em todos os
cursores e usa `step_advance` com `override: true` somente para atravessar a
sequência inteira. A comparação de IDs, cursor, total e proveniência `draft`
prova alcance e roteamento. Esse avanço automatizado não resolve exercícios,
não avalia a ordem pedagógica e não é playtest humano; revisão humana e
playtest ponta a ponta continuam pendentes pelo checklist.

## 7. Sequência de execução e gates

1. **Congelar o ledger:** manter o resultado pós-remediação por pack e por
   `macro/meso/micro`; manter uma linha separada para o protótipo histórico e
   outra para as oito variantes. A disposição dos dois protótipos continua uma
   decisão humana da spec, sem quebrar IDs de sessões.
2. **Conceitos e competências:** os 34 conceitos e quatro competências das
   seções 3 e 4 estão reconciliados; manter exemplos, relações e referências
   ligados ao conteúdo real ao autorar os nodes restantes.
3. **Completar árvores de bases:** a matriz da seção 5 foi aplicada; após a
   correção do micropasso redundante, a recontagem alcança 303 nodes
   creditáveis no conjunto. Variantes
   podem receber revisão de equivalência, mas não podem reduzir o trabalho das
   bases nem pagar a meta.
4. **Verificar:** executar, no estado atual,
   `go run ./cmd/codinho catalog validate --json` e
   `go run ./cmd/codinho catalog validate --checks --json`. Registrar
   `declared == verified` para checks e separar o `--v1-gate` de publicação dos
   gates de conteúdo draft. Rodar também `go test ./...` e os exemplos
   isolados; não confundir check aprovado com revisão pedagógica.
5. **Pré-revisão e humano:** submeter o lote à pré-revisão automatizada
   independente e depois ao revisor humano. Fazer playtest ponta a ponta com
   `workspace prepare` e sessão real antes de qualquer publicação.

## 8. Decisões históricas e pendências preservadas

- A árvore é reutilizada por granularidade e por trilhas; não criar packs por
  senioridade nem duplicar desafios para obter mais nodes.
- O conteúdo histórico mantém IDs, autores, fixtures e sessões. Em particular,
  `go-data.slice-filter-preserve-input` não será apagado por estar além da
  distribuição prevista, e `go-debug.slice-off-by-one` permanece no pack
  `go-debugging`, fora dos sete packs fundamentais.
- As oito variantes continuam contextualizadas, mas com `variant_of` e
  `canonical: false`; não serão contadas como novos desafios, competências ou
  profundidade das bases.
- A classificação, eventual migração e publicação dos dois protótipos, assim
  como o aceite humano da distribuição 32/10/2, continuam pendentes. Até essa
  decisão, nenhum protótipo será reclassificado ou excluído silenciosamente
  para atingir a meta.
- Nenhum `publication.reviewed_by`, `playtested: true` ou
  `status: published` será preenchido por autoria automatizada. O checklist
  exige pessoa diferente do autor e playtest honesto; ambos estão pendentes.
- A spec continua `in-progress` até a disposição dos protótipos, a revisão de
  duplicações/leaks e os gates humanos; a confirmação de 303 nodes creditáveis
  no ledger, 104 conceitos e 60 competências já está reconciliada.
  O `--v1-gate` global é uma etapa posterior e seus limiares maiores não devem
  ser alegados como satisfeitos por este plano.

## 9. Auditoria da remediação de profundidade (`audit-depth-remediation`)

Esta seção registra a evidência desta tarefa no estado medido em 2026-09-12
UTC. Os únicos paths autorizados e editados neste lote foram este documento e
`internal/curriculum/foundation_variants_test.go`; não houve commit, integração,
alteração de estado do coordenador ou promoção editorial. A árvore de entrada
das remediações é `e7ba4533a31b0806b62ec7af89839d72a47ade7d`; o estado após os
dois lotes aprovados está em `135f532a47e47c66d93a54726ba5e157bc9df9d4`
(HEAD). O commit de remediação contém os
dois lotes no mesmo batch editorial; esta auditoria não o altera.

### 9.0 Rastreabilidade da auditoria global anterior

A auditoria global `audit-global-gates` foi integrada no commit
`655e9611278fb2ce67fd4de0f1d7b1d7e3a6091a` e registrou o consumo de seus 61
dependentes: 52 auditorias de desafios, 7 lotes conceituais e 2 reconciliações.
Esta seção não duplica aquele ledger; preserva sua proveniência e o atualiza
com o commit `135f532a47e47c66d93a54726ba5e157bc9df9d4`, que entregou as duas
remediações de profundidade e registrou, naquele snapshot, os números
367/305/304.

### 9.1 Lotes aprovados, versões, contagens e decisões

Os dois lotes de remediação aprovados foram consumidos. O commit é único porque
o batch editorial aplicou os dois paths de pack juntos; cada linha mantém o ID
da tarefa, a versão observada e a contagem independente antes/depois:

| Tarefa | Commit aprovado | Pack/versão observada | Antes M/Meso/Micro | Depois M/Meso/Micro | Decisões cobertas |
|---|---|---|---:|---:|---|
| `remediate-depth-go-core` | `135f532a47e47c66d93a54726ba5e157bc9df9d4` | `go-core` 1.4.0 | 0/0/33 | 5/5/38 (+5/+5/+5) | Escopo de inicializador, ciclo de `for`, avaliação de argumentos de `defer`, captura de closure e ligação de method value; a fronteira com package consumidor permaneceu coberta pelo conceito/exemplo existente. |
| `remediate-depth-go-testing` | `135f532a47e47c66d93a54726ba5e157bc9df9d4` | `go-testing` 1.3.0 | 4/4/11 | 5/5/18 (+1/+1/+7) | Nomes de subtestes, diagnóstico de helper, registro e ordem de cleanup, seeds de fronteira, entradas geradas e estabilidade da propriedade de fuzz. |

O delta combinado é `+6/+6/+12`: o escopo fundamental passou de 281 nodes
brutos/280 creditáveis para 305/304. Nenhum lote criou desafio, variante ou
protótipo novo; ambos preservaram IDs, fixtures, relações e o status draft.

As versões carregadas no manifest são:

| Pack | Versão | Escopo nesta auditoria |
|---|---:|---|
| `go-first-steps` | 1.9.0 | fundamental; também hospeda as cinco trilhas |
| `go-core` | 1.4.0 | fundamental |
| `go-data-text` | 1.3.0 | fundamental |
| `go-errors` | 1.5.0 | fundamental |
| `go-io` | 1.3.0 | fundamental; quatro variantes |
| `go-testing` | 1.3.0 | fundamental; quatro variantes |
| `go-type-design` | 1.2.0 | fundamental |
| `go-debugging` | 1.2.0 | fora do escopo; um protótipo preservado |

O schema do manifest é 1.0.0. As versões dos packs e os IDs dos protótipos
foram auditados no estado pós-remediação; os lotes não alteraram os números de
versão nem renomearam, removeram ou promoveram registros.

### 9.2 Ledger global, fundamental, variantes e publicação

`catalog validate --json` reporta a projeção completa abaixo. A segunda linha
separa os sete packs fundamentais; a terceira mantém o pack de debugging fora
do cálculo de profundidade fundamental.

| Projeção | Conceitos | Competências | Trilhas | Desafios | Nodes |
|---|---:|---:|---:|---:|---:|
| Manifest completo | 105 | 64 | 5 | 54 | 367 |
| Sete packs fundamentais | 104 | 60 | 5 reutilizadas | 45 bases + 8 variantes | 305 bases brutos; 304 creditáveis |
| `go-debugging` fora do escopo | 1 | 4 | — | 1 protótipo | 6 |
| Variantes fundamentais | — | — | — | 8 (`canonical: false`) | 56 (`9/9/38`) |
| Publicados | 0 | 0 | 0 | 0 | 0 |
| Elegíveis para publicação | 0 | 0 | 0 | 0 | 0 |

O inventário global por tipo é 41 atômicos, 10 combinados, 1 debug e 2
`functional_slice`. Nos sete packs, as 45 bases incluem 33 atômicos, 10
combinados e 2 fatias funcionais; removendo conceitualmente o protótipo
histórico `go-data.slice-filter-preserve-input`, o plano fundamental fica em
44 entradas: 32 atômicos, 10 combinados e 2 fatias funcionais. As oito
variantes continuam fora da contagem de desafios novos e de profundidade; os
305 nodes brutos incluem o micro histórico, e os 304 creditáveis não.

Os protótipos estão nos escopos corretos e permanecem distintos:

- `go-data.slice-filter-preserve-input` está em `go-first-steps`, é a entrada
  histórica com um micro-node não creditado e não participa de `go-from-zero`;
- `go-debug.slice-off-by-one` está em `go-debugging`, tem 6 nodes e não é
  contado na cobertura fundamental;
- nenhuma das duas entradas é variante, e a classificação, disposição e
  publicação de ambas continua pendente de decisão humana.

### 9.3 Gates quantitativos

| Gate da spec `go-foundations-packs` | Resultado | Evidência |
|---|---|---|
| Pelo menos 100 conceitos fundamentais | PASS — 104 | `TestFoundationCatalogAudit` |
| Pelo menos 60 competências fundamentais | PASS — 60 | `TestFoundationCatalogAudit` |
| Pelo menos 300 nodes contextualizados creditáveis | PASS — 304 | `TestFoundationCatalogAudit`; bases `45/45/214`, sem variantes e sem o micro histórico |

O ledger de planejamento da seção 5 aponta o alvo operacional de 304 nodes de
bases creditáveis; após os dois lotes, tanto esse alvo quanto o limiar de 300
da spec estão satisfeitos. O teste mantém a separação estrutural: variantes
somam 56 nodes apenas na projeção própria, o micro do protótipo histórico não
entra nos 304, e o protótipo de `go-debugging` permanece fora do escopo.

### 9.4 Auditoria de cobertura técnica e contextos essenciais

Os seguintes contextos foram auditados nos briefs, critérios, fixtures,
referências privadas e conteúdo conceitual. O validador de checks executou os
checks declarados; o teste novo também impede que variantes ou o protótipo
histórico paguem a meta de bases.

| Tema auditado | Contextos representativos |
|---|---|
| `nil` e zero values | `go-core.find-first-value-at-least`, `go-type-design.sum-list-with-nil-safe-receiver`, `go-type-design.generic-node-values-nil-safe`, `go-data-text.lookup-map-value-with-comma-ok` |
| Aliasing e isolamento | `go-data-text.filter-without-mutating-input`, `go-type-design.clone-inventory-independently`, `go-type-design.first-circle-clone`, `go-testing.snapshot-readings` |
| UTF-8 e fronteiras de rune | `go-data-text.truncate-bytes-at-rune-boundary`, `go-data-text.title-first-letters-preserving-spacing`, `go-io.read-fragmented-note`, `go-io.read-quoted-attendees` |
| Method sets e addressability | `go-core.increment-counter-with-pointer-receiver`, `go-core.hide-counter-type-behind-interface`, conceito `method-set-addressability` |
| Interface nil/typed nil | `go-type-design.select-validator-true-nil`, `go-type-design.first-circle-clone`, `go-type-design.describe-if-circle-safely` |
| Wrapping e identidade de erro | `go-errors.wrap-sentinel-error-with-context`, `go-errors.extract-field-with-errors-as`, `go-errors.classify-wrapped-errors`, `go-errors.translate-error-at-boundary` |
| Testabilidade e controle do tempo/lifecycle | `go-testing.select-active-tokens-with-injected-clock`, `go-testing.measure-reservation-remaining`, `go-testing.count-deadline-clock-calls`, `go-testing.check-ticket-boundary` |

As dez competências essenciais com profundidade contextual mínima têm pelo
menos dois desafios distintos cada. O conjunto abaixo é a saída registrada por
`TestFoundationCatalogAudit`; variantes podem aparecer como um contexto
pedagógico, mas não entram no crédito de profundidade das bases.

| Competência essencial | Dois ou mais contextos verificados |
|---|---|
| `clone-struct-with-independent-slice` | `go-type-design.clone-inventory-independently`; `go-type-design.first-circle-clone` |
| `control-closure-capture` | `go-core.configure-server-with-functional-options`; `go-core.run-commands-with-cleanup-and-recovery` |
| `control-test-cleanup-lifecycle` | `go-testing.measure-reservation-remaining`; `go-testing.select-active-tokens-with-injected-clock` |
| `define-package-api-boundary` | `go-core.configure-server-with-functional-options`; `go-core.hide-counter-type-behind-interface` |
| `distinguish-zero-value-from-absence` | `go-data-text.build-scores-from-entries`; `go-data-text.lookup-map-value-with-comma-ok` |
| `inspect-go-toolchain-context` | `go-first-steps.declare-a-minimal-module`; `go-first-steps.format-rate-limit-message` |
| `return-true-nil-interface` | `go-type-design.first-circle-clone`; `go-type-design.select-validator-true-nil` |
| `testing-design-deterministic-time` | `go-testing.measure-reservation-remaining`; `go-testing.select-active-tokens-with-injected-clock` |
| `testing-detect-aliasing` | `go-testing.select-active-tokens-with-injected-clock`; `go-testing.snapshot-readings` |
| `wrap-sentinel-with-context` | `go-errors.load-and-validate-required-keys`; `go-errors.wrap-sentinel-error-with-context` |

### 9.5 Comandos, gates técnicos e interpretação

O help atual do CLI mostra `catalog validate` como subcomando. Portanto, a
forma válida de exercer o flag pedido é `catalog validate --v1-gate`; não há
um flag equivalente em `catalog` sem o subcomando `validate`.

| Comando | Resultado | Interpretação |
|---|---|---|
| `go run ./cmd/codinho help` | PASS | help raiz disponível |
| `go run ./cmd/codinho catalog` | Uso exibido, exit 2 | confirma que o próximo comando válido é `validate` |
| `go run ./cmd/codinho catalog validate --json` | PASS, exit 0; `diagnostics: null`, `editorial: null`, `StepNodes: 367` | catálogo estrutural e editorialmente carregável; tudo ainda draft |
| `go run ./cmd/codinho catalog validate --checks` | PASS, exit 0; 54 declarados/54 verificados | baselines produziram os resultados esperados e referências passaram |
| `go run ./cmd/codinho catalog validate --v1-gate` | FAIL, exit 1 | gate global de publicação, detalhado abaixo |
| `pose docs-check` | PASS; resumo `declared=16`, `undeclared=0`, `stale=0`, `errors=0`, `warnings=0`; imprime 2 linhas `REVIEW_PENDING` | documentação sem erro; os dois registros históricos de `docs/compatibility.md` (triggers `concept-content-authoring` e `learning-track-composition`) continuam pendentes |
| `pose check --strict` | PASS | estrutura POSE válida em modo estrito |

O `--v1-gate` usa a projeção elegível publicada, não o inventário draft. Como
há zero desafios publicados, ele reporta 5 `v1_gate_below_threshold` (0 contra
160/100/84/12/500), 36 `type_distribution_mismatch` (a política global inclui
packs futuros e seus tipos) e 1 `v1_check_proof_incomplete` (zero checks
publicados). Esses findings são gates de publicação/draft e decisão humana,
não defeitos de conteúdo nos sete packs; o gate editorial de 300/304 nodes
está satisfeito.

### 9.6 Fixtures, compatibilidade, duplicação e vazamento de solução

- Compatibilidade: o `go.mod` raiz e as fixtures positivas/de referência
  declaram Go 1.25.0. Há uma única fixture de baseline deliberadamente
  inválida em `go-first-steps.declare-a-minimal-module`, com `go 1.24.0`; a
  própria aceitação do exercício exige corrigir esse valor para 1.25.0, então
  ela é um caso negativo pedagógico e não uma dependência publicada. O runner
  observado foi `go1.26.5-X:nodwarf5 linux/amd64`.
- Rede, segredos e paths: a varredura das fixtures YAML em `packs/` não
  encontrou URL, `network: true`, chave privada, token ou segredo material,
  nem path absoluto, traversal ou path externo. O único match lexical de
  `secret` foi o campo sintético `secret-label` em
  `go-errors.yaml`, sem credencial. Os checks usam fixtures locais e biblioteca
  padrão.
- Duplicação: `LoadPacks`/`catalog validate` carregou IDs e relações sem
  diagnóstico; `TestFoundationVariantOrigins` confirma origem, classificação e
  competência das quatro variantes de cada pack autorizado. Repetições nas
  cinco trilhas são posições pedagógicas intencionais, não desafios duplicados.
  A revisão semântica humana de sinônimos continua sendo uma limitação do
  validador automático.
- Vazamento de solução: `editorial: null` e os checks automatizados não
  encontraram solução exposta em brief, critérios ou hints. As implementações
  de `reference_fixture` são privadas por contrato do runner; stubs com
  `TODO`/`panic("not implemented")` pertencem ao baseline dos exercícios e não
  foram tratados como divulgação de solução. Ainda falta a pré-revisão
  automatizada independente, a revisão por outra pessoa e o playtest humano.
- Publicação: nenhum `status: published`, `publication.reviewed_by` ou
  `playtested: true` foi preenchido. A quantidade publicada permanece zero por
  decisão de governança, não por falha de carregamento.

### 9.7 Checks desta execução e limitações remanescentes

Também foram executados os checks aplicáveis ao código e à documentação:

- `go test ./internal/curriculum -run
  '^TestFoundation(VariantOrigins|CatalogAudit)$' -count=1 -v`: PASS; ledger
  104/60, 45 bases, 8 variantes, 305 brutos, 304 creditáveis e dez
  competências essenciais com pelo menos dois contextos.
- `go test ./...`: FAIL ambiental nos testes de integração que constroem o CLI;
  o ambiente não conseguiu obter o status VCS e sugeriu `-buildvcs=false`.
- `go test -race ./...`: FAIL pelo mesmo erro ambiental de stamping VCS.
- `env GOFLAGS='-mod=readonly -buildvcs=false' go test ./...`: PASS.
- `env GOFLAGS='-mod=readonly -buildvcs=false' go test -race ./...`: PASS.
- `go vet ./...`: PASS.
- `bash scripts/ci/check-format.sh`: PASS.
- `pose docs-check`: PASS (`declared=16`, `undeclared=0`, `stale=0`,
  `errors=0`, `warnings=0`); os dois avisos históricos `REVIEW_PENDING` de
  `docs/compatibility.md` permanecem fora dos paths autorizados.
- `pose check --strict`: PASS (`valid POSE structure`).
- `git diff --check`: PASS.

As limitações que permanecem para a revisão primária são os gates globais de
publicação que corretamente veem zero conteúdo elegível, o `--v1-gate` abaixo
dos limiares por causa do estado draft, a ausência de `govulncheck` nesta
execução e a falha ambiental de stamping VCS nos comandos Go sem
`-buildvcs=false`. A execução com essa flag passou nos testes, vet/race e
checks editoriais; ela só desabilita metadado VCS do build e não altera o
conteúdo auditado. Revisão humana, playtest, classificação dos protótipos e
publicação continuam pendentes.

## 10. Roteiro de revisão humana e playtest (`prepare-human-review-playtest`)

Esta seção é um artefato de preparação para a revisão primária, não um relato
de revisão concluída. Ela identifica as 53 entradas dos sete packs
fundamentais (45 bases, incluindo o protótipo histórico, e oito variantes) e
mantém também o protótipo de `go-debugging` para que nenhum desafio do
manifesto fique sem destino explícito. Cada linha é uma ficha independente:
deve ser executada por uma pessoa que não seja o autor da linha e deve guardar
a evidência da sessão real.

O estado observado antes desta tarefa permanece intacto: os 54 desafios são
rascunhos na projeção do catálogo (`publication.status` ausente conta como
rascunho), nenhum tem `publication.reviewed_by`, nenhum tem
`playtested: true` e a contagem publicada é zero. Esta seção não preenche
esses campos e não recomenda publicação. `PENDENTE` abaixo significa que a
execução humana ainda não ocorreu; não é um resultado positivo automatizado.

### 10.1 Protocolo comum da ficha

Use os seguintes passos para cada ID, sem compartilhar o workspace entre
itens:

1. Confirme o ID, a versão, o título, o `author`, `variant_of` (se houver),
   as fixtures e os checks com `catalog_get` ou `catalog show`. O autor é o
   valor já existente no pack; quando o pack não declara autoria histórica,
   registre essa lacuna em vez de inventar um nome.
2. A partir da raiz do repositório, crie um destino temporário novo e rode
   `go run ./cmd/codinho workspace prepare <challenge-id> --dest <destino>`.
   Não use `--overwrite`. A exceção `WS-empty` abaixo é intencional para o
   único protótipo sem fixture distribuída.
3. Inicie o servidor com `go run ./cmd/codinho serve --authoring`, conecte um
   host MCP real e execute `catalog_get` seguido de
   `session_start({challenge_id: <ID>, mode: practice, depth: micro,
   help: progressive, evaluation: on_demand, request_id: <ID>})`. Para
   `go-debug.slice-off-by-one`, use `mode: debug` e preserve a ordem de
   reprodução, localização, hipótese, observação e menor correção. O modo
   `--authoring` é necessário porque o conteúdo ainda é draft; registre o
   aviso `content_provenance: draft` retornado pelo servidor.
4. Guarde `session_id` e a revisão inicial. Consulte `session_get` antes de
   inferir o passo ativo e `instruction_get` para cada instrução. Para cada
   passo, faça `workspace_observe` antes de qualquer avaliação, deixe o
   participante tentar, observe novamente, rode somente o check declarado
   com `check_run`, leia a evidência com `evidence_get`, e só então registre
   avaliação/reflexão e avance. Não use `override` para fabricar conclusão.
5. Encerre com `session_finish` por escolha explícita da pessoa. Na ficha,
   registre o `session_id`, os `evidence_id`, os checks executados, as pistas
   usadas, a resposta às reflexões, os achados e a decisão do revisor.
   Registre também `started_at_utc`, `finished_at_utc` e
   `observed_duration`; nunca copie `estimated_minutes` para esse último
   campo.

O campo `A/H` das tabelas é `autor no pack / revisor humano`. Enquanto a
execução não existir, o segundo valor fica como `PENDENTE (distinto de A)`.
O revisor primário deve substituir o marcador pelo handle real somente depois
de executar a sessão e julgar o conteúdo; esse registro editorial não deve ser
copiado para `publication.reviewed_by` antes da decisão de publicação. A
mesma pessoa pode revisar mais de uma linha, desde que continue distinta do
autor de cada linha.

Os marcadores de evidência automatizada usados por item são:

- `E1`: `go run ./cmd/codinho catalog validate --json`, executado nesta
  rodada com `diagnostics: null`, 54 desafios, 367 nodes, 54 drafts e zero
  publicados.
- `E2-F`: `go run ./cmd/codinho catalog validate --checks --json`, check
  `focused-tests`, baseline esperado `test_failure` e referência `pass`;
  ambos foram verificados.
- `E2-C`: o mesmo comando e check, com baseline esperado `compile_failure` e
  referência `pass`.
- `E2-M`: o mesmo comando, `module-contract`, baseline esperado
  `test_failure` e referência `pass`.
- `E2-MF`: o mesmo comando, `message-content` (`test_failure` → `pass`) e
  `format-safety` (`pass` → `pass`), todos verificados.
- `E4`: `go test ./internal/curriculum -run
  '^TestFoundation(VariantOrigins|CatalogAudit)$' -count=1 -v`, que
  reconcilia origem/classificação das variantes e o ledger fundamental.
- `E5`: `env GOFLAGS='-mod=readonly -buildvcs=false' go test
  ./cmd/codinho -run '^TestFoundationTracksOverRealStdio$' -count=1`,
  verificação de roteamento MCP dos cinco percursos. Como a própria
  implementação do teste registra, travessia com `override` é uma prova de
  alcance, não um playtest humano.

Essas evidências automatizadas só provam estrutura, expectativa de fixture,
checks ou roteamento. As perguntas nas linhas são deliberadamente deixadas
para o revisor humano, que deve registrar se o brief não entrega a solução,
se cada micro tem uma intenção, se os hints sobem em custo real, se os
critérios são observáveis e se a reflexão exige causalidade. Para variantes,
compare ainda o contexto com `variant_of` sem contar a variante como novo
desafio canônico.

### 10.2 Fichas por desafio e variante

`WS(ID)` significa o comando de preparação do passo 2 com o ID da própria
linha e um destino temporário exclusivo. `WS-empty(ID)` significa iniciar um
diretório temporário vazio, sem copiar `validation.reference_fixture`; o
protótipo histórico não declara `fixture`, e `workspace prepare` deve ser
registrado como não aplicável por `ErrEmptyFixture`. `S(ID)` significa a
sessão MCP do passo 3 para esse ID. Em `tempo`, `E` é a estimativa registrada
no catálogo e `O` é a duração observada exigida, ainda pendente.

#### `go-core`

| Conteúdo | Workspace | Sessão | Tempo | A/H | Evidência automatizada | Dúvida pedagógica para registrar |
|---|---|---|---|---|---|---|
| `go-core.close-resources-in-defer-order` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa distingue a ordem LIFO do momento em que o argumento de `defer` é capturado e explica por que a entrada vazia conserva o contrato? |
| `go-core.sum-variadic-numbers` (base) | `WS(ID)` | `S(ID)` | E 10 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-C` — `focused-tests` | A pessoa reconhece a diferença entre uma assinatura variádica e um slice obrigatório e explica o caso sem argumentos? |
| `go-core.stop-processing-commands-with-labeled-break` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa consegue explicar por que um `break` sem label não encerra o laço externo ao sair do `switch`? |
| `go-core.increment-counter-with-pointer-receiver` (base) | `WS(ID)` | `S(ID)` | E 10 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-C` — `focused-tests` | A pessoa relaciona o receiver escolhido à persistência do estado entre chamadas e identifica a cópia indevida? |
| `go-core.recover-from-panic-in-safe-call` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica em que fronteira a recuperação precisa ocorrer e como o caminho normal se diferencia do caminho de `panic`? |
| `go-core.guard-invariant-with-unexported-field` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-C` — `focused-tests` | A pessoa explica como a API impede estados inválidos sem depender de acesso externo ao campo e por que a validação ocorre no ponto observado? |
| `go-core.find-first-value-at-least` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa separa o valor zero do indicador de que um nó foi encontrado, inclusive para `min: 0` e lista nil? |
| `go-core.hide-counter-type-behind-interface` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa consegue justificar o menor contrato que o consumidor precisa e o que deixa de ser parte da API pública? |
| `go-core.run-commands-with-cleanup-and-recovery` (base) | `WS(ID)` | `S(ID)` | E 35 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa distingue limpeza garantida, recuperação local e parada do laço; a técnica pedida está compreensível sem o teste fingir provar a sintaxe? |
| `go-core.configure-server-with-functional-options` (base) | `WS(ID)` | `S(ID)` | E 30 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa percebe por que uma opção intermediariamente inválida pode ser corrigida por uma posterior, enquanto a configuração final inválida deve falhar? |

#### `go-data-text`

| Conteúdo | Workspace | Sessão | Tempo | A/H | Evidência automatizada | Dúvida pedagógica para registrar |
|---|---|---|---|---|---|---|
| `go-data-text.filter-without-mutating-input` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa testa independência da saída, preserva ordem e explica por que uma igualdade inicial não prova ausência de aliasing? |
| `go-data-text.truncate-bytes-at-rune-boundary` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica a diferença entre limite em bytes e limite em runes e o que torna uma fronteira UTF-8 válida? |
| `go-data-text.lookup-map-value-with-comma-ok` (base) | `WS(ID)` | `S(ID)` | E 10 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa distingue chave ausente de valor zero e prevê o comportamento de uma consulta em map nil sem escrever nele? |
| `go-data-text.quote-csv-field-when-needed` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa identifica quais caracteres mudam a representação CSV e explica por que a decisão não deve ser uma simples concatenação? |
| `go-data-text.dedupe-preserving-first-occurrence` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa consegue explicar como a primeira ocorrência vence sem perder a ordem de entrada nem confundir deduplicação com ordenação? |
| `go-data-text.sum-valid-integers-safely` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa distingue zero válido de erro de parsing e explica como o contrato trata entradas válidas e inválidas na mesma sequência? |
| `go-data-text.preallocate-slice-with-zero-length` (base) | `WS(ID)` | `S(ID)` | E 10 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica por que capacidade e comprimento têm papéis diferentes e por que uma slice de comprimento zero não pode ser indexada? |
| `go-data-text.title-first-letters-preserving-spacing` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa preserva espaços e fronteiras Unicode ao capitalizar e consegue explicar por que normalizar whitespace seria outro comportamento? |
| `go-data-text.build-scores-from-entries` (base) | `WS(ID)` | `S(ID)` | E 30 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica como o placar escolhe o primeiro valor válido e separa ausência, valor zero e erro de uma entrada? |
| `go-data-text.format-names-as-csv-fields` (base) | `WS(ID)` | `S(ID)` | E 30 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa consegue explicar a interação entre capitalização, preservação de campos e escaping CSV sem misturar responsabilidades? |

#### `go-first-steps`

| Conteúdo | Workspace | Sessão | Tempo | A/H | Evidência automatizada | Dúvida pedagógica para registrar |
|---|---|---|---|---|---|---|
| `go-data.slice-filter-preserve-input` (protótipo histórico) | `WS-empty(ID)` | `S(ID)` | E 15 min; O PENDENTE | `autoria histórica não declarada / PENDENTE (atribuir autor antes)` | `E1`, `E2-F` — `focused-tests` com baseline privado | A pessoa entende que preservar a entrada também exige discutir o armazenamento compartilhado da saída; a ausência de starter torna o briefing executável sem revelar referência? |
| `go-first-steps.enumerate-weekdays-with-iota` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica a sequência produzida por `iota` e por que a regra de dia útil não precisa comparar strings? |
| `go-first-steps.convert-celsius-to-fahrenheit` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa identifica onde a conversão inteira perderia precisão e justifica o tipo usado no cálculo? |
| `go-first-steps.declare-a-minimal-module` (base) | `WS(ID)` | `S(ID)` | E 10 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-M` — `module-contract` | A pessoa diferencia o módulo declarado, a versão da linguagem e o diretório efetivo, sem tratar a criação do arquivo como compreensão do contrato? |
| `go-first-steps.clamp-int-to-byte` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica a diferença entre limitar o domínio antes da conversão e aceitar wraparound silencioso? |
| `go-first-steps.format-rate-limit-message` (base) | `WS(ID)` | `S(ID)` | E 10 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-MF` — `message-content`, `format-safety` | A pessoa distingue conteúdo correto da mensagem de segurança de tipos dos verbos de formatação e explica o que cada check prova? |
| `go-first-steps.avoid-shadowing-named-returns` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa localiza a variável que esconde o retorno nomeado e explica por que o fluxo de erro pode devolver um valor diferente do esperado? |

#### `go-errors`

| Conteúdo | Workspace | Sessão | Tempo | A/H | Evidência automatizada | Dúvida pedagógica para registrar |
|---|---|---|---|---|---|---|
| `go-errors.wrap-sentinel-error-with-context` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica por que `errors.Is` preserva a categoria mesmo quando a mensagem ganha contexto e por que comparar texto seria frágil? |
| `go-errors.extract-field-with-errors-as` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa consegue explicar como dados estruturados atravessam o encadeamento e por que a extração não depende da mensagem? |
| `go-errors.validate-user-joining-all-errors` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica o benefício de retornar todas as violações e como a árvore de erros afeta a classificação? |
| `go-errors.translate-error-at-boundary` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa identifica quais detalhes devem parar na fronteira e quais propriedades semânticas o consumidor ainda precisa observar? |
| `go-errors.load-and-validate-required-keys` (base) | `WS(ID)` | `S(ID)` | E 25 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa diferencia chave ausente de chave presente com valor vazio e mantém a identidade do erro de entrada ao agregar faltas? |
| `go-errors.classify-wrapped-errors` (base) | `WS(ID)` | `S(ID)` | E 35 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests` | A pessoa consegue explicar a precedência entre identidade, classificação e tipo ao atravessar wrappers e joins, sem usar a mensagem como protocolo? |
| `go-errors.maintain-memory-catalog` (base) | `WS(ID)` | `S(ID)` | E 60 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica por que validar antes de mutar, retornar snapshots isolados e distinguir zero de ausência são contratos diferentes do catálogo? |

#### `go-io`

| Conteúdo | Workspace | Sessão | Tempo | A/H | Evidência automatizada | Dúvida pedagógica para registrar |
|---|---|---|---|---|---|---|
| `go-io.decode-strict-config` (base) | `WS(ID)` | `S(ID)` | E 30 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests` | A pessoa separa erro do Reader, documento adicional, campo desconhecido e regra de domínio sem transformar um `Decode` bem-sucedido em documento completo? |
| `go-io.import-csv-inventory` (base) | `WS(ID)` | `S(ID)` | E 50 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica como validar largura/valores antes de escrever e o que o contrato garante quando o writer faz uma escrita curta? |
| `go-io.read-fragmented-note` (variante de `go-io.decode-strict-config`) | `WS(ID)` | `S(ID)` | E 20 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests`; `E4` origem | A variante torna visível a diferença entre leitura parcial e EOF sem repetir integralmente o desafio de configuração? |
| `go-io.decode-weather-fields` (variante de `go-io.decode-strict-config`) | `WS(ID)` | `S(ID)` | E 20 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests`; `E4` origem | O contexto de previsão deixa claro o contrato de campo desconhecido e acrescenta decisão real, sem apenas renomear o exemplo da origem? |
| `go-io.read-quoted-attendees` (variante de `go-io.import-csv-inventory`) | `WS(ID)` | `S(ID)` | E 20 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests`; `E4` origem | A pessoa percebe por que nomes citados exigem parser CSV e por que separar por vírgula não preserva campos? |
| `go-io.write-result-report` (variante de `go-io.import-csv-inventory`) | `WS(ID)` | `S(ID)` | E 20 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests`; `E4` origem | A variante evidencia a propagação de falha do writer e o estado parcial sem duplicar a validação de entrada da origem? |

#### `go-testing`

| Conteúdo | Workspace | Sessão | Tempo | A/H | Evidência automatizada | Dúvida pedagógica para registrar |
|---|---|---|---|---|---|---|
| `go-testing.select-active-tokens-with-injected-clock` (base) | `WS(ID)` | `S(ID)` | E 35 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica a fronteira estrita antes/igual/depois, uma única observação do relógio e o isolamento da saída? |
| `go-testing.measure-reservation-remaining` (variante de `go-testing.select-active-tokens-with-injected-clock`) | `WS(ID)` | `S(ID)` | E 20 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests`; `E4` origem | O cálculo do tempo restante exige distinguir duração negativa, igualdade e instante injetado, em vez de só trocar o domínio do exemplo? |
| `go-testing.check-ticket-boundary` (variante de `go-testing.select-active-tokens-with-injected-clock`) | `WS(ID)` | `S(ID)` | E 20 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests`; `E4` origem | A pessoa explica por que validade na igualdade é uma decisão de contrato e como fusos diferentes podem representar o mesmo instante? |
| `go-testing.snapshot-readings` (variante de `go-testing.select-active-tokens-with-injected-clock`) | `WS(ID)` | `S(ID)` | E 20 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests`; `E4` origem | A pessoa prova isolamento mutando entrada e saída em momentos distintos, sem confundir igualdade de conteúdo com independência? |
| `go-testing.count-deadline-clock-calls` (variante de `go-testing.select-active-tokens-with-injected-clock`) | `WS(ID)` | `S(ID)` | E 20 min; O PENDENTE | `codex / PENDENTE (≠ codex)` | `E1`, `E2-F` — `focused-tests`; `E4` origem | O dublê que muda a cada chamada torna claro por que o lote deve compartilhar um instante e ainda consultar o relógio no lote vazio? |

#### `go-type-design`

| Conteúdo | Workspace | Sessão | Tempo | A/H | Evidência automatizada | Dúvida pedagógica para registrar |
|---|---|---|---|---|---|---|
| `go-type-design.clone-inventory-independently` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa distingue cópia rasa da independência da slice interna e testa mutações nos dois sentidos? |
| `go-type-design.select-validator-true-nil` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica como um valor tipado nil pode tornar uma interface não nil e por que o contrato exige nil verdadeiro? |
| `go-type-design.sum-list-with-nil-safe-receiver` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa justifica o comportamento de um receiver nil e diferencia segurança deliberada de um panic oculto? |
| `go-type-design.contains-generic-comparable` (base) | `WS(ID)` | `S(ID)` | E 10 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica por que a constraint precisa permitir comparação e quais tipos não podem ser aceitos pelo contrato? |
| `go-type-design.describe-if-circle-safely` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa usa o resultado da assertion para separar tipo esperado, tipo diferente e ausência sem provocar panic? |
| `go-type-design.map-generic-transform` (base) | `WS(ID)` | `S(ID)` | E 15 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica o papel dos dois tipos genéricos e como a transformação preserva ordem e contrato de entrada? |
| `go-type-design.first-circle-clone` (base) | `WS(ID)` | `S(ID)` | E 30 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa consegue separar clonagem independente, assertion segura e retorno nil sem tratar um caso combinado como três soluções desconexas? |
| `go-type-design.generic-node-values-nil-safe` (base) | `WS(ID)` | `S(ID)` | E 25 min; O PENDENTE | `claude / PENDENTE (≠ claude)` | `E1`, `E2-F` — `focused-tests` | A pessoa explica como generics e receiver nil-safe interagem e qual parte do contrato pertence à lista, não à constraint? |

#### `go-debugging` fora do escopo fundamental

| Conteúdo | Workspace | Sessão | Tempo | A/H | Evidência automatizada | Dúvida pedagógica para registrar |
|---|---|---|---|---|---|---|
| `go-debug.slice-off-by-one` (protótipo debug) | `WS(ID)` — fixture `main.go` | `S(ID; mode=debug)` | E 20 min; O PENDENTE | `autoria histórica não declarada / PENDENTE (atribuir autor antes)` | `E1` — sem check declarado; a execução do programa e a menor correção são evidência humana exigida | A pessoa reproduz antes de editar, localiza a primeira divergência, formula a hipótese e considera honestamente a fronteira vazia sem ampliar a correção? |

### 10.3 Modelo de resultado humano e gate

Para cada linha, o revisor humano deve anexar ao relatório primário:

```text
ID/versão:
autor no pack:
revisor humano:
session_id:
content_provenance:
started_at_utc:
finished_at_utc:
observed_duration:
workspace_manifest/evidence_ids:
checks e resultados:
pistas usadas e nível máximo:
reflexões respondidas:
achados pedagógicos:
decisão: pendente | pedir correção | pronto para revisão de publicação
```

O gate desta tarefa permanece `HUMAN_REVIEW_PLAYTEST_PENDING` para todas as
54 linhas. A validação automatizada `E1`/`E2`/`E4`/`E5` não altera esse estado,
não autentica uma pessoa, não mede duração observada e não autoriza
publicação. A pré-revisão automatizada independente também continua
pendente nesta execução, pois este autor não a substitui nem inicia outro
agente. Até que cada ficha tenha uma sessão MCP real, um autor e um revisor
humano distintos, início/fim observados e julgamento pedagógico registrado,
`reviewed_by`, `playtested` e `publication.status: published` devem continuar
inalterados.

## 11. Preparação de reconciliação e fechamento (`prepare-publication-closeout`)

Esta seção é a entrega documental da tarefa `prepare-publication-closeout`,
dependente de `prepare-human-review-playtest`. Ela prepara o material para a
revisão primária e o coordenador, mas não fecha a spec, não altera seu estado e
não recomenda publicação. A ausência de uma decisão humana ou de evidência real
mantém `go-foundations-packs` em `in-progress`.

### 11.1 Estado reconciliado em 2026-09-15 UTC

A medição desta execução separa inventário, rascunhos, publicação e conteúdo
elegível. O comando global encontrou `105` conceitos, `64` competências, `5`
trilhas, `54` desafios e `366` nodes; todos os `54` desafios continuam draft,
com `0` publicados e `0` elegíveis. No recorte dos sete packs, o teste de
auditoria preserva `104` conceitos, `60` competências, `45` bases, `8`
variantes, `304` nodes brutos e `303` nodes creditáveis. A base ainda contém o
protótipo histórico de filtro; sua exclusão conceitual produz a distribuição
pretendida de `32` atômicos, `10` combinados e `2` fatias funcionais.

| Projeção | Resultado atual | Interpretação de fechamento |
|---|---:|---|
| Manifest completo | `105` conceitos / `64` competências / `54` desafios / `366` nodes | Inventário carregável, ainda sem conteúdo elegível. |
| Sete packs fundamentais | `104` conceitos / `60` competências / `45` bases + `8` variantes | Metas quantitativas de conteúdo atendidas no recorte documentado; variantes não pagam profundidade. |
| Profundidade fundamental | `304` brutos / `303` creditáveis | O limiar de `300` da spec continua atendido; o alvo operacional anterior foi reduzido ao remover um micropasso redundante. |
| Publicação | `0` publicados / `0` elegíveis / `54` drafts | Gate de publicação não atendido. |
| Revisão/playtest | `0` fichas humanas concluídas nesta execução | Gate `HUMAN_REVIEW_PLAYTEST_PENDING`. |

Evidência executável desta reconciliação:

- `go run ./cmd/codinho catalog validate --json`: PASS, `diagnostics: null`,
  `editorial: null`, cobertura global `105/64/5/54/366`.
- `go run ./cmd/codinho catalog validate --checks --json`: PASS, `54`
  checks declarados e `54` verificados; cada resultado mantém cenário,
  expectativa, resultado e digest.
- `go run ./cmd/codinho catalog validate --v1-gate`: FAIL, exit `1`, como
  esperado para `0` conteúdo publicado/elegível; o relatório mostra os
  limiares globais, a distribuição publicada e a prova de checks publicados
  ainda não satisfeitos.
- `go test ./internal/curriculum -run
  '^TestFoundation(VariantOrigins|CatalogAudit)$' -count=1 -v`: PASS; o
  `TestFoundationCatalogAudit` registra `104/60`, `45` bases, `8` variantes,
  `304` nodes e `303` creditáveis.
- `pose docs-check`: PASS com resumo `declared=16`, `undeclared=0`, `stale=0`,
  `errors=0` e `warnings=0`; a saída imprime dois registros históricos
  `REVIEW_PENDING`, que não são convertidos em aprovação.
- `pose check --strict` e `git diff --check`: PASS.

As duas linhas `REVIEW_PENDING` impressas pelo `docs-check` pertencem aos
triggers históricos de `docs/compatibility.md`; elas não são erros do documento
nem evidência de revisão concluída. A execução de 2026-09-15 acima é a medição
corrente e deixa essa pendência visível.

### 11.2 Trace R1–R10 da spec `go-foundations-packs`

O trace abaixo usa os R-IDs da seção Requirements da [spec de fundamentos](../../.pose/specs/2026-08-22-go-foundations-packs.md#2-requirements).
“Atendido no draft” descreve conteúdo ou prova técnica disponível; não equivale
a publicação nem a aceite pedagógico.

| R-ID | Situação atual | Evidência dos lotes | Gate ou decisão que falta |
|---|---|---|---|
| R1 | **Pendente para publicação.** Os sete packs e as cinco trilhas existem como draft; `published=0` e `eligible=0`. | [go-first-steps: estrutura](../../.pose/reports/2026-09-10-first-steps-structure-content.md), [go-core](../../.pose/reports/2026-09-10-core-concept-content.md), [go-data-text](../../.pose/reports/2026-09-10-data-text-concept-content.md), [go-type-design](../../.pose/reports/2026-09-10-type-design-concept-content.md), [go-errors](../../.pose/reports/2026-09-10-go-errors-authoring-batch.md), [go-io](../../.pose/reports/2026-09-09-go-io-authoring-batch.md), [go-testing](../../.pose/reports/2026-09-10-go-testing-authoring-batch.md), [trilhas](../../.pose/reports/2026-09-10-foundation-track-routing.md). | Completar pré-revisão independente, revisão/playtest humanos e os metadados verdadeiros antes de recomendar `published`. |
| R2 | **Atendido no draft.** A seção 2 reconcilia os temas A–M e os conceitos nos sete packs. | [go-first-steps](../../.pose/reports/2026-09-10-first-steps-concept-content.md), [runtime](../../.pose/reports/2026-09-10-first-steps-runtime-content.md), [estrutura](../../.pose/reports/2026-09-10-first-steps-structure-content.md), [go-core](../../.pose/reports/2026-09-10-core-concept-content.md), [go-data-text](../../.pose/reports/2026-09-10-data-text-concept-content.md), [go-type-design](../../.pose/reports/2026-09-10-type-design-concept-content.md), [go-errors](../../.pose/reports/2026-09-10-error-concept-content.md), [go-io](../../.pose/reports/2026-09-10-io-concept-content.md), [go-testing](../../.pose/reports/2026-09-10-testing-concept-content.md). | Revisar semanticamente se cada exemplo conceitual permanece distinto do desafio e não revela a solução. |
| R3 | **Parcial.** O recorte conceitual tem `32` atômicos depois da exclusão do protótipo histórico; o armazenamento ainda mostra `33` bases atômicas até a decisão de disposição. | [inventário editorial histórico](../editorial/go-foundations-inventory.md), [ledger de protótipos e distribuição](#1-ledger-reconciliado-após-a-remediação), [fila de disposição](../editorial/go-foundations-queue.json). | Registrar a decisão humana sobre `go-data.slice-filter-preserve-input` e manter seus IDs/sessões sem reclassificação silenciosa. |
| R4 | **Atendido no draft.** Há `10` desafios combinados no recorte fundamental, com critérios e árvores verificadas. | [complemento go-errors](../../.pose/reports/2026-09-10-go-errors-authoring-batch.md), [lote go-io](../../.pose/reports/2026-09-09-go-io-authoring-batch.md), [lote go-testing](../../.pose/reports/2026-09-10-go-testing-authoring-batch.md), [auditoria global](#9-auditoria-da-remediação-de-profundidade-audit-depth-remediation). | Fazer revisão humana das intenções, hints, critérios e reflexões; a contagem não autoriza publicação. |
| R5 | **Atendido no draft.** As duas fatias funcionais são o importador validado de `go-io` e o catálogo em memória de `go-errors`. | [importador CSV/JSON](../../.pose/reports/2026-09-09-go-io-authoring-batch.md), [catálogo em memória](../../.pose/reports/2026-09-10-go-errors-authoring-batch.md), [nodes do catálogo](../../.pose/reports/2026-09-10-errors-catalog-nodes.md). | Confirmar em playtest que o escopo funcional é compreensível, observável e não exige arquitetura não declarada. |
| R6 | **Atendido quantitativamente no draft.** A remediação chegou a `45/45/214` nodes brutos e `303` creditáveis após a remoção do micropasso redundante; a suficiência pedagógica continua humana. | [auditoria de profundidade](#9-auditoria-da-remediação-de-profundidade-audit-depth-remediation), [nodes de go-core](../../.pose/reports/2026-09-10-core-defer-empty-node.md), [nodes de go-testing](../../.pose/reports/2026-09-10-testing-clock-call-nodes.md), [regra de autoria](../content-authoring.md#sequência-e-alternativas-da-árvore). | Revisar se cada micro tem uma intenção, evidência e critério próprios e se os hints sobem em custo real. |
| R7 | **Atendido mecanicamente no draft.** `catalog validate --checks` verificou `54/54`; as oito variantes têm `variant_of`, `canonical: false` e fixtures próprias. | [revisão de hints](../../.pose/reports/2026-09-10-catalog-hint-review.md), [revisão de regressões](../../.pose/reports/2026-09-10-catalog-luna-review.md), [lote de variantes](../../.pose/reports/2026-09-10-foundation-variants-authoring-batch.md). | Completar julgamento humano sobre brief, leakage, micropassos, reflexões e equivalência das variantes. |
| R8 | **Atendido tecnicamente no draft.** As cinco trilhas têm `47` posições e o roteamento MCP real foi atravessado em modo de autoria. | [roteamento das cinco trilhas](../../.pose/reports/2026-09-10-foundation-track-routing.md), [auditoria das trilhas](#6-auditoria-das-cinco-trilhas-audit-tracks). | Fazer playtest humano da progressão; a travessia com `override` prova alcance, não aprendizagem nem publicação. |
| R9 | **Atendido tecnicamente no draft.** Nil, zero values, aliasing, UTF-8, method sets, interface nil, wrapping e tempo determinístico têm contextos registrados; as dez competências essenciais têm pelo menos dois contextos. | [cobertura técnica](#94-auditoria-de-cobertura-técnica-e-contextos-essenciais), [conteúdo conceitual revisado](../../.pose/reports/2026-09-09-review-concept-content-authoring.md), [relações de aplicação](../../.pose/reports/2026-09-10-catalog-relations-review.md). | Confirmar por revisão humana que os contextos ensinam decisões causais sem duplicação semântica. |
| R10 | **Pendente.** Não há sessões MCP humanas, durações observadas, `reviewed_by` distinto ou `playtested: true` registrados para as `54` linhas. | [roteiro por desafio](#10-roteiro-de-revisão-humana-e-playtest-prepare-human-review-playtest), [checklist editorial](../content-review-checklist.md), [revisão técnica de publication-integrity](../../.pose/reports/2026-09-08-review-catalog-publication-integrity.md). | Exigir um revisor humano diferente de cada autor, playtest ponta a ponta e evidência real antes de publicar. |

### 11.3 Gates pendentes de `concept-content-authoring` e `catalog-publication-integrity`

As duas specs de infraestrutura editorial estão `done` quanto à implementação
e à revisão técnica registrada, mas isso não fecha a curadoria dos sete packs.
Os gates abaixo pertencem ao uso dessas capacidades no catálogo fundamental.

| Spec / gate | Evidência disponível | Estado e condição de saída |
|---|---|---|
| `concept-content-authoring`: conteúdo canônico versionado, schema/loader, projeção MCP e não divulgação | [spec](../../.pose/specs/2026-09-07-concept-content-authoring.md), [revisão independente](../../.pose/reports/2026-09-09-review-concept-content-authoring.md) | Implementação **done**. Não provar, por contagem ou teste estrutural, a qualidade semântica de todos os exemplos; manter revisão humana dos conceitos como gate do pack. |
| `concept-content-authoring`: distinção dos exemplos contra briefs, fixtures e soluções | A própria revisão registra que a amostra não substitui a curadoria integral; [guia de autoria](../content-authoring.md#conteúdo-de-conceitos-content) delimita o contrato. | **Pendente.** Revisar cada exemplo autorado em contexto diferente e registrar leaks, sobreposição e decisão do revisor. |
| `catalog-publication-integrity`: projeções, distribuição, provas privadas e MCP publicado | [spec](../../.pose/specs/2026-09-07-catalog-publication-integrity.md), [revisão técnica aprovada](../../.pose/reports/2026-09-08-review-catalog-publication-integrity.md) | Implementação **done** e tecnicamente aprovada; o catálogo real ainda tem zero elegíveis. Não usar a aprovação da capability como aprovação dos packs. |
| `catalog-publication-integrity`: pré-revisão automatizada independente dos `54` itens | [checklist](../content-review-checklist.md), [roteiro da dependência](#10-roteiro-de-revisão-humana-e-playtest-prepare-human-review-playtest) | **Pendente nesta execução.** O autor não substitui nem inicia outro agente para preencher esse gate. |
| `catalog-publication-integrity`: revisão e playtest humanos distintos | [fichas de resultado](#103-modelo-de-resultado-humano-e-gate) e [checklist](../content-review-checklist.md#playtest) | **Pendente.** Completar as `54` fichas com sessão MCP real, workspace isolado, início/fim, duração observada, checks, pistas, reflexões e achados. |
| `catalog-publication-integrity`: checks publicados e gate V1 | Execução corrente de `catalog validate --v1-gate`: `0` publicados, `0` elegíveis e `0` checks publicados; [guia de publicação](../content-authoring.md#publicação-publication) | **Bloqueado por estado draft.** Só reexecutar como gate de publicação depois de existir conteúdo humano efetivamente revisado e publicado. |
| Fechamento da spec fundamental e reconciliação de follow-ups | [follow-up da spec de conceitos](../../.pose/specs/2026-09-07-concept-content-authoring.md#7-final-report), [follow-up da spec de publicação](../../.pose/specs/2026-09-07-catalog-publication-integrity.md#7-final-report), [follow-up da spec fundamental](../../.pose/specs/2026-08-22-go-foundations-packs.md#7-final-report) | **Pendente.** O coordenador deve atualizar a spec fundamental, preservar os follow-ups aplicáveis e só mudar para `done` quando nenhum gate acima faltar. |

### 11.4 Decisão requerida antes de qualquer recomendação de publicação

Não recomendar publicação neste estado. Antes de emitir uma recomendação, o
coordenador deve obter e registrar uma decisão humana, na spec ou em decisão
referenciada por ela, com pelo menos:

- responsável, data UTC, opções consideradas e racional;
- aceite ou alteração explícita da distribuição `32` atômicos, `10`
  combinados e `2` fatias funcionais;
- disposição de `go-data.slice-filter-preserve-input` e
  `go-debug.slice-off-by-one`, incluindo escopo, autor, versão e impacto nos
  IDs/sessões;
- tratamento das oito variantes como derivados `canonical: false`, sem
  contá-las como novos desafios;
- efeito da decisão nas contagens, na spec e nos follow-ups.

Também deve exigir revisão/playtest humanos reais para as `54` linhas do
manifesto: `53` entradas dos sete packs fundamentais (45 bases e 8 variantes)
mais o protótipo de debugging. O revisor deve ser diferente do autor, preparar
um workspace isolado, executar a sessão MCP real sem `override` para fabricar
conclusão, observar a duração e registrar a ficha da seção 10.3. Só depois de
todas as fichas aprovadas será legítimo preencher `reviewed_by`,
`playtested: true` e `publication.status: published`, seguido da validação dos
checks publicados. Uma pré-revisão automatizada, `catalog validate` ou uma
travessia MCP automatizada não substitui nenhuma dessas decisões.

### 11.5 Reconciliação documental e follow-ups preservados

Os documentos antigos continuam úteis como histórico imutável do estágio em que
foram escritos. A tabela abaixo evita que seus números intermediários sejam
apresentados como a medição corrente.

| Fonte histórica | Registro antigo | Reconciliação atual |
|---|---|---|
| [Inventário editorial de 2026-09-11](../editorial/go-foundations-inventory.md) no commit `9caef5c1a585813076d4a4bc3957d71c0e1a2fe2` | `70` conceitos, `56` competências, `53` entradas nos sete packs, `167` nodes de bases e `53` de variantes; `220` nodes nos sete packs. | Depois dos lotes conceituais, variantes e remediações: `104/60`, `45` bases, `8` variantes, `304` brutos e `303` creditáveis no escopo fundamental. Preservar a linha histórica e usar este documento para a leitura atual. |
| [Handoff `go-foundations-io-batch`](../../.pose/knowledge/2026-09-10-handoff-go-foundations-io-batch.md) | Estágio intermediário com `45` bases, `70` conceitos, `51` competências e `118` nodes de bases, além de `32` nodes das variantes. | Os lotes posteriores de conceitos e profundidade estão consumidos nas seções 3–5 e 9; o handoff não é evidência de fechamento. |
| [Spec `go-foundations-packs` §6/§7](../../.pose/specs/2026-08-22-go-foundations-packs.md#6-validation) | O `Results summary`, o trace e o `Final Report` ainda dizem autoria parcial, `70/51/118` e requisito pendente. A lista de tarefas também conserva checkboxes históricos não atualizados. | O estado corrente foi medido em 2026-09-15 e está neste trace; o coordenador deve reconciliar `Execution log`, `Results summary`, `Requirement trace`, `Known gaps` e `Follow-ups` da spec sem reescrever o histórico. |
| [Spec `concept-content-authoring`](../../.pose/specs/2026-09-07-concept-content-authoring.md) e [spec `catalog-publication-integrity`](../../.pose/specs/2026-09-07-catalog-publication-integrity.md) | Implementações `done` com follow-up `v1-integrated-acceptance`; as revisões técnicas explicitam que amostra, inventário e provas automatizadas não são curadoria humana do catálogo. | Tratar esses follow-ups como dependências de integração da V1, não como autorização para publicar os drafts fundamentais. Manter os gates da seção 11.3 abertos. |

Follow-ups que permanecem acionáveis:

- Reexecutar `v1-integrated-acceptance` no candidato composto, como já
  referenciado pelas specs de conteúdo conceitual e integridade de publicação.
- Registrar as sessões reais e a amostra humana dos sete packs, conforme o
  follow-up de `go-foundations-packs` e o roteiro da seção 10.
- Decidir a distribuição/protótipos e atualizar a spec fundamental somente
  depois de haver evidência para essa decisão.

### 11.6 Encaminhamento ao coordenador

Não executar ações de coordenador nesta tarefa. Após a entrega real da revisão,
do playtest e da decisão de distribuição, o coordenador deve:

1. atualizar a [spec fundamental](../../.pose/specs/2026-08-22-go-foundations-packs.md)
   com o trace final, o `completed_at` somente se aplicável e a disposição de
   cada follow-up;
2. executar `pose assess discover --update-state` para recalcular métricas e
   atualizar os artefatos de assessment/state;
3. executar os checks de fechamento aplicáveis, incluindo
   `pose artifact-check --spec go-foundations-packs --strict`,
   `pose lint-spec go-foundations-packs --ready-check`,
   `pose validate --strict --json .pose/results/delivery-validation.json`,
   `pose docs-check`, `pose check --strict`, `catalog validate --checks` e o
   `catalog validate --v1-gate` final;
4. comparar claims e observed, versões, commits e evidência humana, sem contar
   inventário draft como publicação;
5. manter `status: in-progress` e devolver pendências se qualquer gate faltar;
   só então considerar closeout e recomendação de publicação.

Este lote não executou `assess discover --update-state`, não alterou a spec nem
o estado do coordenador e não preencheu qualquer metadado humano. A medição e os
checks desta seção são evidência de reconciliação documental; não são um
closeout da spec.
