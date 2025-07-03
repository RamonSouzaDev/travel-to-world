<template>
  <div class="bg-gray-50 p-4 rounded-lg mb-6">
    <h4 class="font-medium mb-3 text-gray-900">Filtros de Busca</h4>
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <!-- Status -->
      <div>
        <label class="block text-xs font-medium text-gray-700 mb-1">Status</label>
        <select v-model="filters.status" class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
          <option value="">Todos</option>
          <option value="solicitado">Solicitado</option>
          <option value="aprovado">Aprovado</option>
          <option value="cancelado">Cancelado</option>
        </select>
      </div>
      
      <!-- Destino -->
      <div>
        <label class="block text-xs font-medium text-gray-700 mb-1">Destino</label>
        <input
          v-model="filters.destination"
          type="text"
          placeholder="Ex: São Paulo"
          class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      
      <!-- Data de Ida (Início) -->
      <div>
        <label class="block text-xs font-medium text-gray-700 mb-1">Ida Após</label>
        <input
          v-model="filters.start_date"
          type="date"
          class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      
      <!-- Data de Volta (Fim) -->
      <div>
        <label class="block text-xs font-medium text-gray-700 mb-1">Volta Antes</label>
        <input
          v-model="filters.end_date"
          type="date"
          class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      
      <!-- Criado Após -->
      <div>
        <label class="block text-xs font-medium text-gray-700 mb-1">Criado Após</label>
        <input
          v-model="filters.created_after"
          type="date"
          class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      
      <!-- Criado Antes -->
      <div>
        <label class="block text-xs font-medium text-gray-700 mb-1">Criado Antes</label>
        <input
          v-model="filters.created_before"
          type="date"
          class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
    </div>
    
    <div class="flex gap-3 mt-4">
      <button 
        @click="applyFilters" 
        class="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
      >
        Aplicar Filtros
      </button>
      <button 
        @click="clearFilters" 
        class="px-4 py-2 bg-gray-200 text-gray-800 text-sm font-medium rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
      >
        Limpar
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['filter'])

const filters = ref({
  status: '',
  destination: '',
  start_date: '',
  end_date: '',
  created_after: '',
  created_before: ''
})

const applyFilters = () => {
  // Remove campos vazios
  const activeFilters = {}
  Object.keys(filters.value).forEach(key => {
    if (filters.value[key]) {
      activeFilters[key] = filters.value[key]
    }
  })
  
  emit('filter', activeFilters)
}

const clearFilters = () => {
  filters.value = {
    status: '',
    destination: '',
    start_date: '',
    end_date: '',
    created_after: '',
    created_before: ''
  }
  emit('filter', {})
}
</script>
