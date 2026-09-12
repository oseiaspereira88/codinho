---
title: Catálogo fundamental de Go
doc_type: reference
---

# Catálogo fundamental de Go — plano e auditoria de trilhas

Tarefas editoriais: `plan-concept-depth` (baseline de expansão) e
`audit-tracks` (reconciliação dos percursos). Inventário auditado em
2026-09-12 UTC contra o manifest atual, depois dos lotes conceituais listados
como dependências. Este documento é um plano de autoria e auditoria, não uma
declaração de publicação: a spec `go-foundations-packs` continua `in-progress`,
os sete packs continuam em rascunho e revisão humana/playtest continuam
pendentes.

## 1. Baseline reconciliado

Os números abaixo vêm de `go run ./cmd/codinho catalog validate --json` e de
uma contagem independente das árvores `layers[*].macro_steps`, recursiva. O nó
de `layer` não entra na contagem. `Base` inclui o conteúdo histórico; entre
parênteses, quando necessário, aparece o crédito conservador usado para a
meta: o protótipo histórico `go-data.slice-filter-preserve-input` tem um
micro-node, mas esse node não é creditado antes de sua disposição explícita.

| Pack | Temas A–M | Entradas de base/variantes | Conceitos | Competências | Nodes de base M/Meso/Micro | Nodes de variantes M/Meso/Micro |
|---|---|---:|---:|---:|---:|---:|
| `go-first-steps` | A, B, com E/C introdutórios | 7 bases | 22 | 9 | 6/6/31 = 43 (crédito 6/6/30) | 0/0/0 |
| `go-core` | C, D, K | 10 bases | 15 | 10 | 0/0/33 = 33 | 0/0/0 |
| `go-data-text` | E, F | 10 bases | 13 | 8 | 12/12/40 = 64 | 0/0/0 |
| `go-type-design` | G, H, I | 8 bases | 12 | 6 | 7/7/34 = 48 | 0/0/0 |
| `go-errors` | J | 7 bases | 17 | 10 | 5/5/35 = 45 | 0/0/0 |
| `go-io` | L | 2 bases + 4 variantes | 14 | 11 | 5/5/19 = 29 | 4/4/16 = 24 |
| `go-testing` | M | 1 base + 4 variantes | 11 | 6 | 4/4/11 = 19 | 5/5/22 = 32 |
| **Sete packs fundamentais** | **A–M** | **45 bases + 8 variantes** | **104** | **60** | **39/39/203 = 281** (**crédito 280**) | **9/9/38 = 56** |
| `go-debugging` (fora do escopo) | protótipo | 1 protótipo | 1 | 4 | 1/1/4 = 6 | 0/0/0 |
| **Manifest completo** | — | **46 bases/protótipo + 8 variantes** | **105** | **64** | **40/40/207 = 287** | **9/9/38 = 56** |

As 44 bases canônicas preservam a distribuição histórica de 32 atômicos, 10
combinados e 2 fatias funcionais. O inventário carregado tem uma 45ª base
atômica histórica, `go-data.slice-filter-preserve-input`; ela
continua no pack para preservar IDs e histórico, mas não é mais membro da
trilha `go-from-zero`, que usa o desafio canônico
`go-data-text.filter-without-mutating-input`. As quatro variantes de I/O e as
quatro de testes mantêm `variant_of` e `canonical: false` e não financiam a
meta de profundidade. O protótipo de depuração e suas quatro competências
também não financiam a meta fundamental.

### Regra de crédito

Há três projeções que devem permanecer visíveis em toda auditoria:

1. **Escopo fundamental:** sete packs, 104 conceitos e 60 competências hoje;
   45 bases e 8 variantes permanecem inventariadas separadamente.
2. **Crédito de profundidade:** nodes de bases, sem as oito variantes e sem o
   micro-node do protótipo histórico até a decisão de disposição.
3. **Manifesto:** inclui `go-debugging` e as variantes para fins de inventário,
   mas não transforma esses itens em conteúdo fundamental novo.

Os lotes conceituais já alcançaram **104 conceitos e 60 competências nos sete
packs**. A meta de profundidade continua sendo **304 nodes de bases
creditáveis**; o baseline atual credita 280, portanto ainda há um déficit
editorial real de 24 nodes. Com `go-debugging` preservado, o manifesto tem 105
conceitos, 64 competências e 287 nodes de base; os 56 nodes de variantes
continuam reportados separadamente.

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

