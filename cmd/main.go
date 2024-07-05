/*
=====================================================================================================
  - main.go : Deverá usar o que aprendemos com Multithreading e APIs para buscar o resultado mais
  - rápido entre duas APIs distintas.
  - As duas requisições serão feitas simultaneamente para as seguintes APIs:
  - https://brasilapi.com.br/api/cep/v1/ + cep
  - http://viacep.com.br/ws/" + cep + "/json/
  - Os requisitos para este desafio são:
  - Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
  - O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual
  - API a enviou.
  - Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.

=====================================================================================================
*/

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Definição de constantes para chamada das APIs de busca de CEP
const (
	BrasilApi_URL string = "https://brasilapi.com.br/api/cep/v1/%s"
	ViaCep_URL    string = "https://viacep.com.br/ws/%s/json/"
)

// Estruct que receberá o resultado da chamada das APIs
type ResultCEP struct {
	URL_Vencedora string //Nome da URL que respondeu mais rápido
	Dados         string //Dados retornados pela API
}

func main() {

	//Verifica se foi informado o cep na command line
	if len(os.Args) < 2 {
		fmt.Println("ATENÇÃO\nCEP inválido!\nUse: go run main.go <cep>")
		os.Exit(1)
	}
	cep := os.Args[1]

	//Cria os canais para execução das treads
	canalBrasilApi := make(chan ResultCEP)
	canalViaCep := make(chan ResultCEP)

	//Chama CepWorker com a URL da BrasilApi em uma Tread separada
	go CepWorker(BrasilApi_URL, cep, canalBrasilApi)

	//Chama CepWorker com a URL da ViaApi em uma Tread separada
	go CepWorker(ViaCep_URL, cep, canalViaCep)

	//Select para esperar a execução das treads ou timeout
	select {
	case brasilApi := <-canalBrasilApi:
		fmt.Printf("URL: %s\n\nRESPOSTA: %s\n", brasilApi.URL_Vencedora, brasilApi.Dados)
	case viaCep := <-canalViaCep:
		fmt.Printf("URL: %s\n\nRESPOSTA: %s\n", viaCep.URL_Vencedora, viaCep.Dados)
	case <-time.After(time.Second):
		log.Fatalln("Tempo de resposta excedido")
	}

}

/*
==========================================================
  - Função: CepWorker
  - Descrição : Função que executa a requisição para as
  - APIs de busca de CEP.
  - Parametros :
  - url - URL da API a ser executada - tipo String
  - cep - CEP a ser passado para a API - tipo: String
  - canalCEP - Canal no qual será devolvido o resultado da
  - API chamada
  - Retorno:

==========================================================
*/
func CepWorker(url string, cep string, canalCEP chan<- ResultCEP) {
	//Cria uma estrutura com a URL e o CEP recebido
	resultCEP := ResultCEP{URL_Vencedora: fmt.Sprintf(url, cep)}

	//Cria uma nova requisição GET com a URL
	requestCEP, err := http.NewRequest("GET", resultCEP.URL_Vencedora, nil)
	if err != nil {
		close(canalCEP) //Fecha o canal
		return
	}

	//Executa a requisição criada
	resultado, err := http.DefaultClient.Do(requestCEP)
	if err != nil {
		close(canalCEP) //Fecha o canal
		return
	}

	//Lê o resultado da requisição
	body, err := io.ReadAll(resultado.Body)
	_ = resultado.Body.Close()
	if err != nil {
		close(canalCEP) //Fecha o canal
		return
	}

	//Joga o body para Dados
	resultCEP.Dados = string(body)

	//Joga a estrutura no canal
	canalCEP <- resultCEP
}
