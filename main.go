package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Documentacao struct {
	Titulo      string   `json:"titulo"`
	Docente     string   `Json:"Docente"`
	Disciplina  string   `json:"disciplina"`
	Integrantes []string `json:"integrantes"`
	Secao1      Secao1   `json:"secao_1_pontos_fortes_e_filosofia"`
	Secao2      Secao2   `json:"secao_2_pre_requisitos_e_setup"`
	Secao3      Secao3   `json:"secao_3_codigo_do_servidor"`
	Secao4      Secao4   `json:"secao_4_conclusao"`
}

type Secao1 struct {
	Descricao      string   `json:"descricao"`
	Vantagens      []string `json:"vantagens"`
	CasoDeUsoIdeal string   `json:"caso_de_uso_ideal"`
}

type Secao2 struct {
	PreRequisitos []string `json:"pre_requisitos"`
	Comandos      []string `json:"comandos"`
}

type Secao3 struct {
	Arquivo        string `json:"arquivo"`
	Comando        string `json:"comando_para_rodar"`
	Explicacao     string `json:"explicacao"`
	CodigoCompleto string `json:"codigo_completo"`
}

type Secao4 struct {
	Dificuldade  string `json:"principal_dificuldade"`
	Recomendacao string `json:"recomendacao_final"`
}

func paginaInicial(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "O metodo foi negado, parceiro!", http.StatusMethodNotAllowed)
		return
	}
	templateHTML, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Erro ao carregar o HTML:", err)
		http.Error(w, "Erro interno ao carregar a página", http.StatusInternalServerError)
		return
	}

	documentacao, err := criarDocumentacao()
	if err != nil {
		log.Println("Erro ao criar a documentação:", err)
		http.Error(w, "Erro interno ao carregar a documentação", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := templateHTML.ExecuteTemplate(w, "index.html", documentacao); err != nil {
		log.Println("Erro ao renderizar o HTML:", err)
	}
}

func criarDocumentacao() (Documentacao, error) {
	codigoFonte, err := os.ReadFile("main.go")
	if err != nil {
		return Documentacao{}, err
	}

	documentacao := Documentacao{
		Titulo:     "Go (Golang) - Servidor Web",
		Docente:    "Welington Gadelha",
		Disciplina: "Programacao Web - UNIVALI",
		Integrantes: []string{
			"Arthur Sperandio Marcelino",
			"Anderson Marques",
		},

		Secao1: Secao1{
			Descricao: "Go é uma linguagem de programacao criada pelo Google. Ela foi desenvolvida com foco em simplicidade, desempenho e facilidade para trabalhar com sistemas e servidores web.",
			Vantagens: []string{
				"Sintaxe simples e organizada",
				"Bom desempenho",
				"Compilação rapida",
				"Pacote net/http disponivel na biblioteca padrão",
			},

			CasoDeUsoIdeal: "Go é muito indicado para APIs, servidores web e sistemas que precisam lidar com varias requisições ao mesmo tempo.",
		},

		Secao2: Secao2{
			PreRequisitos: []string{
				"Go instalado no computador",
				"Visual Studio Code ou algum outro editor de codigo",
				"Postman(opcional) ou navegador para testar as requisições",
			},
			Comandos: []string{
				"go version",
				"go mod init servidor-web-go",
				"go run main.go",
			},
		},

		Secao3: Secao3{
			Arquivo:    "main.go",
			Comando:    "go run main.go",
			Explicacao: "O servidor utiliza o pacote net/http. A função http.HandleFunc associa a rota principal / a função paginaInicial. Essa função recebe a requisicao HTTP e retorna uma resposta em JSON e o comando http.ListenAndServe inicia o servidor na porta 8083.",

			CodigoCompleto: string(codigoFonte),
		},

		Secao4: Secao4{
			Dificuldade:  "As principais dificuldades foram entender como criar uma rota HTTP, a lógica da linguagem em si, como, por exemplo, oque cada componente faz e os seus relacionamentos. Além de como montar isso de forma que funcione onde o servidor receba uma requisição e devolva uma resposta.",
			Recomendacao: "Go é uma boa opção para desenvolvimento back-end pois possui uma sintaxe compacta e de fácil compreensão (de acordo com a documentação), bom desempenho comparado a outras linguagens e ferramentas disponíveis para trabalhar com servidores HTTP.",
		},
	}

	return documentacao, nil
}

func apiDocumentacao(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/docs" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "O método foi negado", http.StatusMethodNotAllowed)
		return
	}

	documentacao, err := criarDocumentacao()
	if err != nil {
		log.Println("Erro ao criar a documentação:", err)
		http.Error(w, "Erro interno ao carregar a documentação", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(documentacao); err != nil {
		log.Println("Erro ao enviar a documentação em JSON:", err)
	}
}

func main() {
	http.HandleFunc("/api/docs", apiDocumentacao)
	http.HandleFunc("/", paginaInicial)
	log.Println("Servidor iniciado em http://localhost:8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
