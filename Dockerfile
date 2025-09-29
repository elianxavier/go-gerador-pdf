# ESTÁGIO 1: BUILDER (Compilação)
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o api_pdf_softprev main.go


# --------------------------------------------------------------------------------

# ESTÁGIO 2: RUNNER (Foco na Estabilidade e Instalação)
# Trocamos para debian:bookworm-slim. É o mínimo necessário para rodar wkhtmltopdf de forma estável.
FROM debian:bookworm-slim

# Instalação Robusta do wkhtmltopdf:
# Usamos o gerenciador de pacotes 'apt-get' do Debian, que garante a instalação.
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    wkhtmltopdf \
    # Fontes essenciais para o wkhtmltopdf funcionar corretamente
    xfonts-base \
    # Limpa o cache imediatamente para reduzir o tamanho da imagem e a superfície de ataque
    && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copia APENAS o binário compilado
COPY --from=builder /app/api_pdf_softprev .

EXPOSE 5555

CMD ["./api_pdf_softprev"]
