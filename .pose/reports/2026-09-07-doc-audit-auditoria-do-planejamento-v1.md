# Auditoria do planejamento V1 — 2026-09-07 UTC

## Scope

Revisão de planejamento solicitada pelo mantenedor: PROJECT.md, as 29 specs
preexistentes, o roadmap codinho-v1, seis ADRs, conhecimentos ativos, matriz
de validação, políticas POSE, CI, catálogo e caminhos relevantes do código.
Baseline: HEAD d7587442bc9f, com alteração preexistente em packs/go-errors.yaml.
A data UTC equivale a 2026-09-06 em America/Recife durante esta execução.

Resultado: planejamento revisado; decisão técnica sobre prontidão V1:
**changes requested**. Nenhuma spec de implementação foi encerrada, nenhum
bundle histórico foi reescrito e nenhuma publicação/commit foi realizado.
Os gates verdes atuais não cobrem todas as promessas do produto.

## Estado observado

Há um módulo Go, CLI executável, servidor MCP stdio, casos de uso, domínio,
event store, executor, observação de workspace e projeção de mastery.
A descoberta POSE registrou 10.002 LOC de produção e 8.931 LOC de testes
(heurística do engine, não indicador de completude).

Antes da revisão: 29 specs, sendo 22 done, uma in-progress e seis draft.
Depois: 34 specs, mantendo as 22 done e uma in-progress; 11 draft.
Status done é histórico de escopo/revisão, não comprovação automática de
comportamento composto hoje.

| Dimensão do catálogo do worktree | Inventário atual | Gate V1 |
|---|---:|---:|
| Conceitos | 45 | pelo menos 160 |
| Competências | 35 | pelo menos 100 |
| Desafios carregados | 41 | 84 curados, na distribuição definida |
| Desafios com publication.status published | 0 | todos os contabilizados |
| Trilhas | 1 | pelo menos 12 |
| Step nodes | 89 | pelo menos 500 |

Os 41 itens são 33 atomic, sete combined e um debug. Incluem dois protótipos
históricos e um combinado de go-errors ainda no diff do usuário. A autoria
fundamental canônica é 38 no HEAD e 39 no worktree. Os números de conceitos,
competências e nodes também incluem conteúdo sem publicação; não representam
cobertura editorial entregue. Não há evidência de revisão humana inferida a
partir de campos de autoria ou pré-revisão automatizada.

## Findings

### F1 — high: reinício quebra a sessão e impede novo início

Reprodução em diretório temporário, com cópia do catálogo e dois processos
reais de codinho serve via JSON-RPC stdio:

1. Inicializar MCP e chamar session_start com request_id audit-start.
2. Receber status ok, session_id ses_1, revision 1.
3. Fechar stdin e aguardar término limpo do processo.
4. Iniciar outro processo sobre o mesmo diretório.
5. session_get ses_1 retorna SESSION_NOT_ACTIVE.
6. session_start com request_id diferente retorna STATE_CONFLICT.

Evidência: internal/session/service.go New inicializa mapas vazios e nextID;
Start aloca ses_N e persiste apenas challenge_id/mode. cmd/codinho/main.go
conecta serviços novos sem reconstruir sessões. O log confirmado existe;
a projeção pedagógica não é reidratada e o contador volta a colidir.

Encaminhamento: [session-recovery-version-pinning](../specs/2026-09-07-session-recovery-version-pinning.md).
Exigir política completa, consentimentos, eventos legados, versão recuperável
do conteúdo, baselines e escopo de evidências. O teste do event store isolado
não substitui esse cenário.

### F2 — high: navegação não cobre a árvore prometida

internal/session/progression.go rootSteps retorna a primeira layer não vazia.
granularity.go deriveWindow parte do primeiro ramo; challenge/layer podem
virar IDs ativos que findStep não resolve. Start seleciona firstStep mesmo
com profundidade micro. O problema já aparece como limitação em
feedback-evaluation-progression, mas não tinha remediação executável própria.

Encaminhamento: [session-tree-progression](../specs/2026-09-07-session-tree-progression.md).
Testar múltiplas camadas, agrupamento sem perda de posição e escolha de
ramificação pela API. Esta constatação vem da leitura de código; o probe
executado nesta rodada demonstrou restart, não toda a navegação.

