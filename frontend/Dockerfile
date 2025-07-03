FROM node:18-alpine AS builder

WORKDIR /app

# Copiar package.json
COPY package.json ./

# Instalar dependências
RUN npm install

# Copiar código fonte
COPY . .

# Build para produção
RUN npm run build

# Estágio de produção com nginx
FROM nginx:alpine

# Remover configuração padrão
RUN rm /etc/nginx/conf.d/default.conf

# Copiar nossa configuração
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copiar build
COPY --from=builder /app/dist /usr/share/nginx/html

# Expor porta
EXPOSE 80

# Iniciar nginx
CMD ["nginx", "-g", "daemon off;"]
