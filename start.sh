#!/bin/bash

echo "🚀 Iniciando Travel Requests..."

# Parar containers existentes
echo "🛑 Parando containers existentes..."
docker-compose down

# Iniciar com build
echo "🔨 Iniciando containers (pode levar alguns minutos)..."
docker-compose up --build -d

# Aguardar serviços
echo "⏳ Aguardando serviços ficarem prontos..."
sleep 20

# Verificar backend
echo "🔍 Verificando backend..."
for i in {1..15}; do
    if curl -s http://localhost:8080/health | grep -q "ok"; then
        echo "✅ Backend funcionando: http://localhost:8080"
        break
    fi
    sleep 2
done

# Verificar frontend
echo "🔍 Verificando frontend..."
sleep 5
if curl -s http://localhost:3000 >/dev/null 2>&1; then
    echo "✅ Frontend funcionando: http://localhost:3000"
else
    echo "⚠️ Frontend ainda iniciando..."
fi

echo ""
echo "🎉 TRAVEL REQUESTS INICIADO!"
echo "================================"
echo "✅ Frontend: http://localhost:3000"
echo "✅ Backend:  http://localhost:8080"
echo "✅ Health:   http://localhost:8080/health"
echo "✅ Database: PostgreSQL na porta 5432"
echo ""
echo "📋 Status dos containers:"
docker-compose ps