### F3 — high: inventário pode passar por conteúdo entregue

internal/curriculum/index.go indexa todos os desafios; coverage.go conta esse
inventário. O filtro de publicação não está presente na busca padrão.
CheckTypeDistribution só tem chamadas nos testes, sem composição com a CLI.
O gate V1 atual verifica mínimos totais, não a distribuição exata nem
elegibilidade editorial de cada item contado.

Encaminhamento: [catalog-publication-integrity](../specs/2026-09-07-catalog-publication-integrity.md).
Separar inventário/publicação, ligar distribuição à CLI e testar catálogo
numericamente suficiente mas sem revisão. Não publicar os rascunhos atuais
para fazer os números passarem.

### F4 — high: a matriz editorial previa 52 fundamentais, não 44

A tabela de go-foundations-packs somava 6+8+8+6+4+4+4 = 40 atômicos;
com dez combinados e duas fatias, eram 52 desafios. A declaração dizia
32+10+2 = 44; o restante da V1 soma 23+13+4, totalizando 84 apenas se a
fundação entregar 44.

Correção aplicada: preservar os 32 atômicos dos cinco primeiros packs e
planejar as oito práticas adicionais de I/O/testes como variantes, sujeitas
à equivalência pedagógica, sem duplicar a contagem canônica. Não foram
removidos nem reclassificados desafios existentes. Se as competências
exigirem itens independentes, revisar explicitamente a distribuição antes
da autoria. Corrigidos também paths de diretórios fictícios para YAML e
comandos que davam a entender seleção de packs por posicionais ignorados.

A profundidade atual (89 nodes em todo catálogo) merece atenção equivalente
à quantidade de desafios: não preencher o déficit com nós mecânicos.

### F5 — high: as specs adaptativas omitiam estado e transições essenciais

learning-track-composition pressupunha Track com sequência, mas o modelo
tem somente ID/Title/Themes. Não definia aceite da composição ad hoc,
retomada, cobertura impossível nem todos os consumidores (a CLI usa Query).

agent-authored-catalog-drafts tratava eventos sem proveniência como revisados,
embora o catálogo atual já permita sessões sem publicação. Havia ambiguidade
entre bloquear todo conteúdo executável e receber fixtures, e entre validar
draft e exigir condições de publicação.

Correção aplicada: R6–R9 de trilhas e R7–R11 de rascunhos explicitam esses
contratos, compatibilidade, limites e provas negativas. Decisões de formato
e eventual revalidação pós-publicação devem atualizar os ADRs antes do código.
Não foi tomada decisão irreversível de armazenamento nesta auditoria.

### F6 — high: auxílio canônico ainda é apenas título

ConceptAuthoring e assistance.ConceptContent têm ID/Title. RF-021 e §15.7
pedem explicações canônicas, relações e exemplos contextualizáveis. O
follow-up de assistance-hints-detours não foi resolvido por
catalog-authoring-quality.

Encaminhamento: [concept-content-authoring](../specs/2026-09-07-concept-content-authoring.md).
O servidor deve entregar conteúdo autorado; o agente adapta a linguagem.
Não adicionar geração de texto ao núcleo.

### F7 — high: CI e gates declarados não equivalem à validação executada

installation-documentation-ci R4 pede pose check/validate estritos.
.github/workflows/ci.yml não os executa nem instala POSE; skills-check pula
quando a ferramenta não está presente. A matriz POSE atual não inclui
formatação com saída de falha, race, catálogo nem checks de fixtures.
A CI executa parte desses checks separadamente, mas não constitui um único
contrato de validação consistente.

Há apenas execução Linux nesta rodada; RNF-005 exige Linux e macOS.
Windows é desejável, não requisito terminal implícito. Build cruzado de
arm64 não prova execução arm64. Actions por tags e scanner @latest não são
pins imutáveis. O input version do release é interpolado diretamente em
shell; planejar validação/variável de ambiente no mesmo hardening.

Encaminhamento: [v1-delivery-ci-assurance](../specs/2026-09-07-v1-delivery-ci-assurance.md).
Não foram executadas CI remota, instalação macOS/Windows nem hosts reais.

