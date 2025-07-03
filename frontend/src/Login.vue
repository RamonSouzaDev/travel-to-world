<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
    <div class="max-w-md w-full mx-4">
      <!-- Header -->
      <div class="text-center mb-8">
        <div class="mx-auto h-12 w-12 bg-blue-600 rounded-lg flex items-center justify-center mb-4">
          <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path>
          </svg>
        </div>
        <h2 class="text-3xl font-bold text-gray-900">Travel Requests</h2>
        <p class="mt-2 text-gray-600">Sistema de Gerenciamento de Viagens Corporativas</p>
      </div>
      
      <!-- Card principal -->
      <div class="bg-white rounded-lg shadow-xl p-8">
        <!-- Mensagens -->
        <div v-if="error" class="mb-4 p-4 bg-red-50 border-l-4 border-red-400 text-red-700">
          <div class="flex">
            <div class="ml-3">
              <p class="text-sm">{{ error }}</p>
            </div>
          </div>
        </div>
        
        <div v-if="success" class="mb-4 p-4 bg-green-50 border-l-4 border-green-400 text-green-700">
          <div class="flex">
            <div class="ml-3">
              <p class="text-sm">{{ success }}</p>
            </div>
          </div>
        </div>
        
        <!-- Toggle Login/Register -->
        <div class="mb-6">
          <div class="flex bg-gray-100 rounded-lg p-1">
            <button 
              @click="isLogin = true"
              :class="isLogin ? 'bg-white shadow-sm' : 'text-gray-500'"
              class="flex-1 py-2 px-4 rounded-md text-sm font-medium transition-all"
            >
              Login
            </button>
            <button 
              @click="isLogin = false"
              :class="!isLogin ? 'bg-white shadow-sm' : 'text-gray-500'"
              class="flex-1 py-2 px-4 rounded-md text-sm font-medium transition-all"
            >
              Registrar
            </button>
          </div>
        </div>
        
        <!-- Formulário -->
        <form @submit.prevent="isLogin ? login() : register()" class="space-y-6">
          <!-- Nome (só no registro) -->
          <div v-if="!isLogin">
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Nome Completo
            </label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Seu nome completo"
            >
          </div>
          
          <!-- Email -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              v-model="form.email"
              type="email"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="seu@email.com"
            >
          </div>
          
          <!-- Senha -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Senha
            </label>
            <input
              v-model="form.password"
              type="password"
              required
              :minlength="isLogin ? 1 : 6"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :placeholder="isLogin ? 'Sua senha' : 'Mínimo 6 caracteres'"
            >
          </div>
          
          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <svg v-if="loading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ loading ? 'Processando...' : (isLogin ? 'Entrar' : 'Criar Conta') }}
          </button>
        </form>
        
        <!-- Demo info -->
        <div class="mt-6 p-4 bg-blue-50 rounded-md">
          <h4 class="text-sm font-medium text-blue-800 mb-2">💡 Para testar rapidamente:</h4>
          <p class="text-xs text-blue-700">
            1. Clique em "Registrar" e crie um usuário<br>
            2. Faça login com as credenciais criadas<br>
            3. Crie e gerencie pedidos de viagem
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

// Estado reativo
const isLogin = ref(true)
const loading = ref(false)
const error = ref('')
const success = ref('')

const form = ref({
  name: '',
  email: '',
  password: ''
})

// Configurar axios
const api = axios.create({
  baseURL: 'http://localhost:8080/api'
})

// Funções
const clearMessages = () => {
  error.value = ''
  success.value = ''
}

const resetForm = () => {
  form.value = {
    name: '',
    email: '',
    password: ''
  }
}

const login = async () => {
  loading.value = true
  clearMessages()
  
  try {
    const response = await api.post('/auth/login', {
      email: form.value.email,
      password: form.value.password
    })
    
    // Salvar dados no localStorage
    localStorage.setItem('token', response.data.token)
    localStorage.setItem('user', JSON.stringify(response.data.user))
    
    success.value = 'Login realizado com sucesso!'
    
    // Redirect para dashboard
    setTimeout(() => {
      router.push('/dashboard')
    }, 1000)
    
  } catch (err) {
    error.value = err.response?.data?.error || 'Erro ao fazer login. Verifique suas credenciais.'
  } finally {
    loading.value = false
  }
}

const register = async () => {
  loading.value = true
  clearMessages()
  
  try {
    await api.post('/auth/register', form.value)
    
    success.value = 'Usuário criado com sucesso! Agora você pode fazer login.'
    resetForm()
    isLogin.value = true
    
  } catch (err) {
    error.value = err.response?.data?.error || 'Erro ao criar usuário. Tente novamente.'
  } finally {
    loading.value = false
  }
}
</script>