## 5. Nodes e decisões pedagógicas que faltam

O alvo abaixo mede somente árvores de bases, excluindo variantes e não
creditando o micro-node do protótipo histórico. O baseline já incorpora os
lotes conceituais; macro e meso ainda podem organizar micros existentes quando
a decisão já está presente, mas nenhum micro novo pode repetir uma intenção
existente. As adições são números de nodes `macro/meso/micro`; cada micro
restante terá um verbo, um alvo, uma evidência e um critério observável.

| Pack | Crédito atual M/Meso/Micro | Alvo M/Meso/Micro | Adição | Decisões ausentes a autorar |
|---|---:|---:|---:|---|
| `go-first-steps` | 6/6/30 | 6/6/30 | 0/0/0 | Contexto da toolchain, formatação/build, default de constantes e representabilidade já entregues; o micro histórico de filtro não é creditado. |
| `go-core` | 0/0/33 | 5/5/38 | +5/+5/+5 | Escopo de inicializador, ciclo de `for`, avaliação de argumentos de `defer`, captura de closure, ligação de method value e package consumidor. |
| `go-data-text` | 12/12/40 | 12/12/40 | 0/0/0 | Cópia de array, realocação de `append`, remoção/limpeza de map, builder e contrato de submatch já entregues. |
| `go-type-design` | 7/7/34 | 7/7/34 | 0/0/0 | Embedding/promotion, tags refletidas, addressability/method sets, satisfação em compilação, interface mínima e constraint com `~` já entregues. |
| `go-errors` | 5/5/35 | 5/5/35 | 0/0/0 | Tipo de erro com dados, hooks customizados `Is`/`As`/`Unwrap` e sucesso nil já entregues. |
| `go-io` | 5/5/19 | 5/5/19 | 0/0/0 | Buffer/token, short write, limite explícito de Reader e ausência/`null`/zero em JSON já entregues; variantes não contam. |
| `go-testing` | 4/4/11 | 5/5/18 | +1/+1/+7 | Helper, cleanup, propriedade de fuzz e registro de lifecycle/diagnóstico ainda precisam de decomposição adicional; variantes não contam. |
| **Total creditável** | **39/39/202 = 280** | **45/45/214 = 304** | **+6/+6/+12 = +24** | **Faltam 24 nodes de bases para a meta; variantes permanecem separadas.** |

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
- A projeção original de 108 nodes adicionais será distribuída nas bases já
  existentes e nas duas fatias funcionais, sem novos desafios canônicos. Os
  lotes já entregaram 84 nodes creditáveis desde o baseline anterior; restam 24
  para a meta. A decisão histórica de 44 desafios, cinco trilhas e oito
  variantes permanece intacta.

## 6. Auditoria das cinco trilhas (`audit-tracks`)

As cinco trilhas vivem em `packs/go-first-steps.yaml` versão 1.8.0 e reutilizam
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

1. **Congelar o ledger:** manter o baseline acima por pack e por
   `macro/meso/micro`; manter uma linha separada para o protótipo histórico e
   outra para as oito variantes. A disposição dos dois protótipos continua uma
   decisão humana da spec, sem quebrar IDs de sessões.
2. **Conceitos e competências:** os 34 conceitos e quatro competências das
   seções 3 e 4 estão reconciliados; manter exemplos, relações e referências
   ligados ao conteúdo real ao autorar os nodes restantes.
3. **Completar árvores de bases:** aplicar a matriz da seção 5; recontar
   macro/meso/micro após cada pack e alcançar pelo menos 304 creditáveis no
   conjunto. Variantes podem receber revisão de equivalência, mas não podem
   reduzir o trabalho das bases nem pagar o déficit.
4. **Verificar:** executar, no estado atual,
   `go run ./cmd/codinho catalog validate --json` e
   `go run ./cmd/codinho catalog validate --checks --json`. Registrar
   `declared == verified` para checks e guardar as falhas esperadas de baseline
   contra referências aprovadas. Rodar também `go test ./...` e os exemplos
   isolados; não confundir check aprovado com revisão pedagógica.