### F8 — high: gates de composição/proveniência não estão operacionais

.pose/policy/artifacts.json e delivery.json estão disabled. O roadmap tinha
critérios só em português livre; o parser requer Cut criteria e C<N>.
O resultado inicial foi roadmap.criteria=0. O review bundle do roadmap
também encontra cmd/ailearn/main.go histórico inexistente e nenhum bundle
selado dos milestones. Componentes explícitos e paths geram avisos de
mapeamento; a descoberta agrupa o produto em um módulo Go na raiz.

Correção aplicada: oito critérios declarativos C1–C8. Agora o check reconhece
oito e aponta referências/evidências pendentes. Habilitação cuidadosa de
políticas, mapping e reconciliação histórica pertencem a
v1-delivery-ci-assurance. Não alterar atestações antigas nem desativar gates
para contornar erros. Mesmo com políticas permissivas, a revisão semântica
continua exigindo evidência real.

### F9 — medium: checks editoriais provam infraestrutura limitada

internal/cli/editorial.go ignora desafios sem fixture/check e aceita fail e
skipped como resultados de infraestrutura. Isso não prova que o teste
selecionado existe ou que uma solução correta passa. O catálogo também não
tem gate de equivalência entre schemas JSON e loader Go.

Encaminhamento em catalog-publication-integrity: baseline esperado, solução
de referência reservada, teste inexistente/skip/timeout negativos e corpus
compartilhado de schemas. A execução de todos os checks de fixtures não foi
repetida nesta revisão; não se declara o catálogo playtested.

### F10 — medium: aceite e sequência não acompanhavam o escopo

A spec terminal não dependia de learning-track-composition nem
agent-authored-catalog-drafts. O roadmap colocava autoria atrás dessas
features e hardening atrás de todo o currículo, embora já houvesse autoria
em andamento e hardening done. O plano poderia atrasar feedback real e
misturar preferência de agenda com pré-requisito técnico.

Correção aplicada: incluir remediações e features adaptativas como
dependências reais do aceite; separar autoria, continuidade e integridade
editorial; antecipar CI e piloto incremental. O gate terminal do roadmap
ocorre depois de fechar o aceite, evitando depender do próprio done para
produzir sua evidência.

### F11 — medium: memória e follow-ups escondiam o estado real

O resumo curado de project-state dizia não existir aplicação; o refresh
nativo preserva texto curado, portanto sozinho não corrigiria isso.
Havia seis follow-ups abertos sem owner/SLA e dois marcadores [open] vazios
nas specs adaptativas. Alguns Known gaps históricos já foram corrigidos
por specs posteriores, mas outros permaneceram sem dono executável.

Correção aplicada: resumo curado atualizado, seis follow-ups com owner
@oseiaspereira e triagem em 2026-09-21, handoff e referências de destino.
Mantidas as disposições open quando implementação/revisão humana segue
pendente; não se marcou risco como resolvido por criar uma spec.

## Cobertura das specs preexistentes

A tabela registra o destino de planejamento de todas as 29 specs; não é uma
nova atestação de cada implementação.

