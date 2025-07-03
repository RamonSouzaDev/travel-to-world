<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow-lg">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <div class="flex-shrink-0 flex items-center">
              <div class="h-8 w-8 bg-blue-600 rounded-lg flex items-center justify-center mr-3">
                <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path>
                </svg>
              </div>
              <h1 class="text-xl font-bold text-gray-900">Travel Requests</h1>
            </div>
          </div>
          
          <div class="flex items-center space-x-4">
            <span class="text-gray-700">Olá, <span class="font-medium">{{ user?.name }}</span></span>
            <button 
              @click="logout" 
              class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
            >
              Sair
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <!-- Header Actions -->
      <div class="mb-8">
        <div class="sm:flex sm:items-center sm:justify-between">
          <div>
            <h2 class="text-2xl font-bold text-gray-900">Meus Pedidos de Viagem</h2>
            <p class="mt-1 text-sm text-gray-600">Gerencie suas solicitações de viagem corporativa</p>
          </div>
          <div class="mt-4 sm:mt-0">
            <button 
              @click="showCreateForm = !showCreateForm" 
              class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white transition-colors"
              :class="showCreateForm ? 'bg-gray-600 hover:bg-gray-700' : 'bg-blue-600 hover:bg-blue-700'"
            >
              <svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="showCreateForm ? 'M6 18L18 6M6 6l12 12' : 'M12 4v16m8-8H4'"></path>
              </svg>
              {{ showCreateForm ? 'Cancelar' : 'Novo Pedido' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Filtros Avançados -->
      <FilterComponent @filter="handleFilter" />

      <!-- Create Form -->
      <div v-if="showCreateForm" class="bg-white shadow rounded-lg mb-8">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">Novo Pedido de Viagem</h3>
        </div>
        
        <form @submit.prevent="createRequest" class="p-6">
          <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Nome do Solicitante
              </label>
              <input
                v-model="newRequest.requester_name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Nome completo do solicitante"
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Destino
              </label>
              <input
                v-model="newRequest.destination"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Cidade/País de destino"
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Data de Ida
              </label>
              <input
                v-model="newRequest.departure_date"
                type="date"
                required
                :min="today"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Data de Volta
              </label>
              <input
                v-model="newRequest.return_date"
                type="date"
                required
                :min="newRequest.departure_date || today"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
            </div>
          </div>
          
          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              @click="showCreateForm = false"
              class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Cancelar
            </button>
            <button
              type="submit"
              :disabled="loadingCreate"
              class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
            >
              {{ loadingCreate ? 'Criando...' : 'Criar Pedido' }}
            </button>
          </div>
        </form>
      </div>

      <!-- Requests List -->
      <div class="bg-white shadow rounded-lg">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-medium text-gray-900">
              Lista de Pedidos
              <span v-if="filteredCount > 0" class="text-sm text-gray-500 font-normal">
                ({{ filteredCount }} {{ filteredCount === 1 ? 'resultado' : 'resultados' }})
              </span>
            </h3>
            <button 
              @click="loadRequests" 
              class="text-blue-600 hover:text-blue-700 text-sm font-medium"
            >
              Atualizar
            </button>
          </div>
        </div>
        
        <!-- Loading State -->
        <div v-if="loadingList" class="p-8 text-center">
          <div class="inline-flex items-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Carregando pedidos...
          </div>
        </div>
        
        <!-- Empty State -->
        <div v-else-if="requests.length === 0" class="p-8 text-center">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">Nenhum pedido encontrado</h3>
          <p class="mt-1 text-sm text-gray-500">
            {{ appliedFilters ? 'Tente ajustar os filtros ou' : '' }} Comece criando seu primeiro pedido de viagem.
          </p>
        </div>
        
        <!-- Requests Grid -->
        <div v-else class="divide-y divide-gray-200">
          <div 
            v-for="request in requests" 
            :key="request.id"
            class="p-6 hover:bg-gray-50 transition-colors"
          >
            <div class="flex items-center justify-between">
              <div class="flex-1">
                <div class="flex items-center justify-between">
                  <div>
                    <h4 class="text-lg font-medium text-gray-900">{{ request.destination }}</h4>
                    <p class="text-sm text-gray-600">{{ request.requester_name }}</p>
                  </div>
                  <span 
                    :class="getStatusClass(request.status)"
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ getStatusText(request.status) }}
                  </span>
                </div>
                
                <div class="mt-2 text-sm text-gray-500">
                  <p>
                    <span class="font-medium">Período:</span> 
                    {{ formatDate(request.departure_date) }} até {{ formatDate(request.return_date) }}
                  </p>
                  <p>
                    <span class="font-medium">Criado em:</span> 
                    {{ formatDateTime(request.created_at) }}
                  </p>
                  <p v-if="request.created_by_id === user?.id" class="text-yellow-600 text-xs mt-1">
                    ⚠️ Você criou este pedido - outros usuários podem alterar o status
                  </p>
                </div>
                
                <div class="mt-4 flex space-x-2">
                  <!-- Só mostra botões se NÃO foi o usuário atual que criou -->
                  <template v-if="request.created_by_id !== user?.id">
                    <button
                      v-if="request.status === 'solicitado'"
                      @click="updateStatus(request.id, 'aprovado')"
                      class="inline-flex items-center px-3 py-1 border border-transparent text-xs font-medium rounded text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                    >
                      Aprovar
                    </button>
                    
                    <button
                      v-if="request.status === 'solicitado'"
                      @click="updateStatus(request.id, 'cancelado')"
                      class="inline-flex items-center px-3 py-1 border border-transparent text-xs font-medium rounded text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                    >
                      Rejeitar
                    </button>
                  </template>
                  
                  <button
                    v-if="request.status === 'aprovado'"
                    @click="cancelRequest(request.id)"
                    class="inline-flex items-center px-3 py-1 border border-transparent text-xs font-medium rounded text-white bg-orange-600 hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-500"
                  >
                    Cancelar Viagem
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import FilterComponent from './FilterComponent.vue'

