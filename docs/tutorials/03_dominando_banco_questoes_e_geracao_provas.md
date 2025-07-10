# Tutorial 03: Dominando o Banco de Questões e a Geração de Provas

Este tutorial explora em profundidade como construir um banco de questões robusto no Vigenda e como utilizar seus recursos para gerar provas personalizadas. Cobriremos a criação de um arquivo JSON de questões mais elaborado, a importação para o banco e a geração de provas com diversos critérios.

**Pré-requisitos:**
*   Vigenda instalado e funcionando. Consulte [**INSTALLATION.MD**](../../INSTALLATION.MD).
*   Uma ou mais **disciplinas** já devem ter sido criadas através da **Interface de Texto do Usuário (TUI)** do Vigenda.
    *   Para iniciar a TUI: execute `vigenda` (ou `go run ./cmd/vigenda/main.go`).
    *   No menu da TUI, crie as disciplinas necessárias (ex: "Geografia", "História"). Anote os IDs das disciplinas criadas (vamos assumir "Geografia" com `ID = 2` e "História" com `ID = 1` para este tutorial).
*   Compreensão básica do formato JSON para importação de questões e dos comandos CLI `vigenda bancoq add` e `vigenda prova gerar`. Consulte o [Manual do Usuário](../../docs/user_manual/README.md) para detalhes.

## Passo 1: Criar um Arquivo JSON de Questões Detalhado

Vamos criar um arquivo chamado `banco_questoes_completo.json`. Este arquivo conterá questões de diferentes disciplinas, tópicos, tipos e dificuldades.

```json
[
  {
    "disciplina": "Geografia",
    "topico": "Relevo e Hidrografia",
    "tipo": "multipla_escolha",
    "dificuldade": "facil",
    "enunciado": "Qual é o maior rio do mundo em volume de água?",
    "opcoes": ["Nilo", "Amazonas", "Mississipi", "Yangtzé"],
    "resposta_correta": "Amazonas"
  },
  {
    "disciplina": "Geografia",
    "topico": "Relevo e Hidrografia",
    "tipo": "dissertativa",
    "dificuldade": "media",
    "enunciado": "Descreva as principais formas de relevo encontradas no Brasil e dê um exemplo de cada.",
    "resposta_correta": "As principais formas são planaltos (ex: Planalto Central), planícies (ex: Planície Amazônica) e depressões (ex: Depressão Sertaneja). Montanhas de formação recente não são expressivas no Brasil."
  },
  {
    "disciplina": "Geografia",
    "topico": "Clima",
    "tipo": "multipla_escolha",
    "dificuldade": "media",
    "enunciado": "Qual fator climático é o principal responsável pelas estações do ano?",
    "opcoes": ["Rotação da Terra", "Translação da Terra e inclinação do eixo", "Correntes marítimas", "Altitude"],
    "resposta_correta": "Translação da Terra e inclinação do eixo"
  },
  {
    "disciplina": "Geografia",
    "topico": "Clima",
    "tipo": "dissertativa",
    "dificuldade": "dificil",
    "enunciado": "Explique o fenômeno El Niño e suas principais consequências para o clima global.",
    "resposta_correta": "O El Niño é um aquecimento anômalo das águas superficiais do Oceano Pacífico Equatorial, alterando os padrões de vento e chuva em diversas regiões do mundo, causando secas em alguns lugares e inundações em outros."
  },
  {
    "disciplina": "História",
    "topico": "Idade Média",
    "tipo": "multipla_escolha",
    "dificuldade": "facil",
    "enunciado": "Qual era o sistema econômico, político e social predominante na Europa Ocidental durante a Idade Média?",
    "opcoes": ["Capitalismo", "Socialismo", "Feudalismo", "Mercantilismo"],
    "resposta_correta": "Feudalismo"
  },
  {
    "disciplina": "História",
    "topico": "Idade Média",
    "tipo": "dissertativa",
    "dificuldade": "media",
    "enunciado": "Discorra sobre o papel da Igreja Católica na sociedade feudal.",
    "resposta_correta": "A Igreja Católica detinha grande poder espiritual, cultural, econômico e político, influenciando todos os aspectos da vida medieval, desde a educação até a organização social e as relações de poder."
  },
  {
    "disciplina": "História",
    "topico": "Grandes Navegações",
    "tipo": "multipla_escolha",
    "dificuldade": "media",
    "enunciado": "Qual navegador português foi o primeiro a completar o contorno da África, chegando às Índias?",
    "opcoes": ["Cristóvão Colombo", "Fernão de Magalhães", "Vasco da Gama", "Pedro Álvares Cabral"],
    "resposta_correta": "Vasco da Gama"
  }
]
```
Salve este arquivo.

## Passo 2: Importar as Questões para o Banco

Agora, adicione essas questões ao banco de dados do Vigenda:

```bash
./vigenda bancoq add banco_questoes_completo.json
```

O sistema deve confirmar a importação. Ex: `7 questões importadas com sucesso de banco_questoes_completo.json.`