| Spec | Estado inicial | Avaliação/encaminhamento |
|---|---|---|
| architecture-decision-baseline | done | Preservar fronteiras dos ADRs; nenhum banco/LLM novo |
| go-runtime-foundation | done | Build/entrada verificados; reconciliar rename histórico |
| learning-domain-model | done | Invariantes úteis; compor com cursor e recuperação |
| catalog-schema-loader | done | Loader existe; paridade schema e snapshot durável faltam |
| local-event-store | done | Log/replay de baixo nível existem; não restauram sessão |
| eventstore-idempotency-scope | done | Correção por stream preservada; IDs de sessão pós-restart ainda falham |
| mcp-stdio-foundation | done | Contratos/stdio passam; detector POSE não inventaria tools |
| session-orchestration-disclosure | done | Fechar limites por specs de recuperação e árvore |
| assistance-hints-detours | done | Escada/detours existem; conteúdo canônico em nova spec |
| feedback-evaluation-progression | done | Separação preservada; navegação inter-layer ainda incompleta |
| workspace-observation-baselines | done | Observação/testes existem; persistir baseline e escopo |
| safe-check-executor | done | Registry/limites testados; não é sandbox forte |
| mastery-review-scheduling | done | Regras/replay existem; proveniência e piloto necessários |
| curriculum-graph-path-recommendation | done | Grafo/busca existem; seleção múltipla e trilhas continuam draft |
| tutor-skill-host-integration | done | Roteamento stdio testado; hosts reais no piloto/aceite |
| administrative-cli-fixtures | done | CLI/fixture e2e passam; comandos documentados alinhados no plano |
| learning-practice-debug-modes | done | Defaults/adaptação testados; EvidenceThreshold requer piloto |
| interview-mode | done | Mecânica testada; realismo/duração dependem de conteúdo e humanos |
| catalog-authoring-quality | done | Mecanismos existem; publicação/distribuição/composição incompletas |
| go-foundations-packs | in-progress | Corrigir aritmética, profundidade, protótipos, autoria e playtests |
| go-backend-packs | draft | Preservar 23; fixar driver SQL/hermeticidade antes da autoria |
| go-production-architecture-packs | draft | Preservar 13; tolerâncias locais e rubricas reais |
| go-interviews-pack | draft | Preservar quatro simulações e oito playtests independentes |
| security-privacy-hardening | done | Revisão humana de threat model ainda obrigatória |
| reliability-observability-compatibility | done | Separar recuperação de log de sessão; macOS obrigatório |
| installation-documentation-ci | done | CI estrita/multiplataforma e proveniência em remediação |
| learning-track-composition | draft | Ampliada com membership, cursor, cobertura e compatibilidade |
| agent-authored-catalog-drafts | draft | Ampliada com proveniência, consentimento, fixtures e quotas |
| v1-integrated-acceptance | draft | DAG completo, R14–R17 e fechamento em duas etapas |

## Sequência recomendada

1. Implementar recuperação durável e IDs sem colisão; atacar navegação da
   árvore em incremento separado, coordenando os arquivos de sessão.
2. Implementar integridade de publicação e paridade do catálogo. Preparar
   CI/POSE/macOS em paralelo e introduzir conteúdo canônico antes da expansão.
3. Executar learning-track-composition sobre sessão e catálogo corrigidos;
   depois agent-authored-catalog-drafts, sem redefinir os ADRs implicitamente.
4. Continuar autoria dos sete packs fundamentais, executar pré-revisão e
   playtest humano por lote; validar a nova distribuição e a profundidade.
5. Expandir os 23 backend, 13 produção/arquitetura e quatro entrevistas.
   Iniciar piloto limitado assim que um lote curado estiver utilizável;
   medir autonomia, clareza, duração e retorno após intervalo.
6. No candidato completo, executar aceite nos dois hosts e plataformas
   declaradas, fechar specs/milestones com revisão independente e então
   comprovar C1–C8. Preparar release governado sem publicação automática.

Evitar estimar porcentagem de conclusão por 22/29 specs: volume editorial,
gates humanos e integração têm custos diferentes. Não abrir agora specs de
GUI, nuvem, SQLite ou outras linguagens; não há evidência de necessidade
dentro da V1. As cinco specs novas fecham obrigações já existentes.

## Rules applied during review

- Change type: documentation-update com diagnóstico de implementação.
- Workflows: review.md, documentation-update.md; feature.md em planner mode.
- security.md: dados sintéticos, confinamento, contratos e vulnerabilidades.
- documentation-style.md: referências, distinção intenção/evidência e comandos.
- delivery-evidence.md e delivery-surface.md: nenhuma entrega por contagem,
  links de reports não substituem reachability/integração.
- knowledge-governance.md: handoff com owner/TTL, sem atestar trabalho futuro.
- _base-recurrence.md: recorrência por causa/domínio; threshold não atingido.
- Backend-Go e CI de stack: extensões não instaladas; não inventar regras.
  Registrar sua adoção/mapeamento na remediação de governança.
- Frontend/Kubernetes/Cloudflare: não aplicáveis ao produto local inspecionado.

## Checks e evidências

- pose assess discover: um componente; métricas acima.
- pose assess tech-debt: zero marcadores; não prova ausência de defeitos.
- pose assess integrate: zero contratos detectados; resultado inconclusivo
  para MCP Go. Contribuição sanitizada preparada apenas localmente.
