# go-gerador-pdf
API que recebe JSON e gera um PDF

## Como consultar:
String de consexão
```js
    const url = 'http://localhost:5555/rota_do_relatorio';

    const driver = "sqlserver";
    const usuario = "root";
    const senha = "12345";
    const servidor = "nome_do_servidor";
    const banco = "nome_do_banco";
    const options = "&TrustServerCertificate=true";

    const connectionString = `${driver}://${usuario}:${senha}@${servidor}?database=${banco}${options}`;

    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'X-DB-Connection-String': connectionString
            }
        });
    } catch (error) {
        console.error(error);
    }
```