---
title: Catálogo fundamental de Go
doc_type: reference
---

# Plano de expansão não redundante do catálogo fundamental de Go

Tarefa editorial: `plan-concept-depth`. Inventário auditado em 2026-09-11 UTC
contra o manifest atual, depois das auditorias listadas como dependências desta
tarefa. Este documento é um plano de autoria e não uma declaração de publicação:
a spec `go-foundations-packs` continua `in-progress`, todos os sete packs
continuam em rascunho e revisão humana/playtest continuam pendentes.

## 1. Baseline reconciliado

Os números abaixo vêm de `go run ./cmd/codinho catalog validate --json` e de
uma contagem independente das árvores `layers[*].macro_steps`, recursiva. O nó
de `layer` não entra na contagem. `Base` inclui o conteúdo histórico; entre
parênteses, quando necessário, aparece o crédito conservador usado para a
meta: o protótipo histórico `go-data.slice-filter-preserve-input` tem um
micro-node, mas esse node não é creditado antes de sua disposição explícita.

| Pack | Temas A–M | Bases atuais | Conceitos | Competências | Nodes de base M/Meso/Micro | Nodes de variantes |
|---|---|---:|---:|---:|---:|---:|
| `go-first-steps` | A, B, com E/C introdutórios | 7 | 17 | 8 | 0/0/22 (crédito 0/0/21) | 0 |
| `go-core` | C, D, K | 10 | 9 | 8 | 0/0/33 | 0 |
| `go-data-text` | E, F | 10 | 8 | 8 | 10/10/34 | 0 |
| `go-type-design` | G, H, I | 8 | 6 | 6 | 0/0/30 | 0 |
| `go-errors` | J | 7 | 12 | 10 | 2/2/32 | 0 |
| `go-io` | L | 2 bases | 10 | 11 | 2/2/8 | 4 variantes: 4/4/16 |
| `go-testing` | M | 1 base | 8 | 5 | 1/1/8 | 4 variantes: 4/4/21 |
| **Sete packs fundamentais** | **A–M** | **45 entradas** | **70** | **56** | **15/15/167 = 197** (**crédito 196**) | **8 variantes = 53** |
| `go-debugging` (fora do escopo) | protótipo | 1 | 1 | 4 | 1/1/4 = 6 | 0 |
| **Manifest completo** | — | **46 entradas** | **71** | **60** | **16/16/171 = 203 de base; 24/24/208 = 256 com variantes** | **53** |

As 45 entradas fundamentais preservam a distribuição histórica de 33
atômicos, 10 combinados e 2 fatias funcionais. A diferença em relação à meta
de 44 é somente a base histórica de filtro em `go-first-steps`; ela não será
apagada, renomeada ou silenciosamente reclassificada. As quatro variantes de
I/O e as quatro de testes mantêm `variant_of` e `canonical: false` e não
financiam a meta de profundidade. O protótipo de depuração e suas quatro
competências também não financiam a meta fundamental.

### Regra de crédito

Há três projeções que devem permanecer visíveis em toda auditoria:

1. **Escopo fundamental:** sete packs, 70 conceitos e 56 competências hoje.
2. **Crédito de profundidade:** nodes de bases, sem as oito variantes e sem o
   micro-node do protótipo histórico até a decisão de disposição.
3. **Manifesto:** inclui `go-debugging` e as variantes para fins de inventário,
   mas não transforma esses itens em conteúdo fundamental novo.

Assim, a meta deste plano é **104 conceitos e 60 competências nos sete packs**
e **304 nodes de bases creditáveis**, já com margem sobre o mínimo de 300.
Com `go-debugging` preservado, o manifesto projetado teria 105 conceitos, 64
competências e 310 nodes de base; os 53 nodes de variantes continuam
reportados separadamente. A margem garante a meta mesmo se o protótipo
histórico permanecer fora do crédito.

## 2. Mapa dos 70 conceitos existentes e lacunas reais

O mapeamento usa um tema primário por conceito; alguns conceitos têm relação
secundária indicada na lacuna. A enumeração soma exatamente os 70 conceitos
atuais, sem contar `go-debugging`.

