Programa GO do curso de PÓS Graduação GOEXPERTS que utiliza Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.
  As duas requisições serão feitas simultaneamente para as seguintes APIs:
  https://brasilapi.com.br/api/cep/v1/ + cep
  http://viacep.com.br/ws/" + cep + "/json/
Os requisitos para este desafio são:
  - Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
  - O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.
  - Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.