const router = useRouter()

// Estado
const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))
const requests = ref([])
const showCreateForm = ref(false)
const loadingList = ref(false)
const loadingCreate = ref(false)
const appliedFilters = ref(false)

const newRequest = ref({
  requester_name: '',
  destination: '',
  departure_date: '',
  return_date: ''
})

// Computed
const today = computed(() => {
  return new Date().toISOString().split('T')[0]
})

const filteredCount = computed(() => {
  return requests.value.length
})

// Configurar axios
const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('token')}`
  }
})

// Interceptor para lidar com erros de autenticação
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      localStorage.clear()
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

// Funções
const loadRequests = async (filters = {}) => {
  loadingList.value = true
  try {
    const params = new URLSearchParams(filters).toString()
    const url = '/travel-requests' + (params ? `?${params}` : '')
    const response = await api.get(url)
    requests.value = response.data || []
    appliedFilters.value = Object.keys(filters).length > 0
  } catch (error) {
    console.error('Erro ao carregar pedidos:', error)
    alert('Erro ao carregar pedidos: ' + (error.response?.data?.error || 'Erro desconhecido'))
  } finally {
    loadingList.value = false
  }
}

const handleFilter = (filters) => {
  loadRequests(filters)
}

const createRequest = async () => {
  loadingCreate.value = true
  try {
    await api.post('/travel-requests', newRequest.value)
    showCreateForm.value = false
    newRequest.value = {
      requester_name: '',
      destination: '',
      departure_date: '',
      return_date: ''
    }
    await loadRequests()
    alert('Pedido criado com sucesso!')
  } catch (error) {
    console.error('Erro ao criar pedido:', error)
    alert('Erro ao criar pedido: ' + (error.response?.data?.error || 'Erro desconhecido'))
  } finally {
    loadingCreate.value = false
  }
}

const updateStatus = async (id, status) => {
  try {
    await api.put(`/travel-requests/${id}/status`, { status })
    await loadRequests()
    alert(`Status atualizado para: ${getStatusText(status)}`)
  } catch (error) {
    console.error('Erro ao atualizar status:', error)
    const errorMsg = error.response?.data?.error || 'Erro desconhecido'
    alert('Erro ao atualizar status: ' + errorMsg)
  }
}

const cancelRequest = async (id) => {
  if (confirm('Tem certeza que deseja cancelar este pedido de viagem?')) {
    try {
      await api.delete(`/travel-requests/${id}`)
      await loadRequests()
      alert('Pedido cancelado com sucesso!')
    } catch (error) {
      console.error('Erro ao cancelar pedido:', error)
      alert('Erro ao cancelar pedido: ' + (error.response?.data?.error || 'Erro desconhecido'))
    }
  }
}

const logout = () => {
  if (confirm('Tem certeza que deseja sair?')) {
    localStorage.clear()
    router.push('/login')
  }
}

const getStatusClass = (status) => {
  const classes = {
    'solicitado': 'bg-yellow-100 text-yellow-800',
    'aprovado': 'bg-green-100 text-green-800',
    'cancelado': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  const texts = {
    'solicitado': 'Pendente',
    'aprovado': 'Aprovado',
    'cancelado': 'Cancelado'
  }
  return texts[status] || status
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('pt-BR')
}

const formatDateTime = (dateString) => {
  return new Date(dateString).toLocaleString('pt-BR')
}

// Lifecycle
onMounted(() => {
  if (!localStorage.getItem('token')) {
    router.push('/login')
    return
  }
  loadRequests()
})
</script>