| Tema | Pack(s) e conceitos existentes | Lacuna que precisa de decisão editorial |
|---|---|---|
| **A — Ambiente e tooling** | `go-first-steps`: `go-modules`, `go-version-directive`, `printf-verb-matching` | Falta observar o contexto efetivo de `go env`, tornar formatação repetível e distinguir seleção de pacotes no `go build`. Propostas: `go-toolchain-environment`, `gofmt-idempotent-source`, `go-build-package-selection`. Workspaces, tags e cross-compile ficam para um lote posterior, pois não há desafio fundamental correspondente. |
| **B — Declarações e tipos** | `go-first-steps`: `named-types`, `field-types`, `zero-values`, `short-variable-declaration`, `constant-declaration`, `iota-sequence`, `type-conversion`, `narrowing-conversion-overflow`, `named-return-values`, `variable-shadowing`, `block-scope` | Falta separar tipo padrão de constante não tipada e representabilidade em tempo de compilação de conversão em runtime. Propostas: `untyped-constant-defaulting`, `numeric-constant-representability`. |
| **C — Controle de fluxo** | `go-first-steps`: `range-loop`; `go-core`: `defer-execution-order`, `labeled-break-scope`, `short-circuit-nil-guard` | `range` já existe, mas faltam escopo de inicialização de `if`, ciclo completo de `for` e avaliação dos argumentos de `defer`. Propostas: `if-initializer-scope`, `for-clause-lifecycle`, `defer-argument-evaluation`. |
| **D — Funções e métodos** | `go-core`: `variadic-parameters`, `pointer-vs-value-receiver`, `panic-recover-boundary` | Faltam captura independente de closures e semântica de method value. `pointer-vs-value-receiver` é reutilizado como relação secundária em G; `panic-recover-boundary` também se relaciona com J, sem criar cópia. Propostas: `closure-capture-lifetime`, `method-value-binding`. |
| **E — Arrays, slices e maps** | `go-first-steps`: `slice-declaration`; `go-data-text`: `slice-backing-array-aliasing`, `map-comma-ok-lookup`, `set-membership-with-empty-struct-map`, `slice-preallocation-length-vs-capacity` | Não há decisão explícita sobre semântica de cópia de array, realocação de `append` ou `delete`/`clear` de map. Propostas: `array-value-copy`, `append-reallocation-boundary`, `map-delete-and-clear`. |
| **F — Texto e dados binários** | `go-data-text`: `utf8-rune-boundary`, `csv-field-quoting`, `strconv-error-handling`, `rune-aware-whitespace-preservation` | UTF-8, CSV e parsing já têm cobertura, mas faltam construção incremental de texto e contrato de submatches. Propostas: `strings-builder-output`, `regexp-submatch-contract`. |
| **G — Structs e ponteiros** | `go-first-steps`: `structs`; `go-type-design`: `struct-shallow-copy-aliasing`, `nil-safe-pointer-receiver` | Faltam embedding/promotion, tags observáveis e o efeito da addressability sobre method sets. Propostas: `struct-embedding-promotion`, `struct-tags-reflection`, `method-set-addressability`. `pointer-vs-value-receiver` e `unexported-field-invariant` entram como relações cruzadas, não como conceitos duplicados. |
| **H — Interfaces** | `go-type-design`: `typed-nil-interface-trap`, `safe-type-assertion` | A armadilha de nil e a assertion segura existem, mas faltam satisfação verificada em compilação e composição por consumidor mínimo. Propostas: `compile-time-interface-satisfaction`, `interface-composition-consumer`. |
| **I — Generics** | `go-type-design`: `generic-comparable-constraint`, `generic-two-type-parameters` | Falta o contraste entre tipo definido e seu tipo subjacente usando `~`, incluindo o limite da constraint. Proposta: `generic-underlying-type-constraint`; não criar outro conceito de “comparável”. |
| **J — Erros** | `go-errors`: `error-wrapping-with-percent-w`, `errors-as-type-extraction`, `errors-join-aggregation`, `error-translation-at-boundary`, `error-classification-precedence`, `error-identity-not-message`, `error-chain-type-traversal`, `error-joined-tree-order`, `catalog-zero-value`, `catalog-validation-before-mutation`, `catalog-deterministic-list`, `catalog-value-snapshots` | Wrapping, busca, agregação e catálogo já têm profundidade; faltam contratos de tipos de erro customizados, hooks customizados de `Is`/`As`/`Unwrap` e retorno nil de sucesso. Propostas: `custom-error-type-contract`, `errors-is-custom-matcher`, `errors-as-custom-matcher`, `errors-unwrap-method-contract`, `nil-error-success-contract`. |
| **K — Packages e APIs** | `go-core`: `unexported-field-invariant`, `constructor-returns-interface`, `functional-options-pattern` | Há visibilidade de campo e encapsulamento por construtor, mas falta uma fronteira entre package produtor e consumidor externo com exportação mínima. Proposta: `package-import-export-boundary`; ela se relaciona a `go-modules` sem repetir o conceito de módulo. |
| **L — I/O e serialização** | `go-io`: `io-reader-stream`, `io-partial-reads`, `io-eof-boundary`, `io-json-decoder`, `io-json-strict-fields`, `io-csv-quoted-fields`, `io-writer-errors`, `io-validate-before-write`, `io-csv-record-width`, `io-json-encode` | Reader/Writer, JSON e CSV estão cobertos, mas faltam limite de token com buffering, short write, limite explícito de leitura e distinção JSON entre campo ausente, `null` e zero. Propostas: `bufio-scanner-token-boundary`, `io-short-write-contract`, `io-limit-reader-boundary`, `json-null-versus-absent`. |
| **M — Testes** | `go-testing`: `testing-table-cases`, `testing-named-subtests`, `testing-clock-injection`, `testing-counting-double`, `testing-time-boundaries`, `testing-instant-versus-zone`, `testing-input-isolation`, `testing-repeatability` | Tempo, fronteiras, doubles e isolamento estão cobertos; faltam contexto de helper, lifecycle de cleanup e invariantes de fuzz. Paralelismo, golden, benchmark e race detector continuam temas posteriores, não serão simulados por sinônimos neste lote. Propostas: `testing-helper-context`, `testing-cleanup-lifecycle`, `testing-fuzz-invariant`. |

