# Weather Service - Integração Open Telemetry e Zapkin
- Este projeto é um serviço web desenvolvido em Go que recebe um CEP válido, identifica a cidade correspondente e retorna o clima atual com as temperaturas formatadas em **Celsius**, **Fahrenhei**t e **Kelvin**.

- Utilizando **Open Telemetry** e **Zapkin** para observabilidade da aplicação!

## Requisitos do Sistema

Receber um CEP válido de 8 dígitos.
Consultar a localização usando a API ViaCEP.

Consultar as temperaturas da localização usando a API WeatherAPI.

Retornar as temperaturas nos seguintes formatos:

- Celsius
- Fahrenheit
- Kelvin

Registrar o track da aplicação no Zapkin.

# **Executar com Docker Compose**
Você pode executar o sistema dentro de um Docker Compose. Para isso, siga os passos abaixo.

### **Configure suas credenciais de API:**
- WeatherAPI: Obtenha uma chave em [WeatherAPI](https://www.weatherapi.com/).
- Copie o arquivo **.env.exemplo** com o nome **.env**
```
WEATHER_API_KEY="SUA_CHAVE_AQUI"
```

## **Execute o Docker Compose**
```bash
docker compose up --build
```

### Exemplo de requisição

```bash
curl --location 'http://localhost:8080/weather' \
--header 'Content-Type: application/json' \
--data '{
    "cep": "81900120"
}'
```

### **Acesse o Zapkin**
O zapkin estará rodando em http://localhost:9411

# **Licença**
Este projeto está licenciado sob a [Licença MIT](LICENSE).