5. **Pré-revisão e humano:** submeter o lote à pré-revisão automatizada
   independente e depois ao revisor humano. Fazer playtest ponta a ponta com
   `workspace prepare` e sessão real antes de qualquer publicação.

## 8. Decisões históricas e pendências preservadas

- A árvore é reutilizada por granularidade e por trilhas; não criar packs por
  senioridade nem duplicar desafios para obter mais nodes.
- O conteúdo histórico mantém IDs, autores e fixtures. Em particular,
  `go-data.slice-filter-preserve-input` não será apagado por estar além da
  distribuição prevista, e `go-debug.slice-off-by-one` permanece fora do
  escopo fundamental.
- As oito variantes continuam contextualizadas, mas com `variant_of` e
  `canonical: false`; não serão contadas como novos desafios, competências ou
  profundidade das bases.
- Nenhum `publication.reviewed_by`, `playtested: true` ou
  `status: published` será preenchido por autoria automatizada. O checklist
  exige pessoa diferente do autor e playtest honesto; ambos estão pendentes.
- A spec continua `in-progress` até a disposição dos protótipos, a revisão de
  duplicações/leaks, a confirmação de 304 nodes no ledger e os gates humanos;
  104 conceitos e 60 competências já estão reconciliados.
  O `--v1-gate` global é uma etapa posterior e seus limiares maiores não devem
  ser alegados como satisfeitos por este plano.

## 9. Evidência desta tarefa

Nesta tarefa foram editados somente os três paths autorizados: a membership da
trilha em `packs/go-first-steps.yaml`, a prova MCP em
`cmd/codinho/foundation_tracks_integration_test.go` e este documento. A medição
que sustenta a reconciliação foi:

- `go run ./cmd/codinho catalog validate --json`: catálogo estruturalmente
  carregável; 14 temas, 105 conceitos, 64 competências, 5 trilhas, 54
  entradas de desafio e 343 step nodes no manifest.
- Contagem reconciliada: sete packs com 104 conceitos, 60 competências, 45
  bases e 8 variantes; 281 nodes de bases, 280 creditáveis sem o micro-node
  histórico, e 56 nodes de variantes. O protótipo `go-debugging` permanece
  separado do escopo fundamental.
- As cinco trilhas somam 47 posições e 37 IDs distintos; o percurso `go-from-zero`
  está em `go-first-steps` 1.8.0 e usa o desafio canônico de filtro.

Checks executados neste lote:

- `go run ./cmd/codinho catalog validate --json`: passou, sem diagnósticos ou
  findings editoriais; o catálogo atual permanece em draft.
- `go run ./cmd/codinho catalog validate --checks --json`: passou com 54 checks
  declarados e 54 verificados; cada baseline e referência teve o resultado
  esperado.
- `env GOFLAGS=-buildvcs=false go test ./cmd/codinho -run
  '^TestFoundationTracksOverRealStdio$' -count=1 -v`: passou; cada uma das
  cinco trilhas iniciou e alcançou todos os seus IDs esperados via MCP stdio.
- `pose docs-check`: passou (`declared=16`, `undeclared=0`, `stale=0`,
  `errors=0`, `warnings=0`). O comando ainda reporta dois avisos
  `REVIEW_PENDING` históricos em `docs/compatibility.md`, fora deste path.
- `env GOFLAGS=-buildvcs=false go test ./...` e
  `env GOFLAGS=-buildvcs=false go test -race -count=1 ./...`: passaram; também
  passaram `env GOFLAGS=-buildvcs=false go vet ./...` e
  `bash scripts/ci/check-format.sh`.
- `env GOFLAGS=-buildvcs=false pose validate --strict --stack go`: a matriz
  executou os checks de Go, docs, catálogo, integrações e Python, mas terminou
  com falha requerida porque `govulncheck ./...` não está instalado no ambiente
  (`command not found`).
- `git diff --check`: passou.

Uma execução sem `GOFLAGS=-buildvcs=false` falhou somente nos testes de
integração que fazem build do binário: o checkout local devolveu exit 128 ao
obter o status VCS. A flag desabilita apenas o stamping VCS do build e deixou a
mesma suíte passar; essa limitação ambiental permanece registrada para o
revisor. A aprovação mecânica não completa os 24 nodes restantes nem comprova
revisão humana/playtest.
