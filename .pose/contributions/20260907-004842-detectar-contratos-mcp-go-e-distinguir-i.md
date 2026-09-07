---
title: "Detectar contratos MCP Go e distinguir inventario de integracao verificada"
type: "enhancement"
module: ""
created_at: "2026-09-07T00:48:42Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: Detectar contratos MCP Go e distinguir inventario de integracao verificada

## Problem / Observation

Em uma aplicação Go com tools registradas pelo SDK MCP oficial, assess
integrate pode retornar zero contratos avaliados/ativos e zero gaps. Esse
resultado não distingue ausência real de contratos de ausência de cobertura
do detector. O relato usa apenas um exemplo sintético; não contém código ou
dados da aplicação em que a limitação foi observada.

## Reproduction (Synthetic)

1. Crie um módulo Go sintético que dependa do SDK MCP oficial.
2. Registre uma tool fictícia hello usando mcp.AddTool em uma função Go.
3. Adicione um teste de contrato que enumere e chame hello por stdio.
4. Execute pose assess integrate e compare o inventário estático à tool
   conhecida. Na observação original, o detector retornou zero contratos.
5. Use o exemplo para reproduzir e acrescentar um fixture de regressão no
   engine; a reprodução sintética upstream ainda não foi executada aqui.

## Proposed Solution

Reconheça registros mcp.AddTool via AST/import resolvido. Declare provedores,
consumidores e cobertura do detector separadamente; zero itens reconhecidos
com SDK MCP presente deve produzir unknown/partial, nunca sugerir prova de
integração. Permita inventário declarativo complementar e exija evidência de
teste para afirmar integração. Não execute código da aplicação para descobri-lo.

## Submission

Somente staging local. Nenhuma issue ou mensagem upstream foi enviada.
