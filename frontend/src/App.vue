<script setup>
import { ref, onMounted } from 'vue';

import axios from 'axios/dist/axios.min.js';

const radios = ref([]);
const form = ref({ name: '', address: '' });

// Загрузка раций
onMounted(async () => {
  const res = await axios.get('/api/radios');
  radios.value = res.data;
});

// Добавление рации
const addRadio = async () => {
  await axios.post('/api/radios', form.value);
  form.value = { name: '', address: '' };
};
</script>

<template>
  <div>
    <input v-model="form.name" placeholder="Название">
    <input v-model="form.address" placeholder="IP:Порт">
    <button @click="addRadio">Добавить</button>
    <div v-for="radio in radios" :key="radio.id">
      {{ radio.name }} - {{ radio.address }}
    </div>
  </div>
</template>