# Tutorial 03: Dominando o Banco de Questões e a Geração de Provas

Este tutorial explora em profundidade como construir um banco de questões robusto no Vigenda e como utilizar seus recursos para gerar provas personalizadas. Cobriremos a criação de um arquivo JSON de questões mais elaborado, a importação para o banco e a geração de provas com diversos critérios.

**Pré-requisitos:**
*   Vigenda instalado e configurado (consulte o [Guia de Introdução](../../docs/getting_started/README.md)).
*   Familiaridade com o formato JSON para importação de questões, conforme detalhado na seção [Importação de Questões (JSON)](../../user_manual/README.md#importacao-de-questoes-json) do Manual do Usuário.
*   Conhecimento dos seguintes comandos do Vigenda:
    *   `vigenda bancoq add`
    *   `vigenda bancoq listar` (Para verificar as questões importadas e obter seus IDs, se necessário para edição/remoção)
    *   `vigenda prova gerar`
    *   `vigenda prova gabarito` (Para gerar o gabarito correspondente à prova)
*   IDs de disciplinas existentes no sistema. Para este tutorial, vamos assumir que:
    *   Disciplina "Geografia" tem `ID = 2`.
    *   Disciplina "História" tem `ID = 1`.
    (Se você não tem disciplinas, o Vigenda pode criá-las implicitamente durante a importação de questões se o JSON contiver nomes de disciplinas, ou você pode precisar de um mecanismo para gerenciá-las, dependendo da versão do Vigenda).

## Passo 1: Criar um Arquivo JSON de Questões Detalhado

Vamos elaborar um arquivo JSON chamado `banco_questoes_avancado.json`. Este arquivo incluirá questões com mais detalhes, como `id_externo` (para seu controle), `tags` e `tempo_estimado_min`.

```json
[
  {
    "id_externo": "GEO-RH-001",
    "disciplina": "Geografia",
    "subject_id": 2,
    "topico": "Relevo e Hidrografia",
    "tipo": "multipla_escolha",
    "dificuldade": "facil",
    "enunciado": "Qual é o maior rio do mundo em volume de água?",
    "opcoes": ["Nilo", "Amazonas", "Mississipi", "Yangtzé"],
    "resposta_correta": "Amazonas",
    "tags": ["hidrografia", "mundo", "rios"],
    "tempo_estimado_min": 1
  },
  {
    "id_externo": "GEO-RH-002",
    "disciplina": "Geografia",
    "subject_id": 2,
    "topico": "Relevo e Hidrografia",
    "tipo": "dissertativa",
    "dificuldade": "media",
    "enunciado": "Descreva as principais formas de relevo encontradas no Brasil e dê um exemplo de cada.",
    "resposta_correta": "As principais formas são planaltos (ex: Planalto Central), planícies (ex: Planície Amazônica) e depressões (ex: Depressão Sertaneja). Montanhas de formação recente não são expressivas no Brasil.",
    "tags": ["relevo", "brasil", "geomorfologia"],
    "tempo_estimado_min": 3
  },
  {
    "id_externo": "GEO-CL-001",
    "disciplina": "Geografia",
    "subject_id": 2,
    "topico": "Clima",
    "tipo": "multipla_escolha",
    "dificuldade": "media",
    "enunciado": "Qual fator climático é o principal responsável pelas estações do ano?",
    "opcoes": ["Rotação da Terra", "Translação da Terra e inclinação do eixo", "Correntes marítimas", "Altitude"],
    "resposta_correta": "Translação da Terra e inclinação do eixo",
    "tags": ["clima", "estações", "astronomia"],
    "tempo_estimado_min": 2
  },
  {
    "id_externo": "GEO-CL-002",
    "disciplina": "Geografia",
    "subject_id": 2,
    "topico": "Clima",
    "tipo": "verdadeiro_falso",
    "dificuldade": "facil",
    "enunciado": "O efeito estufa é um fenômeno exclusivamente causado pela ação humana.",
    "opcoes": ["Verdadeiro", "Falso"],
    "resposta_correta": "Falso",
    "tags": ["clima", "efeito estufa", "meio ambiente"],
    "tempo_estimado_min": 1
  },
  {
    "id_externo": "GEO-CL-003",
    "disciplina": "Geografia",
    "subject_id": 2,
    "topico": "Clima",
    "tipo": "dissertativa",
    "dificuldade": "dificil",
    "enunciado": "Explique o fenômeno El Niño e suas principais consequências para o clima global.",
    "resposta_correta": "O El Niño é um aquecimento anômalo das águas superficiais do Oceano Pacífico Equatorial, alterando os padrões de vento e chuva em diversas regiões do mundo, causando secas em alguns lugares e inundações em outros.",
    "tags": ["clima", "el niño", "oceanografia", "impactos globais"],
    "tempo_estimado_min": 5
  },
  {
    "id_externo": "HIST-IM-001",
    "disciplina": "História",
    "subject_id": 1,
    "topico": "Idade Média",
    "tipo": "multipla_escolha",
    "dificuldade": "facil",
    "enunciado": "Qual era o sistema econômico, político e social predominante na Europa Ocidental durante a Idade Média?",
    "opcoes": ["Capitalismo", "Socialismo", "Feudalismo", "Mercantilismo"],
    "resposta_correta": "Feudalismo",
    "tags": ["idade media", "europa", "sistemas economicos"],
    "tempo_estimado_min": 1
  },
  {
    "id_externo": "HIST-IM-002",
    "disciplina": "História",
    "subject_id": 1,
    "topico": "Idade Média",
    "tipo": "dissertativa",
    "dificuldade": "media",
    "enunciado": "Discorra sobre o papel da Igreja Católica na sociedade feudal, abordando aspectos políticos, culturais e sociais.",
    "resposta_correta": "A Igreja Católica detinha grande poder espiritual (salvação da alma), cultural (monopólio do conhecimento, universidades), econômico (terras, dízimos) e político (influência sobre reis e nobres, mediação de conflitos). Ela moldava a visão de mundo, a moral e o cotidiano da sociedade feudal.",
    "tags": ["idade media", "igreja catolica", "poder"],
    "tempo_estimado_min": 4
  },
  {
    "id_externo": "HIST-GN-001",
    "disciplina": "História",
    "subject_id": 1,
    "topico": "Grandes Navegações",
    "tipo": "multipla_escolha",
    "dificuldade": "media",
    "enunciado": "Qual navegador português foi o primeiro a completar o contorno da África, chegando às Índias?",
    "opcoes": ["Cristóvão Colombo", "Fernão de Magalhães", "Vasco da Gama", "Pedro Álvares Cabral"],
    "resposta_correta": "Vasco da Gama",
    "tags": ["grandes navegacoes", "portugal", "india"],
    "tempo_estimado_min": 2
  }
]
```
Salve este arquivo como `banco_questoes_avancado.json`.

## Passo 2: Importar as Questões para o Banco

Agora, vamos adicionar essas questões ao banco de dados do Vigenda:

```bash
./vigenda bancoq add banco_questoes_avancado.json
```

O Vigenda deve confirmar a importação, idealmente informando quantas questões foram adicionadas ou atualizadas.
Exemplo de feedback:
```
Importação de 'banco_questoes_avancado.json' concluída.
8 questões processadas.
8 questões novas adicionadas.
0 questões atualizadas.
0 questões com erro.
```

**Verificação:**
Para garantir que as questões foram importadas corretamente, utilize o comando `bancoq listar`. Você pode filtrar para ver as questões de uma disciplina específica:
```bash
./vigenda bancoq listar --subjectid 2 --limit 5
```
Isso deve listar até 5 questões da disciplina Geografia (ID 2). Verifique se os detalhes como tópico, tipo e dificuldade correspondem ao seu arquivo JSON.
```
ID Questão | Disciplina ID | Tópico                | Tipo              | Dificuldade | Enunciado (início)
-----------|---------------|-----------------------|-------------------|-------------|-----------------------
201        | 2             | Relevo e Hidrografia  | multipla_escolha  | facil       | Qual é o maior rio...
202        | 2             | Relevo e Hidrografia  | dissertativa      | media       | Descreva as principais...
203        | 2             | Clima                 | multipla_escolha  | media       | Qual fator climático...
204        | 2             | Clima                 | verdadeiro_falso  | facil       | O efeito estufa é ...
205        | 2             | Clima                 | dissertativa      | dificil     | Explique o fenômeno...
```
(Os IDs de Questão são exemplos e serão gerados pelo sistema).

## Passo 3: Gerar Provas com Diferentes Critérios

Com nosso banco de questões populado, vamos explorar diferentes formas de gerar provas.

1.  **Prova de Geografia com foco em "Relevo e Hidrografia" (2 questões):**
    Queremos uma prova curta sobre um tópico específico.
    ```bash
    ./vigenda prova gerar --subjectid 2 --topic "Relevo e Hidrografia" --total 2 --output prova_geo_relevo_hidro.txt --title "Avaliação: Relevo e Hidrografia"
    ```
    *   `--subjectid 2`: Filtra pela disciplina Geografia.
    *   `--topic "Relevo e Hidrografia"`: Filtra pelo tópico.
    *   `--total 2`: Solicita um total de 2 questões. O Vigenda selecionará aleatoriamente entre as disponíveis que atendam aos critérios.
    *   `--output prova_geo_relevo_hidro.txt`: Salva a prova no arquivo especificado.
    *   `--title "Avaliação: Relevo e Hidrografia"`: Define um título para a prova.

    Abra o arquivo `prova_geo_relevo_hidro.txt` para verificar. Ele deve conter 2 questões do tópico especificado.

2.  **Prova de História sobre "Idade Média" com distribuição de dificuldade:**
    Queremos 1 questão fácil e 1 questão média.
    ```bash
    ./vigenda prova gerar --subjectid 1 --topic "Idade Média" --easy 1 --medium 1 --total 2 --output prova_hist_idademedia_diff.txt --title "Avaliação: Idade Média (Variada)"
    ```
    *   `--easy 1 --medium 1`: Especifica a quantidade por dificuldade.
    *   `--total 2`: Confirma o número total de questões.

    Verifique o arquivo `prova_hist_idademedia_diff.txt`.

3.  **Prova de Geografia geral com 3 questões, incluindo uma do tipo "Verdadeiro/Falso":**
    ```bash
    ./vigenda prova gerar --subjectid 2 --total 3 --type verdadeiro_falso --output prova_geo_geral_vf.txt --title "Teste Rápido: Geografia (V/F)"
    ```
    *   `--type verdadeiro_falso`: Tenta incluir questões desse tipo. Se não houver questões suficientes do tipo `verdadeiro_falso` para atingir `--total 3`, o Vigenda pode complementar com outros tipos ou informar a limitação. (O comportamento exato de como `--type` interage com `--total` quando há insuficiência de um tipo específico deve ser verificado no Manual do Usuário ou com `--help`). Para garantir que *apenas* questões V/F sejam incluídas, você pode precisar ajustar o `--total` para o número de questões V/F que você sabe que existem ou omitir outras especificações de dificuldade/tipo se o filtro for exclusivo.
    *   Uma abordagem mais precisa se você quiser *uma* questão V/F e as outras de qualquer tipo:
        ```bash
        # Primeiro, gere uma prova com uma questão V/F (se houver)
        ./vigenda prova gerar --subjectid 2 --type verdadeiro_falso --total 1 --output parte1_vf.txt
        # Depois, gere mais duas questões de qualquer tipo, excluindo as já usadas (funcionalidade avançada)
        # ./vigenda prova gerar --subjectid 2 --total 2 --exclude-ids <IDs_da_parte1> --output parte2_outras.txt
        # E então junte os arquivos.
        # Alternativamente, se o sistema permitir, uma combinação mais complexa de filtros.
        # Por simplicidade, o exemplo inicial assume que o sistema tentará incluir o tipo especificado.
        ```
    O Vigenda tentará selecionar uma mistura de questões, priorizando o tipo especificado se possível.

4.  **Tentativa de gerar uma prova com mais questões difíceis de História do que o disponível:**
    No nosso JSON `banco_questoes_avancado.json`, não temos questões de História com dificuldade "dificil".
    ```bash
    ./vigenda prova gerar --subjectid 1 --hard 3 --total 3 --output prova_hist_hard_fail.txt --title "Teste Desafio: História (Difícil)"
    ```
    Observe a saída no console. O Vigenda deve informar que não há questões suficientes:
    ```
    Aviso: Não foi possível encontrar 3 questões do tipo 'hard' para a disciplina ID 1 e tópico(s) especificado(s). Foram encontradas 0.
    Prova "Teste Desafio: História (Difícil)" gerada com 0 questões e salva em prova_hist_hard_fail.txt.
    ```
    O arquivo `prova_hist_hard_fail.txt` provavelmente estará vazio ou conterá apenas o cabeçalho.

## Passo 4: Analisando o Arquivo de Saída da Prova e Gerando Gabarito

1.  **Analisando a Prova Gerada:**
    Abra um dos arquivos de prova, por exemplo, `prova_geo_relevo_hidro.txt`.
    A estrutura esperada é:
    *   Cabeçalho: Título da Prova, Disciplina, Data (se o Vigenda adicionar).
    *   Questões numeradas, com enunciados.
    *   Opções listadas para questões de múltipla escolha ou V/F.
    *   Espaço para resposta para questões dissertativas.
    *   **Importante:** O arquivo da prova gerado por `vigenda prova gerar` **não deve conter as respostas corretas**.

2.  **Gerando o Gabarito:**
    Para cada prova gerada, você deve gerar o gabarito correspondente usando `vigenda prova gabarito`. É crucial usar **exatamente os mesmos critérios de filtragem de questões** usados para `prova gerar` para garantir que o gabarito corresponda à prova.

    Para a prova `prova_geo_relevo_hidro.txt`:
    ```bash
    ./vigenda prova gabarito --subjectid 2 --topic "Relevo e Hidrografia" --total 2 --output gabarito_geo_relevo_hidro.txt --title "GABARITO - Avaliação: Relevo e Hidrografia"
    ```
    Abra `gabarito_geo_relevo_hidro.txt`. Ele deve listar as questões (talvez apenas o número ou início do enunciado) e suas respectivas respostas corretas.

    Para a prova `prova_hist_idademedia_diff.txt`:
    ```bash
    ./vigenda prova gabarito --subjectid 1 --topic "Idade Média" --easy 1 --medium 1 --total 2 --output gabarito_hist_idademedia_diff.txt --title "GABARITO - Avaliação: Idade Média (Variada)"
    ```

**Layout e Formatação:**
Os arquivos `.txt` gerados são de texto simples. Para uma formatação mais elaborada (negrito, itálico, tabelas, etc.), você pode copiar o conteúdo do arquivo de prova e do gabarito para um editor de texto avançado (como Microsoft Word, Google Docs, LibreOffice Writer) ou usar ferramentas que convertem Markdown/texto para PDF com estilos.

## Dicas para um Banco de Questões Eficaz

*   **Use `id_externo`:** Facilita a referência e atualização de questões específicas se você mantiver seus próprios identificadores.
*   **Tags Detalhadas (`tags`):** Quanto mais descritivas e consistentes forem suas tags, mais fácil será encontrar e filtrar questões muito específicas no futuro (assumindo que `bancoq listar` e `prova gerar` possam filtrar por tags).
*   **Tempo Estimado (`tempo_estimado_min`):** Registrar isso pode ajudar a montar provas com uma duração total previsível.
*   **Revisão e Atualização:**
    *   Use `vigenda bancoq listar` com vários filtros para revisar suas questões periodicamente.
    *   Se o comando `vigenda bancoq editar <ID_QUESTAO_INTERNA>` estiver disponível, use-o para correções. Caso contrário, pode ser necessário remover a questão (com `vigenda bancoq remover`) e reimportar a versão corrigida do seu JSON.
*   **Backup dos Arquivos JSON Fonte:** Seus arquivos JSON são a "verdade" do seu banco de questões. Mantenha-os seguros e versionados (usando Git, por exemplo).

## Conclusão

Dominar o banco de questões e a geração de provas no Vigenda envolve criar arquivos JSON bem estruturados e utilizar as diversas opções de filtragem dos comandos `vigenda bancoq add`, `vigenda bancoq listar`, `vigenda prova gerar` e `vigenda prova gabarito`. Com um banco de questões rico e bem organizado, você pode gerar avaliações diversificadas e adequadas às suas necessidades pedagógicas de forma muito mais ágil.

Lembre-se de consultar o `--help` de cada comando para entender todas as suas capacidades e opções mais recentes!