**Verificação (Opcional):**
O `AGENTS.md` não especifica um comando para listar questões do banco. Se houvesse (ex: `vigenda bancoq listar --disciplina Geografia`), seria útil para verificar. Sem isso, a confirmação da geração de provas (próximo passo) será nosso principal indicador.

## Passo 3: Gerar Provas com Diferentes Critérios

Vamos gerar algumas provas usando as questões que importamos. Assumimos que a disciplina "Geografia" tem `ID = 2` e "História" tem `ID = 1`.

1.  **Prova de Geografia com foco em Relevo e Hidrografia (2 questões):**
    ```bash
    ./vigenda prova gerar --subjectid 2 --topic "Relevo e Hidrografia" --total 2 --output prova_geo_relevo.txt
    ```
    *   `--subjectid 2`: Seleciona a disciplina Geografia.
    *   `--topic "Relevo e Hidrografia"`: Filtra pelo tópico.
    *   `--total 2`: Pede 2 questões.
    *   `--output prova_geo_relevo.txt`: Salva a prova em um arquivo.

    Abra `prova_geo_relevo.txt` para ver o resultado. Deve conter as duas questões de Geografia sobre "Relevo e Hidrografia" do nosso JSON.

2.  **Prova de História com 1 questão fácil e 1 média sobre Idade Média:**
    ```bash
    ./vigenda prova gerar --subjectid 1 --topic "Idade Média" --easy 1 --medium 1 --output prova_hist_idademedia.txt
    ```
    *   `--subjectid 1`: Seleciona a disciplina História.
    *   `--topic "Idade Média"`: Filtra pelo tópico.
    *   `--easy 1 --medium 1`: Especifica a quantidade por dificuldade.

    Verifique o arquivo `prova_hist_idademedia.txt`.

3.  **Prova de Geografia geral com 3 questões, balanceando dificuldades:**
    ```bash
    ./vigenda prova gerar --subjectid 2 --total 3 --output prova_geo_geral.txt
    ```
    O Vigenda tentará selecionar uma mistura de dificuldades entre as questões de Geografia disponíveis para compor as 3 questões.

4.  **Tentativa de gerar uma prova com mais questões do que o disponível:**
    Vamos pedir 3 questões difíceis de História. No nosso JSON, não temos nenhuma questão difícil de História.
    ```bash
    ./vigenda prova gerar --subjectid 1 --hard 3 --output prova_hist_dificil_teste.txt
    ```
    Observe a saída no console. O Vigenda deve informar que não há questões suficientes para atender ao pedido ou gerar uma prova com menos questões, acompanhada de um aviso. O conteúdo de `prova_hist_dificil_teste.txt` refletirá isso.

    Exemplo de aviso no console (hipotético):
    ```
    Aviso: Não foi possível encontrar 3 questões do tipo 'hard' para a disciplina ID 1. Foram encontradas 0.
    Prova gerada com 0 questões.
    ```

## Passo 4: Analisando o Arquivo de Saída da Prova

Abra qualquer um dos arquivos de prova gerados (ex: `prova_geo_relevo.txt`). O formato típico é um texto simples, pronto para impressão ou cópia.

**Estrutura esperada:**
*   Um cabeçalho (ex: Nome da Disciplina, Título da Prova).
*   As questões numeradas.
*   Para questões de múltipla escolha, as opções listadas (ex: a), b), c), d)).
*   Para questões dissertativas, um espaço para resposta.

**Importante:**
*   **Respostas Corretas:** Verifique se as respostas corretas estão incluídas no arquivo da prova. Para uma prova a ser entregue aos alunos, elas **não devem** estar presentes. Pode haver uma opção separada no Vigenda para gerar um gabarito (não especificado no `AGENTS.md` atual).
*   **Layout:** O layout é básico. Para formatação avançada, você pode copiar o texto para um editor de documentos (Word, Google Docs, etc.).

## Dicas para um Banco de Questões Eficaz

*   **Consistência nos Tópicos:** Use nomes de tópicos consistentes para facilitar a filtragem.
*   **Variedade de Dificuldades:** Tente ter um bom equilíbrio de questões fáceis, médias e difíceis para cada tópico/disciplina.
*   **Revisão Regular:** Periodicamente, revise as questões no seu banco para garantir que ainda são relevantes e precisas. (Isso exigiria uma forma de listar ou exportar questões do banco).
*   **Backup dos Arquivos JSON:** Mantenha backups dos seus arquivos JSON de questões. Eles são a fonte do seu banco.

## Conclusão

Este tutorial mostrou como criar um banco de questões mais completo e como usar as opções do comando `vigenda prova gerar` para criar provas altamente personalizadas. Você aprendeu a filtrar por disciplina e tópico, especificar o número de questões por dificuldade ou um total, e como o sistema lida com a insuficiência de questões.

Com um banco de questões bem curado, gerar avaliações se torna uma tarefa muito mais rápida e eficiente!
