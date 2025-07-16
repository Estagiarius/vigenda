# Tutorial 03: Dominando o Banco de Questões e a Geração de Provas

Este tutorial explora como construir um banco de questões robusto no Vigenda e como utilizar a **Interface de Texto do Usuário (TUI)** para gerar provas personalizadas. Cobriremos a criação de um arquivo JSON de questões, a importação via linha de comando (CLI) e a geração de provas com diversos critérios através da TUI.

**Pré-requisitos:**
*   Vigenda instalado e funcionando. Consulte [**INSTALLATION.MD**](../../INSTALLATION.MD).
*   Uma ou mais **disciplinas** já devem ter sido criadas na TUI (ex: "Geografia", "História"). Se não, siga o guia em [**Exemplos de Uso da TUI**](../../docs/user_manual/TUI_EXAMPLES.md).
*   Compreensão básica do formato JSON.

## Passo 1: Criar um Arquivo JSON de Questões Detalhado

A maneira mais eficiente de adicionar muitas questões ao Vigenda é através de um arquivo JSON. Crie um arquivo chamado `banco_questoes_completo.json`.

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
    "enunciado": "Descreva as principais formas de relevo encontradas no Brasil.",
    "resposta_correta": "As principais formas são planaltos, planícies e depressões."
  },
  {
    "disciplina": "História",
    "topico": "Idade Média",
    "tipo": "multipla_escolha",
    "dificuldade": "facil",
    "enunciado": "Qual era o sistema predominante na Europa durante a Idade Média?",
    "opcoes": ["Capitalismo", "Socialismo", "Feudalismo", "Mercantilismo"],
    "resposta_correta": "Feudalismo"
  }
]
```
Salve este arquivo.

## Passo 2: Importar as Questões para o Banco (via CLI)

Para este passo, você precisará usar o terminal com a aplicação Vigenda já compilada.

1.  **Feche a TUI** se ela estiver aberta.
2.  Execute o seguinte comando no seu terminal, no mesmo diretório do executável `vigenda`:
    ```bash
    ./vigenda bancoq add banco_questoes_completo.json
    ```
3.  O sistema deve confirmar a importação. Ex: `3 questões importadas com sucesso.`

**Nota:** A TUI é ideal para gerar provas, mas a importação em lote via CLI é a forma recomendada para popular seu banco de dados rapidamente.

## Passo 3: Gerar Provas Usando a TUI

Agora, vamos usar a interface interativa para criar as provas.

1.  **Inicie a TUI do Vigenda:**
    ```bash
    ./vigenda
    ```
2.  No "Menu Principal", navegue com as setas até **"Gerar Provas"** e pressione **Enter**.

### Exemplo 1: Prova de Geografia sobre Relevo e Hidrografia

1.  A TUI pedirá para você **selecionar a disciplina**. Escolha "Geografia" e pressione **Enter**.
2.  Em seguida, você verá a tela de **"Critérios da Prova"**. Preencha os campos:
    *   **Tópico (opcional):** Digite `Relevo e Hidrografia`.
    *   **Número de questões fáceis:** Digite `1`.
    *   **Número de questões médias:** Digite `1`.
    *   **Número de questões difíceis:** Deixe `0`.
    *   **Nome do arquivo de saída:** Digite `prova_geo_relevo.txt`.
3.  Navegue até o botão **"Gerar Prova"** e pressione **Enter**.
4.  A TUI confirmará que a prova foi gerada no arquivo `prova_geo_relevo.txt`.

Abra o arquivo `prova_geo_relevo.txt` para ver o resultado. Ele deve conter as duas questões de Geografia sobre "Relevo e Hidrografia" do nosso JSON.

### Exemplo 2: Prova de História sobre a Idade Média

1.  Se você ainda estiver na tela de geração, pressione **Esc** para voltar ao menu principal e entre em **"Gerar Provas"** novamente.
2.  **Selecione a disciplina** "História".
3.  Na tela de **"Critérios da Prova"**:
    *   **Tópico (opcional):** `Idade Média`
    *   **Número de questões fáceis:** `1`
    *   **Número de questões médias:** `0`
    *   **Número de questões difíceis:** `0`
    *   **Nome do arquivo de saída:** `prova_hist_idademedia.txt`
4.  Selecione **"Gerar Prova"** e pressione **Enter**.

Verifique o arquivo `prova_hist_idademedia.txt`.

### Exemplo 3: Tentativa de Gerar Prova com Mais Questões do que o Disponível

Vamos tentar gerar uma prova com 5 questões difíceis de História.

1.  Volte para a tela de critérios de prova para a disciplina "História".
2.  Preencha os campos:
    *   **Número de questões difíceis:** `5`
    *   **Nome do arquivo de saída:** `prova_hist_teste.txt`
3.  Selecione **"Gerar Prova"**.

A TUI provavelmente exibirá um **aviso ou erro** na tela, informando que não há questões suficientes para atender ao pedido. O arquivo `prova_hist_teste.txt` pode não ser criado ou pode estar vazio.

## Passo 4: Analisando o Arquivo de Saída da Prova

Abra qualquer um dos arquivos de prova gerados. O formato é um texto simples, pronto para impressão.

**Estrutura esperada:**
*   Um cabeçalho (Nome da Disciplina, etc.).
*   As questões numeradas.
*   Para questões de múltipla escolha, as opções listadas (a, b, c, d).

**Importante:** As **respostas corretas não são incluídas** no arquivo da prova, garantindo que ele possa ser entregue diretamente aos alunos.

## Dicas para um Banco de Questões Eficaz

*   **Consistência nos Tópicos:** Use nomes de tópicos consistentes em seus arquivos JSON para facilitar a filtragem na TUI.
*   **Variedade de Dificuldades:** Crie um bom equilíbrio de questões fáceis, médias e difíceis.
*   **Backup dos Arquivos JSON:** Mantenha seus arquivos `.json` seguros. Eles são a fonte do seu banco de questões.

## Conclusão

Este tutorial mostrou o fluxo de trabalho combinado para o banco de questões: **importação em massa via CLI** e **geração de provas personalizadas via TUI**. Essa abordagem aproveita o melhor de ambos os mundos: a eficiência da linha de comando para entrada de dados e a facilidade de uso da interface interativa para a criação de avaliações.