## 3. Distribuição dos 34 conceitos novos

Cada linha define um conteúdo futuro, não apenas um título: haverá um
`content.example` em um contexto diferente dos desafios, uma relação inicial
no grafo, uma competência/critério observável e a referência primária a
consultar. O exemplo será um mini-módulo ou teste Go executável com biblioteca
padrão, Go mínimo 1.25, sem rede e sem arquivos externos. O comando de prova
de cada exemplo será `go test ./... -count=1` no diretório temporário da
fixture; quando a decisão for de compilação ou tooling, o plano registra
também o comando determinístico de `go/types`, `go env`, `gofmt` ou `go build`.

### `go-first-steps` — 17 → 22 conceitos

| Novo conceito (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `go-toolchain-environment` (`inspect-go-toolchain-context`) | Mini-módulo que registra `go env GOMOD GOOS GOARCH` e compara o build com `runtime.GOOS`/`runtime.GOARCH`. | `go-modules` — `relates_to` | O teste identifica o módulo e o alvo efetivos sem depender do diretório corrente ou alterar o ambiente. | [CMD], [MOD], [RUNTIME] |
| `gofmt-idempotent-source` | Arquivo Go deliberadamente desalinhado; `gofmt -d` após formatação deve produzir diff vazio e a segunda execução não pode mudar bytes. | `printf-verb-matching` — `relates_to` | Mede idempotência da ferramenta; não confunde formato com correção de verbos de `fmt`. | [EG], [CMD] |
| `go-build-package-selection` | Mini-módulo com library, `cmd/` e package interno; `go list ./...` e `go build ./...` devem enumerar e compilar os três contextos permitidos. | `go-modules` — `relates_to`; `package-import-export-boundary` — `deepens_into` | Verifica seleção/compilação de packages, não só a existência de `go.mod`. | [S], [CMD], [MOD] |
| `untyped-constant-defaulting` | Teste atribui a mesma constante não tipada a `int`, `float64`, `rune` e `string`, verificando tipo estático por compilação e valor observado. | `constant-declaration` — `relates_to`; `type-conversion` — `contrasts_with` | O resultado nasce do contexto de uso; não há conversão explícita para mascarar o default. | [S] |
| `numeric-constant-representability` | `go/types` verifica snippets: uma constante representável passa e uma que excede o tipo falha antes da execução. | `narrowing-conversion-overflow` — `contrasts_with`; `type-conversion` — `relates_to` | Distingue erro de compilação de wraparound em conversão executada. | [S], [TYPES] |

### `go-core` — 9 → 15 conceitos

| Novo conceito (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `if-initializer-scope` | Teste de uma decisão com inicializador local e um snippet em `go/types` que falha ao usar o nome depois do `if`. | `block-scope` — `deepens_into`; `short-circuit-nil-guard` — `relates_to` | A variável do inicializador existe somente na condição e nos ramos do `if`. | [S] |
| `for-clause-lifecycle` | Tabela executa `init`, condição, corpo, `continue` e `post` com um trace esperado, sem usar `range`. | `range-loop` — `contrasts_with`; `labeled-break-scope` — `relates_to` | O trace prova a ordem do `post` inclusive quando há `continue`; não é apenas uma contagem de iterações. | [S] |
| `defer-argument-evaluation` | Função adia uma função com argumento escalar e outra closure; o trace verifica o valor capturado na inscrição e o valor lido na chamada. | `defer-execution-order` — `deepens_into` | Isola tempo de avaliação do argumento, sem repetir a ordem LIFO dos `defer` existentes. | [S] |
| `closure-capture-lifetime` (`control-closure-capture`) | Fábrica retorna duas closures independentes; chamadas intercaladas devem manter contadores separados e estado após o retorno da fábrica. | `block-scope` — `relates_to`; `named-return-values` — `relates_to` | Observa ambiente capturado e independência entre instâncias, não apenas uma função anônima que retorna valor. | [S], [EG] |
| `method-value-binding` | Teste cria method values antes e depois de alterar um receiver e compara o efeito para receiver por valor e por ponteiro. | `pointer-vs-value-receiver` — `deepens_into` | Verifica quando o receiver é ligado ao method value; não repete a escolha de receiver do desafio-base. | [S] |
| `package-import-export-boundary` (`define-package-api-boundary`) | Mini-módulo tem teste em package externo: somente identificadores exportados compilam, e o consumidor usa a menor interface pública necessária. | `constructor-returns-interface` — `relates_to`; `go-modules` — `relates_to` | A fronteira é observada por um consumidor de outro package; não é apenas campo não exportado no mesmo package. | [S], [MOD], [EG] |

### `go-data-text` — 8 → 13 conceitos

| Novo conceito (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `array-value-copy` | Teste atribui `[3]int`, altera a cópia e contrasta com uma slice sobre o mesmo conteúdo. | `slice-backing-array-aliasing` — `contrasts_with`; `slice-declaration` — `relates_to` | A atribuição de array copia valores; a comparação não permite concluir que toda slice seja cópia. | [S] |
| `append-reallocation-boundary` | Slice com capacidade conhecida sofre `append` antes e depois de cruzar a capacidade; o teste observa quando a mutação deixa de compartilhar backing array. | `slice-backing-array-aliasing` — `deepens_into`; `slice-preallocation-length-vs-capacity` — `relates_to` | O critério separa aliasing enquanto há capacidade de realocação e independência depois dela. | [S] |
| `map-delete-and-clear` | Teste insere chaves, usa `delete` para uma só, `clear` para todas e consulta mapa nil sem panic. | `map-comma-ok-lookup` — `deepens_into`; `set-membership-with-empty-struct-map` — `relates_to` | Verifica remoção seletiva e limpeza total, não apenas presença/ausência por comma-ok. | [S], [MAPS] |
| `strings-builder-output` | `strings.Builder` recebe fragmentos, `String` é comparada com a concatenação esperada e o caso vazio é verificado. | `rune-aware-whitespace-preservation` — `relates_to`; `utf8-rune-boundary` — `relates_to` | O contrato é acumulação e saída, não normalização de espaços nem corte de runes. | [STR] |
| `regexp-submatch-contract` | Teste usa uma regex compilada em linhas de inventário e verifica match completo, submatch e ausência (`nil`) separadamente. | `strconv-error-handling` — `relates_to`; `rune-aware-whitespace-preservation` — `relates_to` | Distingue ausência de match e grupo vazio sem transformar regex em mais um conceito de parsing numérico. | [REGEXP] |

### `go-type-design` — 6 → 12 conceitos

| Novo conceito (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `struct-embedding-promotion` | Struct externa embute uma struct de endereço; teste acessa campo/método promovido e registra a ambiguidade quando dois campos têm o mesmo nome. | `struct-shallow-copy-aliasing` — `relates_to`; `structs` — `deepens_into` | O foco é seleção/promoção de membros, não cópia independente de slices. | [S] |
| `struct-tags-reflection` | `reflect.Type.FieldByName` e `StructTag.Get` leem uma tag `json` de um registro de catálogo. | `structs` — `relates_to`; `field-types` — `relates_to` | A tag altera metadado de serialização observado por reflexão; não altera nome nem tipo do campo. | [S], [REFLECT] |
| `method-set-addressability` | Teste compila chamada de método de ponteiro em valor endereçável e usa `go/types` para rejeitar a mesma chamada em elemento não endereçável de map. | `pointer-vs-value-receiver` — `deepens_into`; `nil-safe-pointer-receiver` — `relates_to` | O critério é addressability/method set, não somente o comportamento de receiver nil. | [S], [TYPES] |
| `compile-time-interface-satisfaction` | Atribuição de uma implementação a uma interface compila; um snippet mutante sem método é rejeitado por `go/types`. | `safe-type-assertion` — `contrasts_with`; `typed-nil-interface-trap` — `relates_to` | A falha é descoberta na compilação, antes de qualquer type assertion ou valor nil. | [S], [TYPES] |
| `interface-composition-consumer` | Package consumidor aceita uma interface de um método; um fake mínimo passa no teste sem conhecer o tipo concreto do produtor. | `constructor-returns-interface` — `deepens_into`; `safe-type-assertion` — `relates_to` | Mede dependência no menor contrato do consumidor, não só satisfação acidental de uma interface grande. | [S], [EG] |
| `generic-underlying-type-constraint` | Função com constraint `~int` aceita `type Score int` e rejeita `string`; o teste instancia ambos e mantém o erro no caso inválido. | `generic-comparable-constraint` — `deepens_into`; `generic-two-type-parameters` — `relates_to` | O critério é aceitar tipos definidos pelo tipo subjacente, sem duplicar a constraint `comparable`. | [S] |

### `go-errors` — 12 → 17 conceitos

| Novo conceito (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `custom-error-type-contract` | Tipo `FieldError` carrega campo e valor inválido; teste verifica `Error` e extrai os dados com `errors.As`. | `errors-as-type-extraction` — `deepens_into`; `error-identity-not-message` — `relates_to` | O contrato do tipo inclui dados estruturados; a mensagem não é usada como protocolo. | [ERRORS] |
| `errors-is-custom-matcher` | Erro wrapper implementa `Is` para uma categoria estável; duas instâncias com mensagens diferentes devem satisfazer `errors.Is` somente pela categoria declarada. | `error-identity-not-message` — `deepens_into`; `error-classification-precedence` — `relates_to` | A equivalência semântica é explícita no método `Is`, sem comparar `Error()`. | [ERRORS] |
| `errors-as-custom-matcher` | Wrapper implementa `As` para expor apenas um tipo público; o teste verifica match do tipo permitido e ausência de match de um tipo interno. | `errors-as-type-extraction` — `deepens_into`; `error-translation-at-boundary` — `relates_to` | O limite do tipo exposto é controlado pelo método `As`, não pela profundidade casual da cadeia. | [ERRORS] |
| `errors-unwrap-method-contract` | Erro customizado implementa `Unwrap`; `errors.Is` atravessa uma causa, enquanto `Unwrap` nil encerra a cadeia. | `error-chain-type-traversal` — `deepens_into`; `errors-join-aggregation` — `relates_to` | Testa contrato de uma causa linear e seu término, sem recontar a ordem de árvore de `Join`. | [ERRORS] |
| `nil-error-success-contract` | Função de validação retorna `(value, nil)` em sucesso e um erro não nil em falha; o teste inclui entrada vazia e sucesso parcial. | `error-translation-at-boundary` — `relates_to`; `errors-join-aggregation` — `relates_to` | O critério é a identidade nil do caminho de sucesso, não a aparência da mensagem nem o typed nil de interface. | [ERRORS], [S] |

### `go-io` — 10 → 14 conceitos

| Novo conceito (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
|---|---|---|---|---|
| `bufio-scanner-token-boundary` | Scanner lê linhas curtas e uma entrada que excede o buffer configurado; o teste distingue tokens lidos de `ErrTooLong`. | `io-reader-stream` — `deepens_into`; `io-partial-reads` — `relates_to` | A fronteira é do token/buffer, não uma leitura parcial genérica ou EOF. | [BUFIO], [IO] |
| `io-short-write-contract` | Writer de teste aceita somente um prefixo; a operação verifica `n < len(p)` e preserva `io.ErrShortWrite`/erro declarado. | `io-writer-errors` — `deepens_into`; `io-validate-before-write` — `relates_to` | O caso cobre escrita parcial e sua sinalização, não somente Writer que rejeita tudo. | [IO] |
| `io-limit-reader-boundary` | `io.LimitReader` recebe mais bytes que o limite; o teste lê exatamente o prefixo e confirma que o limite termina a leitura observável. | `io-reader-stream` — `relates_to`; `io-eof-boundary` — `deepens_into` | O limite imposto pelo consumidor é separado do EOF originado pela fonte. | [IO] |
| `json-null-versus-absent` | Decodifica três documentos usando presença em `map[string]json.RawMessage` e valor opcional: campo ausente, `null` e zero explícito produzem estados distinguíveis. | `io-json-decoder` — `deepens_into`; `io-json-strict-fields` — `relates_to` | Verifica semântica de presença/nulidade, não apenas rejeição de campo desconhecido. | [JSON] |

### `go-testing` — 8 → 11 conceitos

| Novo conceito (competência) | Exemplo executável | Relação inicial | Critério distinto | Ref. |
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

## 4. Quatro competências fundamentais ausentes

O baseline fundamental é 56. A meta é adicionar exatamente as quatro
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

Distribuição resultante: `go-first-steps` 9, `go-core` 10, `go-data-text` 8,
`go-type-design` 6, `go-errors` 10, `go-io` 11 e `go-testing` 6, totalizando
60. `go-debugging` permanece em 4 e aparece somente na projeção global.

## 5. Nodes e decisões pedagógicas que faltam

O alvo abaixo mede somente árvores de bases, excluindo variantes e não
creditando o micro-node do protótipo histórico. Macro e meso novos podem
organizar micros existentes quando a decisão já está presente, mas nenhum
micro novo pode repetir uma intenção existente. As adições são números de
nodes `macro/meso/micro`; cada micro futuro terá um verbo, um alvo, uma
evidência e um critério observável.

| Pack | Crédito atual M/Meso/Micro | Alvo M/Meso/Micro | Adição | Decisões ausentes a autorar |
|---|---:|---:|---:|---|
| `go-first-steps` | 0/0/21 | 6/6/30 | +6/+6/+9 | Contexto da toolchain; formatação/build observáveis; default de constantes; representabilidade antes da conversão. |
| `go-core` | 0/0/33 | 5/5/38 | +5/+5/+5 | Escopo de inicializador; ciclo de `for`; avaliação de argumentos de `defer`; captura de closure; ligação de method value e package consumidor. |
| `go-data-text` | 10/10/34 | 12/12/40 | +2/+2/+6 | Cópia de array versus slice; realocação de `append`; remoção/limpeza de map; builder; contrato de submatch. |
| `go-type-design` | 0/0/30 | 7/7/34 | +7/+7/+4 | Embedding/promotion; tags refletidas; addressability e method sets; satisfação em compilação; interface mínima; constraint com `~`. |
| `go-errors` | 2/2/32 | 5/5/35 | +3/+3/+3 | Tipo de erro com dados; hooks customizados `Is`/`As`/`Unwrap`; sucesso nil. |
| `go-io` | 2/2/8 | 5/5/19 | +3/+3/+11 | Buffer/token; short write; limite explícito de Reader; ausência/`null`/zero em JSON. Somente as duas bases, nunca as quatro variantes. |
| `go-testing` | 1/1/8 | 5/5/18 | +4/+4/+10 | Helper; cleanup; propriedade de fuzz; registro de lifecycle e diagnóstico. Somente a base combinada, nunca as quatro variantes. |
| **Total creditável** | **15/15/166 = 196** | **45/45/214 = 304** | **+30/+30/+48 = +108** | **304 nodes de bases, com margem de quatro sobre a meta.** |

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
- Os 108 nodes adicionais serão distribuídos nas bases já existentes e nas
  duas fatias funcionais, sem novos desafios canônicos. A decisão histórica de
  44 desafios, cinco trilhas e oito variantes permanece intacta.

## 6. Sequência de execução e gates

1. **Congelar o ledger:** registrar o baseline acima por pack e por
   `macro/meso/micro`; manter uma linha separada para o protótipo histórico e
   outra para as oito variantes. A disposição dos dois protótipos continua uma
   decisão humana da spec, sem quebrar IDs de sessões.
2. **Autorar conceitos por pack:** implementar os 34 itens nas respectivas
   `packs/*.yaml`, com `content.version: 1`, explicação, exemplo, relações e
   referências. Executar cada exemplo isolado e revisar duplicação semântica
   antes de contar o conceito.
3. **Autorar competências e relações:** adicionar somente as quatro
   competências da seção 4; ligar conceitos e desafios existentes por
   `applies_in` quando houver contexto real. Recalcular 104/60 antes de iniciar
   a etapa de nodes.
4. **Completar árvores de bases:** aplicar a matriz da seção 5; recontar
   macro/meso/micro após cada pack e exigir pelo menos 304 creditáveis no
   conjunto. Variantes podem receber revisão de equivalência, mas não podem
   reduzir o trabalho das bases nem pagar o déficit.
5. **Verificar:** executar, no estado que contiver a autoria futura,
   `go run ./cmd/codinho catalog validate --json` e
   `go run ./cmd/codinho catalog validate --checks --json`. Registrar
   `declared == verified` para checks e guardar as falhas esperadas de baseline
   contra referências aprovadas. Rodar também `go test ./...` e os exemplos
   isolados; não confundir check aprovado com revisão pedagógica.
6. **Pré-revisão e humano:** submeter o lote à pré-revisão automatizada
   independente e depois ao revisor humano. Fazer playtest ponta a ponta com
   `workspace prepare` e sessão real antes de qualquer publicação.

## 7. Decisões históricas e pendências preservadas

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
  duplicações/leaks, a confirmação de 104/60/304 no ledger e os gates humanos.
  O `--v1-gate` global é uma etapa posterior e seus limiares maiores não devem
  ser alegados como satisfeitos por este plano.

## 8. Evidência desta tarefa

Nesta tarefa foi editado somente este documento. A medição que sustenta o plano
foi:

- `go run ./cmd/codinho catalog validate --json`: catálogo estruturalmente
  carregável; 71 conceitos, 60 competências, 54 entradas de desafio no
  inventário e 256 step nodes no manifest.
- Contagem por pack: sete packs com 70/56 e 250 nodes incluindo variantes;
  crédito conservador de 196 nodes nas bases, mais 53 nodes de variantes e
  um micro-node histórico não creditado.

Checks executados neste lote:

- `go run ./cmd/codinho catalog validate --json`: passou, sem diagnósticos ou
  findings editoriais; o catálogo atual permanece em draft.
- `go run ./cmd/codinho catalog validate --checks --json`: passou com 54 checks
  declarados e 54 verificados; cada baseline e referência teve o resultado
  esperado.
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
revisor. A aprovação mecânica não completa os conceitos propostos nem
comprova revisão humana/playtest.