- pose check --strict: SUCCESS antes e após os ajustes de DAG.
- pose lint-spec --all --ready-check: 34 specs, zero falhas após alterações.
- pose lint-spec --all --strict: 34 specs no diff final, zero falhas;
  permanecem três avisos de âncoras históricas de follow-ups covered.
- pose validate --strict: SUCCESS, 11 checks; [resultado estruturado](../results/planning-audit-validation.json).
  Foi necessário incluir /home/go/go/bin no PATH e executar fora do sandbox
  para o cache Go/scanner; a execução final reportou No vulnerabilities found.
- go test -race ./...: PASS em todos os pacotes com testes.
- catalog validate --json: sem diagnóstico estrutural bloqueante, 16 avisos
  relation_isolated. --v1-gate falha nas cinco dimensões quantitativas.
- pose skills-check --strict: 12 skills, zero erros/avisos.
- pose knowledge-check --strict: oito artefatos no diff final, zero erros
  e zero vencidos; followups --open registra seis abertos, zero sem owner.
- pose index e pose state refresh: executados; resumo curado preservado
  com o diagnóstico atualizado. git diff --check: sem erros.
- pose docs-check: não executa auditoria porque não existe docs manifest.
- pose recurrence-check --tolerant --window-days 14: zero chaves sinalizadas.
- Probe de restart: reproduziu os erros descritos em F1, sem usar estado do aluno.

### Plano e bundles

Scope: roadmap:codinho-v1. Digest inicial do plano:
sha256:5aee1d13feadac5e04ec1513e7fca1850315f201debc360b15175feeb840209d.
Após incluir as remediações e regenerar os índices:
sha256:c3134d8d1920da6966fac418d3b51b9c35cfedf6f2867fe872d4c0e97d744232.
O digest identifica critérios/ferramentas, não aprovação da entrega.

Required ativos: assess-integrate, skills-check, validate e roadmap-check
executados. O plano pede validação por cmd/internal/packs/schemas/scripts/
testdata; executada matriz completa do módulo Go raiz que os contém, sem
afirmar que pastas sem matriz tiveram check de conteúdo próprio.
Recommended: suggest-review em todos esses paths; discover global e
tech-debt executados. Descobertas por diretório adicionais não repetidas:
o índice atual resolve um único componente; seu mapeamento é um achado.
Completion: review-check/closeout-check/attestation terminal permanecem
sem aprovação enquanto requisitos e milestones estão abertos.

review bundle --explain não produz candidato selável: paths históricos
ausentes e child bundles inexistentes. review verify retorna needs-validation;
review-check não aprova o roadmap. Não se executou auto-attest nem se criou
aprovação fictícia. A revisão de planejamento está concluída; a revisão de
entrega permanece pendente das specs e dos gates enumerados.

## Fixes applied

- Cinco specs draft novas com escopo, dependências, requisitos, arquivos,
  integração e testes negativos.
- Roadmap com três milestones adicionais e oito critérios declarativos.
- Specs adaptativas e aceite com dependências/casos faltantes.
- Quatro specs de packs com paths/CLI corrigidos; aritmética reconciliada.
- Seis follow-ups com owner e SLA; resumo curado do estado corrigido.
- [Handoff](../knowledge/2026-09-07-handoff-planning-audit-2026-09.md) para continuidade.
- Assessments e índices/estado derivados atualizados por comandos nativos.
- Nenhuma alteração no código do produto ou em packs/go-errors.yaml.

## Residual risks e próximos responsáveis

Owner de triagem: @oseiaspereira, revisão até 2026-09-21.
Os achados high são bloqueadores da V1, não riscos aceitos.
Humanos ainda precisam executar playtests, revisar threat model e validar
experiência nos hosts. Nenhum desses gates foi substituído por agente.
As opções de formato/ADR dos novos contratos serão resolvidas antes do código.
Renames históricos e políticas disabled podem exigir correção do engine;
não fabricar evidência de proveniência.

Limite da auditoria: inspeção local e testes acima. Não é análise exaustiva
de segurança de cada função, avaliação pedagógica humana de cada desafio
nem consulta do estado atual de GitHub/hosts externos.
